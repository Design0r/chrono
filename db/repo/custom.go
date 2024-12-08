package repo

import (
	"context"
	"fmt"
	"time"
)

const getEventsForMonth = `-- name: GetEventsForMonth :many
SELECT id, scheduled_at, created_at, edited_at
FROM events
WHERE strftime('%Y-%m', scheduled_at) = ?
`

func (q *Queries) CustomGetEventsForMonth(
	ctx context.Context,
	scheduledAt time.Time,
) ([]Event, error) {
	strtime := fmt.Sprintf("%v-%v", scheduledAt.Year(), int(scheduledAt.Month()))
	rows, err := q.db.QueryContext(ctx, getEventsForMonth, strtime)
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
			&i.CreatedAt,
			&i.EditedAt,
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
