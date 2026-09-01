package storage

import (
	"testing"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/stretchr/testify/require"
)

func Test__Deck_Remaining(t *testing.T) {
	ID := uuid.New()

	t.Run("empty deck -> 0", func(t *testing.T) {
		deck := &Deck{DeckID: &ID, Cards: []cards.Card{}}
		require.Equal(t, 0, deck.Remaining())
	})

	t.Run("deck with cards -> number of cards left", func(t *testing.T) {
		deck := &Deck{
			DeckID: &ID,
			Cards: []cards.Card{
				{Suit: cards.CardSuitSpades, Rank: cards.CardRank(1)},
				{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(13)},
			},
		}

		require.Equal(t, 2, deck.Remaining())
	})
}
