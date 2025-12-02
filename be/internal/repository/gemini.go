package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type GeminiRepository interface {
	EstimateStartingPrice(req PriceEstimationRequest) (float64, error)
}

type geminiRepo struct {
	logger *slog.Logger
	APIKey string
	client *http.Client
}

func NewGeminiRepository(logger *slog.Logger, apiKey string) GeminiRepository {
	return &geminiRepo{
		logger: logger,
		APIKey: apiKey,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type PriceEstimationRequest struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Condition   string `json:"condition"`
	Description string `json:"description"`
}

type respBody struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (r *geminiRepo) EstimateStartingPrice(req PriceEstimationRequest) (float64, error) {
	price, err := r.callGeminiModel(req, "gemini-2.5-flash")
	if err == nil {
		return price, nil
	}
	r.logger.Error("Primary model failed", "error", err)

	var lastErr error
	for i := 1; i <= 3; i++ {
		wait := time.Duration(i*i) * time.Second
		time.Sleep(wait)

		r.logger.Warn("Retrying Gemini...", "attempt", i)
		retryPrice, retryErr := r.callGeminiModel(req, "gemini-2.5-flash")

		if retryErr == nil {
			return retryPrice, nil
		}

		lastErr = retryErr
	}

	r.logger.Warn("Switching to fallback model gemini-1.5-flash", "error", lastErr)

	priceFallback, errFallback := r.callGeminiModel(req, "gemini-1.5-flash")
	if errFallback == nil {
		return priceFallback, nil
	}

	r.logger.Error("AI failed completely, requiring manual input", "error", errFallback)
	return 0, errors.New("AI failed completely, manual price needed")
}

func (r *geminiRepo) callGeminiModel(req PriceEstimationRequest, model string) (float64, error) {
	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent",
		model,
	)

	prompt := fmt.Sprintf(`
		Anda adalah AI untuk menentukan start price lelang barang.

		Kembalikan JSON saja:
		{ "starting_price": 123456 }

		Nama: %s
		Kategori: %s
		Kondisi: %s
		Deskripsi: %s
		`, req.Name, req.Category, req.Condition, req.Description)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
		"generationConfig": map[string]interface{}{
			"responseMimeType": "application/json",
			"responseJsonSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"starting_price": map[string]interface{}{"type": "integer"},
				},
				"required": []string{"starting_price"},
			},
		},
	}

	body, _ := json.Marshal(payload)
	reqHTTP, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	reqHTTP.Header.Set("x-goog-api-key", r.APIKey)
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(reqHTTP)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	var data respBody
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return 0, err
	}

	if len(data.Candidates) == 0 {
		return 0, fmt.Errorf("no AI response")
	}

	raw := strings.TrimSpace(data.Candidates[0].Content.Parts[0].Text)

	var result struct {
		StartingPrice float64 `json:"starting_price"`
	}

	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return 0, fmt.Errorf("invalid AI JSON: %w", err)
	}

	return result.StartingPrice, nil
}
