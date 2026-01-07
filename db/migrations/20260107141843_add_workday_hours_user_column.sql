-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN workday_hours FLOAT NOT NULL DEFAULT 8.0;
ALTER TABLE users ADD COLUMN workdays_week FLOAT NOT NULL DEFAULT 5.0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN workday_hours;
ALTER TABLE users DROP COLUMN workdays_week;
-- +goose StatementEnd
