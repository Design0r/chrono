-- +goose Up
-- Unfinished timestamps that were superseded by a newer one can never be
-- stopped and are ignored everywhere, remove them so every user has at most
-- one running timestamp.
DELETE FROM timestamps
WHERE end_time IS NULL
AND EXISTS (
    SELECT 1 FROM timestamps newer
    WHERE newer.user_id = timestamps.user_id
    AND newer.id > timestamps.id
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_timestamps_running_user
ON timestamps(user_id) WHERE end_time IS NULL;


-- +goose Down
DROP INDEX IF EXISTS idx_timestamps_running_user;
