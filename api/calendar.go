package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/schemas"
	"calendar/service"
)

func InitCalendarRoutes(group *echo.Group, db *sql.DB) {
	group.GET(
		"/:year/:month",
		func(c echo.Context) error { return MonthCalendarHandler(c, db) },
	)
}

func MonthCalendarHandler(c echo.Context, db *sql.DB) error {
	var date schemas.YMDate
	if err := c.Bind(&date); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter")
	}

	month := service.GetDaysOfMonth(time.Month(date.Month), date.Year)
	d := time.Date(date.Year, time.Month(date.Month), 0, 0, 0, 0, 0, time.Now().Local().Location())
	event, err := service.GetEventsForMonth(db, d)
	if err != nil {
		return err
	}
	fmt.Println(month.Offset)

	templates.Index(month, event).
		Render(context.Background(), c.Response().Writer)
	return nil
}
