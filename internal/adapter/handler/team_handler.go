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
	adminGrp.GET("/team/form", h.TeamForm)
	adminGrp.PATCH("/team", h.TeamEdit)
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

func (h *TeamHandler) TeamForm(c echo.Context) error {
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

	return Render(c, http.StatusOK, templates.TeamForm(allUserswithVac, currUser, notifs))
}

func (h *TeamHandler) TeamEdit(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	form, err := c.FormParams()
	if err != nil {
		return err
	}

	for k, v := range form {
		userId, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			continue
		}
		vacation, err := strconv.Atoi(v[0])
		if err != nil {
			continue
		}

		err = h.user.SetVacation(ctx, userId, vacation, domain.CurrentYear())
		if err != nil {
			continue
		}
	}

	allUserswithVac, err := h.event.GetAllUsersWithVacation(ctx, domain.CurrentYear())
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user data.")
	}

	notifs, err := h.notif.GetByUserId(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(c, http.StatusOK, templates.TeamHTMX(allUserswithVac, currUser, notifs))
}
