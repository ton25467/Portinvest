-- name: CreateTransaction :one
INSERT INTO transactions (holding_id, type, quantity, price, fee, notes, executed_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, holding_id, type, quantity, price, fee, notes, executed_at, created_at;

-- name: ListTransactionsByHoldingID :many
SELECT id, holding_id, type, quantity, price, fee, notes, executed_at, created_at
FROM transactions WHERE holding_id = $1 ORDER BY executed_at DESC;
