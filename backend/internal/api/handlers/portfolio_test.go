package handlers

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/yourusername/crypto-monitor/internal/models"
    "github.com/yourusername/crypto-monitor/internal/services/portfolio"
)

func TestGetPortfolio(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    r := gin.New()
    mockService := portfolio.NewMockService()
    handler := NewPortfolioHandler(mockService)

    r.GET("/portfolio", func(c *gin.Context) {
        c.Set("userID", uint(1))
        handler.GetPortfolio(c)
    })

    // Test cases
    tests := []struct {
        name           string
        expectedCode   int
        expectedItems  int
        setupMock      func()
    }{
        {
            name:          "successful retrieval",
            expectedCode:  http.StatusOK,
            expectedItems: 2,
            setupMock: func() {
                mockService.On("GetPortfolio", uint(1)).Return([]models.Portfolio{
                    {Symbol: "BTC", Quantity: 1.0},
                    {Symbol: "ETH", Quantity: 5.0},
                }, nil)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.setupMock()
            
            w := httptest.NewRecorder()
            req, _ := http.NewRequest(http.MethodGet, "/portfolio", nil)
            r.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedCode, w.Code)
            
            if tt.expectedCode == http.StatusOK {
                var response []models.Portfolio
                err := json.Unmarshal(w.Body.Bytes(), &response)
                assert.NoError(t, err)
                assert.Len(t, response, tt.expectedItems)
            }
        })
    }
}