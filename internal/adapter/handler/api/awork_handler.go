package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APIAworkHandler struct {
	user  *service.UserService
	event *service.EventService
	awork *service.AworkService
	log   *slog.Logger
}

func NewAPIAworkHandler(
	u *service.UserService,
	e *service.EventService,
	aw *service.AworkService,
	log *slog.Logger,
) APIAworkHandler {
	return APIAworkHandler{user: u, event: e, awork: aw, log: log}
}

func (h *APIAworkHandler) RegisterRoutes(auth *echo.Group, admin *echo.Group) {
	auth.GET("/awork/:year", h.GetWorkHoursForYear)
	// the awork directory is fetched with the company api key, keep it to admins
	admin.GET("/awork/users", h.GetAworkUsers)
}

func (h *APIAworkHandler) GetWorkHoursForYear(c echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "year parameter is missing")
	}

	if currUser.AworkID == nil || *currUser.AworkID == "" {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "awork id is missing")
	}

	work, err := h.awork.GetWorkHoursForYear(ctx, &currUser, year)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, work)
}

func (h *APIAworkHandler) GetAworkUsers(c echo.Context) error {
	users, err := h.awork.GetUsers(c.Request().Context())
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, users)
}
