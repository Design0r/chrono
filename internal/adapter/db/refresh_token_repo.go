package db

import (
	"context"
	"log/slog"

	"chrono/db/repo"
	"chrono/internal/domain"
)

type SQLRefreshTokenRepo struct {
	q   *repo.Queries
	log *slog.Logger
}

func NewSQLRefreshTokenRepo(q *repo.Queries, log *slog.Logger) domain.RefreshTokenRepository {
	return &SQLRefreshTokenRepo{q: q, log: log}
}

func (r *SQLRefreshTokenRepo) Create(ctx context.Context, year int, userId int64) (*domain.RefreshToken, error) {
	params := repo.CreateRefreshTokenParams{UserID: userId, Year: int64(year)}
	token, err := r.q.CreateRefreshToken(ctx, params)
	if err != nil {
		r.log.Error(
			"CreateRefreshToken failed",
			slog.Int64("user_id", userId),
			slog.String("error", err.Error()))

		return nil, err
	}

	return &domain.RefreshToken{
		ID:     token.ID,
		UserID: token.UserID,
		Year:   token.Year,
	}, nil
}

func (r *SQLRefreshTokenRepo) DeleteAll(ctx context.Context) error {
	err := r.q.DeleteAllRefreshTokens(ctx)
	if err != nil {
		r.log.Error(
			"DeleteAllRefreshTokens failed",
			slog.String("error", err.Error()))

		return err
	}

	return nil
}

func (r *SQLRefreshTokenRepo) ExistsForUser(
	ctx context.Context,
	userId int64,
	year int,
) (bool, error) {
	params := repo.GetRefreshTokenParams{UserID: userId, Year: int64(year)}
	tokenCount, err := r.q.GetRefreshToken(ctx, params)
	if err != nil {
		r.log.Error(
			"GetRefreshToken failed",
			slog.Int64("user_id", userId),
			slog.String("error", err.Error()))

		return false, err
	}

	return tokenCount > 0, nil
}
