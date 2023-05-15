package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Bontl3/data_ingestion_microservice/internal/models"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter ticker symbol (or 'exit' to quit): ")
		ticker, _ := reader.ReadString('\n')
		ticker = strings.TrimSpace(ticker)
		//fmt.Printf("Ticker: '%s'\n", ticker) // debugging
		if ticker == "exit" {
			break
		}
		/*
			if len(ticker) == 0 {
				fmt.Println("Ticker symbol cannot be empty.")
				continue
			}
		*/
		resp, err := http.Get("http://localhost:8080/market-data/?ticker=" + ticker)
		if err != nil {
			fmt.Printf("An error occurred: %s\n", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Server returned an error: %s\n", resp.Request.URL)

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Failed to read response body: %v", err)
			} else {
				fmt.Printf("Response body: %s\n", string(bodyBytes))
			}
			continue
		}

		// Read and print the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response: %s\n", err)
			continue
		}
		var data []models.MarketData
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Printf("Failed to parse JSON response: %s\n", err)
			continue
		}

		for _, item := range data {
			fmt.Printf("Data: %+v\n", item)
		}
	}
}
