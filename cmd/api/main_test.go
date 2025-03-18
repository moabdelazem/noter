package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleHealthCheck)

	// Call the handler directly and pass in our request and response recorder
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Could not unmarshal response: %v", err)
	}

	// Validate response fields
	if response["status"] != "ok" {
		t.Errorf("handler returned unexpected status: got %v want %v",
			response["status"], "ok")
	}

	if response["success"] != "true" {
		t.Errorf("handler returned unexpected success value: got %v want %v",
			response["success"], "true")
	}

	// Check time field exists and is in expected format
	_, err = time.Parse(time.RFC3339, response["time"])
	if err != nil {
		t.Errorf("handler returned invalid time format: %v", err)
	}
}

// TestRouter tests the entire router setup including middleware and routes
func TestRouter(t *testing.T) {
	// Create a new router with the middleware
	router := mux.NewRouter()
	router.Use(LoggerMiddleware)
	router.HandleFunc("/health", handleHealthCheck)

	// Create a test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Send a request to the server
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Could not send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	// Check response body
	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Could not parse response body: %v", err)
	}

	// Basic validation
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok'; got %v", response["status"])
	}
}
