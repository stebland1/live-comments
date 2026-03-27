package stream

import (
	"context"
)

type Subscriber interface {
	Subscribe(ctx context.Context, videoID string) <-chan string
}

type Service struct {
	subscriber Subscriber
}

func NewService(subscriber Subscriber) *Service {
	return &Service{subscriber: subscriber}
}

func (s *Service) Subscribe(ctx context.Context, videoID string) <-chan string {
	return s.subscriber.Subscribe(ctx, videoID)
}
