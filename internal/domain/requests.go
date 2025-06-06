package domain

import "time"

type Request struct {
	ID        int64     `json:"id"`
	Message   *string   `json:"message"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	EditedAt  time.Time `json:"edited_at"`
	UserID    int64     `json:"user_id"`
	EditedBy  *int64    `json:"edited_by"`
	EventID   int64     `json:"event_id"`
}

type RequestRepository interface {
	Create(msg string, user *User, event *Event) (*Request, error)
	Update(req *Request) error
	GetPending() ([]Request, error)
	GetEventNameFrom(req *Request) (string, error)
}
