package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/AelcioJozias/vibe-invest/backend/internal/shared/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) List(ctx context.Context, search string) ([]Response, error) {
	const query = `
		SELECT a.id, a.name, COALESCE(SUM(i.current_amount_cents), 0) AS amount, a.is_active
		FROM accounts a
		LEFT JOIN investments i
		  ON i.account_id = a.id
		 AND i.is_active = TRUE
		WHERE a.is_active = TRUE
		  AND ($1 = '' OR a.name ILIKE '%' || $1 || '%')
		GROUP BY a.id, a.name, a.is_active
		ORDER BY a.id`

	rows, err := r.db.Query(ctx, query, search)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	defer rows.Close()

	accounts := make([]Response, 0)
	for rows.Next() {
		var account Response
		if err := rows.Scan(&account.ID, &account.Name, &account.Amount, &account.IsActive); err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate accounts: %w", err)
	}

	return accounts, nil
}

func (r *PostgresRepository) Create(ctx context.Context, name string) (Response, error) {
	const query = `
		INSERT INTO accounts (name, is_active, created_at, updated_at)
		VALUES ($1, TRUE, NOW(), NOW())
		RETURNING id, name, is_active`

	var account Response
	if err := r.db.QueryRow(ctx, query, name).Scan(&account.ID, &account.Name, &account.IsActive); err != nil {
		return Response{}, fmt.Errorf("create account: %w", err)
	}

	account.Amount = 0
	return account, nil
}

func (r *PostgresRepository) UpdateName(ctx context.Context, id int64, name string) (Response, error) {
	const query = `
		UPDATE accounts
		SET name = $2, updated_at = NOW()
		WHERE id = $1 AND is_active = TRUE
		RETURNING id, name, is_active`

	var account Response
	if err := r.db.QueryRow(ctx, query, id, name).Scan(&account.ID, &account.Name, &account.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Response{}, apperrors.ErrNotFound
		}
		return Response{}, fmt.Errorf("update account: %w", err)
	}

	amount, err := r.totalAmountByAccount(ctx, id)
	if err != nil {
		return Response{}, err
	}

	account.Amount = amount
	return account, nil
}

func (r *PostgresRepository) Deactivate(ctx context.Context, id int64) error {
	const query = `
		UPDATE accounts
		SET is_active = FALSE, updated_at = NOW()
		WHERE id = $1 AND is_active = TRUE`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deactivate account: %w", err)
	}

	if result.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *PostgresRepository) ExistsActive(ctx context.Context, id int64) (bool, error) {
	const query = `SELECT 1 FROM accounts WHERE id = $1 AND is_active = TRUE`

	var marker int
	err := r.db.QueryRow(ctx, query, id).Scan(&marker)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	return false, fmt.Errorf("check active account: %w", err)
}

func (r *PostgresRepository) totalAmountByAccount(ctx context.Context, accountID int64) (int64, error) {
	const query = `
		SELECT COALESCE(SUM(current_amount_cents), 0)
		FROM investments
		WHERE account_id = $1 AND is_active = TRUE`

	var amount int64
	if err := r.db.QueryRow(ctx, query, accountID).Scan(&amount); err != nil {
		return 0, fmt.Errorf("load account amount: %w", err)
	}

	return amount, nil
}
