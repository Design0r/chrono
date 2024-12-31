package views

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/service"
)

func InitTeamRoutes(group *echo.Group, r *repo.Queries) {
	group.GET("/team", func(c echo.Context) error { return HandleTeam(c, r) })
}

func HandleTeam(c echo.Context, r *repo.Queries) error {
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}
	users, err := service.GetAllVacUsers(r)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	templates.Team(users, currUser, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}
