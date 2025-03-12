package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/yourusername/crypto-monitor/internal/market"
)

type MarketHandler struct {
    marketService *market.Service
}

func NewMarketHandler(marketService *market.Service) *MarketHandler {
    return &MarketHandler{
        marketService: marketService,
    }
}

func (h *MarketHandler) GetMarketData(c *gin.Context) {
    symbol := c.Param("symbol")
    data, err := h.marketService.GetMarketData(symbol)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, data)
}

func (h *MarketHandler) GetMultipleMarketData(c *gin.Context) {
    symbols := c.QueryArray("symbols")
    if len(symbols) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no symbols provided"})
        return
    }

    data, err := h.marketService.GetMultipleMarketData(symbols)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, data)
}