package service

import (
	"log"
	"os"

	"chrono/db/repo"
)

type APIBot struct {
	Name        string
	Email       string
	Password    string
	IsSuperuser bool
}

func NewAPIBotFromEnv() APIBot {
	name, exists := os.LookupEnv("BOT_NAME")
	if !exists {
		log.Fatal("BOT_NAME env var missing.")
	}

	pw, exists := os.LookupEnv("BOT_PASSWORD")
	if !exists {
		log.Fatal("BOT_PASSWORD env var missing.")
	}

	email, exists := os.LookupEnv("BOT_EMAIL")
	if !exists {
		log.Fatal("BOT_EMAIL env var missing.")
	}

	return APIBot{Name: name, Email: email, Password: pw, IsSuperuser: true}
}

func (self *APIBot) Register(r *repo.Queries) error {
	_, err := GetUserByEmail(r, self.Email)
	if err == nil {
		return nil
	}

	hashedPw, err := HashPassword(self.Password)
	if err != nil {
		log.Fatalf("Failed hashing %v password", self.Name)
	}
	_, err = CreateUser(
		r,
		repo.CreateUserParams{
			Username:     self.Name,
			Email:        self.Email,
			Password:     hashedPw,
			IsSuperuser:  self.IsSuperuser,
			VacationDays: 0,
		},
	)
	if err != nil {
		log.Printf("User with email %v already exists", self.Email)
		return err
	}
	log.Printf("Created %v Bot user", self.Name)
	return nil
}
