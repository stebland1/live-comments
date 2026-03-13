package repo

import (
	"database/sql"
	"time"

	"github.com/stebland1/live-comments/internal/config"
)

type PostgresStore struct {
	db      *sql.DB
	timeout time.Duration
}

func NewPostgresStore(cfg config.Config) (*PostgresStore, error) {
	db, err := sql.Open("postgres", cfg.PostgresDSN())

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db, timeout: cfg.Postgres.Timeout}, nil
}
