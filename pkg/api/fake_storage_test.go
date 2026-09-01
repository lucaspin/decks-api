package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/lucaspin/decks-api/pkg/storage"
	"github.com/stretchr/testify/require"
)

// errBoom is a generic, non-sentinel error used to exercise the
// "unknown error" branches of the HTTP handlers, which are otherwise
// unreachable when using the in-memory storage implementation.
var errBoom = errors.New("boom")

// stubStorage is a minimal storage.Storage implementation whose
// methods return a configurable error, so we can exercise the
// handlers' 500 error paths without relying on a real backend.
type stubStorage struct {
	createErr error
	getErr    error
	drawErr   error
}

func (s *stubStorage) Create(ctx context.Context, list []cards.Card, shuffled bool) (*storage.Deck, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}

	ID := uuid.New()
	return &storage.Deck{DeckID: &ID, Shuffled: shuffled, Cards: list}, nil
}

func (s *stubStorage) Get(ctx context.Context, deckID *uuid.UUID) (*storage.Deck, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}

	return &storage.Deck{DeckID: deckID}, nil
}

func (s *stubStorage) Draw(ctx context.Context, deckID *uuid.UUID, count int) ([]cards.Card, error) {
	if s.drawErr != nil {
		return nil, s.drawErr
	}

	return nil, nil
}

func Test__CreateDeck__StorageError(t *testing.T) {
	testServer := NewServer(&stubStorage{createErr: errBoom})
	response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks", nil)
	require.Equal(t, 500, response.Code)
	require.Equal(t, "boom\n", response.Body.String())
}

func Test__OpenDeck__UnknownStorageError(t *testing.T) {
	testServer := NewServer(&stubStorage{getErr: errBoom})
	ID := uuid.New()
	response := execRequest(testServer, http.MethodGet, "/api/v1alpha/decks/"+ID.String(), nil)
	require.Equal(t, 500, response.Code)
	require.Equal(t, "boom\n", response.Body.String())
}

func Test__DrawCards__UnknownStorageError(t *testing.T) {
	testServer := NewServer(&stubStorage{drawErr: errBoom})
	ID := uuid.New()
	response := execRequest(testServer, http.MethodPost, "/api/v1alpha/decks/"+ID.String()+"/draw?count=1", nil)
	require.Equal(t, 500, response.Code)
	require.Equal(t, "unknown error\n", response.Body.String())
}
