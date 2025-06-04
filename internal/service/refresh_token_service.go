package service

import (
	"context"
	"log/slog"

	"chrono/internal/domain"
)

type RefreshTokenService interface {
	Create(ctx context.Context, t domain.CreateVacationToken) (*domain.VacationToken, error)
	DeleteAll(ctx context.Context) error
	ExistsForUser(ctx context.Context, userId int64, year int) (bool, error)
}

type refreshTokenService struct {
	r   domain.RefreshTokenRepository
	log *slog.Logger
}

func NewRefreshTokenService(
	r domain.RefreshTokenRepository,
	log *slog.Logger,
) refreshTokenService {
	return refreshTokenService{r: r, log: log}
}

func (s *refreshTokenService) Create(
	ctx context.Context,
	t domain.CreateVacationToken,
) (*domain.VacationToken, error) {
	return s.r.Create(ctx, t)
}

func (s *refreshTokenService) DeleteAll(ctx context.Context) error {
	return s.r.DeleteAll(ctx)
}

func (s *refreshTokenService) ExistsForUser(
	ctx context.Context,
	userId int64,
	year int,
) (bool, error) {
	return s.r.ExistsForUser(ctx, userId, year)
}
