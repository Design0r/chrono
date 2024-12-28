package service

import (
	"time"

	"chrono/schemas"
)

var monthDays map[time.Month]int = map[time.Month]int{
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

var weekdays map[time.Weekday]string = map[time.Weekday]string{
	time.Monday:    "Monday",
	time.Tuesday:   "Tuesday",
	time.Wednesday: "Wednesday",
	time.Thursday:  "Thursday",
	time.Friday:    "Friday",
	time.Saturday:  "Saturday",
	time.Sunday:    "Sunday",
}

func GetNumDaysOfMonth(month time.Month, year int) int {
	if month == time.February {
		if IsLeapYear(year) {
			return 29
		}
	}

	return monthDays[month]
}

func GetDaysOfMonth(month time.Month, year int) schemas.Month {
	numDays := GetNumDaysOfMonth(month, year)
	days := make([]schemas.Day, numDays)
	for i := 0; i < numDays; i++ {
		date := time.Date(year, month, i+1, 0, 0, 0, 0, time.Now().Local().Location())
		day := schemas.Day{Number: i + 1, Name: weekdays[date.Weekday()], Date: date}
		days[i] = day
	}

	return schemas.Month{
		Name:   days[0].Date.Month().String(),
		Days:   days,
		Offset: getMonthOffset(days[0].Date.Weekday() + 6), // Weekday 0 = Sunday
		Year:   year,
		Number: int(month),
	}
}

func getMonthOffset(weekday time.Weekday) int {
	return int(weekday) % 7
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

func YearProgress(year int) int {
	now := time.Now().YearDay()

	return now
}

func NumWorkDays(year int) int {
	counter := 0
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.Now().Location())
	for i := 0; i < NumDaysInYear(year); i++ {
		if start.Weekday() == time.Saturday || start.Weekday() == time.Sunday {
			counter++
		}
	}

	return counter
}
