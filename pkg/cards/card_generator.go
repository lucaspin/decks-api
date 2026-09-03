package cards

import (
	"math/rand"
	"strings"
	"time"
)

type CardGenerator struct {
	rand *rand.Rand
}

type GeneratorConfig struct {
	Shuffled bool
	Codes    string
}

func NewCardGenerator() *CardGenerator {
	return &CardGenerator{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (g *CardGenerator) NewListWithConfig(config GeneratorConfig) ([]Card, error) {
	var list []Card
	if config.Codes == "" {
		list = g.FullCardList()
	} else {
		l, err := CodesToCardList(strings.Split(config.Codes, ","))
		if err != nil {
			return nil, err
		}

		list = l
	}

	if config.Shuffled {
		return g.Shuffle(list), nil
	}

	return list, nil
}

func (g *CardGenerator) FullCardList() []Card {
	cards := []Card{}
	for _, suit := range AllSuits() {
		for i := 1; i <= 13; i++ {
			cards = append(cards, Card{Suit: suit, Rank: CardRank(i)})
		}
	}

	return cards
}

func (g *CardGenerator) Shuffle(list []Card) []Card {
	return ShuffleList(list)
}

// ShuffleList performs an in-place Fisher-Yates shuffle of the given list of
// cards, using a locally seeded random source. It returns the same list,
// shuffled, so callers can chain it if convenient.
//
// This is exposed as a package-level function (rather than only a method on
// CardGenerator) so that other packages, such as storage, can reshuffle
// cards without needing to depend on a CardGenerator instance.
func ShuffleList(list []Card) []Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range list {
		j := r.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}

	return list
}
