package domain

import (
	"context"
	"time"
)

type ApiCache struct {
	ID        int64     `json:"id"`
	Year      int64     `json:"year"`
	CreatedAt time.Time `json:"created_at"`
}

type Holidays = map[string]map[string]string

type ApiCacheRepository interface {
	Exists(ctx context.Context, year int64) (int64, error)
	GetAll(ctx context.Context) ([]int64, error)
	Create(ctx context.Context, year int64) error
}
