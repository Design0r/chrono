-- name: CreateEvent :one
INSERT INTO events (name, user_id, scheduled_at, state)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetEventById :one
SELECT * FROM events
WHERE id = ?;

-- name: GetEventsForDay :many
SELECT * FROM events 
WHERE Date(scheduled_at) = ?;

-- name: GetEventsForYear :many
SELECT * FROM events e
JOIN users u ON e.user_id = u.id
WHERE e.scheduled_at >= ? 
  AND e.scheduled_at < ?
  AND e.state = "accepted"
  AND (e.name IN ('urlaub', 'urlaub halbtags') OR e.user_id = 1)

ORDER BY scheduled_at;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = ?;


-- name: GetEventsForMonth :many
SELECT *
FROM events e
JOIN users u ON e.user_id = u.id
WHERE scheduled_at >= ? AND scheduled_at < ?;

-- name: GetPendingEventsForYear :one
SELECT Count(id) from events
WHERE state = "pending"
AND scheduled_at >= ?
AND scheduled_at < ?
AND user_id = ?;

-- name: GetVacationCountForUser :one 
SELECT 
  SUM(
    CASE
      WHEN name = 'urlaub'          THEN 1
      WHEN name = 'urlaub halbtags' THEN 0.5
      ELSE 0
    END
  ) 
FROM events
WHERE user_id = ?
  AND scheduled_at >= ?
  AND scheduled_at < ?
  AND name IN ('urlaub', 'urlaub halbtags')
  AND state = 'accepted';

-- name: UpdateEventState :one
UPDATE events
SET state = ?,
edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateEventsRange :exec
UPDATE events
SET state = ?, 
edited_at = CURRENT_TIMESTAMP
WHERE user_id = ? 
AND scheduled_at >= ?
AND scheduled_at <= ?;

-- name: GetConflictingEventUsers :many
SELECT DISTINCT u.* FROM events e
JOIN users u on e.user_id = u.id
WHERE u.id != ? 
AND e.scheduled_at >= ?
AND e.scheduled_at <= ?;
