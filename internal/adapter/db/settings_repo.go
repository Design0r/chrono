package db

import (
	"context"
	"log/slog"

	"chrono/db/repo"
	"chrono/internal/domain"
)

type SQLSettingsRepo struct {
	q   repo.Querier
	log *slog.Logger
}

func NewSQLSettingsRepo(q repo.Querier, log *slog.Logger) domain.SettingsRepository {
	return &SQLSettingsRepo{q: q, log: log}
}

func (r *SQLSettingsRepo) GetById(ctx context.Context, id int64) (domain.Settings, error) {
	s, err := r.q.GetSettingsById(ctx, id)
	if err != nil {
		r.log.Error("repo.GetSettingsById failed:", slog.String("error", err.Error()))
		return domain.Settings{}, err
	}

	return (domain.Settings)(s), nil
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
	settings, err := r.q.CreateSettings(ctx, s.SignupEnabled)
	if err != nil {
		r.log.Error("repo.GetSettingsById failed:", slog.String("error", err.Error()))
		return domain.Settings{}, err
	}

	return (domain.Settings)(settings), nil
}

func (r *SQLSettingsRepo) Update(ctx context.Context, s domain.Settings) (domain.Settings, error) {
	settings, err := r.q.UpdateSettings(
		ctx,
		repo.UpdateSettingsParams{ID: s.ID, SignupEnabled: s.SignupEnabled},
	)
	if err != nil {
		r.log.Error("repo.GetSettingsById failed:", slog.String("error", err.Error()))
		return domain.Settings{}, err
	}

	return (domain.Settings)(settings), nil
}
