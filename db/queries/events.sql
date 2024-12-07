-- name: CreateEvent :one
INSERT INTO events (scheduled_at)
VALUES (?)
RETURNING *;

