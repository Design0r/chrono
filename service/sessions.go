package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"chrono/db/repo"
)

func CreateSession(db *sql.DB, userId int64) (repo.Session, error) {
	r := repo.New(db)

	secRand := SecureRandom(64)

	date := time.Now().Add(time.Hour * 24 * 7) // 1 week
	data := repo.CreateSessionParams{ID: secRand, ValidUntil: date, UserID: userId}
	session, err := r.CreateSession(context.Background(), data)
	if err != nil {
		log.Printf("Failed creating Session: %v", err)
		return repo.Session{}, err
	}

	return session, nil
}

func GetUserFromSession(db *sql.DB, id string) (repo.User, error) {
	r := repo.New(db)

	user, err := r.GetUserFromSession(context.Background(), id)
	if err != nil {
		log.Printf("Failed getting session: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func DeleteSession(db *sql.DB, id string) error {
	r := repo.New(db)

	err := r.DeleteSession(context.Background(), id)
	if err != nil {
		log.Printf("Failed deleting session: %v", err)
		return err
	}

	return nil
}

func IsValidSession(db *sql.DB, id string) bool {
	r := repo.New(db)

	_, err := r.GetValidSession(
		context.Background(),
		repo.GetValidSessionParams{ID: id, ValidUntil: time.Now()},
	)
	if err != nil {
		log.Printf("Failed getting valid session: %v", err)
		return false
	}

	return true
}
