package views

import (
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
		return Render(
			c,
			http.StatusInternalServerError,
			templates.Error(http.StatusInternalServerError, err.Error()),
		)
	}

	if !currUser.IsSuperuser {
		return Render(
			c,
			http.StatusForbidden,
			templates.Error(
				http.StatusForbidden,
				"This page is only accessible by admins",
			),
		)
	}

	requests, _ := service.GetPendingRequests(r)

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return Render(
			c,
			http.StatusInternalServerError,
			templates.Error(http.StatusInternalServerError, err.Error()),
		)
	}

	return Render(c, http.StatusOK, templates.Requests(&currUser, requests, notifications))
}

func HandlePatchRequests(c echo.Context, r *repo.Queries) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return Render(c, http.StatusBadRequest, htmx.ErrorMessage("Invalid request id", c))
	}

	stateParam := c.FormValue("state")

	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		return Render(
			c,
			http.StatusInternalServerError,
			htmx.ErrorMessage("Internal server error", c),
		)
	}

	if !currUser.IsSuperuser {
		return Render(
			c,
			http.StatusForbidden,
			htmx.ErrorMessage("Not authorized", c),
		)
	}

	err = service.UpdateRequestState(r, stateParam, currUser, int64(id))
	if err != nil {
		return Render(
			c,
			http.StatusForbidden,
			htmx.ErrorMessage("Failed updating request", c),
		)
	}

	return Render(
		c,
		http.StatusOK,
		htmx.SuccessMessage(fmt.Sprintf("%v %v", strings.Title(stateParam), "Request"), c),
	)
}
