package service

import (
	"context"
	"log/slog"
	"time"

	"chrono/internal/domain"
)

type TokenService struct {
	refresh domain.RefreshTokenRepository
	vac     domain.VacationTokenRepository
	log     *slog.Logger
}

func NewTokenService(
	r domain.RefreshTokenRepository,
	v domain.VacationTokenRepository,
	log *slog.Logger,
) TokenService {
	return TokenService{refresh: r, vac: v, log: log}
}

func (svc *TokenService) InitYearlyTokens(ctx context.Context, user *domain.User, year int) error {
	exists, err := svc.CreateRefreshTokenIfNotExists(ctx, user.ID, year)
	if err != nil {
		svc.log.Error("failed to get refresh token")
		return err
	}

	if exists || user.VacationDays <= 0 {
		return nil
	}

	_, err = svc.CreateVacationToken(ctx, float64(user.VacationDays), year, user.ID)
	if err != nil {
		svc.log.Error("failed to create vac tokens")
		return err
	}

	return nil
}

func (svc *TokenService) UpdateYearlyTokens(
	ctx context.Context,
	userId int64,
	vacation, year int,
) error {
	_, err := svc.CreateRefreshTokenIfNotExists(ctx, userId, year)
	if err != nil {
		return err
	}

	_, err = svc.CreateVacationToken(ctx, float64(vacation), year, userId)
	if err != nil {
		return err
	}

	return nil
}

func (svc *TokenService) DeleteAll(ctx context.Context) error {
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

func (svc *TokenService) CreateRefreshToken(
	ctx context.Context,
	year int,
	userId int64,
) (*domain.RefreshToken, error) {
	return svc.refresh.Create(ctx, year, userId)
}

func (svc *TokenService) DeleteAllRefreshToken(ctx context.Context) error {
	return svc.refresh.DeleteAll(ctx)
}

func (svc *TokenService) RefreshTokenExistsForUser(
	ctx context.Context,
	userId int64,
	year int,
) (bool, error) {
	return svc.refresh.ExistsForUser(ctx, userId, year)
}

func (svc *TokenService) CreateRefreshTokenIfNotExists(
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

	_, err = svc.CreateRefreshToken(ctx, year, userId)
	if err != nil {
		return false, err
	}

	return false, nil
}

func (svc *TokenService) CreateVacationToken(
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
	end := start.AddDate(1, 2, 0)
	return svc.vac.Create(
		ctx,
		domain.CreateVacationToken{StartDate: start, EndDate: end, UserID: userId, Value: value},
	)
}

func (svc *TokenService) DeleteVacationToken(ctx context.Context, id int64) error {
	return svc.vac.Delete(ctx, id)
}

func (svc *TokenService) DeleteAllVacationToken(ctx context.Context) error {
	return svc.vac.DeleteAll(ctx)
}

func (svc *TokenService) GetRemainingVacationForUser(
	ctx context.Context,
	userId int64,
	start time.Time,
	end time.Time,
) (float64, error) {
	return svc.vac.GetRemainingVacationForUser(ctx, userId, start, end)
}
