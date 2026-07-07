package postgres

import (
	"context"

	"portinves/internal/domain"
)

type UptimeLogRepo struct {
	db DBTX
}

func NewUptimeLogRepo(db DBTX) *UptimeLogRepo {
	return &UptimeLogRepo{db: db}
}

func (r *UptimeLogRepo) Create(ctx context.Context, ul *domain.UptimeLog) error {
	const q = `INSERT INTO uptime_logs (service_check_id, status, status_code, response_time_ms, error_message, checked_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	err := r.db.QueryRow(ctx, q, ul.ServiceCheckID, ul.Status, ul.StatusCode, ul.ResponseTimeMs, ul.ErrorMessage, ul.CheckedAt).
		Scan(&ul.ID)
	return wrapError(err, "uptime log")
}

func (r *UptimeLogRepo) ListByServiceCheckID(ctx context.Context, serviceCheckID string, limit int) ([]domain.UptimeLog, error) {
	const q = `SELECT id, service_check_id, status, status_code, response_time_ms, error_message, checked_at
		FROM uptime_logs WHERE service_check_id = $1 ORDER BY checked_at DESC LIMIT $2`
	rows, err := r.db.Query(ctx, q, serviceCheckID, limit)
	if err != nil {
		return nil, wrapError(err, "uptime log")
	}
	defer rows.Close()

	var logs []domain.UptimeLog
	for rows.Next() {
		var ul domain.UptimeLog
		if err := rows.Scan(
			&ul.ID, &ul.ServiceCheckID, &ul.Status, &ul.StatusCode, &ul.ResponseTimeMs, &ul.ErrorMessage, &ul.CheckedAt,
		); err != nil {
			return nil, wrapError(err, "uptime log")
		}
		logs = append(logs, ul)
	}
	return logs, rows.Err()
}

func (r *UptimeLogRepo) GetUptimeByServerID(ctx context.Context, serverID string, limit int) ([]domain.UptimeLog, error) {
	const q = `SELECT ul.id, ul.service_check_id, ul.status, ul.status_code, ul.response_time_ms, ul.error_message, ul.checked_at
		FROM uptime_logs ul
		JOIN service_checks sc ON ul.service_check_id = sc.id
		WHERE sc.server_id = $1
		ORDER BY ul.checked_at DESC
		LIMIT $2`
	rows, err := r.db.Query(ctx, q, serverID, limit)
	if err != nil {
		return nil, wrapError(err, "uptime log")
	}
	defer rows.Close()

	var logs []domain.UptimeLog
	for rows.Next() {
		var ul domain.UptimeLog
		if err := rows.Scan(
			&ul.ID, &ul.ServiceCheckID, &ul.Status, &ul.StatusCode, &ul.ResponseTimeMs, &ul.ErrorMessage, &ul.CheckedAt,
		); err != nil {
			return nil, wrapError(err, "uptime log")
		}
		logs = append(logs, ul)
	}
	return logs, rows.Err()
}
