package auth

import (
	"context"

	"github.com/zeroCoder1/feedjson/internal/cache"
)

const tokenSet = "feedjson:tokens"

// AddToken saves a new token (no expiry) to Redis
func AddToken(ctx context.Context, token string) error {
	return cache.RDB().SAdd(ctx, tokenSet, token).Err()
}

// IsValidToken checks whether the given token exists
func IsValidToken(ctx context.Context, token string) (bool, error) {
	return cache.RDB().SIsMember(ctx, tokenSet, token).Result()
}
