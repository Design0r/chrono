package domain

import (
	"context"
	"time"
)

type RefreshToken struct {
	ID        int64     `json:"id"`
	Year      int64     `json:"year"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, year int, userId int64) (*RefreshToken, error)
	DeleteAll(ctx context.Context) error
	ExistsForUser(ctx context.Context, userId int64, year int) (bool, error)
}
