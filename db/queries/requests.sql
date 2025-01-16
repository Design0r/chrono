-- name: CreateRequest :one 
INSERT INTO requests (message, state, user_id, event_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetUserRequests :many
SELECT * FROM requests
WHERE user_id = ?;

-- name: GetPendingRequests :many
SELECT * FROM requests r
JOIN users u ON r.user_id = u.id
JOIN events e ON r.event_id = e.id
WHERE r.state = "pending"
ORDER BY e.scheduled_at;

-- name: UpdateRequestState :one
UPDATE requests
SET state = ?,
edited_by = ?,
edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateRequestStateRange :exec
UPDATE requests
SET state = ?,
    edited_by = ?,
    edited_at = CURRENT_TIMESTAMP
WHERE requests.user_id = ?
  AND event_id IN (
    SELECT e.id
    FROM events e
    WHERE e.scheduled_at >= ?
      AND e.scheduled_at <= ?
  );

-- name: GetRequestRange :many
SELECT * FROM requests r
JOIN users u ON r.user_id = u.id
JOIN events e ON r.event_id = e.id
WHERE r.user_id = ?
AND e.scheduled_at >= ?
AND e.scheduled_at <= ?
ORDER BY e.scheduled_at;
