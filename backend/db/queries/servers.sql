-- name: CreateServer :one
INSERT INTO servers (user_id, name, host, port, type, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, name, host, port, type, status, last_checked_at, created_at, updated_at;

-- name: GetServerByID :one
SELECT id, user_id, name, host, port, type, status, last_checked_at, created_at, updated_at
FROM servers WHERE id = $1;

-- name: ListServersByUserID :many
SELECT id, user_id, name, host, port, type, status, last_checked_at, created_at, updated_at
FROM servers WHERE user_id = $1 ORDER BY name ASC;

-- name: UpdateServer :exec
UPDATE servers SET name = $2, host = $3, port = $4, type = $5, status = $6,
    last_checked_at = $7, updated_at = NOW()
WHERE id = $1;

-- name: DeleteServer :exec
DELETE FROM servers WHERE id = $1;

-- name: GetDashboardOverview :one
SELECT
    (SELECT COUNT(*)::INTEGER FROM portfolios WHERE portfolios.user_id = $1) AS total_portfolios,
    COALESCE((SELECT SUM(h.quantity * h.current_price) FROM holdings h
        JOIN portfolios p ON p.id = h.portfolio_id WHERE p.user_id = $1), 0)::DOUBLE PRECISION AS total_value,
    (SELECT COUNT(*)::INTEGER FROM servers WHERE servers.user_id = $1) AS total_servers,
    (SELECT COUNT(*)::INTEGER FROM servers WHERE servers.user_id = $1 AND status = 'online') AS servers_online,
    (SELECT COUNT(*)::INTEGER FROM servers WHERE servers.user_id = $1 AND status = 'offline') AS servers_offline,
    (SELECT COUNT(*)::INTEGER FROM servers WHERE servers.user_id = $1 AND status = 'degraded') AS servers_degraded;
