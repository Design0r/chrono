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
		"%vT00:00",
		end.Format(time.DateOnly),
	)

	u, _ := url.Parse(fmt.Sprintf("%v/%v", AWORK_API_URL, "timeentries"))

	q := u.Query()
	q.Set("pageSize", "1000")
	q.Set(
		"filterby",
		fmt.Sprintf(
			"userId eq guid'%v' and startDateLocal ge datetime'%v' and endDateLocal le datetime'%v'",
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

func (a *AworkService) GetWorkHoursForWeek(
	userId string,
	isoWeek int,
	year int,
) (domain.WorkHours, error) {
	weekStart := domain.FirstDayOfISOWeek(year, isoWeek, time.Now().Location())
	weekEnd := weekStart.AddDate(0, 0, 7)

	entries, err := a.GetTimeEntries(userId, weekStart, weekEnd)
	if err != nil {
		return domain.WorkHours{}, err
	}

	workSecs := 0
	for _, entry := range entries {
		workSecs += entry.Duration
	}

	workHours := float64(workSecs) / 60

	return domain.WorkHours{Worked: workHours, Expected: 40, Vacation: 0}, nil
}

func (a *AworkService) GetWorkHoursForMonth(
	userId string,
	month int,
	year int,
) (domain.WorkHours, error) {
	monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	entries, err := a.GetTimeEntries(userId, monthStart, monthEnd)
	if err != nil {
		return domain.WorkHours{}, err
	}

	workSecs := 0
	for _, entry := range entries {
		workSecs += entry.Duration
	}

	workHours := float64(workSecs) / 60

	expected := 0
	start := monthStart
	for range domain.GetNumDaysOfMonth(time.Month(month), year) {
		weekday := start.Weekday()
		if weekday != time.Sunday && weekday != time.Saturday {
			expected += 1
		}
		start = start.AddDate(0, 0, 1)
	}

	return domain.WorkHours{Worked: workHours, Expected: float64(expected), Vacation: 0}, nil
}

func (a *AworkService) GetWorkHoursForYear(
	aworkUserId string,
	userId int64,
	year int,
) (domain.WorkHours, error) {
	ctx := context.Background()
	yearStart := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Now().Location())
	yearEnd := time.Now()

	entries, err := a.GetTimeEntries(aworkUserId, yearStart, yearEnd)
	if err != nil {
		return domain.WorkHours{}, err
	}

	holidays, err := a.event.GetNonWeekendCountHolidaysForYear(ctx, year)
	if err != nil {
		return domain.WorkHours{}, err
	}

	vacation, err := a.event.GetUsedVacationTilNow(ctx, userId, year)
	if err != nil {
		return domain.WorkHours{}, err
	}

	workSecs := 0
	for _, entry := range entries {
		workSecs += entry.Duration
	}

	workHours := float64(workSecs) / 60 / 60

	expected := 0
	start := yearStart
	for range yearEnd.YearDay() {
		weekday := start.Weekday()
		if weekday != time.Sunday && weekday != time.Saturday {
			expected += 1
		}
		start = start.AddDate(0, 0, 1)
	}

	return domain.WorkHours{
		Worked:   workHours,
		Expected: (float64(expected-holidays) - vacation) * 8,
		Holidays: float64(holidays) * 8,
		Vacation: vacation * 8,
	}, nil
}
