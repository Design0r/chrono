package api

import (
	"net/http"
	"strconv"

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

func (s *APITimestampsHandler) RegisterRoutes(group *echo.Group) {
	g := group.Group("/timestamps")
	g.POST("", s.Start)
	g.PATCH("/:id", s.Stop)
	g.PUT("/:id", s.Update)
	g.GET("/day", s.GetTimestampsForToday)
	g.GET("/latest", s.GetLatestTimestamp)
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
	ctx := c.Request().Context()

	var tsForm domain.Timestamp
	if err := c.Bind(&tsForm); err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid form parameters")
	}

	t, err := h.timestamps.Update(ctx, &tsForm)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return NewJsonResponse(c, t)
}
