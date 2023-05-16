-- +goose Up
CREATE TABLE IF NOT EXISTS market_data_db (
	ticker VARCHAR(10) NOT NULL,
	date TIMESTAMP NOT NULL,
	open NUMERIC NOT NULL,
	high NUMERIC NOT NULL,
	low NUMERIC NOT NULL,
	close NUMERIC NOT NULL,
	volume BIGINT NOT NULL
);
-- +goose Down
DROP TABLE IF EXISTS market_data_db;