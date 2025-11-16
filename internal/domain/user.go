package domain

import (
	"context"
	"slices"
	"time"
)

type Role string

var (
	AdminRole Role = "admin"
	UserRole  Role = "user"
	GuestRole Role = "guest"
)

func UserRoles() []Role {
	return []Role{AdminRole, UserRole, GuestRole}
}

func IsValidRole(role Role) bool {
	return slices.Contains(UserRoles(), role)
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	VacationDays int64     `json:"vacation_days"`
	IsSuperuser  bool      `json:"is_superuser"`
	CreatedAt    time.Time `json:"created_at"`
	EditedAt     time.Time `json:"edited_at"`
	Color        string    `json:"color"`
	Role         string    `json:"role"`
	Enabled      bool      `json:"enabled"`
}

func (u *User) IsAdmin() bool {
	return u.IsSuperuser
}

type UserWithVacation struct {
	User
	VacationRemaining float64 `json:"vacation_remaining"`
	VacationUsed      float64 `json:"vacation_used"`
	PendingEvents     int     `json:"pending_events"`
}

type PatchUser struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Color    string `form:"color"`
	Password string `form:"password"`
	Role     string `form:"role"`
	Enabled  bool   `form:"enabled"`
}

type ApiPatchUser struct {
	Name         string `form:"username"`
	Email        string `form:"email"`
	Color        string `form:"color"`
	Password     string `form:"password"`
	Role         string `form:"role"`
	Enabled      *bool  `form:"enabled"`
	VacationDays *int64 `form:"vacation_days"`
}

type CreateUser struct {
	Username     string `json:"username"      form:"qwenameasd"`
	Color        string `json:"color"`
	VacationDays int64  `json:"vacation_days"`
	Email        string `json:"email"         form:"qweemailasd"`
	Password     string `json:"password"      form:"qwepasswordasd"`
	IsSuperuser  bool   `json:"is_superuser"`
}

type ApiCreateUser struct {
	Username     string `json:"username"      form:"username"`
	Color        string `json:"color"`
	VacationDays int64  `json:"vacation_days"`
	Email        string `json:"email"         form:"email"`
	Password     string `json:"password"      form:"password"`
	IsSuperuser  bool   `json:"is_superuser"`
}

type Login struct {
	Email    string `form:"qweemailasd"`
	Password string `form:"qwepasswordasd"`
}

type APILogin struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type Honeypot struct {
	Name     string `form:"name"     json:"name"`
	Email    string `form:"email"    json:"email"`
	Password string `form:"password" json:"password"`
}

func (h *Honeypot) IsFilled() bool {
	return h.Name != "" || h.Email != "" || h.Password != ""
}

type UserRepository interface {
	Create(ctx context.Context, user *CreateUser) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	GetById(ctx context.Context, id int64) (*User, error)
	GetByName(ctx context.Context, name string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	GetAdmins(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, id int64) error
	GetConflicting(
		ctx context.Context,
		userId int64,
		start time.Time,
		end time.Time,
	) ([]User, error)
}
