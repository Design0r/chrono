package views

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/htmx"
	"chrono/middleware"
	"chrono/schemas"
	"chrono/service"
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
		htmx.ErrorMessage(err.Error(), c)
		return err
	}
	eventName := c.FormValue("name")

	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	event, err := service.CreateEvent(db, date, currUser, eventName)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	e := schemas.Event{Username: currUser.Username, Event: event}

	vacationUsed, err := service.GetVacationCountForUserYear(db, int(currUser.ID), date.Year)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	pendingEvents, err := service.GetPendingEventsForYear(
		db,
		currUser.ID,
		event.ScheduledAt.Year(),
	)

	templates.CreateEventUpdate(e, currUser, vacationUsed, pendingEvents, len(notifications)).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func DeleteEventHandler(c echo.Context, db *sql.DB) error {
	event := c.Param("id")
	eventId, err := strconv.Atoi(event)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	deletedEvent, err := service.DeleteEvent(db, eventId)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	e := schemas.Event{Username: currUser.Username, Event: deletedEvent}

	vacationUsed, err := service.GetVacationCountForUserYear(
		db,
		int(currUser.ID),
		deletedEvent.ScheduledAt.Year(),
	)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	pendingEvents, err := service.GetPendingEventsForYear(
		db,
		currUser.ID,
		deletedEvent.ScheduledAt.Year(),
	)

	templates.CreateEventUpdate(e, currUser, vacationUsed, pendingEvents, len(notifications)).
		Render(context.Background(), c.Response().Writer)
	return nil
}
