-- name: CreateNotification :one
INSERT INTO notifications (message)
VALUES (?)
RETURNING *;

-- name: ClearNotification :exec
UPDATE notifications
SET viewed_at = CURRENT_TIMESTAMP
WHERE id = ?;
