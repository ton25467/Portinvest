-- name: CreateUptimeLog :one
INSERT INTO uptime_logs (service_check_id, status, status_code, response_time_ms, error_message, checked_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, service_check_id, status, status_code, response_time_ms, error_message, checked_at;

-- name: ListUptimeLogsByServiceCheckID :many
SELECT id, service_check_id, status, status_code, response_time_ms, error_message, checked_at
FROM uptime_logs WHERE service_check_id = $1
ORDER BY checked_at DESC LIMIT $2;

-- name: GetUptimeByServerID :many
SELECT ul.id, ul.service_check_id, ul.status, ul.status_code, ul.response_time_ms, ul.error_message, ul.checked_at
FROM uptime_logs ul
JOIN service_checks sc ON sc.id = ul.service_check_id
WHERE sc.server_id = $1
ORDER BY ul.checked_at DESC LIMIT $2;
