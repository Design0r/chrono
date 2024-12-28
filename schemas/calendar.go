package schemas

import (
	"time"

	"chrono/db/repo"
)

type Month struct {
	Name   string
	Number int
	Year   int
	Days   []Day
	Offset int
}

type Day struct {
	Number int
	Name   string
	Events []Event
	Date   time.Time
}

type Event struct {
	Username string
	repo.Event
}

type YearProgress struct {
	NumDays           int
	NumWorkDays       int
	NumDaysPassed     int
	DaysPassedPercent float32
	NumHolidays       int
	NumWastedHolidays int
}
