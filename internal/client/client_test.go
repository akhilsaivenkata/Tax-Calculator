package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akhilsaivenkata/go-tax-calculator/internal/model"
	"github.com/akhilsaivenkata/go-tax-calculator/pkg/logger"
)

func TestGetTaxBrackets(t *testing.T) {
	logger.Init()
	tests := []struct {
		name           string
		mockStatusCode int
		mockResponse   interface{}
		expectError    bool
		skipServer     bool
	}{
		{
			name:           "Success - valid JSON",
			mockStatusCode: http.StatusOK,
			mockResponse: model.TaxAPIResponse{
				TaxBrackets: []model.TaxBracket{
					{Min: 0, Max: 50000, Rate: 0.10},
				},
			},
			expectError: false,
		},
		{
			name:           "Failure - API returns 500",
			mockStatusCode: http.StatusInternalServerError,
			mockResponse:   `Internal Server Error`,
			expectError:    true,
		},
		{
			name:           "Failure - invalid JSON response",
			mockStatusCode: http.StatusOK,
			mockResponse:   `not-a-json`,
			expectError:    true,
		},
		{
			name:           "Failure - unreachable server",
			mockStatusCode: 0, // ignored
			mockResponse:   nil,
			expectError:    true,
			skipServer:     true, // <-- custom flag
		},
		{
			name:           "Failure - non-200 with valid JSON body",
			mockStatusCode: http.StatusNotFound,
			mockResponse: model.TaxAPIResponse{
				TaxBrackets: []model.TaxBracket{},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var baseURL string
			if tc.skipServer {
				// Simulate DNS failure / unreachable host
				baseURL = "http://localhost:9999" // invalid or closed port
			} else {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tc.mockStatusCode)
					switch body := tc.mockResponse.(type) {
					case string:
						fmt.Fprint(w, body)
					default:
						json.NewEncoder(w).Encode(body)
					}
				}))
				defer server.Close()
				baseURL = server.URL
			}

			client := NewTaxAPIClient(baseURL)

			_, err := client.GetTaxBrackets(2022)
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}
		})
	}
}
