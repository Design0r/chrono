package api

import (
	"fmt"
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
	ctx := c.Request().Context()

	id := c.Param("id")
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid user id")
	}

	userToEdit, err := h.user.GetById(ctx, userId)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, "user id does not exist")
	}

	patchedData := domain.ApiPatchUser{}
	if err := c.Bind(&patchedData); err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "Invalid data")
	}

	username := userToEdit.Username
	if patchedData.Name != "" {
		username = patchedData.Name
	}

	email := userToEdit.Email
	if patchedData.Email != "" {
		email = patchedData.Email
	}

	aworkId := userToEdit.AworkID
	if patchedData.AworkID != nil {
		aworkId = patchedData.AworkID
	}

	color := userToEdit.Color
	if patchedData.Color != "" {
		color = patchedData.Color
	}

	vacDays := userToEdit.VacationDays
	if currUser.IsAdmin() && patchedData.VacationDays != nil {
		vacDays = *patchedData.VacationDays
	}

	role := userToEdit.Role
	if currUser.IsAdmin() && patchedData.Role != "" {
		if !domain.IsValidRole((domain.Role)(patchedData.Role)) {
			return NewErrorResponse(c, http.StatusUnprocessableEntity, "Invalid user role")
		}
		role = patchedData.Role
	}

	superuser := role == "admin"

	enabled := userToEdit.Enabled
	if currUser.IsAdmin() && patchedData.Enabled != nil {
		enabled = *patchedData.Enabled
	}

	u := &domain.User{
		ID:           userToEdit.ID,
		Username:     username,
		Email:        email,
		Color:        color,
		Role:         role,
		AworkID:      aworkId,
		Enabled:      enabled,
		IsSuperuser:  superuser,
		VacationDays: vacDays,
		Password:     currUser.Password,
	}

	fmt.Println(u)

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
