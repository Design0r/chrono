-- name: CreateEvent :one
INSERT INTO events (scheduled_at)
VALUES (?)
RETURNING *;

-- name: GetEventsForDay :many
SELECT * FROM events 
WHERE Date(scheduled_at) = ?;

-- name: DeleteEvent :exec 
DELETE from events
WHERE id = ?;
