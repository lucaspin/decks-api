package cards

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__Generator(t *testing.T) {
	generator := NewCardGenerator()

	t.Run("generate with default arguments -> full, unshuffled deck", func(t *testing.T) {
		cards, err := generator.NewListWithConfig(GeneratorConfig{})
		require.NoError(t, err)
		require.Len(t, cards, 52)
		requireFullUnshuffledDeck(t, cards)
	})

	t.Run("generate full shuffled deck", func(t *testing.T) {
		cards, err := generator.NewListWithConfig(GeneratorConfig{Shuffled: true})
		require.NoError(t, err)
		require.Len(t, cards, 52)
		requireFullShuffledDeck(t, cards)
	})

	t.Run("create with specific valid cards -> not full deck", func(t *testing.T) {
		cards, err := generator.NewListWithConfig(GeneratorConfig{Codes: "AS,KD,AC,2C,KH,10D"})
		require.NoError(t, err)
		require.Equal(t, []Card{
			{Suit: CardSuitSpades, Rank: 1},
			{Suit: CardSuitDiamonds, Rank: 13},
			{Suit: CardSuitClubs, Rank: 1},
			{Suit: CardSuitClubs, Rank: 2},
			{Suit: CardSuitHearts, Rank: 13},
			{Suit: CardSuitDiamonds, Rank: 10},
		}, cards)
	})

	t.Run("create with specific invalid rank -> error", func(t *testing.T) {
		_, err := generator.NewListWithConfig(GeneratorConfig{Codes: "AS,14C"})
		require.ErrorContains(t, err, "invalid rank code '14'")
	})

	t.Run("create with specific invalid suit -> error", func(t *testing.T) {
		_, err := generator.NewListWithConfig(GeneratorConfig{Codes: "AS,10C,Q?"})
		require.ErrorContains(t, err, "invalid suit code '?'")
	})
}

func requireFullUnshuffledDeck(t *testing.T, list []Card) {
	codes := make([]string, len(list))
	for i, card := range list {
		codes[i] = card.Code()
	}

	require.Equal(t, []string{
		"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "10S", "JS", "QS", "KS",
		"AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "10D", "JD", "QD", "KD",
		"AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "10C", "JC", "QC", "KC",
		"AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "10H", "JH", "QH", "KH",
	}, codes)
}

func requireFullShuffledDeck(t *testing.T, list []Card) {
	codes := make([]string, len(list))
	for i, card := range list {
		codes[i] = card.Code()
	}

	// cards are not following the unshuffled order.
	require.NotEqual(t, []string{
		"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "10S", "JS", "QS", "KS",
		"AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "10D", "JD", "QD", "KD",
		"AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "10C", "JC", "QC", "KC",
		"AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "10H", "JH", "QH", "KH",
	}, codes)
}
