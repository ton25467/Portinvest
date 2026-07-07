package postgres

import (
	"context"

	"portinves/internal/domain"
)

// ServerRepo implements domain.ServerRepository using PostgreSQL.
type ServerRepo struct {
	db DBTX
}

// NewServerRepo creates a new ServerRepo.
func NewServerRepo(db DBTX) *ServerRepo {
	return &ServerRepo{db: db}
}

func (r *ServerRepo) Create(ctx context.Context, s *domain.Server) error {
	const q = `INSERT INTO servers (user_id, name, host, port, type, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, last_checked_at, created_at, updated_at`
	err := r.db.QueryRow(ctx, q, s.UserID, s.Name, s.Host, s.Port, s.Type, s.Status).
		Scan(&s.ID, &s.LastCheckedAt, &s.CreatedAt, &s.UpdatedAt)
	return wrapError(err, "server")
}

func (r *ServerRepo) GetByID(ctx context.Context, id string) (*domain.Server, error) {
	const q = `SELECT id, user_id, name, host, port, type, status, last_checked_at, created_at, updated_at
		FROM servers WHERE id = $1`
	s := &domain.Server{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&s.ID, &s.UserID, &s.Name, &s.Host, &s.Port, &s.Type,
		&s.Status, &s.LastCheckedAt, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, wrapError(err, "server")
	}
	return s, nil
}

func (r *ServerRepo) ListByUserID(ctx context.Context, userID string) ([]domain.Server, error) {
	const q = `SELECT id, user_id, name, host, port, type, status, last_checked_at, created_at, updated_at
		FROM servers WHERE user_id = $1 ORDER BY name ASC`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, wrapError(err, "server")
	}
	defer rows.Close()

	var servers []domain.Server
	for rows.Next() {
		var s domain.Server
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.Name, &s.Host, &s.Port, &s.Type,
			&s.Status, &s.LastCheckedAt, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, wrapError(err, "server")
		}
		servers = append(servers, s)
	}
	return servers, rows.Err()
}

func (r *ServerRepo) Update(ctx context.Context, s *domain.Server) error {
	const q = `UPDATE servers SET name = $2, host = $3, port = $4, type = $5,
		status = $6, last_checked_at = $7, updated_at = NOW()
		WHERE id = $1`
	cmd, err := r.db.Exec(ctx, q,
		s.ID, s.Name, s.Host, s.Port, s.Type, s.Status, s.LastCheckedAt,
	)
	if err != nil {
		return wrapError(err, "server")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("server")
	}
	return nil
}

func (r *ServerRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM servers WHERE id = $1`
	cmd, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return wrapError(err, "server")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("server")
	}
	return nil
}
