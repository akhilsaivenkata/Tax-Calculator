package service

import (
	"math"
	"testing"

	"github.com/akhilsaivenkata/go-tax-calculator/internal/model"
	"github.com/akhilsaivenkata/go-tax-calculator/pkg/logger"
)

func TestCalculateTax(t *testing.T) {
	logger.Init()
	service := NewTaxService()

	brackets := []model.TaxBracket{
		{Min: 0, Max: 50000, Rate: 0.10},
		{Min: 50000, Max: 100000, Rate: 0.20},
		{Min: 100000, Rate: 0.30}, // top slab, no Max
	}

	tests := []struct {
		name          string
		income        float64
		expectedTax   float64
		expectedRate  float64
		expectedBands int
	}{
		{
			name:          "Zero income",
			income:        0,
			expectedTax:   0,
			expectedRate:  0,
			expectedBands: 0,
		},
		{
			name:          "Income within first bracket",
			income:        30000,
			expectedTax:   3000,
			expectedRate:  0.10,
			expectedBands: 1,
		},
		{
			name:          "Income spans two brackets",
			income:        75000,
			expectedTax:   (50000 * 0.10) + (25000 * 0.20),
			expectedRate:  ((50000 * 0.10) + (25000 * 0.20)) / 75000,
			expectedBands: 2,
		},
		{
			name:          "Income spans all brackets",
			income:        150000,
			expectedTax:   (50000 * 0.10) + (50000 * 0.20) + (50000 * 0.30),
			expectedRate:  ((50000 * 0.10) + (50000 * 0.20) + (50000 * 0.30)) / 150000,
			expectedBands: 3,
		},
	}

	const floatTolerance = 0.0001

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := service.CalculateTax(tc.income, brackets)

			if math.Abs(result.TotalTax-tc.expectedTax) > floatTolerance {
				t.Errorf("Expected tax %.2f, got %.2f", tc.expectedTax, result.TotalTax)
			}

			if math.Abs(result.EffectiveTax-tc.expectedRate) > floatTolerance {
				t.Errorf("Expected effective rate %.5f, got %.5f", tc.expectedRate, result.EffectiveTax)
			}

			if len(result.Breakdown) != tc.expectedBands {
				t.Errorf("Expected %d bands, got %d", tc.expectedBands, len(result.Breakdown))
			}
		})
	}
}
