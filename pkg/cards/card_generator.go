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
	return shuffleWithRand(list, g.rand)
}

// ShuffleList shuffles the given list of cards in place, using a locally
// seeded random source. This is used by callers that don't have access to
// a CardGenerator (e.g. the storage layer), so there is a single shuffle
// implementation shared across the codebase.
func ShuffleList(list []Card) []Card {
	return shuffleWithRand(list, rand.New(rand.NewSource(time.Now().UnixNano())))
}

func shuffleWithRand(list []Card, r *rand.Rand) []Card {
	for i := range list {
		j := r.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}

	return list
}
