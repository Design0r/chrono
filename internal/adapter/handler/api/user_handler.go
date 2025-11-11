package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/internal/service"
)

type APIUserHandler struct {
	user service.UserService
	log  *slog.Logger
}

func NewAPIUserHandler(
	u service.UserService,
	log *slog.Logger,
) APIUserHandler {
	return APIUserHandler{user: u, log: log}
}

func (h *APIUserHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/users/:id", h.GetUserById)
}

func (h *APIUserHandler) GetUserById(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
	}

	user, err := h.user.GetById(ctx, userId)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, "user not found")
	}

	return NewJsonResponse(c, user)
}
