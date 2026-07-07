package postgres

import (
	"context"

	"portinves/internal/domain"
)

type NotificationRepo struct {
	db DBTX
}

func NewNotificationRepo(db DBTX) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func (r *NotificationRepo) Create(ctx context.Context, n *domain.Notification) error {
	const q = `INSERT INTO notifications (id, user_id, title, message, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	
	_, err := r.db.Exec(ctx, q, n.ID, n.UserID, n.Title, n.Message, n.IsRead, n.CreatedAt)
	return wrapError(err, "notification")
}

func (r *NotificationRepo) ListByUserID(ctx context.Context, userID string, unreadOnly bool) ([]domain.Notification, error) {
	const q = `SELECT id, user_id, title, message, is_read, created_at
		FROM notifications
		WHERE user_id = $1 AND ($2::boolean = false OR is_read = false)
		ORDER BY created_at DESC`
	
	rows, err := r.db.Query(ctx, q, userID, unreadOnly)
	if err != nil {
		return nil, wrapError(err, "notification")
	}
	defer rows.Close()

	var notifs []domain.Notification
	for rows.Next() {
		var n domain.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, wrapError(err, "notification")
		}
		notifs = append(notifs, n)
	}
	return notifs, rows.Err()
}

func (r *NotificationRepo) MarkAsRead(ctx context.Context, id, userID string) error {
	const q = `UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`
	cmd, err := r.db.Exec(ctx, q, id, userID)
	if err != nil {
		return wrapError(err, "notification")
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrNotFound("notification")
	}
	return nil
}
