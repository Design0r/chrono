package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APISettingsHandler struct {
	settings service.SettingsService
}

func NewAPISettingsHandler(s service.SettingsService) APISettingsHandler {
	return APISettingsHandler{settings: s}
}

func (s *APISettingsHandler) RegisterRoutes(group *echo.Group) {
	g := group.Group("/settings")
	g.GET("", s.Settings)
	g.PATCH("", s.PatchSettings)
}

func (h *APISettingsHandler) Settings(c echo.Context) error {
	s, err := h.settings.GetFirst(c.Request().Context())
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, s)
}

func (h *APISettingsHandler) PatchSettings(c echo.Context) error {
	settings := domain.Settings{ID: 1, SignupEnabled: c.FormValue("signup_enabled") == "on"}
	s, err := h.settings.Update(c.Request().Context(), settings)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, "Failed to update settings.")
	}

	return NewJsonResponse(c, s)
}
