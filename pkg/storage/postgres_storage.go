package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lucaspin/decks-api/pkg/cards"
)

// An implementation of the Storage interface backed by PostgreSQL.
//
// A deck is stored across two tables:
//   - 'decks' holds the deck metadata (its ID and whether it is shuffled).
//   - 'deck_cards' holds one row per card, with a 'position' column that
//     preserves insertion order (0-indexed), so that drawing cards is
//     equivalent to taking the rows with the lowest positions.
//
// Both tables are created idempotently (CREATE TABLE IF NOT EXISTS) the
// first time this storage is instantiated, so no separate migration step
// is needed to run the server or the tests.
//
// Unlike RedisStorage, draws done through this implementation are atomic:
// a single SQL statement (a CTE combined with SELECT ... FOR UPDATE) locks
// and removes the drawn rows in one round-trip, wrapped in a transaction
// that also checks for the deck's existence. This means concurrent draws
// for the same deck never hand out the same card twice.
type PostgresStorage struct {
	Pool *pgxpool.Pool
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func NewPostgresConfigFromEnvironment() (*PostgresConfig, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		return nil, fmt.Errorf("no POSTGRES_HOST set")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}

	database := os.Getenv("POSTGRES_DB")
	if database == "" {
		database = "decks"
	}

	sslMode := os.Getenv("POSTGRES_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	return &PostgresConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: database,
		SSLMode:  sslMode,
	}, nil
}

func NewPostgresStorage(config *PostgresConfig) (Storage, error) {
	// If no config is passed,
	// we try to create it from environment variables.
	if config == nil {
		c, err := NewPostgresConfigFromEnvironment()
		if err != nil {
			return nil, err
		}

		config = c
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.Database, config.SSLMode,
	)

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	// Make sure we have a valid connection before proceeding.
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	if err := ensureSchema(ctx, pool); err != nil {
		return nil, err
	}

	log.Printf("Successfully connected to Postgres")
	return &PostgresStorage{Pool: pool}, nil
}

func ensureSchema(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS decks (
			deck_id    UUID PRIMARY KEY,
			shuffled   BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS deck_cards (
			deck_id  UUID NOT NULL REFERENCES decks(deck_id) ON DELETE CASCADE,
			position INT NOT NULL,
			code     VARCHAR(3) NOT NULL,
			PRIMARY KEY (deck_id, position)
		);
	`)

	return err
}

func (s *PostgresStorage) Create(ctx context.Context, list []cards.Card, shuffled bool) (*Deck, error) {
	ID := uuid.New()

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `INSERT INTO decks (deck_id, shuffled) VALUES ($1, $2)`, ID, shuffled)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		codes := cards.CardListToCodes(list)
		_, err = tx.Exec(ctx, `
			INSERT INTO deck_cards (deck_id, position, code)
			SELECT $1, position - 1, code
			FROM unnest($2::text[]) WITH ORDINALITY AS c(code, position)
		`, ID, codes)

		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &Deck{
		DeckID:   &ID,
		Shuffled: shuffled,
		Cards:    list,
	}, nil
}

func (s *PostgresStorage) Get(ctx context.Context, deckID *uuid.UUID) (*Deck, error) {
	var shuffled bool
	err := s.Pool.QueryRow(ctx, `SELECT shuffled FROM decks WHERE deck_id = $1`, deckID).Scan(&shuffled)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrDeckNotFound
	}

	if err != nil {
		return nil, err
	}

	rows, err := s.Pool.Query(ctx, `SELECT code FROM deck_cards WHERE deck_id = $1 ORDER BY position ASC`, deckID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	codes := []string{}
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}

		codes = append(codes, code)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// We know this is valid because we validate it before inserting.
	cardList, _ := cards.CodesToCardList(codes)

	return &Deck{
		DeckID:   deckID,
		Shuffled: shuffled,
		Cards:    cardList,
	}, nil
}

func (s *PostgresStorage) Draw(ctx context.Context, deckID *uuid.UUID, count int) ([]cards.Card, error) {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	// Lock the deck row first, so we can tell apart a deck that does not
	// exist from a deck that simply has no cards left.
	var exists int
	err = tx.QueryRow(ctx, `SELECT 1 FROM decks WHERE deck_id = $1 FOR UPDATE`, deckID).Scan(&exists)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrDeckNotFound
	}

	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `
		WITH to_draw AS (
			SELECT deck_id, position, code
			FROM deck_cards
			WHERE deck_id = $1
			ORDER BY position ASC
			LIMIT $2
			FOR UPDATE
		)
		DELETE FROM deck_cards
		USING to_draw
		WHERE deck_cards.deck_id = to_draw.deck_id
		  AND deck_cards.position = to_draw.position
		RETURNING deck_cards.code, deck_cards.position
	`, deckID, count)

	if err != nil {
		return nil, err
	}

	type drawnCard struct {
		code     string
		position int
	}

	drawn := []drawnCard{}
	for rows.Next() {
		var c drawnCard
		if err := rows.Scan(&c.code, &c.position); err != nil {
			rows.Close()
			return nil, err
		}

		drawn = append(drawn, c)
	}

	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(drawn) == 0 {
		return nil, ErrEmptyDeck
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// Rows come back from the driver in no guaranteed order,
	// so we sort them by position before converting them back to cards.
	sort.Slice(drawn, func(i, j int) bool { return drawn[i].position < drawn[j].position })

	codes := make([]string, len(drawn))
	for i, c := range drawn {
		codes[i] = c.code
	}

	// We know this is valid because we validate it before inserting.
	cardList, _ := cards.CodesToCardList(codes)
	return cardList, nil
}
