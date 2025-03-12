package market

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/yourusername/crypto-monitor/internal/cache"
)

type MarketService struct {
    cache  *cache.RedisCache
    client *http.Client
}

type PriceData struct {
    Symbol string  `json:"symbol"`
    Price  float64 `json:"price"`
    Time   int64   `json:"timestamp"`
}

func NewMarketService(cache *cache.RedisCache) *MarketService {
    return &MarketService{
        cache: cache,
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (s *MarketService) GetPrice(symbol string) (*PriceData, error) {
    // Try cache first
    var priceData PriceData
    cacheKey := fmt.Sprintf("price:%s", symbol)
    
    err := s.cache.Get(cacheKey, &priceData)
    if err == nil {
        return &priceData, nil
    }

    // Fetch from external API
    url := fmt.Sprintf("https://api.example.com/v1/prices/%s", symbol)
    resp, err := s.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
    }

    if err := json.NewDecoder(resp.Body).Decode(&priceData); err != nil {
        return nil, err
    }

    // Cache the result
    s.cache.Set(cacheKey, priceData, 30*time.Second)

    return &priceData, nil
}
