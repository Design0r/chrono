package service

import (
	"context"
	"log/slog"
	"os"

	"chrono/internal/domain"
	"chrono/internal/service/auth"
)

type APIBot struct {
	Name        string
	Email       string
	Password    string
	IsSuperuser bool
	log         *slog.Logger
}

func NewAPIBotFromEnv(log *slog.Logger) APIBot {
	name, exists := os.LookupEnv("BOT_NAME")
	if !exists {
		log.Error("BOT_NAME env var missing.")
	}

	pw, exists := os.LookupEnv("BOT_PASSWORD")
	if !exists {
		log.Error("BOT_PASSWORD env var missing.")
	}

	email, exists := os.LookupEnv("BOT_EMAIL")
	if !exists {
		log.Error("BOT_EMAIL env var missing.")
	}

	return APIBot{Name: name, Email: email, Password: pw, IsSuperuser: true, log: log}
}

func (a *APIBot) Register(svc *UserService, pw auth.PasswordHasher) {
	ctx := context.Background()
	_, err := svc.GetByEmail(ctx, a.Email)
	if err == nil {
		a.log.Error("User with email already exists", slog.String("email", a.Email))
		return
	}

	hashedPw, err := pw.Hash(a.Password)
	if err != nil {
		a.log.Error("Failed hashing password", slog.String("name", a.Name))
		return
	}
	_, err = svc.Create(
		ctx,
		&domain.CreateUser{
			Username:     a.Name,
			Email:        a.Email,
			Color:        domain.Color.RandomHexColor(),
			Password:     hashedPw,
			IsSuperuser:  a.IsSuperuser,
			VacationDays: 0,
		},
	)
	if err != nil {
		a.log.Error("User with email already exists", slog.String("email", a.Email))
		return
	}
	a.log.Info("Created Bot user", slog.String("name", a.Name))
}
