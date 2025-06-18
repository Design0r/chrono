package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/adapter/htmx"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type RequestHandler struct {
	log      *slog.Logger
	request  service.RequestService
	event    service.EventService
	notif    service.NotificationService
	vacation service.VacationTokenService
}

func NewRequestHandler(
	r service.RequestService,
	n service.NotificationService,
	e service.EventService,
	v service.VacationTokenService,
	log *slog.Logger,
) RequestHandler {
	return RequestHandler{request: r, notif: n, event: e, vacation: v, log: log}
}

func (h *RequestHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/requests", h.Requests)
	group.GET("/requests/modal", h.RejectModal)
	group.PATCH("/requests", h.PatchRequests)
}

func (h *RequestHandler) Requests(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	requests, err := h.request.GetPending(ctx)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get pending requests.")
	}
	fmt.Println(requests[1].Conflicts)

	notifications, err := h.notif.GetByUserId(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(c, http.StatusOK, templates.Requests(&currUser, requests, notifications))
}

func (h *RequestHandler) RejectModal(c echo.Context) error {
	var form domain.RejectModalForm
	if err := c.Bind(&form); err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid parameter.")
	}

	startDate := time.Unix(form.StartDate, 0).UTC()
	endDate := time.Unix(form.EndDate, 0).UTC()

	requests, err := h.request.GetInRange(
		c.Request().Context(),
		form.UserID,
		startDate,
		endDate,
	)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed getting requests")
	}

	return Render(
		c,
		http.StatusOK,
		templates.RejectModal(
			requests[0].Message,
			startDate,
			endDate,
			form.UserID,
			form.RequestID,
		),
	)
}

func (h *RequestHandler) PatchRequests(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	form := domain.PatchRequestForm{}
	if err := c.Bind(&form); err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid parameter.")
	}
	startDate := time.Unix(form.StartDate, 0).UTC()
	endDate := time.Unix(form.EndDate, 0).UTC()
	fmt.Println(form, startDate, endDate)

	reqId, err := h.request.UpdateInRange(ctx, currUser.ID, form)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed updating request")
	}

	err = h.event.UpdateInRange(ctx, form.UserID, form.State, startDate, endDate)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed updating events")
	}

	eventName, err := h.request.GetEventName(ctx, reqId)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed getting event name")
	}

	if form.State == "accepted" {
		days := (startDate.Sub(endDate).Hours() / 24) - 1.0
		if eventName == "urlaub halbtags" {
			days /= 2
		}
		h.vacation.Create(ctx, days, startDate.Year(), form.UserID)
	}

	return Render(
		c,
		http.StatusOK,
		htmx.SuccessMessage(fmt.Sprintf("%v %v", strings.Title(form.State), "Request")),
	)
}
