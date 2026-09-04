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

// Shuffles the given list of cards in place, using the CardGenerator's own
// random source, and returns it for convenience.
func (g *CardGenerator) Shuffle(list []Card) []Card {
	shuffleWithRand(list, g.rand)
	return list
}

// A package-local random source used by the standalone Shuffle helper below,
// so callers that don't have (or want) a CardGenerator instance - e.g. the
// storage layer - can still shuffle a list of cards.
var packageRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Shuffles the given list of cards in place, using the Fisher-Yates algorithm.
// This is a package-level equivalent of CardGenerator.Shuffle, kept separate
// so it can be used without depending on a CardGenerator instance.
func Shuffle(list []Card) {
	shuffleWithRand(list, packageRand)
}

func shuffleWithRand(list []Card, r *rand.Rand) {
	for i := range list {
		j := r.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}
}
