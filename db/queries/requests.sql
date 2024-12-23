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
AND state = "pending";

-- name: UpdateRequestState :one
UPDATE requests
SET state = ?
WHERE user_id = ?
AND id = ?
RETURNING *;

