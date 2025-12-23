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
	g.PATCH("/:id", s.Start)
	g.GET("/day", s.GetTimestampsForToday)
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
