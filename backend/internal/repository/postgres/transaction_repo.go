package postgres

import (
	"context"

	"portinves/internal/domain"
)

// TransactionRepo implements domain.TransactionRepository using PostgreSQL.
type TransactionRepo struct {
	db DBTX
}

// NewTransactionRepo creates a new TransactionRepo.
func NewTransactionRepo(db DBTX) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) Create(ctx context.Context, t *domain.Transaction) error {
	const q = `INSERT INTO transactions (holding_id, type, quantity, price, fee, notes, executed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`
	err := r.db.QueryRow(ctx, q,
		t.HoldingID, t.Type, t.Quantity, t.Price, t.Fee, t.Notes, t.ExecutedAt,
	).Scan(&t.ID, &t.CreatedAt)
	return wrapError(err, "transaction")
}

func (r *TransactionRepo) ListByHoldingID(ctx context.Context, holdingID string) ([]domain.Transaction, error) {
	const q = `SELECT id, holding_id, type, quantity, price, fee, notes, executed_at, created_at
		FROM transactions WHERE holding_id = $1 ORDER BY executed_at DESC`
	rows, err := r.db.Query(ctx, q, holdingID)
	if err != nil {
		return nil, wrapError(err, "transaction")
	}
	defer rows.Close()

	var txns []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(
			&t.ID, &t.HoldingID, &t.Type, &t.Quantity, &t.Price, &t.Fee, &t.Notes, &t.ExecutedAt, &t.CreatedAt,
		); err != nil {
			return nil, wrapError(err, "transaction")
		}
		txns = append(txns, t)
	}
	return txns, rows.Err()
}
