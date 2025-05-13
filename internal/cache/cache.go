package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// InitRedis initializes the global Redis client.
// Call this once (e.g. in cmd/feedjson/main.go) before using Get/Set or RDB.
func InitRedis(addr, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

// Get retrieves a string value by key.
// Returns redis.Nil error if key does not exist.
func Get(ctx context.Context, key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

// Set stores a string value with the given TTL.
func Set(ctx context.Context, key, val string, ttl time.Duration) error {
	return rdb.Set(ctx, key, val, ttl).Err()
}

// RDB returns the underlying Redis client for direct use
// (e.g. in auth/store.go or advanced operations).
func RDB() *redis.Client {
	return rdb
}

// GetClient returns the global Redis client.
func GetClient() *redis.Client {
	return rdb
}
