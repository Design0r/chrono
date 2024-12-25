package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"calendar/db/repo"
)

func CreateNotification(db *sql.DB, msg string, userId int64) (repo.Notification, error) {
	r := repo.New(db)
	data := repo.CreateNotificationParams{Message: msg, UserID: userId}

	n, err := r.CreateNotification(context.Background(), data)
	if err != nil {
		log.Printf("Failed to create notification: %v", err)
		return repo.Notification{}, err
	}

	return n, nil
}

func ClearAllNotifications(db *sql.DB, userId int64) error {
	r := repo.New(db)

	err := r.ClearAllNotification(context.Background(), userId)
	if err != nil {
		log.Printf("Failed to clear notifications: %v", err)
		return err
	}

	return nil
}

func ClearNotification(db *sql.DB, notifId int64) error {
	r := repo.New(db)

	err := r.ClearNotification(context.Background(), notifId)
	if err != nil {
		log.Printf("Failed to clear notification: %v", err)
		return err
	}

	return nil
}

func GetUserNotifications(db *sql.DB, userId int64) ([]repo.Notification, error) {
	r := repo.New(db)

	n, err := r.GetUserNotifications(context.Background(), userId)
	if err != nil {
		log.Printf("Failed to clear notification: %v", err)
		return []repo.Notification{}, err
	}

	return n, nil
}

func GenerateRequestMsg(username string, event repo.Event) string {
	return fmt.Sprintf("%v sent a new request for %v!", username, event.Name)
}

func GenerateAcceptMsg(username string, event repo.Event) string {
	return fmt.Sprintf("%v accepted your %v request!", username, event.Name)
}

func GenerateRejectMsg(username string, event repo.Event) string {
	return fmt.Sprintf("%v rejected your %v request!", username, event.Name)
}

func GenerateUpdateMsg(username string, state string, event repo.Event) string {
	return fmt.Sprintf("%v %v your %v request!", username, state, event.Name)
}
