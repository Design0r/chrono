package service

import (
	"database/sql"
	"log"

	"calendar/db/repo"
)

func InitAPIBot(db *sql.DB) error {
	botName := "Chrono Bot"
	_, err := GetUserByName(db, botName)
	if err != nil {
		_, err = CreateUser(
			db,
			repo.CreateUserParams{
				Username:     botName,
				Password:     "chrono",
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
