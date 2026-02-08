package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"velox-raptor-saham/include/models"

	"github.com/spf13/viper"
)

func CallPythonEngine(symbol string) (*models.AnalysisResult, error) {
	apiURL := viper.GetString("python_service.url")
	if apiURL == "" {
		return nil, fmt.Errorf("python_service.url not set in config")
	}

	timeout := viper.GetInt("python_service.timeout_second")

	reqBody := map[string]string{
		"symbol": symbol,
	}
	jsonBody, _ := json.Marshal(reqBody)

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Post(apiURL+"/analyze", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to connect python engine: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result models.AnalysisResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed parsing json: %v (Raw: %s)", err, string(body))
	}

	return &result, nil
}
