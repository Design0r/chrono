package db

import (
	"context"
	"log/slog"
	"strings"
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
	eventId int64,
	state string,
) (*domain.Event, error) {
	e, err := r.r.UpdateEventState(
		ctx,
		repo.UpdateEventStateParams{
			ID:    eventId,
			State: state,
		},
	)
	if err != nil {
		r.log.Error(
			"UpdateEvent failed",
			slog.Int64("eventId", eventId),
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

	events, err := r.r.GetEventsForDay(ctx, date)
	if err != nil {
		r.log.Error("GetEventsForDay failed", slog.String("error", err.Error()))
		return []domain.Event{}, err
	}

	e := make([]domain.Event, len(events))
	for i := range events {
		e[i] = (domain.Event)(events[i])
	}

	return e, nil
}

func (r *SQLEventRepo) GetForMonth(
	ctx context.Context,
	data domain.YMDate,
	botName string,
	userFilter *domain.User,
	eventFilter string,
) (domain.Month, error) {
	date := time.Date(
		data.Year,
		time.Month(data.Month),
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	events, err := r.r.GetEventsForMonth(
		ctx,
		repo.GetEventsForMonthParams{ScheduledAt: date, ScheduledAt_2: date.AddDate(0, 1, 0)},
	)
	if err != nil {
		r.log.Error(
			"GetEventsForMonth failed",
			slog.Int("year", data.Year),
			slog.String("month", time.Month(data.Month).String()),
			slog.String("error", err.Error()),
		)
		return domain.Month{}, err
	}
	month := domain.GetDaysOfMonth(date.Month(), data.Year)

	for _, event := range events {
		idx := event.ScheduledAt.Day() - 1
		if userFilter != nil && event.Username != userFilter.Username &&
			event.Username != botName {
			continue
		}
		if eventFilter != "" && !strings.Contains(event.Name, eventFilter) &&
			eventFilter != "all" &&
			event.Username != botName {
			continue
		}

		newEvent := domain.EventUser{
			User: domain.User{
				ID:           event.ID_2,
				Username:     event.Username,
				Email:        event.Email,
				Password:     event.Password,
				VacationDays: event.VacationDays,
				IsSuperuser:  event.IsSuperuser,
				CreatedAt:    event.CreatedAt_2,
				EditedAt:     event.EditedAt_2,
				Color:        event.Color,
			},
			Event: domain.Event{
				Name:        event.Name,
				ID:          event.ID,
				State:       event.State,
				ScheduledAt: event.ScheduledAt,
				CreatedAt:   event.CreatedAt,
				EditedAt:    event.EditedAt,
				UserID:      event.UserID,
			},
		}
		month.Days[idx].Events = append(month.Days[idx].Events, newEvent)
	}

	return month, nil
}

func (r *SQLEventRepo) GetForYear(
	ctx context.Context,
	year int,
) ([]domain.EventUser, error) {
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)

	params := repo.GetEventsForYearParams{
		ScheduledAt:   yearStart,
		ScheduledAt_2: yearStart.AddDate(1, 0, 0),
	}

	events, err := r.r.GetEventsForYear(ctx, params)
	if err != nil {
		r.log.Error(
			"GetEventsForYear failed",
			slog.Int("year", year),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	e := make([]domain.EventUser, len(events))
	for i, event := range events {
		e[i] = domain.EventUser{
			User: domain.User{
				ID:           event.ID_2,
				Username:     event.Username,
				Email:        event.Email,
				Password:     event.Password,
				VacationDays: event.VacationDays,
				IsSuperuser:  event.IsSuperuser,
				CreatedAt:    event.CreatedAt_2,
				EditedAt:     event.EditedAt_2,
				Color:        event.Color,
			},
			Event: domain.Event{
				Name:        event.Name,
				ID:          event.ID,
				State:       event.State,
				ScheduledAt: event.ScheduledAt,
				CreatedAt:   event.CreatedAt,
				EditedAt:    event.EditedAt,
				UserID:      event.UserID,
			},
		}
	}

	return e, nil
}

func (r *SQLEventRepo) GetPendingForUser(
	ctx context.Context,
	userId int64,
	year int,
) (int, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	params := repo.GetPendingEventsForYearParams{
		ScheduledAt:   start,
		ScheduledAt_2: start.AddDate(1, 0, 0),
		UserID:        userId,
	}
	count, err := r.r.GetPendingEventsForYear(ctx, params)
	if err != nil {
		r.log.Error(
			"GetPendingEventsForYear failed",
			slog.Int64("userId", userId),
			slog.Int("year", year),
			slog.String("error", err.Error()),
		)
		return 0, err
	}

	return int(count), nil
}

func (r *SQLEventRepo) GetUsedVacationForUser(
	ctx context.Context,
	userId int64,
	year int,
) (float64, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	params := repo.GetVacationCountForUserParams{
		ScheduledAt:   start,
		ScheduledAt_2: start.AddDate(1, 0, 0),
		UserID:        userId,
	}

	count, err := r.r.GetVacationCountForUser(ctx, params)
	if err != nil {
		r.log.Error(
			"GetVacationCountForUser failed",
			slog.Int64("userId", userId),
			slog.Int("year", year),
			slog.String("error", err.Error()),
		)
		return 0, err
	}

	if count == nil {
		return 0.0, nil
	}

	return *count, nil
}

func (r *SQLEventRepo) GetById(ctx context.Context, eventId int64) (*domain.Event, error) {
	event, err := r.r.GetEventById(ctx, eventId)
	if err != nil {
		r.log.Error(
			"GetEventById failed",
			slog.Int64("eventId", eventId),
			slog.String("error", err.Error()),
		)

		return nil, err

	}

	return (*domain.Event)(&event), nil
}

func (r *SQLEventRepo) UpdateInRange(
	ctx context.Context,
	userId int64,
	state string,
	start, end time.Time,
) error {
	err := r.r.UpdateEventsRange(
		ctx,
		repo.UpdateEventsRangeParams{
			UserID:        userId,
			State:         state,
			ScheduledAt:   start,
			ScheduledAt_2: end,
		},
	)
	if err != nil {
		r.log.Error(
			"UpdateEventsRange failed",
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}
