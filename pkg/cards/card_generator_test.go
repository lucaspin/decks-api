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
