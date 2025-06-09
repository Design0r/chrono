package service

import (
	"context"
	"fmt"

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
	ToggleAdmin(
		ctx context.Context,
		userToUpdate int64,
		currUser *domain.User,
	) (*domain.User, error)
}

type userService struct {
	user  domain.UserRepository
	notif NotificationService
}

func NewUserService(r domain.UserRepository) userService {
	return userService{user: r}
}

func (svc *userService) Create(ctx context.Context, user *domain.CreateUser) (*domain.User, error) {
	return svc.user.Create(ctx, user)
}

func (svc *userService) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	return svc.user.Update(ctx, user)
}

func (svc *userService) Delete(ctx context.Context, id int64) error {
	return svc.user.Delete(ctx, id)
}

func (svc *userService) GetById(ctx context.Context, id int64) (*domain.User, error) {
	return svc.user.GetById(ctx, id)
}

func (svc *userService) GetByName(ctx context.Context, name string) (*domain.User, error) {
	return svc.user.GetByName(ctx, name)
}

func (svc *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return svc.user.GetByName(ctx, email)
}

func (svc *userService) GetAll(ctx context.Context) ([]*domain.User, error) {
	return svc.user.GetAll(ctx)
}

func (svc *userService) GetUsersWithVacation(
	ctx context.Context,
) ([]*domain.UserWithVacation, error) {
	return nil, nil
}

func (svc *userService) ToggleAdmin(
	ctx context.Context,
	userToChange int64,
	currUser *domain.User,
) (*domain.User, error) {
	if !currUser.IsSuperuser {
		return nil, nil
	}
	user, err := svc.GetById(ctx, userToChange)
	if err != nil {
		return nil, err
	}
	user.IsSuperuser = true

	updatedUser, err := svc.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("%v changed your admin status to %v", currUser.Username, user.IsSuperuser)
	err = svc.notif.CreateAndNotify(ctx, msg, []domain.User{*user})
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
