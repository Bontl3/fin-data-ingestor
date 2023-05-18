package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Bontl3/data_ingestion_microservice/internal/config"
	"github.com/Bontl3/data_ingestion_microservice/internal/database"
	pkg "github.com/Bontl3/data_ingestion_microservice/pkg/http"
)

func main() {
	// Load the configuration from the user-provided YAML config file
	currentDir, err := os.Getwd()
	configFile := "config.yml"
	configFilePath := filepath.Join(currentDir, "..", "cmd", configFile)
	log.Println(configFilePath)
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("Failed to load configuration :%s", err)
	}

	// Update the configuration with environment variables
	err = config.UpdateConfigWithEnv(cfg)
	if err != nil {
		log.Fatalf("Failed to update configuration with environment variables: %s", err)
	}

	// Set up the database connection
	db, err := database.NewDBConn(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to establish database connection: %s", err)
	}

	// Create an instance of the HTTP server
	server := pkg.NewServer(cfg, db)

	// Start the HTTP server in a separate goroutine
	go func() {
		log.Printf("Starting server on port %s...", cfg.Server.Port)
		err := http.ListenAndServe(":"+cfg.Server.Port, server.Handler)
		if err != nil {
			log.Fatalf("Failed to start HTTP server: %s", err)
		}
	}()

	// Create a channel to listen for OS signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan // Wait for OS signal

	// Create a context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shut down the server
	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Failed to gracefully shut down the server: %s", err)
	}

	log.Println("Server stopped")
}
