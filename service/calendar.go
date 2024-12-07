package service

import "time"

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

func GetNumDaysOfMonth(month time.Month, year int) int {
	if month == time.February {
		if IsLeapYear(year) {
			return 29
		}
	}

	return monthDays[month]
}

func GetDaysOfMonth(month time.Month, year int) []time.Time {
	numDays := GetNumDaysOfMonth(month, year)
	days := make([]time.Time, numDays)
	for i := 0; i <= numDays; i++ {
		day := time.Date(year, month, 0, 0, 0, 0, 0, time.Now().Local().Location())
		days[i] = day
	}

	return days
}

func IsLeapYear(year int) bool {
	t := time.Date(year, time.February, 29, 0, 0, 0, 0, time.UTC)
	return t.Day() == 29
}

func CurrentYear() int {
	return time.Now().Year()
}
