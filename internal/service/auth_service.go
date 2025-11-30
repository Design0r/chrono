package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"chrono/internal/domain"
	"chrono/internal/service/auth"
)

type AuthService struct {
	user            domain.UserRepository
	session         domain.SessionRepository
	pw              auth.PasswordHasher
	sessionDuration time.Duration
	log             *slog.Logger
}

func NewAuthService(
	u domain.UserRepository,
	s domain.SessionRepository,
	sessionDuration time.Duration,
	pw auth.PasswordHasher,
	log *slog.Logger,
) AuthService {
	return AuthService{user: u, session: s, log: log, sessionDuration: sessionDuration, pw: pw}
}

func (svc *AuthService) createSessionCookie(session domain.Session) *http.Cookie {
	return &http.Cookie{
		Path:     "/",
		Name:     "session",
		Value:    session.ID,
		HttpOnly: true,
		Secure:   false,
		Expires:  session.ValidUntil,
		SameSite: http.SameSiteStrictMode,
	}
}

func (svc *AuthService) deleteSessionCookie() *http.Cookie {
	return &http.Cookie{
		Path:     "/",
		Name:     "session",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	}
}

func (svc *AuthService) Login(ctx context.Context, email, pw string) (*http.Cookie, error) {
	user, err := svc.user.GetByEmail(ctx, email)
	if err != nil {
		svc.log.Error(
			"Login failed",
			slog.String("email", email),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	ok := svc.pw.Compare(user.Password, pw)
	if !ok {
		svc.log.Error("Login failed, incorrect password or email", slog.String("email", email))
		return nil, errors.New("passwords do not match")
	}

	session, err := svc.session.Create(ctx, user.ID, svc.pw.SecureRandom64(), svc.sessionDuration)
	if err != nil {
		svc.log.Error(
			"Login failed",
			slog.String("email", email),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return svc.createSessionCookie(*session), nil
}

func (svc *AuthService) Logout(ctx context.Context, cookie string) (*http.Cookie, error) {
	err := svc.session.Delete(ctx, cookie)
	if err != nil {
		return nil, err
	}

	return svc.deleteSessionCookie(), nil
}

func (svc *AuthService) Signup(
	ctx context.Context,
	userParams domain.CreateUser,
) (*http.Cookie, error) {
	_, err := svc.user.GetByEmail(ctx, userParams.Email)
	if err == nil {
		svc.log.Error("User with this email already exists", slog.String("email", userParams.Email))
		return nil, errors.New("User with this email already exists")
	}

	hashedPw, err := svc.pw.Hash(userParams.Password)
	if err != nil {
		svc.log.Error(
			"Signup failed",
			slog.String("email", userParams.Email),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	userParams.Password = hashedPw
	userParams.Color = domain.Color.RandomHexColor()

	user, err := svc.user.Create(ctx, &userParams)
	if err != nil {
		svc.log.Error(
			"Signup Failed",
			slog.String("email", userParams.Email),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	session, err := svc.session.Create(ctx, user.ID, svc.pw.SecureRandom64(), svc.sessionDuration)
	if err != nil {
		svc.log.Error(
			"Login failed",
			slog.String("email", userParams.Email),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return svc.createSessionCookie(*session), nil
}

func (svc *AuthService) GetCurrentUser(ctx context.Context, cookie string) (*domain.User, error) {
	return svc.session.GetSessionUser(ctx, cookie)
}

func (svc *AuthService) HashPassword(password string) (string, error) {
	return svc.pw.Hash(password)
}
