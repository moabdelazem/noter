package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// Health represents the health check response
type Health struct {
	Status  string `json:"status"`
	Success string `json:"success"`
	Time    string `json:"time"`
}

// HealthHandler handles the health check endpoint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	health := Health{
		Status:  "ok",
		Success: "true",
		Time:    time.Now().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(health)
}
