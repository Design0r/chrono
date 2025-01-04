package views

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/schemas"
	"chrono/service"
)

func InitEventRoutes(group *echo.Group, r *repo.Queries) {
	group.POST(
		"/:year/:month/:day",
		func(c echo.Context) error { return CreateEventHandler(c, r) },
	)
	group.DELETE(
		"/events/:id",
		func(c echo.Context) error { return DeleteEventHandler(c, r) },
	)
}

func CreateEventHandler(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	var date schemas.YMDDate
	if err := c.Bind(&date); err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}
	eventName := c.FormValue("name")

	event, err := service.CreateEvent(r, date, currUser, eventName)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	e := schemas.Event{Username: currUser.Username, Event: event}

	vacationUsed, err := service.GetVacationCountForUserYear(r, int(currUser.ID), date.Year)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	pendingEvents, err := service.GetPendingEventsForYear(
		r,
		currUser.ID,
		event.ScheduledAt.Year(),
	)

	return Render(
		c,
		http.StatusOK,
		templates.CreateEventUpdate(e, currUser, vacationUsed, pendingEvents, len(notifications)),
	)
}

func DeleteEventHandler(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	event := c.Param("id")
	eventId, err := strconv.Atoi(event)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	deletedEvent, err := service.DeleteEvent(r, eventId)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	e := schemas.Event{Username: currUser.Username, Event: deletedEvent}

	vacationUsed, err := service.GetVacationCountForUserYear(
		r,
		int(currUser.ID),
		deletedEvent.ScheduledAt.Year(),
	)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	pendingEvents, err := service.GetPendingEventsForYear(
		r,
		currUser.ID,
		deletedEvent.ScheduledAt.Year(),
	)

	return Render(
		c,
		http.StatusOK,
		templates.CreateEventUpdate(e, currUser, vacationUsed, pendingEvents, len(notifications)),
	)
}
