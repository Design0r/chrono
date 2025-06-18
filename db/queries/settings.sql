-- name: CreateSettings :one
INSERT INTO settings (signup_enabled)
VALUES (?)
RETURNING *;

-- name: GetSettingsById :one
SELECT * FROM settings 
WHERE id = ?; 

-- name: UpdateSettings :one
UPDATE settings
SET signup_enabled = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSettings :exec
DELETE FROM settings 
WHERE id = ?;
