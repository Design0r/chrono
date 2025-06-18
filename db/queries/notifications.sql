-- name: CreateNotification :one
INSERT INTO notifications (message)
VALUES (?)
RETURNING *;

-- name: ClearNotification :one
UPDATE notifications
SET viewed_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateNotification :one
update notifications
SET viewed_at = CURRENT_TIMESTAMP,
message = ?
RETURNING *;



