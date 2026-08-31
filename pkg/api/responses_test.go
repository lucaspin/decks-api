package api

import (
	"testing"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/lucaspin/decks-api/pkg/storage"
	"github.com/stretchr/testify/require"
)

func Test__NewCreateDeckResponse(t *testing.T) {
	ID := uuid.New()
	deck := &storage.Dec{
		DeckID:   &ID,
		Shuffled: true,
		Cards: []cards.Card{
			{Suit: cards.CardSuitSpades, Rank: cards.CardRank(1)},
			{Suit: cards.CardSuitHearts, Rank: cards.CardRank(13)},
		},
	}

	response := newCreateDeckResponse(deck)
	require.Equal(t, CreateDeckResponse{
		DeckID:    &ID,
		Shuffled:  true,
		Remaining: 2,
	}, response)
}

func Test__NewOpenDeckResponse(t *testing.T) {
	ID := uuid.New()
	deck := &storage.Deck{
		DeckID:   &ID,
		Shuffled: false,
		Cards: []cards.Card{
			{Suit: cards.CardSuitSpades, Rank: cards.CardRank(1)},
			{Suit: cards.CardSuitDiamonds, Rank: cards.CardRank(11)},
		},
	}

	response := newOpenDeckResponse(deck)
	require.Equal(t, OpenDeckResponse{
		DeckID:    &ID,
		Shuffled:  false,
		Remaining: 2,
		Cards: []Card{
			{Value: "ACE", Suit: "SPADES", Code: "AS"},
			{Value: "JACK", Suit: "DIAMONDS", Code: "JD"},
		},
	}, response)
}

func Test__NewOpenDeckResponse_EmptyDeck(t *testing.T) {
	ID := uuid.New()
	deck := &storage.Deck{DeckID: &ID, Shuffled: false, Cards: []cards.Card{}}

	response := newOpenDeckResponse(deck)
	require.Equal(t, 0, response.Remaining)
	require.Empty(t, response.Cards)
}

func Test__NewDrawCardsResponse(t *testing.T) {
	deckCards := []cards.Card{
		{Suit: cards.CardSuitClubs, Rank: cards.CardRank(12)},
		{Suit: cards.CardSuitHearts, Rank: cards.CardRank(7)},
	}

	response := newDrawCardsResponse(deckCards)
	require.Equal(t, DrawCardsResponse{
		Cards: []Card{
			{Value: "QUEEN", Suit: "CLUBS", Code: "QC"},
			{Value: "7", Suit: "HEARTS", Code: "7H"},
		},
	}, response)
}

func Test__NewDrawCardsResponse_Empty(t *testing.T) {
	response := newDrawCardsResponse([]cards.Card{})
	require.Empty(t, response.Cards)
}
