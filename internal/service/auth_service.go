package service

import (
	"chrono/internal/domain"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	HashPassword(pw string) (string, error)
	ComparePasswords(hashedPw, pw string) bool
	SecureRandom(length int) string
	SecureRandom64() string
	GetCurrentUser()
	CreateSessionCookie(s domain.Session) *http.Cookie
	DeleteSessionCookie(s domain.Session) *http.Cookie
	Login(ctx context.Context, email, password string) (*http.Cookie, error)
	Logout(ctx context.Context) (*http.Cookie, error)
	Signup(ctx context.Context, userParams domain.CreateUser) (*http.Cookie, error)
}

type authService struct {
	u               domain.UserRepository
	s               domain.SessionRepository
	sessionDuration time.Duration
	log             *slog.Logger
}

func (svc *authService) NewAuthService(u domain.UserRepository, s domain.SessionRepository, sessionDuration time.Duration, log *slog.Logger) authService {
	return authService{u: u, s: s, log: log, sessionDuration: sessionDuration}
}

func (svc *authService) SecureRandom(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(randomBytes)
}

func (svc *authService) SecureRandom64() string {
	return svc.SecureRandom(64)
}

func (svc *authService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (svc *authService) ComparePasswords(hashedPw, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(pw))
	return err == nil // Returns true if the password matches
}

func (s *authService) CreateSessionCookie(session domain.Session) *http.Cookie {
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

func (s *authService) DeleteSessionCookie() *http.Cookie {
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

func (svc *authService) Login(ctx context.Context, email, pw string) (*http.Cookie, error) {
	user, err := svc.u.GetByEmail(ctx, email)
	if err != nil {
		svc.log.Error("Login failed", slog.String("email", email), slog.String("error", err.Error()))
		return nil, err
	}

	ok := svc.ComparePasswords(user.Password, pw)
	if !ok {
		svc.log.Error("Login failed, incorrect password or email", slog.String("email", email), slog.String("error", err.Error()))
		return nil, errors.New("passwords do not match")
	}

	session, err := svc.s.Create(ctx, user.ID, svc.SecureRandom64(), time.Now().Add(time.Hour*24*7))
	if err != nil {
		svc.log.Error("Login failed", slog.String("email", email), slog.String("error", err.Error()))
		return nil, err
	}

	return svc.CreateSessionCookie(*session), nil
}

func (svc *authService) Logout(ctx context.Context, cookie string) (*http.Cookie, error) {
	err := svc.s.Delete(ctx, cookie)
	if err != nil {
		return nil, err
	}

	return svc.DeleteSessionCookie(), nil
}

func (svc *authService) Signup(ctx context.Context, userParams domain.CreateUser) (*http.Cookie, error) {
	_, err := svc.u.GetByEmail(ctx, userParams.Email)
	if err == nil {
		svc.log.Error("User with this email already exists", slog.String("email", userParams.Email), slog.String("error", err.Error()))
		return nil, err
	}

	hashedPw, err := svc.HashPassword(userParams.Password)
	if err != nil {
		svc.log.Error("Signup failed", slog.String("email", userParams.Email), slog.String("error", err.Error()))
		return nil, err
	}
	userParams.Password = hashedPw
	userParams.Color = domain.Color.RandomHexColor()

	user, err := svc.u.Create(ctx, &userParams)
	if err != nil {
		svc.log.Error("Signup Failed", slog.String("email", userParams.Email), slog.String("error", err.Error()))
		return nil, err
	}

	session, err := svc.s.Create(ctx, user.ID, svc.SecureRandom64(), time.Now().Add(time.Hour*24*7))
	if err != nil {
		svc.log.Error("Login failed", slog.String("email", userParams.Email), slog.String("error", err.Error()))
		return nil, err
	}

	return svc.CreateSessionCookie(*session), nil

}
