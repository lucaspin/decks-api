package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/lucaspin/decks-api/pkg/storage"
	"github.com/stretchr/testify/require"
)

func Test__HealthCheckEndpointRespondsWith200(t *testing.T) {
	testServer := NewServer(storage.NewInMemoryStorage())
	response := execRequest(testServer, http.MethodGet, "/", nil)
	require.Equal(t, response.Code, 200)
}

func Test__CreateDeck(t *testing.T) {
	testServer := NewServer(storage.NewInMemoryStorage())

	t.Run("default deck created", func(t *testing.T) {
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks", nil)
		require.Equal(t, response.Code, 201)

		r := &CreateDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&r))
		require.NotNil(t, r.DeckID)
		require.False(t, r.Shuffled)
		require.Equal(t, r.Remaining, 52)
	})

	t.Run("deck can be created with specific cards", func(t *testing.T) {
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks?cards=AS,KD,AC,7H", nil)
		require.Equal(t, response.Code, 201)

		createResponse := &CreateDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&createResponse))
		require.NotNil(t, createResponse.DeckID)
		require.False(t, createResponse.Shuffled)
		require.Equal(t, createResponse.Remaining, 4)

		response = execRequest(testServer, http.MethodGet, "/api/v1alpha/decks/"+createResponse.DeckID.String(), nil)
		require.Equal(t, response.Code, 200)
		openResponse := &OpenDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&openResponse))
		require.Equal(t, &OpenDeckResponse{
			DeckID:    createResponse.DeckID,
			Shuffled:  false,
			Remaining: 4,
			Cards: []Card{
				{Value: "ACE", Suit: "SPADES", Code: "AS"},
				{Value: "KING", Suit: "DIAMONDS", Code: "KD"},
				{Value: "ACE", Suit: "CLUBS", Code: "AC"},
				{Value: "7", Suit: "HEARTS", Code: "7H"},
			},
		}, openResponse)
	})

	t.Run("deck cannot be created with invalid cards", func(t *testing.T) {
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks?cards=AS,KD,14C", nil)
		require.Equal(t, response.Code, 400)
		require.Equal(t, response.Body.String(), "invalid rank code '14'\n")
	})

	t.Run("deck can be created shuffled", func(t *testing.T) {
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks?shuffled=true", nil)
		require.Equal(t, response.Code, 201)

		r := &CreateDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&r))
		require.NotNil(t, r.DeckID)
		require.True(t, r.Shuffled)
		require.Equal(t, r.Remaining, 52)
	})

	t.Run("storage error on create -> 500", func(t *testing.T) {
		failingServer := NewServer(&failingStorage{})
		response := execRequest(failingServer, http.MethodPost, "/api/v1alpha/decks", nil)
		require.Equal(t, response.Code, 500)
		require.Equal(t, response.Body.String(), "storage error: create\n")
	})
}

func Test__OpenDeck(t *testing.T) {
	testServer := NewServer(storage.NewInMemoryStorage())

	t.Run("invalid deck ID -> 400", func(t *testing.T) {
		response := execRequest(testServer, http.MethodGet, "/api/v1alpha/decks/not-a-valid-uuid", nil)
		require.Equal(t, response.Code, 400)
		require.Equal(t, response.Body.String(), "invalid deck ID\n")
	})

	t.Run("deck that does not exist -> 404", func(t *testing.T) {
		ID := uuid.New()
		response := execRequest(testServer, http.MethodGet, "/api/v1alpha/decks/"+ID.String(), nil)
		require.Equal(t, response.Code, 404)
	})

	t.Run("deck that exists -> 200 with proper response", func(t *testing.T) {
		// deck is created
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks", nil)
		require.Equal(t, response.Code, 201)
		createResponse := &CreateDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&createResponse))

		// deck is opened
		response = execRequest(testServer, http.MethodGet, "/api/v1alpha/decks/"+createResponse.DeckID.String(), nil)
		require.Equal(t, response.Code, 200)
		openResponse := &OpenDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&openResponse))
		require.Equal(t, createResponse.DeckID.String(), openResponse.DeckID.String())
		require.False(t, openResponse.Shuffled)
		require.Equal(t, openResponse.Remaining, len(openResponse.Cards))
		requireFullUnshuffledDeck(t, openResponse.Cards)
	})

	t.Run("storage error on get -> 500", func(t *testing.T) {
		failingServer := NewServer(&failingStorage{})
		ID := uuid.New()
		response := execRequest(failingServer, http.MethodGet, "/api/v1alpha/decks/"+ID.String(), nil)
		require.Equal(t, response.Code, 500)
		require.Equal(t, response.Body.String(), "storage error: get\n")
	})
}

