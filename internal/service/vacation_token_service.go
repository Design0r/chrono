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
	r   domain.VacationTokenRepository
	log *slog.Logger
}

func NewVacationTokenService(
	r domain.VacationTokenRepository,
	log *slog.Logger,
) vacationTokenService {
	return vacationTokenService{r: r, log: log}
}

func (s *vacationTokenService) Create(
	ctx context.Context,
	t domain.CreateVacationToken,
) (*domain.VacationToken, error) {
	return s.r.Create(ctx, t)
}

func (s *vacationTokenService) Delete(ctx context.Context, id int64) error {
	return s.r.Delete(ctx, id)
}

func (s *vacationTokenService) DeleteAll(ctx context.Context) error {
	return s.r.DeleteAll(ctx)
}

func (s *vacationTokenService) GetRemainingVacationForUser(
	ctx context.Context,
	userId int64,
	start time.Time,
	end time.Time,
) (float64, error) {
	return s.r.GetRemainingVacationForUser(ctx, userId, start, end)
}
