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
			req[endIndex].ScheduledAt.YearDay()+1 == req[endIndex+1].ScheduledAt.YearDay() &&
			req[endIndex].UserID == req[endIndex+1].UserID {
			endIndex++
		}

		startDate := req[startIndex].ScheduledAt
		endDate := req[endIndex].ScheduledAt

		fmt.Println(startDate, endDate)

		confilctingUsers, _ := GetConflictingEventUsers(
			r,
			req[startIndex].UserID,
			startDate,
			endDate,
		)

		requestsToShow = append(requestsToShow, schemas.BatchRequest{
			StartDate:  startDate,
			EndDate:    endDate,
			EventCount: endIndex - startIndex + 1,
			Request:    &req[endIndex],
			Conflicts:  &confilctingUsers,
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
	editor repo.User,
	userId int64,
	state string,
	startDate time.Time,
	endDate time.Time,
	reason string,
) error {
	params := repo.UpdateRequestStateRangeParams{
		UserID:        userId,
		EditedBy:      &editor.ID,
		State:         state,
		ScheduledAt:   startDate,
		ScheduledAt_2: endDate,
	}
	err := r.UpdateRequestStateRange(context.Background(), params)
	if err != nil {
		log.Printf("Failed to update requests: %v", err)
	}

	params2 := repo.UpdateEventsRangeParams{
		UserID:        userId,
		State:         state,
		ScheduledAt:   startDate,
		ScheduledAt_2: endDate,
	}

	err = r.UpdateEventsRange(context.Background(), params2)
	if err != nil {
		log.Printf("Failed to update events: %v", err)
	}

	msg := GenerateBatchUpdateMsg(editor.Username, state)
	if reason != "" {
		msg = GenerateBatchUpdateReasonMsg(editor.Username, state, reason)
	}

	_, err = CreateUserNotification(r, msg, userId)

	return err
}

func GetRequestRange(
	r *repo.Queries,
	startDate time.Time,
	endDate time.Time,
	userID int64,
) ([]repo.GetRequestRangeRow, error) {
	params := repo.GetRequestRangeParams{
		UserID:        userID,
		ScheduledAt:   startDate,
		ScheduledAt_2: endDate,
	}

	requests, err := r.GetRequestRange(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func GetConflictingEventUsers(
	r *repo.Queries,
	userId int64,
	startDate time.Time,
	endDate time.Time,
) ([]repo.User, error) {
	params := repo.GetConflictingEventUsersParams{
		ID:            userId,
		ScheduledAt:   startDate,
		ScheduledAt_2: endDate,
	}
	users, err := r.GetConflictingEventUsers(context.Background(), params)
	if err != nil {
		log.Printf("Failed to fetch conflicting users: %v", err)
		return []repo.User{}, err
	}

	return users, nil
}
