package service

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"chrono/calendar"
	"chrono/config"
	"chrono/db/repo"
	"chrono/schemas"
)

var vacationNames []string = []string{"urlaub", "urlaub halbtags"}

func IsVacation(name string) bool {
	return slices.Contains(vacationNames, name)
}

func CreateEvent(
	r *repo.Queries,
	data schemas.YMDDate,
	user repo.User,
	name string,
) (repo.Event, error) {
	if IsVacation(name) && user.IsSuperuser {
		CreateToken(r, user.ID, data.Year, -1.0)
		return createEvent(r, data, user, name)
	}

	if !IsVacation(name) {
		return createEvent(r, data, user, name)
	}

	return createRequestEvent(r, data, user, name)
}

func createEvent(
	r *repo.Queries,
	data schemas.YMDDate,
	user repo.User,
	name string,
) (repo.Event, error) {
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

	state := "pending"
	if !IsVacation(name) || user.IsSuperuser {
		state = "accepted"
	}

	event, err := r.CreateEvent(
		context.Background(),
		repo.CreateEventParams{Name: name, UserID: user.ID, ScheduledAt: date, State: state},
	)
	if err != nil {
		log.Printf("Failed creating event: %v", err)
		return repo.Event{}, err
	}

	return event, nil
}

func createRequestEvent(
	r *repo.Queries,
	data schemas.YMDDate,
	user repo.User,
	name string,
) (repo.Event, error) {
	event, err := createEvent(r, data, user, name)
	if err != nil {
		return repo.Event{}, err
	}

	_, err = CreateRequest(r, GenerateRequestMsg(user.Username, event), user, event)
	if err != nil {
		return repo.Event{}, err
	}

	return event, nil
}

func DeleteEvent(r *repo.Queries, eventId int) (repo.Event, error) {
	event, err := r.DeleteEvent(context.Background(), int64(eventId))
	if err != nil {
		log.Printf("Failed deleting event: %v", err)
		return repo.Event{}, err
	}

	return event, nil
}

func GetEventsForDay(r *repo.Queries, data schemas.YMDDate) ([]repo.Event, error) {
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

	events, err := r.GetEventsForDay(context.Background(), date)
	if err != nil {
		log.Printf("Failed getting event: %v", err)
		return []repo.Event{}, err
	}

	return events, nil
}

func GetEventsForMonth(
	r *repo.Queries,
	month *schemas.Month,
	filter *repo.User,
	eventFilter string,
) error {
	cfg := config.GetConfig()
	date := time.Date(
		month.Year,
		time.Month(month.Number),
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	events, err := r.GetEventsForMonth(
		context.Background(),
		repo.GetEventsForMonthParams{ScheduledAt: date, ScheduledAt_2: date.AddDate(0, 1, 0)},
	)
	if err != nil {
		log.Printf("Failed getting events: %v", err)
		return err
	}

	for _, event := range events {
		idx := event.ScheduledAt.Day() - 1
		user, err := GetUserById(r, event.UserID)
		if err != nil {
			continue
		}
		if filter != nil && user.Username != filter.Username &&
			user.Username != cfg.BotName {
			continue
		}
		if eventFilter != "" && !strings.Contains(event.Name, eventFilter) &&
			eventFilter != "all" &&
			user.Username != cfg.BotName {
			continue
		}

		newEvent := schemas.Event{
			Username: user.Username,
			Color:    user.Color,
			Event: repo.Event{
				Name:        event.Name,
				ID:          event.ID,
				State:       event.State,
				ScheduledAt: event.ScheduledAt,
				CreatedAt:   event.CreatedAt,
				EditedAt:    event.EditedAt,
				UserID:      event.UserID,
			},
		}
		fmt.Println(newEvent.Event.Name, newEvent.Event.ScheduledAt)
		month.Days[idx].Events = append(month.Days[idx].Events, newEvent)
	}

	return nil
}

func GetEventsForYear(r *repo.Queries, year int) ([]repo.GetEventsForYearRow, error) {
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Now().Location())

	params := repo.GetEventsForYearParams{
		ScheduledAt:   yearStart,
		ScheduledAt_2: yearStart.AddDate(1, 0, 0),
	}

	events, err := r.GetEventsForYear(context.Background(), params)
	if err != nil {
		log.Printf("Failed to get events for year: %v", err)
		return nil, err
	}

	return events, err
}

func GetEventCountForYear(r *repo.Queries, year int) ([]schemas.YearHistogram, error) {
	events, err := GetEventsForYear(r, year)
	if err != nil {
		return nil, err
	}

	numDays := calendar.NumDaysInYear(year)
	eventList := make([]schemas.YearHistogram, numDays)

	for _, event := range events {
		i := event.ScheduledAt.YearDay() - 1
		date := event.ScheduledAt
		days := calendar.GetNumDaysOfMonth(date.Month(), date.Year())

		eventList[i].Count += 1
		eventList[i].IsHoliday = event.UserID == 1
		eventList[i].LastDayOfMonth = date.Day() == days
		_, dateWeek := date.ISOWeek()
		_, currWeek := time.Now().ISOWeek()
		eventList[i].IsCurrentWeek = dateWeek == currWeek
		eventList[i].Usernames = append(eventList[i].Usernames, event.Username)
		s := strings.Split(date.Format(time.DateOnly), "-")
		slices.Reverse(s)
		eventList[i].Date = strings.Join(s, ".")
	}

	return eventList, err
}

func GetRemainingVacation(r *repo.Queries, userId int64, year int, month int) (float64, error) {
	yearStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())

	value, err := GetValidUserTokenSum(r, userId, yearStart)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func UpdateEventState(r *repo.Queries, state string, eventId int64) (repo.Event, error) {
	params := repo.UpdateEventStateParams{State: state, ID: eventId}

	event, err := r.UpdateEventState(context.Background(), params)
	if err != nil {
		log.Printf("Failed to update event state: %v", err)
		return repo.Event{}, err
	}

	return event, nil
}

func GetPendingEventsForYear(r *repo.Queries, userId int64, year int) (int, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())

	params := repo.GetPendingEventsForYearParams{
		ScheduledAt:   start,
		ScheduledAt_2: start.AddDate(1, 0, 0),
		UserID:        userId,
	}
	count, err := r.GetPendingEventsForYear(context.Background(), params)
	if err != nil {
		log.Printf("Failed getting pending events: %v", count)
		return 0, err
	}

	return int(count), nil
}

func GetVacationCountForUser(r *repo.Queries, userId int64, year int) (float64, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())

	params := repo.GetVacationCountForUserParams{
		ScheduledAt:   start,
		ScheduledAt_2: start.AddDate(1, 0, 0),
		UserID:        userId,
	}

	count, err := r.GetVacationCountForUser(context.Background(), params)
	if err != nil {
		log.Printf("Failed getting vacation count: %v", err)
		return 0, err
	}

	if count != nil {
		return *count, nil
	}

	return 0, nil
}
