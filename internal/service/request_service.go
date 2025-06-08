package service

import (
	"context"
	"log/slog"

	"chrono/internal/domain"
)

type RequestService interface {
	Create(
		ctx context.Context,
		msg string,
		user *domain.User,
		event *domain.Event,
	) (*domain.Request, error)
	GetPending(ctx context.Context) ([]domain.RequestEventUser, error)
}

type requestService struct {
	request domain.RequestRepository
	notif   NotificationService
	user    domain.UserRepository
	log     *slog.Logger
}

func NewRequestService(
	r domain.RequestRepository,
	u domain.UserRepository,
	n NotificationService,
	log *slog.Logger,
) requestService {
	return requestService{request: r, notif: n, log: log, user: u}
}

func (svc *requestService) Create(
	ctx context.Context,
	msg string,
	user *domain.User,
	event *domain.Event,
) (*domain.Request, error) {
	req, err := svc.request.Create(ctx, msg, user, event)
	if err != nil {
		return nil, err
	}

	admins, err := svc.user.GetAdmins(ctx)
	if err != nil {
		return nil, err
	}

	err = svc.notif.CreateAndNotify(ctx, msg, admins)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (svc *requestService) GetPending(ctx context.Context) ([]domain.RequestEventUser, error) {
	return svc.request.GetPending(ctx)
}

func (svc *requestService) GetEventNameFrom(
	ctx context.Context,
	req *domain.Request,
) (string, error) {
	return svc.request.GetEventNameFrom(ctx, req)
}
