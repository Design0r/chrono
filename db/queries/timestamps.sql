-- name: CreateTimestamp :one
INSERT INTO timestamps (user_id, start_time, end_time)
VALUES (?,?,?)
RETURNING *;

-- Only inserts if the user has no running timestamp, so a device with
-- outdated state cannot start a second timer next to it.
-- Returns no row when a running timestamp already exists.
-- name: StartTimestamp :one
INSERT INTO timestamps (user_id)
SELECT u.id FROM users u
WHERE u.id = @user_id
AND NOT EXISTS (
  SELECT 1 FROM timestamps t
  WHERE t.user_id = u.id
  AND t.end_time IS NULL
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

-- The running timestamp wins over newer finished ones, e.g. entries an admin
-- created for the user while the timer was running.
-- name: GetLatestTimestamp :one
SELECT * FROM timestamps
WHERE user_id = ?
ORDER BY (end_time IS NULL) DESC, id DESC
LIMIT 1;

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

