package views

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/middleware"
	"chrono/schemas"
	"chrono/service"
)

func InitEventRoutes(group *echo.Group, r *repo.Queries) {
	group.POST(
		"/:year/:month/:day",
		func(c echo.Context) error { return CreateEventHandler(c, r) },
		middleware.SessionMiddleware(r),
	)
	group.DELETE(
		"/events/:id",
		func(c echo.Context) error { return DeleteEventHandler(c, r) },
		middleware.SessionMiddleware(r),
	)
}

func CreateEventHandler(c echo.Context, r *repo.Queries) error {
	var date schemas.YMDDate
	if err := c.Bind(&date); err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}
	eventName := c.FormValue("name")

	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	event, err := service.CreateEvent(r, date, currUser, eventName)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	e := schemas.Event{Username: currUser.Username, Event: event}

	vacationUsed, err := service.GetVacationCountForUserYear(r, int(currUser.ID), date.Year)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	pendingEvents, err := service.GetPendingEventsForYear(
		r,
		currUser.ID,
		event.ScheduledAt.Year(),
	)

	templates.CreateEventUpdate(e, currUser, vacationUsed, pendingEvents, len(notifications)).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func DeleteEventHandler(c echo.Context, r *repo.Queries) error {
	event := c.Param("id")
	eventId, err := strconv.Atoi(event)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	deletedEvent, err := service.DeleteEvent(r, eventId)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	e := schemas.Event{Username: currUser.Username, Event: deletedEvent}

	vacationUsed, err := service.GetVacationCountForUserYear(
		r,
		int(currUser.ID),
		deletedEvent.ScheduledAt.Year(),
	)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	pendingEvents, err := service.GetPendingEventsForYear(
		r,
		currUser.ID,
		deletedEvent.ScheduledAt.Year(),
	)

	templates.CreateEventUpdate(e, currUser, vacationUsed, pendingEvents, len(notifications)).
		Render(context.Background(), c.Response().Writer)
	return nil
}
