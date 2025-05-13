package config

import (
	"os"
)

// Config holds service configuration
type Config struct {
	Port     string
	RedisURL string
}

// LoadConfig reads environment variables
func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	redisURL := os.Getenv("REDIS_ADDR")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}
	return Config{Port: port, RedisURL: redisURL}
}
