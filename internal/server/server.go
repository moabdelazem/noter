package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/moabdelazem/noter/internal/config"
	"github.com/moabdelazem/noter/internal/database"
	"github.com/moabdelazem/noter/internal/routes"
)

type Server struct {
	router *mux.Router
	config *config.Config
	db     *database.DB
}

func New(cfg *config.Config) *Server {
	router := mux.NewRouter()

	// Setup all routes
	routes.SetupRoutes(router)

	return &Server{
		router: router,
		config: cfg,
	}
}

// InitDB initializes the database connection
func (s *Server) InitDB() error {
	// Initialize database connection
	db, err := database.New(s.config)
	if err != nil {
		return fmt.Errorf("error initializing database: %w", err)
	}
	s.db = db

	// Run migrations
	if err := database.RunMigrations(s.config); err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	// Setup database-specific routes
	routes.SetupDBRoutes(s.router, s.db)

	return nil
}

func (s *Server) Start() error {
	// Initialize database connection
	if err := s.InitDB(); err != nil {
		return err
	}
	// Ensure we close the database connection when the server shuts down
	defer s.db.Close()

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.config.ServerPort),
		Handler: s.router,
	}

	// Graceful Shutdown
	// Start a goroutine to handle graceful shutdown
	go func() {
		// Create a buffered channel to receive OS signals
		quitChan := make(chan os.Signal, 1)
		// Listen for SIGTERM and SIGINT signals and send them to quitChan
		signal.Notify(quitChan, syscall.SIGTERM, syscall.SIGINT)
		// Block until a signal is received
		<-quitChan

		log.Println("Shutting down server...")

		// Create a context with 5 second timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// Ensure the cancel function is called to release resources
		defer cancel()

		// Attempt to gracefully shutdown the server
		if err := server.Shutdown(ctx); err != nil {
			// Log error if shutdown fails and server is forced to stop
			log.Printf("Server forced to shutdown: %v\n", err)
		}
	}()

	log.Printf("Starting the server at port %s\n", s.config.ServerPort)
	return server.ListenAndServe()
}
