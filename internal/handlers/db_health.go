package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/moabdelazem/noter/internal/database"
)

// DBHealth represents the database health check response
type DBHealth struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// DBHealthHandler handles the database health check endpoint
func DBHealthHandler(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		// Ping the database
		err := db.Ping(ctx)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(DBHealth{
				Status:    "error",
				Message:   "Database connection failed: " + err.Error(),
				Timestamp: time.Now().Format(time.RFC3339),
			})
			return
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(DBHealth{
			Status:    "ok",
			Message:   "Database connection is healthy",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}
}
