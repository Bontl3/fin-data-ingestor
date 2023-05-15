package models

import "time"

type MarketData struct {
	Ticker string
	Date   time.Time
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume int64
}
