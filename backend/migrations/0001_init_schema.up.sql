CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(120) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE investments (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id),
    current_amount_cents BIGINT NOT NULL DEFAULT 0,
    yield_rate VARCHAR(60) NOT NULL,
    observation TEXT NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE investment_movements (
    id BIGSERIAL PRIMARY KEY,
    investment_id BIGINT NOT NULL REFERENCES investments(id),
    reference_month DATE NOT NULL,
    movement_type VARCHAR(32) NOT NULL,
    amount_cents BIGINT NOT NULL,
    observation TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_investment_movements_type
        CHECK (movement_type IN ('INVESTMENT_CREATED', 'CONTRIBUTION', 'INTEREST', 'ADJUSTMENT'))
);

CREATE TABLE investment_monthly_balances (
    id BIGSERIAL PRIMARY KEY,
    investment_id BIGINT NOT NULL REFERENCES investments(id),
    reference_month DATE NOT NULL,
    opening_amount_cents BIGINT NOT NULL DEFAULT 0,
    closing_amount_cents BIGINT NOT NULL DEFAULT 0,
    contribution_amount_cents BIGINT NOT NULL DEFAULT 0,
    interest_amount_cents BIGINT NOT NULL DEFAULT 0,
    adjustment_amount_cents BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_investment_monthly_balances UNIQUE (investment_id, reference_month)
);

CREATE INDEX idx_investments_account_id ON investments (account_id);
CREATE INDEX idx_investments_is_active ON investments (is_active);
CREATE INDEX idx_investment_movements_lookup ON investment_movements (investment_id, reference_month);
CREATE INDEX idx_investment_monthly_balances_reference_month ON investment_monthly_balances (reference_month);
CREATE INDEX idx_investment_monthly_balances_lookup ON investment_monthly_balances (investment_id, reference_month);
