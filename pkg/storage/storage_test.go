package storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/stretchr/testify/require"
)

func Test__InMemoryStorage(t *testing.T) {
	storage := NewInMemoryStorage()

	t.Run("get with deck that does not exist -> ErrDeckNotFound error", func(t *testing.T) {
		ID := uuid.New()
		_, err := storage.Get(context.Background(), &ID)
		require.ErrorIs(t, err, ErrDeckNotFound)
	})

	t.Run("get with existing deck -> returns deck", func(t *testing.T) {
		cards := []cards.Card{
			{Suit: cards.CardSuitClubs, Rank: cards.CardRank(3)},
			{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(8)},
		}

		d1, err := storage.Create(context.Background(), cards, false)
		require.NoError(t, err)

		d2, err := storage.Get(context.Background(), d1.DeckID)
		require.NoError(t, err)
		require.Equal(t, d1, d2)
	})
}
