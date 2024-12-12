package service

import (
	"context"
	"database/sql"
	"log"

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
