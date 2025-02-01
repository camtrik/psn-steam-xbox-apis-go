package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PSNRefreshToken string
	PSNAccountID    string
}

func Load() *Config {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Printf("Error loading .env.local file: %v", err)
	}

	return &Config{
		PSNRefreshToken: os.Getenv("PSN_REFRESH_TOKEN"),
		PSNAccountID:    os.Getenv("PSN_ACCOUNT_ID"),
	}
}
