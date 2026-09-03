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

// ShuffleList shuffles the given list of cards in place using a Fisher-Yates
// shuffle, seeded from the current time. It is exposed as a package-level
// helper so that callers without access to a CardGenerator (e.g. the storage
// layer) can shuffle a list of cards without duplicating the algorithm.
func ShuffleList(list []Card) []Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return shuffleWithRand(list, r)
}

func shuffleWithRand(list []Card, r *rand.Rand) []Card {
	for i := range list {
		j := r.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}

	return list
}
