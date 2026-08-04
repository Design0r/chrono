-- +goose Up
ALTER TABLE users ADD COLUMN start_date DATETIME NOT NULL DEFAULT '2026-01-01 00:00:00';

-- +goose Down
ALTER TABLE users DROP COLUMN start_date;
