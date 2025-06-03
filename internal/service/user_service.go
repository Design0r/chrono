package service

import (
	"context"

	"chrono/internal/domain"
)

type UserService interface {
	Create(ctx context.Context, user *domain.CreateUser) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	GetById(ctx context.Context, id int64) (*domain.User, error)
	GetByName(ctx context.Context, name string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	Delete(ctx context.Context, id int64) error
	GetAllVacUsers
}
