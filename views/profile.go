package views

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/middleware"
	"chrono/schemas"
	"chrono/service"
)

func InitProfileRoutes(group *echo.Group, r *repo.Queries) {
	group.GET(
		"/profile",
		func(c echo.Context) error { return HandleProfile(c, r) },
		middleware.SessionMiddleware(r),
	)
	group.GET(
		"/profile/edit",
		func(c echo.Context) error { return HandleProfileEditForm(c, r) },
		middleware.SessionMiddleware(r),
	)

	group.PATCH(
		"/profile",
		func(c echo.Context) error { return HandleProfileEdit(c, r) },
		middleware.SessionMiddleware(r),
	)

	group.PUT("profile/:id/admin", func(c echo.Context) error { return HandleToggleAdmin(c, r) })
}

func HandleProfile(c echo.Context, r *repo.Queries) error {
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorPage(http.StatusNotFound, err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorPage(http.StatusNotFound, err.Error(), c)
		return err
	}
	templates.ProfilePage(currUser, notifications).Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleProfileEditForm(c echo.Context, r *repo.Queries) error {
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}
	templates.ProfileEditForm(currUser, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleProfileEdit(c echo.Context, r *repo.Queries) error {
	patchedData := schemas.PatchUser{}
	if err := c.Bind(&patchedData); err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	updatedUser, err := service.UpdateUser(
		r,
		repo.UpdateUserParams{
			VacationDays: int64(patchedData.Vacation),
			Email:        patchedData.Email,
			Username:     patchedData.Name,
			ID:           currUser.ID,
		},
	)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	templates.UpdateProfileWithMessage(updatedUser, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleToggleAdmin(c echo.Context, r *repo.Queries) error {
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	if !currUser.IsSuperuser {
		htmx.ErrorMessage("Admin rights are required to change your teams admin status", c)
		return err
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	user, err := service.ToggleAdmin(r, currUser, int64(id))
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	templates.AdminCheckbox(currUser, user.ID, user.IsSuperuser).
		Render(context.Background(), c.Response().Writer)
	return nil
}
