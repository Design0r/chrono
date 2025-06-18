package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type HomeHandler struct {
	token  service.TokenService
	event  service.EventService
	notifs service.NotificationService
}

func NewHomeHandler(
	t service.TokenService,
	e service.EventService,
	n service.NotificationService,
) HomeHandler {
	return HomeHandler{token: t, event: e, notifs: n}
}

func (h *HomeHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/", h.Home)
}

func (h *HomeHandler) Home(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)
	err := h.token.InitYearlyTokens(ctx, &currUser, domain.CurrentYear())
	if err != nil {
		return RenderError(
			c,
			http.StatusInternalServerError,
			"Failed to initialize vacation tokens",
		)
	}

	userWithVac, err := h.event.GetUserWithVacation(
		ctx,
		currUser.ID,
		time.Now().Year(),
		int(time.Now().Month()),
	)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user data.")
	}

	yearProgress := domain.GetCurrentYearProgress()

	notifs, err := h.notifs.GetByUserId(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	yearOverview, err := h.event.GetHistogramForYear(ctx, time.Now().Year())
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get event data.")
	}

	return Render(c, http.StatusOK, templates.Home(userWithVac, yearProgress, notifs, yearOverview))
}
