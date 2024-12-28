package service

import (
	"database/sql"
	"log"
	"os"

	"chrono/db/repo"
)

func InitAPIBot(db *sql.DB) error {
	botName := "Chrono Bot"
	_, err := GetUserByName(db, botName)
	if err != nil {

		envPw, exists := os.LookupEnv("CHRONO_PASSWORD")
		if !exists {
			log.Fatal("CHRONO_PASSWORD env var missing.")
		}

		envEmail, exists := os.LookupEnv("CHRONO_EMAIL")
		if !exists {
			log.Fatal("CHRONO_EMAIL env var missing.")
		}

		hashedPw, err := HashPassword(envPw)
		if err != nil {
			log.Fatal("Failed hashing Chrono pw")
		}
		_, err = CreateUser(
			db,
			repo.CreateUserParams{
				Username:     botName,
				Email:        envEmail,
				Password:     hashedPw,
				IsSuperuser:  true,
				VacationDays: 0,
			},
		)
		if err != nil {
			log.Println("User Chrono Bot already exists")
			return err
		}
		log.Println("Created Chrono Bot user")
	}

	return nil
}
