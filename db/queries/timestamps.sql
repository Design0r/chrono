-- name: StartTimestamp :one
INSERT INTO timestamps (user_id)
VALUES (?)
RETURNING *;

-- name: StopTimestamp :one
UPDATE timestamps
SET end_time = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: GetTimestampsById :many
SELECT * from timestamps
WHERE id = ?;

-- name: GetTimestampsInRange :many
SELECT * from timestamps
WHERE user_id = ?
AND start_time <= ?
AND end_time IS NOT NULL
AND end_time >= ?;

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

