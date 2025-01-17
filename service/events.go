package service

import (
	"context"
	"log"
	"time"

	"chrono/db/repo"
	"chrono/schemas"
)

func CreateEvent(
	r *repo.Queries,
	data schemas.YMDDate,
	user repo.User,
	name string,
) (repo.Event, error) {
	if name != "urlaub" || user.IsSuperuser {
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
		time.Now().Local().Location(),
	)

	state := "pending"
	if name != "urlaub" || user.IsSuperuser {
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
		time.Now().Local().Location(),
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
) error {
	date := time.Date(
		month.Year,
		time.Month(month.Number),
		1,
		0,
		0,
		0,
		0,
		time.Now().Local().Location(),
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
		newEvent := schemas.Event{
			Username: user.Username,
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
		month.Days[idx].Events = append(month.Days[idx].Events, newEvent)
	}

	return nil
}

func GetVacationCountForUserYear(r *repo.Queries, userId int, year int) (int, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())
	yearEnd := yearStart.AddDate(1, 0, 0)

	count, err := r.GetVacationCountForUser(
		context.Background(),
		repo.GetVacationCountForUserParams{
			UserID:        int64(userId),
			ScheduledAt:   yearStart,
			ScheduledAt_2: yearEnd,
		},
	)
	if err != nil {
		return 0, err
	}

	return int(count), nil
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
