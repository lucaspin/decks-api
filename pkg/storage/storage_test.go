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
	})
}

type StorageImplementation struct {
	CreateFn func() (Storage, error)
}

// Currenly, we only these two implementations.
var storageImplementations = map[string]StorageImplementation{
	"redis": {
		CreateFn: func() (Storage, error) {
			// This requires a redis server to be available in this address.
			// This Redis server is created by docker compose.
			// See the docker-compose.yml file.
			return NewRedisStorage(&RedisConfig{Host: "redis", Port: "6379"})
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
