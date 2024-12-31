package views

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/calendar"
	"chrono/db/repo"
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
		return Render(c, http.StatusBadRequest, templates.Error(http.StatusBadRequest, err.Error()))
	}

	var date schemas.YMDate
	if err := c.Bind(&date); err != nil {
		return Render(c, http.StatusBadRequest, templates.Error(http.StatusBadRequest, err.Error()))
	}
	service.UpdateHolidays(r, date.Year)

	month := calendar.GetDaysOfMonth(time.Month(date.Month), date.Year)
	err = service.GetEventsForMonth(r, &month)
	if err != nil {
		return Render(c, http.StatusBadRequest, templates.Error(http.StatusBadRequest, err.Error()))
	}
	vacationUsed, err := service.GetVacationCountForUserYear(r, int(user.ID), date.Year)
	if err != nil {
		return Render(c, http.StatusBadRequest, templates.Error(http.StatusBadRequest, err.Error()))
	}

	notifications, err := service.GetUserNotifications(r, user.ID)
	if err != nil {
		return Render(c, http.StatusBadRequest, templates.Error(http.StatusBadRequest, err.Error()))
	}

	pendingEvents, err := service.GetPendingEventsForYear(r, user.ID, date.Year)
	if err != nil {
		return Render(c, http.StatusBadRequest, templates.Error(http.StatusBadRequest, err.Error()))
	}

	return Render(
		c,
		http.StatusOK,
		templates.Calendar(user, month, vacationUsed, pendingEvents, notifications),
	)
}
