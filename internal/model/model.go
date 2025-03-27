package model

// TaxBracket maps to each bracket in the Flask API response
type TaxBracket struct {
	Min  float64 `json:"min"`
	Max  float64 `json:"max,omitempty"` // omitempty: max might be missing for top slab
	Rate float64 `json:"rate"`
}

// TaxAPIResponse is the full response from the Flask service
type TaxAPIResponse struct {
	TaxBrackets []TaxBracket `json:"tax_brackets"`
}

// TaxCalculationRequest is your API's input from the client
type TaxCalculationRequest struct {
	Income  float64 `json:"income" binding:"required"`
	TaxYear int     `json:"tax_year" binding:"required"`
}

// TaxBandResult shows how much tax was paid in each bracket
type TaxBandResult struct {
	Min     float64 `json:"min"`
	Max     float64 `json:"max,omitempty"`
	Rate    float64 `json:"rate"`
	TaxPaid float64 `json:"tax_paid"`
}

// TaxCalculationResult is your API's final output
type TaxCalculationResult struct {
	TotalTax     float64         `json:"total_tax"`
	EffectiveTax float64         `json:"effective_tax_rate"`
	Breakdown    []TaxBandResult `json:"tax_by_bracket"`
}
