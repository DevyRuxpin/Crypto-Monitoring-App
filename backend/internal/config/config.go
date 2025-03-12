// backend/internal/config/config.go

package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    DatabaseURL        string
    RedisURL          string
    Port              string
    JWTSecret         string
    JWTExpirationHours int
    Environment       string
}

func LoadConfig() *Config {
    jwtExpHours, _ := strconv.Atoi(getEnvWithDefault("JWT_EXPIRATION_HOURS", "24"))
    
    return &Config{
        DatabaseURL:        getEnvWithDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/crypto_monitor"),
        RedisURL:          getEnvWithDefault("REDIS_URL", "redis://localhost:6379"),
        Port:              getEnvWithDefault("PORT", "8080"),
        JWTSecret:         getEnvWithDefault("JWT_SECRET", "your_super_secret_key_change_this_in_production"),
        JWTExpirationHours: jwtExpHours,
        Environment:       getEnvWithDefault("ENVIRONMENT", "development"),
    }
}

func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
