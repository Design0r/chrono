package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type TeamHandler struct {
	event service.EventService
	notif service.NotificationService
	user  service.UserService
	log   *slog.Logger
}

func NewTeamHandler(
	e service.EventService,
	n service.NotificationService,
	u service.UserService,
	log *slog.Logger,
) TeamHandler {
	return TeamHandler{event: e, notif: n, user: u, log: log}
}

func (h *TeamHandler) RegisterRoutes(authGrp *echo.Group, adminGrp *echo.Group) {
	authGrp.GET("/team", h.Team)
	adminGrp.GET("/team/:id/edit", h.GetTeamRowForm)
	adminGrp.GET("/team/:id", h.GetTeamRow)
	adminGrp.PATCH("/team/:id", h.PatchTeamRow)
}

func (h *TeamHandler) Team(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	allUserswithVac, err := h.event.GetAllUsersWithVacation(ctx, domain.CurrentYear())
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user data.")
	}

	notifs, err := h.notif.GetByUserId(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(c, http.StatusOK, templates.Team(allUserswithVac, currUser, notifs))
}

func (h *TeamHandler) GetTeamRow(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid user id")
	}

	user, err := h.event.GetUserWithVacation(ctx, id, domain.CurrentYear(), 1)
	if err != nil {
		return RenderError(c, http.StatusNotFound, "user not found")
	}

	return Render(c, http.StatusOK, templates.TeamRow(currUser, user, false))
}

func (h *TeamHandler) GetTeamRowForm(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid user id")
	}

	user, err := h.event.GetUserWithVacation(ctx, id, domain.CurrentYear(), 1)
	if err != nil {
		return RenderError(c, http.StatusNotFound, "user not found")
	}

	return Render(c, http.StatusOK, templates.TeamRow(currUser, user, true))
}

func (h *TeamHandler) PatchTeamRow(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid user id")
	}

	vacParam := c.FormValue("vacation")
	vac, err := strconv.Atoi(vacParam)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid vacation number")
	}

	err = h.user.SetVacation(ctx, id, vac, domain.CurrentYear())
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed to update user vacation")
	}

	user, err := h.event.GetUserWithVacation(ctx, id, domain.CurrentYear(), 1)
	if err != nil {
		return RenderError(c, http.StatusNotFound, "user not found")
	}

	return Render(c, http.StatusOK, templates.TeamRow(currUser, user, false))
}
