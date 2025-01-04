package views

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/service"
)

func InitTeamRoutes(group *echo.Group, r *repo.Queries) {
	group.GET("/team", func(c echo.Context) error { return HandleTeam(c, r) })
}

func HandleTeam(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	users, err := service.GetAllVacUsers(r)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Team(users, currUser, notifications))
}
