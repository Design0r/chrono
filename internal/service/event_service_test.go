package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"

	"chrono/internal/domain"
)

type fakeEventRepo struct {
	events map[int64]domain.Event
	nextID int64
}

func newFakeEventRepo() *fakeEventRepo {
	return &fakeEventRepo{events: map[int64]domain.Event{}}
}

func (f *fakeEventRepo) Create(
	_ context.Context,
	data domain.YMDDate,
	eventType string,
	user *domain.User,
) (*domain.Event, error) {
	evt := domain.Event{Name: eventType}

	state := "pending"
	if !evt.IsVacation() || user.IsSuperuser {
		state = "accepted"
	}

	f.nextID++
	e := domain.Event{
		ID:     f.nextID,
		Name:   eventType,
		State:  state,
		UserID: user.ID,
		ScheduledAt: time.Date(
			data.Year, time.Month(data.Month), data.Day, 0, 0, 0, 0, time.UTC),
	}
	f.events[e.ID] = e

	return &e, nil
}

func (f *fakeEventRepo) seed(e domain.Event) domain.Event {
	f.nextID++
	e.ID = f.nextID
	f.events[e.ID] = e

	return e
}

func (f *fakeEventRepo) GetById(_ context.Context, id int64) (*domain.Event, error) {
	e, ok := f.events[id]
	if !ok {
		return nil, fmt.Errorf("no event %d", id)
	}
	return &e, nil
}

func (f *fakeEventRepo) Delete(_ context.Context, id int64) error {
	delete(f.events, id)
	return nil
}

func (f *fakeEventRepo) Update(_ context.Context, _ int64, _ string) (*domain.Event, error) {
	return nil, nil
}
func (f *fakeEventRepo) GetForDay(_ context.Context, _ domain.YMDDate) ([]domain.Event, error) {
	return nil, nil
}
func (f *fakeEventRepo) GetForMonth(
	_ context.Context, _ domain.YMDate, _ string, _ *domain.User, _ string,
) (domain.Month, error) {
	return domain.Month{}, nil
}
func (f *fakeEventRepo) GetForYear(_ context.Context, _ int) ([]domain.EventUser, error) {
	return nil, nil
}
func (f *fakeEventRepo) GetPendingForUser(_ context.Context, _ int64, _ int) (int, error) {
	return 0, nil
}
func (f *fakeEventRepo) GetUsedVacationForUser(
	_ context.Context, _ int64, _ int,
) (float64, error) {
	return 0, nil
}
func (f *fakeEventRepo) UpdateInRange(
	_ context.Context, _ int64, _ string, _, _ time.Time,
) error {
	return nil
}
func (f *fakeEventRepo) GetAllByUserId(_ context.Context, _ int64) ([]domain.Event, error) {
	return nil, nil
}

func newEventService() (*EventService, *fakeEventRepo, *fakeVacationRepo) {
	tokens, vac := newTokenService()
	events := newFakeEventRepo()
	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	return NewEventService(events, nil, nil, tokens, log), events, vac
}

// A superuser's vacation is charged straight away. A half day must cost half a
// day, not a whole one, or it disagrees with the "used" figure counted from the
// events.
func TestSuperuserVacationIsChargedByTheDay(t *testing.T) {
	tests := []struct {
		name     string
		expected float64
	}{
		{"urlaub", -1.0},
		{"urlaub halbtags", -0.5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, _, vac := newEventService()
			admin := &domain.User{ID: 1, IsSuperuser: true}

			_, err := svc.Create(
				context.Background(),
				domain.YMDDate{Year: 2026, Month: 6, Day: 1},
				tc.name,
				admin,
			)
			if err != nil {
				t.Fatalf("Create: %v", err)
			}

			if got := vac.total(); got != tc.expected {
				t.Errorf("charged %v for %q, want %v", got, tc.name, tc.expected)
			}
		})
	}
}

// Deleting vacation has to give back exactly what it cost. Approving a half day
// costs half a day, so refunding a whole one handed out free vacation on every
// book-approve-delete cycle.
func TestDeleteRefundsWhatTheApprovalCharged(t *testing.T) {
	tests := []struct {
		name string
		// what the approval flow books for a single day of this kind
		charged float64
	}{
		{"urlaub", -1.0},
		{"urlaub halbtags", -0.5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, events, vac := newEventService()
			user := &domain.User{ID: 2}
			ctx := context.Background()

			event := events.seed(domain.Event{
				Name:        tc.name,
				State:       "accepted",
				UserID:      user.ID,
				ScheduledAt: time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
			})

			// the request handler charges the allowance when the day is approved
			if _, err := svc.token.CreateVacationToken(ctx, tc.charged, 2026, user.ID); err != nil {
				t.Fatalf("charging the approval: %v", err)
			}

			if _, err := svc.Delete(ctx, event.ID, user); err != nil {
				t.Fatalf("Delete: %v", err)
			}

			if got := vac.total(); got != 0 {
				t.Errorf("charged %v then deleted, balance = %v, want 0", tc.charged, got)
			}
		})
	}
}

// Nothing was charged for vacation that was never accepted, so nothing may be
// given back for it either.
func TestDeleteOfPendingVacationRefundsNothing(t *testing.T) {
	svc, events, vac := newEventService()
	user := &domain.User{ID: 2}
	ctx := context.Background()

	event, err := events.Create(
		ctx, domain.YMDDate{Year: 2026, Month: 6, Day: 1}, "urlaub halbtags", user)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if event.State != "pending" {
		t.Fatalf("expected a pending event, got %q", event.State)
	}

	if _, err := svc.Delete(ctx, event.ID, user); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if got := vac.total(); got != 0 {
		t.Errorf("refunded %v for a pending event, want 0", got)
	}
}
