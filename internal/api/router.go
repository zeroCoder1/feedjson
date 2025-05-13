package api

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/google/uuid"
	"github.com/zeroCoder1/feedjson/internal/auth"
	"github.com/zeroCoder1/feedjson/internal/model"
	"github.com/zeroCoder1/feedjson/internal/parser"
)

// RegisterRoutes sets up all API endpoints
func RegisterRoutes(r *gin.Engine, rdb *redis.Client) {
	v1 := r.Group("/v1")
	{
		// Token management
		v1.POST("/tokens", adminAuth(), func(c *gin.Context) {
			// generate
			tok := uuid.NewString()
			if err := auth.AddToken(context.Background(), tok); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save token"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"token": tok})
		})

		// Feed endpoint stays protected by RequireAuth + rate-limit
		feeds := v1.Group("/feed")
		feeds.Use(auth.RequireAuth())
		{
			feeds.GET("", apiFeedHandler(rdb)) // your existing handler
		}
	}
}

// apiFeedHandler returns our /v1/feed handler
func apiFeedHandler(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse query params
		rssURL := c.Query("rss_url")
		if rssURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "rss_url query param required"})
			return
		}

		countStr := c.DefaultQuery("count", "0")
		count, err := strconv.Atoi(countStr)
		if err != nil || count < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid count param"})
			return
		}

		ctx := c.Request.Context()

		// Cache lookup
		key := "feed:" + hash(rssURL+"|"+countStr)
		if data, err := rdb.Get(ctx, key).Result(); err == nil {
			var resp model.FeedResponse
			if err := json.Unmarshal([]byte(data), &resp); err == nil {
				c.JSON(http.StatusOK, resp)
				return
			}
		}

		// Fetch and parse fresh feed
		respPtr, err := parser.FetchFeed(ctx, rssURL, count)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		// Cache the result
		bytes, _ := json.Marshal(respPtr)
		rdb.Set(ctx, key, bytes, 15*time.Minute)

		// Return JSON
		c.JSON(http.StatusOK, respPtr)
	}
}

// Admin-only middleware (checks ADMIN_SECRET)
func adminAuth() gin.HandlerFunc {
	secret := os.Getenv("ADMIN_SECRET")
	return func(c *gin.Context) {
		if c.GetHeader("X-Admin-Token") != secret {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin creds required"})
			return
		}
		c.Next()
	}
}

// hash returns an MD5 hex string for a given input
func hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
