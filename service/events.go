package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"calendar/db/repo"
	"calendar/schemas"
)

func CreateEvent(db *sql.DB, data schemas.YMDDate) (repo.Event, error) {
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

	event, err := r.CreateEvent(context.Background(), date)
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

func GetEventsForMonth(db *sql.DB, date time.Time) (map[int][]repo.Event, error) {
	r := repo.New(db)

	events, err := r.CustomGetEventsForMonth(context.Background(), date)
	if err != nil {
		log.Printf("Failed getting events: %v", err)
		return map[int][]repo.Event{}, err
	}

	eventMap := map[int][]repo.Event{}
	for _, event := range events {
		eventMap[event.ScheduledAt.Day()] = append(eventMap[event.ScheduledAt.Day()], event)
	}

	return eventMap, nil
}
