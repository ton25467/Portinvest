package postgres

import (
	"context"

	"portinves/internal/domain"
)

// HoldingRepo implements domain.HoldingRepository using PostgreSQL.
type HoldingRepo struct {
	db DBTX
}

// NewHoldingRepo creates a new HoldingRepo.
func NewHoldingRepo(db DBTX) *HoldingRepo {
	return &HoldingRepo{db: db}
}

func (r *HoldingRepo) Create(ctx context.Context, h *domain.Holding) error {
	const q = `INSERT INTO holdings (portfolio_id, symbol, name, asset_type, quantity, avg_buy_price, current_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, q,
		h.PortfolioID, h.Symbol, h.Name, h.AssetType, h.Quantity, h.AvgBuyPrice, h.CurrentPrice,
	).Scan(&h.ID, &h.CreatedAt, &h.UpdatedAt)
	return wrapError(err, "holding")
}

func (r *HoldingRepo) GetByID(ctx context.Context, id string) (*domain.Holding, error) {
	const q = `SELECT id, portfolio_id, symbol, name, asset_type, quantity, avg_buy_price, current_price, created_at, updated_at
		FROM holdings WHERE id = $1`
	h := &domain.Holding{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&h.ID, &h.PortfolioID, &h.Symbol, &h.Name, &h.AssetType,
		&h.Quantity, &h.AvgBuyPrice, &h.CurrentPrice, &h.CreatedAt, &h.UpdatedAt,
	)
	if err != nil {
		return nil, wrapError(err, "holding")
	}
	return h, nil
}

func (r *HoldingRepo) ListByPortfolioID(ctx context.Context, portfolioID string) ([]domain.Holding, error) {
	const q = `SELECT id, portfolio_id, symbol, name, asset_type, quantity, avg_buy_price, current_price, created_at, updated_at
		FROM holdings WHERE portfolio_id = $1 ORDER BY symbol ASC`
	rows, err := r.db.Query(ctx, q, portfolioID)
	if err != nil {
		return nil, wrapError(err, "holding")
	}
	defer rows.Close()

	var holdings []domain.Holding
	for rows.Next() {
		var h domain.Holding
		if err := rows.Scan(
			&h.ID, &h.PortfolioID, &h.Symbol, &h.Name, &h.AssetType,
			&h.Quantity, &h.AvgBuyPrice, &h.CurrentPrice, &h.CreatedAt, &h.UpdatedAt,
		); err != nil {
			return nil, wrapError(err, "holding")
		}
		holdings = append(holdings, h)
	}
	return holdings, rows.Err()
}

func (r *HoldingRepo) Update(ctx context.Context, h *domain.Holding) error {
	const q = `UPDATE holdings SET symbol = $2, name = $3, asset_type = $4, quantity = $5,
		avg_buy_price = $6, current_price = $7, updated_at = NOW()
		WHERE id = $1`
	cmd, err := r.db.Exec(ctx, q,
		h.ID, h.Symbol, h.Name, h.AssetType, h.Quantity, h.AvgBuyPrice, h.CurrentPrice,
	)
	if err != nil {
		return wrapError(err, "holding")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("holding")
	}
	return nil
}

func (r *HoldingRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM holdings WHERE id = $1`
	cmd, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return wrapError(err, "holding")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("holding")
	}
	return nil
}
