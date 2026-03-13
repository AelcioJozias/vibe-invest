# Dashboard Calculation Rules

## Purpose

Define deterministic calculation rules for GET /api/v1/portfolio/summary.

This document removes ambiguity between business intent and SQL implementation.

## Input

- referenceMonth (optional): YYYY-MM

Behavior:

- If referenceMonth is omitted, use current month reference.
- If referenceMonth is invalid, return 400.

## Data Sources

- investment_monthly_balances
- investment_movements (support/audit source)
- investments (active filter)
- accounts (active filter)

## Active Resource Filters

Default dashboard considers only:

- accounts where is_active = true
- investments where is_active = true

## Metric Definitions

For the selected month M:

- totalInvestedAmount
  - Sum of closing_amount_cents for month M

- totalMonthlyYieldAmount
  - Sum of interest_amount_cents for month M
  - Equivalent business rule: only INTEREST enters yield amount

- totalMonthlyContributions
  - Sum of contribution_amount_cents for month M

- previousMonthTotalAmount
  - Sum of opening_amount_cents for month M
  - Equivalent to previous month carry-over total

- portfolioGrowthAmount
  - Sum of (closing_amount_cents - opening_amount_cents) for month M

- averageMonthlyYieldRate
  - Average yield rate percentage for month M across considered investments
  - Computation strategy in phase 1:
    - per investment rate = (interest_amount_cents / NULLIF(opening_amount_cents, 0)) * 100
    - global average = average of per investment rate values where opening_amount_cents > 0

- topYieldInvestment
  - Investment with max interest_amount_cents in month M
  - Tie-breaker: lower investment_id first

## First Month Rule

When there is no previous month baseline:

- opening_amount_cents can be 0
- previousMonthTotalAmount can be 0
- rates requiring division by opening amount should ignore zero denominators
- dashboard must still return valid payload with numeric defaults

## Movement Type Rules

Allowed values:

- INVESTMENT_CREATED
- CONTRIBUTION
- INTEREST
- ADJUSTMENT

Calculation implications:

- CONTRIBUTION does not count as yield
- INTEREST counts as yield
- ADJUSTMENT is operational correction and does not count as yield by default

## Numerical Representation

All amount fields in API response are integer cents.

Examples:

- 120.00 -> 12000
- 400.00 -> 40000

## Determinism and Recalculation

- Monthly balances are treated as source for summary metrics.
- If a past movement is corrected, the affected monthly balances must be recalculated before querying dashboard.

## Validation Checklist

For each dashboard query, ensure:

1. month format normalized
2. active filters applied
3. sums computed from *_cents columns
4. INTEREST separated from CONTRIBUTION and ADJUSTMENT
5. zero-division guarded in rate calculations
6. deterministic tie-break for topYieldInvestment

## Example

Given month M:

- Sum closing_amount_cents = 10057820
- Sum interest_amount_cents = 12000
- Sum contribution_amount_cents = 40000
- Sum opening_amount_cents = 9957820

Then:

- totalInvestedAmount = 10057820
- totalMonthlyYieldAmount = 12000
- totalMonthlyContributions = 40000
- previousMonthTotalAmount = 9957820
- portfolioGrowthAmount = 600000

Note:

- portfolioGrowthAmount includes all net change between opening and closing values.
- totalMonthlyYieldAmount remains isolated to INTEREST only.
