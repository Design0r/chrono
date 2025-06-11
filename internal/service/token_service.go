package service

import (
	"context"
	"log/slog"
	"time"

	"chrono/internal/domain"
)

type TokenService interface {
	InitYearlyTokens(ctx context.Context, user *domain.User) error
	UpdateYearlyTokens(ctx context.Context, userId int64, vacation, year int) error
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

func (svc *tokenService) InitYearlyTokens(ctx context.Context, user *domain.User) error {
	exists, err := svc.refresh.CreateIfNotExists(ctx, user.ID)
	if err != nil {
		return err
	}

	if exists || user.VacationDays <= 0 {
		return nil
	}

	currYear := domain.CurrentYear()
	start := time.Date(currYear, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 3, 0)

	_, err = svc.vac.Create(
		ctx,
		domain.CreateVacationToken{
			StartDate: start,
			EndDate:   end,
			Value:     float64(user.VacationDays),
			UserID:    user.ID,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (svc *tokenService) UpdateYearlyTokens(
	ctx context.Context,
	userId int64,
	vacation, year int,
) error {
	_, err := svc.refresh.CreateIfNotExists(ctx, userId)
	if err != nil {
		return err
	}

	start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 3, 0)
	params := domain.CreateVacationToken{
		StartDate: start,
		EndDate:   end,
		UserID:    userId,
		Value:     float64(vacation),
	}
	_, err = svc.vac.Create(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
