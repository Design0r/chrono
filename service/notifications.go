package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"chrono/db/repo"
)

func _createNotification(db *sql.DB, msg string) (repo.Notification, error) {
	r := repo.New(db)

	n, err := r.CreateNotification(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to create notification: %v", err)
		return repo.Notification{}, err
	}

	return n, nil
}

func CreateUserNotification(db *sql.DB, msg string, userId int64) (repo.Notification, error) {
	r := repo.New(db)

	n, err := r.CreateNotification(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to create notification: %v", err)
		return repo.Notification{}, err
	}

	params := repo.CreateNotificationUserParams{NotificationID: n.ID, UserID: userId}
	err = r.CreateNotificationUser(context.Background(), params)
	if err != nil {
		log.Printf("Failed to create notification user association: %v", err)
	}

	return n, nil
}

func CreateAdminNotification(db *sql.DB, msg string) (repo.Notification, error) {
	r := repo.New(db)
	ctx := context.Background()

	n, err := _createNotification(db, msg)
	if err != nil {
		return repo.Notification{}, err
	}

	admins, err := r.GetAdmins(ctx)
	if err != nil {
		log.Printf("Failed to get admin users: %v", admins)
		return repo.Notification{}, err
	}

	for _, a := range admins {
		params := repo.CreateNotificationUserParams{NotificationID: n.ID, UserID: a.ID}
		err := r.CreateNotificationUser(ctx, params)
		if err != nil {
			log.Printf("Failed to create notification user association: %v", err)
		}
	}

	return n, nil
}

func ClearAllNotifications(db *sql.DB, userId int64) error {
	r := repo.New(db)

	err := r.ClearAllUserNotifications(context.Background(), userId)
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
