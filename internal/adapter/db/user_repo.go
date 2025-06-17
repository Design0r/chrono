package db

import (
	"context"
	"log/slog"
	"time"

	"chrono/db/repo"
	"chrono/internal/domain"
)

type SQLUserRepo struct {
	q   *repo.Queries
	log *slog.Logger
}

func NewSQLUserRepo(q *repo.Queries, l *slog.Logger) domain.UserRepository {
	return &SQLUserRepo{q: q, log: l}
}

func (r *SQLUserRepo) Create(ctx context.Context, user *domain.CreateUser) (*domain.User, error) {
	u, err := r.q.CreateUser(
		ctx,
		repo.CreateUserParams{
			Username:     user.Username,
			Color:        user.Color,
			VacationDays: user.VacationDays,
			Email:        user.Email,
			Password:     user.Password,
			IsSuperuser:  user.IsSuperuser,
		},
	)
	if err != nil {
		r.log.Error(
			"CreateUser failed:",
			slog.String("email", user.Email),
			slog.String("error", err.Error()),
		)
		return &domain.User{}, err
	}

	return (*domain.User)(&u), nil
}

func (r *SQLUserRepo) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	u, err := r.q.UpdateUser(
		ctx,
		repo.UpdateUserParams{
			Username:     user.Username,
			Color:        user.Color,
			VacationDays: user.VacationDays,
			Email:        user.Email,
			Password:     user.Password,
			IsSuperuser:  user.IsSuperuser,
			ID:           user.ID,
		},
	)
	if err != nil {
		r.log.Error(
			"UpdateUser failed:",
			slog.String("email", user.Email),
			slog.String("error", err.Error()),
		)
		return &domain.User{}, err
	}

	return (*domain.User)(&u), nil
}

func (r *SQLUserRepo) Delete(ctx context.Context, id int64) error {
	err := r.q.DeleteUser(ctx, id)
	if err != nil {
		r.log.Error(
			"DeleteUser failed:",
			slog.Int64("user_id", id),
			slog.String("error", err.Error()),
		)
		return err
	}

	return nil
}

func (r *SQLUserRepo) GetById(ctx context.Context, id int64) (*domain.User, error) {
	u, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		r.log.Error(
			"GetUserByID failed:",
			slog.Int64("user_id", id),
			slog.String("error", err.Error()),
		)
		return &domain.User{}, err
	}

	return (*domain.User)(&u), nil
}

func (r *SQLUserRepo) GetByName(ctx context.Context, name string) (*domain.User, error) {
	u, err := r.q.GetUserByName(ctx, name)
	if err != nil {
		r.log.Error(
			"GetUserByName failed:",
			slog.String("username", name),
			slog.String("error", err.Error()),
		)
		return &domain.User{}, err
	}

	return (*domain.User)(&u), nil
}

func (r *SQLUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		r.log.Error(
			"GetUserByEmail failed:",
			slog.String("email", email),
			slog.String("error", err.Error()),
		)
		return &domain.User{}, err
	}

	return (*domain.User)(&u), nil
}

func (r *SQLUserRepo) GetAll(ctx context.Context) ([]domain.User, error) {
	u, err := r.q.GetAllUsers(ctx)
	if err != nil {
		r.log.Error(
			"GetAllUsers failed:",
			slog.String("error", err.Error()),
		)
		return []domain.User{}, err
	}

	users := make([]domain.User, len(u))
	for i := range u {
		users[i] = (domain.User)(u[i])
	}

	return users, nil
}

func (r *SQLUserRepo) GetAdmins(ctx context.Context) ([]domain.User, error) {
	u, err := r.q.GetAdmins(ctx)
	if err != nil {
		r.log.Error(
			"GetAdmins failed:",
			slog.String("error", err.Error()),
		)
		return []domain.User{}, err
	}

	admins := make([]domain.User, len(u))
	for i := range u {
		admins[i] = (domain.User)(u[i])
	}

	return admins, nil
}

func (r *SQLUserRepo) GetConflicting(
	ctx context.Context,
	userId int64,
	start time.Time,
	end time.Time,
) ([]domain.User, error) {
	users, err := r.q.GetConflictingEventUsers(
		ctx,
		repo.GetConflictingEventUsersParams{ID: userId, ScheduledAt: start, ScheduledAt_2: end},
	)
	if err != nil {
		r.log.Error(
			"GetConflictingEventUsers failed",
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	u := make([]domain.User, len(users))
	for i := range u {
		u[i] = (domain.User)(u[i])
	}

	return u, nil
}
