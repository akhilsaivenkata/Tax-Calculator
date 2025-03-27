package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akhilsaivenkata/go-tax-calculator/internal/model"
	"github.com/akhilsaivenkata/go-tax-calculator/internal/service"
	"github.com/akhilsaivenkata/go-tax-calculator/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	logger.Init()
	gin.SetMode(gin.TestMode)
}

// --- Mock Client Implementation ---
type mockClient struct {
	mockResponse *model.TaxAPIResponse
	mockError    error
}

func (m *mockClient) GetTaxBrackets(year int) (*model.TaxAPIResponse, error) {
	return m.mockResponse, m.mockError
}

func setupRouter(client *mockClient) *gin.Engine {
	r := gin.Default()
	svc := service.NewTaxService()
	h := NewTaxHandler(client, svc)
	h.RegisterRoutes(r)
	return r
}

// --- Tests ---

func TestHandler_CalculateTax_ValidRequest(t *testing.T) {
	mock := &mockClient{
		mockResponse: &model.TaxAPIResponse{
			TaxBrackets: []model.TaxBracket{
				{Min: 0, Max: 50000, Rate: 0.10},
				{Min: 50000, Rate: 0.20},
			},
		},
	}

	router := setupRouter(mock)

	reqBody := model.TaxCalculationRequest{
		Income:  75000,
		TaxYear: 2022,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.TaxCalculationResult
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.TotalTax > 0)
	assert.Equal(t, 2, len(resp.Breakdown))
}

func TestHandler_CalculateTax_InvalidInput(t *testing.T) {
	mock := &mockClient{}
	router := setupRouter(mock)

	badPayload := `{"income": -1000, "tax_year": 2022}`
	req, _ := http.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBufferString(badPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Income must be a positive number")
}

func TestHandler_CalculateTax_MissingFields(t *testing.T) {
	mock := &mockClient{}
	router := setupRouter(mock)

	reqBody := `{"tax_year": 2022}` // missing income
	req, _ := http.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request payload")
}

func TestHandler_CalculateTax_InvalidTaxYear(t *testing.T) {
	mock := &mockClient{}
	router := setupRouter(mock)

	reqBody := model.TaxCalculationRequest{
		Income:  50000,
		TaxYear: 2025, // unsupported
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Tax year must be one of")
}

func TestHandler_CalculateTax_ExternalAPIFailure(t *testing.T) {
	mock := &mockClient{
		mockError: assert.AnError, // simulate failure
	}
	router := setupRouter(mock)

	reqBody := model.TaxCalculationRequest{
		Income:  50000,
		TaxYear: 2022,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/calculate-tax", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to fetch tax brackets")
}
