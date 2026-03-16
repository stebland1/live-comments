package comment

import (
	"context"
	"fmt"
)

type Repository interface {
	CreateComment(ctx context.Context, userID int64, videoID int64, content string) (int64, error)
}
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateComment(ctx context.Context, userID int64, videoID int64, content string) (int64, error) {
	commentID, err := s.repo.CreateComment(ctx, userID, videoID, content)
	if err != nil {
		return 0, fmt.Errorf("creating comment: %w", err)
	}

	// TODO: Publish the comment into a message queue.
	return commentID, nil
}
