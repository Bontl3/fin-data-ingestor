package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"github.com/Bontl3/data_ingestion_microservice/internal/config"
	"github.com/Bontl3/data_ingestion_microservice/internal/database"
	"github.com/Bontl3/data_ingestion_microservice/internal/ingestion"
	"github.com/Bontl3/data_ingestion_microservice/internal/models"
)

type MarketDataHandler struct {
	cfg *config.Config
	db  *database.Postgres
}

func NewMarketDataHandler(cfg *config.Config, db *database.Postgres) *MarketDataHandler {
	return &MarketDataHandler{
		cfg: cfg,
		db:  db,
	}
}

// HandleMarketDataRequest handles the request for market data.
// It fetches the ticker symbol either from the database or from AlphaVantage,
// and writes the data as JSON to the response.
func (h *MarketDataHandler) HandleMarketDataRequest(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	length := r.URL.Query().Get("length")

	if ticker == "" {
		http.Error(w, "Missing ticker parameter", http.StatusBadRequest)
		return
	}

   // Set a default limit
   limit := 50

   // Override the limit if the length parameter is provided
   if length != "" {
	   var err error
	   limit, err = strconv.Atoi(length)
	   if err != nil {
		   http.Error(w, "Invalid length parameter", http.StatusBadRequest)
		   return
	   }
   }

   
	// Check if ticker data exists in the database
	exists, err := h.db.TickerExists(ticker)

	if err != nil {
		log.Printf("Error checking if ticker exists: %v", err) // debuging
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	var data []*models.MarketData
	if exists {
		// Fetch data from the database
		data, err = h.db.FetchMarketData(ticker, limit)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else {
		// Fetch data from Yahoo Finance and store in the database
		data, err = ingestion.FetchDataFromTwelveData(ticker)
		if err != nil {
			http.Error(w, "Failed to fetch market data", http.StatusInternalServerError)
			return
		}
	}

	// Write data as JSON to the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
