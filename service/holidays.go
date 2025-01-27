package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"chrono/db/repo"
	"chrono/schemas"
)

type Holidays = map[string]map[string]string

func UpdateHolidays(r *repo.Queries, year int) error {
	if HolidayCacheExists(r, year) {
		return nil
	}
	bot, err := GetUserByName(r, "Chrono Bot")
	if err != nil {
		return err
	}

	holidays, err := FetchHolidays(r, year)
	if err != nil {
		return err
	}

	for name, data := range holidays {
		date, err := time.Parse(time.DateOnly, data["datum"])
		if err != nil {
			log.Println(err)
			continue
		}
		CreateEvent(
			r,
			schemas.YMDDate{Year: date.Year(), Month: int(date.Month()), Day: date.Day()},
			bot,
			name,
		)
	}

	return CreateCache(r, year)
}

func FetchHolidays(r *repo.Queries, year int) (Holidays, error) {
	holidays := Holidays{}
	resp, err := http.Get(fmt.Sprintf("https://feiertage-api.de/api/?jahr=%v&nur_land=BW", year))
	if err != nil {
		log.Printf("Error fetching feiertage: %v", err)
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

func HolidayCacheExists(r *repo.Queries, year int) bool {
	count, err := r.CacheExists(context.Background(), int64(year))
	if err != nil {
		return false
	}
	return count > 0
}

func CreateCache(r *repo.Queries, year int) error {
	err := r.CreateCache(context.Background(), int64(year))
	if err != nil {
		return err
	}

	return nil
}
