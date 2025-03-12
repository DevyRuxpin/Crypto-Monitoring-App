package tests

import (
    "testing"
    "context"
    "github.com/stretchr/testify/assert"
    "github.com/yourusername/crypto-monitor/internal/market"
)

func TestMarketService(t *testing.T) {
    cache, err := setupTestCache()
    assert.NoError(t, err)

    service := market.NewMarketService(cache)

    t.Run("GetPrice returns correct data", func(t *testing.T) {
        price, err := service.GetPrice("BTC")
        assert.NoError(t, err)
        assert.NotNil(t, price)
        assert.Equal(t, "BTC", price.Symbol)
        assert.Greater(t, price.Price, 0.0)
    })

    t.Run("Cache is working correctly", func(t *testing.T) {
        // First call should hit the API
        price1, err := service.GetPrice("ETH")
        assert.NoError(t, err)

        // Second call should hit the cache
        price2, err := service.GetPrice("ETH")
        assert.NoError(t, err)

        assert.Equal(t, price1.Price, price2.Price)
    })
}
