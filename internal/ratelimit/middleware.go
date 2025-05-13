package ratelimit

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	limiter "github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	limiterRedis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// NewRateLimiter returns a Gin middleware enforcing the RATE_LIMIT env var.
func NewRateLimiter() (gin.HandlerFunc, error) {
	// Parse rate, e.g. "1000-H"
	rate, err := limiter.NewRateFromFormatted(os.Getenv("RATE_LIMIT"))
	if err != nil {
		return nil, err
	}

	// Init Redis client v9
	dbNum := 0
	if v := os.Getenv("REDIS_DB"); v != "" {
		dbNum, _ = strconv.Atoi(v)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNum,
	})

	// Use default Redis store (prefix, retries, cleanup all come from limiter defaults)
	store, err := limiterRedis.NewStore(client)
	if err != nil {
		return nil, err
	}

	// Create limiter instance
	instance := limiter.New(store, rate)

	// Return Gin middleware
	return ginmiddleware.NewMiddleware(instance), nil
}
