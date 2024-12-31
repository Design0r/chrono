package views

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/calendar"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/schemas"
	"chrono/service"
)

func InitCalendarRoutes(group *echo.Group, r *repo.Queries) {
	group.GET(
		"/:year/:month",
		func(c echo.Context) error { return MonthCalendarHandler(c, r) },
	)
}

func MonthCalendarHandler(c echo.Context, r *repo.Queries) error {
	user, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorPage(http.StatusBadRequest, err.Error(), c)
		return err
	}

	var date schemas.YMDate
	if err := c.Bind(&date); err != nil {
		htmx.ErrorPage(http.StatusBadRequest, "Invalid parameter", c)
		return err
	}
	service.UpdateHolidays(r, date.Year)

	month := calendar.GetDaysOfMonth(time.Month(date.Month), date.Year)
	err = service.GetEventsForMonth(r, &month)
	if err != nil {
		return err
	}
	vacationUsed, err := service.GetVacationCountForUserYear(r, int(user.ID), date.Year)
	if err != nil {
		htmx.ErrorPage(http.StatusBadRequest, err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(r, user.ID)
	if err != nil {
		htmx.ErrorPage(http.StatusBadRequest, err.Error(), c)
		return err
	}

	pendingEvents, err := service.GetPendingEventsForYear(r, user.ID, date.Year)
	if err != nil {
		htmx.ErrorPage(http.StatusBadRequest, err.Error(), c)
		return err
	}

	templates.Calendar(user, month, vacationUsed, pendingEvents, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}
