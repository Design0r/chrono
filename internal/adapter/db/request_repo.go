package db

import (
	"chrono/db/repo"
	"chrono/internal/domain"
	"log/slog"
)

type SQLRequestRepo struct {
	r   *repo.Queries
	log *slog.Logger
}

func NewSQLRequestRepo(r *repo.Queries, log *slog.Logger) SQLRequestRepo {
	return SQLRequestRepo{r: r, log: log}
}

func (r *SQLRequestRepo) Create(msg string, user *domain.User, event *domain.Event) (*domain.Request, error)
