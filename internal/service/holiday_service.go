package service

import (
	"chrono/internal/domain"
	"chrono/internal/service"
	"context"
)

type HolidayService struct {
	user  service.UserService
	event service.EventService
	api   domain.ApiCacheRepository
}

func NewHolidayService(user service.UserService, event service.EventService, api domain.ApiCacheRepository) *HolidayService {
	return &HolidayService{user: user, event: event, api: api}
}

func (h *HolidayService) GetHolidays(ctx context.Context, year int) ([]domain.Holiday, error) {
	return h.api.GetHolidays(ctx, year)
}
