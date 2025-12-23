package domain

import (
	"context"
	"time"
)

type Timestamp struct {
	ID        int64      `json:"id"`
	StartTime time.Time  `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	UserID    int64      `json:"user_id"`
}

type TimestampsRepository interface {
	GetById(ctx context.Context, id int64) (Timestamp, error)
	Start(ctx context.Context, userId int64) (Timestamp, error)
	Stop(ctx context.Context, id int64) (Timestamp, error)
	Delete(ctx context.Context, id int64) error
	GetInRange(
		ctx context.Context,
		userId int64,
		start time.Time,
		stop time.Time,
	) ([]Timestamp, error)
	GetTotalSecondsInRange(
		ctx context.Context,
		userId int64,
		start time.Time,
		stop time.Time,
	) (float64, error)
}
