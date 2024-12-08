-- name: CreateEvent :one
INSERT INTO events (scheduled_at)
VALUES (?)
RETURNING *;

-- name: GetEventsForDay :many
SELECT * FROM events 
WHERE Date(scheduled_at) = ?;

-- name: GetEventsForMonth :many
-- SELECT * FROM events
-- WHERE strftime('%Y-%m', scheduled_at) = ?;
