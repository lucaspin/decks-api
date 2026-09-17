package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/stretchr/testify/require"
)

func Test__StorageTest(t *testing.T) {
	runTestForAllImplementations(t, func(storageName string, storage Storage) {
		t.Run(fmt.Sprintf("%s - get with deck that does not exist -> ErrDeckNotFound error", storageName), func(t *testing.T) {
			ID := uuid.New()
			_, err := storage.Get(context.Background(), &ID)
			require.ErrorIs(t, err, ErrDeckNotFound)
		})

		t.Run(fmt.Sprintf("%s - get with existing deck -> returns deck", storageName), func(t *testing.T) {
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

		t.Run(fmt.Sprintf("%s - drawing from empty deck -> error", storageName), func(t *testing.T) {
			initial := []cards.Card{{Suit: cards.CardSuitClubs, Rank: cards.CardRank(3)}}
			deck, err := storage.Create(context.Background(), initial, false)
			require.NoError(t, err)

			// draw all the cards
			_, err = storage.Draw(context.Background(), deck.DeckID, 1)
			require.NoError(t, err)

			// deck is empty
			_, err = storage.Draw(context.Background(), deck.DeckID, 1)
			require.ErrorIs(t, err, ErrEmptyDeck)
		})

		t.Run(fmt.Sprintf("%s - drawing more cards than deck has -> error", storageName), func(t *testing.T) {
			initial := []cards.Card{
				{Suit: cards.CardSuitClubs, Rank: cards.CardRank(3)},
				{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(8)},
			}

			deck, err := storage.Create(context.Background(), initial, false)
			require.NoError(t, err)

			cards, err := storage.Draw(context.Background(), deck.DeckID, 3)
			require.NoError(t, err)
			require.Len(t, cards, 2)
		})

		t.Run(fmt.Sprintf("%s - drawing removes cards from deck", storageName), func(t *testing.T) {
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

		t.Run(fmt.Sprintf("%s - reshuffle with deck that does not exist -> ErrDeckNotFound error", storageName), func(t *testing.T) {
			ID := uuid.New()
			_, err := storage.Reshuffle(context.Background(), &ID, reverseShuffle)
			require.ErrorIs(t, err, ErrDeckNotFound)
		})

		t.Run(fmt.Sprintf("%s - reshuffle reorders the remaining cards and marks the deck as shuffled", storageName), func(t *testing.T) {
			initial := []cards.Card{
				{Suit: cards.CardSuitClubs, Rank: cards.CardRank(3)},
				{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(8)},
				{Suit: cards.CardSuitHearts, Rank: cards.CardRank(1)},
			}

			deck, err := storage.Create(context.Background(), initial, false)
			require.NoError(t, err)

			// draw a card, so only some cards are left to be reshuffled.
			_, err = storage.Draw(context.Background(), deck.DeckID, 1)
			require.NoError(t, err)

			reshuffled, err := storage.Reshuffle(context.Background(), deck.DeckID, reverseShuffle)
			require.NoError(t, err)
			require.True(t, reshuffled.Shuffled)
			require.Equal(t, []cards.Card{
				{Suit: cards.CardSuitHearts, Rank: cards.CardRank(1)},
				{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(8)},
			}, reshuffled.Cards)

			fetched, err := storage.Get(context.Background(), deck.DeckID)
			require.NoError(t, err)
			require.Equal(t, reshuffled, fetched)
		})

		t.Run(fmt.Sprintf("%s - delete with deck that does not exist -> ErrDeckNotFound error", storageName), func(t *testing.T) {
			ID := uuid.New()
			err := storage.Delete(context.Background(), &ID)
			require.ErrorIs(t, err, ErrDeckNotFound)
		})

		t.Run(fmt.Sprintf("%s - delete with existing deck -> deck is removed", storageName), func(t *testing.T) {
			initial := []cards.Card{{Suit: cards.CardSuitClubs, Rank: cards.CardRank(3)}}
			deck, err := storage.Create(context.Background(), initial, false)
			require.NoError(t, err)

			err = storage.Delete(context.Background(), deck.DeckID)
			require.NoError(t, err)

			_, err = storage.Get(context.Background(), deck.DeckID)
			require.ErrorIs(t, err, ErrDeckNotFound)
		})

		t.Run(fmt.Sprintf("%s - deleting the same deck twice -> second call returns ErrDeckNotFound", storageName), func(t *testing.T) {
			initial := []cards.Card{{Suit: cards.CardSuitClubs, Rank: cards.CardRank(3)}}
			deck, err := storage.Create(context.Background(), initial, false)
			require.NoError(t, err)

			err = storage.Delete(context.Background(), deck.DeckID)
			require.NoError(t, err)

			err = storage.Delete(context.Background(), deck.DeckID)
			require.ErrorIs(t, err, ErrDeckNotFound)
		})
	})
}

// reverseShuffle is a deterministic ShuffleFunc used in tests, so that
// assertions about the resulting order don't have to deal with randomness.
func reverseShuffle(list []cards.Card) []cards.Card {
	reversed := make([]cards.Card, len(list))
	for i, card := range list {
		reversed[len(list)-1-i] = card
	}

	return reversed
}

type StorageImplementation struct {
	CreateFn func() (Storage, error)
}

// Currenly, we only these two implementations.
var storageImplementations = map[string]StorageImplementation{
	"redis": {
		CreateFn: func() (Storage, error) {
			// Requires Redis, configured via REDIS_HOST / REDIS_PORT
			// (docker-compose sets these; CI points them at a service container).
			return NewRedisStorage(nil)
		},
	},
	"in-memory": {
		CreateFn: func() (Storage, error) {
			return NewInMemoryStorage(), nil
		},
	},
}

// Easy way to run a bunch of tests for all available storage implementations.
func runTestForAllImplementations(t *testing.T, test func(string, Storage)) {
	for name, implementation := range storageImplementations {
		storage, err := implementation.CreateFn()
		require.Nil(t, err)
		test(name, storage)
	}
}
