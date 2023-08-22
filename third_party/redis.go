package thirdparty

import (
	"context"
	"time"

	"github.com/akshitbansal-1/async-testing/be/config"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type CacheClient interface {
	Get(key string) (string, error)
	Set(key string, val []byte, exp time.Duration) error
	Delete(key string) error
	Reset() error
	Close() error
}

type cacheClient struct {
	client *redis.Client
}

// Close implements CacheClient.
func (client *cacheClient) Close() error {
	return client.client.Close()
}

// Delete implements CacheClient.
func (client *cacheClient) Delete(key string) error {
	return client.client.Del(context.Background(), key).Err()
}

// Get implements CacheClient.
func (client *cacheClient) Get(key string) (string, error) {
	return client.client.Get(context.Background(), key).Result()
}

// Reset implements CacheClient.
func (client *cacheClient) Reset() error {
	return nil
}

// Set implements CacheClient.
func (client *cacheClient) Set(key string, val []byte, exp time.Duration) error {
	return client.client.Set(context.Background(), key, val, exp).Err()
}

func NewCacheClient(redisConfig config.RedisConfiguration) CacheClient {
	return &cacheClient{
		redis.NewClient(&redis.Options{
			Addr:     redisConfig.Hosts,
			Password: redisConfig.Password,
			DB:       0,
		}),
	}
}
