package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lucaspin/decks-api/pkg/storage"
)

type Server struct {
	router     *mux.Router
	httpServer *http.Server
	storage    storage.Storage
}

func NewServer() *Server {
	server := &Server{
		storage: storage.NewStorage(),
	}

	server.InitRouter()
	return server
}

func (s *Server) InitRouter() {
	basePath := "/api/v1alpha"
	s.router = mux.NewRouter().StrictSlash(true)
	s.router.HandleFunc("/", s.HealthCheck).Methods(http.MethodGet)
	s.router.HandleFunc(basePath+"/decks", s.CreateDeck).Methods(http.MethodPost)
	s.router.HandleFunc(basePath+"/decks/{deck_id}", s.OpenDeck).Methods(http.MethodGet)
}

func (s *Server) CreateDeck(w http.ResponseWriter, r *http.Request) {
	deck, err := s.storage.Create(r.Context())
	if err == nil {
		response := newCreateDeckResponse(deck)
		respondWithJSON(w, http.StatusCreated, &response)
		return
	}

	http.Error(w, err.Error(), http.StatusBadRequest)
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

// An endpoint used to check if the server is running.
// Mostly used for Kubernetes probes or Docker health checks.
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (s *Server) Serve(host string, port int) error {
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
