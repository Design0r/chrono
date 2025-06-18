package domain

import (
	"context"
	"fmt"
	"slices"
	"time"
)

var vacationNames = []string{"urlaub", "urlaub halbtags"}

type Event struct {
	ID          int64     `json:"id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Name        string    `json:"name"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	EditedAt    time.Time `json:"edited_at"`
	UserID      int64     `json:"user_id"`
}

func (e *Event) IsVacation() bool {
	return slices.Contains(vacationNames, e.Name)
}

func (e *Event) IsAccepted() bool {
	return e.State == "accepted"
}

func (e *Event) RequestMsg(username string) string {
	return fmt.Sprintf("%v sent a new request for %v.", username, e.Name)
}

func (e *Event) AcceptMsg(username string) string {
	return fmt.Sprintf("%v accepted your %v request.", username, e.Name)
}

func (e *Event) RejectMsg(username string) string {
	return fmt.Sprintf("%v rejected your %v request.", username, e.Name)
}

func (e *Event) UpdateMsg(username string, state string) string {
	return fmt.Sprintf("%v %v your %v request.", username, state, e.Name)
}

type EventRepository interface {
	Create(ctx context.Context, data YMDDate, eventType string, user *User) (*Event, error)
	Update(ctx context.Context, eventId int64, state string) (*Event, error)
	Delete(ctx context.Context, id int64) (*Event, error)
	GetForDay(ctx context.Context, data YMDDate) ([]Event, error)
	GetForMonth(
		ctx context.Context,
		data YMDate,
		botName string,
		userFiler *User,
		eventFilter string,
	) (Month, error)
	GetForYear(ctx context.Context, year int) ([]EventUser, error)
	GetPendingForUser(ctx context.Context, userId int64, year int) (int, error)
	GetUsedVacationForUser(ctx context.Context, userId int64, year int) (float64, error)
	GetById(ctx context.Context, eventId int64) (*Event, error)
	// GetInRange(ctx context.Context, userId int64, start, end time.Time) ([]Event, error)
	UpdateInRange(ctx context.Context, userId int64, state string, start, end time.Time) error
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
