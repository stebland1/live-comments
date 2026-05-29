package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type fakePubSub struct {
	ch chan *redis.Message
}

func (ps *fakePubSub) Channel(...redis.ChannelOption) <-chan *redis.Message {
	return ps.ch
}

func (ps *fakePubSub) Close() error {
	close(ps.ch)
	return nil
}

type fakeClient struct {
	ps *fakePubSub
}

func (f *fakeClient) Subscribe(ctx context.Context, channels ...string) PubSub {
	return f.ps
}

func TestSubscribe_ForwardsMessages(t *testing.T) {
	ctx := t.Context()

	ch := make(chan *redis.Message)
	ps := &fakePubSub{ch: ch}

	client := &fakeClient{ps: ps}
	s := &CommentSubscriber{
		client:  client,
		timeout: time.Hour,
	}

	message := "Hello"
	videoID := "vid123"

	out := s.Subscribe(ctx, videoID)
	ch <- &redis.Message{
		Payload: message,
	}

	select {
	case got := <-out:
		if got != message {
			t.Fatalf("expected %s, got %s", message, got)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestSubscribe_CancelsCleanly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan *redis.Message)
	ps := &fakePubSub{ch: ch}
	client := &fakeClient{ps: ps}
	s := &CommentSubscriber{
		client:  client,
		timeout: time.Second,
	}

	videoID := "vid123"
	out := s.Subscribe(ctx, videoID)

	cancel()

	select {
	case _, ok := <-out:
		if ok {
			t.Fatalf("expected channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for channel close")
	}
}
