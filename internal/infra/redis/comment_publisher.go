package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stebland1/live-comments/internal/comment"
	"github.com/stebland1/live-comments/internal/config"
)

type CommentPublisher struct {
	client  *redis.Client
	timeout time.Duration
}

func NewCommentPublisher(cfg config.Config) *CommentPublisher {
	return &CommentPublisher{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		}),
		timeout: cfg.Redis.Timeout,
	}
}

func (cp *CommentPublisher) PublishComment(ctx context.Context, c comment.Comment) error {
	ctx, cancel := context.WithTimeout(ctx, cp.timeout)
	defer cancel()

	event := comment.CommentEvent{
		ID:      c.ID,
		VideoID: c.VideoID,
		Content: c.Content,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshalling comment: %w", err)
	}

	channel := fmt.Sprintf("comment:%d", c.VideoID)
	err = cp.client.Publish(ctx, channel, payload).Err()
	if err != nil {
		return fmt.Errorf("publishing to channel %s: %w", channel, err)
	}

	return nil
}
