package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"

	"github.com/AelcioJozias/vibe-invest/backend/migrations"
)

// This command is the Go equivalent of running Flyway or Liquibase as a
// separate step before starting the application.
//
// Usage:
//
//	go run ./cmd/migrate up    → apply all pending migrations
//	go run ./cmd/migrate down  → roll back all migrations
func main() {
	// Load .env for local dev, same as config.Load() does.
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// golang-migrate pgx/v5 driver is registered under the "pgx5" scheme.
	dbURL = strings.Replace(dbURL, "postgres://", "pgx5://", 1)
	dbURL = strings.Replace(dbURL, "postgresql://", "pgx5://", 1)

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalf("load migrations source: %v", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, dbURL)
	if err != nil {
		log.Fatalf("create migrator: %v", err)
	}
	defer m.Close()

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate up: %v", err)
		}
		fmt.Println("migrations applied successfully")
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migrate down: %v", err)
		}
		fmt.Println("migrations rolled back successfully")
	default:
		log.Fatalf("unknown command %q — use 'up' or 'down'", cmd)
	}
}
