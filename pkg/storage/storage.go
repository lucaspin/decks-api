package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
)

var ErrDeckNotFound = errors.New("deck not found")

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
}

func NewStorage() Storage {
	return NewInMemoryStorage()
}
