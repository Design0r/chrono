package views

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/htmx"
	"calendar/middleware"
	"calendar/schemas"
	"calendar/service"
)

func InitEventRoutes(group *echo.Group, db *sql.DB) {
	group.POST(
		"/:year/:month/:day",
		func(c echo.Context) error { return CreateEventHandler(c, db) },
		middleware.SessionMiddleware(db),
	)
	group.DELETE(
		"/events/:id",
		func(c echo.Context) error { return DeleteEventHandler(c, db) },
		middleware.SessionMiddleware(db),
	)
}

func CreateEventHandler(c echo.Context, db *sql.DB) error {
	var date schemas.YMDDate
	if err := c.Bind(&date); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter")
	}
	eventName := c.FormValue("name")

	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Current user not found")
	}

	event, err := service.CreateEvent(db, date, currUser, eventName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}

	e := schemas.Event{Username: currUser.Username, Event: event}

	vacationUsed, err := service.GetVacationCountForUserYear(db, int(currUser.ID), date.Year)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorMessage("Failed to get notifications", c)
		return err
	}

	templates.CreateEventUpdate(e, currUser, vacationUsed, len(notifications)).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func DeleteEventHandler(c echo.Context, db *sql.DB) error {
	event := c.Param("id")
	eventId, err := strconv.Atoi(event)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid event id")
	}
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Current user not found")
	}

	deletedEvent, err := service.DeleteEvent(db, eventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid event id")
	}

	e := schemas.Event{Username: currUser.Username, Event: deletedEvent}

	vacationUsed, err := service.GetVacationCountForUserYear(
		db,
		int(currUser.ID),
		deletedEvent.ScheduledAt.Year(),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorMessage("Failed getting notifications", c)
		return err
	}

	templates.CreateEventUpdate(e, currUser, vacationUsed, len(notifications)).
		Render(context.Background(), c.Response().Writer)
	return nil
}
