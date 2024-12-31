package service

import (
	"context"
	"log"
	"time"

	"chrono/db/repo"
)

func CreateSession(r *repo.Queries, userId int64) (repo.Session, error) {
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

func GetUserFromSession(r *repo.Queries, id string) (repo.User, error) {
	user, err := r.GetUserFromSession(context.Background(), id)
	if err != nil {
		log.Printf("Failed getting session: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func DeleteSession(r *repo.Queries, id string) error {
	err := r.DeleteSession(context.Background(), id)
	if err != nil {
		log.Printf("Failed deleting session: %v", err)
		return err
	}

	return nil
}

func IsValidSession(r *repo.Queries, id string) bool {
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
