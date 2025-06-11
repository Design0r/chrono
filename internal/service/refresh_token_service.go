package service

import (
	"context"
	"log/slog"

	"chrono/internal/domain"
)

type RefreshTokenService interface {
	Create(ctx context.Context, year int, userId int64) (*domain.RefreshToken, error)
	CreateIfNotExists(ctx context.Context, userId int64) (bool, error)
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
	year int,
	userId int64,
) (*domain.RefreshToken, error) {
	return svc.refresh.Create(ctx, year, userId)
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

func (svc *refreshTokenService) CreateIfNotExists(ctx context.Context, userId int64) (bool, error) {
	currYear := domain.CurrentYear()
	exists, err := svc.refresh.ExistsForUser(ctx, userId, currYear)
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	_, err = svc.Create(ctx, currYear, userId)
	if err != nil {
		return false, err
	}

	return true, nil
}
