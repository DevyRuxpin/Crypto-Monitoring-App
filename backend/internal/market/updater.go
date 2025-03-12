package market

import (
    "context"
    "time"
    "log"
)

type Updater struct {
    service     *Service
    hub         *websocket.Hub
    symbols     []string
    interval    time.Duration
    logger      *log.Logger
}

func NewUpdater(service *Service, hub *websocket.Hub, symbols []string, interval time.Duration, logger *log.Logger) *Updater {
    return &Updater{
        service:  service,
        hub:      hub,
        symbols:  symbols,
        interval: interval,
        logger:   logger,
    }
}

func (u *Updater) Start(ctx context.Context) {
    ticker := time.NewTicker(u.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            u.updateMarketData()
        }
    }
}

func (u *Updater) updateMarketData() {
    for _, symbol := range u.symbols {
        data, err := u.service.GetMarketData(symbol)
        if err != nil {
            u.logger.Printf("Failed to update market data for %s: %v", symbol, err)
            continue
        }

        message := websocket.Message{
            Type:    websocket.MessageTypeMarketData,
            Payload: data,
        }

        u.hub.Broadcast(message)
    }
}