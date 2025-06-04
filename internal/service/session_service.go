package service

import (
	"context"
	"time"

	"chrono/internal/domain"
)

type SessionService interface {
	Create(ctx context.Context, userId int64, secureRand string, validUntil time.Time) (*domain.Session, error)
	Delete(ctx context.Context, cookie string) error
	DeleteAll(ctx context.Context) error
	IsValidSession(ctx context.Context, cookie string, timestamp time.Time) bool
	GetUserFromSession(ctx context.Context, cookie string) (*domain.User, error)
}

type sessionService struct {
	r domain.SessionRepository
}

func NewSessionService(r domain.SessionRepository) sessionService {
	return sessionService{r: r}
}

func (s *sessionService) Create(ctx context.Context, userId int64, secureRand string, validUntil time.Time) (*domain.Session, error) {
	return s.r.Create(ctx, userId, secureRand, validUntil)
}

func (s *sessionService) Delete(ctx context.Context, cookie string) error {
	return s.r.Delete(ctx, cookie)
}

func (s *sessionService) DeleteAll(ctx context.Context) error {
	return s.r.DeleteAll(ctx)
}

func (s *sessionService) IsValidSession(ctx context.Context, cookie string, timestamp time.Time) bool {
	session, err := s.r.GetById(ctx, cookie)
	if err != nil {
		return false
	}

	// if timestamp before ValidUntil -> true
	return timestamp.Compare(session.ValidUntil) <= 0

}

func (s *sessionService) GetUserFromSession(ctx context.Context, cookie string) (*domain.User, error) {
	return s.r.GetUserFromSession(ctx, cookie)
}
