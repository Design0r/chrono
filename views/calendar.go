package views

import (
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
		func(c echo.Context) error { return HandleCalendar(c, r) },
	)
}

func HandleCalendar(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	var date schemas.YMDate
	if err := c.Bind(&date); err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	userFilter := c.QueryParam("filter")
	var filtered *repo.User
	if userFilter != "" {
		filteredUser, err := service.GetUserByName(r, userFilter)
		if err == nil {
			filtered = &filteredUser
		}
	}

	if date.Year >= 1900 {
		service.UpdateHolidays(r, date.Year)
	}
	service.InitYearlyTokens(r, currUser, date.Year)

	month := calendar.GetDaysOfMonth(time.Month(date.Month), date.Year)
	eventFilter := c.QueryParam("event-filter")
	err := service.GetEventsForMonth(r, &month, filtered, eventFilter)
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

	vacTaken, err := service.GetVacationCountForUser(r, currUser.ID, calendar.CurrentYear())
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	pendingEvents, err := service.GetPendingEventsForYear(r, currUser.ID, date.Year)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	allUsers, err := service.GetAllUsers(r)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	if htmx.IsHTMXRequest(c) {
		return Render(
			c,
			http.StatusOK,
			templates.CalendarCoreResponse(month, currUser, userFilter, eventFilter),
		)
	}

	return Render(
		c,
		http.StatusOK,
		templates.Calendar(
			currUser,
			month,
			vacationRemaining,
			vacTaken,
			pendingEvents,
			notifications,
			allUsers,
			userFilter,
			eventFilter,
		),
	)
}
