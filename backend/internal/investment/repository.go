package investment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/apperrors"
	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/timeutil"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) ListByAccount(ctx context.Context, accountID int64) ([]record, error) {
	const query = `
		SELECT id, account_id, current_amount_cents, yield_rate, observation, is_active
		FROM investments
		WHERE account_id = $1 AND is_active = TRUE
		ORDER BY id`

	rows, err := r.db.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("list investments by account: %w", err)
	}
	defer rows.Close()

	items := make([]record, 0)
	for rows.Next() {
		var item record
		if err := rows.Scan(&item.ID, &item.AccountID, &item.Amount, &item.YieldRate, &item.Observation, &item.IsActive); err != nil {
			return nil, fmt.Errorf("scan investment: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate investments: %w", err)
	}

	return items, nil
}

func (r *PostgresRepository) Create(ctx context.Context, accountID int64, input CreateRequest, referenceMonth time.Time) (record, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return record{}, fmt.Errorf("begin create investment transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := r.ensureActiveAccount(ctx, tx, accountID); err != nil {
		return record{}, err
	}

	const insertInvestment = `
		INSERT INTO investments (account_id, current_amount_cents, yield_rate, observation, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, TRUE, NOW(), NOW())
		RETURNING id, account_id, current_amount_cents, yield_rate, observation, is_active`

	var item record
	if err := tx.QueryRow(ctx, insertInvestment, accountID, input.Amount, input.YieldRate, input.Observation).Scan(
		&item.ID,
		&item.AccountID,
		&item.Amount,
		&item.YieldRate,
		&item.Observation,
		&item.IsActive,
	); err != nil {
		return record{}, fmt.Errorf("insert investment: %w", err)
	}

	if err := r.applyMovement(ctx, tx, item.ID, MovementTypeInvestmentCreated, input.Amount, input.Observation, referenceMonth); err != nil {
		return record{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return record{}, fmt.Errorf("commit create investment transaction: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, investmentID int64) (record, error) {
	const query = `
		SELECT id, account_id, current_amount_cents, yield_rate, observation, is_active
		FROM investments
		WHERE id = $1 AND is_active = TRUE`

	var item record
	if err := r.db.QueryRow(ctx, query, investmentID).Scan(
		&item.ID,
		&item.AccountID,
		&item.Amount,
		&item.YieldRate,
		&item.Observation,
		&item.IsActive,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return record{}, apperrors.ErrNotFound
		}
		return record{}, fmt.Errorf("get investment by id: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) Update(ctx context.Context, investmentID int64, input UpdateRequest, referenceMonth time.Time) (record, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return record{}, fmt.Errorf("begin update investment transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const currentQuery = `
		SELECT account_id, current_amount_cents
		FROM investments
		WHERE id = $1 AND is_active = TRUE
		FOR UPDATE`

	var currentAccountID int64
	var currentAmount int64
	if err := tx.QueryRow(ctx, currentQuery, investmentID).Scan(&currentAccountID, &currentAmount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return record{}, apperrors.ErrNotFound
		}
		return record{}, fmt.Errorf("load current investment state: %w", err)
	}

	const updateQuery = `
		UPDATE investments
		SET current_amount_cents = $2,
		    yield_rate = $3,
		    observation = $4,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, account_id, current_amount_cents, yield_rate, observation, is_active`

	var item record
	if err := tx.QueryRow(ctx, updateQuery, investmentID, input.Amount, input.YieldRate, input.Observation).Scan(
		&item.ID,
		&item.AccountID,
		&item.Amount,
		&item.YieldRate,
		&item.Observation,
		&item.IsActive,
	); err != nil {
		return record{}, fmt.Errorf("update investment: %w", err)
	}

	delta := input.Amount - currentAmount
	if delta != 0 {
		if err := r.applyMovement(ctx, tx, item.ID, MovementTypeAdjustment, delta, input.Observation, referenceMonth); err != nil {
			return record{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return record{}, fmt.Errorf("commit update investment transaction: %w", err)
	}

	_ = currentAccountID
	return item, nil
}

func (r *PostgresRepository) Deactivate(ctx context.Context, investmentID int64) error {
	const query = `
		UPDATE investments
		SET is_active = FALSE, updated_at = NOW()
		WHERE id = $1 AND is_active = TRUE`

	result, err := r.db.Exec(ctx, query, investmentID)
	if err != nil {
		return fmt.Errorf("deactivate investment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *PostgresRepository) IncrementFees(ctx context.Context, investmentID int64, amount int64, referenceMonth time.Time) (record, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return record{}, fmt.Errorf("begin increment fees transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const lockQuery = `
		SELECT current_amount_cents
		FROM investments
		WHERE id = $1 AND is_active = TRUE
		FOR UPDATE`

	var currentAmount int64
	if err := tx.QueryRow(ctx, lockQuery, investmentID).Scan(&currentAmount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return record{}, apperrors.ErrNotFound
		}
		return record{}, fmt.Errorf("lock investment for fee increment: %w", err)
	}

	updatedAmount := currentAmount + amount

	const updateQuery = `
		UPDATE investments
		SET current_amount_cents = $2,
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, account_id, current_amount_cents, yield_rate, observation, is_active`

	var item record
	if err := tx.QueryRow(ctx, updateQuery, investmentID, updatedAmount).Scan(
		&item.ID,
		&item.AccountID,
		&item.Amount,
		&item.YieldRate,
		&item.Observation,
		&item.IsActive,
	); err != nil {
		return record{}, fmt.Errorf("update investment amount after fees: %w", err)
	}

	if err := r.applyMovement(ctx, tx, item.ID, MovementTypeInterest, amount, "fees increment", referenceMonth); err != nil {
		return record{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return record{}, fmt.Errorf("commit increment fees transaction: %w", err)
	}

	return item, nil
}

func (r *PostgresRepository) ensureActiveAccount(ctx context.Context, tx pgx.Tx, accountID int64) error {
	const query = `SELECT 1 FROM accounts WHERE id = $1 AND is_active = TRUE`

	var marker int
	err := tx.QueryRow(ctx, query, accountID).Scan(&marker)
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrNotFound
	}

	return fmt.Errorf("ensure active account: %w", err)
}

func (r *PostgresRepository) applyMovement(ctx context.Context, tx pgx.Tx, investmentID int64, movementType string, amount int64, observation string, when time.Time) error {
	referenceMonth := timeutil.CurrentReferenceMonth(when)
	if err := r.ensureMonthlyBalance(ctx, tx, investmentID, referenceMonth); err != nil {
		return err
	}

	const insertMovement = `
		INSERT INTO investment_movements (
			investment_id,
			reference_month,
			movement_type,
			amount_cents,
			observation,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, NOW())`

	if _, err := tx.Exec(ctx, insertMovement, investmentID, referenceMonth, movementType, amount, observation); err != nil {
		return fmt.Errorf("insert investment movement: %w", err)
	}

	var updateBalanceQuery string
	switch movementType {
	case MovementTypeContribution:
		updateBalanceQuery = `
			UPDATE investment_monthly_balances
			SET contribution_amount_cents = contribution_amount_cents + $3,
			    closing_amount_cents = closing_amount_cents + $3,
			    updated_at = NOW()
			WHERE investment_id = $1 AND reference_month = $2`
	case MovementTypeInterest:
		updateBalanceQuery = `
			UPDATE investment_monthly_balances
			SET interest_amount_cents = interest_amount_cents + $3,
			    closing_amount_cents = closing_amount_cents + $3,
			    updated_at = NOW()
			WHERE investment_id = $1 AND reference_month = $2`
	case MovementTypeAdjustment:
		updateBalanceQuery = `
			UPDATE investment_monthly_balances
			SET adjustment_amount_cents = adjustment_amount_cents + $3,
			    closing_amount_cents = closing_amount_cents + $3,
			    updated_at = NOW()
			WHERE investment_id = $1 AND reference_month = $2`
	default:
		updateBalanceQuery = `
			UPDATE investment_monthly_balances
			SET closing_amount_cents = closing_amount_cents + $3,
			    updated_at = NOW()
			WHERE investment_id = $1 AND reference_month = $2`
	}

	if _, err := tx.Exec(ctx, updateBalanceQuery, investmentID, referenceMonth, amount); err != nil {
		return fmt.Errorf("update monthly balance after movement: %w", err)
	}

	return nil
}

func (r *PostgresRepository) ensureMonthlyBalance(ctx context.Context, tx pgx.Tx, investmentID int64, referenceMonth time.Time) error {
	previousMonth := timeutil.PreviousMonth(referenceMonth)

	const query = `
		INSERT INTO investment_monthly_balances (
			investment_id,
			reference_month,
			opening_amount_cents,
			closing_amount_cents,
			contribution_amount_cents,
			interest_amount_cents,
			adjustment_amount_cents,
			created_at,
			updated_at
		)
		VALUES (
			$1,
			$2,
			COALESCE((
				SELECT closing_amount_cents
				FROM investment_monthly_balances
				WHERE investment_id = $1 AND reference_month = $3
			), 0),
			COALESCE((
				SELECT closing_amount_cents
				FROM investment_monthly_balances
				WHERE investment_id = $1 AND reference_month = $3
			), 0),
			0,
			0,
			0,
			NOW(),
			NOW()
		)
		ON CONFLICT (investment_id, reference_month) DO NOTHING`

	if _, err := tx.Exec(ctx, query, investmentID, referenceMonth, previousMonth); err != nil {
		return fmt.Errorf("ensure monthly balance row: %w", err)
	}

	return nil
}
