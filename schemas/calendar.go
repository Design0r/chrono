package schemas

import (
	"time"

	"calendar/db/repo"
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
	Events []repo.Event
	Date   time.Time
}
