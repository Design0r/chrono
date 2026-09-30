-- name: CreateTimestamp :one
INSERT INTO timestamps (user_id, start_time, end_time)
VALUES (?,?,?)
RETURNING *;

-- Only inserts if the user's latest timestamp is not still running, so a
-- device with outdated state cannot start a second timer next to it.
-- Returns no row when a running timestamp already exists.
-- name: StartTimestamp :one
INSERT INTO timestamps (user_id)
SELECT u.id FROM users u
WHERE u.id = @user_id
AND NOT EXISTS (
  SELECT 1 FROM timestamps t
  WHERE t.end_time IS NULL
  AND t.id = (SELECT MAX(id) FROM timestamps WHERE user_id = @user_id)
)
RETURNING *;

-- name: UpdateTimestamp :one
UPDATE timestamps
SET start_time = ?,
end_time = ?
WHERE id = ?
RETURNING *;

-- name: StopTimestamp :one
UPDATE timestamps
SET end_time = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteTimestamp :exec
DELETE FROM timestamps
WHERE id = ?;

-- name: GetTimestampById :one
SELECT * FROM timestamps
WHERE id = ?;

-- name: GetAllTimestampsForUser :many
SELECT * FROM timestamps
WHERE user_id = ?;

-- name: GetTimestampsInRange :many
SELECT * FROM timestamps
WHERE user_id = ?
AND start_time < @end_time
AND end_time IS NOT NULL
AND end_time > @start_time;

-- name: GetAllTimestampsInRange :many
SELECT * FROM timestamps
WHERE start_time < @end_time
AND end_time IS NOT NULL
AND end_time > @start_time;

-- name: GetLatestTimestamp :one
SELECT * FROM timestamps
WHERE user_id = ?
ORDER BY id DESC;

-- name: GetTotalSecondsInRange :one
SELECT
  SUM(
    MAX(
      0,
      strftime('%s', MIN(end_time, @range_end))
      - strftime('%s', MAX(start_time, @range_start))
    )
  ) AS total_seconds
FROM timestamps
WHERE user_id = @user_id
  AND end_time IS NOT NULL
  AND start_time < @range_end
  AND end_time > @range_start;

