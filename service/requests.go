package service

import (
	"context"
	"database/sql"
	"log"

	"calendar/db/repo"
)

func CreateRequest(db *sql.DB, msg string, user repo.User, event repo.Event) (repo.Request, error) {
	r := repo.New(db)
	data := repo.CreateRequestParams{
		Message: &msg,
		State:   "pending",
		UserID:  user.ID,
		EventID: event.ID,
	}

	request, err := r.CreateRequest(context.Background(), data)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return repo.Request{}, err
	}

	_, err = CreateNotification(db, GenerateRequestMsg(user.Username), user.ID)
	if err != nil {
		return repo.Request{}, err
	}

	return request, nil
}

func GetPendingRequests(db *sql.DB) ([]repo.GetPendingRequestsRow, error) {
	r := repo.New(db)

	req, err := r.GetPendingRequests(context.Background())
	if err != nil {
		return nil, err
	}

	return req, nil
}
