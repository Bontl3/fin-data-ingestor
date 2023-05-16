package main

import (
	"context"
	"io/ioutil"
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
	"github.com/pressly/goose"
)

func main() {
	// Get the current file's directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %s", err)
	}

	// Path to the plain configuration file
	configFile := filepath.Join(currentDir, "config.yml")

	encryptionKey := os.Getenv("CONFIG_ENCRYPTION_KEY")

	// read the plain configuration
	plainConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read plain configuartion file: %s", err)
	}

	// Encrpt the plain configuration file
	encryptedConfig, err := config.Encrypt(plainConfig, encryptionKey)
	if err != nil {
		log.Fatalf("Failed to encrpt configuration: %s", err)
	}

	// Write the encrypted configurayion to a new file
	encryptedConfigFile := filepath.Join(currentDir, "encrypted_config.yml")
	err = ioutil.WriteFile(encryptedConfigFile, encryptedConfig, 0644)
	if err != nil {
		log.Fatalf("Failed to encrypt configuration File: %s", err)
	}

	// Load the encrypyted configuration
	config, err := config.LoadEncryptedConfig(encryptedConfigFile, encryptionKey)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	// Set up the database connection
	db, err := database.NewDBConn(&config.DB)

	if err != nil {
		log.Fatalf("Failed to establish database connection: %s", err)
	}

	// Run migrations
	migrationsDir := filepath.Join(currentDir, "migrations")
	migrationsDir = filepath.ToSlash(migrationsDir) // Convert path to use forward slashes
	err = runMigrations(db, migrationsDir)
	if err != nil {
		log.Printf(migrationsDir)
		log.Fatalf("Failed to apply migrations: %s", err)
	}

	// Start the HTTP server
	server := pkg.NewServer(config, db)
	log.Printf("Starting server on port %s...", config.Server.Port)
	// Create a channel to listen for OS signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server exited with error: %s", err)
		}
	}()

	// Block until an OS signal is received
	<-stopChan
	log.Println("Shutting down the server...")

	// Create a context with a timeout for the server's shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %s", err)
	}

	log.Println("Server stopped")
}

func runMigrations(db *database.Postgres, migrationsDir string) error {
	// Run the migrations using Goose
	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	gooseDB, err := goose.OpenDBWithDriver("postgres", db.GetDatabaseURL())
	if err != nil {
		return err
	}

	err = goose.Up(gooseDB, migrationsDir)
	if err != nil {
		return err
	}

	return nil
}
