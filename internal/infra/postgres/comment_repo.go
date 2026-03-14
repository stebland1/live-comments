package postgres

import (
	"database/sql"
	"time"

	"github.com/stebland1/live-comments/internal/config"
)

type CommentRepo struct {
	db      *sql.DB
	timeout time.Duration
}

func NewCommentRepo(cfg config.Config) (*CommentRepo, error) {
	db, err := sql.Open("postgres", cfg.PostgresDSN())

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &CommentRepo{db: db, timeout: cfg.Postgres.Timeout}, nil
}
