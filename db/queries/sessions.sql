-- name: CreateSession :one
INSERT INTO sessions (id, valid_until, user_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions
WHERE id = ?;

-- name: GetValidSession :one
SELECT * FROM sessions
WHERE id = ?
AND valid_until >= ?;

-- name: GetUserFromSession :one
SELECT u.* FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.id = ?;

-- name: DeleteSession :exec
DELETE from sessions 
WHERE id = ?;

-- name: DeleteAllSessions :exec
DELETE from sessions;
