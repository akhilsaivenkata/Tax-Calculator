package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/akhilsaivenkata/go-tax-calculator/internal/client"
	"github.com/akhilsaivenkata/go-tax-calculator/internal/model"
	"github.com/akhilsaivenkata/go-tax-calculator/internal/service"
	"github.com/akhilsaivenkata/go-tax-calculator/pkg/logger"
)

// TaxHandler handles HTTP requests for tax calculation.
type TaxHandler struct {
	Client  client.TaxClient
	Service *service.TaxService
}

// NewTaxHandler returns a new instance of TaxHandler.
func NewTaxHandler(c client.TaxClient, s *service.TaxService) *TaxHandler {
	return &TaxHandler{
		Client:  c,
		Service: s,
	}
}

// RegisterRoutes connects the routes to Gin
func (h *TaxHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/calculate-tax", h.CalculateTax)
}

// CalculateTax is the main endpoint logic
func (h *TaxHandler) CalculateTax(c *gin.Context) {
	var req model.TaxCalculationRequest

	// Parsing request body
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.WithError(err).Warn("Failed to parse request body")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload. Make sure 'income' (float) and 'tax_year' (int) are provided.",
		})
		return
	}

	// Manual input validation
	if req.Income <= 0 {
		logger.Log.WithField("income", req.Income).Warn("Income must be positive")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Income must be a positive number greater than 0.",
		})
		return
	}

	if req.TaxYear < 2019 || req.TaxYear > 2022 {
		logger.Log.WithField("taxYear", req.TaxYear).Warn("Unsupported tax year")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tax year must be one of: 2019, 2020, 2021, 2022.",
		})
		return
	}

	logger.Log.WithFields(map[string]interface{}{
		"income":  req.Income,
		"taxYear": req.TaxYear,
	}).Info("Processing tax calculation request")

	// Fetching tax brackets from external API
	bracketsResp, err := h.Client.GetTaxBrackets(req.TaxYear)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to fetch tax brackets from API")
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to fetch tax brackets from external service. Please try again.",
		})
		return
	}

	// Calculating tax
	result := h.Service.CalculateTax(req.Income, bracketsResp.TaxBrackets)

	logger.Log.WithFields(map[string]interface{}{
		"total_tax":     result.TotalTax,
		"effective_tax": result.EffectiveTax,
	}).Info("Successfully calculated tax")

	// Returning result
	c.JSON(http.StatusOK, result)
}
