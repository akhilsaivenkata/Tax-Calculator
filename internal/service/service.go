package service

import (
	"github.com/akhilsaivenkata/go-tax-calculator/internal/model"
	"github.com/akhilsaivenkata/go-tax-calculator/pkg/logger"
)

// TaxService handles the tax calculation
type TaxService struct{}

// NewTaxService returns a new instance
func NewTaxService() *TaxService {
	return &TaxService{}
}

// CalculateTax method applies the marginal tax logic
func (s *TaxService) CalculateTax(income float64, brackets []model.TaxBracket) model.TaxCalculationResult {
	logger.Log.WithField("income", income).Info("Starting tax calculation")

	var totalTax float64
	var breakdown []model.TaxBandResult

	for _, bracket := range brackets {
		var upper float64
		if bracket.Max == 0 {
			// No max means it's the highest bracket (infinite upper bound)
			upper = income
		} else {
			upper = bracket.Max
		}

		// Taxable amount = income in this bracket
		if income > bracket.Min {
			taxable := min(income, upper) - bracket.Min
			taxPaid := taxable * bracket.Rate

			logger.Log.WithFields(map[string]interface{}{
				"bracket_min": bracket.Min,
				"bracket_max": bracket.Max,
				"rate":        bracket.Rate,
				"taxable":     taxable,
				"tax_paid":    taxPaid,
			}).Info("Calculated tax for bracket")

			totalTax += taxPaid

			breakdown = append(breakdown, model.TaxBandResult{
				Min:     bracket.Min,
				Max:     bracket.Max,
				Rate:    bracket.Rate,
				TaxPaid: taxPaid,
			})
		}
	}

	// Avoid division by zero
	var effectiveRate float64
	if income > 0 {
		effectiveRate = totalTax / income
	}

	logger.Log.WithFields(map[string]interface{}{
		"total_tax":      totalTax,
		"effective_rate": effectiveRate,
	}).Info("Finished tax calculation")

	return model.TaxCalculationResult{
		TotalTax:     totalTax,
		EffectiveTax: effectiveRate,
		Breakdown:    breakdown,
	}
}

// min returns the smaller of two floats
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
