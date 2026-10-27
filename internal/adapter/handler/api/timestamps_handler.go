package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APITimestampsHandler struct {
	timestamps *service.TimestampsService
	user       *service.UserService
}

func NewAPITimestampsHandler(t *service.TimestampsService, u *service.UserService) APITimestampsHandler {
	return APITimestampsHandler{timestamps: t, user: u}
}

func (s *APITimestampsHandler) RegisterRoutes(auth *echo.Group, admin *echo.Group) {
	g := auth.Group("/timestamps")
	g.POST("", s.Start)
	g.POST("/new", s.Create)
	g.PATCH("/:id", s.Stop)
	g.PUT("/:id", s.Update)
	g.GET("/day", s.GetTimestampsForToday)
	g.GET("", s.GetTimestamps)
	g.GET("/latest", s.GetLatestTimestamp)
	g.GET("/worked/:year", s.GetWorkHoursForYear)
	g.GET("/worked/:year/all", s.GetWorkHoursForYearForAllUsers)

	admin.GET("/timestamps/all", s.GetAllTimestamps)
	admin.GET("/timestamps/all/summary", s.GetWorkedSecondsPerUser)
}

func (h *APITimestampsHandler) Start(c *echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	t, err := h.timestamps.Start(ctx, currUser.ID)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) Stop(c *echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid id")
	}

	// same response for foreign and missing ids, so ids of others can't be probed
	ts, err := h.timestamps.GetById(ctx, id)
	if err != nil || (ts.UserID != currUser.ID && !currUser.IsAdmin()) {
		return NewErrorResponse(c, http.StatusNotFound, "timestamp not found")
	}

	t, err := h.timestamps.Stop(ctx, id)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) GetTimestampsForToday(c *echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	t, err := h.timestamps.GetForToday(ctx, currUser.ID)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) GetLatestTimestamp(c *echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	t, err := h.timestamps.GetLatest(ctx, currUser.ID)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) Create(c *echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	if !currUser.IsAdmin() {
		return NewErrorResponse(c, http.StatusForbidden, "only allowed for admins")
	}

	var tsForm domain.Timestamp
	if err := c.Bind(&tsForm); err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid form parameters")
	}

	if tsForm.UserID == 0 {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "select a valid user id")
	}

	t, err := h.timestamps.Create(ctx, &tsForm)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) Update(c *echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	if !currUser.IsAdmin() {
		return NewErrorResponse(c, http.StatusForbidden, "only allowed for admins")
	}

	var tsForm domain.Timestamp
	if err := c.Bind(&tsForm); err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid form parameters")
	}

	t, err := h.timestamps.Update(ctx, &tsForm)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) GetTimestamps(c *echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	t, err := h.timestamps.GetInRange(ctx, currUser.ID, startDate, endDate)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) GetWorkHoursForYear(c *echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "year parameter is missing")
	}

	work, err := h.timestamps.GetWorkHoursForYear(ctx, currUser.ID, year, currUser.WorkdayHours)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, work)
}

func (h *APITimestampsHandler) GetWorkHoursForYearForAllUsers(c *echo.Context) error {
	currUser := c.Get("user").(domain.User)

	if !currUser.IsAdmin() {
		return NewErrorResponse(c, http.StatusForbidden, "only allowed for admins")
	}

	ctx := c.Request().Context()

	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "year parameter is missing")
	}

	inactiveParam := c.QueryParam("includeInactive")
	includeInactive, err := strconv.ParseBool(inactiveParam)
	if err != nil {
		includeInactive = false
	}

	users, err := h.user.GetAll(ctx, includeInactive)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, "failed to get users")
	}

	work := h.timestamps.GetWorkHoursForYearForAllUsers(ctx, users, year)

	return NewJsonResponse(c, work)
}

func (h *APITimestampsHandler) GetAllTimestamps(c *echo.Context) error {
	ctx := c.Request().Context()

	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	userParam := c.QueryParam("userId")
	if userParam == "" {
		t, err := h.timestamps.GetAllInRange(ctx, startDate, endDate)
		if err != nil {
			return NewErrorResponse(c, http.StatusNotFound, err.Error())
		}
		return NewJsonResponse(c, t)
	}

	userId, err := strconv.ParseInt(userParam, 10, 64)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "invalid userId")
	}

	t, err := h.timestamps.GetInRange(ctx, userId, startDate, endDate)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	return NewJsonResponse(c, t)
}

func (h *APITimestampsHandler) GetWorkedSecondsPerUser(c *echo.Context) error {
	ctx := c.Request().Context()

	startDate, endDate, err := parseDateRange(c)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	worked, err := h.timestamps.GetWorkedSecondsPerUserInRange(ctx, startDate, endDate)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return NewJsonResponse(c, worked)
}

func parseDateRange(c *echo.Context) (time.Time, time.Time, error) {
	startDate := time.UnixMilli(0)
	endDate := time.Now()

	if startParam := c.QueryParam("startDate"); startParam != "" {
		s, err := time.Parse(time.DateOnly, startParam)
		if err != nil {
			return startDate, endDate, errors.New("invalid startDate")
		}
		startDate = s
	}

	if endParam := c.QueryParam("endDate"); endParam != "" {
		e, err := time.Parse(time.DateOnly, endParam)
		if err != nil {
			return startDate, endDate, errors.New("invalid endDate")
		}
		endDate = e
	}

	return startDate, endDate, nil
}
