# Technical Specification

## Scope

This document defines the technical conventions for implementing the backend.

It complements:

- backend/docs/architecture.md
- backend/docs/database-model.md
- backend/docs/business-rules/business-rules-01.md

## Runtime and Stack

- Language: Go
- Go style: standard library first
- HTTP: net/http and http.ServeMux
- Database: PostgreSQL
- Driver: pgx with connection pool
- Logging: log/slog
- Migrations: SQL migration tool (up/down files)

## Architecture and Package Convention

The project uses feature-oriented organization.

- Each feature owns its handler, service, repository, and model.
- Shared code must only exist when reused by multiple features.
- Dependency direction:
  - handler -> service -> repository -> database

No framework and no ORM in this first phase.

## API and Contract Source of Truth

OpenAPI contracts under backend-contracts are the source of truth for:

- endpoint paths
- payload shape
- status codes
- required fields

Implementation must follow contracts first.

## Naming and Language

- Code names and constants are in English.
- Database table and column names use snake_case.
- JSON payload fields use camelCase.
- Soft delete flag name:
  - API: isActive
  - Database: is_active

## Money and Numeric Precision

All monetary values are represented as integer cents.

Examples:

- 1000.12 BRL -> 100012
- 25.75 BRL -> 2575

Rules:

- Do not use float or double for persisted monetary values.
- Database money columns must use BIGINT with *_cents suffix.
- API documentation must describe that values are in cents.

## Soft Delete Convention

Accounts and investments use logical deletion.

Rules:

- Physical delete is not performed by default flows.
- Soft delete updates is_active to false.
- Read queries must filter active records when listing default resources.
- Updates on inactive resources should return not found or conflict based on endpoint rule.

## Date and Month Reference

Dashboard and monthly aggregations use referenceMonth in YYYY-MM format.

Rules:

- Persist movement month as DATE representing the month reference.
- Service layer normalizes month input before query.
- Invalid month format returns 400.

## Movement Types

Allowed movement_type values:

- INVESTMENT_CREATED
- CONTRIBUTION
- INTEREST
- ADJUSTMENT

Rules:

- Store as controlled string value with database constraint.
- Business calculations must separate each type.
- Only INTEREST contributes to yield amount metrics.

## Error Handling Convention

- Use RFC 7807 style payload for handled errors.
- Return deterministic status codes according to OpenAPI.
- Validation errors return 400.
- Missing resource returns 404.
- Conflict scenarios return 409.

## Repository Rules

- Repositories execute explicit SQL.
- SQL must be parameterized.
- No string concatenation with user input.
- Transaction boundaries are controlled by service use cases that require atomicity.

## Testing Strategy

- Unit tests for service business rules.
- Handler tests using httptest.
- Repository tests against PostgreSQL (integration profile).
- Priority tests:
  - movement classification
  - month-over-month calculation
  - soft delete filtering
  - cents conversion and aggregation

## Non-goals for Phase 1

- Pagination
- Authentication and authorization
- Caching layer
- Background jobs
- Distributed architecture concerns

## Decision Update Rule

If a new decision affects contracts, database schema, or business rules, update this file in the same change set.
