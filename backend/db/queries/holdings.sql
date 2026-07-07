-- name: CreateHolding :one
INSERT INTO holdings (portfolio_id, symbol, name, asset_type, quantity, avg_buy_price, current_price)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, portfolio_id, symbol, name, asset_type, quantity, avg_buy_price, current_price, created_at, updated_at;

-- name: GetHoldingByID :one
SELECT id, portfolio_id, symbol, name, asset_type, quantity, avg_buy_price, current_price, created_at, updated_at
FROM holdings WHERE id = $1;

-- name: ListHoldingsByPortfolioID :many
SELECT id, portfolio_id, symbol, name, asset_type, quantity, avg_buy_price, current_price, created_at, updated_at
FROM holdings WHERE portfolio_id = $1 ORDER BY symbol ASC;

-- name: UpdateHolding :exec
UPDATE holdings SET symbol = $2, name = $3, asset_type = $4, quantity = $5,
    avg_buy_price = $6, current_price = $7, updated_at = NOW()
WHERE id = $1;

-- name: DeleteHolding :exec
DELETE FROM holdings WHERE id = $1;
