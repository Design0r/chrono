package views

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/middleware"
	"chrono/schemas"
	"chrono/service"
)

func InitProfileRoutes(group *echo.Group, r *repo.Queries) {
	group.GET(
		"/profile",
		func(c echo.Context) error { return HandleProfile(c, r) },
	)
	group.GET(
		"/profile/edit",
		func(c echo.Context) error { return HandleProfileEditForm(c, r) },
	)

	group.PATCH(
		"/profile",
		func(c echo.Context) error { return HandleProfileEdit(c, r) },
	)

	group.PUT(
		"profile/:id/admin",
		func(c echo.Context) error { return HandleToggleAdmin(c, r) },
		middleware.AdminMiddleware(r),
	)
}

func HandleProfile(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(
		c,
		http.StatusOK,
		templates.ProfilePage(currUser, notifications),
	)
}

func HandleProfileEditForm(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.ProfileEditForm(currUser, notifications))
}

func HandleProfileEdit(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	patchedData := schemas.PatchUser{}
	if err := c.Bind(&patchedData); err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	updatedUser, err := service.UpdateUser(
		r,
		repo.UpdateUserParams{
			Color:    patchedData.Color,
			Email:    patchedData.Email,
			Username: patchedData.Name,
			ID:       currUser.ID,
		},
	)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.UpdateProfileWithMessage(updatedUser, notifications))
}

func HandleToggleAdmin(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	user, err := service.ToggleAdmin(r, currUser, int64(id))
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(
		c,
		http.StatusOK,
		templates.AdminCheckbox(currUser, user.ID, user.IsSuperuser, true),
	)
}
