package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"chrono/db/repo"
	"chrono/schemas"
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

func GetPendingRequests(r *repo.Queries) ([]schemas.BatchRequest, error) {
	req, err := r.GetPendingRequests(context.Background())
	if err != nil {
		return nil, err
	}

	if len(req) == 0 {
		return nil, nil
	}

	var requestsToShow []schemas.BatchRequest

	startIndex := 0
	for startIndex < len(req) {
		endIndex := startIndex

		for endIndex+1 < len(req) &&
			req[endIndex].ScheduledAt.Year() == req[endIndex+1].ScheduledAt.Year() &&
			req[endIndex].ScheduledAt.YearDay()+1 == req[endIndex+1].ScheduledAt.YearDay() {
			endIndex++
		}

		requestsToShow = append(requestsToShow, schemas.BatchRequest{
			StartDate:  req[startIndex].ScheduledAt,
			EndDate:    req[endIndex].ScheduledAt,
			EventCount: endIndex - startIndex + 1,
			Request:    &req[endIndex],
		})

		startIndex = endIndex + 1
	}

	return requestsToShow, nil
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

func UpdateRequestStateRange(
	r *repo.Queries,
	userId int64,
	state string,
	startDate time.Time,
	endDate time.Time,
) error {
	params := repo.UpdateRequestStateRangeParams{
		UserID:        userId,
		State:         state,
		ScheduledAt:   startDate,
		ScheduledAt_2: endDate,
	}
	err := r.UpdateRequestStateRange(context.Background(), params)
	if err != nil {
		log.Printf("Failed to update requests: %v", err)
	}

	return err
}
