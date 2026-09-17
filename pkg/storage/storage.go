package storage

import (
	"context"
	"errors"
	"log"
	"os"

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

// ShuffleFunc reorders the given list of cards, returning the reordered list.
// It exists so the storage layer does not need to own a card generator/RNG itself;
// the caller (the server, which already holds a generator) provides the shuffling logic.
type ShuffleFunc func(list []cards.Card) []cards.Card

type Storage interface {
	Create(ctx context.Context, cards []cards.Card, shuffled bool) (*Deck, error)
	Get(ctx context.Context, deckID *uuid.UUID) (*Deck, error)
	Draw(ctx context.Context, deckID *uuid.UUID, count int) ([]cards.Card, error)
	Delete(ctx context.Context, deckID *uuid.UUID) error
	Reshuffle(ctx context.Context, deckID *uuid.UUID, shuffle ShuffleFunc) (*Deck, error)
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
