package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APITimestampsHandler struct {
	timestamps *service.TimestampsService
}

func NewAPITimestampsHandler(t *service.TimestampsService) APITimestampsHandler {
	return APITimestampsHandler{timestamps: t}
}

func (s *APITimestampsHandler) RegisterRoutes(auth *echo.Group, admin *echo.Group) {
	g := auth.Group("/timestamps")
	g.POST("", s.Start)
	g.PATCH("/:id", s.Stop)
	g.PUT("/:id", s.Update)
	g.GET("/day", s.GetTimestampsForToday)
	g.GET("", s.GetTimestamps)
	g.GET("/latest", s.GetLatestTimestamp)
	g.GET("/worked/:year", s.GetWorkHoursForYear)

	admin.GET("/timestamps/all", s.GetAllTimestamps)
}

func (h *APITimestampsHandler) Start(c echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	t, err := h.timestamps.Start(ctx, currUser.ID)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) Stop(c echo.Context) error {
	ctx := c.Request().Context()
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid id")
	}

	t, err := h.timestamps.Stop(ctx, id)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) GetTimestampsForToday(c echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	t, err := h.timestamps.GetForToday(ctx, currUser.ID)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) GetLatestTimestamp(c echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	t, err := h.timestamps.GetLatest(ctx, currUser.ID)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) Update(c echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	var tsForm domain.Timestamp
	if err := c.Bind(&tsForm); err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid form parameters")
	}

	if !currUser.IsAdmin() && tsForm.UserID != currUser.ID {
		return NewErrorResponse(c, http.StatusMethodNotAllowed, "not allowed")
	}

	t, err := h.timestamps.Update(ctx, &tsForm)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) GetTimestamps(c echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	startParam := c.QueryParam("startDate")
	endParam := c.QueryParam("endDate")

	startDate := time.UnixMilli(0)
	endDate := time.Now()

	if startParam != "" {
		s, err := time.Parse(time.DateOnly, startParam)
		if err != nil {
			return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid startDate")
		}
		startDate = s
	}

	if endParam != "" {
		e, err := time.Parse(time.DateOnly, endParam)
		if err != nil {
			return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid startDate")
		}
		endDate = e
	}

	t, err := h.timestamps.GetInRange(ctx, currUser.ID, startDate, endDate)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) GetWorkHoursForYear(c echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "year parameter is missing")
	}

	work, err := h.timestamps.GetWorkHoursForYear(ctx, currUser.ID, year, currUser.WorkdayHours)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, work)
}

func (h *APITimestampsHandler) GetAllTimestamps(c echo.Context) error {
	ctx := c.Request().Context()

	startParam := c.QueryParam("startDate")
	endParam := c.QueryParam("endDate")

	startDate := time.UnixMilli(0)
	endDate := time.Now()

	if startParam != "" {
		s, err := time.Parse(time.DateOnly, startParam)
		if err != nil {
			return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid startDate")
		}
		startDate = s
	}

	if endParam != "" {
		e, err := time.Parse(time.DateOnly, endParam)
		if err != nil {
			return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid startDate")
		}
		endDate = e
	}

	t, err := h.timestamps.GetAllInRange(ctx, startDate, endDate)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return NewJsonResponse(c, t)
}
