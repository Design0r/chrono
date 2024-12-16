package service

import (
	"database/sql"
	"log"

	"calendar/db/repo"
)

func InitAPIBot(db *sql.DB) error {
	botName := "APIBot"
	_, err := GetUserByName(db, botName)
	if err != nil {
		_, err = CreateUser(
			db,
			repo.CreateUserParams{
				Username:     botName,
				Password:     "apibot",
				IsSuperuser:  false,
				VacationDays: 100000000,
			},
		)
		if err != nil {
			log.Println("User APIBot already exists")
			return err
		}
		log.Println("Created APIBot user")
	}

	return nil
}
