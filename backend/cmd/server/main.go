package main

import (
    "log"
    "os"
    "github.com/yourusername/crypto-monitor/internal/api"
    "github.com/yourusername/crypto-monitor/internal/database"
    "github.com/yourusername/crypto-monitor/internal/cache"
    "github.com/yourusername/crypto-monitor/internal/market"
    "github.com/yourusername/crypto-monitor/internal/portfolio"
    "github.com/yourusername/crypto-monitor/internal/websocket"
)

func main() {
    // Initialize logger
    logger := log.New(os.Stdout, "[CRYPTO-MONITOR] ", log.LstdFlags)

    // Load configuration
    config, err := LoadConfig()
    if err != nil {
        logger.Fatalf("Failed to load config: %v", err)
    }

    // Initialize database
    db, err := database.NewDB(config.DatabaseURL)
    if err != nil {
        logger.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Initialize Redis cache
    redisCache, err := cache.NewRedisCache(config.RedisURL)
    if err != nil {
        logger.Fatalf("Failed to connect to Redis: %v", err)
    }
    defer redisCache.Close()

    // Initialize services
    marketService := market.NewMarketService(redisCache)
    portfolioService := portfolio.NewService(db, marketService)
    
    // Initialize WebSocket hub
    hub := websocket.NewHub()
    go hub.Run()

    // Initialize API server
    server := api.NewServer(
        db,
        marketService,
        portfolioService,
        hub,
        logger,
    )

    // Start server
    logger.Printf("Starting server on %s", config.ServerAddress)
    if err := server.Start(config.ServerAddress); err != nil {
        logger.Fatalf("Server failed: %v", err)
    }
}

type Config struct {
    ServerAddress string
    DatabaseURL   string
    RedisURL      string
    JWTSecret     string
}

func LoadConfig() (*Config, error) {
    return &Config{
        ServerAddress: os.Getenv("SERVER_ADDRESS"),
        DatabaseURL:   os.Getenv("DATABASE_URL"),
        RedisURL:      os.Getenv("REDIS_URL"),
        JWTSecret:     os.Getenv("JWT_SECRET"),
    }, nil
}