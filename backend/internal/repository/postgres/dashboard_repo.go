package postgres

import (
	"context"

	"portinves/internal/domain"
)

type DashboardRepo struct {
	db DBTX
}

func NewDashboardRepo(db DBTX) *DashboardRepo {
	return &DashboardRepo{db: db}
}

func (r *DashboardRepo) GetOverview(ctx context.Context, userID string) (*domain.DashboardOverview, error) {
	overview := &domain.DashboardOverview{}

	// 1. Total Portfolios
	const qPortfolios = `SELECT COUNT(id) FROM portfolios WHERE user_id = $1`
	err := r.db.QueryRow(ctx, qPortfolios, userID).Scan(&overview.TotalPortfolios)
	if err != nil {
		return nil, wrapError(err, "dashboard overview: portfolios")
	}

	// 2. Total Portfolio Value
	const qValue = `SELECT COALESCE(SUM(h.quantity * h.current_price), 0)
		FROM holdings h
		JOIN portfolios p ON h.portfolio_id = p.id
		WHERE p.user_id = $1`
	err = r.db.QueryRow(ctx, qValue, userID).Scan(&overview.TotalValue)
	if err != nil {
		return nil, wrapError(err, "dashboard overview: holdings value")
	}

	// 3. Servers status counts
	const qServers = `SELECT
		COUNT(id) as total_servers,
		COUNT(id) FILTER (WHERE status = 'online') as servers_online,
		COUNT(id) FILTER (WHERE status = 'offline') as servers_offline,
		COUNT(id) FILTER (WHERE status = 'degraded') as servers_degraded
		FROM servers
		WHERE user_id = $1`
	err = r.db.QueryRow(ctx, qServers, userID).Scan(
		&overview.TotalServers,
		&overview.ServersOnline,
		&overview.ServersOffline,
		&overview.ServersDegraded,
	)
	if err != nil {
		return nil, wrapError(err, "dashboard overview: servers status")
	}

	return overview, nil
}
