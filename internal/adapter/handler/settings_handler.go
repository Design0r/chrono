package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type SettingsHandler struct {
	s service.SettingsService
}

func NewSettingsHandler(s service.SettingsService) SettingsHandler {
	return SettingsHandler{s: s}
}

func RegisterSettingsRoutes(group *echo.Group, handler *SettingsHandler) {
	g := group.Group("/settings")
	g.GET("", handler.Settings)
}

func (h *SettingsHandler) Settings(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	s, err := h.s.GetFirst(c.Request().Context())
	if err != nil {
		RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Settings(s, &currUser, []domain.Notification{}))
}
