package views

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/service"
)

func InitRequestRoutes(group *echo.Group, r *repo.Queries) {
	group.GET(
		"/requests",
		func(c echo.Context) error { return HandleRequests(c, r) },
	)
	group.PATCH(
		"/requests/:id",
		func(c echo.Context) error { return HandlePatchRequests(c, r) },
	)
}

func HandleRequests(c echo.Context, r *repo.Queries) error {
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	if !currUser.IsSuperuser {
		htmx.ErrorPage(http.StatusForbidden, "This page is only accessible by admins", c)
		return nil
	}

	requests, _ := service.GetPendingRequests(r)

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	templates.Requests(&currUser, requests, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandlePatchRequests(c echo.Context, r *repo.Queries) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		htmx.ErrorMessage("Invalid request id", c)
		return err
	}

	stateParam := c.FormValue("state")

	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorMessage("Internal Error", c)
		return err
	}

	if !currUser.IsSuperuser {
		htmx.ErrorMessage("Not authorized", c)
		return err
	}

	err = service.UpdateRequestState(r, stateParam, currUser, int64(id))
	if err != nil {
		htmx.ErrorMessage("Failed updating request", c)
		return err
	}

	htmx.SuccessMessage(fmt.Sprintf("%v %v", strings.Title(stateParam), "Request"), c)
	return nil
}
