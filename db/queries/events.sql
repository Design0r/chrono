-- name: CreateEvent :one
INSERT INTO events (name, user_id, scheduled_at, state)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetEventsForDay :many
SELECT * FROM events 
WHERE Date(scheduled_at) = ?;

-- name: DeleteEvent :one
DELETE from events
WHERE id = ?
RETURNING *;

-- name: GetUserPendingEvents :many
SELECT * FROM events
WHERE user_id = ?
AND state = "pending";

-- name: GetEventsForMonth :many
SELECT *
FROM events e
JOIN users u ON e.user_id = u.id
WHERE scheduled_at >= ? AND scheduled_at < ?;

-- name: GetAcceptedEventsForMonth :many
SELECT *
FROM events e
JOIN users u ON e.user_id = u.id
WHERE scheduled_at >= ? AND scheduled_at < ?
AND state = "accepted";

-- name: GetPendingEventsForYear :one
SELECT Count(id) from events
WHERE state = "pending"
AND scheduled_at >= ?
AND scheduled_at < ?
AND user_id = ?;

-- name: GetVacationCountForUser :one 
SELECT Count(*) from events
WHERE user_id = ?
AND scheduled_at >= ?
AND scheduled_at < ?
AND name = "urlaub"
AND state = "accepted";

-- name: UpdateEventState :one
UPDATE events
SET state = ?,
edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;
