-- name: CreateRefreshToken :one
INSERT INTO token_refresh (user_id, year)
VALUES (?,?)
RETURNING *;

-- name: GetRefreshToken :one
SELECT Count(*) FROM token_refresh
WHERE user_id = ?
AND year = ?;

-- name: DeleteAllRefreshTokens :exec
DELETE FROM token_refresh;
