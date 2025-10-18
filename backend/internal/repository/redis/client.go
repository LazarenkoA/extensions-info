package redis

import (
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"log"
)

const (
	channel = "channel-logs"
)

type Redis struct {
	client *redis.Client
}

func New(ctx context.Context, host string) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: host,
		DB:   0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "regis is not available")
	}

	return &Redis{client: rdb}, nil
}

func (r *Redis) Subscribe(ctx context.Context, channel string, f func(msg string)) error {
	pubsub := r.client.Subscribe(ctx, channel)
	defer pubsub.Close()

	go func() {
		<-ctx.Done() // при закрытии канала закрывем подписку (что б закрылся канал)
		_ = pubsub.Close()
	}()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		return errors.Wrap(err, "subscribe error")
	}

	ch := pubsub.Channel()

	log.Printf("subscribe for channel: %s", channel)
	for msg := range ch {
		f(msg.Payload)
	}

	return nil
}

func (r *Redis) LogPush(ctx context.Context, msg string) {
	r.client.Publish(ctx, channel, msg)
}

func (r *Redis) LogSubscribe(ctx context.Context, f func(msg string)) {
	go func() {
		err := r.Subscribe(ctx, channel, f)
		if err != nil {
			log.Printf("LogSubscribe error: %v", err)
		}
	}()
}
