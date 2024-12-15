package views

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/middleware"
	"calendar/schemas"
	"calendar/service"
)

func InitCalendarRoutes(group *echo.Group, db *sql.DB) {
	group.GET(
		"/:year/:month",
		func(c echo.Context) error { return MonthCalendarHandler(c, db) },
		middleware.SessionMiddleware(db),
	)
}

func MonthCalendarHandler(c echo.Context, db *sql.DB) error {
	var date schemas.YMDate
	if err := c.Bind(&date); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter")
	}
	service.UpdateHolidays(db, date.Year)

	month := service.GetDaysOfMonth(time.Month(date.Month), date.Year)
	err := service.GetEventsForMonth(db, &month)
	if err != nil {
		return err
	}

	templates.Calendar(month).
		Render(context.Background(), c.Response().Writer)
	return nil
}
