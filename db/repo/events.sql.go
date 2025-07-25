// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: events.sql

package repo

import (
	"context"
	"time"
)

const CreateEvent = `-- name: CreateEvent :one
INSERT INTO events (name, user_id, scheduled_at, state)
VALUES (?, ?, ?, ?)
RETURNING id, scheduled_at, name, state, created_at, edited_at, user_id
`

type CreateEventParams struct {
	Name        string    `json:"name"`
	UserID      int64     `json:"user_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	State       string    `json:"state"`
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, CreateEvent,
		arg.Name,
		arg.UserID,
		arg.ScheduledAt,
		arg.State,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.ScheduledAt,
		&i.Name,
		&i.State,
		&i.CreatedAt,
		&i.EditedAt,
		&i.UserID,
	)
	return i, err
}

const DeleteEvent = `-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = ?
`

func (q *Queries) DeleteEvent(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, DeleteEvent, id)
	return err
}

const GetConflictingEventUsers = `-- name: GetConflictingEventUsers :many
SELECT DISTINCT u.id, u.username, u.email, u.password, u.vacation_days, u.is_superuser, u.created_at, u.edited_at, u.color, u.role, u.enabled FROM events e
JOIN users u on e.user_id = u.id
WHERE u.id != ? 
AND e.scheduled_at >= ?
AND e.scheduled_at <= ?
`

type GetConflictingEventUsersParams struct {
	ID            int64     `json:"id"`
	ScheduledAt   time.Time `json:"scheduled_at"`
	ScheduledAt_2 time.Time `json:"scheduled_at_2"`
}

func (q *Queries) GetConflictingEventUsers(ctx context.Context, arg GetConflictingEventUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, GetConflictingEventUsers, arg.ID, arg.ScheduledAt, arg.ScheduledAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.VacationDays,
			&i.IsSuperuser,
			&i.CreatedAt,
			&i.EditedAt,
			&i.Color,
			&i.Role,
			&i.Enabled,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetEventById = `-- name: GetEventById :one
SELECT id, scheduled_at, name, state, created_at, edited_at, user_id FROM events
WHERE id = ?
`

func (q *Queries) GetEventById(ctx context.Context, id int64) (Event, error) {
	row := q.db.QueryRowContext(ctx, GetEventById, id)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.ScheduledAt,
		&i.Name,
		&i.State,
		&i.CreatedAt,
		&i.EditedAt,
		&i.UserID,
	)
	return i, err
}

const GetEventsForDay = `-- name: GetEventsForDay :many
SELECT id, scheduled_at, name, state, created_at, edited_at, user_id FROM events 
WHERE Date(scheduled_at) = ?
`

func (q *Queries) GetEventsForDay(ctx context.Context, scheduledAt time.Time) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, GetEventsForDay, scheduledAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.ScheduledAt,
			&i.Name,
			&i.State,
			&i.CreatedAt,
			&i.EditedAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetEventsForMonth = `-- name: GetEventsForMonth :many
SELECT e.id, scheduled_at, name, state, e.created_at, e.edited_at, user_id, u.id, username, email, password, vacation_days, is_superuser, u.created_at, u.edited_at, color, role, enabled
FROM events e
JOIN users u ON e.user_id = u.id
WHERE scheduled_at >= ? AND scheduled_at < ?
`

type GetEventsForMonthParams struct {
	ScheduledAt   time.Time `json:"scheduled_at"`
	ScheduledAt_2 time.Time `json:"scheduled_at_2"`
}

type GetEventsForMonthRow struct {
	ID           int64     `json:"id"`
	ScheduledAt  time.Time `json:"scheduled_at"`
	Name         string    `json:"name"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	EditedAt     time.Time `json:"edited_at"`
	UserID       int64     `json:"user_id"`
	ID_2         int64     `json:"id_2"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	VacationDays int64     `json:"vacation_days"`
	IsSuperuser  bool      `json:"is_superuser"`
	CreatedAt_2  time.Time `json:"created_at_2"`
	EditedAt_2   time.Time `json:"edited_at_2"`
	Color        string    `json:"color"`
	Role         string    `json:"role"`
	Enabled      bool      `json:"enabled"`
}

func (q *Queries) GetEventsForMonth(ctx context.Context, arg GetEventsForMonthParams) ([]GetEventsForMonthRow, error) {
	rows, err := q.db.QueryContext(ctx, GetEventsForMonth, arg.ScheduledAt, arg.ScheduledAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEventsForMonthRow
	for rows.Next() {
		var i GetEventsForMonthRow
		if err := rows.Scan(
			&i.ID,
			&i.ScheduledAt,
			&i.Name,
			&i.State,
			&i.CreatedAt,
			&i.EditedAt,
			&i.UserID,
			&i.ID_2,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.VacationDays,
			&i.IsSuperuser,
			&i.CreatedAt_2,
			&i.EditedAt_2,
			&i.Color,
			&i.Role,
			&i.Enabled,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetEventsForYear = `-- name: GetEventsForYear :many
SELECT e.id, scheduled_at, name, state, e.created_at, e.edited_at, user_id, u.id, username, email, password, vacation_days, is_superuser, u.created_at, u.edited_at, color, role, enabled FROM events e
JOIN users u ON e.user_id = u.id
WHERE e.scheduled_at >= ? 
  AND e.scheduled_at < ?
  AND e.state = "accepted"
  AND (e.name IN ('urlaub', 'urlaub halbtags') OR e.user_id = 1)

ORDER BY scheduled_at
`

type GetEventsForYearParams struct {
	ScheduledAt   time.Time `json:"scheduled_at"`
	ScheduledAt_2 time.Time `json:"scheduled_at_2"`
}

type GetEventsForYearRow struct {
	ID           int64     `json:"id"`
	ScheduledAt  time.Time `json:"scheduled_at"`
	Name         string    `json:"name"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	EditedAt     time.Time `json:"edited_at"`
	UserID       int64     `json:"user_id"`
	ID_2         int64     `json:"id_2"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	VacationDays int64     `json:"vacation_days"`
	IsSuperuser  bool      `json:"is_superuser"`
	CreatedAt_2  time.Time `json:"created_at_2"`
	EditedAt_2   time.Time `json:"edited_at_2"`
	Color        string    `json:"color"`
	Role         string    `json:"role"`
	Enabled      bool      `json:"enabled"`
}

func (q *Queries) GetEventsForYear(ctx context.Context, arg GetEventsForYearParams) ([]GetEventsForYearRow, error) {
	rows, err := q.db.QueryContext(ctx, GetEventsForYear, arg.ScheduledAt, arg.ScheduledAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetEventsForYearRow
	for rows.Next() {
		var i GetEventsForYearRow
		if err := rows.Scan(
			&i.ID,
			&i.ScheduledAt,
			&i.Name,
			&i.State,
			&i.CreatedAt,
			&i.EditedAt,
			&i.UserID,
			&i.ID_2,
			&i.Username,
			&i.Email,
			&i.Password,
			&i.VacationDays,
			&i.IsSuperuser,
			&i.CreatedAt_2,
			&i.EditedAt_2,
			&i.Color,
			&i.Role,
			&i.Enabled,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const GetPendingEventsForYear = `-- name: GetPendingEventsForYear :one
SELECT Count(id) from events
WHERE state = "pending"
AND scheduled_at >= ?
AND scheduled_at < ?
AND user_id = ?
`

type GetPendingEventsForYearParams struct {
	ScheduledAt   time.Time `json:"scheduled_at"`
	ScheduledAt_2 time.Time `json:"scheduled_at_2"`
	UserID        int64     `json:"user_id"`
}

func (q *Queries) GetPendingEventsForYear(ctx context.Context, arg GetPendingEventsForYearParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, GetPendingEventsForYear, arg.ScheduledAt, arg.ScheduledAt_2, arg.UserID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const GetVacationCountForUser = `-- name: GetVacationCountForUser :one
SELECT 
  SUM(
    CASE
      WHEN name = 'urlaub'          THEN 1
      WHEN name = 'urlaub halbtags' THEN 0.5
      ELSE 0
    END
  ) 
FROM events
WHERE user_id = ?
  AND scheduled_at >= ?
  AND scheduled_at < ?
  AND name IN ('urlaub', 'urlaub halbtags')
  AND state = 'accepted'
`

type GetVacationCountForUserParams struct {
	UserID        int64     `json:"user_id"`
	ScheduledAt   time.Time `json:"scheduled_at"`
	ScheduledAt_2 time.Time `json:"scheduled_at_2"`
}

func (q *Queries) GetVacationCountForUser(ctx context.Context, arg GetVacationCountForUserParams) (*float64, error) {
	row := q.db.QueryRowContext(ctx, GetVacationCountForUser, arg.UserID, arg.ScheduledAt, arg.ScheduledAt_2)
	var sum *float64
	err := row.Scan(&sum)
	return sum, err
}

const UpdateEventState = `-- name: UpdateEventState :one
UPDATE events
SET state = ?,
edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, scheduled_at, name, state, created_at, edited_at, user_id
`

type UpdateEventStateParams struct {
	State string `json:"state"`
	ID    int64  `json:"id"`
}

func (q *Queries) UpdateEventState(ctx context.Context, arg UpdateEventStateParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, UpdateEventState, arg.State, arg.ID)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.ScheduledAt,
		&i.Name,
		&i.State,
		&i.CreatedAt,
		&i.EditedAt,
		&i.UserID,
	)
	return i, err
}

const UpdateEventsRange = `-- name: UpdateEventsRange :exec
UPDATE events
SET state = ?, 
edited_at = CURRENT_TIMESTAMP
WHERE user_id = ? 
AND scheduled_at >= ?
AND scheduled_at <= ?
`

type UpdateEventsRangeParams struct {
	State         string    `json:"state"`
	UserID        int64     `json:"user_id"`
	ScheduledAt   time.Time `json:"scheduled_at"`
	ScheduledAt_2 time.Time `json:"scheduled_at_2"`
}

func (q *Queries) UpdateEventsRange(ctx context.Context, arg UpdateEventsRangeParams) error {
	_, err := q.db.ExecContext(ctx, UpdateEventsRange,
		arg.State,
		arg.UserID,
		arg.ScheduledAt,
		arg.ScheduledAt_2,
	)
	return err
}
