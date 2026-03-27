package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stebland1/live-comments/internal/config"
)

type CommentSubscriber struct {
	client  *redis.Client
	timeout time.Duration
}

func NewCommentSubscriber(cfg config.Config) *CommentSubscriber {
	return &CommentSubscriber{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		}),
		timeout: cfg.Redis.Timeout,
	}
}

func (cs *CommentSubscriber) Subscribe(ctx context.Context, videoID string) <-chan string {
	// TODO: centralise this in one place. GetChannel() -> string (for example)
	channel := fmt.Sprintf("comment:%s", videoID)

	subCtx, cancel := context.WithTimeout(ctx, cs.timeout)
	defer cancel()

	pubsub := cs.client.Subscribe(subCtx, channel)
	out := make(chan string)

	go func() {
		defer close(out)
		defer pubsub.Close()

		ch := pubsub.Channel()

		// two levels of indirection is necessary here;
		// to decouple the service that calls this method from the infra.
		// this way we can return a channel receiving strings
		// instead of returning a channel receiving *redis.Message
		for {
			select {
			case msg := <-ch:
				select {
				case out <- msg.Payload:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
