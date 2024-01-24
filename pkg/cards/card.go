package cards

import (
	"fmt"
	"strconv"
)

type CardSuit int
type CardRank int

const (
	CardSuitClubs CardSuit = iota
	CardSuitDiamonds
	CardSuitHearts
	CardSuitSpades
	CardSuitUnknown
)

type Card struct {
	Suit CardSuit
	Rank CardRank
}

func AllSuits() []CardSuit {
	return []CardSuit{CardSuitSpades, CardSuitDiamonds, CardSuitClubs, CardSuitHearts}
}

func NewCardFromCode(code string) (*Card, error) {
	suit, err := CardSuitFromCode(code[len(code)-1])
	if err != nil {
		return nil, err
	}

	rank, err := CardRankFromCode(code[0 : len(code)-1])
	if err != nil {
		return nil, err
	}

	return &Card{Suit: suit, Rank: rank}, nil
}

func (c *Card) Code() string {
	return c.Rank.Code() + c.Suit.Code()
}

func (s *CardSuit) String() string {
	switch *s {
	case CardSuitClubs:
		return "CLUBS"
	case CardSuitDiamonds:
		return "DIAMONDS"
	case CardSuitHearts:
		return "HEARTS"
	case CardSuitSpades:
		return "SPADES"
	default:
		return ""
	}
}

func (s *CardSuit) Code() string {
	switch *s {
	case CardSuitClubs:
		return "C"
	case CardSuitDiamonds:
		return "D"
	case CardSuitHearts:
		return "H"
	case CardSuitSpades:
		return "S"
	default:
		return ""
	}
}

func CardSuitFromCode(code byte) (CardSuit, error) {
	switch code {
	case 'C':
		return CardSuitClubs, nil
	case 'D':
		return CardSuitDiamonds, nil
	case 'H':
		return CardSuitHearts, nil
	case 'S':
		return CardSuitSpades, nil
	default:
		return CardSuitUnknown, fmt.Errorf("invalid suit code '%s'", string(code))
	}
}

func (r *CardRank) String() string {
	switch int(*r) {
	case 1:
		return "ACE"
	case 11:
		return "QUEEN"
	case 12:
		return "JACK"
	case 13:
		return "KING"
	default:
		return fmt.Sprintf("%d", int(*r))
	}
}

func (r *CardRank) Code() string {
	switch int(*r) {
	case 1:
		return "A"
	case 11:
		return "Q"
	case 12:
		return "J"
	case 13:
		return "K"
	default:
		return fmt.Sprintf("%d", int(*r))
	}
}

func CardRankFromCode(code string) (CardRank, error) {
	switch code[0] {
	case 'A':
		return CardRank(1), nil
	case 'Q':
		return CardRank(11), nil
	case 'J':
		return CardRank(12), nil
	case 'K':
		return CardRank(13), nil
	default:
		n, err := strconv.Atoi(string(code))
		if err != nil {
			return CardRank(-1), fmt.Errorf("invalid rank code '%s'", code)
		}

		if n >= 2 && n <= 10 {
			return CardRank(n), nil
		}

		return CardRank(-1), fmt.Errorf("invalid rank code '%s'", code)
	}
}
