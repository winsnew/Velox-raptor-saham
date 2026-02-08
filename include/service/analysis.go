package service

import (
	"fmt"
	"time"
	"velox-raptor-saham/include/models"
	"velox-raptor-saham/include/storage"
)

func AnalyzeStock(telegramID int64, symbol string) (*models.AnalysisResult, error) {
	var user models.User
	storage.DB.Where("telegram_id = ?", telegramID).FirstOrCreate(&user, models.User{TelegramID: telegramID})

	if user.Plan == "free" {
		now := time.Now()
		if now.Day() != user.LastReqAt.Day() || now.Month() != user.LastReqAt.Month() || now.Year() != user.LastReqAt.Year() {
			user.QuotaUsed = 0
		}

		if user.QuotaUsed >= 5 {
			return nil, fmt.Errorf("Daily Limit User")
		}
	}

	result, err := CallPythonEngine(symbol)
	if err != nil {
		return nil, err
	}

	if result.Status == "error" {
		return nil, fmt.Errorf(result.Error)
	}

	user.QuotaUsed++
	user.LastReqAt = time.Now()
	storage.DB.Save(&user)

	return result, nil
}
