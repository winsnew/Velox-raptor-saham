package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramID int64 `gorm:"uniqueIndex"`
	Username   string
	Plan       string `gorm:"default:free"`
	QuotaUsed  int    `gorm:"default:0"`
	LastReqAt  time.Time
}

type AnalysisResult struct {
	Status             string  `json:"status"`
	Symbol             string  `json:"symbol"`
	CurrentPrice       float64 `json:"current_price"`
	Prediction         float64 `json:"prediction"`
	ChangePct          float64 `json:"change_pct"`
	Signal             string  `json:"signal"`
	Support            float64 `json:"support"`
	Resistance         float64 `json:"resistance"`
	RMSE               float64 `json:"rmse"`
	ConfidenceInterval string  `json:"confidence_interval"`
	ConfidenceLevel    string  `json:"confidence_level"`
	Error              string  `json:"error,omitempty"`
}
