// Package domain defines the core business entities and repository interfaces.
// This package has zero external dependencies — it only uses the standard library.
package domain

import (
	"time"
)

// Role represents a user's authorization level.
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleViewer Role = "viewer"
)

// AssetType classifies a financial instrument.
type AssetType string

const (
	AssetTypeStock  AssetType = "stock"
	AssetTypeBond   AssetType = "bond"
	AssetTypeCrypto AssetType = "crypto"
	AssetTypeETF    AssetType = "etf"
)

// TransactionType is either a buy or sell action.
type TransactionType string

const (
	TransactionTypeBuy  TransactionType = "buy"
	TransactionTypeSell TransactionType = "sell"
)

// CashflowType is either income, expense, deposit, or withdrawal.
type CashflowType string

const (
	CashflowTypeIncome     CashflowType = "income"
	CashflowTypeExpense    CashflowType = "expense"
	CashflowTypeDeposit    CashflowType = "deposit"
	CashflowTypeWithdrawal CashflowType = "withdrawal"
)

// ServerType categorises a managed server.
type ServerType string

const (
	ServerTypeWeb   ServerType = "web"
	ServerTypeDB    ServerType = "db"
	ServerTypeAPI   ServerType = "api"
	ServerTypeCache ServerType = "cache"
)

// ServerStatus represents the operational state of a server.
type ServerStatus string

const (
	ServerStatusOnline   ServerStatus = "online"
	ServerStatusOffline  ServerStatus = "offline"
	ServerStatusDegraded ServerStatus = "degraded"
)

// CheckStatus represents the result of a service health check.
type CheckStatus string

const (
	CheckStatusUp       CheckStatus = "up"
	CheckStatusDown     CheckStatus = "down"
	CheckStatusDegraded CheckStatus = "degraded"
)

// HTTPMethod represents allowed HTTP methods for service checks.
type HTTPMethod string

const (
	HTTPMethodGET  HTTPMethod = "GET"
	HTTPMethodPOST HTTPMethod = "POST"
)

// --- Core Entities ---

// User represents an authenticated platform user.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never exposed via JSON
	Name         string    `json:"name"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Portfolio groups financial holdings for a user.
type Portfolio struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Holding represents a single asset position within a portfolio.
type Holding struct {
	ID           string    `json:"id"`
	PortfolioID  string    `json:"portfolio_id"`
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	AssetType    AssetType `json:"asset_type"`
	Quantity     float64   `json:"quantity"`
	AvgBuyPrice  float64   `json:"avg_buy_price"`
	CurrentPrice float64   `json:"current_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Transaction records a buy or sell action on a holding.
type Transaction struct {
	ID         string          `json:"id"`
	HoldingID  string          `json:"holding_id"`
	Type       TransactionType `json:"type"`
	Quantity   float64         `json:"quantity"`
	Price      float64         `json:"price"`
	Fee        float64         `json:"fee"`
	Notes      string          `json:"notes"`
	ExecutedAt time.Time       `json:"executed_at"`
	CreatedAt  time.Time       `json:"created_at"`
}

// Server represents a managed infrastructure server.
type Server struct {
	ID            string       `json:"id"`
	UserID        string       `json:"user_id"`
	Name          string       `json:"name"`
	Host          string       `json:"host"`
	Port          int          `json:"port"`
	Type          ServerType   `json:"type"`
	Status        ServerStatus `json:"status"`
	LastCheckedAt *time.Time   `json:"last_checked_at,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// ServiceCheck defines a periodic health-check against a server endpoint.
type ServiceCheck struct {
	ID              string     `json:"id"`
	ServerID        string     `json:"server_id"`
	Name            string     `json:"name"`
	Endpoint        string     `json:"endpoint"`
	Method          HTTPMethod `json:"method"`
	ExpectedStatus  int        `json:"expected_status"`
	IntervalSeconds int        `json:"interval_seconds"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
}

// UptimeLog records the result of a single service check execution.
type UptimeLog struct {
	ID             string      `json:"id"`
	ServiceCheckID string      `json:"service_check_id"`
	Status         CheckStatus `json:"status"`
	StatusCode     int         `json:"status_code"`
	ResponseTimeMs int         `json:"response_time_ms"`
	ErrorMessage   string      `json:"error_message,omitempty"`
	CheckedAt      time.Time   `json:"checked_at"`
}

// Cashflow records a money movement.
type Cashflow struct {
	ID          string       `json:"id"`
	UserID      string       `json:"user_id"`
	PortfolioID *string      `json:"portfolio_id,omitempty"`
	Type        CashflowType `json:"type"`
	Amount      float64      `json:"amount"`
	Currency    string       `json:"currency"`
	Description string       `json:"description"`
	ExecutedAt  time.Time    `json:"executed_at"`
	CreatedAt   time.Time    `json:"created_at"`
}

// Notification represents an in-app notification.
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type BankCredential struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	BankName          string    `json:"bank_name"`
	Username          string    `json:"username"`
	PasswordEncrypted string    `json:"password_encrypted"` // Will not be sent to frontend
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// --- Aggregate / Summary Types ---

// PortfolioSummary provides computed portfolio statistics.
type PortfolioSummary struct {
	PortfolioID  string  `json:"portfolio_id"`
	TotalValue   float64 `json:"total_value"`
	TotalCost    float64 `json:"total_cost"`
	TotalGainPct float64 `json:"total_gain_pct"`
	HoldingCount int     `json:"holding_count"`
}

// DashboardOverview provides a high-level snapshot for the dashboard.
type DashboardOverview struct {
	TotalPortfolios int     `json:"total_portfolios"`
	TotalValue      float64 `json:"total_value"`
	TotalServers    int     `json:"total_servers"`
	ServersOnline   int     `json:"servers_online"`
	ServersOffline  int     `json:"servers_offline"`
	ServersDegraded int     `json:"servers_degraded"`
}
