-- +goose Up
CREATE TABLE transactions (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    holding_id  UUID NOT NULL REFERENCES holdings(id) ON DELETE CASCADE,
    type        VARCHAR(10) NOT NULL CHECK (type IN ('buy','sell')),
    quantity    DOUBLE PRECISION NOT NULL,
    price       DOUBLE PRECISION NOT NULL,
    fee         DOUBLE PRECISION NOT NULL DEFAULT 0,
    notes       TEXT NOT NULL DEFAULT '',
    executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_holding_id ON transactions(holding_id);
CREATE INDEX idx_transactions_executed_at ON transactions(executed_at);

-- +goose Down
DROP TABLE IF EXISTS transactions;
