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

type ReshuffleDeckResponse struct {
	DeckID    *uuid.UUID `json:"deck_id"`
	Shuffled  bool       `json:"shuffled"`
	Remaining int        `json:"remaining"`
	Cards     []Card     `json:"cards"`
}

type Card struct {
	Value string
	Suit  string
	Code  string
}

func newOpenDeckResponse(deck *storage.Deck) OpenDeckResponse {
	return OpenDeckResponse{
		DeckID:    deck.DeckID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining(),
		Cards:     cardListToResponseCards(deck.Cards),
	}
}

func newDrawCardsResponse(deckCards []cards.Card) DrawCardsResponse {
	return DrawCardsResponse{Cards: cardListToResponseCards(deckCards)}
}

func newReshuffleDeckResponse(deck *storage.Deck) ReshuffleDeckResponse {
	return ReshuffleDeckResponse{
		DeckID:    deck.DeckID,
		Shuffled:  deck.Shuffled,
		Remaining: deck.Remaining(),
		Cards:     cardListToResponseCards(deck.Cards),
	}
}

func cardListToResponseCards(deckCards []cards.Card) []Card {
	list := make([]Card, len(deckCards))
	for i, c := range deckCards {
		list[i] = Card{
			Value: c.Rank.String(),
			Suit:  c.Suit.String(),
			Code:  c.Code(),
		}
	}

	return list
}
