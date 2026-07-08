-- +goose Up
CREATE TABLE cashflows (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    portfolio_id UUID REFERENCES portfolios(id) ON DELETE SET NULL,
    type VARCHAR(50) NOT NULL, -- 'income', 'expense', 'deposit', 'withdrawal'
    amount DECIMAL(19, 4) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    description TEXT NOT NULL,
    executed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE cashflows;
