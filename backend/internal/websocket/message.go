package websocket

type MessageType string

const (
    MessageTypeSubscribe   MessageType = "SUBSCRIBE"
    MessageTypeUnsubscribe MessageType = "UNSUBSCRIBE"
    MessageTypeMarketData  MessageType = "MARKET_DATA"
    MessageTypeError       MessageType = "ERROR"
)

type Message struct {
    Type    MessageType  `json:"type"`
    Payload interface{} `json:"payload"`
}

type SubscriptionPayload struct {
    Symbol string `json:"symbol"`
}

type MarketDataPayload struct {
    Symbol    string  `json:"symbol"`
    Price     float64 `json:"price"`
    Change24h float64 `json:"change24h"`
    Volume24h float64 `json:"volume24h"`
    Timestamp int64   `json:"timestamp"`
}
