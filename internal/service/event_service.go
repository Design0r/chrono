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

type EventService interface {
	Create(
		ctx context.Context,
		data domain.YMDDate,
		eventType string,
		user *domain.User,
	) (*domain.Event, error)
	Update(ctx context.Context, eventId int64, state string) (*domain.Event, error)
	Delete(ctx context.Context, id int64, currUser *domain.User) (*domain.Event, error)
	GetForDay(ctx context.Context, data domain.YMDDate) ([]domain.Event, error)
	GetForMonth(
		ctx context.Context,
		data domain.YMDate,
		userFilter *domain.User,
		eventFilter string,
	) (domain.Month, error)
	GetHistogramForYear(ctx context.Context, year int) ([]domain.YearHistogram, error)
	GetPendingForUser(ctx context.Context, userId int64, year int) (int, error)
	GetUsedVacationForUser(ctx context.Context, userId int64, year int) (float64, error)
	GetUserWithVacation(
		ctx context.Context,
		userId int64,
		year int,
		month int,
	) (domain.UserWithVacation, error)
	GetAllUsersWithVacation(ctx context.Context, year int) ([]domain.UserWithVacation, error)
	UpdateInRange(ctx context.Context, userId int64, state string, start, end time.Time) error
}

type eventService struct {
	log      *slog.Logger
	event    domain.EventRepository
	vacation VacationTokenService
	request  RequestService
	user     UserService
}

func NewEventService(
	e domain.EventRepository,
	r RequestService,
	u UserService,
	v VacationTokenService,
	log *slog.Logger,
) eventService {
	return eventService{log: log, event: e, request: r, user: u, vacation: v}
}

func (svc *eventService) Create(
	ctx context.Context,
	data domain.YMDDate,
	eventType string,
	user *domain.User,
) (*domain.Event, error) {
	evt := domain.Event{Name: eventType}

	if evt.IsVacation() && user.IsSuperuser {
		_, err := svc.vacation.Create(ctx, -1, data.Year, user.ID)
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

func (svc *eventService) Update(
	ctx context.Context,
	eventId int64,
	state string,
) (*domain.Event, error) {
	return svc.event.Update(ctx, eventId, state)
}

func (svc *eventService) Delete(
	ctx context.Context,
	eventId int64,
	currUser *domain.User,
) (*domain.Event, error) {
	event, err := svc.event.GetById(ctx, eventId)
	if err != nil {
		return nil, err
	}

	if !currUser.IsAdmin() && currUser.ID != eventId {
		return nil, fmt.Errorf("User: %v has no permission to delete the event.", currUser.ID)
	}

	_, err = svc.event.Delete(ctx, eventId)
	if err != nil {
		return nil, err
	}

	if event.IsVacation() && event.IsAccepted() {
		_, err := svc.vacation.Create(ctx, 1.0, event.ScheduledAt.Year(), event.UserID)
		if err != nil {
			return nil, err
		}
	}

	return event, nil
}

func (svc *eventService) GetForDay(
	ctx context.Context,
	data domain.YMDDate,
) ([]domain.Event, error) {
	return svc.event.GetForDay(ctx, data)
}

func (svc *eventService) GetForMonth(
	ctx context.Context,
	data domain.YMDate,
	userFilter *domain.User,
	eventFilter string,
) (domain.Month, error) {
	cfg := config.GetConfig()
	return svc.event.GetForMonth(ctx, data, cfg.BotName, userFilter, eventFilter)
}

func (svc *eventService) GetHistogramForYear(
	ctx context.Context,
	year int,
) ([]domain.YearHistogram, error) {
	events, err := svc.event.GetForYear(ctx, year)
	if err != nil {
		return nil, nil
	}
	numDays := domain.NumDaysInYear(year)
	eventList := make([]domain.YearHistogram, numDays)

	for _, event := range events {
		i := event.Event.ScheduledAt.YearDay() - 1
		date := event.Event.ScheduledAt
		days := domain.GetNumDaysOfMonth(date.Month(), date.Year())

		eventList[i].Count += 1
		eventList[i].IsHoliday = event.User.ID == 1
		eventList[i].LastDayOfMonth = date.Day() == days
		_, dateWeek := date.ISOWeek()
		_, currWeek := time.Now().ISOWeek()
		eventList[i].IsCurrentWeek = dateWeek == currWeek
		eventList[i].Usernames = append(eventList[i].Usernames, event.User.Username)
		s := strings.Split(date.Format(time.DateOnly), "-")
		slices.Reverse(s)
		eventList[i].Date = strings.Join(s, ".")
	}

	return eventList, err
}

func (svc *eventService) GetPendingForUser(
	ctx context.Context,
	userId int64,
	year int,
) (int, error) {
	return svc.event.GetPendingForUser(ctx, userId, year)
}

func (svc *eventService) GetUsedVacationForUser(
	ctx context.Context,
	userId int64,
	year int,
) (float64, error) {
	return svc.event.GetUsedVacationForUser(ctx, userId, year)
}

func (svc *eventService) GetUserWithVacation(
	ctx context.Context,
	userId int64,
	year int,
	month int,
) (domain.UserWithVacation, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 2, 0)
	fmt.Println(end)
	remaining, err := svc.vacation.GetRemainingVacationForUser(ctx, userId, start, start)
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

func (svc *eventService) GetAllUsersWithVacation(
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

func (svc *eventService) UpdateInRange(
	ctx context.Context,
	userId int64,
	state string,
	start, end time.Time,
) error {
	return svc.event.UpdateInRange(ctx, userId, state, start, end)
}
