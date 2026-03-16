package comment

import (
	"context"
	"fmt"
)

type Repository interface {
	CreateComment(ctx context.Context, userID int64, videoID int64, content string) (int64, error)
}

type Publisher interface {
	PublishComment(ctx context.Context, comment Comment) error
}

type Service struct {
	repo      Repository
	publisher Publisher
}

func NewService(repo Repository, publisher Publisher) *Service {
	return &Service{repo: repo, publisher: publisher}
}

func (s *Service) CreateComment(ctx context.Context, userID int64, videoID int64, content string) (int64, error) {
	commentID, err := s.repo.CreateComment(ctx, userID, videoID, content)
	if err != nil {
		return 0, fmt.Errorf("creating comment: %w", err)
	}

	comment := Comment{
		ID:      commentID,
		UserID:  userID,
		VideoID: videoID,
		Content: content,
	}

	err = s.PublishComment(ctx, comment)
	if err != nil {
		return 0, fmt.Errorf("publishing comment: %w", err)
	}

	return commentID, nil
}

func (s *Service) PublishComment(ctx context.Context, comment Comment) error {
	return s.publisher.PublishComment(ctx, comment)
}
