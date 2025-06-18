package db

import (
	"context"
	"log/slog"

	"chrono/db/repo"
	"chrono/internal/domain"
)

type SQLNotificationRepo struct {
	r   *repo.Queries
	log *slog.Logger
}

type SQLUserNotificationRepo struct {
	r   *repo.Queries
	log *slog.Logger
}

func NewSQLNotificationRepo(r *repo.Queries, log *slog.Logger) domain.NotificationRepository {
	return &SQLNotificationRepo{r: r, log: log}
}

func NewSQLUserNotificationRepo(r *repo.Queries, log *slog.Logger) domain.NotificationUserRepository {
	return &SQLUserNotificationRepo{r: r, log: log}
}

func (r *SQLNotificationRepo) Create(ctx context.Context, msg string) (domain.Notification, error) {
	notif, err := r.r.CreateNotification(ctx, msg)
	if err != nil {
		r.log.Error("CreateNotification failed", slog.String("error", err.Error()))
		return domain.Notification{}, err
	}

	return (domain.Notification)(notif), nil
}

func (r *SQLNotificationRepo) Update(ctx context.Context, n domain.Notification) error {
	_, err := r.r.UpdateNotification(ctx, n.Message)
	if err != nil {
		r.log.Error("UpdateNotification failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *SQLNotificationRepo) Clear(ctx context.Context, notifId int64) error {
	_, err := r.r.ClearNotification(ctx, notifId)
	if err != nil {
		r.log.Error("ClearNotification failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *SQLNotificationRepo) ClearAll(ctx context.Context, userId int64) error {
	err := r.r.ClearAllUserNotifications(ctx, userId)
	if err != nil {
		r.log.Error("ClearAllUserNotifications failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *SQLUserNotificationRepo) Create(
	ctx context.Context,
	userId int64,
	notifId int64,
) error {
	params := repo.CreateNotificationUserParams{UserID: userId, NotificationID: notifId}
	err := r.r.CreateNotificationUser(ctx, params)
	if err != nil {
		r.log.Error("CreateNotificationUser failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *SQLUserNotificationRepo) GetByUserId(
	ctx context.Context,
	userId int64,
) ([]domain.Notification, error) {
	notif, err := r.r.GetUserNotifications(ctx, userId)
	if err != nil {
		r.log.Error("GetUserNotifications failed", slog.String("error", err.Error()))
		return []domain.Notification{}, err
	}

	n := make([]domain.Notification, len(notif))
	for i := range notif {
		n[i] = (domain.Notification)(notif[i])
	}

	return n, nil
}

func (r *SQLUserNotificationRepo) UpdateByUserId(
	ctx context.Context,
	userId int64,
) error {
	err := r.r.ClearAllUserNotifications(ctx, userId)
	if err != nil {
		r.log.Error(
			"ClearAllUserNotifications failed",
			slog.Int64("userId", userId),
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}
