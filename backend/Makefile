.PHONY: run test migrate-up migrate-down

# Start the API server (reads .env automatically)
run:
	go run ./cmd/api

# Run all unit tests
test:
	go test ./...

# Apply all pending SQL migrations
# Equivalent to Flyway's migrate goal in Spring Boot
migrate-up:
	go run ./cmd/migrate up

# Roll back all applied migrations
migrate-down:
	go run ./cmd/migrate down
