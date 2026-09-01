package cards

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__NewCardFromCode(t *testing.T) {
	type testCase struct {
		code         string
		expectedRank CardRank
		expectedSuit CardSuit
		expectErr    bool
		errMessage   string
	}

	for _, tc := range []testCase{
		{code: "AL", expectErr: true, errMessage: "invalid suit code 'L'"},
		{code: "A9", expectErr: true, errMessage: "invalid suit code '9'"},
		{code: "11D", expectErr: true, errMessage: "invalid rank code '11'"},
		{code: "-11D", expectErr: true, errMessage: "invalid rank code '-11'"},
		{code: "99D", expectErr: true, errMessage: "invalid rank code '99'"},
		{code: "AD", expectErr: false, expectedRank: CardRank(1), expectedSuit: CardSuitDiamonds},
		{code: "2D", expectErr: false, expectedRank: CardRank(2), expectedSuit: CardSuitDiamonds},
		{code: "JD", expectErr: false, expectedRank: CardRank(11), expectedSuit: CardSuitDiamonds},
		{code: "QD", expectErr: false, expectedRank: CardRank(12), expectedSuit: CardSuitDiamonds},
		{code: "KD", expectErr: false, expectedRank: CardRank(13), expectedSuit: CardSuitDiamonds},
		{code: "AH", expectErr: false, expectedRank: CardRank(1), expectedSuit: CardSuitHearts},
		{code: "2H", expectErr: false, expectedRank: CardRank(2), expectedSuit: CardSuitHearts},
		{code: "JH", expectErr: false, expectedRank: CardRank(11), expectedSuit: CardSuitHearts},
		{code: "QH", expectErr: false, expectedRank: CardRank(12), expectedSuit: CardSuitHearts},
		{code: "KH", expectErr: false, expectedRank: CardRank(13), expectedSuit: CardSuitHearts},
		{code: "AS", expectErr: false, expectedRank: CardRank(1), expectedSuit: CardSuitSpades},
		{code: "2S", expectErr: false, expectedRank: CardRank(2), expectedSuit: CardSuitSpades},
		{code: "JS", expectErr: false, expectedRank: CardRank(11), expectedSuit: CardSuitSpades},
		{code: "QS", expectErr: false, expectedRank: CardRank(12), expectedSuit: CardSuitSpades},
		{code: "KS", expectErr: false, expectedRank: CardRank(13), expectedSuit: CardSuitSpades},
		{code: "AC", expectErr: false, expectedRank: CardRank(1), expectedSuit: CardSuitClubs},
		{code: "2C", expectErr: false, expectedRank: CardRank(2), expectedSuit: CardSuitClubs},
		{code: "JC", expectErr: false, expectedRank: CardRank(11), expectedSuit: CardSuitClubs},
		{code: "QC", expectErr: false, expectedRank: CardRank(12), expectedSuit: CardSuitClubs},
		{code: "KC", expectErr: false, expectedRank: CardRank(13), expectedSuit: CardSuitClubs},
	} {
		card, err := NewCardFromCode(tc.code)
		if tc.expectErr {
			require.ErrorContains(t, err, tc.errMessage)
		} else {
			require.NoError(t, err)
			require.Equal(t, &Card{Rank: tc.expectedRank, Suit: tc.expectedSuit}, card)
		}
	}
}

func Test__CardSuit__String(t *testing.T) {
	type testCase struct {
		suit     CardSuit
		expected string
	}

	for _, tc := range []testCase{
		{suit: CardSuitClubs, expected: "CLUBS"},
		{suit: CardSuitDiamonds, expected: "DIAMONDS"},
		{suit: CardSuitHearts, expected: "HEARTS"},
		{suit: CardSuitSpades, expected: "SPADES"},
		{suit: CardSuitUnknown, expected: ""},
		{suit: CardSuit(99), expected: ""},
	} {
		suit := tc.suit
		require.Equal(t, tc.expected, suit.String())
	}
}

