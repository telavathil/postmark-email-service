package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostmarkToken string
	Port          string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	postmarkToken := os.Getenv("POSTMARK_SERVER_TOKEN")
	if postmarkToken == "" {
		log.Fatal("POSTMARK_SERVER_TOKEN environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		PostmarkToken: postmarkToken,
		Port:          port,
	}
}
