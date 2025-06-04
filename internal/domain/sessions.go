package domain

import (
	"context"
	"time"
)

type Session struct {
	ID         string    `json:"id"`
	ValidUntil time.Time `json:"valid_until"`
	CreatedAt  time.Time `json:"created_at"`
	EditedAt   time.Time `json:"edited_at"`
	UserID     int64     `json:"user_id"`
}

type SessionRepository interface {
	Create(ctx context.Context, userId int64, secureRand string, validUntil time.Time) (*Session, error)
	Delete(ctx context.Context, cookie string) error
	DeleteAll(ctx context.Context) error
	IsValidSession(ctx context.Context, cookie string, timestamp time.Time) bool
	GetUserFromSession(ctx context.Context, cookie string) (*User, error)
	GetById(ctx context.Context, cookie string) (*Session, error)
}
