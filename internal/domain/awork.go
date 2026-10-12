package domain

import (
	"regexp"
	"time"
)

var aworkIDPattern = regexp.MustCompile(
	`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`,
)

// IsValidAworkID reports whether id is a GUID as used by awork for user ids.
func IsValidAworkID(id string) bool {
	return aworkIDPattern.MatchString(id)
}

type AworkUser struct {
	Id        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type TimeBooking struct {
	Id        string    `json:"id"`
	Duration  float32   `json:"duration"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type Project struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TimeEntry struct {
	Id             string `json:"id"`
	Duration       int    `json:"duration"`
	StartDateLocal string `json:"startDateLocal"`
	EndDateLocal   string `json:"endDateLocal"`
	Task           Task
	Project        Project
}

type TimeBookingResponse struct {
	UserId       string        `json:"userId"`
	TimeBookings []TimeBooking `json:"timeBookings"`
}

type WorkHours struct {
	Worked   float64 `json:"worked"`
	Expected float64 `json:"expected"`
	Holidays float64 `json:"holidays"`
	Vacation float64 `json:"vacation"`
}
