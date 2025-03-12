// backend/internal/cache/redis.go

package cache

import (
    "github.com/go-redis/redis/v8"
    "context"
)

func InitRedis(redisURL string) (*redis.Client, error) {
    opt, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, err
    }

    client := redis.NewClient(opt)
    
    // Test the connection
    ctx := context.Background()
    _, err = client.Ping(ctx).Result()
    if err != nil {
        return nil, err
    }

    return client, nil
}
