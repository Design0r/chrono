package handler

import (
	"chrono/internal/domain"
	"chrono/internal/service"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CalendarHandler struct {
	log  *slog.Logger
	user service.UserService
}

func NewCalendarHandler(log *slog.Logger, user service.UserService) CalendarHandler {
	return CalendarHandler{log: log, user: user}
}

func (h *CalendarHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/calendar", h.HandleCalendar)
}

func (h *CalendarHandler) HandleCalendar(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	var date domain.YMDate
	if err := c.Bind(&date); err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid date")
	}

	userFilter := c.QueryParam("filter")
	var filtered *domain.User
	if userFilter != "" {
		filteredUser, err := h.user.GetByName(c.Request().Context(), userFilter)
		if err == nil {
			filtered = filteredUser
		}
	}

	month, err := service.GetMonth(date.Year, date.Month)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get month")
	}

	return nil
}
