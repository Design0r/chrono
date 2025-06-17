package db

import (
	"context"
	"log/slog"
	"time"

	"chrono/db/repo"
	"chrono/internal/domain"
)

type SQLVacationTokenRepo struct {
	q   *repo.Queries
	log *slog.Logger
}

func NewSQLVacationTokenRepo(q *repo.Queries, log *slog.Logger) domain.VacationTokenRepository {
	return &SQLVacationTokenRepo{q: q, log: log}
}

func (r *SQLVacationTokenRepo) Create(
	ctx context.Context,
	t domain.CreateVacationToken,
) (*domain.VacationToken, error) {
	params := repo.CreateVacationTokenParams{
		UserID:    t.UserID,
		Value:     t.Value,
		StartDate: t.StartDate,
		EndDate:   t.EndDate,
	}
	token, err := r.q.CreateVacationToken(ctx, params)
	if err != nil {
		r.log.Error(
			"CreateVacationToken failed",
			slog.Int64("user_id", t.UserID), slog.String("error", err.Error()))

		return nil, err
	}

	return &domain.VacationToken{
		ID:        token.ID,
		StartDate: token.StartDate,
		EndDate:   token.EndDate,
		Value:     token.Value,
		UserID:    token.UserID,
	}, nil
}

func (r *SQLVacationTokenRepo) Delete(ctx context.Context, id int64) error {
	err := r.q.DeleteVacationToken(ctx, id)
	if err != nil {
		r.log.Error(
			"DeleteVacationToken failed",
			slog.String("error", err.Error()))

		return err
	}

	return nil
}

func (r *SQLVacationTokenRepo) DeleteAll(ctx context.Context) error {
	err := r.q.DeleteAllVacationTokens(ctx)
	if err != nil {
		r.log.Error(
			"DeleteAllVacationTokens failed",
			slog.String("error", err.Error()))

		return err
	}

	return nil
}

func (r *SQLVacationTokenRepo) GetRemainingVacationForUser(
	ctx context.Context,
	userId int64,
	start time.Time,
	end time.Time,
) (float64, error) {
	params := repo.GetRemainingVacationForUserParams{UserID: userId, StartDate: start, EndDate: end}
	vac, err := r.q.GetRemainingVacationForUser(ctx, params)
	if err != nil {
		r.log.Error(
			"GetRemainingVacationForUser failed",
			slog.Int64("user_id", userId),
			slog.String("error", err.Error()))

		return 0, err
	}

	if vac == nil {
		return 0, nil
	}

	return *vac, nil
}
