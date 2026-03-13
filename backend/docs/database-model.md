# Database Model

This first database model is based on:

- Current OpenAPI contracts for accounts, investments, and portfolio summary
- Business rules for monthly comparison and movement types
- Need to keep implementation simple, explicit, and SQL-friendly in Go

## Mermaid ER Diagram

```mermaid
erDiagram
    ACCOUNTS ||--o{ INVESTMENTS : owns
    INVESTMENTS ||--o{ INVESTMENT_MOVEMENTS : receives
    INVESTMENTS ||--o{ INVESTMENT_MONTHLY_BALANCES : snapshots

    ACCOUNTS {
      BIGINT id PK
      VARCHAR name
      BOOLEAN is_active
      TIMESTAMP created_at
      TIMESTAMP updated_at
    }

    INVESTMENTS {
      BIGINT id PK
      BIGINT account_id FK
      BIGINT current_amount_cents
      VARCHAR yield_rate
      TEXT observation
      BOOLEAN is_active
      TIMESTAMP created_at
      TIMESTAMP updated_at
    }

    INVESTMENT_MOVEMENTS {
      BIGINT id PK
      BIGINT investment_id FK
      DATE reference_month
      VARCHAR movement_type
      BIGINT amount_cents
      TEXT observation
      TIMESTAMP created_at
    }

    INVESTMENT_MONTHLY_BALANCES {
      BIGINT id PK
      BIGINT investment_id FK
      DATE reference_month
      BIGINT opening_amount_cents
      BIGINT closing_amount_cents
      BIGINT contribution_amount_cents
      BIGINT interest_amount_cents
      BIGINT adjustment_amount_cents
      TIMESTAMP created_at
      TIMESTAMP updated_at
    }
```

## Notes for Implementation

- Movement types should be restricted to: INVESTMENT_CREATED, CONTRIBUTION, INTEREST, ADJUSTMENT.
- Use boolean logical delete flag with `is_active` (default true) for accounts and investments.
- Store financial amounts as integer cents to avoid floating-point precision issues.
- Use one row per investment per month in INVESTMENT_MONTHLY_BALANCES.
- Add unique constraint for monthly balance: unique(investment_id, reference_month).
- Add indexes:
  - investments(account_id)
  - investment_movements(investment_id, reference_month)
  - investment_monthly_balances(reference_month)
  - investment_monthly_balances(investment_id, reference_month)

## Why This Model

- Keeps CRUD straightforward for accounts and investments.
- Preserves event history in INVESTMENT_MOVEMENTS.
- Makes dashboard queries simple and fast with monthly snapshots.
- Supports the business rule that interest must be separated from contribution and adjustment.

## Dashboard Derivation (Portfolio Summary)

For a reference month:

- totalInvestedAmount: sum(closing_amount)
- totalMonthlyYieldAmount: sum(interest_amount)
- totalMonthlyContributions: sum(contribution_amount)
- previousMonthTotalAmount: sum(opening_amount)
- portfolioGrowthAmount: sum(closing_amount - opening_amount)
- averageMonthlyYieldRate: average by investment for the month

All amount fields above are represented in cents.
