package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	DatabaseURL     string
	CORSAllowOrigin string
}

func Load() (Config, error) {
	// Similar to Spring's application-local profile file, .env is loaded for local dev only.
	// Existing process environment variables keep precedence over file values.
	_ = godotenv.Load()

	cfg := Config{
		Port:            envOrDefault("PORT", "8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		CORSAllowOrigin: os.Getenv("CORS_ALLOW_ORIGIN"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	return cfg, nil
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
