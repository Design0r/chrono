package service

import (
	"context"
	"log/slog"
	"time"

	"chrono/internal/domain"
)

type VacationTokenService interface {
	Create(ctx context.Context, t domain.CreateVacationToken) (*domain.VacationToken, error)
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
	ctx context.Context,
	t domain.CreateVacationToken,
) (*domain.VacationToken, error) {
	return svc.vacation.Create(ctx, t)
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
