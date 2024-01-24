package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test__HealthCheckEndpointRespondsWith200(t *testing.T) {
	testServer := NewServer()
	response := execRequest(testServer, http.MethodGet, "/", nil)
	require.Equal(t, response.Code, 200)
}

func Test__CreateDeck(t *testing.T) {
	testServer := NewServer()

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
	codes := make([]string, len(list))
	for i, card := range list {
		codes[i] = card.Code
	}

	require.Equal(t, []string{
		"AS", "2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "10S", "QS", "JS", "KS",
		"AD", "2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "10D", "QD", "JD", "KD",
		"AC", "2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "10C", "QC", "JC", "KC",
		"AH", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "10H", "QH", "JH", "KH",
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
