package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"chrono/internal/domain"
)

type RequestService interface {
	Create(
		ctx context.Context,
		msg string,
		user *domain.User,
		event *domain.Event,
	) (*domain.Request, error)
	GetPending(ctx context.Context) ([]domain.BatchRequest, error)
	GetInRange(ctx context.Context, userId int64, start, end time.Time) ([]domain.Request, error)
	UpdateInRange(
		ctx context.Context,
		editor int64,
		form domain.PatchRequestForm,
	) (int64, error)
	GetEventName(ctx context.Context, reqId int64) (string, error)
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

func (svc *requestService) GetPending(ctx context.Context) ([]domain.BatchRequest, error) {
	req, err := svc.request.GetPending(ctx)
	if err != nil {
		return nil, err
	}

	requestsToShow := []domain.BatchRequest{}

	startIndex := 0
	for startIndex < len(req) {
		endIndex := startIndex

		for endIndex+1 < len(req) &&
			req[endIndex].ScheduledAt.Year() == req[endIndex+1].ScheduledAt.Year() &&
			req[endIndex].ScheduledAt.YearDay()+1 == req[endIndex+1].ScheduledAt.YearDay() &&
			req[endIndex].Name == req[endIndex+1].Name &&
			req[endIndex].UserID == req[endIndex+1].UserID {
			endIndex++
		}

		startDate := req[startIndex].ScheduledAt
		endDate := req[endIndex].ScheduledAt

		confilctingUsers, err := svc.user.GetConflicting(
			ctx,
			req[startIndex].UserID,
			startDate,
			endDate,
		)
		if err != nil {
			return nil, err
		}

		requestsToShow = append(requestsToShow, domain.BatchRequest{
			StartDate:  startDate,
			EndDate:    endDate,
			EventCount: endIndex - startIndex + 1,
			Request:    &req[endIndex],
			Conflicts:  &confilctingUsers,
		})

		startIndex = endIndex + 1
	}

	return requestsToShow, nil
}

func (svc *requestService) GetEventNameFrom(
	ctx context.Context,
	req int64,
) (string, error) {
	return svc.request.GetEventNameFrom(ctx, req)
}

func (svc *requestService) GetInRange(
	ctx context.Context,
	userId int64,
	start, end time.Time,
) ([]domain.Request, error) {
	return svc.request.GetInRange(ctx, userId, start, end)
}

func (svc *requestService) UpdateInRange(
	ctx context.Context,
	editorId int64,
	form domain.PatchRequestForm,
) (int64, error) {
	reqId, err := svc.request.UpdateInRange(
		ctx,
		form.State,
		editorId,
		form.UserID,
		form.StartDate,
		form.EndDate,
	)
	if err != nil {
		return 0, err
	}

	editor, err := svc.user.GetById(ctx, editorId)
	if err != nil {
		return 0, err
	}

	msg := fmt.Sprintf("%v %v your request.", editor.Username, form.State)
	if form.Reason != "" {
		msg = fmt.Sprintf("%v %v your request: %v.", editor.Username, form.State, form.Reason)
	}

	err = svc.notif.CreateAndNotify(ctx, msg, []domain.User{{ID: form.UserID}})
	if err != nil {
		return 0, err
	}

	return reqId, nil
}

func (svc *requestService) GetEventName(ctx context.Context, reqId int64) (string, error) {
	return svc.request.GetEventNameFrom(ctx, reqId)
}
