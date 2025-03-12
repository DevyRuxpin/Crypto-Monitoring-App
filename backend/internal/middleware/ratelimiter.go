package middleware

import (
    "context"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
)

type RateLimiterConfig struct {
    RequestsPerSecond int
    BurstSize        int
    ExpirationTime   time.Duration
}

func RateLimiter(client *redis.Client, config RateLimiterConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := "ratelimit:" + c.ClientIP()
        ctx := context.Background()

        // Use Redis INCR for counting requests
        current, err := client.Incr(ctx, key).Result()
        if err != nil {
            c.AbortWithStatus(500)
            return
        }

        // Set expiration on first request
        if current == 1 {
            client.Expire(ctx, key, time.Second)
        }

        if current > int64(config.RequestsPerSecond) {
            c.AbortWithStatusJSON(429, gin.H{
                "error": "Too many requests",
                "retry_after": 1,
            })
            return
        }

        c.Next()
    }
}

func NewRateLimiterConfig() RateLimiterConfig {
    return RateLimiterConfig{
        RequestsPerSecond: 10,
        BurstSize:        20,
        ExpirationTime:   time.Second,
    }
}