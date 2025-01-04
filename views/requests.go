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
	currUser := c.Get("user").(repo.User)

	requests, _ := service.GetPendingRequests(r)

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Requests(&currUser, requests, notifications))
}

func HandlePatchRequests(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid request id")
	}

	stateParam := c.FormValue("state")

	err = service.UpdateRequestState(r, stateParam, currUser, int64(id))
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed updating request")
	}

	return Render(
		c,
		http.StatusOK,
		htmx.SuccessMessage(fmt.Sprintf("%v %v", strings.Title(stateParam), "Request")),
	)
}
