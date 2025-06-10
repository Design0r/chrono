package domain

import (
	"context"
	"slices"
	"time"
)

type Event struct {
	ID          int64     `json:"id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Name        string    `json:"name"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	EditedAt    time.Time `json:"edited_at"`
	UserID      int64     `json:"user_id"`
}

var vacationNames []string = []string{"urlaub", "urlaub halbtags"}

func (e *Event) IsVacation() bool {
	return slices.Contains(vacationNames, e.Name)
}

type EventRepository interface {
	Create(ctx context.Context, data YMDate, eventType string, user *User) (*Event, error)
	Delete(ctx context.Context, id int64) (*Event, error)
	GetForDay(ctx context.Context, data YMDDate) ([]Event, error)
	GetForMonth(ctx context.Context, data YMDDate) ([]Event, error)
}

type EventUser struct {
	Event Event
	User  User
}

type YMDDate struct {
	Year  int `param:"year"`
	Month int `param:"month"`
	Day   int `param:"day"`
}

type YMDate struct {
	Year  int `param:"year"`
	Month int `param:"month"`
}

type YDate struct {
	Year int `param:"year"`
}
