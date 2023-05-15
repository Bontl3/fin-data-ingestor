package main

import (
	"log"

	"os"
	"path/filepath"

	"github.com/Bontl3/data_ingestion_microservice/internal/config"
	"github.com/Bontl3/data_ingestion_microservice/internal/database"
	pkg "github.com/Bontl3/data_ingestion_microservice/pkg/http"
)

func main() {
	// Get the current file's directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %s", err)
	}

	// Construct the absolute file path to config.yml
	configFile := filepath.Join(currentDir, "config.yml")

	// Load the configuration from the encrypted config file
	config, err := config.LoadEncryptedConfig(configFile, os.Getenv("CONFIG_ENCRYPTION_KEY"))
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	// Set up the database connection
	db, err := database.NewDBConn(&config.DB)
	if err != nil {
		log.Fatalf("Failed to establish database connection: %s", err)
	}

	// Start the HTTP server
	server := pkg.NewServer(config, db)
	log.Printf("Starting server on port %s...", config.Server.Port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server exited with error: %s", err)
	}
}
