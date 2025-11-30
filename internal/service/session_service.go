package service

import (
	"context"
	"log/slog"
	"time"

	"chrono/internal/domain"
)

type SessionService struct {
	session domain.SessionRepository
	log     *slog.Logger
}

func NewSessionService(r domain.SessionRepository, log *slog.Logger) SessionService {
	return SessionService{session: r, log: log}
}

func (svc *SessionService) Create(
	ctx context.Context,
	userId int64,
	secureRand string,
	duration time.Duration,
) (*domain.Session, error) {
	return svc.session.Create(ctx, userId, secureRand, duration)
}

func (svc *SessionService) Delete(ctx context.Context, cookie string) error {
	return svc.session.Delete(ctx, cookie)
}

func (svc *SessionService) DeleteAll(ctx context.Context) error {
	return svc.session.DeleteAll(ctx)
}

func (svc *SessionService) IsValidSession(
	ctx context.Context,
	cookie string,
	timestamp time.Time,
) bool {
	session, err := svc.session.GetById(ctx, cookie)
	if err != nil {
		return false
	}

	// if timestamp before ValidUntil -> true
	return timestamp.Compare(session.ValidUntil) <= 0
}

func (svc *SessionService) GetUserFromSession(
	ctx context.Context,
	cookie string,
) (*domain.User, error) {
	return svc.session.GetSessionUser(ctx, cookie)
}
