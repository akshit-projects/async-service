package pubsub

import (
	"context"

	"github.com/akshitbansal-1/async-testing/worker/config"
	"github.com/redis/go-redis/v9"
)

type PubSub interface {
	PublishMessage(ctx context.Context, topic string, msg string) error
}

type pubsub struct {
	client *redis.Client
}

func NewPubSubClient(config *config.Configuration) PubSub {
	redisConfig := config.RedisConfiguration
	return &pubsub{
		redis.NewClient(&redis.Options{
			Addr:     redisConfig.Hosts,
			Password: redisConfig.Password,
			DB:       0,
		}),
	}
}

func (p *pubsub) PublishMessage(ctx context.Context, topic string, msg string) error {
	cmd := p.client.Publish(ctx, topic, msg)
	return cmd.Err()
}