func Test__DrawCards(t *testing.T) {
	testServer := NewServer(storage.NewInMemoryStorage())

	t.Run("invalid deck ID -> 400", func(t *testing.T) {
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/not-a-valid-uuid/draw?count=1", nil)
		require.Equal(t, response.Code, 400)
		require.Equal(t, response.Body.String(), "invalid deck ID\n")
	})

	t.Run("deck that does not exist -> 404", func(t *testing.T) {
		ID := uuid.New()
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/"+ID.String()+"/draw?count=1", nil)
		require.Equal(t, response.Code, 404)
	})

	t.Run("negative count -> 400", func(t *testing.T) {
		deckID := createDeck(t, testServer)
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/"+deckID+"/draw?count=-1", nil)
		require.Equal(t, response.Code, 400)
		require.Equal(t, response.Body.String(), "count must be positive\n")
	})

	t.Run("invalid count -> 400", func(t *testing.T) {
		deckID := createDeck(t, testServer)
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/"+deckID+"/draw?count=not-a-number", nil)
		require.Equal(t, response.Code, 400)
		require.Equal(t, response.Body.String(), "invalid count\n")
	})

	t.Run("missing count -> 400", func(t *testing.T) {
		deckID := createDeck(t, testServer)
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/"+deckID+"/draw", nil)
		require.Equal(t, response.Code, 400)
		require.Equal(t, response.Body.String(), "count is required\n")
	})

	t.Run("count=0 -> reaches storage, 200 with no cards drawn", func(t *testing.T) {
		deckID := createDeck(t, testServer)
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/"+deckID+"/draw?count=0", nil)
		require.Equal(t, response.Code, 200)

		drawResponse := &DrawCardsResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&drawResponse))
		require.Len(t, drawResponse.Cards, 0)
	})

	t.Run("successful draw -> 200 with the drawn cards, removed from the deck", func(t *testing.T) {
		// deck is created with known cards
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks?cards=AS,KD,AC,7H", nil)
		require.Equal(t, response.Code, 201)
		createResponse := &CreateDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&createResponse))
		deckID := createResponse.DeckID.String()

		// draw the first two cards
		response = execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/"+deckID+"/draw?count=2", nil)
		require.Equal(t, response.Code, 200)

		drawResponse := &DrawCardsResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&drawResponse))
		require.Equal(t, []Card{
			{Value: "ACE", Suit: "SPADES", Code: "AS"},
			{Value: "KING", Suit: "DIAMONDS", Code: "KD"},
		}, drawResponse.Cards)

		// the drawn cards are gone from the deck, and the remaining ones stay
		response = execRequest(testServer, http.MethodGet, "/api/v1alpha/decks/"+deckID, nil)
		require.Equal(t, response.Code, 200)
		openResponse := &OpenDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&openResponse))
		require.Equal(t, openResponse.Remaining, 2)
		require.Equal(t, []Card{
			{Value: "ACE", Suit: "CLUBS", Code: "AC"},
			{Value: "7", Suit: "HEARTS", Code: "7H"},
		}, openResponse.Cards)
	})

	t.Run("draw more than available -> 200 with only what is left", func(t *testing.T) {
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks?cards=AS,KD", nil)
		require.Equal(t, response.Code, 201)
		createResponse := &CreateDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&createResponse))
		deckID := createResponse.DeckID.String()

		response = execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/"+deckID+"/draw?count=5", nil)
		require.Equal(t, response.Code, 200)

		drawResponse := &DrawCardsResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&drawResponse))
		require.Equal(t, []Card{
			{Value: "ACE", Suit: "SPADES", Code: "AS"},
			{Value: "KING", Suit: "DIAMONDS", Code: "KD"},
		}, drawResponse.Cards)
	})

	t.Run("drawing from an empty deck -> 400", func(t *testing.T) {
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks?cards=AS", nil)
		require.Equal(t, response.Code, 201)
		createResponse := &CreateDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&createResponse))
		deckID := createResponse.DeckID.String()

		// draw the only card
		response = execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/"+deckID+"/draw?count=1", nil)
		require.Equal(t, response.Code, 200)

		// deck is now empty
		response = execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/"+deckID+"/draw?count=1", nil)
		require.Equal(t, response.Code, 400)
		require.Equal(t, response.Body.String(), "deck has no more cards\n")
	})

	t.Run("storage error on draw -> 500", func(t *testing.T) {
		failingServer := NewServer(&failingStorage{})
		ID := uuid.New()
		response := execRequest(failingServer, http.MethodPost, "/api/v1alpha/decks/"+ID.String()+"/draw?count=1", nil)
		require.Equal(t, response.Code, 500)
		require.Equal(t, response.Body.String(), "unknown error\n")
	})
}

// failingStorage is a storage.Storage implementation that always returns
// a generic (non-sentinel) error. It's used to exercise the unknown-error
// (500) branches of the handlers, which the sentinel-error-driven tests
// above cannot reach.
type failingStorage struct{}

func (s *failingStorage) Create(ctx context.Context, list []cards.Card, shuffled bool) (*storage.Deck, error) {
	return nil, fmt.Errorf("storage error: create")
}

func (s *failingStorage) Get(ctx context.Context, deckID *uuid.UUID) (*storage.Deck, error) {
	return nil, fmt.Errorf("storage error: get")
}

func (s *failingStorage) Draw(ctx context.Context, deckID *uuid.UUID, count int) ([]cards.Card, error) {
	return nil, fmt.Errorf("storage error: draw")
}

func requireFullUnshuffledDeck(t *testing.T, list []Card) {
	codes := make([]string, len(list))
	for i, card := range list {
		codes[i] = card.Code
	}

	require.Equal(t, []string{
		"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "10S", "JS", "QS", "KS",
		"AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "10D", "JD", "QD", "KD",
		"AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "10C", "JC", "QC", "KC",
		"AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "10H", "JH", "QH", "KH",
	}, codes)
}

func execRequest(server *Server, method, path string, body interface{}) *httptest.ResponseRecorder {
	stringBody := ""

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		stringBody = string(jsonBody)
	} else {
		stringBody = ""
	}

	bodyReader := strings.NewReader(stringBody)
	req, _ := http.NewRequest(method, path, bodyReader)
	rr := httptest.NewRecorder()
	server.router.ServeHTTP(rr, req)
	return rr
}

func createDeck(t *testing.T, server *Server) string {
	response := execRequest(server, http.MethodPost, "/api/v1alpha/decks", nil)
	require.Equal(t, response.Code, 201)
	createResponse := &CreateDeckResponse{}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&createResponse))
	return createResponse.DeckID.String()
}
