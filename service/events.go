package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"calendar/db/repo"
	"calendar/schemas"
)

func CreateEvent(db *sql.DB, data schemas.YMDDate, id int64) (repo.Event, error) {
	r := repo.New(db)

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

	event, err := r.CreateEvent(
		context.Background(),
		repo.CreateEventParams{UserID: id, ScheduledAt: date},
	)
	if err != nil {
		log.Printf("Failed creating event: %v", err)
		return repo.Event{}, err
	}

	return event, nil
}

func GetEventsForDay(db *sql.DB, data schemas.YMDDate) ([]repo.Event, error) {
	r := repo.New(db)

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
	db *sql.DB,
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
	r := repo.New(db)

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
		newEvent := schemas.Event{
			Username: event.Username,
			Event: repo.Event{
				ID:          event.ID,
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
