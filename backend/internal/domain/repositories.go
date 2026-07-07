package domain

import "context"

// UserRepository defines persistence operations for users.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

// PortfolioRepository defines persistence operations for portfolios.
type PortfolioRepository interface {
	Create(ctx context.Context, portfolio *Portfolio) error
	GetByID(ctx context.Context, id string) (*Portfolio, error)
	ListByUserID(ctx context.Context, userID string) ([]Portfolio, error)
	Update(ctx context.Context, portfolio *Portfolio) error
	Delete(ctx context.Context, id string) error
	GetSummary(ctx context.Context, id string) (*PortfolioSummary, error)
}

// HoldingRepository defines persistence operations for holdings.
type HoldingRepository interface {
	Create(ctx context.Context, holding *Holding) error
	GetByID(ctx context.Context, id string) (*Holding, error)
	ListByPortfolioID(ctx context.Context, portfolioID string) ([]Holding, error)
	Update(ctx context.Context, holding *Holding) error
	Delete(ctx context.Context, id string) error
}

// TransactionRepository defines persistence operations for transactions.
type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	ListByHoldingID(ctx context.Context, holdingID string) ([]Transaction, error)
}

// ServerRepository defines persistence operations for servers.
type ServerRepository interface {
	Create(ctx context.Context, server *Server) error
	GetByID(ctx context.Context, id string) (*Server, error)
	ListByUserID(ctx context.Context, userID string) ([]Server, error)
	Update(ctx context.Context, server *Server) error
	Delete(ctx context.Context, id string) error
}

// ServiceCheckRepository defines persistence operations for service checks.
type ServiceCheckRepository interface {
	Create(ctx context.Context, check *ServiceCheck) error
	GetByID(ctx context.Context, id string) (*ServiceCheck, error)
	ListByServerID(ctx context.Context, serverID string) ([]ServiceCheck, error)
	ListActive(ctx context.Context) ([]ServiceCheck, error)
}

// UptimeLogRepository defines persistence operations for uptime logs.
type UptimeLogRepository interface {
	Create(ctx context.Context, log *UptimeLog) error
	ListByServiceCheckID(ctx context.Context, serviceCheckID string, limit int) ([]UptimeLog, error)
	GetUptimeByServerID(ctx context.Context, serverID string, limit int) ([]UptimeLog, error)
}

// DashboardRepository provides aggregate queries for the dashboard.
type DashboardRepository interface {
	GetOverview(ctx context.Context, userID string) (*DashboardOverview, error)
}

// CashflowRepository defines persistence operations for cashflows.
type CashflowRepository interface {
	Create(ctx context.Context, cashflow *Cashflow) error
	ListByUserID(ctx context.Context, userID string) ([]Cashflow, error)
	Delete(ctx context.Context, id string) error
}

// NotificationRepository defines persistence operations for notifications.
type NotificationRepository interface {
	Create(ctx context.Context, n *Notification) error
	ListByUserID(ctx context.Context, userID string, unreadOnly bool) ([]Notification, error)
	MarkAsRead(ctx context.Context, id, userID string) error
}

type BankCredentialRepository interface {
	Save(ctx context.Context, cred *BankCredential) error
	GetByUserIDAndBank(ctx context.Context, userID, bankName string) (*BankCredential, error)
	Delete(ctx context.Context, userID, bankName string) error
}
