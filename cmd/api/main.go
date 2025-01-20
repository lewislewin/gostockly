package main

import (
	"log"
	"net/http"

	"gostockly/config"
	"gostockly/pkg/api"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialise router
	router := api.NewRouter(cfg)

	// Start the server
	log.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
