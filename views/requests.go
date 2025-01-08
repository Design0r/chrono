package views

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
		"/requests",
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
	log.Println("hello")
	currUser := c.Get("user").(repo.User)

	stateParam := c.FormValue("state")
	startDateParam := c.FormValue("start_date")
	endDateParam := c.FormValue("end_date")
	startDate, err := time.Parse("2022-01-02 00:00:00 +0100 +0100", startDateParam)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed parsing start date")
	}
	endDate, err := time.Parse("2022-01-02 00:00:00 +0100 +0100", endDateParam)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed parsing end date")
	}
	log.Println(startDate, endDate)

	err = service.UpdateRequestStateRange(r, currUser.ID, stateParam, startDate, endDate)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed updating request")
	}

	return Render(
		c,
		http.StatusOK,
		htmx.SuccessMessage(fmt.Sprintf("%v %v", strings.Title(stateParam), "Request")),
	)
}
