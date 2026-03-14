package dashboard

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Summary(ctx context.Context, referenceMonth time.Time) (Response, error) {
	const summaryQuery = `
		SELECT
			COALESCE(SUM(b.closing_amount_cents), 0) AS total_invested_amount,
			COALESCE(SUM(b.interest_amount_cents), 0) AS total_monthly_yield_amount,
			COALESCE(SUM(b.contribution_amount_cents), 0) AS total_monthly_contributions,
			COALESCE(SUM(b.opening_amount_cents), 0) AS previous_month_total_amount,
			COALESCE(SUM(b.closing_amount_cents - b.opening_amount_cents), 0) AS portfolio_growth_amount,
			COALESCE(
				AVG(
					CASE
						WHEN b.opening_amount_cents > 0 THEN (b.interest_amount_cents::numeric / b.opening_amount_cents::numeric) * 100
						ELSE NULL
					END
				),
				0
			) AS average_monthly_yield_rate
		FROM investment_monthly_balances b
		INNER JOIN investments i ON i.id = b.investment_id
		INNER JOIN accounts a ON a.id = i.account_id
		WHERE b.reference_month = $1
		  AND a.is_active = TRUE
		  AND i.is_active = TRUE`

	var summary Response
	if err := r.db.QueryRow(ctx, summaryQuery, referenceMonth).Scan(
		&summary.TotalInvestedAmount,
		&summary.TotalMonthlyYieldAmount,
		&summary.TotalMonthlyContributions,
		&summary.PreviousMonthTotalAmount,
		&summary.PortfolioGrowthAmount,
		&summary.AverageMonthlyYieldRate,
	); err != nil {
		return Response{}, fmt.Errorf("query dashboard summary: %w", err)
	}

	const topYieldQuery = `
		SELECT
			i.id,
			COALESCE(NULLIF(i.observation, ''), 'Investment #' || i.id::text) AS investment_name,
			b.interest_amount_cents
		FROM investment_monthly_balances b
		INNER JOIN investments i ON i.id = b.investment_id
		INNER JOIN accounts a ON a.id = i.account_id
		WHERE b.reference_month = $1
		  AND a.is_active = TRUE
		  AND i.is_active = TRUE
		ORDER BY b.interest_amount_cents DESC, i.id ASC
		LIMIT 1`

	var top TopYieldInvestment
	err := r.db.QueryRow(ctx, topYieldQuery, referenceMonth).Scan(
		&top.InvestmentID,
		&top.Name,
		&top.YieldAmount,
	)
	if err == nil {
		summary.TopYieldInvestment = &top
		return summary, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return summary, nil
	}

	return Response{}, fmt.Errorf("query top yield investment: %w", err)
}
