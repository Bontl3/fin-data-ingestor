package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Bontl3/data_ingestion_microservice/internal/models"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Read the server address from the environment variable
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		fmt.Println("SERVER_ADDRESS environment variable is not set.")
		return
	}

	for {
		fmt.Print("Enter ticker symbol (or 'exit' to quit): ")
		ticker, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break // Exit on EOF (Ctrl+D)
			}
			fmt.Printf("Failed to read input: %s\n", err)
			continue
		}

		ticker = strings.TrimSpace(ticker)
		//fmt.Printf("Ticker: '%s'\n", ticker) // debugging
		if ticker == "exit" {
			break
		}

		if ticker == "" {
			fmt.Println("Ticker symbol cannot be empty.")
			continue
		}

		resp, err := http.Get(serverAddress + "/market-data/?ticker=" + ticker)
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
