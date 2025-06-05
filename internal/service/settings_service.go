package service

import (
	"context"

	"chrono/internal/domain"
)

type SettingsService interface {
	Create(ctx context.Context, s domain.Settings) (domain.Settings, error)
	Update(ctx context.Context, s domain.Settings) (domain.Settings, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (domain.Settings, error)
	GetFirst(ctx context.Context) (domain.Settings, error)
}

type settingsService struct {
	settings domain.SettingsRepository
}

func NewSettingsService(r domain.SettingsRepository) settingsService {
	return settingsService{settings: r}
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
