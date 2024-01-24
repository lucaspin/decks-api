package api

import (
	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/lucaspin/decks-api/pkg/storage"
)

type CreateDeckResponse struct {
	DeckID    *uuid.UUID `json:"deck_id"`
	Shuffled  bool       `json:"shuffled"`
	Remaining int        `json:"remaining"`
}

func newCreateDeckResponse(deck *storage.Deck) CreateDeckResponse {
	return CreateDeckResponse{
		DeckID:    deck.DeckID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining(),
	}
}

type OpenDeckResponse struct {
	DeckID    *uuid.UUID `json:"deck_id"`
	Shuffled  bool       `json:"shuffled"`
	Remaining int        `json:"remaining"`
	Cards     []Card     `json:"cards"`
}

type DrawCardsResponse struct {
	Cards []Card `json:"cards"`
}

type Card struct {
	Value string
	Suit  string
	Code  string
}

func newOpenDeckResponse(deck *storage.Deck) OpenDeckResponse {
	cards := make([]Card, len(deck.Cards))
	for i, c := range deck.Cards {
		cards[i] = Card{
			Value: c.Rank.String(),
			Suit:  c.Suit.String(),
			Code:  c.Code(),
		}
	}

	return OpenDeckResponse{
		DeckID:    deck.DeckID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining(),
		Cards:     cards,
	}
}

func newDrawCardsResponse(deckCards []cards.Card) DrawCardsResponse {
	cards := make([]Card, len(deckCards))
	for i, c := range deckCards {
		cards[i] = Card{
			Value: c.Rank.String(),
			Suit:  c.Suit.String(),
			Code:  c.Code(),
		}
	}

	return DrawCardsResponse{Cards: cards}
}
