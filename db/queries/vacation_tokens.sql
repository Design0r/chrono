-- name: CreateToken :one 
INSERT INTO vacation_tokens (user_id, start_date, end_date, value)
VALUES (?,?,?,?)
RETURNING *;

-- name: DeleteToken :exec
DELETE FROM vacation_tokens
WHERE id = ?;

-- name: GetValidUserTokens :many
SELECT * from vacation_tokens
WHERE user_id = ? 
AND ? >= start_date
AND ? <= end_date;

-- name: GetValidUserTokenSum :one
SELECT SUM(value) FROM vacation_tokens
WHERE user_id = ? 
AND start_date <= ?
AND end_date >= ?;

-- name: UpdateYearlyTokens :exec
UPDATE vacation_tokens
SET value = ?
WHERE user_id = ? 
AND start_date <= ?
AND end_date >= ?;

