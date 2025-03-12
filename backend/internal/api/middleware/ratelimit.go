package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
    "sync"
)

type RateLimiter struct {
    ips map[string]*rate.Limiter
    mu  *sync.RWMutex
}

func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        ips: make(map[string]*rate.Limiter),
        mu:  &sync.RWMutex{},
    }
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    limiter, exists := rl.ips[ip]
    if !exists {
        limiter = rate.NewLimiter(rate.Every(1*time.Second), 10)
        rl.ips[ip] = limiter
    }

    return limiter
}

func RateLimit(rl *RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        limiter := rl.getLimiter(ip)
        if !limiter.Allow() {
            c.JSON(429, gin.H{
                "error": "Too many requests",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}