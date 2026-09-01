package storage

import (
	"context"
	"testing"

	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/stretchr/testify/require"
)

// These tests cover the in-memory storage implementation directly,
// so they run without needing a live Redis server, unlike the shared
// suite in storage_test.go (which iterates over all implementations,
// including Redis).
func Test__InMemoryStorage__CreateGetDraw(t *testing.T) {
	initial := []cards.Card{
		{Suit: cards.CardSuitSpades, Rank: cards.CardRank(1)},
		{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(13)},
	}

	storage := NewInMemoryStorage()

	deck, err := storage.Create(context.Background(), initial, false)
	require.NoError(t, err)
	require.NotNil(t, deck.DeckID)
	require.False(t, deck.Shuffled)
	require.Equal(t, initial, deck.Cards)

	fetched, err := storage.Get(context.Background(), deck.DeckID)
	require.NoError(t, err)
	require.Equal(t, deck, fetched)

	drawn, err := storage.Draw(context.Background(), deck.DeckID, 1)
	require.NoError(t, err)
	require.Equal(t, []cards.Card{{Suit: cards.CardSuitSpades, Rank: cards.CardRank(1)}}, drawn)

	fetched, err = storage.Get(context.Background(), deck.DeckID)
	require.NoError(t, err)
	require.Equal(t, 1, fetched.Remaining())
}

func Test__InMemoryStorage__DrawWithCountZero(t *testing.T) {
	initial := []cards.Card{
		{Suit: cards.CardSuitSpades, Rank: cards.CardRank(1)},
		{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(13)},
	}

	storage := NewInMemoryStorage()
	deck, err := storage.Create(context.Background(), initial, false)
	require.NoError(t, err)

	drawn, err := storage.Draw(context.Background(), deck.DeckID, 0)
	require.NoError(t, err)
	require.Empty(t, drawn)

	// the deck is left untouched
	fetched, err := storage.Get(context.Background(), deck.DeckID)
	require.NoError(t, err)
	require.Equal(t, 2, fetched.Remaining())
}

func Test__InMemoryStorage__GetAfterDrawingAllCards(t *testing.T) {
	initial := []cards.Card{
		{Suit: cards.CardSuitSpades, Rank: cards.CardRank(1)},
		{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(13)},
	}

	storage := NewInMemoryStorage()
	deck, err := storage.Create(context.Background(), initial, false)
	require.NoError(t, err)

	_, err = storage.Draw(context.Background(), deck.DeckID, 2)
	require.NoError(t, err)

	fetched, err := storage.Get(context.Background(), deck.DeckID)
	require.NoError(t, err)
	require.Equal(t, 0, fetched.Remaining())
	require.Empty(t, fetched.Cards)
}

func Test__InMemoryStorage__RemainingReflectsDraws(t *testing.T) {
	initial := []cards.Card{
		{Suit: cards.CardSuitSpades, Rank: cards.CardRank(1)},
		{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(13)},
		{Suit: cards.CardSuitClubs, Rank: cards.CardRank(7)},
	}

	storage := NewInMemoryStorage()
	deck, err := storage.Create(context.Background(), initial, false)
	require.NoError(t, err)
	require.Equal(t, 3, deck.Remaining())

	_, err = storage.Draw(context.Background(), deck.DeckID, 1)
	require.NoError(t, err)

	fetched, err := storage.Get(context.Background(), deck.DeckID)
	require.NoError(t, err)
	require.Equal(t, 2, fetched.Remaining())
}
