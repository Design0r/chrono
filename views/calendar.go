package views

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/htmx"
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
		htmx.ErrorPage(http.StatusBadRequest, "Invalid parameter", c)
		return nil
	}
	service.UpdateHolidays(db, date.Year)

	month := service.GetDaysOfMonth(time.Month(date.Month), date.Year)
	err = service.GetEventsForMonth(db, &month)
	if err != nil {
		return err
	}
	vacationUsed, err := service.GetVacationCountForUserYear(db, int(user.ID), date.Year)
	if err != nil {
		htmx.ErrorPage(http.StatusBadRequest, err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(db, user.ID)
	if err != nil {
		htmx.ErrorPage(http.StatusBadRequest, err.Error(), c)
		return err
	}

	pendingEvents, err := service.GetPendingEventsForYear(db, user.ID, date.Year)
	if err != nil {
		htmx.ErrorPage(http.StatusBadRequest, err.Error(), c)
		return err
	}

	templates.Calendar(user, month, vacationUsed, pendingEvents, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}
