-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD color TEXT NOT NULL DEFAULT '#000000';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
