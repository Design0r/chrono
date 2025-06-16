package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type CalendarHandler struct {
	log     *slog.Logger
	user    service.UserService
	notif   service.NotificationService
	event   service.EventService
	token   service.TokenService
	holiday service.HolidayService
}

func NewCalendarHandler(
	user service.UserService,
	n service.NotificationService,
	e service.EventService,
	t service.TokenService,
	h service.HolidayService,
	log *slog.Logger,
) CalendarHandler {
	return CalendarHandler{log: log, user: user, notif: n, event: e, token: t, holiday: h}
}

func (h *CalendarHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/:year/:month", h.Calendar)
	group.POST("/:year/:month/:day", h.CreateEvent)
	group.DELETE("/events/:id", h.CreateEvent)
}

func (h *CalendarHandler) Calendar(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	var date domain.YMDate
	if err := c.Bind(&date); err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid date")
	}
	if date.Year >= 1900 {
		h.holiday.Update(ctx, date.Year)
	}

	err := h.token.InitYearlyTokens(ctx, &currUser, date.Year)
	if err != nil {
		return RenderError(
			c,
			http.StatusInternalServerError,
			"Failed to initialize vacation tokens.",
		)
	}

	eventFilter := c.QueryParam("event-filter")
	userFilter := c.QueryParam("filter")
	var filtered *domain.User
	if userFilter != "" {
		filteredUser, err := h.user.GetByName(c.Request().Context(), userFilter)
		if err == nil {
			filtered = filteredUser
		}
	}

	month, err := h.event.GetForMonth(ctx, date, filtered, eventFilter)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get month")
	}

	userWithVac, err := h.event.GetUserWithVacation(ctx, currUser.ID, date.Year)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user data")
	}

	notifs, err := h.notif.GetByUserId(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications")
	}

	allUsers, err := h.user.GetAll(ctx)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get users")
	}

	return Render(
		c,
		http.StatusOK,
		templates.Calendar(userWithVac, month, notifs, allUsers, userFilter, eventFilter),
	)
}

func (h *CalendarHandler) CreateEvent(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	var date domain.YMDDate
	if err := c.Bind(&date); err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}
	eventName := c.FormValue("name")

	event, err := h.event.Create(ctx, date, eventName, &currUser)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed to create event.")
	}

	userWithVac, err := h.event.GetUserWithVacation(ctx, currUser.ID, date.Year)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user data.")
	}

	notifs, err := h.notif.GetByUserId(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(
		c,
		http.StatusOK,
		templates.CreateEventUpdate(
			domain.EventUser{Event: *event, User: currUser},
			userWithVac,
			len(notifs),
		),
	)
}

func (h *CalendarHandler) DeleteEvent(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid event id.")
	}

	event, err := h.event.Delete(ctx, eventId, &currUser)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed to delete event.")
	}

	userWithVac, err := h.event.GetUserWithVacation(ctx, currUser.ID, event.ScheduledAt.Year())
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user data.")
	}

	notifs, err := h.notif.GetByUserId(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(
		c,
		http.StatusOK,
		templates.CreateEventUpdate(
			domain.EventUser{Event: *event, User: currUser},
			userWithVac,
			len(notifs),
		),
	)
}
