package postgres

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *CommentRepo) CreateComment(ctx context.Context, userID int64, videoID int64, content string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO comments (user_id, video_id, content) VALUES ($1, $2, $3) RETURNING id",
		userID,
		videoID,
		content,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("inserting into comments table: %w", err)
	}

	return id, nil
}
