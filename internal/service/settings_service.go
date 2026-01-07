package service

import (
	"context"
	"log/slog"

	"chrono/internal/domain"
)

type SettingsService struct {
	settings domain.SettingsRepository
	log      *slog.Logger
}

func NewSettingsService(r domain.SettingsRepository, log *slog.Logger) *SettingsService {
	return &SettingsService{settings: r, log: log}
}

func (svc *SettingsService) Create(
	ctx context.Context,
	settings domain.Settings,
) (domain.Settings, error) {
	return svc.settings.Create(ctx, settings)
}

func (svc *SettingsService) Update(
	ctx context.Context,
	settings domain.Settings,
) (domain.Settings, error) {
	return svc.settings.Update(ctx, settings)
}

func (svc *SettingsService) Delete(ctx context.Context, id int64) error {
	return svc.settings.Delete(ctx, id)
}

func (svc *SettingsService) GetById(ctx context.Context, id int64) (domain.Settings, error) {
	return svc.settings.GetById(ctx, id)
}

func (svc *SettingsService) GetFirst(ctx context.Context) (domain.Settings, error) {
	return svc.settings.GetById(ctx, 1)
}

func (svc *SettingsService) Init(
	ctx context.Context,
	settings domain.Settings,
) (domain.Settings, error) {
	s, err := svc.GetFirst(ctx)
	if err != nil {
		_, err := svc.settings.Create(ctx, settings)
		if err != nil {
			return domain.Settings{}, err
		}
		svc.log.Info("Initialized new settings.")
	}

	svc.log.Info("Using existing settings.")
	return s, nil
}
