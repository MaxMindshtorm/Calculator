package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaxMindshtorm/calculator/internal/application"
)

type Expression struct {
	Expression string `json:"expression"`
}

type Result struct {
	Res float64 `json:"result"`
}

type Error struct {
	Err string `json:"error"`
}

func TestApplicationHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           Expression
		expectedStatus int
		expectedResult interface{}
	}{
		{
			name:           "Valid expression",
			method:         http.MethodPost,
			body:           Expression{Expression: "3 + 5"},
			expectedStatus: http.StatusOK,
			expectedResult: Result{Res: 8},
		},
		{
			name:           "Invalid method",
			method:         http.MethodGet,
			body:           Expression{},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedResult: Error{Err: "Method is not allowed"},
		},
		{
			name:           "Empty expression",
			method:         http.MethodPost,
			body:           Expression{Expression: ""},
			expectedStatus: http.StatusBadRequest,
			expectedResult: Error{Err: "Bad request"},
		},
		{
			name:           "Invalid expression",
			method:         http.MethodPost,
			body:           Expression{Expression: "3 +"},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResult: Error{Err: "Expression is not valid: неверное выражение"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, _ := json.Marshal(test.body)
			req := httptest.NewRequest(test.method, "/api/v1/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			http.HandlerFunc(application.ApplicationHandler).ServeHTTP(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, rec.Code)
			}

			var result Result
			if rec.Code != http.StatusOK && rec.Code == test.expectedStatus {
				return
			}
			err := json.NewDecoder(rec.Body).Decode(&result)
			if err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if result != test.expectedResult {
				t.Errorf("Expected result %+v, got %+v", test.expectedResult, result)
			}
		})
	}
}
