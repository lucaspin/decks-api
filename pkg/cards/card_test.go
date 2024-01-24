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
		{code: "QD", expectErr: false, expectedRank: CardRank(11), expectedSuit: CardSuitDiamonds},
		{code: "JD", expectErr: false, expectedRank: CardRank(12), expectedSuit: CardSuitDiamonds},
		{code: "KD", expectErr: false, expectedRank: CardRank(13), expectedSuit: CardSuitDiamonds},
		{code: "AH", expectErr: false, expectedRank: CardRank(1), expectedSuit: CardSuitHearts},
		{code: "2H", expectErr: false, expectedRank: CardRank(2), expectedSuit: CardSuitHearts},
		{code: "QH", expectErr: false, expectedRank: CardRank(11), expectedSuit: CardSuitHearts},
		{code: "JH", expectErr: false, expectedRank: CardRank(12), expectedSuit: CardSuitHearts},
		{code: "KH", expectErr: false, expectedRank: CardRank(13), expectedSuit: CardSuitHearts},
		{code: "AS", expectErr: false, expectedRank: CardRank(1), expectedSuit: CardSuitSpades},
		{code: "2S", expectErr: false, expectedRank: CardRank(2), expectedSuit: CardSuitSpades},
		{code: "QS", expectErr: false, expectedRank: CardRank(11), expectedSuit: CardSuitSpades},
		{code: "JS", expectErr: false, expectedRank: CardRank(12), expectedSuit: CardSuitSpades},
		{code: "KS", expectErr: false, expectedRank: CardRank(13), expectedSuit: CardSuitSpades},
		{code: "AC", expectErr: false, expectedRank: CardRank(1), expectedSuit: CardSuitClubs},
		{code: "2C", expectErr: false, expectedRank: CardRank(2), expectedSuit: CardSuitClubs},
		{code: "QC", expectErr: false, expectedRank: CardRank(11), expectedSuit: CardSuitClubs},
		{code: "JC", expectErr: false, expectedRank: CardRank(12), expectedSuit: CardSuitClubs},
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
