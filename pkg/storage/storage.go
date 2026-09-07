package storage

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
)

var ErrDeckNotFound = errors.New("deck not found")
var ErrEmptyDeck = errors.New("deck has no more cards")

type Deck struct {
	DeckID   *uuid.UUID
	Shuffled bool
	Cards    []cards.Card
}

func (d *Deck) Remaining() int {
	return len(d.Cards)
}

type Storage interface {
	Create(ctx context.Context, cards []cards.Card, shuffled bool) (*Deck, error)
	Get(ctx context.Context, deckID *uuid.UUID) (*Deck, error)
	Draw(ctx context.Context, deckID *uuid.UUID, count int) ([]cards.Card, error)
	Delete(ctx context.Context, deckID *uuid.UUID) error
	Shuffle(ctx context.Context, deckID *uuid.UUID) (*Deck, error)
}

// Shuffles the given list of cards in place, using a Fisher-Yates shuffle.
// This is intentionally kept separate from cards.CardGenerator.Shuffle so that
// the storage package doesn't need to depend on a generator instance.
func shuffleCards(list []cards.Card) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(list) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}
}

func NewStorage() (Storage, error) {
	switch os.Getenv("DECK_STORAGE_TYPE") {
	case "redis":
		return NewRedisStorage(nil)
	default:
		log.Printf("No DECK_STORAGE_TYPE set, using in-memory default\n")
		return NewInMemoryStorage(), nil
	}
}
