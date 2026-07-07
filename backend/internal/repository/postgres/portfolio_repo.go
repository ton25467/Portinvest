package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"portinves/internal/domain"
)

// PortfolioRepo implements domain.PortfolioRepository using PostgreSQL.
type PortfolioRepo struct {
	db DBTX
}

// NewPortfolioRepo creates a new PortfolioRepo.
func NewPortfolioRepo(db DBTX) *PortfolioRepo {
	return &PortfolioRepo{db: db}
}

func (r *PortfolioRepo) Create(ctx context.Context, p *domain.Portfolio) error {
	const q = `INSERT INTO portfolios (user_id, name, description, currency)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, q, p.UserID, p.Name, p.Description, p.Currency).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	return wrapError(err, "portfolio")
}

func (r *PortfolioRepo) GetByID(ctx context.Context, id string) (*domain.Portfolio, error) {
	const q = `SELECT id, user_id, name, description, currency, created_at, updated_at
		FROM portfolios WHERE id = $1`
	p := &domain.Portfolio{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&p.ID, &p.UserID, &p.Name, &p.Description, &p.Currency, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, wrapError(err, "portfolio")
	}
	return p, nil
}

func (r *PortfolioRepo) ListByUserID(ctx context.Context, userID string) ([]domain.Portfolio, error) {
	const q = `SELECT id, user_id, name, description, currency, created_at, updated_at
		FROM portfolios WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, wrapError(err, "portfolio")
	}
	defer rows.Close()

	var portfolios []domain.Portfolio
	for rows.Next() {
		var p domain.Portfolio
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.Currency, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, wrapError(err, "portfolio")
		}
		portfolios = append(portfolios, p)
	}
	return portfolios, rows.Err()
}

func (r *PortfolioRepo) Update(ctx context.Context, p *domain.Portfolio) error {
	const q = `UPDATE portfolios SET name = $2, description = $3, currency = $4, updated_at = NOW()
		WHERE id = $1`
	cmd, err := r.db.Exec(ctx, q, p.ID, p.Name, p.Description, p.Currency)
	if err != nil {
		return wrapError(err, "portfolio")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("portfolio")
	}
	return nil
}

func (r *PortfolioRepo) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM portfolios WHERE id = $1`
	cmd, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return wrapError(err, "portfolio")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("portfolio")
	}
	return nil
}

func (r *PortfolioRepo) GetSummary(ctx context.Context, id string) (*domain.PortfolioSummary, error) {
	const q = `SELECT
		p.id,
		COALESCE(SUM(h.quantity * h.current_price), 0) AS total_value,
		COALESCE(SUM(h.quantity * h.avg_buy_price), 0) AS total_cost,
		COUNT(h.id)::INTEGER AS holding_count
		FROM portfolios p
		LEFT JOIN holdings h ON h.portfolio_id = p.id
		WHERE p.id = $1
		GROUP BY p.id`

	s := &domain.PortfolioSummary{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&s.PortfolioID, &s.TotalValue, &s.TotalCost, &s.HoldingCount,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound("portfolio")
		}
		return nil, wrapError(err, "portfolio summary")
	}

	// Calculate gain percentage
	if s.TotalCost > 0 {
		s.TotalGainPct = ((s.TotalValue - s.TotalCost) / s.TotalCost) * 100
	}

	return s, nil
}
