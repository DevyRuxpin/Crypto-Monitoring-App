// backend/cmd/server/main.go

package main

import (
    "log"
    "your-project/internal/config"
    "your-project/internal/database"
    "your-project/internal/cache"
    "your-project/internal/api"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize database
    db, err := database.InitDB(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize Redis
    redis, err := cache.InitRedis(cfg.RedisURL)
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    defer redis.Close()

    // Initialize router
    router := api.SetupRouter(db, redis, cfg)

    // Start server
    log.Printf("Server starting on port %s", cfg.Port)
    if err := router.Run(":" + cfg.Port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
