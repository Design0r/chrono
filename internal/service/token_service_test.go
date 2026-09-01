package service

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"chrono/internal/domain"
)

type fakeRefreshRepo struct {
	years map[int]bool
}

func (f *fakeRefreshRepo) Create(_ context.Context, year int, userId int64) (*domain.RefreshToken, error) {
	f.years[year] = true
	return &domain.RefreshToken{Year: int64(year), UserID: userId}, nil
}

func (f *fakeRefreshRepo) DeleteAll(_ context.Context) error { return nil }

func (f *fakeRefreshRepo) ExistsForUser(_ context.Context, _ int64, year int) (bool, error) {
	return f.years[year], nil
}

type fakeVacationRepo struct {
	booked []float64
}

func (f *fakeVacationRepo) Create(
	_ context.Context,
	t domain.CreateVacationToken,
) (*domain.VacationToken, error) {
	f.booked = append(f.booked, t.Value)
	return &domain.VacationToken{Value: t.Value, UserID: t.UserID}, nil
}

func (f *fakeVacationRepo) Delete(_ context.Context, _ int64) error { return nil }
func (f *fakeVacationRepo) DeleteAll(_ context.Context) error       { return nil }

func (f *fakeVacationRepo) GetRemainingVacationForUser(
	_ context.Context,
	_ int64,
	_ time.Time,
	_ time.Time,
) (float64, error) {
	return f.total(), nil
}

func (f *fakeVacationRepo) total() float64 {
	sum := 0.0
	for _, v := range f.booked {
		sum += v
	}
	return sum
}

func newTokenService() (*TokenService, *fakeVacationRepo) {
	vac := &fakeVacationRepo{}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	return NewTokenService(&fakeRefreshRepo{years: map[int]bool{}}, vac, log), vac
}

// An allowance change before the user's yearly grant has been handed out must
// still leave them with the full new allowance. The delta used to be booked on
// top of nothing while marking the year as done, so the grant was lost for good.
func TestUpdateYearlyTokensGrantsTheBaseFirst(t *testing.T) {
	svc, vac := newTokenService()
	user := &domain.User{ID: 1, VacationDays: 30}

	// admin raises Bob 30 -> 32 in January, before Bob has opened the app
	if err := svc.UpdateYearlyTokens(context.Background(), user, 2, 2026); err != nil {
		t.Fatalf("UpdateYearlyTokens: %v", err)
	}

	if got := vac.total(); got != 32 {
		t.Errorf("balance after raising 30 -> 32 = %v, want 32", got)
	}

	// Bob opens the app later that year: the year is done, nothing more is granted
	user.VacationDays = 32
	if err := svc.InitYearlyTokens(context.Background(), user, 2026); err != nil {
		t.Fatalf("InitYearlyTokens: %v", err)
	}

	if got := vac.total(); got != 32 {
		t.Errorf("balance after the user's first visit = %v, want 32", got)
	}
}

// The same edit once the year has been granted must book only the difference.
func TestUpdateYearlyTokensAfterTheGrant(t *testing.T) {
	svc, vac := newTokenService()
	user := &domain.User{ID: 1, VacationDays: 30}
	ctx := context.Background()

	if err := svc.InitYearlyTokens(ctx, user, 2026); err != nil {
		t.Fatalf("InitYearlyTokens: %v", err)
	}

	if err := svc.UpdateYearlyTokens(ctx, user, 2, 2026); err != nil {
		t.Fatalf("UpdateYearlyTokens: %v", err)
	}

	if got := vac.total(); got != 32 {
		t.Errorf("balance = %v, want 32", got)
	}
}

// A profile edit that leaves the allowance alone must not mark the year as
// granted, or the user never receives it.
func TestUpdateYearlyTokensWithNoChange(t *testing.T) {
	svc, vac := newTokenService()
	user := &domain.User{ID: 1, VacationDays: 30}
	ctx := context.Background()

	if err := svc.UpdateYearlyTokens(ctx, user, 0, 2026); err != nil {
		t.Fatalf("UpdateYearlyTokens: %v", err)
	}

	if got := vac.total(); got != 30 {
		t.Errorf("balance after a no-op edit = %v, want 30", got)
	}

	if err := svc.InitYearlyTokens(ctx, user, 2026); err != nil {
		t.Fatalf("InitYearlyTokens: %v", err)
	}

	if got := vac.total(); got != 30 {
		t.Errorf("balance after the user's first visit = %v, want 30", got)
	}
}

// A user with no allowance still has the year marked as handled, so that a
// later increase is booked as a delta rather than granted twice.
func TestInitYearlyTokensWithoutAllowance(t *testing.T) {
	svc, vac := newTokenService()
	user := &domain.User{ID: 1, VacationDays: 0}
	ctx := context.Background()

	if err := svc.InitYearlyTokens(ctx, user, 2026); err != nil {
		t.Fatalf("InitYearlyTokens: %v", err)
	}

	if len(vac.booked) != 0 {
		t.Errorf("booked %v for a user with no allowance, want nothing", vac.booked)
	}

	if err := svc.UpdateYearlyTokens(ctx, user, 10, 2026); err != nil {
		t.Fatalf("UpdateYearlyTokens: %v", err)
	}

	if got := vac.total(); got != 10 {
		t.Errorf("balance after being given 10 days = %v, want 10", got)
	}
}
