package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/internal/service"
)

type APIExportHandler struct {
	krank service.ExportService
}

func NewAPIExportHandler(k service.ExportService) APIExportHandler {
	return APIExportHandler{krank: k}
}

func (s *APIExportHandler) RegisterRoutes(group *echo.Group) {
	g := group.Group("/export")
	g.GET("/:year", s.ExportYear)
}

func (h *APIExportHandler) ExportYear(c echo.Context) error {
	yearParam := c.Param("year")
	year, err := strconv.Atoi(yearParam)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	s, err := h.krank.ExportAll(c.Request().Context(), year)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	filename := fmt.Sprintf("chrono-krankheitstage-export-%v.csv", year)
	c.Response().Header().Set(echo.HeaderContentType, "text/csv; charset=utf-8")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+filename)
	return c.Blob(http.StatusOK, "text/csv; charset=utf-8", []byte(s))
}
