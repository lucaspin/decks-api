package storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test__InMemoryStorage(t *testing.T) {
	storage := NewInMemoryStorage()

	t.Run("create returns new deck", func(t *testing.T) {
		deck, err := storage.Create(context.Background())
		require.NoError(t, err)
		require.NotNil(t, deck)
		require.NotNil(t, deck.DeckID)
		require.False(t, deck.Shuffled)
		require.Equal(t, 52, deck.Remaining())
		require.Len(t, deck.Cards, 52)
	})

	t.Run("get returns not found", func(t *testing.T) {
		ID := uuid.New()
		_, err := storage.Get(context.Background(), &ID)
		require.ErrorIs(t, err, ErrDeckNotFound)
	})

	t.Run("get returns not found", func(t *testing.T) {
		d1, err := storage.Create(context.Background())
		require.NoError(t, err)

		d2, err := storage.Get(context.Background(), d1.DeckID)
		require.NoError(t, err)
		require.Equal(t, d1, d2)
	})
}
