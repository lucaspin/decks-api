package main

import (
	"log"

	"github.com/lucaspin/decks-api/pkg/api"
	"github.com/lucaspin/decks-api/pkg/storage"
)

func main() {
	store, err := storage.NewStorage()
	if err != nil {
		log.Fatalf("error initializing storage: %v", err)
	}

	server := api.NewServer(store)
	err = server.Serve("0.0.0.0", 4000)
	if err != nil {
		log.Fatal(err)
	}
}
