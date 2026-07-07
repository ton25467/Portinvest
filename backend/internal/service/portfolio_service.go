package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"portinves/internal/domain"
	"portinves/internal/repository/postgres"
)

type PortfolioService struct {
	pool          *pgxpool.Pool
	portfolioRepo domain.PortfolioRepository
	holdingRepo   domain.HoldingRepository
	txRepo        domain.TransactionRepository
}

func NewPortfolioService(
	pool *pgxpool.Pool,
	portfolioRepo domain.PortfolioRepository,
	holdingRepo domain.HoldingRepository,
	txRepo domain.TransactionRepository,
) *PortfolioService {
	return &PortfolioService{
		pool:          pool,
		portfolioRepo: portfolioRepo,
		holdingRepo:   holdingRepo,
		txRepo:        txRepo,
	}
}

func (s *PortfolioService) CreatePortfolio(ctx context.Context, userID, name, description, currency string) (*domain.Portfolio, error) {
	if name == "" {
		return nil, domain.ErrBadRequest("portfolio name is required")
	}
	if currency == "" {
		currency = "USD"
	}
	p := &domain.Portfolio{
		UserID:      userID,
		Name:        name,
		Description: description,
		Currency:    currency,
	}
	if err := s.portfolioRepo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PortfolioService) GetPortfolio(ctx context.Context, id string) (*domain.Portfolio, error) {
	return s.portfolioRepo.GetByID(ctx, id)
}

func (s *PortfolioService) ListPortfolios(ctx context.Context, userID string) ([]domain.Portfolio, error) {
	return s.portfolioRepo.ListByUserID(ctx, userID)
}

func (s *PortfolioService) GetPortfolioSummary(ctx context.Context, id string) (*domain.PortfolioSummary, error) {
	return s.portfolioRepo.GetSummary(ctx, id)
}

func (s *PortfolioService) GetHoldings(ctx context.Context, portfolioID string) ([]domain.Holding, error) {
	return s.holdingRepo.ListByPortfolioID(ctx, portfolioID)
}

func (s *PortfolioService) AddTransaction(ctx context.Context, portfolioID string, symbol, name string, assetType domain.AssetType, txType domain.TransactionType, quantity, price, fee float64, notes string, executedAt time.Time) (*domain.Transaction, error) {
	if quantity <= 0 || price <= 0 {
		return nil, domain.ErrBadRequest("quantity and price must be greater than zero")
	}
	if executedAt.IsZero() {
		executedAt = time.Now()
	}

	// Begin DB Transaction
	dbTx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer dbTx.Rollback(ctx)

	// Build repo instances inside dbTx context
	hRepoTx := postgres.NewHoldingRepo(dbTx)
	tRepoTx := postgres.NewTransactionRepo(dbTx)

	// Find if holding already exists in portfolio
	holdings, err := hRepoTx.ListByPortfolioID(ctx, portfolioID)
	if err != nil {
		return nil, err
	}

	var existingHolding *domain.Holding
	for i := range holdings {
		if holdings[i].Symbol == symbol {
			existingHolding = &holdings[i]
			break
		}
	}

	var holdingID string
	if txType == domain.TransactionTypeBuy {
		if existingHolding == nil {
			// Create new holding
			h := &domain.Holding{
				PortfolioID:  portfolioID,
				Symbol:       symbol,
				Name:         name,
				AssetType:    assetType,
				Quantity:     quantity,
				AvgBuyPrice:  price,
				CurrentPrice: price, // Initialize current price as buy price
			}
			if err := hRepoTx.Create(ctx, h); err != nil {
				return nil, err
			}
			holdingID = h.ID
		} else {
			// Update existing holding quantity and avg buy price
			newQty := existingHolding.Quantity + quantity
			newAvgPrice := ((existingHolding.Quantity * existingHolding.AvgBuyPrice) + (quantity * price)) / newQty

			existingHolding.Quantity = newQty
			existingHolding.AvgBuyPrice = newAvgPrice
			existingHolding.CurrentPrice = price // Update current price to latest buy price

			if err := hRepoTx.Update(ctx, existingHolding); err != nil {
				return nil, err
			}
			holdingID = existingHolding.ID
		}
	} else if txType == domain.TransactionTypeSell {
		if existingHolding == nil {
			return nil, domain.ErrBadRequest(fmt.Sprintf("no holdings found for symbol %s to sell", symbol))
		}
		if existingHolding.Quantity < quantity {
			return nil, domain.ErrBadRequest(fmt.Sprintf("insufficient holding quantity: you have %f, but tried to sell %f", existingHolding.Quantity, quantity))
		}

		newQty := existingHolding.Quantity - quantity
		existingHolding.Quantity = newQty
		existingHolding.CurrentPrice = price // Update current price to latest sell price

		if err := hRepoTx.Update(ctx, existingHolding); err != nil {
			return nil, err
		}
		holdingID = existingHolding.ID
	} else {
		return nil, domain.ErrBadRequest("invalid transaction type")
	}

	// Create transaction record
	t := &domain.Transaction{
		HoldingID:  holdingID,
		Type:       txType,
		Quantity:   quantity,
		Price:      price,
		Fee:        fee,
		Notes:      notes,
		ExecutedAt: executedAt,
	}

	if err := tRepoTx.Create(ctx, t); err != nil {
		return nil, err
	}

	// Commit DB Transaction
	if err := dbTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return t, nil
}

func (s *PortfolioService) GetTransactions(ctx context.Context, holdingID string) ([]domain.Transaction, error) {
	return s.txRepo.ListByHoldingID(ctx, holdingID)
}

func (s *PortfolioService) UpdatePrices(ctx context.Context, prices map[string]float64) error {
	// prices: symbol -> currentPrice
	// Begin DB Transaction
	dbTx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer dbTx.Rollback(ctx)

	// Since we don't have a single UpdateCurrentPrice query, let's query all holdings and update them if symbol matches.
	// In production, we would run: UPDATE holdings SET current_price = $1 WHERE symbol = $2
	// Let's implement this efficiently by just running the UPDATE query directly.
	for symbol, price := range prices {
		const q = `UPDATE holdings SET current_price = $1, updated_at = NOW() WHERE symbol = $2`
		if _, err := dbTx.Exec(ctx, q, price, symbol); err != nil {
			return fmt.Errorf("update price for %s: %w", symbol, err)
		}
	}

	return dbTx.Commit(ctx)
}
