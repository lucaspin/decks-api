package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
)

// A naive implementation of a deck storage using Redis.
// The operations done here are not atomic, so this implementation is not safe
// to be used if the decks API are supposed to serve multiple requests for the same deck at the same time.
//
// To make if safe, we'd need to use distributed locks or lua scripts.
// See: https://lucaspin.github.io/redis/databases/2021/07/21/atomicity-in-redis-operations.html.
//
// In Redis, a deck is composed of two keys:
// 'decks:{deckID}:cards' - a Redis list that holds the current list of card codes for the deck.
// 'decks:{deckID}:shuffled' - a Redis value indicating if the deck is shuffled or not.
//
// This makes it easy to draw cards from the deck:
// just use the LPOP operation on the cards key.

type RedisStorage struct {
	Client *redis.Client
}

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

func NewRedisStorage(config *RedisConfig) (Storage, error) {
	// If no config is passed,
	// we try to create it from environment variables.
	if config == nil {
		c, err := NewRedisConfigFromEnvironment()
		if err != nil {
			return nil, err
		}

		config = c
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Username: config.Username,
		Password: config.Password,
		DB:       0,
	})

	// Make sure we have a valid connection before proceeding.
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	log.Printf("Successfully connected to Redis")
	return &RedisStorage{Client: rdb}, nil
}

func NewRedisConfigFromEnvironment() (*RedisConfig, error) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return nil, fmt.Errorf("no REDIS_HOST set")
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	return &RedisConfig{
		Host:     host,
		Port:     port,
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}, nil
}

func (s *RedisStorage) Create(ctx context.Context, list []cards.Card, shuffled bool) (*Deck, error) {
	ID := uuid.New()
	deck := Deck{
		DeckID:   &ID,
		Shuffled: shuffled,
		Cards:    list,
	}

	cardsKey := keyForAttribute(&ID, "cards")
	shuffledKey := keyForAttribute(&ID, "shuffled")

	// Push the cards
	_, err := s.Client.RPush(ctx, cardsKey, cards.CardListToCodes(list)).Result()
	if err != nil {
		return nil, err
	}

	// Add shuffled attribute.
	// If this fails, we make a small effort
	// to delete the key we added previously for the cards.
	_, err = s.Client.Set(ctx, shuffledKey, shuffled, 0).Result()
	if err != nil {
		if _, err := s.Client.Del(ctx, cardsKey).Result(); err != nil {
			log.Printf("Error rolling back key: %v", err)
		}

		return nil, err
	}

	return &deck, nil
}

func (s *RedisStorage) Get(ctx context.Context, deckID *uuid.UUID) (*Deck, error) {
	shuffled, err := s.getShuffledAttribute(ctx, deckID)
	if errors.Is(err, ErrDeckNotFound) {
		return nil, ErrDeckNotFound
	}

	// Unknown error
	if err != nil {
		return nil, err
	}

	// We know the deck exists, so let's grab all cards for it.
	cardsKey := keyForAttribute(deckID, "cards")
	list, err := s.Client.LRange(ctx, cardsKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// We know this is valid because we validate it before inserting.
	cardList, _ := cards.CodesToCardList(list)

	return &Deck{
		DeckID:   deckID,
		Shuffled: shuffled,
		Cards:    cardList,
	}, nil
}

func (s *RedisStorage) Draw(ctx context.Context, deckID *uuid.UUID, count int) ([]cards.Card, error) {
	// We don't really need the shuffled attribute here,
	// but LPOP returns redis.Nil when the deck does not exist,
	// and when the deck is empty, so we use this key to check for existence.
	_, err := s.getShuffledAttribute(ctx, deckID)
	if errors.Is(err, ErrDeckNotFound) {
		return nil, err
	}

	// Unknown error
	if err != nil {
		return nil, err
	}

	cardsKey := keyForAttribute(deckID, "cards")
	list, err := s.Client.LPopCount(ctx, cardsKey, count).Result()

	// If we receive a redis.Nil, we know that means an empty deck,
	// because we checked for deck existence above.
	if errors.Is(err, redis.Nil) {
		return nil, ErrEmptyDeck
	}

	// Unknown error
	if err != nil {
		return nil, err
	}

	// We know this is valid because we validate it before inserting.
	cardList, _ := cards.CodesToCardList(list)
	return cardList, nil
}

func keyForAttribute(deckID *uuid.UUID, attrName string) string {
	return fmt.Sprintf("decks:%s:%s", deckID.String(), attrName)
}

func (s *RedisStorage) getShuffledAttribute(ctx context.Context, deckID *uuid.UUID) (bool, error) {
	shuffledKey := keyForAttribute(deckID, "shuffled")
	_, err := s.Client.Get(ctx, shuffledKey).Result()

	// When a key does not exist, Redis gives us a Nil reply
	if errors.Is(err, redis.Nil) {
		return false, ErrDeckNotFound
	}

	// Some other unknown error.
	if err != nil {
		return false, err
	}

	return shuffledKey == "true", nil
}
