package service

import (
	"chrono/internal/domain"
)

type HolidayService struct {
	user  UserService
	event EventService
	api   domain.ApiCacheRepository
}

func NewHolidayService(user UserService, event EventService, api domain.ApiCacheRepository) *HolidayService {
	return &HolidayService{user: user, event: event, api: api}
}

/* func (h *HolidayService) GetHolidays(ctx context.Context, year int) (domain.Holidays, error) {
	return h.api.GetHolidays(ctx, year)
} */
