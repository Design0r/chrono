package views

import (
	"fmt"
	"net/http"
	"strconv"
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
	group.GET(
		"/requests/modal",
		func(c echo.Context) error { return HandleRequestModal(c, r) },
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

func HandleRequestModal(c echo.Context, r *repo.Queries) error {
	reqUserIdParam := c.FormValue("user_id")
	startDateParam := c.FormValue("start_date")
	endDateParam := c.FormValue("end_date")
	reqIdParam := c.FormValue("request_id")

	reqUserId, err := strconv.ParseInt(reqUserIdParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed parsing params")
	}
	reqId, err := strconv.ParseInt(reqIdParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed parsing params")
	}
	startUnix, err := strconv.ParseInt(startDateParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed parsing params")
	}
	endUnix, err := strconv.ParseInt(endDateParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed parsing params")
	}
	startDate := time.Unix(startUnix, 0)
	endDate := time.Unix(endUnix, 0)

	requests, err := service.GetRequestRange(r, startDate, endDate, reqUserId)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed getting requests")
	}

	return Render(
		c,
		http.StatusOK,
		templates.RejectModal(requests[0].Message, startDate, endDate, reqUserId, reqId),
	)
}

func HandlePatchRequests(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	reqUserIdParam := c.FormValue("user_id")
	stateParam := c.FormValue("state")
	startDateParam := c.FormValue("start_date")
	endDateParam := c.FormValue("end_date")
	reasonParam := c.FormValue("reason")

	reqUserId, err := strconv.ParseInt(reqUserIdParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed parsing start date")
	}
	startUnix, err := strconv.ParseInt(startDateParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed parsing start date")
	}
	endUnix, err := strconv.ParseInt(endDateParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed parsing start date")
	}
	startDate := time.Unix(startUnix, 0)
	endDate := time.Unix(endUnix, 0)

	err = service.UpdateRequestStateRange(
		r,
		currUser,
		reqUserId,
		stateParam,
		startDate,
		endDate,
		reasonParam,
	)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed updating request")
	}

	return Render(
		c,
		http.StatusOK,
		htmx.SuccessMessage(fmt.Sprintf("%v %v", strings.Title(stateParam), "Request")),
	)
}
