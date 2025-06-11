package service

import (
	"chrono/internal/domain"
	"context"
	"log/slog"
	"time"
)

type EventService interface {
	Create(ctx context.Context, data domain.YMDDate, eventType string, user *domain.User) (*domain.Event, error)
	Update(ctx context.Context, eventId int64, state string) (*domain.Event, error)
	Delete(ctx context.Context, id int64) (*domain.Event, error)
	GetForDay(ctx context.Context, data domain.YMDDate) ([]domain.Event, error)
	GetForMonth(ctx context.Context, data domain.YMDate) (domain.Month, error)
	GetForYear(ctx context.Context, year int) ([]domain.EventUser, error)
	GetPendingForUser(ctx context.Context, userId int64, year int) (int, error)
	GetUsedVacationForUser(ctx context.Context, userId int64, year int) (float64, error)
}

type eventService struct {
	log      *slog.Logger
	event    domain.EventRepository
	vacation VacationTokenService
	request  RequestService
}

func NewEventService(e domain.EventRepository, log *slog.Logger, v VacationTokenService) eventService {
	return eventService{log: log, event: e}
}

func (svc *eventService) Create(ctx context.Context, data domain.YMDDate, eventType string, user *domain.User) (*domain.Event, error) {
	evt := domain.Event{Name: eventType}
	start := time.Date(data.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 3, 0)

	if evt.IsVacation() && user.IsSuperuser {
		_, err := svc.vacation.Create(ctx, domain.CreateVacationToken{StartDate: start, EndDate: end, UserID: user.ID, Value: -1})
		if err != nil {
			return nil, err
		}
		return svc.event.Create(ctx, data, eventType, user)
	}

	if !evt.IsVacation() {
		return svc.event.Create(ctx, data, eventType, user)
	}

	event, err := svc.event.Create(ctx, data, eventType, user)
	if err != nil {
		return nil, err
	}

	_, err = svc.request.Create(ctx, event.RequestMsg(user.Username), user, event)
	return event, nil
}

func (svc *eventService) Update(ctx context.Context, eventId int64, state string) (*domain.Event, error) {
	return svc.event.Update(ctx, eventId, state)
}

func (svc *eventService) Delete(ctx context.Context, eventId int64) (*domain.Event, error) {
	return svc.event.Delete(ctx, eventId)
}

func (svc *eventService) GetForDay(ctx context.Context, data domain.YMDDate) ([]domain.Event, error) {
	return svc.event.GetForDay(ctx, data)
}

func (svc *eventService) GetForMonth(ctx context.Context, data domain.YMDate) (domain.Month, error) {
	return svc.event.GetForMonth(ctx, data)
}

func (svc *eventService) GetForYear(ctx context.Context, year int) ([]domain.EventUser, error) {
	return svc.event.GetForYear(ctx, year)
}

func (svc *eventService) GetPendingForUser(ctx context.Context, userId int64, year int) (int, error) {
	return svc.event.GetPendingForUser(ctx, userId, year)
}

func (svc *eventService) GetUsedVacationForUser(ctx context.Context, userId int64, year int) (float64, error) {
	return svc.event.GetUsedVacationForUser(ctx, userId, year)
}
