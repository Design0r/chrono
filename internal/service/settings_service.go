package service

import (
	"context"
	"log/slog"

	"chrono/internal/domain"
)

type SettingsService interface {
	Create(ctx context.Context, s domain.Settings) (domain.Settings, error)
	Update(ctx context.Context, s domain.Settings) (domain.Settings, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (domain.Settings, error)
	GetFirst(ctx context.Context) (domain.Settings, error)
	Init(ctx context.Context, s domain.Settings) (domain.Settings, error)
}

type settingsService struct {
	settings domain.SettingsRepository
	log      *slog.Logger
}

func NewSettingsService(r domain.SettingsRepository, log *slog.Logger) settingsService {
	return settingsService{settings: r, log: log}
}

func (svc *settingsService) Create(
	ctx context.Context,
	settings domain.Settings,
) (domain.Settings, error) {
	return svc.settings.Create(ctx, settings)
}

func (svc *settingsService) Update(
	ctx context.Context,
	settings domain.Settings,
) (domain.Settings, error) {
	return svc.settings.Update(ctx, settings)
}

func (svc *settingsService) Delete(ctx context.Context, id int64) error {
	return svc.settings.Delete(ctx, id)
}

func (svc *settingsService) GetById(ctx context.Context, id int64) (domain.Settings, error) {
	return svc.settings.GetById(ctx, id)
}

func (svc *settingsService) GetFirst(ctx context.Context) (domain.Settings, error) {
	return svc.settings.GetById(ctx, 1)
}
func (svc *settingsService) Init(
	ctx context.Context,
	settings domain.Settings,
) (domain.Settings, error) {
	s, err := svc.GetFirst(ctx)
	if err != nil {
		s, err := svc.settings.Create(ctx, settings)
		if err != nil {
			return domain.Settings{}, err
		}
		svc.log.Info("Initialized new settings.", "settings", s)
	}

	svc.log.Info("Initialized existing settings.", "settings", s)
	return s, nil
}
