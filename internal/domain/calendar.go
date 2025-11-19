package domain

import (
	"time"
)

var (
	monthDays map[time.Month]int = map[time.Month]int{
		time.January:   31,
		time.February:  28,
		time.March:     31,
		time.April:     30,
		time.May:       31,
		time.June:      30,
		time.July:      31,
		time.August:    31,
		time.September: 30,
		time.October:   31,
		time.November:  30,
		time.December:  31,
	}

	MonthDaysList [12]int = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	weekdays map[time.Weekday]string = map[time.Weekday]string{
		time.Monday:    "Monday",
		time.Tuesday:   "Tuesday",
		time.Wednesday: "Wednesday",
		time.Thursday:  "Thursday",
		time.Friday:    "Friday",
		time.Saturday:  "Saturday",
		time.Sunday:    "Sunday",
	}
)

type Month struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
	Year   int    `json:"year"`
	Days   []Day  `json:"days"`
	Offset int    `json:"offset"`
}

type Day struct {
	Number int         `json:"number"`
	Name   string      `json:"name"`
	Events []EventUser `json:"events"`
	Date   time.Time   `json:"date"`
}

type YearProgress struct {
	NumDays           int
	NumWorkDays       int
	NumDaysPassed     int
	DaysPassedPercent float32
	NumHolidays       int
	NumWastedHolidays int
}

type YearHistogram struct {
	IsHoliday      bool     `json:"is_holiday"`
	Count          int      `json:"count"`
	LastDayOfMonth bool     `json:"last_day_of_month"`
	IsCurrentWeek  bool     `json:"is_current_week"`
	Usernames      []string `json:"usernames"`
	Date           string   `json:"date"`
}

func GetNumDaysOfMonth(month time.Month, year int) int {
	if month == time.February {
		if IsLeapYear(year) {
			return 29
		}
	}

	return monthDays[month]
}

func GetDaysOfMonth(month time.Month, year int) Month {
	numDays := GetNumDaysOfMonth(month, year)
	days := make([]Day, numDays)
	for i := range numDays {
		date := time.Date(year, month, i+1, 0, 0, 0, 0, time.Local)
		day := Day{Number: i + 1, Name: weekdays[date.Weekday()], Date: date}
		days[i] = day
	}

	return Month{
		Name:   days[0].Date.Month().String(),
		Days:   days,
		Offset: getMonthOffset(days[0].Date.Weekday()),
		Year:   year,
		Number: int(month),
	}
}

func getMonthOffset(weekday time.Weekday) int {
	return (int(weekday) + 6) % 7 // Weekday 0 = Sunday
}

func GetYearOffset(year int) int {
	firstDay := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	offset := int(firstDay.Weekday()) - 1
	if offset < 0 {
		offset = 6
	}

	return offset
}

func GetMonthGaps(year int) []int {
	list := make([]int, 12)

	for i := range 12 {
		month := time.Month(i + 1)
		firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		offset := int(firstDay.Weekday()+6) - 1 // -1 because sunday is considered the first day
		numDays := GetNumDaysOfMonth(month, year)

		cols := int((offset + numDays) / 7)

		if offset > 0 && i > 0 {
			cols--
		}

		list[i] = max(cols-1, 0)

	}

	return list
}

func IsLeapYear(year int) bool {
	t := time.Date(year, time.February, 29, 0, 0, 0, 0, time.UTC)
	return t.Day() == 29
}

func CurrentYear() int {
	return time.Now().Year()
}

func NumDaysInYear(year int) int {
	if IsLeapYear(year) {
		return 366
	}
	return 365
}

func YearProgressPercent(year int) float32 {
	now := time.Now().YearDay()
	days := NumDaysInYear(year)

	return (float32(now) / float32(days)) * float32(100)
}

func CurrentYearDay(year int) int {
	now := time.Now().YearDay()

	return now
}

func NumWorkDays(year int) int {
	counter := 0
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())
	for range NumDaysInYear(year) {
		if start.Weekday() == time.Saturday || start.Weekday() == time.Sunday {
			counter++
		}
	}

	return counter
}

func GetCurrentYearProgress() YearProgress {
	currYear := CurrentYear()
	return YearProgress{
		NumDays:           NumDaysInYear(currYear),
		NumWorkDays:       NumWorkDays(currYear),
		NumDaysPassed:     CurrentYearDay(currYear),
		DaysPassedPercent: YearProgressPercent(currYear),
	}
}

func GetStrWeekday(day time.Weekday) string {
	return weekdays[day]
}

func FirstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	return date
}
