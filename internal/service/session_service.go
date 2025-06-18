package service

import (
	"context"
	"log/slog"
	"time"

	"chrono/internal/domain"
)

type SessionService interface {
	Create(ctx context.Context, userId int64, secureRand string, duration time.Duration) (*domain.Session, error)
	Delete(ctx context.Context, cookie string) error
	DeleteAll(ctx context.Context) error
	IsValidSession(ctx context.Context, cookie string, timestamp time.Time) bool
	GetUserFromSession(ctx context.Context, cookie string) (*domain.User, error)
}

type sessionService struct {
	session domain.SessionRepository
	log     *slog.Logger
}

func NewSessionService(r domain.SessionRepository, log *slog.Logger) sessionService {
	return sessionService{session: r, log: log}
}

func (svc *sessionService) Create(ctx context.Context, userId int64, secureRand string, duration time.Duration) (*domain.Session, error) {
	return svc.session.Create(ctx, userId, secureRand, duration)
}

func (svc *sessionService) Delete(ctx context.Context, cookie string) error {
	return svc.session.Delete(ctx, cookie)
}

func (svc *sessionService) DeleteAll(ctx context.Context) error {
	return svc.session.DeleteAll(ctx)
}

func (svc *sessionService) IsValidSession(ctx context.Context, cookie string, timestamp time.Time) bool {
	session, err := svc.session.GetById(ctx, cookie)
	if err != nil {
		return false
	}

	// if timestamp before ValidUntil -> true
	return timestamp.Compare(session.ValidUntil) <= 0

}

func (svc *sessionService) GetUserFromSession(ctx context.Context, cookie string) (*domain.User, error) {
	return svc.session.GetSessionUser(ctx, cookie)
}
