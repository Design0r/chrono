package service

import (
	"context"
	"log/slog"
	"time"

	"chrono/internal/domain"
)

type TimestampsService struct {
	timestamps domain.TimestampsRepository
	event      *EventService
	log        *slog.Logger
}

func NewTimestampsService(
	r domain.TimestampsRepository,
	e *EventService,
	log *slog.Logger,
) *TimestampsService {
	return &TimestampsService{timestamps: r, log: log, event: e}
}

func (r *TimestampsService) GetById(ctx context.Context, id int64) (domain.Timestamp, error) {
	return r.timestamps.GetById(ctx, id)
}

func (r *TimestampsService) Start(ctx context.Context, userId int64) (domain.Timestamp, error) {
	return r.timestamps.Start(ctx, userId)
}

func (r *TimestampsService) Stop(ctx context.Context, id int64) (domain.Timestamp, error) {
	t, err := r.GetById(ctx, id)
	if err != nil {
		return domain.Timestamp{}, err
	}

	if t.EndTime != nil {
		return t, nil
	}

	return r.timestamps.Stop(ctx, id)
}

func (r *TimestampsService) Delete(ctx context.Context, id int64) error {
	return r.timestamps.Delete(ctx, id)
}

func (r *TimestampsService) GetInRange(
	ctx context.Context,
	userId int64,
	start time.Time,
	stop time.Time,
) ([]domain.Timestamp, error) {
	return r.timestamps.GetInRange(ctx, userId, start, stop)
}

func (r *TimestampsService) GetAllInRange(
	ctx context.Context,
	start time.Time,
	stop time.Time,
) ([]domain.Timestamp, error) {
	return r.timestamps.GetAllInRange(ctx, start, stop)
}

func (r *TimestampsService) GetForToday(
	ctx context.Context,
	userId int64,
) ([]domain.Timestamp, error) {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	stop := start.AddDate(0, 0, 1)

	return r.timestamps.GetInRange(ctx, userId, start, stop)
}

func (r *TimestampsService) GetForYear(
	ctx context.Context,
	userId int64,
	year int,
) ([]domain.Timestamp, error) {
	start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	stop := start.AddDate(1, 0, 0)

	return r.timestamps.GetInRange(ctx, userId, start, stop)
}

func (r *TimestampsService) GetForMonth(
	ctx context.Context,
	userId int64,
	year int,
	month int,
) ([]domain.Timestamp, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	stop := start.AddDate(0, 1, 0)

	return r.timestamps.GetInRange(ctx, userId, start, stop)
}

func (r *TimestampsService) GetTotalSecondsInRange(
	ctx context.Context,
	userId int64,
	start time.Time,
	stop time.Time,
) (float64, error) {
	return r.timestamps.GetTotalSecondsInRange(ctx, userId, start, stop)
}

func (r *TimestampsService) GetLatest(
	ctx context.Context, userId int64,
) (domain.Timestamp, error) {
	return r.timestamps.GetLatest(ctx, userId)
}

func (r *TimestampsService) Update(
	ctx context.Context,
	ts *domain.Timestamp,
) (domain.Timestamp, error) {
	return r.timestamps.Update(ctx, ts)
}

func (r *TimestampsService) GetAllForUser(
	ctx context.Context,
	userId int64,
) ([]domain.Timestamp, error) {
	return r.timestamps.GetAllForUser(ctx, userId)
}

func (r *TimestampsService) GetWorkHoursForYearForAllUsers(
	ctx context.Context,
	users []domain.User,
	year int,
	workDayHours float64,
) map[int64]domain.WorkHours {
	workHours := map[int64]domain.WorkHours{}

	for _, user := range users {
		result, err := r.GetWorkHoursForYear(ctx, user.ID, year, workDayHours)
		if err != nil {
			continue
		}
		workHours[user.ID] = result
	}

	return workHours
}

func (r *TimestampsService) GetWorkHoursForYear(
	ctx context.Context,
	userId int64,
	year int,
	workDayHours float64,
) (domain.WorkHours, error) {
	now := time.Now()
	loc := now.Location()

	// Start: Jan 1 of the given year
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, loc)

	var periodEnd time.Time

	switch {
	case year < now.Year():
		// Past year: full year until Dec 31 23:59:59
		periodEnd = time.Date(year, time.December, 31, 23, 59, 59, 0, loc)

	case year == now.Year():
		// Current year: until "yesterday evening"
		yesterday := now.AddDate(0, 0, -1)
		// If it's Jan 1, yesterday is previous year → no days in this year yet
		if yesterday.Before(yearStart) {
			return domain.WorkHours{}, nil
		}
		periodEnd = time.Date(
			yesterday.Year(),
			yesterday.Month(),
			yesterday.Day(),
			23,
			59,
			59,
			0,
			loc,
		)

	default: // year > now.Year()
		// Future year: nothing to count yet
		return domain.WorkHours{}, nil
	}

	// If somehow end is before start, just bail out
	if periodEnd.Before(yearStart) {
		return domain.WorkHours{}, nil
	}
	// ---- HOLIDAYS / VACATION / SICKNESS ----
	// NOTE: these must use the same period (yearStart..periodEnd) to be fully consistent.
	holidays, err := r.event.GetNonWeekendCountHolidays(ctx, yearStart, periodEnd)
	if err != nil {
		return domain.WorkHours{}, err
	}

	vacation, err := r.event.GetUsedVacation(ctx, userId, yearStart, periodEnd)
	if err != nil {
		return domain.WorkHours{}, err
	}

	sickDays := 0
	allEvents, err := r.event.GetAllByUserId(ctx, userId)
	if err != nil {
		return domain.WorkHours{}, err
	}

	for _, e := range allEvents {
		if e.Name != "krank" {
			continue
		}
		// only count sick days in [yearStart, periodEnd]
		if !e.ScheduledAt.Before(yearStart) && !e.ScheduledAt.After(periodEnd) {
			sickDays++
		}
	}

	// ---- WORKED HOURS ----

	worked, err := r.timestamps.GetTotalSecondsInRange(ctx, userId, yearStart, periodEnd)
	if err != nil {
		return domain.WorkHours{}, err
	}

	workedHours := worked / 60 / 60

	// ---- EXPECTED HOURS (weekdays between yearStart and periodEnd) ----

	expectedDays := 0
	for d := yearStart; !d.After(periodEnd); d = d.AddDate(0, 0, 1) {
		wd := d.Weekday()
		if wd != time.Saturday && wd != time.Sunday {
			expectedDays++
		}
	}

	expectedHours := (float64(expectedDays-holidays) - vacation - float64(sickDays)) * workDayHours
	holidayHours := float64(holidays) * workDayHours
	vacationHours := vacation * workDayHours

	return domain.WorkHours{
		Worked:   workedHours,
		Expected: expectedHours,
		Holidays: holidayHours,
		Vacation: vacationHours,
	}, nil
}
