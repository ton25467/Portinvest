-- name: CreateCashflow :exec
INSERT INTO cashflows (
    id, user_id, portfolio_id, type, amount, currency, description, executed_at, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
);

-- name: ListCashflowsByUserID :many
SELECT id, user_id, portfolio_id, type, amount, currency, description, executed_at, created_at 
FROM cashflows
WHERE user_id = $1
ORDER BY executed_at DESC;

-- name: DeleteCashflow :exec
DELETE FROM cashflows WHERE id = $1;
