package domain

import (
	"context"
	"time"
)

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

// ID           int64     `json:"id"`
// Message      *string   `json:"message"`
// State        string    `json:"state"`
// CreatedAt    time.Time `json:"created_at"`
// EditedAt     time.Time `json:"edited_at"`
// UserID       int64     `json:"user_id"`
// EditedBy     *int64    `json:"edited_by"`
// EventID      int64     `json:"event_id"`
// ID_2         int64     `json:"id_2"`
// Username     string    `json:"username"`
// Email        string    `json:"email"`
// Password     string    `json:"password"`
// VacationDays int64     `json:"vacation_days"`
// IsSuperuser  bool      `json:"is_superuser"`
// CreatedAt_2  time.Time `json:"created_at_2"`
// EditedAt_2   time.Time `json:"edited_at_2"`
// Color        string    `json:"color"`
// Role         string    `json:"role"`
// Enabled      bool      `json:"enabled"`
// AworkID      *string   `json:"awork_id"`
// ID_3         int64     `json:"id_3"`
// ScheduledAt  time.Time `json:"scheduled_at"`
// Name         string    `json:"name"`
// State_2      string    `json:"state_2"`
// CreatedAt_3  time.Time `json:"created_at_3"`
// EditedAt_3   time.Time `json:"edited_at_3"`
// UserID_2     int64     `json:"user_id_2"`

type RequestEventUser struct {
	ID           int64     `json:"request_id"`
	Message      *string   `json:"message"`
	State        string    `json:"request_state"`
	CreatedAt    time.Time `json:"created_at"`
	EditedAt     time.Time `json:"edited_at"`
	UserID       int64     `json:"-"`
	EditedBy     *int64    `json:"edited_by"`
	EventID      int64     `json:"-"`
	ID_2         int64     `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	VacationDays int64     `json:"vacation_days"`
	IsSuperuser  bool      `json:"is_superuser"`
	CreatedAt_2  time.Time `json:"user_created_at"`
	EditedAt_2   time.Time `json:"user_edited_at"`
	Color        string    `json:"color"`
	Role         string    `json:"role"`
	Enabled      bool      `json:"enabled"`
	AworkID      *string   `json:"awork_id"`
	WorkdayHours float64   `json:"workday_hours"`
	WorkdaysWeek float64   `json:"workdays_week"`
	ID_3         int64     `json:"event_id"`
	ScheduledAt  time.Time `json:"scheduled_at"`
	Name         string    `json:"name"`
	State_2      string    `json:"event_state"`
	CreatedAt_3  time.Time `json:"event_created_at"`
	EditedAt_3   time.Time `json:"event_edited_at"`
	UserID_2     int64     `json:"-"`
}

type BatchRequest struct {
	StartDate  time.Time         `json:"start_date"`
	EndDate    time.Time         `json:"end_date"`
	EventCount int               `json:"event_count"`
	Request    *RequestEventUser `json:"request"`
	Conflicts  *[]User           `json:"conflicts"`
}

type RejectModalForm struct {
	UserID    int64 `query:"user_id"`
	StartDate int64 `query:"start_date"`
	EndDate   int64 `query:"end_date"`
	RequestID int64 `query:"request_id"`
}

type PatchRequestForm struct {
	UserID    int64  `form:"user_id"`
	State     string `form:"state"`
	Reason    string `form:"reason"`
	StartDate int64  `form:"start_date"`
	EndDate   int64  `form:"end_date"`
}

type ApiPatchRequestForm struct {
	UserID    int64     `form:"user_id"`
	State     string    `form:"state"`
	Reason    string    `form:"reason"`
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
}

type RequestRepository interface {
	Create(ctx context.Context, msg string, user *User, event *Event) (*Request, error)
	Update(ctx context.Context, editor *User, req *Request) (*Request, error)
	GetPending(ctx context.Context) ([]RequestEventUser, error)
	GetEventNameFrom(ctx context.Context, reqId int64) (string, error)
	GetInRange(ctx context.Context, userId int64, start, end time.Time) ([]Request, error)
	UpdateInRange(
		ctx context.Context,
		state string,
		editor, reqUserId int64,
		start, end time.Time,
	) (int64, error)
}
