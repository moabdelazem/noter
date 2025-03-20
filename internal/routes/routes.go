package routes

import (
	"github.com/gorilla/mux"
	"github.com/moabdelazem/noter/internal/database"
	"github.com/moabdelazem/noter/internal/handlers"
	"github.com/moabdelazem/noter/internal/middleware"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(router *mux.Router) {
	// Add middleware
	router.Use(middleware.Logger)

	// Home route
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	// Health check route
	router.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
}

// SetupDBRoutes configures routes that require a database connection
func SetupDBRoutes(router *mux.Router, db *database.DB) {
	// Database health check route
	router.HandleFunc("/db/health", handlers.DBHealthHandler(db)).Methods("GET")

	// Create note repository
	noteRepo := database.NewNoteRepository(db)

	// Create note handler
	noteHandler := handlers.NewNoteHandler(noteRepo)

	// Note routes
	notesRouter := router.PathPrefix("/notes").Subrouter()
	notesRouter.HandleFunc("", noteHandler.GetAllNotes).Methods("GET")      // GET /notes - get all notes
	notesRouter.HandleFunc("", noteHandler.CreateNote).Methods("POST")      // POST /notes - create a new note
	notesRouter.HandleFunc("/{id}", noteHandler.GetNoteByID).Methods("GET") // GET /notes/{id} - get a note by ID
}
