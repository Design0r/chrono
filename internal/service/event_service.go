package service

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"chrono/config"
	"chrono/internal/domain"
)

type EventService struct {
	log     *slog.Logger
	event   domain.EventRepository
	token   *TokenService
	request *RequestService
	user    *UserService
}

func NewEventService(
	e domain.EventRepository,
	r *RequestService,
	u *UserService,
	t *TokenService,
	log *slog.Logger,
) EventService {
	return EventService{log: log, event: e, request: r, user: u, token: t}
}

func (svc *EventService) Create(
	ctx context.Context,
	data domain.YMDDate,
	eventType string,
	user *domain.User,
) (*domain.Event, error) {
	evt := domain.Event{Name: eventType}

	if evt.IsVacation() && user.IsSuperuser {
		_, err := svc.token.CreateVacationToken(ctx, -1, data.Year, user.ID)
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
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (svc *EventService) Update(
	ctx context.Context,
	eventId int64,
	state string,
) (*domain.Event, error) {
	return svc.event.Update(ctx, eventId, state)
}

func (svc *EventService) Delete(
	ctx context.Context,
	eventId int64,
	currUser *domain.User,
) (*domain.Event, error) {
	event, err := svc.event.GetById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	if !currUser.IsAdmin() && currUser.ID != event.UserID {
		return nil, fmt.Errorf("User: %v has no permission to delete the event.", currUser.Username)
	}

	err = svc.event.Delete(ctx, eventId)
	if err != nil {
		return nil, err
	}

	if event.IsVacation() && event.IsAccepted() {
		_, err := svc.token.CreateVacationToken(ctx, 1.0, event.ScheduledAt.Year(), event.UserID)
		if err != nil {
			return nil, err
		}
	}

	return event, nil
}

func (svc *EventService) GetForDay(
	ctx context.Context,
	data domain.YMDDate,
) ([]domain.Event, error) {
	return svc.event.GetForDay(ctx, data)
}

func (svc *EventService) GetForMonth(
	ctx context.Context,
	data domain.YMDate,
	userFilter *domain.User,
	eventFilter string,
) (domain.Month, error) {
	cfg := config.GetConfig()
	return svc.event.GetForMonth(ctx, data, cfg.BotName, userFilter, eventFilter)
}

func (svc *EventService) GetForYear(
	ctx context.Context,
	year int,
) ([]domain.EventUser, error) {
	return svc.event.GetForYear(ctx, year)
}

func (svc *EventService) GetHistogramForYear(
	ctx context.Context,
	year int,
) ([]domain.YearHistogram, error) {
	events, err := svc.event.GetForYear(ctx, year)
	if err != nil {
		return nil, nil
	}
	numDays := domain.NumDaysInYear(year)
	eventList := make([]domain.YearHistogram, numDays)

	for i := range eventList {
		date := time.Date(year, time.Month(1), i+1, 0, 0, 0, 0, time.Local)
		days := domain.GetNumDaysOfMonth(date.Month(), date.Year())
		eventList[i].LastDayOfMonth = date.Day() == days

		s := strings.Split(date.Format(time.DateOnly), "-")
		slices.Reverse(s)
		eventList[i].Date = strings.Join(s, ".")
	}

	for _, event := range events {
		i := event.Event.ScheduledAt.YearDay() - 1
		date := event.Event.ScheduledAt

		eventList[i].Count += 1
		eventList[i].IsHoliday = event.User.ID == 1
		_, dateWeek := date.ISOWeek()
		_, currWeek := time.Now().ISOWeek()
		eventList[i].IsCurrentWeek = dateWeek == currWeek
		eventList[i].Usernames = append(eventList[i].Usernames, event.User.Username)
	}

	return eventList, err
}

func (svc *EventService) GetPendingForUser(
	ctx context.Context,
	userId int64,
	year int,
) (int, error) {
	return svc.event.GetPendingForUser(ctx, userId, year)
}

func (svc *EventService) GetUsedVacationForUser(
	ctx context.Context,
	userId int64,
	year int,
) (float64, error) {
	return svc.event.GetUsedVacationForUser(ctx, userId, year)
}

func (svc *EventService) GetUserWithVacation(
	ctx context.Context,
	userId int64,
	year int,
	month int,
) (domain.UserWithVacation, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	remaining, err := svc.token.GetRemainingVacationForUser(ctx, userId, start, start)
	if err != nil {
		return domain.UserWithVacation{}, err
	}

	used, err := svc.GetUsedVacationForUser(ctx, userId, year)
	if err != nil {
		return domain.UserWithVacation{}, err
	}

	user, err := svc.user.GetById(ctx, userId)
	if err != nil {
		return domain.UserWithVacation{}, err
	}

	pending, err := svc.GetPendingForUser(ctx, user.ID, year)
	if err != nil {
		return domain.UserWithVacation{}, err
	}

	return domain.UserWithVacation{
		VacationRemaining: remaining,
		VacationUsed:      used,
		User:              *user,
		PendingEvents:     pending,
	}, nil
}

func (svc *EventService) GetAllUsersWithVacation(
	ctx context.Context,
	year int,
) ([]domain.UserWithVacation, error) {
	allUsers, err := svc.user.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	allUsersWithVac := make([]domain.UserWithVacation, len(allUsers))
	for i, user := range allUsers {
		u, err := svc.GetUserWithVacation(ctx, user.ID, year, 1)
		if err != nil {
			svc.log.Error(
				"Unable to get user with vacation, skipping user.",
				slog.String("username", user.Username),
				slog.String("error", err.Error()),
			)
			continue
		}
		allUsersWithVac[i] = u
	}

	return allUsersWithVac, nil
}

func (svc *EventService) UpdateInRange(
	ctx context.Context,
	userId int64,
	state string,
	start, end time.Time,
) error {
	return svc.event.UpdateInRange(ctx, userId, state, start, end)
}

func (svc *EventService) GetAllByUserId(
	ctx context.Context,
	userId int64,
) ([]domain.Event, error) {
	return svc.event.GetAllByUserId(ctx, userId)
}

func (svc *EventService) GetNonWeekendCountHolidays(
	ctx context.Context,
	start, end time.Time,
) (int, error) {
	cfg := config.GetConfig()

	// "Bot" user that holds holiday events
	bot, err := svc.user.GetByName(ctx, cfg.BotName)
	if err != nil {
		return 0, err
	}

	holidays, err := svc.event.GetAllByUserId(ctx, bot.ID)
	if err != nil {
		return 0, err
	}

	nonWeekendHolidays := 0

	for _, h := range holidays {
		t := h.ScheduledAt

		// only count events within [start, end]
		if t.Before(start) || t.After(end) {
			continue
		}

		wd := t.Weekday()
		if wd == time.Saturday || wd == time.Sunday {
			continue
		}

		nonWeekendHolidays++
	}

	return nonWeekendHolidays, nil
}

func (svc *EventService) GetUsedVacation(
	ctx context.Context,
	userId int64,
	start, end time.Time,
) (float64, error) {
	events, err := svc.event.GetAllByUserId(ctx, userId)
	if err != nil {
		return 0, err
	}

	count := 0.0

	for _, e := range events {
		t := e.ScheduledAt

		// only count events within [start, end]
		if t.Before(start) || t.After(end) {
			continue
		}

		// only accepted vacation
		if e.State != "accepted" {
			continue
		}

		switch e.Name {
		case "urlaub":
			count += 1.0
		case "urlaub halbtags":
			count += 0.5
		}
	}

	return count, nil
}
