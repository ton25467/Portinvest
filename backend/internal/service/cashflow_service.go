package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"portinves/internal/domain"
)

type CashflowService struct {
	repo    domain.CashflowRepository
	notifSrv *NotificationService
}

func NewCashflowService(repo domain.CashflowRepository, notifSrv *NotificationService) *CashflowService {
	return &CashflowService{
		repo:    repo,
		notifSrv: notifSrv,
	}
}

func (s *CashflowService) CreateCashflow(ctx context.Context, userID string, portfolioID *string, cfType domain.CashflowType, amount float64, currency, description string, executedAt time.Time) (*domain.Cashflow, error) {
	cf := &domain.Cashflow{
		ID:          uuid.New().String(),
		UserID:      userID,
		PortfolioID: portfolioID,
		Type:        cfType,
		Amount:      amount,
		Currency:    currency,
		Description: description,
		ExecutedAt:  executedAt,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, cf); err != nil {
		return nil, err
	}

	// Trigger a notification if it's a deposit or withdrawal, or just generally notify them
	if cfType == domain.CashflowTypeDeposit || cfType == domain.CashflowTypeWithdrawal {
		title := "Portfolio Update"
		msg := fmt.Sprintf("A new %s of %s %.2f was recorded.", cfType, currency, amount)
		_, _ = s.notifSrv.CreateNotification(ctx, userID, title, msg)
	} else {
		title := "Cashflow Update"
		msg := fmt.Sprintf("A new %s of %s %.2f was recorded.", cfType, currency, amount)
		_, _ = s.notifSrv.CreateNotification(ctx, userID, title, msg)
	}

	return cf, nil
}

func (s *CashflowService) ListUserCashflows(ctx context.Context, userID string) ([]domain.Cashflow, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *CashflowService) DeleteCashflow(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
