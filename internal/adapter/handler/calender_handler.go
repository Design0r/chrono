package handler

import (
	"chrono/assets/templates"
	"chrono/internal/domain"
	"chrono/internal/service"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CalendarHandler struct {
	log   *slog.Logger
	user  service.UserService
	notif service.NotificationService
	event service.EventService
	token service.TokenService
}

func NewCalendarHandler(user service.UserService, n service.NotificationService, e service.EventService, t service.TokenService, log *slog.Logger) CalendarHandler {
	return CalendarHandler{log: log, user: user, notif: n, event: e, token: t}
}

func (h *CalendarHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/:year/:month", h.HandleCalendar)
}

func (h *CalendarHandler) HandleCalendar(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	var date domain.YMDate
	if err := c.Bind(&date); err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid date")
	}
	if date.Year >= 1900 {
		//TODO update holidays
	}

	err := h.token.InitYearlyTokens(ctx, &currUser)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to initialize vacation tokens.")
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

	return Render(c, http.StatusOK, templates.Calendar(userWithVac, month, notifs, allUsers, userFilter, eventFilter))
}
