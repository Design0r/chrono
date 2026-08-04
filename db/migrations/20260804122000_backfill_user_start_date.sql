-- +goose Up
-- +goose StatementBegin
UPDATE users SET start_date = created_at WHERE start_date = '2026-01-01 00:00:00';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE users SET start_date = '2026-01-01 00:00:00';
-- +goose StatementEnd
