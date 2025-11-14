package api

import (
	"log/slog"
	"net/http"
	"strconv"

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
	group.GET("/events/:year", h.GetVacationGraph)
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

func (h *APIEventHandler) GetVacationGraph(c echo.Context) error {
	ctx := c.Request().Context()

	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, "invalid year parameter")
	}

	data, err := h.event.GetHistogramForYear(ctx, year)
	if err != nil {
		return NewErrorResponse(
			c,
			http.StatusInternalServerError,
			"failed to fetch data for vacation graph",
		)
	}

	yearOffset := domain.GetYearOffset(year)
	monthGaps := domain.GetMonthGaps(year)

	response := map[string]any{
		"year_offset":   yearOffset,
		"month_gaps":    monthGaps,
		"vacation_data": data,
	}

	return NewJsonResponse(c, response)
}
