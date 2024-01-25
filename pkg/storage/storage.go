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

type Storage interface {
	Create(ctx context.Context, cards []cards.Card, shuffled bool) (*Deck, error)
	Get(ctx context.Context, deckID *uuid.UUID) (*Deck, error)
	Draw(ctx context.Context, deckID *uuid.UUID, count int) ([]cards.Card, error)
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
