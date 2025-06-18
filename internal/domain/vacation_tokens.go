package domain

import (
	"context"
	"time"
)

type VacationToken struct {
	ID        int64     `json:"id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Value     float64   `json:"value"`
	UserID    int64     `json:"user_id"`
}

type CreateVacationToken struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Value     float64   `json:"value"`
	UserID    int64     `json:"user_id"`
}

type VacationTokenRepository interface {
	Create(ctx context.Context, t CreateVacationToken) (*VacationToken, error)
	Delete(ctx context.Context, id int64) error
	DeleteAll(ctx context.Context) error
	GetRemainingVacationForUser(
		ctx context.Context,
		userId int64,
		start time.Time,
		end time.Time,
	) (float64, error)
}
