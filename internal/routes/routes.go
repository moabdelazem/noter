package routes

import (
	"github.com/gorilla/mux"
	"github.com/moabdelazem/noter/internal/handlers"
	"github.com/moabdelazem/noter/internal/middleware"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(router *mux.Router) {
	// Add middleware
	router.Use(middleware.Logger)

	// Home route
	router.HandleFunc("/", handlers.HomeHandler)

	// Health check route
	router.HandleFunc("/health", handlers.HealthHandler)

}
