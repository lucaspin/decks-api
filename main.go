package main

import (
	"log"

	"github.com/lucaspin/decks-api/pkg/api"
)

func main() {
	server := api.NewServer()
	err := server.Serve("0.0.0.0", 4000)
	if err != nil {
		log.Fatal(err)
	}
}
