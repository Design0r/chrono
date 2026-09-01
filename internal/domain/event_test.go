package domain_test

import (
	"testing"

	"chrono/internal/domain"
)

// TestVacationDays pins the weighting of each event name. Half days used to be
// charged and refunded as whole days in places that only asked IsVacation.
func TestVacationDays(t *testing.T) {
	tests := []struct {
		name     string
		vacation bool
		days     float64
	}{
		{name: "urlaub", vacation: true, days: 1.0},
		{name: "urlaub halbtags", vacation: true, days: 0.5},
		{name: "krank", vacation: false, days: 0},
		{name: "krank halbtags", vacation: false, days: 0},
		{name: "Neujahrstag", vacation: false, days: 0},
		{name: "", vacation: false, days: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := domain.Event{Name: tc.name}

			if got := e.IsVacation(); got != tc.vacation {
				t.Errorf("IsVacation() = %v, want %v", got, tc.vacation)
			}

			if got := e.VacationDays(); got != tc.days {
				t.Errorf("VacationDays() = %v, want %v", got, tc.days)
			}
		})
	}
}
