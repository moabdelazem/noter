package main

import (
	"log"

	"github.com/moabdelazem/noter/internal/config"
	"github.com/moabdelazem/noter/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Create and start server
	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
