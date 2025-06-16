package service

import (
	"context"
	"log/slog"

	"chrono/internal/domain"
)

type TokenService interface {
	InitYearlyTokens(ctx context.Context, user *domain.User, year int) error
	UpdateYearlyTokens(ctx context.Context, userId int64, vacation, year int) error
	DeleteAll(ctx context.Context) error
}

type tokenService struct {
	refresh RefreshTokenService
	vac     VacationTokenService
	log     *slog.Logger
}

func NewTokenService(
	r RefreshTokenService,
	v VacationTokenService,
	log *slog.Logger,
) tokenService {
	return tokenService{refresh: r, vac: v, log: log}
}

func (svc *tokenService) InitYearlyTokens(ctx context.Context, user *domain.User, year int) error {
	exists, err := svc.refresh.CreateIfNotExists(ctx, user.ID, year)
	if err != nil {
		svc.log.Error("failed to get refresh token")
		return err
	}

	if exists || user.VacationDays <= 0 {
		svc.log.Error("refresh token already exists", "exists", exists)
		return nil
	}

	_, err = svc.vac.Create(ctx, float64(user.VacationDays), year, user.ID)
	if err != nil {
		svc.log.Error("failed to create vac tokens")
		return err
	}

	return nil
}

func (svc *tokenService) UpdateYearlyTokens(
	ctx context.Context,
	userId int64,
	vacation, year int,
) error {
	_, err := svc.refresh.CreateIfNotExists(ctx, userId, year)
	if err != nil {
		return err
	}

	_, err = svc.vac.Create(ctx, float64(vacation), year, userId)
	if err != nil {
		return err
	}

	return nil
}

func (svc *tokenService) DeleteAll(ctx context.Context) error {
	err := svc.vac.DeleteAll(ctx)
	if err != nil {
		return err
	}

	err = svc.refresh.DeleteAll(ctx)
	if err != nil {
		return err
	}

	return nil
}
