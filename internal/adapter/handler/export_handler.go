package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type ExportHandler struct {
	krank service.ExportService
}

func NewExportHandler(k service.ExportService) ExportHandler {
	return ExportHandler{krank: k}
}

func (s *ExportHandler) RegisterRoutes(group *echo.Group) {
	g := group.Group("/export")
	g.GET("", s.Export)
	g.GET("/all", s.ExportAll)
}

func (h *ExportHandler) Export(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	return Render(c, http.StatusOK, templates.Export(currUser, []domain.Notification{}))
}

func (h *ExportHandler) ExportAll(c echo.Context) error {
	yearParam := c.FormValue("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		RenderError(c, http.StatusInternalServerError, err.Error())
	}

	s, err := h.krank.ExportAll(c.Request().Context(), year)
	if err != nil {
		RenderError(c, http.StatusInternalServerError, err.Error())
	}

	filename := fmt.Sprintf("chrono-krankheitstage-export-%v.csv", year)
	c.Response().Header().Set(echo.HeaderContentType, "text/csv; charset=utf-8")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+filename)
	return c.Blob(http.StatusOK, "text/csv; charset=utf-8", []byte(s))
}
