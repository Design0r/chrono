package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"chrono/internal/domain"
)

type UserService struct {
	user  domain.UserRepository
	notif *NotificationService
	token *TokenService
	log   *slog.Logger
}

func NewUserService(
	r domain.UserRepository,
	n *NotificationService,
	t *TokenService,
	log *slog.Logger,
) *UserService {
	return &UserService{user: r, notif: n, token: t, log: log}
}

func (svc *UserService) Create(ctx context.Context, user *domain.CreateUser) (*domain.User, error) {
	return svc.user.Create(ctx, user)
}

func (svc *UserService) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	return svc.user.Update(ctx, user)
}

func (svc *UserService) Delete(ctx context.Context, id int64) error {
	return svc.user.Delete(ctx, id)
}

func (svc *UserService) GetById(ctx context.Context, id int64) (*domain.User, error) {
	return svc.user.GetById(ctx, id)
}

func (svc *UserService) GetByName(ctx context.Context, name string) (*domain.User, error) {
	return svc.user.GetByName(ctx, name)
}

func (svc *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return svc.user.GetByEmail(ctx, email)
}

func (svc *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	return svc.user.GetAll(ctx)
}

func (svc *UserService) GetUsersWithVacation(
	ctx context.Context,
) ([]*domain.UserWithVacation, error) {
	return nil, nil
}

func (svc *UserService) SetUserRole(
	ctx context.Context,
	userToChange int64,
	role domain.Role,
	currUser *domain.User,
) (*domain.User, error) {
	if !currUser.IsSuperuser {
		return nil, nil
	}
	user, err := svc.GetById(ctx, userToChange)
	if err != nil {
		return nil, err
	}

	if !domain.IsValidRole(role) {
		return nil, fmt.Errorf("Invalid user role")
	}
	user.Role = string(role)

	updatedUser, err := svc.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	msg := fmt.Sprintf("%v changed your user role to %v", currUser.Username, user.Role)
	err = svc.notif.CreateAndNotify(ctx, msg, []domain.User{*user})
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (svc *UserService) SetVacation(ctx context.Context, userId int64, vacation, year int) error {
	user, err := svc.GetById(ctx, userId)
	if err != nil {
		return err
	}
	oldVacation := int(user.VacationDays)

	if vacation < 0 {
		svc.log.Error("negative vacation value is not supported", "value", vacation)
		return fmt.Errorf("negative vacation value is not supported %v", vacation)
	}
	user.VacationDays = int64(vacation)
	_, err = svc.Update(ctx, user)
	if err != nil {
		return err
	}

	return svc.token.UpdateYearlyTokens(ctx, userId, vacation-oldVacation, year)
}

func (svc *UserService) GetConflicting(
	ctx context.Context,
	userId int64,
	start time.Time,
	end time.Time,
) ([]domain.User, error) {
	return svc.user.GetConflicting(ctx, userId, start, end)
}
