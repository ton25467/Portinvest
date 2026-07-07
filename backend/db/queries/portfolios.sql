-- name: CreatePortfolio :one
INSERT INTO portfolios (user_id, name, description, currency)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, name, description, currency, created_at, updated_at;

-- name: GetPortfolioByID :one
SELECT id, user_id, name, description, currency, created_at, updated_at
FROM portfolios WHERE id = $1;

-- name: ListPortfoliosByUserID :many
SELECT id, user_id, name, description, currency, created_at, updated_at
FROM portfolios WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdatePortfolio :exec
UPDATE portfolios SET name = $2, description = $3, currency = $4, updated_at = NOW()
WHERE id = $1;

-- name: DeletePortfolio :exec
DELETE FROM portfolios WHERE id = $1;

-- name: GetPortfolioSummary :one
SELECT
    p.id AS portfolio_id,
    COALESCE(SUM(h.quantity * h.current_price), 0)::DOUBLE PRECISION AS total_value,
    COALESCE(SUM(h.quantity * h.avg_buy_price), 0)::DOUBLE PRECISION AS total_cost,
    COUNT(h.id)::INTEGER AS holding_count
FROM portfolios p
LEFT JOIN holdings h ON h.portfolio_id = p.id
WHERE p.id = $1
GROUP BY p.id;
