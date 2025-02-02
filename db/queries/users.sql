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

-- name: UpdateVacationDays :one
UPDATE users
SET vacation_days = ?,
edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET color = ?,
username = ?,
email = ?,
edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: GetUsersWithVacationCount :many
SELECT
    u.*,
    COALESCE(SUM(vt.value), 0.0) AS vac_remaining,
    COALESCE(SUM(0.5), 0.0) AS vac_used
FROM users AS u
LEFT JOIN vacation_tokens vt 
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

-- name: SetUserVacation :one
UPDATE users
SET vacation_days = ?
WHERE id = ?
RETURNING *;

-- name: SetUserColor :exec
UPDATE users
SET color = ?
WHERE id = ?;
