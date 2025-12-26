package service

import (
	"context"
	"log/slog"

	"chrono/internal/domain"
)

type NotificationService struct {
	notif     domain.NotificationRepository
	userNotif domain.NotificationUserRepository
	log       *slog.Logger
}

func NewNotificationService(
	n domain.NotificationRepository,
	un domain.NotificationUserRepository,
	log *slog.Logger,
) *NotificationService {
	return &NotificationService{notif: n, userNotif: un, log: log}
}

func (svc *NotificationService) Create(
	ctx context.Context,
	msg string,
) (domain.Notification, error) {
	return svc.notif.Create(ctx, msg)
}

func (svc *NotificationService) CreateAndNotify(
	ctx context.Context,
	msg string,
	users []domain.User,
) error {
	notif, err := svc.Create(ctx, msg)
	if err != nil {
		return err
	}

	for _, u := range users {
		err = svc.NotifyUser(ctx, u.ID, notif)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *NotificationService) NotifyUser(
	ctx context.Context,
	user int64,
	notif domain.Notification,
) error {
	return svc.userNotif.Create(ctx, user, notif.ID)
}

func (svc *NotificationService) NotifyUsers(
	ctx context.Context,
	users []domain.User,
	notif domain.Notification,
) error {
	for _, u := range users {
		err := svc.NotifyUser(ctx, u.ID, notif)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *NotificationService) GetByUserId(
	ctx context.Context,
	userId int64,
) ([]domain.Notification, error) {
	return svc.userNotif.GetByUserId(ctx, userId)
}

func (svc *NotificationService) Clear(ctx context.Context, notifId int64) error {
	return svc.notif.Clear(ctx, notifId)
}

func (svc *NotificationService) ClearAll(ctx context.Context, userId int64) error {
	return svc.notif.ClearAll(ctx, userId)
}
