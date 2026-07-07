package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"portinves/internal/domain"
	"portinves/internal/websocket"
)

type NotificationService struct {
	repo domain.NotificationRepository
	hub  *websocket.Hub
}

func NewNotificationService(repo domain.NotificationRepository, hub *websocket.Hub) *NotificationService {
	return &NotificationService{
		repo: repo,
		hub:  hub,
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, userID, title, message string) (*domain.Notification, error) {
	n := &domain.Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Title:     title,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}

	// Push real-time notification
	s.hub.BroadcastToUser(userID, map[string]interface{}{
		"type":         "notification",
		"notification": n,
	})

	return n, nil
}

func (s *NotificationService) ListUserNotifications(ctx context.Context, userID string, unreadOnly bool) ([]domain.Notification, error) {
	return s.repo.ListByUserID(ctx, userID, unreadOnly)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id, userID string) error {
	return s.repo.MarkAsRead(ctx, id, userID)
}
