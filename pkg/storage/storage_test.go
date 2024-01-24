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

	t.Run("drawing from empty deck -> error", func(t *testing.T) {
		deck, err := storage.Create(context.Background(), []cards.Card{}, false)
		require.NoError(t, err)

		_, err = storage.Draw(context.Background(), deck.DeckID, 1)
		require.ErrorIs(t, err, ErrEmptyDeck)
	})

	t.Run("drawing more cards than deck has -> error", func(t *testing.T) {
		initial := []cards.Card{
			{Suit: cards.CardSuitClubs, Rank: cards.CardRank(3)},
			{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(8)},
		}

		deck, err := storage.Create(context.Background(), initial, false)
		require.NoError(t, err)

		_, err = storage.Draw(context.Background(), deck.DeckID, 3)
		require.ErrorIs(t, err, ErrNotEnoughCardsInDeck)
	})

	t.Run("drawing removes cards from deck", func(t *testing.T) {
		initial := []cards.Card{
			{Suit: cards.CardSuitClubs, Rank: cards.CardRank(3)},
			{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(8)},
		}

		deck, err := storage.Create(context.Background(), initial, false)
		require.NoError(t, err)

		drawn, err := storage.Draw(context.Background(), deck.DeckID, 1)
		require.NoError(t, err)
		require.Equal(t, []cards.Card{{Suit: cards.CardSuitClubs, Rank: cards.CardRank(3)}}, drawn)

		deck, err = storage.Get(context.Background(), deck.DeckID)
		require.NoError(t, err)
		require.Len(t, deck.Cards, 1)
	})
}
