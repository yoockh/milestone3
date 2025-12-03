package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type AIRepository interface {
	EstimateStartingPrice(req PriceEstimationRequest) (float64, error)
}

type aiRepo struct {
	logger       *slog.Logger
	GeminiAPIKey string
	httpClient   *http.Client
}

func NewAIRepository(logger *slog.Logger, geminiAPIKey string) AIRepository {
	return &aiRepo{
		logger:       logger,
		GeminiAPIKey: geminiAPIKey,
		httpClient:   &http.Client{Timeout: 15 * time.Second},
	}
}

type PriceEstimationRequest struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Condition   string `json:"condition"`
	Description string `json:"description"`
}

type geminiRespBody struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func (r *aiRepo) EstimateStartingPrice(req PriceEstimationRequest) (float64, error) {
	price, err := r.callGeminiModel(req, "gemini-2.5-flash")
	if err == nil {
		return price, nil
	}
	r.logger.Error("Gemini failed completely, manual input for starting price", "error", err)

	return 10000, errors.New("all AI models failed, manual price needed")
}

func (r *aiRepo) callGeminiModel(req PriceEstimationRequest, model string) (float64, error) {
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

	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	reqHTTP, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(body)))
	if err != nil {
		return 0, err
	}
	reqHTTP.Header.Set("x-goog-api-key", r.GeminiAPIKey)
	reqHTTP.Header.Set("Content-Type", "application/json")

	resp, err := r.httpClient.Do(reqHTTP)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var data geminiRespBody
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return 0, err
	}

	if len(data.Candidates) == 0 {
		return 0, errors.New("no AI response")
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
