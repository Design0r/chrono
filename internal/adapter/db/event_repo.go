package db

import (
	"context"
	"log/slog"
	"time"

	"chrono/db/repo"
	"chrono/internal/domain"
)

type SQLEventRepo struct {
	r   *repo.Queries
	log *slog.Logger
}

func NewSQLEventUserRepo(r *repo.Queries, log *slog.Logger) SQLEventRepo {
	return SQLEventRepo{r: r, log: log}
}

func (r *SQLEventRepo) Create(
	ctx context.Context,
	data domain.YMDDate,
	eventType string,
	user *domain.User,
) (*domain.Event, error) {
	date := time.Date(
		data.Year,
		time.Month(data.Month),
		data.Day,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	evt := domain.Event{Name: eventType}
	state := "pending"
	if !evt.IsVacation() || user.IsSuperuser {
		state = "accepted"
	}

	event, err := r.r.CreateEvent(
		ctx,
		repo.CreateEventParams{Name: eventType, UserID: user.ID, ScheduledAt: date, State: state},
	)
	if err != nil {
		r.log.Error(
			"CreateEvent failed",
			slog.String("eventType", eventType),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return (*domain.Event)(&event), nil
}

func (r *SQLEventRepo) Update(
	ctx context.Context,
	event *domain.Event,
) (*domain.Event, error) {
	e, err := r.r.UpdateEvent(
		ctx,
		repo.UpdateEventParams{
			Name:        event.Name,
			ID:          event.ID,
			ScheduledAt: event.ScheduledAt,
			State:       event.State,
		},
	)
	if err != nil {
		r.log.Error(
			"UpdateEvent failed",
			slog.String("eventType", event.Name),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return (*domain.Event)(&e), nil
}

func (r *SQLEventRepo) Delete(
	ctx context.Context,
	eventId int64,
) (*domain.Event, error) {
	e, err := r.r.DeleteEvent(ctx, eventId)
	if err != nil {
		r.log.Error(
			"DeleteEvent failed",
			slog.Int64("eventId", eventId),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return (*domain.Event)(&e), nil
}

func (r *SQLEventRepo) GetForDay(
	ctx context.Context,
	data domain.YMDDate,
) ([]domain.Event, error) {
	date := time.Date(
		data.Year,
		time.Month(data.Month),
		data.Day,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	events, err := r.r.GetEventsForDay(context.Background(), date)
	if err != nil {
		r.log.Error("GetEventsForDay failed", slog.String("error", err.Error()))
		return []domain.Event{}, err
	}

	e := make([]domain.Event, len(events))
	for i := 0; i <= len(events); i++ {
		e[i] = (domain.Event)(events[i])
	}

	return e, nil
}
