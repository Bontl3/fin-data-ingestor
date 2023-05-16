package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Bontl3/data_ingestion_microservice/internal/config"
	"github.com/Bontl3/data_ingestion_microservice/internal/models"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db  *sql.DB
	cfg *config.DatabaseConfig
}

// Connect function establishes a connection to the PostgreSQL database
// using the settings from the DatabaseConfig struct.

func NewDBConn(cfg *config.DatabaseConfig) (*Postgres, error) {
	// Create a connection string using the fields from the DatabaseConfig struct
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	// connect to the database using the connection string and pgxpool.connect
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// if ther's an error establishing the connection, return nil and the error
		return nil, err
	}

	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(8)
	db.SetConnMaxLifetime(time.Minute * 5)

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &Postgres{db: db, cfg: cfg}, nil
}

// Function cheks if ticket data exists in the database
func (p *Postgres) TickerExists(ticker string) (bool, error) {
	var exists bool
	query := `SELECT exists(SELECT 1 FROM market_data WHERE ticker =$1)`
	err := p.db.QueryRow(query, ticker).Scan(&exists)
	return exists, err
}

// Store market data in the database
func (p *Postgres) StoreMarketData(data *models.MarketData) error {
	// SQL statement to insert data into the database
	query := `INSERT INTO market_data (ticker, date, open, high, low, close, volume) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := p.db.Exec(query, data.Ticker, data.Date, data.Open, data.High, data.Low, data.Close, data.Volume)
	if err != nil {
		log.Printf("Error storing market data: %s", err)
		return err
	}
	return nil
}

// FetchMarketData fetches market data for a specific ticker from the database
func (p *Postgres) FetchMarketData(ticker string) ([]*models.MarketData, error) {
	// SQL statement to fetch data from the database
	query := `SELECT ticker, date, open, high, low, close, volume FROM market_data WHERE ticker = $1 ORDER BY date DESC`

	rows, err := p.db.Query(query, ticker)
	if err != nil {
		log.Printf("Error fecthing market data: %s", err)
		return nil, err
	}
	defer rows.Close()

	var marketDataList []*models.MarketData
	for rows.Next() {
		var data models.MarketData
		err := rows.Scan(&data.Ticker, &data.Date, &data.Open, &data.High, &data.Low, &data.Close, &data.Volume)
		if err != nil {
			return nil, err
		}
		marketDataList = append(marketDataList, &data)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return marketDataList, nil
}

// GetDatabaseURL returns a string representation of the database URL
func (p *Postgres) GetDatabaseURL() string {
	fmt.Printf("Host: %s\n", p.cfg.Host)
	fmt.Printf("Port: %d\n", p.cfg.Port)
	fmt.Printf("User: %s\n", p.cfg.User)
	fmt.Printf("Password: %s\n", p.cfg.Password)
	fmt.Printf("DBName: %s\n", p.cfg.DBName)
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.cfg.User, p.cfg.Password, p.cfg.Host, p.cfg.Port, p.cfg.DBName)
}
