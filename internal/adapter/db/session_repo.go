package db

import (
	"context"
	"log/slog"
	"time"

	"chrono/db/repo"
	"chrono/internal/domain"

	"github.com/labstack/gommon/log"
)

type SQLSessionRepo struct {
	q   *repo.Queries
	log *slog.Logger
}

func NewSQLSessionRepo(q *repo.Queries, log *slog.Logger) domain.SessionRepository {
	return &SQLSessionRepo{q: q, log: log}
}

func (r *SQLSessionRepo) Create(ctx context.Context, userId int64, secureRand string, duration time.Duration) (*domain.Session, error) {
	data := repo.CreateSessionParams{ID: secureRand, ValidUntil: time.Now().Add(duration), UserID: userId}
	session, err := r.q.CreateSession(ctx, data)
	if err != nil {
		r.log.Error("CreateSession failed", slog.String("error", err.Error()))
		return &domain.Session{}, err
	}

	return (*domain.Session)(&session), nil
}

func (r *SQLSessionRepo) Delete(ctx context.Context, cookie string) error {
	err := r.q.DeleteSession(ctx, cookie)
	if err != nil {
		r.log.Error("DeleteSession failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *SQLSessionRepo) DeleteAll(ctx context.Context) error {
	err := r.q.DeleteAllSessions(ctx)
	if err != nil {
		r.log.Error("DeleteAllSessions failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *SQLSessionRepo) GetById(ctx context.Context, cookie string) (*domain.Session, error) {
	session, err := r.q.GetSessionById(ctx, cookie)
	if err != nil {
		r.log.Error("GetSessionById failed", slog.String("error", err.Error()))
		return nil, err
	}

	return (*domain.Session)(&session), nil
}

func (r *SQLSessionRepo) GetSessionUser(ctx context.Context, cookie string) (*domain.User, error) {
	u, err := r.q.GetUserFromSession(ctx, cookie)
	if err != nil {
		log.Error("GetValidSession failed", slog.String("error", err.Error()))
		return nil, err
	}

	return (*domain.User)(&u), nil
}
