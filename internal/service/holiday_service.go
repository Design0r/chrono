package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"chrono/config"
	"chrono/internal/domain"
)

type HolidayService struct {
	user  *UserService
	event *EventService
	api   domain.ApiCacheRepository
	log   *slog.Logger
}

func NewHolidayService(
	user *UserService,
	event *EventService,
	api domain.ApiCacheRepository,
	log *slog.Logger,
) *HolidayService {
	return &HolidayService{user: user, event: event, api: api, log: log}
}

func (svc *HolidayService) Update(ctx context.Context, year int) error {
	if year < 1900 {
		return fmt.Errorf("Invalid year %v, must be 1900 and above", year)
	}

	cfg := config.GetConfig()
	if svc.HolidayCacheExists(ctx, year) {
		return nil
	}
	bot, err := svc.user.GetByName(ctx, cfg.BotName)
	if err != nil {
		return err
	}

	holidays, err := svc.FetchHolidays(year)
	if err != nil {
		return err
	}

	holidays = svc.filterHolidays(holidays)

	for name, data := range holidays {
		date, err := time.Parse(time.DateOnly, data["datum"])
		if err != nil {
			svc.log.Error("Failed parsing date", slog.String("error", err.Error()))
			continue
		}
		svc.event.Create(
			ctx,
			domain.YMDDate{Year: date.Year(), Month: int(date.Month()), Day: date.Day()},
			name,
			bot,
		)
	}

	return svc.CreateCache(ctx, year)
}

func (svc *HolidayService) FetchHolidays(year int) (domain.Holidays, error) {
	holidays := domain.Holidays{}
	resp, err := http.Get(fmt.Sprintf("https://feiertage-api.de/api/?jahr=%v&nur_land=BW", year))
	if err != nil {
		svc.log.Error(
			"Error fetching feiertage-api",
			slog.Int("year", year),
			slog.String("error", err.Error()),
		)
		return holidays, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return holidays, err
	}
	err = json.Unmarshal(body, &holidays)
	if err != nil {
		return holidays, err
	}

	return holidays, nil
}

func (svc *HolidayService) filterHolidays(holidays domain.Holidays) domain.Holidays {
	filter := map[string]bool{"Reformationstag": true}

	for holiday := range holidays {
		if _, exists := filter[holiday]; exists {
			delete(holidays, holiday)
		}
	}

	return holidays
}

func (svc *HolidayService) HolidayCacheExists(ctx context.Context, year int) bool {
	count, err := svc.api.Exists(ctx, int64(year))
	if err != nil {
		return false
	}
	return count > 0
}

func (svc *HolidayService) CreateCache(ctx context.Context, year int) error {
	return svc.api.Create(ctx, int64(year))
}

func (svc *HolidayService) GetAPICacheYears(ctx context.Context) ([]int64, error) {
	return svc.api.GetAll(ctx)
}
