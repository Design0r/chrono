package views

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/calendar"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/service"
)

func InitTeamRoutes(group *echo.Group, r *repo.Queries) {
	group.GET("/team", func(c echo.Context) error { return HandleTeam(c, r) })
	group.PATCH("/team", func(c echo.Context) error { return HandleTeamPatch(c, r) })
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

	if htmx.IsHTMXRequest(c) {
		return Render(c, http.StatusOK, templates.TeamForm(users, currUser, notifications))
	}

	return Render(c, http.StatusOK, templates.Team(users, currUser, notifications))
}

func HandleTeamPatch(c echo.Context, r *repo.Queries) error {
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

		err = service.SetUserVacation(r, userId, vacation, calendar.CurrentYear())
		if err != nil {
			continue
		}
	}

	currUser := c.Get("user").(repo.User)

	users, err := service.GetAllVacUsers(r)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.TeamHTMX(users, currUser, notifications))
}
