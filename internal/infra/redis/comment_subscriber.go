package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stebland1/live-comments/internal/config"
)

type RedisClient interface {
	Subscribe(ctx context.Context, channels ...string) PubSub
}

type PubSub interface {
	Channel(...redis.ChannelOption) <-chan *redis.Message
	Close() error
}

type CommentSubscriber struct {
	client  RedisClient
	timeout time.Duration
}

type redisClientAdapter struct {
	client *redis.Client
}

func (r *redisClientAdapter) Subscribe(ctx context.Context, channels ...string) PubSub {
	return r.client.Subscribe(ctx, channels...)
}

func NewCommentSubscriber(cfg config.Config) *CommentSubscriber {
	return &CommentSubscriber{
		client: &redisClientAdapter{client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		})},
		timeout: cfg.Redis.Timeout,
	}
}

func (cs *CommentSubscriber) Subscribe(ctx context.Context, videoID int64) <-chan string {
	// TODO: centralise this in one place. GetChannel() -> string (for example)
	channel := fmt.Sprintf("comment:%d", videoID)

	subCtx, cancel := context.WithTimeout(ctx, cs.timeout)
	defer cancel()

	pubsub := cs.client.Subscribe(subCtx, channel)

	out := make(chan string)
	psChan := pubsub.Channel()

	go func() {
		defer close(out)
		defer pubsub.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-psChan:
				if !ok {
					return
				}

				select {
				case out <- msg.Payload:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}
