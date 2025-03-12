package cache

import (
    "context"
    "encoding/json"
    "time"
    "github.com/go-redis/redis/v8"
)

type RedisCache struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisCache(redisURL string) (*RedisCache, error) {
    opts, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, err
    }

    client := redis.NewClient(opts)
    ctx := context.Background()

    // Test connection
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, err
    }

    return &RedisCache{
        client: client,
        ctx:    ctx,
    }, nil
}

func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
    json, err := json.Marshal(value)
    if err != nil {
        return err
    }

    return c.client.Set(c.ctx, key, json, expiration).Err()
}

func (c *RedisCache) Get(key string, dest interface{}) error {
    val, err := c.client.Get(c.ctx, key).Result()
    if err != nil {
        return err
    }

    return json.Unmarshal([]byte(val), dest)
}

func (c *RedisCache) Delete(key string) error {
    return c.client.Del(c.ctx, key).Err()
}

func (c *RedisCache) Close() error {
    return c.client.Close()
}