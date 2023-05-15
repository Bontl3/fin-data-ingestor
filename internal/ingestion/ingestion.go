package ingestion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Bontl3/data_ingestion_microservice/internal/models"
)

const TwelveDataAPIKey = "be99334f093d41b29242e3485fcf6873"

type TwelveDataTimeSeriesResponse struct {
	Symbol     string `json:"symbol"`
	Exchange   string `json:"exchange"`
	Currency   string `json:"currency"`
	TimeSeries []struct {
		DateTime string `json:"datetime"`
		Open     string `json:"open"`
		High     string `json:"high"`
		Low      string `json:"low"`
		Close    string `json:"close"`
		Volume   string `json:"volume"`
	} `json:"values"`
}

func FetchDataFromTwelveData(symbol string) ([]*models.MarketData, error) {
	url := fmt.Sprintf("https://api.twelvedata.com/time_series?symbol=%s&interval=1day&outputsize=30&apikey=%s", symbol, TwelveDataAPIKey)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data TwelveDataTimeSeriesResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var marketDataList []*models.MarketData
	for _, item := range data.TimeSeries {
		date, _ := time.Parse("2006-01-02 15:04:05", item.DateTime)
		open, _ := strconv.ParseFloat(item.Open, 64)
		high, _ := strconv.ParseFloat(item.High, 64)
		low, _ := strconv.ParseFloat(item.Low, 64)
		close, _ := strconv.ParseFloat(item.Close, 64)
		volume, _ := strconv.ParseInt(item.Volume, 10, 64)

		marketData := &models.MarketData{
			Ticker: data.Symbol,
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		}
		marketDataList = append(marketDataList, marketData)
	}

	return marketDataList, nil
}
