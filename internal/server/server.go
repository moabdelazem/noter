package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/moabdelazem/noter/internal/config"
	"github.com/moabdelazem/noter/internal/routes"
)

type Server struct {
	router *mux.Router
	config *config.Config
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

func (s *Server) Start() error {
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

		// Create a context with 5 second timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// Ensure the cancel function is called to release resources
		defer cancel()

		// Attempt to gracefully shutdown the server
		if err := server.Shutdown(ctx); err != nil {
			// Log error if shutdown fails and server is forced to stop
			fmt.Printf("Server forced to shutdown: %v\n", err)
		}
	}()

	fmt.Printf("Starting the server at port %s\n", s.config.ServerPort)
	return server.ListenAndServe()
}
