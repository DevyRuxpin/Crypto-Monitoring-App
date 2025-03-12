package integration

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "encoding/json"
    "github.com/stretchr/testify/assert"
    "github.com/yourusername/crypto-monitor/internal/api"
    "github.com/yourusername/crypto-monitor/internal/market"
)

func TestMarketIntegration(t *testing.T) {
    // Setup test environment
    db := setupTestDB(t)
    cache := setupTestCache(t)
    marketService := market.NewService(cache)
    server := api.NewServer(db, marketService)

    t.Run("Get Market Data", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/market/BTC", nil)
        w := httptest.NewRecorder()
        server.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

        var response market.MarketData
        err := json.NewDecoder(w.Body).Decode(&response)
        assert.NoError(t, err)
        assert.Equal(t, "BTC", response.Symbol)
        assert.Greater(t, response.Price, 0.0)
    })
}
