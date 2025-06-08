package domain

import (
	"context"
	"time"
)

type Notification struct {
	ID        int64      `json:"id"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	ViewedAt  *time.Time `json:"viewed_at"`
}

type NotificationUser struct {
	NotificationID int64 `json:"notification_id"`
	UserID         int64 `json:"user_id"`
}

type NotificationRepository interface {
	Create(ctx context.Context, msg string) (Notification, error)
	Update(ctx context.Context, n Notification) error
}

type NotificationUserRepository interface {
	Create(ctx context.Context, msg string) (Notification, error)
	GetByUserId(ctx context.Context, userId int64) ([]Notification, error)
	UpdateByUserId(ctx context.Context, userId int64) error
}
