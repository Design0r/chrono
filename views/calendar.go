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
	user, err := service.GetCurrentUser(db, c)
	var date schemas.YMDate
	if err := c.Bind(&date); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter")
	}
	service.UpdateHolidays(db, date.Year)

	month := service.GetDaysOfMonth(time.Month(date.Month), date.Year)
	err = service.GetEventsForMonth(db, &month)
	if err != nil {
		return err
	}
	vacationUsed, err := service.GetVacationCountForUserYear(db, int(user.ID), date.Year)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}

	templates.Calendar(user, month, vacationUsed).
		Render(context.Background(), c.Response().Writer)
	return nil
}
