package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/stretchr/testify/require"
)

func Test__HealthCheckEndpointRespondsWith200(t *testing.T) {
	testServer := NewServer()
	response := execRequest(testServer, http.MethodGet, "/", nil)
	require.Equal(t, response.Code, 200)
}

func Test__CreateDeck(t *testing.T) {
	testServer := NewServer()

	t.Run("deck created -> 201 and proper response", func(t *testing.T) {
		response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks", nil)
		require.Equal(t, response.Code, 201)

		r := &CreateDeckResponse{}
		require.NoError(t, json.NewDecoder(response.Body).Decode(&r))
		require.NotNil(t, r.DeckID)
		require.False(t, r.Shuffled)
		require.Equal(t, r.Remaining, 52)
	})
}

func Test__OpenDeck(t *testing.T) {
	testServer := NewServer()

	t.Run("invalid deck ID -> 400", func(t *testing.T) {
		response := execRequest(testServer, http.MethodGet, "/api/v1alpha/decks/not-a-valid-uuid", nil)
		require.Equal(t, response.Code, 400)
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
}

func requireFullUnshuffledDeck(t *testing.T, list []Card) {
	cardsPerSuit := 13
	for i, suit := range cards.AllSuits() {
		startIndex := 0
		if i > 0 {
			startIndex = cardsPerSuit * i
		}

		endIndex := (cardsPerSuit - 1) * (i + 1)
		cardsForSuit := list[startIndex:endIndex]
		for j, c := range cardsForSuit {
			rank := cards.CardRank(j + 1)
			require.Equal(t, suit.String(), c.Suit)
			require.Equal(t, rank.String(), c.Value)
			require.Equal(t, rank.Code()+suit.Code(), c.Code)
		}
	}
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
