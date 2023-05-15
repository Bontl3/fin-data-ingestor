package pkg

import (
	"net/http"

	"github.com/Bontl3/data_ingestion_microservice/internal/config"
	"github.com/Bontl3/data_ingestion_microservice/internal/database"
	"github.com/Bontl3/data_ingestion_microservice/pkg/http/handler"
	"github.com/Bontl3/data_ingestion_microservice/pkg/http/middleware"
)

// Server struct holds the application's config and database settings
type Server struct {
	cfg *config.Config
	db  *database.Postgres
}

// Creates and configures an HTTP server
func NewServer(cfg *config.Config, db *database.Postgres) *http.Server {
	// Initialise the server struct with config and database information
	server := &Server{
		cfg: cfg,
		db:  db,
	}

	// Initialise HTTP handlers with the shared config and db.
	marketDataHandler := handler.NewMarketDataHandler(cfg, db)

	// Create new ServerMux to handle routing of HTTP requests
	mux := http.NewServeMux()

	// Define Http routes
	// We expecct the ticker to be passed as a URL paremeter, e.g., /market-data/?ticker=MSFT
	mux.HandleFunc("/market-data/", marketDataHandler.HandleMarketDataRequest)

	// Apply the logging middleware to the mux
	loggingMux := middleware.LoggingMiddleware(mux)

	// Return a new HTTP server with the configures settings
	return &http.Server{
		Addr:    ":" + server.cfg.Server.Port,
		Handler: loggingMux,
	}
}
