package thirdparty

import (
	"context"
	"fmt"

	"github.com/akshitbansal-1/async-testing/be/config"
	"github.com/redis/go-redis/v9"
)

type Response struct {
	Err error
	Msg *string
}

type PubSub interface {
	SubscribeTopic(ctx context.Context, topic string) chan Response
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

func (p *pubsub) SubscribeTopic(ctx context.Context, topic string) chan Response {
	subscriber := p.client.Subscribe(ctx, topic)
	responseChan := make(chan Response)
	go func() {
		pubsubChannel := subscriber.Channel()
		for {
			handlePubSubEvent(pubsubChannel, responseChan)
		}
	}()

	return responseChan
}

func handlePubSubEvent(channel <-chan *redis.Message, respChan chan Response) {
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Timeline exceeded")
			respChan <- Response{
				ctx.Err(),
				nil,
			}
			close(respChan)
		}
		return
	case msg := <-channel:
		if msg.Payload == "close" {
			ctx.Done()
			close(respChan)
			return
		}
		respChan <- Response{
			nil,
			&msg.Payload,
		}
	}
}
