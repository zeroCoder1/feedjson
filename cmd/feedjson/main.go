package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/zeroCoder1/feedjson/internal/api"
	"github.com/zeroCoder1/feedjson/internal/cache"
	"github.com/zeroCoder1/feedjson/internal/config"
	"github.com/zeroCoder1/feedjson/internal/ratelimit"
)

func main() {
	// 1. Load configuration from environment
	cfg := config.LoadConfig()

	// 2. Initialize Redis (for cache & rate-limiter store)
	cache.InitRedis(cfg.RedisURL, "", 0)
	rdb := cache.GetClient()

	if os.Getenv("ADMIN_SECRET") == "" {
		fmt.Fprintln(os.Stderr, "⚠️  Warning: ADMIN_SECRET is empty; anyone can issue tokens")
	}

	// 3. Build rate-limit middleware
	rlMiddleware, err := ratelimit.NewRateLimiter()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ failed to init rate limiter: %v\n", err)
		os.Exit(1)
	}

	// 4. Set up Gin router
	r := gin.Default()

	// Serve all files in ./web at /web (e.g. /web/index.html, /web/docs.html)
	r.Static("/docs", "./docs")

	// 4a. Enforce bearer token auth
	// authMiddleware := auth.RequireAuth()
	// r.Use(authMiddleware)

	// 4b. Then rate-limit
	r.Use(rlMiddleware)

	// 5. Register API endpoints
	api.RegisterRoutes(r, rdb)

	// 6. Start listening
	addr := fmt.Sprintf(":%s", cfg.Port)
	if err := r.Run(addr); err != nil {
		fmt.Fprintf(os.Stderr, "❌ server error: %v\n", err)
		os.Exit(1)
	}
}
