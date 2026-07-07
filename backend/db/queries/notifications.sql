-- name: CreateNotification :exec
INSERT INTO notifications (
    id, user_id, title, message, is_read, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: ListNotificationsByUserID :many
SELECT id, user_id, title, message, is_read, created_at
FROM notifications
WHERE user_id = $1 AND ($2::boolean = false OR is_read = false)
ORDER BY created_at DESC;

-- name: MarkNotificationAsRead :exec
UPDATE notifications 
SET is_read = true 
WHERE id = $1 AND user_id = $2;
