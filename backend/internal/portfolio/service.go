package portfolio

import (
    "github.com/yourusername/crypto-monitor/internal/database"
    "github.com/yourusername/crypto-monitor/internal/market"
)

type Portfolio struct {
    ID       uint    `json:"id"`
    UserID   uint    `json:"user_id"`
    Symbol   string  `json:"symbol"`
    Quantity float64 `json:"quantity"`
}

type Service struct {
    db          *database.DB
    marketSvc   *market.MarketService
}

func NewService(db *database.DB, marketSvc *market.MarketService) *Service {
    return &Service{
        db:        db,
        marketSvc: marketSvc,
    }
}

func (s *Service) GetPortfolio(userID uint) ([]Portfolio, error) {
    var portfolio []Portfolio
    err := s.db.Where("user_id = ?", userID).Find(&portfolio).Error
    return portfolio, err
}

func (s *Service) AddToPortfolio(userID uint, symbol string, quantity float64) error {
    var existing Portfolio
    err := s.db.Where("user_id = ? AND symbol = ?", userID, symbol).First(&existing).Error
    
    if err == nil {
        // Update existing position
        existing.Quantity += quantity
        return s.db.Save(&existing).Error
    }

    // Create new position
    portfolio := Portfolio{
        UserID:   userID,
        Symbol:   symbol,
        Quantity: quantity,
    }
    return s.db.Create(&portfolio).Error
}

func (s *Service) GetPortfolioValue(userID uint) (float64, error) {
    portfolio, err := s.GetPortfolio(userID)
    if err != nil {
        return 0, err
    }

    var totalValue float64
    for _, position := range portfolio {
        price, err := s.marketSvc.GetPrice(position.Symbol)
        if err != nil {
            continue
        }
        totalValue += price.Price * position.Quantity
    }

    return totalValue, nil
}