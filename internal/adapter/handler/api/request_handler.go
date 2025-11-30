package api

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APIRequestsHandler struct {
	request  *service.RequestService
	event    *service.EventService
	vacation *service.VacationTokenService
	log      *slog.Logger
}

func NewAPIRequestsHandler(
	r *service.RequestService,
	e *service.EventService,
	v *service.VacationTokenService,
	log *slog.Logger,
) APIRequestsHandler {
	return APIRequestsHandler{request: r, event: e, vacation: v, log: log}
}

func (h *APIRequestsHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/requests", h.Requests)
	group.PATCH("/requests", h.PatchRequests)
}

func (h *APIRequestsHandler) Requests(c echo.Context) error {
	ctx := c.Request().Context()

	requests, err := h.request.GetPending(ctx)
	if err != nil {
		return NewErrorResponse(
			c,
			http.StatusInternalServerError,
			"failed to get pending requests.",
		)
	}

	return NewJsonResponse(c, requests)
}

func (h *APIRequestsHandler) PatchRequests(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	form := domain.ApiPatchRequestForm{}
	if err := c.Bind(&form); err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid parameters")
	}

	oldForm := domain.PatchRequestForm{
		UserID:    form.UserID,
		State:     form.State,
		Reason:    form.Reason,
		StartDate: form.StartDate.Unix(),
		EndDate:   form.EndDate.Unix(),
	}

	reqId, err := h.request.UpdateInRange(ctx, currUser.ID, oldForm)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, "failed updating request")
	}

	err = h.event.UpdateInRange(ctx, form.UserID, form.State, form.StartDate, form.EndDate)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, "Failed updating events")
	}

	eventName, err := h.request.GetEventName(ctx, reqId)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, "Failed getting event name")
	}

	if form.State == "accepted" {
		days := (form.StartDate.Sub(form.EndDate).Hours() / 24) - 1.0
		if eventName == "urlaub halbtags" {
			days /= 2
		}
		h.vacation.Create(ctx, days, form.StartDate.Year(), form.UserID)
	}

	return NewJsonResponse(c, nil)
}
