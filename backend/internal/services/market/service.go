package market

import (
    "context"
    "sync"
    "time"
    "github.com/yourusername/crypto-monitor/pkg/crypto"
)

type Service struct {
    binanceAPI  *crypto.BinanceAPI
    priceCache  map[string]float64
    cacheMutex  sync.RWMutex
    subscribers map[string][]chan float64
}

func NewService(binanceAPI *crypto.BinanceAPI) *Service {
    s := &Service{
        binanceAPI:  binanceAPI,
        priceCache:  make(map[string]float64),
        subscribers: make(map[string][]chan float64),
    }
    go s.startPriceUpdates()
    return s
}

func (s *Service) GetPrice(symbol string) (float64, error) {
    s.cacheMutex.RLock()
    if price, ok := s.priceCache[symbol]; ok {
        s.cacheMutex.RUnlock()
        return price, nil
    }
    s.cacheMutex.RUnlock()

    price, err := s.binanceAPI.GetPrice(symbol)
    if err != nil {
        return 0, err
    }

    s.cacheMutex.Lock()
    s.priceCache[symbol] = price
    s.cacheMutex.Unlock()

    return price, nil
}

func (s *Service) Subscribe(symbol string) chan float64 {
    ch := make(chan float64, 1)
    s.cacheMutex.Lock()
    s.subscribers[symbol] = append(s.subscribers[symbol], ch)
    s.cacheMutex.Unlock()
    return ch
}

func (s *Service) startPriceUpdates() {
    ticker := time.NewTicker(5 * time.Second)
    for range ticker.C {
        s.cacheMutex.RLock()
        symbols := make([]string, 0, len(s.subscribers))
        for symbol := range s.subscribers {
            symbols = append(symbols, symbol)
        }
        s.cacheMutex.RUnlock()

        for _, symbol := range symbols {
            price, err := s.binanceAPI.GetPrice(symbol)
            if err != nil {
                continue
            }

            s.cacheMutex.Lock()
            s.priceCache[symbol] = price
            subs := s.subscribers[symbol]
            s.cacheMutex.Unlock()

            for _, ch := range subs {
                select {
                case ch <- price:
                default:
                }
            }
        }
    }
}