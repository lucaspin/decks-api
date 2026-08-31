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
		{code: "", expectErr: true, errMessage: "invalid card code ''"},
		{code: "A", expectErr: true, errMessage: "invalid card code 'A'"},
		{code: "S", expectErr: true, errMessage: "invalid card code 'S'"},
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

func Test__CardSuit_StringAndCode(t *testing.T) {
	type testCase struct {
		suit           CardSuit
		expectedString string
		expectedCode   string
	}

	for _, tc := range []testCase{
		{suit: CardSuitClubs, expectedString: "CLUBS", expectedCode: "C"},
		{suit: CardSuitDiamonds, expectedString: "DIAMONDS", expectedCode: "D"},
		{suit: CardSuitHearts, expectedString: "HEARTS", expectedCode: "H"},
		{suit: CardSuitSpades, expectedString: "SPADES", expectedCode: "S"},
		{suit: CardSuitUnknown, expectedString: "", expectedCode: ""},
		{suit: CardSuit(999), expectedString: "", expectedCode: ""},
	} {
		suit := tc.suit
		require.Equal(t, tc.expectedString, suit.String())
		require.Equal(t, tc.expectedCode, suit.Code())
	}
}

func Test__CardRank_StringAndCode(t *testing.T) {
	type testCase struct {
		rank           CardRank
		expectedString string
		expectedCode   string
	}

	for _, tc := range []testCase{
		{rank: CardRank(1), expectedString: "ACE", expectedCode: "A"},
		{rank: CardRank(2), expectedString: "2", expectedCode: "2"},
		{rank: CardRank(10), expectedString: "10", expectedCode: "10"},
		{rank: CardRank(11), expectedString: "JACK", expectedCode: "J"},
		{rank: CardRank(12), expectedString: "QUEEN", expectedCode: "Q"},
		{rank: CardRank(13), expectedString: "KING", expectedCode: "K"},
	} {
		rank := tc.rank
		require.Equal(t, tc.expectedString, rank.String())
		require.Equal(t, tc.expectedCode, rank.Code())
	}
}

func Test__CardListToCodesAndCodesToCardList(t *testing.T) {
	t.Run("round trip", func(t *testing.T) {
		codes := []string{"AS", "2S", "10D", "JC", "QH", "KH"}
		list, err := CodesToCardList(codes)
		require.NoError(t, err)
		require.Equal(t, codes, CardListToCodes(list))
	})

	t.Run("trims spaces around codes", func(t *testing.T) {
		list, err := CodesToCardList([]string{" AS ", " KD"})
		require.NoError(t, err)
		require.Equal(t, []string{"AS", "KD"}, CardListToCodes(list))
	})

	t.Run("invalid code returns error", func(t *testing.T) {
		_, err := CodesToCardList([]string{"AS", "invalid"})
		require.Error(t, err)
	})

	t.Run("empty list", func(t *testing.T) {
		list, err := CodesToCardList([]string{})
		require.NoError(t, err)
		require.Empty(t, CardListToCodes(list))
	})
}
