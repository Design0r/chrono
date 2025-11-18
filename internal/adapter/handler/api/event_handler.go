package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APIEventHandler struct {
	user  service.UserService
	event service.EventService
	token service.TokenService
	log   *slog.Logger
}

func NewAPIEventHandler(
	u service.UserService,
	e service.EventService,
	t service.TokenService,
	log *slog.Logger,
) APIEventHandler {
	return APIEventHandler{user: u, event: e, token: t, log: log}
}

func (h *APIEventHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/events/:year/:month", h.GetEventsForMonth)
	group.GET("/events/:year", h.GetVacationGraph)
	group.POST("/events", h.CreateEvent)
	group.DELETE("/events/:id", h.DeleteEvent)
}

func (h *APIEventHandler) GetEventsForMonth(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	var date domain.YMDate
	if err := c.Bind(&date); err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, "Invalid date")
	}

	err := h.token.InitYearlyTokens(ctx, &currUser, date.Year)
	if err != nil {
		return NewErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to initialize vacation tokens",
		)
	}

	month, err := h.event.GetForMonth(ctx, date, nil, "")
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

func (h *APIEventHandler) CreateEvent(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	var eventForm domain.CreateEvent
	if err := c.Bind(&eventForm); err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	event, err := h.event.Create(
		ctx,
		domain.YMDDate{Year: eventForm.Year, Month: eventForm.Month, Day: eventForm.Day},
		strings.ToLower(eventForm.EventName),
		&currUser,
	)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, "Failed to create event.")
	}

	eventUser := domain.EventUser{User: currUser, Event: *event}

	return NewJsonResponse(c, eventUser)
}

func (h *APIEventHandler) DeleteEvent(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid event id")
	}

	_, err = h.event.Delete(ctx, id, &currUser)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, "Failed to delete event.")
	}

	return NewJsonResponse(c, nil)
}
