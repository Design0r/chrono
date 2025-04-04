package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/calendar"
	"chrono/db/repo"
)

func CreateUser(r *repo.Queries, data repo.CreateUserParams) (repo.User, error) {
	user, err := r.CreateUser(context.Background(), data)
	if err != nil {
		log.Printf("Failed creating user: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func UpdateUser(r *repo.Queries, data repo.UpdateUserParams) (repo.User, error) {
	user, err := r.UpdateUser(context.Background(), data)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func GetUserById(r *repo.Queries, id int64) (repo.User, error) {
	user, err := r.GetUserByID(context.Background(), id)
	if err != nil {
		log.Printf("Failed getting user: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func GetUserByName(r *repo.Queries, name string) (repo.User, error) {
	user, err := r.GetUserByName(context.Background(), name)
	if err != nil {
		log.Printf("Failed getting user: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func GetUserByEmail(r *repo.Queries, email string) (repo.User, error) {
	user, err := r.GetUserByEmail(context.Background(), email)
	if err != nil {
		log.Printf("Failed getting user: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func DeleteUser(r *repo.Queries, id int64) error {
	err := r.DeleteUser(context.Background(), id)
	if err != nil {
		log.Printf("Failed getting user: %v", err)
		return err
	}

	return nil
}

func GetCurrentUser(r *repo.Queries, c echo.Context) (repo.User, error) {
	session, err := c.Cookie("session")
	if err != nil {
		return repo.User{}, err
	}
	return GetUserFromSession(r, session.Value)
}

func GetAllVacUsers(r *repo.Queries) ([]repo.GetUsersWithVacationCountRow, error) {
	start := time.Date(calendar.CurrentYear(), 1, 1, 0, 0, 0, 0, time.Now().Location())

	data := repo.GetUsersWithVacationCountParams{
		StartDate: start,
		EndDate:   start,
	}
	users, err := r.GetUsersWithVacationCount(context.Background(), data)
	if err != nil {
		return nil, err
	}

	for i := range users {
		vac, err := GetVacationCountForUser(r, users[i].ID, start.Year())
		if err != nil {
			continue
		}
		users[i].VacUsed = &vac
	}

	return users, err
}

func ToggleAdmin(r *repo.Queries, editor repo.User, userId int64) (repo.User, error) {
	user, err := r.ToggleAdmin(context.Background(), userId)
	if err != nil {
		return repo.User{}, err
	}

	var msg string
	if !user.IsSuperuser {
		msg = "revoked your admin status"
	} else {
		msg = "gave you admin status"
	}

	_, err = CreateUserNotification(r, fmt.Sprintf("%v %v", editor.Username, msg), userId)
	if err != nil {
		return repo.User{}, err
	}

	return user, nil
}

func GetAllUsers(r *repo.Queries) ([]repo.User, error) {
	users, err := r.GetAllUsers(context.Background())
	if err != nil {
		log.Printf("Failed getting all users: %v", err)
		return nil, err
	}

	return users, nil
}

func SetUserVacation(r *repo.Queries, userId int64, vacation int, year int) error {
	oldUser, err := GetUserById(r, userId)
	if err != nil {
		return err
	}

	params := repo.SetUserVacationParams{ID: userId, VacationDays: int64(vacation)}
	updatedUser, err := r.SetUserVacation(context.Background(), params)
	if err != nil {
		log.Printf("Failed updating user vacation: %v", err)
		return err
	}

	err = UpdateYearlyTokens(r, updatedUser, year, vacation-int(oldUser.VacationDays))
	if err != nil {
		return err
	}

	return nil
}

func SetUserColor(r *repo.Queries, userId int64, color string) error {
	params := repo.SetUserColorParams{Color: color, ID: userId}
	err := r.SetUserColor(context.Background(), params)
	if err != nil {
		log.Printf("Failed to set user color: %v", err)
		return err
	}

	return nil
}
