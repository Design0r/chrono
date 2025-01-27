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
		func(c echo.Context) error { return HandleCalendar(c, r) },
	)
}

func HandleCalendar(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	var date schemas.YMDate
	if err := c.Bind(&date); err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	if date.Year >= 1900 {
		service.UpdateHolidays(r, date.Year)
	}
	service.InitYearlyTokens(r, currUser, date.Year)

	month := calendar.GetDaysOfMonth(time.Month(date.Month), date.Year)
	err := service.GetEventsForMonth(r, &month)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}
	vacationRemaining, err := service.GetRemainingVacation(
		r,
		currUser.ID,
		date.Year,
		date.Month,
	)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	pendingEvents, err := service.GetPendingEventsForYear(r, currUser.ID, date.Year)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	return Render(
		c,
		http.StatusOK,
		templates.Calendar(currUser, month, vacationRemaining, pendingEvents, notifications),
	)
}
