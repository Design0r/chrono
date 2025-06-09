package db

import (
	"chrono/db/repo"
	"context"
	"log/slog"
)

type SQLAPICacheRepo struct {
	r   *repo.Queries
	log *slog.Logger
}

func NewSQLAPICacheRepo(r *repo.Queries, log *slog.Logger) SQLAPICacheRepo {
	return SQLAPICacheRepo{r: r, log: log}
}

func (r *SQLAPICacheRepo) Exists(ctx context.Context, year int64) (int64, error) {
	return r.r.CacheExists(ctx, year)
}

func (r *SQLAPICacheRepo) GetAll(ctx context.Context) ([]int64, error) {
	return r.r.GetApiCacheYears(ctx)
}

func (r *SQLAPICacheRepo) Create(ctx context.Context, year int64) error {
	return r.r.CreateCache(ctx, year)
}
