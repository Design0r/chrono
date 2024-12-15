package service

import (
	"database/sql"

	"calendar/db/repo"
)

func InitAPIBot(db *sql.DB) error {
	botName := "APIBot"
	_, err := GetUserByName(db, botName)
	if err != nil {
		CreateUser(
			db,
			repo.CreateUserParams{
				Username:     botName,
				Password:     "apibot",
				IsSuperuser:  false,
				VacationDays: 100000000,
			},
		)
	}

	return nil
}
