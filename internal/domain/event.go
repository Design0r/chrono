package domain

import "time"

type Event struct {
	ID          int64     `json:"id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Name        string    `json:"name"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	EditedAt    time.Time `json:"edited_at"`
	UserID      int64     `json:"user_id"`
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
