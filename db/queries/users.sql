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
SET vacation_days = ?,
edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET vacation_days = ?,
username = ?,
email = ?,
edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: GetUsersWithVacationCount :many
SELECT
    u.*,
    SUM(vt.value) AS vac_remaining,
    SUM(0.5) AS vac_used
FROM users AS u
JOIN vacation_tokens vt 
    ON u.id = vt.user_id
    AND vt.start_date <= ?
    AND vt.end_date   >= ?
GROUP BY u.id
ORDER BY u.id;

-- name: GetAdmins :many
SELECT * FROM users
WHERE is_superuser = true;

-- name: ToggleAdmin :one
UPDATE users
SET is_superuser = NOT is_superuser
WHERE id = ?
RETURNING *;

-- name: GetAllUsers :many
SELECT * FROM users
WHERE id != 1;

-- name: SetUserVacation :exec
UPDATE users
SET vacation_days = ?
WHERE id = ?;
