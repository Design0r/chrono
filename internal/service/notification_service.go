package service

import (
	"context"
	"log/slog"

	"chrono/internal/domain"
)

type NotificationService interface {
	Create(ctx context.Context, msg string) (domain.Notification, error)
	CreateAndNotify(ctx context.Context, msg string, users []domain.User) error
	NotifyUser(ctx context.Context, user *domain.User, notif domain.Notification) error
	NotifyUsers(ctx context.Context, users []domain.User, notif domain.Notification) error
}

type notificationService struct {
	notif     domain.NotificationRepository
	userNotif domain.NotificationUserRepository
	log       *slog.Logger
}

func NewNotificationService(
	n domain.NotificationRepository,
	un domain.NotificationUserRepository,
	log *slog.Logger,
) notificationService {
	return notificationService{notif: n, userNotif: un, log: log}
}

func (svc *notificationService) Create(
	ctx context.Context,
	msg string,
) (domain.Notification, error) {
	return svc.notif.Create(ctx, msg)
}

func (svc *notificationService) CreateAndNotify(
	ctx context.Context,
	msg string,
	users []domain.User,
) error {
	notif, err := svc.Create(ctx, msg)
	if err != nil {
		return err
	}

	for _, u := range users {
		err = svc.NotifyUser(ctx, &u, notif)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *notificationService) NotifyUser(
	ctx context.Context,
	user *domain.User,
	notif domain.Notification,
) error {
	return svc.userNotif.Create(ctx, user.ID, notif.ID)
}

func (svc *notificationService) NotifyUsers(
	ctx context.Context,
	users []domain.User,
	notif domain.Notification,
) error {
	for _, u := range users {
		err := svc.NotifyUser(ctx, &u, notif)
		if err != nil {
			return err
		}
	}

	return nil
}
