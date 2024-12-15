-- name: CreateEvent :one
INSERT INTO events (name, user_id, scheduled_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetEventsForDay :many
SELECT * FROM events 
WHERE Date(scheduled_at) = ?;

-- name: DeleteEvent :one
DELETE from events
WHERE id = ?
RETURNING *;

-- name: GetEventsForMonth :many
SELECT *
FROM events e
JOIN users u ON e.user_id = u.id
WHERE scheduled_at >= ? AND scheduled_at < ?;

-- name: GetVacationCountForUser :one 
SELECT Count(*) from events
WHERE user_id = ?
AND scheduled_at >= ?
AND scheduled_at < ?
AND name = "urlaub";
