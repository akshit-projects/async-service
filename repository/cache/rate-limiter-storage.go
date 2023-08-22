package cache

import (
	"time"

	thirdparty "github.com/akshitbansal-1/async-testing/be/third_party"
)

type RateLimiterStorage struct {
	cacheClient thirdparty.CacheClient
}

func (RateLimiterStorage *RateLimiterStorage) Get(key string) ([]byte, error) {
	val, err := RateLimiterStorage.cacheClient.Get(key)
	return []byte(val), err
}

// Set stores the given value for the given key along
// with an expiration value, 0 means no expiration.
// Empty key or value will be ignored without an error.
func (RateLimiterStorage *RateLimiterStorage) Set(key string, val []byte, exp time.Duration) error {
	return RateLimiterStorage.cacheClient.Set(key, val, exp)
}

// Delete deletes the value for the given key.
// It returns no error if the storage does not contain the key,
func (RateLimiterStorage *RateLimiterStorage) Delete(key string) error {
	return RateLimiterStorage.cacheClient.Delete(key)
}

// Reset resets the storage and delete all keys.
func (*RateLimiterStorage) Reset() error {
	return nil
}

// Close closes the storage and will stop any running garbage
// collectors and open connections.
func (RateLimiterStorage *RateLimiterStorage) Close() error {
	return RateLimiterStorage.cacheClient.Close()
}

func NewCustomRateLimiterStorage(cacheClient thirdparty.CacheClient) *RateLimiterStorage {
	return &RateLimiterStorage{
		cacheClient: cacheClient,
	}
}