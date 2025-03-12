package metrics

import (
    "time"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
    tradingVolume    *prometheus.CounterVec
    portfolioValue   *prometheus.GaugeVec
    apiLatency      *prometheus.HistogramVec
    errorRate       *prometheus.CounterVec
    cacheHitRate    *prometheus.CounterVec
}

func NewMetrics() *Metrics {
    return &Metrics{
        tradingVolume: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "trading_volume_total",
                Help: "Total trading volume by currency",
            },
            []string{"symbol"},
        ),
        portfolioValue: promauto.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "portfolio_current_value",
                Help: "Current portfolio value by user",
            },
            []string{"user_id"},
        ),
        apiLatency: promauto.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "api_request_duration_seconds",
                Help:    "API request latency distribution",
                Buckets: prometheus.DefBuckets,
            },
            []string{"endpoint", "method"},
        ),
        errorRate: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "error_total",
                Help: "Total number of errors by type",
            },
            []string{"type"},
        ),
        cacheHitRate: promauto.NewCounterVec(
            prometheus.CounterOpts{
                Name: "cache_hit_total",
                Help: "Cache hit/miss counter",
            },
            []string{"result"},
        ),
    }
}

func (m *Metrics) RecordTrading(symbol string, amount float64) {
    m.tradingVolume.WithLabelValues(symbol).Add(amount)
}

func (m *Metrics) UpdatePortfolioValue(userID string, value float64) {
    m.portfolioValue.WithLabelValues(userID).Set(value)
}

func (m *Metrics) RecordLatency(endpoint, method string, duration time.Duration) {
    m.apiLatency.WithLabelValues(endpoint, method).Observe(duration.Seconds())
}
