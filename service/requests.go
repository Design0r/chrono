package service

import (
	"context"
	"database/sql"
	"fmt"
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

	_, err = CreateAdminNotification(db, GenerateRequestMsg(user.Username, event))
	if err != nil {
		return repo.Request{}, err
	}

	return request, nil
}

func GetPendingRequests(db *sql.DB) ([]repo.GetPendingRequestsRow, error) {
	r := repo.New(db)

	req, err := r.GetPendingRequests(context.Background())
	if err != nil {
		return []repo.GetPendingRequestsRow{}, err
	}

	return req, nil
}

func UpdateRequestState(db *sql.DB, state string, currUser repo.User, reqId int64) error {
	r := repo.New(db)

	if state != "accepted" && state != "declined" {
		return fmt.Errorf("Invalid state: %v", state)
	}

	data := repo.UpdateRequestStateParams{State: state, EditedBy: &currUser.ID, ID: reqId}
	req, err := r.UpdateRequestState(context.Background(), data)
	if err != nil {
		return err
	}

	event, err := UpdateEventState(db, state, req.EventID)
	if err != nil {
		return err
	}

	_, err = CreateUserNotification(
		db,
		GenerateUpdateMsg(currUser.Username, state, event),
		req.UserID,
	)
	if err != nil {
		return err
	}

	return nil
}
