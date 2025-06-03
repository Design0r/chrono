package db

import (
	"chrono/db/repo"
	"chrono/internal/domain"
	"context"
	"log/slog"
)

type SQLSettingsRepo struct {
	q   *repo.Queries
	log *slog.Logger
}

func NewSQLSettingsRepo(q *repo.Queries) SQLSettingsRepo {
	return SQLSettingsRepo{q: q}
}

func (r *SQLSettingsRepo) GetById(ctx context.Context, id int64) (domain.Settings, error) {
	s, err := r.q.GetSettingsById(ctx, id)
	if err != nil {
		r.log.Error("repo.GetSettingsById failed:", slog.String("error", err.Error()))
		return domain.Settings{}, nil
	}

	return domain.Settings{
		SignupEnabled: s.SignupEnabled,
	}, nil

}

func (r *SQLSettingsRepo) Delete(ctx context.Context, id int64) error {
	err := r.q.DeleteSettings(ctx, id)
	if err != nil {
		r.log.Error("repo.GetSettingsById failed:", slog.String("error", err.Error()))
		return err
	}

	return nil

}

func (r *SQLSettingsRepo) Create(ctx context.Context, s domain.Settings) (domain.Settings, error) {
	s, err := r.q.CreateSettings(ctx, s.SignupEnabled)
	if err != nil {
		r.log.Error("repo.GetSettingsById failed:", slog.String("error", err.Error()))
		return domain.Settings{}, nil
	}

	return domain.Settings{
		SignupEnabled: s.SignupEnabled,
	}, nil

}

func (r *SQLSettingsRepo) Update(ctx context.Context, s domain.Settings) (domain.Settings, error) {
	s, err := r.q.UpdateSettings(ctx, repo.UpdateSettingsParams{ID: s.ID, SignupEnabled: s.SignupEnabled})
	if err != nil {
		r.log.Error("repo.GetSettingsById failed:", slog.String("error", err.Error()))
		return domain.Settings{}, nil
	}

	return domain.Settings{
		SignupEnabled: s.SignupEnabled,
	}, nil

}
