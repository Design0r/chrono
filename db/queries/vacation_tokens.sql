-- name: CreateVacationToken :one 
INSERT INTO vacation_tokens (user_id, start_date, end_date, value)
VALUES (?,?,?,?)
RETURNING *;

-- name: DeleteVacationToken :exec
DELETE FROM vacation_tokens
WHERE id = ?;

-- name: GetRemainingVacationForUser :one
SELECT SUM(value) FROM vacation_tokens
WHERE user_id = ? 
AND start_date <= ?
AND end_date >= ?;

-- name: DeleteAllVacationTokens :exec
DELETE FROM vacation_tokens;
