package integration

import (
    "context"
    "os"
    "testing"
    "github.com/yourusername/crypto-monitor/internal/config"
    "github.com/yourusername/crypto-monitor/internal/database"
)

type TestEnv struct {
    DB     *database.DB
    Config *config.Config
}

func SetupTestEnv(t *testing.T) *TestEnv {
    t.Helper()

    cfg, err := config.LoadConfig()
    if err != nil {
        t.Fatalf("Failed to load config: %v", err)
    }

    // Use test database
    cfg.DatabaseURL = os.Getenv("TEST_DATABASE_URL")

    db, err := database.InitDB(cfg.DatabaseURL)
    if err != nil {
        t.Fatalf("Failed to init test DB: %v", err)
    }

    return &TestEnv{
        DB:     db,
        Config: cfg,
    }
}

func (env *TestEnv) Cleanup(t *testing.T) {
    t.Helper()
    
    // Clean up test data
    if err := env.DB.Exec(context.Background(), "TRUNCATE users, portfolios, alerts CASCADE").Error; err != nil {
        t.Errorf("Failed to cleanup test DB: %v", err)
    }
}