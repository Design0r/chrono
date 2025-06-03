package service

import (
	"chrono/internal/domain"
	"context"
)

type SettingsService interface {
	Create(ctx context.Context, s domain.Settings) (domain.Settings, error)
	Update(ctx context.Context, s domain.Settings) (domain.Settings, error)
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (domain.Settings, error)
	GetFirst(ctx context.Context) (domain.Settings, error)
}

type settingsService struct {
	r domain.SettingsRepository
}

func NewSettingsService(r domain.SettingsRepository) settingsService {
	return settingsService{r: r}
}

func (s *settingsService) Create(ctx context.Context, settings domain.Settings) (domain.Settings, error) {
	return s.r.Create(ctx, settings)
}

func (s *settingsService) Update(ctx context.Context, settings domain.Settings) (domain.Settings, error) {
	return s.r.Update(ctx, settings)
}

func (s *settingsService) Delete(ctx context.Context, id int64) error {
	return s.r.Delete(ctx, id)
}

func (s *settingsService) GetById(ctx context.Context, id int64) (domain.Settings, error) {
	return s.r.GetById(ctx, id)
}

func (s *settingsService) GetFirst(ctx context.Context) (domain.Settings, error) {
	return s.r.GetById(ctx, 1)
}
