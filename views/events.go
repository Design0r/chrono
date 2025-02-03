package views

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/calendar"
	"chrono/db/repo"
	"chrono/schemas"
	"chrono/service"
)

func InitEventRoutes(group *echo.Group, r *repo.Queries) {
	group.POST(
		"/:year/:month/:day",
		func(c echo.Context) error { return HandleCreateEvent(c, r) },
	)
	group.DELETE(
		"/events/:id",
		func(c echo.Context) error { return HandleDeleteEvent(c, r) },
	)
}

func HandleCreateEvent(c echo.Context, r *repo.Queries) error {
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

	e := schemas.Event{Username: currUser.Username, Color: currUser.Color, Event: event}

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

	pendingEvents, err := service.GetPendingEventsForYear(
		r,
		currUser.ID,
		event.ScheduledAt.Year(),
	)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(
		c,
		http.StatusOK,
		templates.CreateEventUpdate(
			e,
			currUser,
			vacationRemaining,
			vacTaken,
			pendingEvents,
			len(notifications),
		),
	)
}

func HandleDeleteEvent(c echo.Context, r *repo.Queries) error {
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

	if deletedEvent.State == "accepted" {
		_, err = service.CreateToken(r, currUser.ID, deletedEvent.ScheduledAt.Year(), 1.0)
		if err != nil {
			return RenderError(c, http.StatusBadRequest, err.Error())
		}
	}

	e := schemas.Event{Username: currUser.Username, Color: currUser.Color, Event: deletedEvent}

	vacRemaining, err := service.GetRemainingVacation(
		r,
		currUser.ID,
		deletedEvent.ScheduledAt.Year(),
		int(deletedEvent.ScheduledAt.Month()),
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

	pendingEvents, err := service.GetPendingEventsForYear(
		r,
		currUser.ID,
		deletedEvent.ScheduledAt.Year(),
	)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(
		c,
		http.StatusOK,
		templates.CreateEventUpdate(
			e,
			currUser,
			vacRemaining,
			vacTaken,
			pendingEvents,
			len(notifications),
		),
	)
}
