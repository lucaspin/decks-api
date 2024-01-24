package cards

import "fmt"

type Card struct {
	Suit CardSuit
	Rank CardRank
}

func (c *Card) Code() string {
	return c.Rank.Code() + c.Suit.Code()
}

type CardSuit int
type CardRank int

const (
	CardSuitClubs CardSuit = iota
	CardSuitDiamonds
	CardSuitHearts
	CardSuitSpades
)

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

func AllSuits() []CardSuit {
	return []CardSuit{CardSuitSpades, CardSuitDiamonds, CardSuitClubs, CardSuitHearts}
}

func New() []Card {
	cards := []Card{}
	for _, suit := range AllSuits() {
		for i := 1; i <= 13; i++ {
			cards = append(cards, Card{Suit: suit, Rank: CardRank(i)})
		}
	}

	return cards
}
