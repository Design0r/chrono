package domain_test

import (
	"chrono/internal/domain"
	"fmt"
	"testing"
	"time"
)

// TestIsLeapYear checks the IsLeapYear function for multiple cases.
func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		year     int
		expected bool
	}{
		{year: 1900, expected: false}, // century not divisible by 400
		{year: 2000, expected: true},  // divisible by 400
		{year: 2020, expected: true},  // typical leap year
		{year: 2021, expected: false}, // typical non-leap year
		{year: 2400, expected: true},  // future, divisible by 400
		{year: 2100, expected: false},
	}

	for _, tc := range tests {
		t.Run(
			funcName(tc.year),
			func(t *testing.T) {
				got := domain.IsLeapYear(tc.year)
				if got != tc.expected {
					t.Errorf("IsLeapYear(%d) = %v, want %v", tc.year, got, tc.expected)
				}
			},
		)
	}
}

// TestGetNumDaysOfMonth checks the correct number of days for months,
// including leap-year handling for February.
func TestGetNumDaysOfMonth(t *testing.T) {
	tests := []struct {
		month    time.Month
		year     int
		expected int
	}{
		{time.January, 2021, 31},
		{time.February, 2021, 28},
		{time.February, 2020, 29}, // leap year
		{time.April, 2021, 30},
		{time.December, 2021, 31},
	}

	for _, tc := range tests {
		t.Run(
			funcName(tc.month.String(), tc.year),
			func(t *testing.T) {
				got := domain.GetNumDaysOfMonth(tc.month, tc.year)
				if got != tc.expected {
					t.Errorf("GetNumDaysOfMonth(%s, %d) = %d, want %d",
						tc.month, tc.year, got, tc.expected)
				}
			},
		)
	}
}

// TestNumDaysInYear checks total days for a normal year vs a leap year.
func TestNumDaysInYear(t *testing.T) {
	tests := []struct {
		year     int
		expected int
	}{
		{2021, 365},
		{2020, 366},
	}

	for _, tc := range tests {
		t.Run(
			funcName(tc.year),
			func(t *testing.T) {
				got := domain.NumDaysInYear(tc.year)
				if got != tc.expected {
					t.Errorf("NumDaysInYear(%d) = %d, want %d", tc.year, got, tc.expected)
				}
			},
		)
	}
}

// TestGetDaysOfMonth checks if we correctly build the Month struct.
func TestGetDaysOfMonth(t *testing.T) {
	monthData := domain.GetDaysOfMonth(time.March, 2021)
	if monthData.Name != "March" {
		t.Errorf("Expected month name to be March, got %s", monthData.Name)
	}
	if len(monthData.Days) != 31 {
		t.Errorf("Expected 31 days, got %d", len(monthData.Days))
	}

	// Check the first day is Monday (2021-03-01 was a Monday).
	if monthData.Days[0].Name != "Monday" {
		t.Errorf("Expected first day to be Monday, got %s", monthData.Days[0].Name)
	}
}

// TestYearProgressPercent is a simple sanity check. This test might be time-dependent.
// Typically you'd freeze time or mock time.Now() in advanced setups.
func TestYearProgressPercent(t *testing.T) {
	year := time.Now().Year()
	got := domain.YearProgressPercent(year)

	if got < 0 || got > 100 {
		t.Errorf("YearProgressPercent(%d) = %f, expected between [0..100]", year, got)
	}
}

// Helper to produce subtest names based on input.
func funcName(args ...interface{}) string {
	var s string
	for _, a := range args {
		s += "_" + time.Now().Format("20060102") + "_" + // for uniqueness if needed
			stringify(a)
	}
	return s
}

func stringify(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
