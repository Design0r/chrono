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
) *TokenService {
	return &TokenService{refresh: r, vac: v, log: log}
}

// InitYearlyTokens grants a user their yearly allowance, once per year. The
// refresh token is what records that the year has been dealt with, so it is
// written only after the grant succeeded: marking the year first would swallow
// the allowance for good if the grant then failed.
func (svc *TokenService) InitYearlyTokens(ctx context.Context, user *domain.User, year int) error {
	exists, err := svc.RefreshTokenExistsForUser(ctx, user.ID, year)
	if err != nil {
		svc.log.Error("failed to get refresh token")
		return err
	}

	if exists {
		return nil
	}

	if user.VacationDays > 0 {
		_, err = svc.CreateVacationToken(ctx, float64(user.VacationDays), year, user.ID)
		if err != nil {
			svc.log.Error("failed to create vac tokens")
			return err
		}
	}

	_, err = svc.CreateRefreshToken(ctx, year, user.ID)
	if err != nil {
		svc.log.Error("failed to create refresh token")
		return err
	}

	return nil
}

// UpdateYearlyTokens books a change of a user's yearly allowance as a delta.
// The user passed in must be the state *before* the change, because the year's
// base grant is materialised first: without that the delta would be booked on
// top of nothing, and InitYearlyTokens would never run for that year again.
func (svc *TokenService) UpdateYearlyTokens(
	ctx context.Context,
	user *domain.User,
	delta, year int,
) error {
	err := svc.InitYearlyTokens(ctx, user, year)
	if err != nil {
		return err
	}

	if delta == 0 {
		return nil
	}

	_, err = svc.CreateVacationToken(ctx, float64(delta), year, user.ID)
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
	end := start.AddDate(1, 3, 0)
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
