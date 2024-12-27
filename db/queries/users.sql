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
    users.*,
    COUNT(events.id) AS vacation_count
FROM 
    users
LEFT JOIN 
    events
ON 
    users.id = events.user_id
    AND events.name = "urlaub"
    AND events.scheduled_at >= ?
    AND events.scheduled_at < ?
GROUP BY 
    users.id, users.username, users.email, users.vacation_days;

-- name: GetAdmins :many
SELECT * FROM users
WHERE is_superuser = true;
