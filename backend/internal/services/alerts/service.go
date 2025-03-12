package alerts

import (
    "context"
    "time"
    "github.com/yourusername/crypto-monitor/internal/models"
    "github.com/yourusername/crypto-monitor/internal/database"
    "github.com/yourusername/crypto-monitor/internal/services/market"
)

type Service struct {
    db           *database.DB
    marketSvc    *market.Service
    notifications chan models.Alert
}

func NewService(db *database.DB, marketSvc *market.Service) *Service {
    s := &Service{
        db:            db,
        marketSvc:     marketSvc,
        notifications: make(chan models.Alert, 100),
    }
    go s.startAlertMonitor()
    return s
}

func (s *Service) CreateAlert(userID uint, symbol string, targetPrice float64) error {
    alert := models.Alert{
        UserID:      userID,
        Symbol:      symbol,
        TargetPrice: targetPrice,
        IsActive:    true,
    }

    return s.db.Create(&alert).Error
}

func (s *Service) GetUserAlerts(userID uint) ([]models.Alert, error) {
    var alerts []models.Alert
    err := s.db.Where("user_id = ? AND is_active = true", userID).Find(&alerts).Error
    return alerts, err
}

func (s *Service) startAlertMonitor() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        var activeAlerts []models.Alert
        if err := s.db.Where("is_active = true").Find(&activeAlerts).Error; err != nil {
            continue
        }

        for _, alert := range activeAlerts {
            price, err := s.marketSvc.GetPrice(alert.Symbol)
            if err != nil {
                continue
            }

            if (alert.TargetPrice > 0 && price >= alert.TargetPrice) ||
               (alert.TargetPrice < 0 && price <= -alert.TargetPrice) {
                alert.IsActive = false
                s.db.Save(&alert)
                s.notifications <- alert
            }
        }
    }
}
