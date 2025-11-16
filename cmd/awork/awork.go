package main

import (
	"fmt"
	"log/slog"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"chrono/config"
	"chrono/internal/service"
)

func main() {
	config.NewConfigFromEnv()
	awork := service.NewAworkService(slog.Default())
	users, _ := awork.GetUsers()

	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), 10, 0, 0, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), 15, 0, 0, 0, 0, now.Location())

	for _, u := range users {
		entries, _ := awork.GetTimeEntries(u.Id, startTime, endTime)

		durationSecs := 0
		for _, e := range entries {
			durationSecs += e.Duration
		}

		fmt.Println("user: ", u.FirstName)
		fmt.Println("seconds: ", durationSecs)
		fmt.Println("hours: ", float32(durationSecs)/60/60)
	}
}
