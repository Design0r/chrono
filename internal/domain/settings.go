package domain

import (
	"context"
)

type Settings struct {
	ID            int64 `json:"id"`
	SignupEnabled bool  `json:"signup_enabled" form:"signup_enabled"`
}

type SettingsPatch struct {
	SignupEnabled string `query:"signup_enabled"`
}

type SettingsRepository interface {
	GetById(ctx context.Context, id int64) (Settings, error)
	Create(ctx context.Context, s Settings) (Settings, error)
	Update(ctx context.Context, s Settings) (Settings, error)
	Delete(ctx context.Context, id int64) error
}