func Test__CardSuit__Code(t *testing.T) {
	type testCase struct {
		suit     CardSuit
		expected string
	}

	for _, tc := range []testCase{
		{suit: CardSuitClubs, expected: "C"},
		{suit: CardSuitDiamonds, expected: "D"},
		{suit: CardSuitHearts, expected: "H"},
		{suit: CardSuitSpades, expected: "S"},
		{suit: CardSuitUnknown, expected: ""},
		{suit: CardSuit(99), expected: ""},
	} {
		suit := tc.suit
		require.Equal(t, tc.expected, suit.Code())
	}
}

func Test__CardRank__String(t *testing.T) {
	type testCase struct {
		rank     CardRank
		expected string
	}

	for _, tc := range []testCase{
		{rank: CardRank(1), expected: "ACE"},
		{rank: CardRank(2), expected: "2"},
		{rank: CardRank(3), expected: "3"},
		{rank: CardRank(4), expected: "4"},
		{rank: CardRank(5), expected: "5"},
		{rank: CardRank(6), expected: "6"},
		{rank: CardRank(7), expected: "7"},
		{rank: CardRank(8), expected: "8"},
		{rank: CardRank(9), expected: "9"},
		{rank: CardRank(10), expected: "10"},
		{rank: CardRank(11), expected: "JACK"},
		{rank: CardRank(12), expected: "QUEEN"},
		{rank: CardRank(13), expected: "KING"},
	} {
		rank := tc.rank
		require.Equal(t, tc.expected, rank.String())
	}
}

func Test__CardRank__Code(t *testing.T) {
	type testCase struct {
		rank     CardRank
		expected string
	}

	for _, tc := range []testCase{
		{rank: CardRank(1), expected: "A"},
		{rank: CardRank(2), expected: "2"},
		{rank: CardRank(3), expected: "3"},
		{rank: CardRank(4), expected: "4"},
		{rank: CardRank(5), expected: "5"},
		{rank: CardRank(6), expected: "6"},
		{rank: CardRank(7), expected: "7"},
		{rank: CardRank(8), expected: "8"},
		{rank: CardRank(9), expected: "9"},
		{rank: CardRank(10), expected: "10"},
		{rank: CardRank(11), expected: "J"},
		{rank: CardRank(12), expected: "Q"},
		{rank: CardRank(13), expected: "K"},
	} {
		rank := tc.rank
		require.Equal(t, tc.expected, rank.Code())
	}
}

func Test__Card__Code(t *testing.T) {
	type testCase struct {
		card     Card
		expected string
	}

	for _, tc := range []testCase{
		{card: Card{Suit: CardSuitSpades, Rank: CardRank(1)}, expected: "AS"},
		{card: Card{Suit: CardSuitClubs, Rank: CardRank(10)}, expected: "10C"},
		{card: Card{Suit: CardSuitDiamonds, Rank: CardRank(13)}, expected: "KD"},
	} {
		card := tc.card
		require.Equal(t, tc.expected, card.Code())
	}
}

func Test__AllSuits(t *testing.T) {
	require.Equal(t, []CardSuit{
		CardSuitSpades, CardSuitDiamonds, CardSuitClubs, CardSuitHearts,
	}, AllSuits())
}

func Test__CardListToCodes(t *testing.T) {
	list := []Card{
		{Suit: CardSuitSpades, Rank: CardRank(1)},
		{Suit: CardSuitDiamonds, Rank: CardRank(13)},
		{Suit: CardSuitClubs, Rank: CardRank(10)},
	}

	require.Equal(t, []string{"AS", "KD", "10C"}, CardListToCodes(list))
}

func Test__CodesToCardList(t *testing.T) {
	t.Run("valid codes, including surrounding whitespace, are converted", func(t *testing.T) {
		list, err := CodesToCardList([]string{"  AS ", "KD ", " 10C"})
		require.NoError(t, err)
		require.Equal(t, []Card{
			{Suit: CardSuitSpades, Rank: CardRank(1)},
			{Suit: CardSuitDiamonds, Rank: CardRank(13)},
			{Suit: CardSuitClubs, Rank: CardRank(10)},
		}, list)
	})

	t.Run("invalid code propagates the error", func(t *testing.T) {
		list, err := CodesToCardList([]string{"AS", "14C"})
		require.Nil(t, list)
		require.ErrorContains(t, err, "invalid rank code '14'")
	})
}
