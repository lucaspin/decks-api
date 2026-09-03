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

		t.Run(fmt.Sprintf("%s - shuffle with deck that does not exist -> ErrDeckNotFound error", storageName), func(t *testing.T) {
			ID := uuid.New()
			_, err := storage.Shuffle(context.Background(), &ID)
			require.ErrorIs(t, err, ErrDeckNotFound)
		})

		t.Run(fmt.Sprintf("%s - shuffle preserves the remaining cards and marks the deck as shuffled", storageName), func(t *testing.T) {
			initial, err := cards.CodesToCardList([]string{"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S"})
			require.NoError(t, err)

			deck, err := storage.Create(context.Background(), initial, false)
			require.NoError(t, err)

			shuffled, err := storage.Shuffle(context.Background(), deck.DeckID)
			require.NoError(t, err)
			require.True(t, shuffled.Shuffled)
			require.ElementsMatch(t, cards.CardListToCodes(initial), cards.CardListToCodes(shuffled.Cards))

			fromStorage, err := storage.Get(context.Background(), deck.DeckID)
			require.NoError(t, err)
			require.True(t, fromStorage.Shuffled)
			require.ElementsMatch(t, cards.CardListToCodes(initial), cards.CardListToCodes(fromStorage.Cards))
		})

		t.Run(fmt.Sprintf("%s - shuffle after a draw only reshuffles the remaining cards", storageName), func(t *testing.T) {
			initial, err := cards.CodesToCardList([]string{"AS", "2S", "3S", "4S"})
			require.NoError(t, err)

			deck, err := storage.Create(context.Background(), initial, false)
			require.NoError(t, err)

			drawn, err := storage.Draw(context.Background(), deck.DeckID, 2)
			require.NoError(t, err)
			require.Len(t, drawn, 2)

			shuffled, err := storage.Shuffle(context.Background(), deck.DeckID)
			require.NoError(t, err)
			require.True(t, shuffled.Shuffled)
			require.Len(t, shuffled.Cards, 2)

			remainingCodes := cards.CardListToCodes(append(drawn, shuffled.Cards...))
			require.ElementsMatch(t, cards.CardListToCodes(initial), remainingCodes)
		})
	})
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
