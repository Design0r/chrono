-- name: CreateSession :one
INSERT INTO sessions (id, valid_until, user_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetSessionById :one
SELECT * FROM sessions
WHERE id = ?;

-- name: GetUserFromSession :one
SELECT u.* FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.id = ?;

-- name: DeleteSession :exec
DELETE from sessions 
WHERE id = ?;

-- name: DeleteAllSessions :exec
DELETE from sessions;

-- name: DeleteSessionsForUser :exec
DELETE FROM sessions
WHERE user_id = ?;

-- name: DeleteOtherSessionsForUser :exec
DELETE FROM sessions
WHERE user_id = @user_id
AND id != @keep_id;

-- datetime() normalizes the stored timezone offsets to UTC before comparing.
-- name: DeleteExpiredSessions :execrows
DELETE FROM sessions
WHERE datetime(valid_until) < datetime('now');
