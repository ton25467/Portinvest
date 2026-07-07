package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"portinves/internal/domain"
)

type BankCredentialRepo struct {
	db DBTX
}

func NewBankCredentialRepo(db DBTX) *BankCredentialRepo {
	return &BankCredentialRepo{db: db}
}

func (r *BankCredentialRepo) Save(ctx context.Context, cred *domain.BankCredential) error {
	const q = `
		INSERT INTO bank_credentials (id, user_id, bank_name, username, password_encrypted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, bank_name) DO UPDATE 
		SET username = EXCLUDED.username,
			password_encrypted = EXCLUDED.password_encrypted,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(ctx, q, cred.ID, cred.UserID, cred.BankName, cred.Username, cred.PasswordEncrypted, cred.CreatedAt, cred.UpdatedAt)
	return wrapError(err, "bank_credential")
}

func (r *BankCredentialRepo) GetByUserIDAndBank(ctx context.Context, userID, bankName string) (*domain.BankCredential, error) {
	const q = `
		SELECT id, user_id, bank_name, username, password_encrypted, created_at, updated_at
		FROM bank_credentials 
		WHERE user_id = $1 AND bank_name = $2
	`
	c := &domain.BankCredential{}
	err := r.db.QueryRow(ctx, q, userID, bankName).Scan(
		&c.ID, &c.UserID, &c.BankName, &c.Username, &c.PasswordEncrypted, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound("bank_credential")
		}
		return nil, wrapError(err, "bank_credential")
	}
	return c, nil
}

func (r *BankCredentialRepo) Delete(ctx context.Context, userID, bankName string) error {
	const q = `DELETE FROM bank_credentials WHERE user_id = $1 AND bank_name = $2`
	cmd, err := r.db.Exec(ctx, q, userID, bankName)
	if err != nil {
		return wrapError(err, "bank_credential")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("bank_credential")
	}
	return nil
}
