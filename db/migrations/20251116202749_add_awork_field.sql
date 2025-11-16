-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN awork_id TEXT;
-- +goose StatementEnd

-- +goose StatementBegin
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
