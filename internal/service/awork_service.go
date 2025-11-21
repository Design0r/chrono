package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"chrono/config"
	"chrono/internal/domain"
)

const AWORK_API_URL = "https://api.awork.com/api/v1"

type AworkService struct {
	client http.Client
	log    *slog.Logger
	event  EventService
	user   UserService
}

func NewAworkService(e EventService, u UserService, s *slog.Logger) AworkService {
	return AworkService{client: http.Client{}, event: e, user: u, log: s}
}

func (a *AworkService) GetUsers() ([]domain.AworkUser, error) {
	cfg := config.GetConfig()
	data := []domain.AworkUser{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%v/users", AWORK_API_URL), nil)
	if err != nil {
		a.log.Error("Failed to create request", slog.String("error", err.Error()))
		return data, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", cfg.AworkApiKey))

	res, err := a.client.Do(req)
	if err != nil {
		a.log.Error("Failed to send request", slog.String("error", err.Error()))
		return data, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		a.log.Error("Failed to read body", slog.String("error", err.Error()))
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		a.log.Error("Failed to unmarshal", slog.String("error", err.Error()))
		return data, err
	}

	return data, nil
}

func (a *AworkService) GetTimeEntries(
	userId string,
	start time.Time,
	end time.Time,
) ([]domain.TimeEntry, error) {
	data := []domain.TimeEntry{}

	cfg := config.GetConfig()
	startDate := fmt.Sprintf(
		"%vT00:00",
		start.Format(time.DateOnly),
	)
	endDate := fmt.Sprintf(
		"%vT23:59",
		end.Format(time.DateOnly),
	)

	u, _ := url.Parse(fmt.Sprintf("%v/%v", AWORK_API_URL, "timeentries"))

	q := u.Query()
	q.Set("pageSize", "1000")
	q.Set(
		"filterby",
		fmt.Sprintf(
			"userId eq guid'%v' and startDateLocal ge datetime'%v' and startDateLocal le datetime'%v'",
			userId,
			startDate,
			endDate,
		),
	)
	q.Set("orderby", "startDateLocal asc")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(
		"GET",
		u.String(),
		nil,
	)
	if err != nil {
		a.log.Error("Failed to create request", slog.String("error", err.Error()))
		return data, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", cfg.AworkApiKey))

	res, err := a.client.Do(req)
	if err != nil {
		a.log.Error("Failed to send request", slog.String("error", err.Error()))
		return data, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return data, errors.New(res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		a.log.Error("Failed to read body", slog.String("error", err.Error()))
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		a.log.Error("Failed to unmarshal", slog.String("error", err.Error()))
		return data, err
	}

	return data, nil
}

func (a *AworkService) ConvertAworkTime(aworkTime string) (time.Time, error) {
	split := strings.Split(aworkTime, "T")
	date := split[0]
	parsed, err := time.Parse(time.DateOnly, date)
	if err != nil {
		a.log.Error("Failed to parse date")
		return time.Time{}, err
	}

	return parsed, nil
}

func (a *AworkService) GetWorkHoursForYear(
	aworkUserId string,
	userId int64,
	year int,
) (domain.WorkHours, error) {
	ctx := context.Background()

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

	// ---- TIME ENTRIES (with basic pagination) ----

	entries, err := a.GetTimeEntries(aworkUserId, yearStart, periodEnd)
	if err != nil {
		return domain.WorkHours{}, err
	}

	if len(entries) == 0 {
		return domain.WorkHours{}, nil
	}

	entryIDs := map[string]struct{}{}
	for _, e := range entries {
		entryIDs[e.Id] = struct{}{}
	}

	// Simple pagination: keep fetching while we get a full page
	for len(entries) == 1000 {
		lastEntry := entries[len(entries)-1]

		newStartDate, err := a.ConvertAworkTime(lastEntry.EndDateLocal)
		if err != nil {
			return domain.WorkHours{}, err
		}

		page, err := a.GetTimeEntries(aworkUserId, newStartDate, periodEnd)
		if err != nil {
			return domain.WorkHours{}, err
		}
		if len(page) == 0 {
			break
		}

		for _, e := range page {
			if _, exists := entryIDs[e.Id]; exists {
				continue
			}
			entryIDs[e.Id] = struct{}{}
			entries = append(entries, e)
		}

		if len(page) < 1000 {
			break
		}
	}

	for i, e := range entries {
		fmt.Println(i, e.StartDateLocal, e.EndDateLocal, e.Task.Name)
	}

	// ---- HOLIDAYS / VACATION / SICKNESS ----
	// NOTE: these must use the same period (yearStart..periodEnd) to be fully consistent.
	holidays, err := a.event.GetNonWeekendCountHolidays(ctx, yearStart, periodEnd)
	if err != nil {
		return domain.WorkHours{}, err
	}

	vacation, err := a.event.GetUsedVacation(ctx, userId, yearStart, periodEnd)
	if err != nil {
		return domain.WorkHours{}, err
	}

	sickDays := 0
	allEvents, err := a.event.GetAllByUserId(ctx, userId)
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

	workSecs := 0
	for _, entry := range entries {
		// Skip vacation/holidays/comp-time entries
		if entry.Task.Name == "Urlaub" ||
			entry.Task.Name == "Feiertag" ||
			entry.Task.Name == "Ausgleich" {
			continue
		}
		workSecs += entry.Duration
	}

	workedHours := float64(workSecs) / 3600.0

	// ---- EXPECTED HOURS (weekdays between yearStart and periodEnd) ----

	expectedDays := 0
	for d := yearStart; !d.After(periodEnd); d = d.AddDate(0, 0, 1) {
		wd := d.Weekday()
		if wd != time.Saturday && wd != time.Sunday {
			expectedDays++
		}
	}

	// final numbers (assuming 8h per working day)
	expectedHours := (float64(expectedDays-holidays) - vacation - float64(sickDays)) * 8
	holidayHours := float64(holidays) * 8
	vacationHours := vacation * 8

	return domain.WorkHours{
		Worked:   workedHours,
		Expected: expectedHours,
		Holidays: holidayHours,
		Vacation: vacationHours,
	}, nil
}
