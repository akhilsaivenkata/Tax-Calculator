package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/akhilsaivenkata/go-tax-calculator/internal/model"
	"github.com/akhilsaivenkata/go-tax-calculator/pkg/logger"
)

// TaxClient defines an interface for fetching tax brackets
type TaxClient interface {
	GetTaxBrackets(year int) (*model.TaxAPIResponse, error)
}

// TaxAPIClient defines the client structure
type TaxAPIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewTaxAPIClient creates a new instance with sensible defaults
func NewTaxAPIClient(baseURL string) *TaxAPIClient {
	return &TaxAPIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 7 * time.Second,
		},
	}
}

// GetTaxBrackets fetches the brackets for a given year from the Flask API
func (c *TaxAPIClient) GetTaxBrackets(year int) (*model.TaxAPIResponse, error) {
	url := fmt.Sprintf("%s/tax-calculator/tax-year/%d", c.BaseURL, year)

	logger.Log.WithFields(map[string]interface{}{
		"year": year,
		"url":  url,
	}).Info("Calling tax API to fetch brackets")

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		logger.Log.WithError(err).WithFields(map[string]interface{}{
			"url": url,
		}).Error("Failed to reach tax API")
		return nil, fmt.Errorf("failed to reach tax API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logger.Log.WithFields(map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
			"url":    url,
		}).Error("Received non-200 response from tax API")
		return nil, fmt.Errorf("API error [%d]: %s", resp.StatusCode, string(body))
	}

	var taxResp model.TaxAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&taxResp); err != nil {
		logger.Log.WithError(err).Error("Failed to decode tax API response")
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	logger.Log.WithField("brackets_count", len(taxResp.TaxBrackets)).Info("Successfully fetched tax brackets")
	return &taxResp, nil
}
