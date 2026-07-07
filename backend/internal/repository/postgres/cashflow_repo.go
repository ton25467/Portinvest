package postgres

import (
	"context"

	"portinves/internal/domain"
)

type CashflowRepo struct {
	db DBTX
}

func NewCashflowRepo(db DBTX) *CashflowRepo {
	return &CashflowRepo{db: db}
}

func (r *CashflowRepo) Create(ctx context.Context, cf *domain.Cashflow) error {
	const q = `INSERT INTO cashflows (id, user_id, portfolio_id, type, amount, currency, description, executed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	
	_, err := r.db.Exec(ctx, q, cf.ID, cf.UserID, cf.PortfolioID, string(cf.Type), cf.Amount, cf.Currency, cf.Description, cf.ExecutedAt, cf.CreatedAt)
	return wrapError(err, "cashflow")
}

func (r *CashflowRepo) ListByUserID(ctx context.Context, userID string) ([]domain.Cashflow, error) {
	const q = `SELECT id, user_id, portfolio_id, type, amount, currency, description, executed_at, created_at
		FROM cashflows WHERE user_id = $1 ORDER BY executed_at DESC`
	
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, wrapError(err, "cashflow")
	}
	defer rows.Close()

	var cfs []domain.Cashflow
	for rows.Next() {
		var cf domain.Cashflow
		var t string
		if err := rows.Scan(&cf.ID, &cf.UserID, &cf.PortfolioID, &t, &cf.Amount, &cf.Currency, &cf.Description, &cf.ExecutedAt, &cf.CreatedAt); err != nil {
			return nil, wrapError(err, "cashflow")
		}
		cf.Type = domain.CashflowType(t)
		cfs = append(cfs, cf)
	}
	return cfs, rows.Err()
}

func (r *CashflowRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM cashflows WHERE id = $1`
	cmd, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return wrapError(err, "cashflow")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("cashflow")
	}
	return nil
}
