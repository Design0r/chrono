package api

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APIEventHandler struct {
	user  service.UserService
	event service.EventService
	log   *slog.Logger
}

func NewAPIEventHandler(
	u service.UserService,
	e service.EventService,
	log *slog.Logger,
) APIEventHandler {
	return APIEventHandler{user: u, event: e, log: log}
}

func (h *APIEventHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/events/:year/:month", h.GetEventsForMonth)
}

func (h *APIEventHandler) GetEventsForMonth(c echo.Context) error {
	ctx := c.Request().Context()

	var date domain.YMDate
	if err := c.Bind(&date); err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, "Invalid date")
	}

	eventFilter := c.QueryParam("event-filter")
	userFilter := c.QueryParam("filter")
	var filtered *domain.User
	if userFilter != "" {
		filteredUser, err := h.user.GetByName(ctx, userFilter)
		if err == nil {
			filtered = filteredUser
		}
	}

	month, err := h.event.GetForMonth(ctx, date, filtered, eventFilter)
	if err != nil {
		return NewErrorResponse(
			c,
			http.StatusInternalServerError,
			"failed to fetch data for month",
		)
	}

	return NewJsonResponse(c, month)
}
