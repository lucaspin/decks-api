package storage

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/stretchr/testify/require"
)

// This test proves that draws against PostgresStorage are atomic:
// firing multiple concurrent draws of 1 card each against a deck
// results in every card being drawn exactly once, with no duplicates
// and no errors - something the Redis backend can't guarantee, since
// its two keys are updated with separate, non-atomic operations.
func Test__PostgresStorage__ConcurrentDraws(t *testing.T) {
	storage, err := NewPostgresStorage(&PostgresConfig{
		Host:     "postgres",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		Database: "decks",
		SSLMode:  "disable",
	})

	require.NoError(t, err)

	const totalCards = 10
	initial := make([]cards.Card, totalCards)
	for i := 0; i < totalCards; i++ {
		initial[i] = cards.Card{Suit: cards.CardSuitSpades, Rank: cards.CardRank((i % 13) + 1)}
	}

	deck, err := storage.Create(context.Background(), initial, false)
	require.NoError(t, err)

	var wg sync.WaitGroup
	drawnCh := make(chan cards.Card, totalCards)
	errCh := make(chan error, totalCards)

	for i := 0; i < totalCards; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			drawn, err := storage.Draw(context.Background(), deck.DeckID, 1)
			if err != nil {
				errCh <- err
				return
			}

			for _, c := range drawn {
				drawnCh <- c
			}
		}()
	}

	wg.Wait()
	close(drawnCh)
	close(errCh)

	for err := range errCh {
		require.NoError(t, err)
	}

	seen := map[string]int{}
	for c := range drawnCh {
		seen[c.Code()]++
	}

	require.Len(t, seen, totalCards, "expected every card to be drawn exactly once")
	for code, count := range seen {
		require.Equal(t, 1, count, fmt.Sprintf("card %s was drawn %d times", code, count))
	}

	// The deck should now be empty.
	_, err = storage.Draw(context.Background(), deck.DeckID, 1)
	require.ErrorIs(t, err, ErrEmptyDeck)
}
