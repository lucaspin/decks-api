package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
)

type InMemoryStorage struct {
	decks map[string]Deck
}

func NewInMemoryStorage() Storage {
	return &InMemoryStorage{decks: map[string]Deck{}}
}

func (s *InMemoryStorage) Create(ctx context.Context, list []cards.Card, shuffled bool) (*Deck, error) {
	ID := uuid.New()
	deck := Deck{
		DeckID:   &ID,
		Shuffled: shuffled,
		Cards:    list,
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

func (s *InMemoryStorage) Draw(ctx context.Context, deckID *uuid.UUID, count int) ([]cards.Card, error) {
	deck, ok := s.decks[deckID.String()]
	if !ok {
		return nil, ErrDeckNotFound
	}

	if len(deck.Cards) == 0 {
		return nil, ErrEmptyDeck
	}

	if len(deck.Cards) < count {
		return nil, ErrNotEnoughCardsInDeck
	}

	// gather cards
	cards := make([]cards.Card, count)
	for i := 0; i < count; i++ {
		cards[i] = deck.Cards[i]
	}

	// remove cards from deck
	s.decks[deckID.String()] = Deck{
		DeckID:   deckID,
		Shuffled: deck.Shuffled,
		Cards:    deck.Cards[count:],
	}

	return cards, nil
}
