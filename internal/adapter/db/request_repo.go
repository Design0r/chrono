package db

import (
	"context"
	"log/slog"
	"time"

	"chrono/db/repo"
	"chrono/internal/domain"
)

type SQLRequestRepo struct {
	r   *repo.Queries
	log *slog.Logger
}

func NewSQLRequestRepo(r *repo.Queries, log *slog.Logger) domain.RequestRepository {
	return &SQLRequestRepo{r: r, log: log}
}

func (r *SQLRequestRepo) Create(
	ctx context.Context,
	msg string,
	user *domain.User,
	event *domain.Event,
) (*domain.Request, error) {
	params := repo.CreateRequestParams{
		Message: &msg,
		State:   "pending",
		UserID:  user.ID,
		EventID: event.ID,
	}

	request, err := r.r.CreateRequest(ctx, params)
	if err != nil {
		r.log.Error(
			"CreateRequest failed",
			slog.String("username", user.Username),
			slog.String("event", event.Name),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return (*domain.Request)(&request), nil
}

func (r *SQLRequestRepo) Update(
	ctx context.Context,
	editor *domain.User,
	req *domain.Request,
) (*domain.Request, error) {
	params := repo.UpdateRequestParams{
		Message: req.Message, State: req.State, EditedBy: &editor.ID,
		EventID: req.EventID,
	}
	request, err := r.r.UpdateRequest(ctx, params)
	if err != nil {
		r.log.Error(
			"UpdateRequest failed",
			slog.String("request", *req.Message),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return (*domain.Request)(&request), nil
}

func (r *SQLRequestRepo) GetPending(ctx context.Context) ([]domain.RequestEventUser, error) {
	result, err := r.r.GetPendingRequests(ctx)
	if err != nil {
		r.log.Error(
			"GetPendingRequests failed",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	requests := make([]domain.RequestEventUser, len(result))
	for i := range result {
		requests[i] = (domain.RequestEventUser)(result[i])
	}

	return requests, nil
}

func (r *SQLRequestRepo) GetEventNameFrom(ctx context.Context, reqId int64) (string, error) {
	result, err := r.r.GetEventNameFromRequest(ctx, reqId)
	if err != nil {
		r.log.Error(
			"GetPendingRequests failed",
			slog.String("error", err.Error()),
		)
		return "", err
	}

	return result, nil
}

func (r *SQLRequestRepo) GetInRange(
	ctx context.Context,
	userId int64,
	start, end time.Time,
) ([]domain.Request, error) {
	params := repo.GetRequestRangeParams{UserID: userId, ScheduledAt: start, ScheduledAt_2: end}
	e, err := r.r.GetRequestRange(ctx, params)
	if err != nil {
		r.log.Error(
			"GetRequestRange failed",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	events := make([]domain.Request, len(e))
	for i := range e {
		events[i] = (domain.Request)(e[i])
	}

	return events, nil
}

func (r *SQLRequestRepo) UpdateInRange(
	ctx context.Context,
	state string,
	editor, reqUserId int64,
	start, end time.Time,
) (int64, error) {
	params := repo.UpdateRequestStateRangeParams{
		UserID:        reqUserId,
		EditedBy:      &editor,
		State:         state,
		ScheduledAt:   start,
		ScheduledAt_2: end,
	}

	reqId, err := r.r.UpdateRequestStateRange(context.Background(), params)
	if err != nil {
		r.log.Error(
			"UpdateRequestStateRange failed",
			slog.String("error", err.Error()),
		)
	}

	return reqId, nil
}
