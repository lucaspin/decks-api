package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lucaspin/decks-api/pkg/cards"
	"github.com/lucaspin/decks-api/pkg/storage"
)

type Server struct {
	router     *mux.Router
	httpServer *http.Server
	storage    storage.Storage
	generator  *cards.CardGenerator
}

func NewServer(storage storage.Storage) *Server {
	server := &Server{
		storage:   storage,
		generator: cards.NewCardGenerator(),
	}

	server.InitRouter()
	return server
}

func (s *Server) InitRouter() {
	basePath := "/api/v1alpha"
	s.router = mux.NewRouter().StrictSlash(true)
	s.router.HandleFunc(basePath+"/decks", s.CreateDeck).Methods(http.MethodPost)
	s.router.HandleFunc(basePath+"/decks/{deck_id}", s.OpenDeck).Methods(http.MethodGet)
	s.router.HandleFunc(basePath+"/decks/{deck_id}/draw", s.DrawCards).Methods(http.MethodPost)
	s.router.HandleFunc("/", s.HealthCheck).Methods(http.MethodGet)
}

func (s *Server) CreateDeck(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	shuffled := queryParams.Get("shuffled") == "true"
	list, err := s.generator.NewListWithConfig(cards.GeneratorConfig{
		Shuffled: shuffled,
		Codes:    queryParams.Get("cards"),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	deck, err := s.storage.Create(r.Context(), list, shuffled)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := newCreateDeckResponse(deck)
	respondWithJSON(w, http.StatusCreated, &response)
}

func (s *Server) OpenDeck(w http.ResponseWriter, r *http.Request) {
	deckID, err := uuid.Parse(mux.Vars(r)["deck_id"])
	if err != nil {
		http.Error(w, "invalid deck ID", http.StatusBadRequest)
		return
	}

	deck, err := s.storage.Get(r.Context(), &deckID)
	if err == nil {
		response := newOpenDeckResponse(deck)
		respondWithJSON(w, http.StatusOK, &response)
		return
	}

	if errors.Is(err, storage.ErrDeckNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Some unknown error happened when trying to find the deck.
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (s *Server) DrawCards(w http.ResponseWriter, r *http.Request) {
	deckID, err := uuid.Parse(mux.Vars(r)["deck_id"])
	if err != nil {
		http.Error(w, "invalid deck ID", http.StatusBadRequest)
		return
	}

	countFromQuery := r.URL.Query().Get("count")
	if countFromQuery == "" {
		http.Error(w, "count is required", http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(countFromQuery)
	if err != nil {
		http.Error(w, "invalid count", http.StatusBadRequest)
		return
	}

	if count < 0 {
		http.Error(w, "count must be positive", http.StatusBadRequest)
		return
	}

	cards, err := s.storage.Draw(r.Context(), &deckID, count)
	if err == nil {
		response := newDrawCardsResponse(cards)
		respondWithJSON(w, http.StatusOK, &response)
		return
	}

	if errors.Is(err, storage.ErrDeckNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if errors.Is(err, storage.ErrEmptyDeck) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Unknown error drawing cards: %v", err)
	http.Error(w, "unknown error", http.StatusInternalServerError)
}

// An endpoint used to check if the server is running.
// Mostly used for Kubernetes probes or Docker health checks.
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (s *Server) Serve(host string, port int) error {
	log.Printf("Starting server on %s:%d\n", host, port)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler: http.TimeoutHandler(
			handlers.LoggingHandler(os.Stdout, s.router),
			15*time.Second,
			"request timed out",
		),
	}

	return s.httpServer.ListenAndServe()
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(response)
	return err
}
