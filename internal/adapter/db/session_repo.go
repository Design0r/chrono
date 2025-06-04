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

func NewSQLSessionRepo(q *repo.Queries) SQLSessionRepo {
	return SQLSessionRepo{q: q}
}

func repoSessionToDomain(s *repo.Session) *domain.Session {
	return &domain.Session{ID: s.ID, ValidUntil: s.ValidUntil, CreatedAt: s.CreatedAt, EditedAt: s.EditedAt, UserID: s.UserID}
}

func (r *SQLSessionRepo) Create(ctx context.Context, userId int64, secureRand string, validUntil time.Time) (*domain.Session, error) {
	data := repo.CreateSessionParams{ID: secureRand, ValidUntil: validUntil, UserID: userId}
	session, err := r.q.CreateSession(ctx, data)
	if err != nil {
		log.Error("CreateSession failed", slog.String("error", err.Error()))
		return &domain.Session{}, err
	}

	return repoSessionToDomain(&session), nil
}

func (r *SQLSessionRepo) Delete(ctx context.Context, cookie string) error {
	err := r.q.DeleteSession(ctx, cookie)
	if err != nil {
		log.Error("DeleteSession failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *SQLSessionRepo) DeleteAll(ctx context.Context) error {
	err := r.q.DeleteAllSessions(ctx)
	if err != nil {
		log.Error("DeleteAllSessions failed", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *SQLSessionRepo) GetById(ctx context.Context, cookie string) (*domain.Session, error) {
	session, err := r.q.GetSessionById(ctx, cookie)
	if err != nil {
		log.Error("GetSessionById failed", slog.String("error", err.Error()))
		return nil, err
	}

	return repoSessionToDomain(&session), nil
}

func (r *SQLSessionRepo) GetSessionUser(ctx context.Context, cookie string) (*domain.User, error) {
	u, err := r.q.GetUserFromSession(ctx, cookie)
	if err != nil {
		log.Error("GetValidSession failed", slog.String("error", err.Error()))
		return nil, err
	}

	return repoUserToDomain(&u), nil
}
