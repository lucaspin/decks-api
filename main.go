package main

import (
	"log"
	"os"
	"strconv"

	"github.com/lucaspin/decks-api/pkg/api"
	"github.com/lucaspin/decks-api/pkg/storage"
)

const defaultPort = 4000

func main() {
	store, err := storage.NewStorage()
	if err != nil {
		log.Fatalf("error initializing storage: %v", err)
	}

	server := api.NewServer(store)
	err = server.Serve("0.0.0.0", getPort())
	if err != nil {
		log.Fatal(err)
	}
}

func getPort() int {
	fromEnv := os.Getenv("API_PORT")
	if fromEnv == "" {
		log.Printf("No API_PORT specified - using 4000")
		return defaultPort
	}

	port, err := strconv.Atoi(fromEnv)
	if err != nil {
		log.Printf("invalid API_PORT '%s' specified: %v - using 4000", fromEnv, err)
		return defaultPort
	}

	return port
}
