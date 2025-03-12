package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/yourusername/crypto-monitor/internal/api"
    "github.com/yourusername/crypto-monitor/internal/database"
    "github.com/yourusername/crypto-monitor/internal/cache"
    "github.com/yourusername/crypto-monitor/internal/market"
    "github.com/yourusername/crypto-monitor/internal/websocket"
)

func main() {
    logger := log.New(os.Stdout, "", log.LstdFlags)

    // Load configuration
    config, err := LoadConfig()
    if err != nil {
        logger.Fatalf("Failed to load configuration: %v", err)
    }

    // Initialize dependencies
    db, err := database.NewDB(config.Database)
    if err != nil {
        logger.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    cache, err := cache.NewRedisCache(config.Redis)
    if err != nil {
        logger.Fatalf("Failed to connect to Redis: %v", err)
    }
    defer cache.Close()

    // Initialize services
    marketService := market.NewService(cache)
    hub := websocket.NewHub()
    go hub.Run()

    // Create and configure the server
    server := api.NewServer(
        db,
        marketService,
        hub,
        logger,
    )

    // Start the server
    srv := &http.Server{
        Addr:    config.Server.Address,
        Handler: server.Router,
    }

    // Graceful shutdown
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatalf("Server forced to shutdown: %v", err)
    }

    logger.Println("Server exiting")
}