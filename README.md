# Vibe Invest

Personal project to learn backend development with Go while building an investment tracking API.

## Project Structure

- backend: Go backend implementation
- backend-contracts: OpenAPI contracts

## Current Status

- API contracts are being modeled first
- Backend architecture and business rules are documented
- Database model is defined in Mermaid

## Docs

Main docs are in backend/docs:

- architecture.md
- technical-spec.md
- database-model.md
- dashboard-calculation-rules.md
- business-rules/

Project-level docs are in docs:

- commit-convention.md

## Commit Pattern

This repository uses Conventional Commits as the required commit message standard.

See docs/commit-convention.md for rules and examples.

## Next Steps

- Create SQL migrations
- Bootstrap HTTP server and database connection
- Implement accounts, investments, and portfolio summary endpoints
