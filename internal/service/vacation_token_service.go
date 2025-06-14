package service

import (
	"context"
	"log/slog"
	"time"

	"chrono/internal/domain"
)

type VacationTokenService interface {
	Create(
		ctx context.Context,
		value float64,
		year int,
		userId int64,
	) (*domain.VacationToken, error)
	Delete(ctx context.Context, id int64) error
	DeleteAll(ctx context.Context) error
	GetRemainingVacationForUser(
		ctx context.Context,
		userId int64,
		start time.Time,
		end time.Time,
	) (float64, error)
}

type vacationTokenService struct {
	vacation domain.VacationTokenRepository
	log      *slog.Logger
}

func NewVacationTokenService(
	r domain.VacationTokenRepository,
	log *slog.Logger,
) vacationTokenService {
	return vacationTokenService{vacation: r, log: log}
}

func (svc *vacationTokenService) Create(
	ctx context.Context, value float64, year int, userId int64,
) (*domain.VacationToken, error) {
	start := time.Date(
		year,
		time.January,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	end := time.Now().AddDate(1, 3, 0)
	return svc.vacation.Create(
		ctx,
		domain.CreateVacationToken{StartDate: start, EndDate: end, UserID: userId, Value: value},
	)
}

func (svc *vacationTokenService) Delete(ctx context.Context, id int64) error {
	return svc.vacation.Delete(ctx, id)
}

func (svc *vacationTokenService) DeleteAll(ctx context.Context) error {
	return svc.vacation.DeleteAll(ctx)
}

func (svc *vacationTokenService) GetRemainingVacationForUser(
	ctx context.Context,
	userId int64,
	start time.Time,
	end time.Time,
) (float64, error) {
	return svc.vacation.GetRemainingVacationForUser(ctx, userId, start, end)
}
