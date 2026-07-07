package postgres

import (
	"context"

	"portinves/internal/domain"
)

type ServiceCheckRepo struct {
	db DBTX
}

func NewServiceCheckRepo(db DBTX) *ServiceCheckRepo {
	return &ServiceCheckRepo{db: db}
}

func (r *ServiceCheckRepo) Create(ctx context.Context, sc *domain.ServiceCheck) error {
	const q = `INSERT INTO service_checks (server_id, name, endpoint, method, expected_status, interval_seconds, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`
	err := r.db.QueryRow(ctx, q, sc.ServerID, sc.Name, sc.Endpoint, sc.Method, sc.ExpectedStatus, sc.IntervalSeconds, sc.IsActive).
		Scan(&sc.ID, &sc.CreatedAt)
	return wrapError(err, "service check")
}

func (r *ServiceCheckRepo) GetByID(ctx context.Context, id string) (*domain.ServiceCheck, error) {
	const q = `SELECT id, server_id, name, endpoint, method, expected_status, interval_seconds, is_active, created_at
		FROM service_checks WHERE id = $1`
	sc := &domain.ServiceCheck{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&sc.ID, &sc.ServerID, &sc.Name, &sc.Endpoint, &sc.Method, &sc.ExpectedStatus, &sc.IntervalSeconds, &sc.IsActive, &sc.CreatedAt,
	)
	if err != nil {
		return nil, wrapError(err, "service check")
	}
	return sc, nil
}

func (r *ServiceCheckRepo) ListByServerID(ctx context.Context, serverID string) ([]domain.ServiceCheck, error) {
	const q = `SELECT id, server_id, name, endpoint, method, expected_status, interval_seconds, is_active, created_at
		FROM service_checks WHERE server_id = $1 ORDER BY name ASC`
	rows, err := r.db.Query(ctx, q, serverID)
	if err != nil {
		return nil, wrapError(err, "service check")
	}
	defer rows.Close()

	var checks []domain.ServiceCheck
	for rows.Next() {
		var sc domain.ServiceCheck
		if err := rows.Scan(
			&sc.ID, &sc.ServerID, &sc.Name, &sc.Endpoint, &sc.Method, &sc.ExpectedStatus, &sc.IntervalSeconds, &sc.IsActive, &sc.CreatedAt,
		); err != nil {
			return nil, wrapError(err, "service check")
		}
		checks = append(checks, sc)
	}
	return checks, rows.Err()
}

func (r *ServiceCheckRepo) ListActive(ctx context.Context) ([]domain.ServiceCheck, error) {
	const q = `SELECT id, server_id, name, endpoint, method, expected_status, interval_seconds, is_active, created_at
		FROM service_checks WHERE is_active = true`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, wrapError(err, "service check")
	}
	defer rows.Close()

	var checks []domain.ServiceCheck
	for rows.Next() {
		var sc domain.ServiceCheck
		if err := rows.Scan(
			&sc.ID, &sc.ServerID, &sc.Name, &sc.Endpoint, &sc.Method, &sc.ExpectedStatus, &sc.IntervalSeconds, &sc.IsActive, &sc.CreatedAt,
		); err != nil {
			return nil, wrapError(err, "service check")
		}
		checks = append(checks, sc)
	}
	return checks, rows.Err()
}
