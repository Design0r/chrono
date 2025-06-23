-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user';
ALTER TABLE users ADD COLUMN enabled BOOLEAN NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose StatementBegin
UPDATE users
SET role = "admin"
WHERE is_superuser = 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
