package views

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/service"
)

func InitDebugRoutes(group *echo.Group, r *repo.Queries) {
	group.GET("/debug", func(c echo.Context) error { return HandleDebug(c, r) })
	group.DELETE("/debug/tokens", func(c echo.Context) error { return HandleDeleteTokens(c, r) })
	group.POST(
		"/debug/tokens",
		func(c echo.Context) error { return HandleCreateTokensForAcceptedEvents(c, r) },
	)
	group.DELETE(
		"/debug/sessions",
		func(c echo.Context) error { return HandleDeleteSessions(c, r) },
	)
	group.PATCH("/debug/color", func(c echo.Context) error { return HandleUserColor(c, r) })
}

func HandleDebug(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(
		c,
		http.StatusOK,
		templates.Debug(&currUser, notifications),
	)
}

func HandleDeleteTokens(c echo.Context, r *repo.Queries) error {
	err := service.DebugResetTokens(r)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	err = service.DebugResetTokenRefresh(r)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Message("Reset token table", "success"))
}

func HandleDeleteSessions(c echo.Context, r *repo.Queries) error {
	err := service.DeleteAllSessions(r)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Message("Reset session table", "success"))
}

func HandleCreateTokensForAcceptedEvents(c echo.Context, r *repo.Queries) error {
	err := service.DebugCreateTokenForAcceptedEvents(r)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(
		c,
		http.StatusOK,
		templates.Message("Created tokens for accepted events", "success"),
	)
}

func HandleUserColor(c echo.Context, r *repo.Queries) error {
	users, err := service.GetAllUsers(r)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	service.SetUserColor(r, 1, service.HSLToHex(service.GenerateHSLFloat(1)))
	for _, user := range users {
		service.SetUserColor(r, user.ID, service.HSLToHex(service.GenerateHSLFloat(int(user.ID))))
	}

	return Render(c, http.StatusOK, templates.Message("Changed user default colors", "success"))
}
