package postgres

import (
	"context"

	"portinves/internal/domain"
)

// UserRepo implements domain.UserRepository using PostgreSQL.
type UserRepo struct {
	db DBTX
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(db DBTX) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	const q = `INSERT INTO users (email, password_hash, name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, q, u.Email, u.PasswordHash, u.Name, u.Role).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	return wrapError(err, "user")
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const q = `SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM users WHERE id = $1`
	u := &domain.User{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, wrapError(err, "user")
	}
	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM users WHERE email = $1`
	u := &domain.User{}
	err := r.db.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, wrapError(err, "user")
	}
	return u, nil
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	const q = `UPDATE users SET name = $2, role = $3, updated_at = NOW()
		WHERE id = $1`
	cmd, err := r.db.Exec(ctx, q, u.ID, u.Name, u.Role)
	if err != nil {
		return wrapError(err, "user")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("user")
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM users WHERE id = $1`
	cmd, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return wrapError(err, "user")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("user")
	}
	return nil
}
