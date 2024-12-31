package views

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/calendar"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/service"
)

func InitIndexRoutes(group *echo.Group, r *repo.Queries) {
	group.GET(
		"",
		func(c echo.Context) error { return HandleIndex(c, r) },
	)
	group.GET("/error", HandleError)
}

func HandleIndex(c echo.Context, r *repo.Queries) error {
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorPage(http.StatusNotFound, err.Error(), c)
		return err
	}
	vacDays, err := service.GetVacationCountForUserYear(
		r,
		int(currUser.ID),
		calendar.CurrentYear(),
	)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	stats := calendar.GetCurrentYearProgress()

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	pendingEvents, err := service.GetPendingEventsForYear(r, currUser.ID, calendar.CurrentYear())
	if err != nil {
		htmx.ErrorPage(http.StatusBadRequest, err.Error(), c)
		return err
	}

	templates.Home(currUser, vacDays, pendingEvents, stats, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleError(c echo.Context) error {
	templates.Error(http.StatusInternalServerError, "").
		Render(context.Background(), c.Response().Writer)
	return nil
}
