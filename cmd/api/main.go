package main

import (
	"log"
	"net/http"

	"gostockly/config"
	"gostockly/pkg/api"
	"gostockly/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg := config.LoadConfig()
	db := database.Connect()

	router := api.NewRouter(cfg, db)

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
