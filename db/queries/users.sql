-- name: CreateUser :one
INSERT INTO users (username, vacation_days, email, password, is_superuser)
VALUES (?, ?, ?, ?, ?)
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

-- name: UpdateVacationDays :one
UPDATE users
SET vacation_days = ?
WHERE id = ?
RETURNING *;

