package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// HealthHandler is a simplified version of the actual health handler for testing
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"status":  "ok",
		"success": "true",
		"time":    time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

func TestHealthEndpoint(t *testing.T) {
	// Create a new router
	router := mux.NewRouter()
	router.HandleFunc("/health", HealthHandler).Methods("GET")

	// Create a request
	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse response
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response fields
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "true", response["success"])

	// Check time field exists and is in expected format
	_, err = time.Parse(time.RFC3339, response["time"])
	assert.NoError(t, err)
}
