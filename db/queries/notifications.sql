-- name: CreateNotification :one
INSERT INTO notifications (message, user_id)
VALUES (?, ?)
RETURNING *;

-- name: ClearNotification :exec
UPDATE notifications
SET viewed_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: ClearAllNotification :exec
UPDATE notifications
SET viewed_at = CURRENT_TIMESTAMP
WHERE user_id = ?;

-- name: GetUserNotifications :many
SELECT * FROM notifications
WHERE user_id = ?
AND viewed_at IS NULL;
