package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/lucaspin/decks-api/pkg/storage"
)

// fakeStorage is a storage.Storage implementation that lets tests force
// arbitrary errors out of each method. It is used to exercise the 500
// (and other error) branches that the real InMemoryStorage can never
// produce on its own.
type fakeStorage struct {
	createErr error
	getErr    error
	drawErr   error

	deck  *storage.Deck
	cards []cards.Card
}

func (s *fakeStorage) Create(ctx context.Context, list []cards.Card, shuffled bool) (*storage.Deck, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}

	ID := uuid.New()
	return &storage.Deck{DeckID: &ID, Shuffled: shuffled, Cards: list}, nil
}

func (s *fakeStorage) Get(ctx context.Context, deckID *uuid.UUID) (*storage.Deck, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}

	if s.deck != nil {
		return s.deck, nil
	}

	return &storage.Deck{DeckID: deckID}, nil
}

func (s *fakeStorage) Draw(ctx context.Context, deckID *uuid.UUID, count int) ([]cards.Card, error) {
	if s.drawErr != nil {
		return nil, s.drawErr
	}

	return s.cards, nil
}
