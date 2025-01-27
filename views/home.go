package views

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/calendar"
	"chrono/db/repo"
	"chrono/service"
)

func InitHomeRoutes(group *echo.Group, r *repo.Queries) {
	group.GET("", func(c echo.Context) error { return HandleHome(c, r) })
}

func HandleHome(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)
	service.InitYearlyTokens(r, currUser, time.Now().Year())
	remainingDays, err := service.GetRemainingVacation(
		r,
		currUser.ID,
		calendar.CurrentYear(),
		0,
	)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	stats := calendar.GetCurrentYearProgress()

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	pendingEvents, err := service.GetPendingEventsForYear(r, currUser.ID, calendar.CurrentYear())
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	return Render(
		c,
		http.StatusOK,
		templates.Home(currUser, remainingDays, pendingEvents, stats, notifications),
	)
}
