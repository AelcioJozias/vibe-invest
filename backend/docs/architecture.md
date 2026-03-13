# Backend Architecture

## Objective

This backend is a personal learning project built with Go, focused on simplicity, clarity, and good engineering practices.

The goal is to build a maintainable API without frameworks, avoiding unnecessary architectural complexity.

## Core Decisions

- Language: Go
- HTTP layer: standard library with `net/http`
- Routing: `http.ServeMux`
- Database: PostgreSQL
- PostgreSQL driver: `pgx` with connection pool
- Migrations: SQL migrations managed by a Go migration tool
- Logging: `log/slog`
- Configuration: environment variables
- Testing: Go standard library with `testing` and `httptest`, plus PostgreSQL integration tests

## What We Will Not Use

- No web framework
- No ORM in the first moment
- No Clean Architecture
- No premature pagination support
- No unnecessary abstraction layers

## Chosen Project Structure

We will organize the backend by feature.

```text
backend/
  cmd/api/main.go
  internal/config/
  internal/database/
  internal/shared/
  internal/account/
    handler.go
    service.go
    repository.go
    model.go
  internal/investment/
    handler.go
    service.go
    repository.go
    model.go
  internal/dashboard/
    handler.go
    service.go
    model.go
  migrations/
  docs/
```

## Why This Structure

This structure keeps each feature grouped in one place.

That makes the code easier to navigate, especially for a small or medium project.

It also helps learning, because each feature shows its full flow in one folder:

- HTTP entry point
- business rules
- data access
- models used by the feature

## Layer Responsibilities

### handler

Responsible for HTTP concerns:

- read path params, query params, and request body
- validate basic input
- call the service layer
- convert result into JSON response
- return correct HTTP status codes

Spring comparison:

This is the closest equivalent to a `@RestController`.

### service

Responsible for business rules and use case orchestration.

Examples:

- validate domain rules
- coordinate multiple repositories
- calculate dashboard values
- decide what can or cannot happen in the flow

Spring comparison:

This is similar to a `@Service`.

### repository

Responsible for persistence access.

Examples:

- run SQL queries
- map rows to Go structs
- persist and retrieve data

Spring comparison:

This plays the role of a repository, but without Spring Data magic. In Go we write the queries and wiring explicitly.

### model

Responsible for feature data structures.

Examples:

- account model
- investment model
- dashboard response model

These models should stay simple and explicit.

## Shared Packages

### internal/config

Will hold application configuration loaded from environment variables.

### internal/database

Will hold PostgreSQL connection setup and database bootstrap code.

### internal/shared

Will hold code reused by multiple features, such as:

- JSON response helpers
- common error payloads
- shared utility functions with clear value

Important rule:

If something is used by only one feature, it should stay inside that feature instead of going to `shared`.

## Coding Principles

We will aim for:

- clean code
- SOLID principles applied with pragmatism
- Object Calisthenics where it improves readability
- explicit dependencies
- low coupling
- high cohesion
- short functions
- clear names
- small interfaces only when necessary

## Practical Go Guidelines

- Prefer simple structs and functions over heavy abstraction
- Prefer explicit wiring in `main.go`
- Keep interfaces near the consumer when they are needed
- Avoid Java-style overengineering
- Avoid generic helper packages with unclear responsibility
- Keep packages focused and small
- Prefer composition over inheritance-style thinking

## Comments Strategy

This project will include targeted comments for learning purposes.

Comments should explain:

- why the code exists
- architectural decisions
- comparison with how the same idea would usually appear in Spring

Comments should not explain obvious syntax.

## Testing Philosophy

Testing is a priority in this project.

Rules:

- TDD is the default workflow for new code.
- BDD is required for business use cases and acceptance behavior.
- Integration tests against PostgreSQL are mandatory for regression safety.
- Every bug fix should include a regression test before the fix is considered done.
- Pull requests should prefer adding tests first, then implementation.

## Initial Development Order

1. Bootstrap HTTP server
2. Add configuration loading
3. Add PostgreSQL connection
4. Create SQL migrations
5. Implement accounts feature
6. Implement investments feature
7. Implement portfolio summary dashboard
8. Implement and keep regression-focused tests (unit, handler, and integration)

## API Direction

The API contract is defined in `backend-contracts`.

The backend implementation should follow the contract first, instead of inventing handlers and payloads ad hoc.

At this moment, the portfolio summary endpoint is modeled as a global aggregated resource.

Example:

```text
GET /api/v1/portfolio/summary
```

## Final Guideline

If a decision makes the project more complex without giving a clear benefit right now, we should avoid it.

The preferred approach is the simplest one that keeps the code clean, testable, and easy to evolve.