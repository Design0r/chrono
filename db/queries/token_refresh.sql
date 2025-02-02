-- name: CreateTokenRefresh :one
INSERT INTO token_refresh (user_id, year)
VALUES (?,?)
RETURNING *;

-- name: GetTokenRefresh :one
SELECT Count(*) FROM token_refresh
WHERE user_id = ?
AND year = ?;

-- name: ResetTokenRefresh :exec
DELETE FROM token_refresh;
