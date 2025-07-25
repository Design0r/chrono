// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: settings.sql

package repo

import (
	"context"
)

const CreateSettings = `-- name: CreateSettings :one
INSERT INTO settings (signup_enabled)
VALUES (?)
RETURNING id, signup_enabled
`

func (q *Queries) CreateSettings(ctx context.Context, signupEnabled bool) (Setting, error) {
	row := q.db.QueryRowContext(ctx, CreateSettings, signupEnabled)
	var i Setting
	err := row.Scan(&i.ID, &i.SignupEnabled)
	return i, err
}

const DeleteSettings = `-- name: DeleteSettings :exec
DELETE FROM settings 
WHERE id = ?
`

func (q *Queries) DeleteSettings(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, DeleteSettings, id)
	return err
}

const GetSettingsById = `-- name: GetSettingsById :one
SELECT id, signup_enabled FROM settings 
WHERE id = ?
`

func (q *Queries) GetSettingsById(ctx context.Context, id int64) (Setting, error) {
	row := q.db.QueryRowContext(ctx, GetSettingsById, id)
	var i Setting
	err := row.Scan(&i.ID, &i.SignupEnabled)
	return i, err
}

const UpdateSettings = `-- name: UpdateSettings :one
UPDATE settings
SET signup_enabled = ?
WHERE id = ?
RETURNING id, signup_enabled
`

type UpdateSettingsParams struct {
	SignupEnabled bool  `json:"signup_enabled"`
	ID            int64 `json:"id"`
}

func (q *Queries) UpdateSettings(ctx context.Context, arg UpdateSettingsParams) (Setting, error) {
	row := q.db.QueryRowContext(ctx, UpdateSettings, arg.SignupEnabled, arg.ID)
	var i Setting
	err := row.Scan(&i.ID, &i.SignupEnabled)
	return i, err
}
