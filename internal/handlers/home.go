package handlers

import (
	"encoding/json"
	"net/http"
)

// HomeHandler handles the root endpoint
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ok, Let's Start!"})
}
