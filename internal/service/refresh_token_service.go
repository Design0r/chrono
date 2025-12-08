package service

import (
	"context"
	"log/slog"

	"chrono/internal/domain"
)

type RefreshTokenService struct {
	refresh domain.RefreshTokenRepository
	log     *slog.Logger
}

func NewRefreshTokenService(
	r domain.RefreshTokenRepository,
	log *slog.Logger,
) RefreshTokenService {
	return RefreshTokenService{refresh: r, log: log}
}

func (svc *RefreshTokenService) Create(
	ctx context.Context,
	year int,
	userId int64,
) (*domain.RefreshToken, error) {
	return svc.refresh.Create(ctx, year, userId)
}

func (svc *RefreshTokenService) DeleteAll(ctx context.Context) error {
	return svc.refresh.DeleteAll(ctx)
}

func (svc *RefreshTokenService) ExistsForUser(
	ctx context.Context,
	userId int64,
	year int,
) (bool, error) {
	return svc.refresh.ExistsForUser(ctx, userId, year)
}

func (svc *RefreshTokenService) CreateIfNotExists(
	ctx context.Context,
	userId int64,
	year int,
) (bool, error) {
	exists, err := svc.refresh.ExistsForUser(ctx, userId, year)
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	_, err = svc.Create(ctx, year, userId)
	if err != nil {
		return false, err
	}

	return false, nil
}
