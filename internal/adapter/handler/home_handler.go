package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type HomeHandler struct {
	token service.TokenService
}

func NewHomeHandler() HomeHandler {
	return HomeHandler{}
}

func RegisterHomeRoutes(group *echo.Group, handler *HomeHandler) {
	group.GET("", handler.Home)
}

func (h *HomeHandler) Home(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)
	err := h.token.InitYearlyTokens(ctx, &currUser)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return nil
}
