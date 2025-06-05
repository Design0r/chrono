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
	refresh domain.RefreshTokenRepository
	log     *slog.Logger
}

func NewRefreshTokenService(
	r domain.RefreshTokenRepository,
	log *slog.Logger,
) refreshTokenService {
	return refreshTokenService{refresh: r, log: log}
}

func (svc *refreshTokenService) Create(
	ctx context.Context,
	t domain.CreateVacationToken,
) (*domain.VacationToken, error) {
	return svc.refresh.Create(ctx, t)
}

func (svc *refreshTokenService) DeleteAll(ctx context.Context) error {
	return svc.refresh.DeleteAll(ctx)
}

func (svc *refreshTokenService) ExistsForUser(
	ctx context.Context,
	userId int64,
	year int,
) (bool, error) {
	return svc.refresh.ExistsForUser(ctx, userId, year)
}
