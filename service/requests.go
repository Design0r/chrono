package service

import (
	"context"
	"fmt"
	"log"

	"chrono/db/repo"
)

func CreateRequest(
	r *repo.Queries,
	msg string,
	user repo.User,
	event repo.Event,
) (repo.Request, error) {
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

	_, err = CreateAdminNotification(r, GenerateRequestMsg(user.Username, event))
	if err != nil {
		return repo.Request{}, err
	}

	return request, nil
}

func GetPendingRequests(r *repo.Queries) ([]repo.GetPendingRequestsRow, error) {
	req, err := r.GetPendingRequests(context.Background())
	if err != nil {
		return []repo.GetPendingRequestsRow{}, err
	}

	return req, nil
}

func UpdateRequestState(r *repo.Queries, state string, currUser repo.User, reqId int64) error {
	if state != "accepted" && state != "declined" {
		return fmt.Errorf("Invalid state: %v", state)
	}

	data := repo.UpdateRequestStateParams{State: state, EditedBy: &currUser.ID, ID: reqId}
	req, err := r.UpdateRequestState(context.Background(), data)
	if err != nil {
		return err
	}

	event, err := UpdateEventState(r, state, req.EventID)
	if err != nil {
		return err
	}

	_, err = CreateUserNotification(
		r,
		GenerateUpdateMsg(currUser.Username, state, event),
		req.UserID,
	)
	if err != nil {
		return err
	}

	return nil
}
