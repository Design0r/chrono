package views

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/calendar"
	"chrono/db/repo"
	"chrono/service"
)

func InitIndexRoutes(group *echo.Group, r *repo.Queries) {
	group.GET("", func(c echo.Context) error { return HandleIndex(c, r) })
}

func HandleIndex(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)
	vacDays, err := service.GetVacationCountForUserYear(
		r,
		int(currUser.ID),
		calendar.CurrentYear(),
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
		templates.Home(currUser, vacDays, pendingEvents, stats, notifications),
	)
}
