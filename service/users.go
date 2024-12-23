package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"calendar/db/repo"
)

func CreateUser(db *sql.DB, data repo.CreateUserParams) (repo.User, error) {
	r := repo.New(db)

	user, err := r.CreateUser(context.Background(), data)
	if err != nil {
		log.Printf("Failed creating user: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func GetUserById(db *sql.DB, id int64) (repo.User, error) {
	r := repo.New(db)

	user, err := r.GetUserByID(context.Background(), id)
	if err != nil {
		log.Printf("Failed getting user: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func GetUserByName(db *sql.DB, name string) (repo.User, error) {
	r := repo.New(db)

	user, err := r.GetUserByName(context.Background(), name)
	if err != nil {
		log.Printf("Failed getting user: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func GetUserByEmail(db *sql.DB, email string) (repo.User, error) {
	r := repo.New(db)

	user, err := r.GetUserByEmail(context.Background(), email)
	if err != nil {
		log.Printf("Failed getting user: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func DeleteUser(db *sql.DB, id int64) error {
	r := repo.New(db)

	err := r.DeleteUser(context.Background(), id)
	if err != nil {
		log.Printf("Failed getting user: %v", err)
		return err
	}

	return nil
}

func GetCurrentUser(db *sql.DB, c echo.Context) (repo.User, error) {
	session, err := c.Cookie("session")
	if err != nil {
		return repo.User{}, err
	}
	return GetUserFromSession(db, session.Value)
}

func GetAllVacUsers(db *sql.DB) ([]repo.GetUsersWithVacationCountRow, error) {
	r := repo.New(db)

	start := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Now().Location())

	data := repo.GetUsersWithVacationCountParams{
		ScheduledAt:   start,
		ScheduledAt_2: start.AddDate(1, 0, 0),
	}
	users, err := r.GetUsersWithVacationCount(context.Background(), data)
	if err != nil {
		return nil, err
	}

	return users, err
}

func HashCode(s string) int {
	hash := 0
	for _, char := range s {
		hash = int(char) + ((hash << 5) - hash)
	}
	return hash
}

func GenerateHSL(seed int) string {
	hue := (seed * 12345) % 360
	saturation := 50 + (seed % 50)
	lightness := 40 + (seed % 20)
	return fmt.Sprintf("hsl(%d, %d%%, %d%%)", hue, saturation, lightness)
}
