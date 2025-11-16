-- name: CreateUser :one
INSERT INTO users (username, color, vacation_days, email, password, is_superuser)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByName :one
SELECT * FROM users
WHERE username = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: UpdateUser :one
UPDATE users
SET color = ?,
username = ?,
email = ?,
password = ?,
role = ?,
vacation_days = ?,
is_superuser = ?,
edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: GetAdmins :many
SELECT * FROM users
WHERE is_superuser = true;

-- name: GetAllUsers :many
SELECT * FROM users
WHERE id != 1;
