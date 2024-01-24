package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	router     *mux.Router
	httpServer *http.Server
}

func NewServer() *Server {
	server := &Server{}
	server.InitRouter()
	return server
}

func (s *Server) InitRouter() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", s.HealthCheck).Methods(http.MethodGet)
	s.router = r
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
