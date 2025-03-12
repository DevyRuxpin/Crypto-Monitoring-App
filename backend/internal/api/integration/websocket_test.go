package integration

import (
    "testing"
    "net/http/httptest"
    "github.com/gorilla/websocket"
    "github.com/stretchr/testify/assert"
    "github.com/yourusername/crypto-monitor/internal/websocket"
)

func TestWebSocketIntegration(t *testing.T) {
    server := httptest.NewServer(nil)
    defer server.Close()

    // Convert http://... to ws://...
    url := "ws" + server.URL[4:] + "/ws"

    // Connect to the server
    ws, _, err := websocket.DefaultDialer.Dial(url, nil)
    assert.NoError(t, err)
    defer ws.Close()

    t.Run("Subscribe to market data", func(t *testing.T) {
        message := websocket.Message{
            Type: websocket.MessageTypeSubscribe,
            Payload: websocket.SubscriptionPayload{
                Symbol: "BTC",
            },
        }

        err := ws.WriteJSON(message)
        assert.NoError(t, err)

        // Read response
        var response websocket.Message
        err = ws.ReadJSON(&response)
        assert.NoError(t, err)
        assert.Equal(t, websocket.MessageTypeMarketData, response.Type)
    })
}