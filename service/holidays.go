package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"calendar/db/repo"
	"calendar/schemas"
)

type Holidays = map[string]map[string]string

func UpdateHolidays(db *sql.DB, year int) error {
	if HolidayCacheExists(db, year) {
		return nil
	}

	holidays, err := FetchHolidays(db, year)
	if err != nil {
		return err
	}

	for name, data := range holidays {
		date, err := time.Parse("2006-01-02", data["datum"])
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = CreateEvent(
			db,
			schemas.YMDDate{Year: date.Year(), Month: int(date.Month()), Day: date.Day()},
			int64(1),
			name,
		)
	}

	CreateCache(db, year)

	return nil
}

func FetchHolidays(db *sql.DB, year int) (Holidays, error) {
	fmt.Println("FETCHINGGGG")
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

func HolidayCacheExists(db *sql.DB, year int) bool {
	r := repo.New(db)

	count, err := r.CacheExists(context.Background(), int64(year))
	if err != nil {
		return false
	}
	return count > 0
}

func CreateCache(db *sql.DB, year int) error {
	r := repo.New(db)

	err := r.CreateCache(context.Background(), int64(year))
	if err != nil {
		return err
	}

	return nil
}
