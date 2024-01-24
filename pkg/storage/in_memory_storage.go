package storage

import (
	"context"

	"github.com/google/uuid"
	cards "github.com/lucaspin/decks-api/pkg/cards"
)

type InMemoryStorage struct {
	decks map[string]Deck
}

func NewInMemoryStorage() Storage {
	return &InMemoryStorage{decks: map[string]Deck{}}
}

func (s *InMemoryStorage) Create(ctx context.Context) (*Deck, error) {
	cards := cards.New()
	ID := uuid.New()

	deck := Deck{
		DeckID:   &ID,
		Shuffled: false,
		Cards:    cards,
	}

	s.decks[deck.DeckID.String()] = deck
	return &deck, nil
}

func (s *InMemoryStorage) Get(ctx context.Context, deckID *uuid.UUID) (*Deck, error) {
	deck, ok := s.decks[deckID.String()]
	if !ok {
		return nil, ErrDeckNotFound
	}

	return &deck, nil
}
