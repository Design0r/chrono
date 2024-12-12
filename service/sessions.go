package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"

	"calendar/db/repo"
)

func CreateSession(db *sql.DB, userId int64) (repo.Session, error) {
	r := repo.New(db)

	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Failed creating uuid: %v", err)
		return repo.Session{}, err
	}

	date := time.Now().Add(time.Hour * 24 * 7) // 1 week
	data := repo.CreateSessionParams{ID: uuid.String(), ValidUntil: date, UserID: userId}
	session, err := r.CreateSession(context.Background(), data)
	if err != nil {
		log.Printf("Failed creating Session: %v", err)
		return repo.Session{}, err
	}

	return session, nil
}

func GetUserFromSession(db *sql.DB, id uuid.UUID) (repo.User, error) {
	r := repo.New(db)

	user, err := r.GetUserFromSession(context.Background(), id.String())
	if err != nil {
		log.Printf("Failed getting session: %v", err)
		return repo.User{}, err
	}

	return user, nil
}

func DeleteSession(db *sql.DB, id uuid.UUID) error {
	r := repo.New(db)

	err := r.DeleteSession(context.Background(), id.String())
	if err != nil {
		log.Printf("Failed deleting session: %v", err)
		return err
	}

	return nil
}
