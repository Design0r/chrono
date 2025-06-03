package domain

import (
	"context"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	VacationDays int64     `json:"vacation_days"`
	IsSuperuser  bool      `json:"is_superuser"`
	CreatedAt    time.Time `json:"created_at"`
	EditedAt     time.Time `json:"edited_at"`
	Color        string    `json:"color"`
}

type CreateUser struct {
	Username     string `json:"username"`
	Color        string `json:"color"`
	VacationDays int64  `json:"vacation_days"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	IsSuperuser  bool   `json:"is_superuser"`
}

type UserRepository interface {
	Create(ctx context.Context, user *CreateUser) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	GetById(ctx context.Context, id int64) (*User, error)
	GetByName(ctx context.Context, name string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id int64) error
}
