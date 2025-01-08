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
WHERE r.state = "pending";

-- name: UpdateRequestState :one
UPDATE requests
SET state = ?,
edited_by = ?,
edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateRequestStateRange :exec
UPDATE events
SET state = ?,
edited_at = CURRENT_TIMESTAMP
WHERE user_id = ?
AND scheduled_at >= ?
AND scheduled_at <= ?
RETURNING *;

