# Backend - Vibe Invest

Backend API built with Go (standard library + PostgreSQL), following the contracts in backend-contracts.

## Stack

- Go 1.26+
- net/http (no framework)
- PostgreSQL
- pgx / pgxpool

## Project Structure

- cmd/api: executable entrypoint
- internal/account: account feature (handler, service, repository, model)
- internal/investment: investment feature (handler, service, repository, model)
- internal/dashboard: portfolio summary feature
- internal/config: env config loading
- internal/database: postgres connection setup
- internal/shared: shared helpers (errors, http, time)
- migrations: SQL schema migrations

## Environment Variables

- DATABASE_URL: required PostgreSQL connection string
- PORT: optional, default is 8080
- CORS_ALLOW_ORIGIN: optional, empty by default (CORS disabled)

For local development, a `.env` file in the `backend` folder is loaded automatically.

Example (`.env`):

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/app_db?sslmode=disable
PORT=8080
CORS_ALLOW_ORIGIN=http://localhost:8080
```

If you prefer shell variables, use:

```powershell
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/app_db?sslmode=disable"
$env:PORT="8080"
$env:CORS_ALLOW_ORIGIN="http://localhost:8080"
```

## Run

From backend folder:

```powershell
go run .
```

Or:

```powershell
go run ./cmd/api
```

Health check:

```text
GET /health
```

## Database

Run migrations before starting the API (similar to Flyway in Spring, but as an explicit step):

```powershell
make migrate-up
```

Or without make:

```powershell
go run ./cmd/migrate up
```

Roll back:

```powershell
make migrate-down
```

Migration files:

- migrations/0001_init_schema.up.sql
- migrations/0001_init_schema.down.sql

Main conventions:

- soft delete with is_active
- monetary values stored as integer cents (BIGINT)
- movement_type values:
  - INVESTMENT_CREATED
  - CONTRIBUTION
  - INTEREST
  - ADJUSTMENT

## Implemented Endpoints

Accounts:

- GET /api/v1/accounts
- POST /api/v1/accounts
- PUT /api/v1/accounts/{id}
- DELETE /api/v1/accounts/{id}

Investiments:

- GET /api/v1/accounts/{accountId}/investiments
- POST /api/v1/accounts/{accountId}/investiments
- GET /api/v1/investiments/{investmentId}
- PUT /api/v1/investiments/{investmentId}
- DELETE /api/v1/investiments/{investmentId}
- PUT /api/v1/investiments/{investmentId}/fees

Portfolio:

- GET /api/v1/portfolio/summary?referenceMonth=YYYY-MM

## Tests

Run all tests:

```powershell
go test ./...
```

Testing policy is documented in:

- docs/technical-spec.md
- docs/dashboard-calculation-rules.md
