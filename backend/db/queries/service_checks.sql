-- name: CreateServiceCheck :one
INSERT INTO service_checks (server_id, name, endpoint, method, expected_status, interval_seconds, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, server_id, name, endpoint, method, expected_status, interval_seconds, is_active, created_at;

-- name: GetServiceCheckByID :one
SELECT id, server_id, name, endpoint, method, expected_status, interval_seconds, is_active, created_at
FROM service_checks WHERE id = $1;

-- name: ListServiceChecksByServerID :many
SELECT id, server_id, name, endpoint, method, expected_status, interval_seconds, is_active, created_at
FROM service_checks WHERE server_id = $1 ORDER BY name ASC;

-- name: ListActiveServiceChecks :many
SELECT id, server_id, name, endpoint, method, expected_status, interval_seconds, is_active, created_at
FROM service_checks WHERE is_active = true;
