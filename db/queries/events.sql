-- name: CreateEvent :one
INSERT INTO events (user_id, scheduled_at)
VALUES (?, ?)
RETURNING *;

-- name: GetEventsForDay :many
SELECT * FROM events 
WHERE Date(scheduled_at) = ?;

-- name: DeleteEvent :exec 
DELETE from events
WHERE id = ?;

-- name: GetEventsForMonth :many
SELECT *
FROM events e
JOIN users u ON e.user_id = u.id
WHERE scheduled_at >= ? AND scheduled_at < ?
