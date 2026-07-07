-- +goose Up
CREATE TABLE holdings (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    portfolio_id  UUID NOT NULL REFERENCES portfolios(id) ON DELETE CASCADE,
    symbol        VARCHAR(20) NOT NULL,
    name          VARCHAR(255) NOT NULL,
    asset_type    VARCHAR(20) NOT NULL CHECK (asset_type IN ('stock','bond','crypto','etf')),
    quantity      DOUBLE PRECISION NOT NULL DEFAULT 0,
    avg_buy_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    current_price DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_holdings_portfolio_id ON holdings(portfolio_id);
CREATE INDEX idx_holdings_symbol ON holdings(symbol);

-- +goose Down
DROP TABLE IF EXISTS holdings;
