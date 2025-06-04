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
	GetUsersWithVacation(ctx context.Context) ([]*domain.UserWithVacation, error)
}

type userService struct {
	r domain.UserRepository
}

func NewUserService(r domain.UserRepository) userService {
	return userService{r: r}
}

func (svc *userService) Create(ctx context.Context, user *domain.CreateUser) (*domain.User, error) {
	return svc.r.Create(ctx, user)
}

func (svc *userService) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	return svc.r.Update(ctx, user)
}

func (svc *userService) Delete(ctx context.Context, id int64) error {
	return svc.r.Delete(ctx, id)
}

func (svc *userService) GetById(ctx context.Context, id int64) (*domain.User, error) {
	return svc.r.GetById(ctx, id)
}

func (svc *userService) GetByName(ctx context.Context, name string) (*domain.User, error) {
	return svc.r.GetByName(ctx, name)
}

func (svc *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return svc.r.GetByName(ctx, email)
}

func (svc *userService) GetAll(ctx context.Context) ([]*domain.User, error) {
	return svc.r.GetAll(ctx)
}

func (svc *userService) GetUsersWithVacation(ctx context.Context) ([]*domain.UserWithVacation, error) {
	return nil, nil
}
