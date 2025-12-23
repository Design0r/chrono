package service

import (
	"context"
	"log/slog"
	"time"

	"chrono/internal/domain"
)

type TimestampsService struct {
	timestamps domain.TimestampsRepository
	log        *slog.Logger
}

func NewTimestampsService(r domain.TimestampsRepository, log *slog.Logger) TimestampsService {
	return TimestampsService{timestamps: r, log: log}
}

func (r *TimestampsService) GetById(ctx context.Context, id int64) (domain.Timestamp, error) {
	return r.timestamps.GetById(ctx, id)
}

func (r *TimestampsService) Start(ctx context.Context, userId int64) (domain.Timestamp, error) {
	return r.timestamps.Start(ctx, userId)
}

func (r *TimestampsService) Stop(ctx context.Context, id int64) (domain.Timestamp, error) {
	return r.timestamps.Stop(ctx, id)
}

func (r *TimestampsService) Delete(ctx context.Context, id int64) error {
	return r.timestamps.Delete(ctx, id)
}

func (r *TimestampsService) GetInRange(
	ctx context.Context,
	userId int64,
	start time.Time,
	stop time.Time,
) ([]domain.Timestamp, error) {
	return r.timestamps.GetInRange(ctx, userId, start, stop)
}

func (r *TimestampsService) GetForToday(
	ctx context.Context,
	userId int64,
) ([]domain.Timestamp, error) {
	now := time.Now().UTC()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	stop := start.AddDate(0, 0, 1)

	return r.timestamps.GetInRange(ctx, userId, start, stop)
}

func (r *TimestampsService) GetTotalSecondsInRange(
	ctx context.Context,
	userId int64,
	start time.Time,
	stop time.Time,
) (float64, error) {
	return r.timestamps.GetTotalSecondsInRange(ctx, userId, start, stop)
}
