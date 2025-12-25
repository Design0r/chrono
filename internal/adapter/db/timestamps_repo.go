package db

import (
	"context"
	"log/slog"
	"time"

	"chrono/db/repo"
	"chrono/internal/domain"
)

type SQLTimestampsRepo struct {
	q   *repo.Queries
	log *slog.Logger
}

func NewSQLTimestampsRepo(q *repo.Queries, log *slog.Logger) domain.TimestampsRepository {
	return &SQLTimestampsRepo{q: q, log: log}
}

func (r *SQLTimestampsRepo) GetById(ctx context.Context, id int64) (domain.Timestamp, error) {
	t, err := r.q.GetTimestampById(ctx, id)
	if err != nil {
		r.log.Error("repo.GetTimestampById failed:", slog.String("error", err.Error()))
		return domain.Timestamp{}, err
	}

	return (domain.Timestamp)(t), nil
}

func (r *SQLTimestampsRepo) Start(ctx context.Context, userId int64) (domain.Timestamp, error) {
	t, err := r.q.StartTimestamp(ctx, userId)
	if err != nil {
		r.log.Error("repo.StartTimestamp failed:", slog.String("error", err.Error()))
		return domain.Timestamp{}, err
	}

	return (domain.Timestamp)(t), nil
}

func (r *SQLTimestampsRepo) Stop(ctx context.Context, id int64) (domain.Timestamp, error) {
	t, err := r.q.StopTimestamp(ctx, id)
	if err != nil {
		r.log.Error("repo.StopTimestamp failed:", slog.String("error", err.Error()))
		return domain.Timestamp{}, err
	}

	return (domain.Timestamp)(t), nil
}

func (r *SQLTimestampsRepo) Delete(ctx context.Context, id int64) error {
	err := r.q.DeleteTimestamp(ctx, id)
	if err != nil {
		r.log.Error("repo.DeleteTimestamp failed:", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *SQLTimestampsRepo) GetInRange(
	ctx context.Context,
	userId int64,
	start time.Time,
	stop time.Time,
) ([]domain.Timestamp, error) {
	params := repo.GetTimestampsInRangeParams{UserID: userId, StartTime: &start, EndTime: stop}
	t, err := r.q.GetTimestampsInRange(ctx, params)
	if err != nil {
		r.log.Error("repo.GetTimestampsInRange failed:", slog.String("error", err.Error()))
		return []domain.Timestamp{}, err
	}

	timestamps := make([]domain.Timestamp, len(t))
	for i, x := range t {
		timestamps[i] = (domain.Timestamp)(x)
	}

	return timestamps, nil
}

func (r *SQLTimestampsRepo) GetTotalSecondsInRange(
	ctx context.Context,
	userId int64,
	start time.Time,
	stop time.Time,
) (float64, error) {
	s := 0.0
	params := repo.GetTotalSecondsInRangeParams{UserID: userId, RangeEnd: stop, RangeStart: start}
	seconds, err := r.q.GetTotalSecondsInRange(ctx, params)
	if err != nil {
		r.log.Error("repo.GetTimestampsInRange failed:", slog.String("error", err.Error()))
		return s, err
	}

	if seconds != nil {
		s = *seconds
	}

	return s, nil
}

func (r *SQLTimestampsRepo) GetLatest(ctx context.Context, userId int64) (domain.Timestamp, error) {
	t, err := r.q.GetLatestTimestamp(ctx, userId)
	if err != nil {
		r.log.Error("repo.GetLatestTimestamp failed:", slog.String("error", err.Error()))
		return domain.Timestamp{}, err
	}

	return (domain.Timestamp)(t), nil
}

func (r *SQLTimestampsRepo) Update(
	ctx context.Context,
	ts *domain.Timestamp,
) (domain.Timestamp, error) {
	params := repo.UpdateTimestampParams{ID: ts.ID, StartTime: ts.StartTime, EndTime: ts.EndTime}
	t, err := r.q.UpdateTimestamp(ctx, params)
	if err != nil {
		r.log.Error("repo.UpdateTimestamp failed:", slog.String("error", err.Error()))
		return domain.Timestamp{}, err
	}

	return (domain.Timestamp)(t), nil
}

func (r *SQLTimestampsRepo) GetAllForUser(
	ctx context.Context,
	userId int64,
) ([]domain.Timestamp, error) {
	t, err := r.q.GetAllTimestampsForUser(ctx, userId)
	if err != nil {
		r.log.Error("repo.GetAllTimestampsForUser failed:", slog.String("error", err.Error()))
		return []domain.Timestamp{}, err
	}

	timestamps := make([]domain.Timestamp, len(t))
	for i, x := range t {
		timestamps[i] = (domain.Timestamp)(x)
	}

	return timestamps, nil
}
