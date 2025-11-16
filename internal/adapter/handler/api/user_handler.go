package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APIUserHandler struct {
	user  service.UserService
	event service.EventService
	auth  service.AuthService
	log   *slog.Logger
}

func NewAPIUserHandler(
	u service.UserService,
	e service.EventService,
	a service.AuthService,
	log *slog.Logger,
) APIUserHandler {
	return APIUserHandler{user: u, event: e, auth: a, log: log}
}

func (h *APIUserHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/users/:id", h.GetUserById)
	group.PATCH("/users/:id", h.ProfileEdit)
	group.GET("/users", h.GetUsers)
}

func (h *APIUserHandler) GetUserById(c echo.Context) error {
	ctx := c.Request().Context()

	id := c.Param("id")
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
	}

	yearParam := c.QueryParam("year")

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		year = domain.CurrentYear()
	}

	vacParam := c.QueryParam("vacation")
	vacation, err := strconv.ParseBool(vacParam)
	if err != nil {
		vacation = false
	}

	if vacation {
		user, err := h.event.GetUserWithVacation(ctx, userId, year, 1)
		if err != nil {
			return NewErrorResponse(c, http.StatusNotFound, "user not found")
		}

		return NewJsonResponse(c, user)
	}

	user, err := h.user.GetById(ctx, userId)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, "user not found")
	}

	return NewJsonResponse(c, user)
}

func (h *APIUserHandler) GetUsers(c echo.Context) error {
	ctx := c.Request().Context()

	yearParam := c.QueryParam("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		year = domain.CurrentYear()
	}

	vacParam := c.QueryParam("vacation")
	vacation, err := strconv.ParseBool(vacParam)
	if err != nil {
		vacation = false
	}

	var users any
	if vacation {
		users, err = h.event.GetAllUsersWithVacation(ctx, year)
		if err != nil {
			return NewErrorResponse(c, http.StatusNotFound, "user not found")
		}
	} else {
		users, err = h.user.GetAll(ctx)
		if err != nil {
			return NewErrorResponse(c, http.StatusInternalServerError, "failed to fetch users")
		}
	}

	return NewJsonResponse(c, users)
}

func (h *APIUserHandler) ProfileEdit(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	patchedData := domain.ApiPatchUser{}
	if err := c.Bind(&patchedData); err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "Invalid data")
	}

	u := &domain.User{
		ID:           currUser.ID,
		Username:     patchedData.Name,
		Email:        patchedData.Email,
		Color:        patchedData.Color,
		Role:         currUser.Role,
		Enabled:      currUser.Enabled,
		IsSuperuser:  currUser.IsSuperuser,
		VacationDays: currUser.VacationDays,
		Password:     currUser.Password,
	}

	if patchedData.Password != "" {
		pw, err := h.auth.HashPassword(patchedData.Password)
		if err != nil {
			return NewErrorResponse(
				c,
				http.StatusInternalServerError,
				"Failed to update user information.",
			)
		}
		u.Password = pw
	}

	updatedUser, err := h.user.Update(
		c.Request().Context(),
		u,
	)
	if err != nil {
		return NewErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to update user information.",
		)
	}

	return NewJsonResponse(c, updatedUser)
}
