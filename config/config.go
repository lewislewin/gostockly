package config

import (
	"log"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
}

func LoadConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	return &Config{
		DatabaseURL: dbURL,
		JWTSecret:   jwtSecret,
	}
}
