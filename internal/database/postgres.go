package database

import (
	"database/sql"
	"fmt"

	"github.com/Bontl3/data_ingestion_microservice/internal/config"
	"github.com/Bontl3/data_ingestion_microservice/internal/models"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
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
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create table if not exists -> look at creating table in onenote
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS market_data (
			ticker VARCHAR(10) NOT NULL,
			date TIMESTAMP NOT NULL,
			open NUMERIC NOT NULL,
			high NUMERIC NOT NULL,
			low NUMERIC NOT NULL,
			close NUMERIC NOT NULL,
			volume BIGINT NOT NULL
		);
		`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &Postgres{db: db}, nil
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
	return err
}

// FetchMarketData fetches market data for a specific ticker from the database
func (p *Postgres) FetchMarketData(ticker string) ([]*models.MarketData, error) {
	// SQL statement to fetch data from the database
	query := `SELECT ticker, date, open, high, low, close, volume FROM market_data WHERE ticker = $1 ORDER BY date DESC`

	rows, err := p.db.Query(query, ticker)
	if err != nil {
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

/*
In this function, we execute a SQL query to select all the data for the given ticker from
the market_data table, ordered by date in descending order (most recent first). We then iterate over
the returned rows, scanning the data into a MarketData struct, and append these structs to a slice.
Once we have iterated over all rows, we return the slice of MarketData structs.

As always, error handling is crucial. We check for errors both when scanning each row and after we
have finished iterating over the rows. If there was an error at any point, we return it along with
a nil slice. If there were no errors, we return the slice and a nil error.

Remember to replace "your_project_path" with the actual import path of your project. The import path
is generally the repository location if you're using version control like git,
such as "github.com/yourusername/yourrepository". If you're not using
 version control, you can replace "your_project_path" with the relative path from postgres.go to
 the models/market_data.go files.
*/
