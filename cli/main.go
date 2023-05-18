package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Gonum/stat"

	"github.com/Bontl3/data_ingestion_microservice/internal/models"
)

// Constant for config template
const configTemplate = `
server:
  port: 8080
db:
  host: {{.DBHost}}
  port: {{.DBPort}}
  user: {{.DBUser}}
  password: {{.DBPassword}}
  dbname: {{.DBName}}
  sslmode: {{.DBSSLMode}}
`

// Struct to hold the data for config template
type ConfigTemplate struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Function to get user input
func getUserInput(prompt string, reader *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}

// Function to generate config file
func generateConfigFile(data ConfigTemplate) error {
	// Create the config file template
	tmpl, err := template.New("config").Parse(configTemplate)
	if err != nil {
		log.Fatalf("Failed to parse config template: %v", err)
	}

	// get the root directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Failed to get current directory: %v", err)
	}

	// Construct path to the root directory
	rootDir := filepath.Join(currentDir, "..")
	// Set the path to the cmd directory
	cmdDirPath := filepath.Join(rootDir, "cmd")
	// Create the config file in the cmd directory
	configFilePath := filepath.Join(cmdDirPath, "config.yml")
	configFile, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("Failed to create a config file: %s", err)
	}

	// Execute the template and write it to the config file
	err = tmpl.Execute(configFile, &data)
	if err != nil {
		fmt.Errorf("Failed to write config file: %s", err)
	}

	// Close the config file
	err = configFile.Close()
	if err != nil {
		return fmt.Errorf("Failed to close config file: %s", err)
	}

	// Communicate with the successful creation of the config file
	fmt.Println("Configuration file (config.yml) created successfully.")
	return nil
}

// Function to run the server
func runServer(rootDir string) error {
	cmdPath := filepath.Join(rootDir, "cmd", "main.go")
	// Execute the cmd/main.go using the `go run` command
	cmd := exec.Command("cmd.exe", "/c", "start", "cmd", "/k", "go", "run", cmdPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start the server: %v", err)
	}
	return nil
}

// Function to handle user input for the ticker symbol
func getTickerSymbol(reader *bufio.Reader) (string, error) {
	fmt.Print("Enter ticker symbol (or 'exit' to quit): ")
	ticker, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "", err // Return error to signal exit
		}
		return "", fmt.Errorf("failed to read input: %v", err)
	}

	ticker = strings.TrimSpace(ticker)
	if ticker == "" {
		return "", fmt.Errorf("ticker symbol cannot be empty")
	}

	return ticker, nil
}

// Function to handle the ticker symbol entered by the user
func handleTickerSymbol(reader *bufio.Reader, serverAddress string, ticker string, stats bool) {
	for {
		// If ticker symbol is specified in command line flags, use it
		// Otherwise, ask the user for the ticker symbol
		var err error
		if ticker == "" {
			ticker, err = getTickerSymbol(reader)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println(err)
				continue
			}
		}

		if ticker == "exit" {
			break
		}

		// Ask the user for the desired data length
		length, err := getUserInput("Enter desired data length (leave empty for default): ", reader)
		if err != nil {
			fmt.Printf("Failed to read input for length: %v\n", err)
			continue
		}

		// Construct the request URL
		requestURL := serverAddress + "/market-data/?ticker=" + ticker
		if length != "" {
			requestURL += "&length=" + length
		}

		resp, err := http.Get(requestURL)
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
			fmt.Printf("Date: %s, Close: %f\n", item.Date, item.Close)
		}

		if stats {
			calculateStats(data)
		} else {
			for _, item := range data {
				fmt.Printf("Date: %s, Close: %f\n", item.Date, item.Close)
			}
		}

		// If ticker symbol is specified in command line flags, break the loop
		// after handling it once
		if ticker != "" {
			break
		}
	}
}

func calculateStats(data []models.MarketData) {
	// Create a slice to hold the closing prices
	closes := make([]float64, len(data))

	// Populate the slice with closing prices
	for i, item := range data {
		closes[i] = item.Close
	}

	// Calculate and print the mean and standard deviation of the closing prices
	mean := stat.Mean(closes, nil)
	stdDev := stat.StdDev(closes, nil)

	fmt.Printf("Mean: %f\n", mean)
	fmt.Printf("Standard Deviation: %f\n", stdDev)
}

func main() {
	// Define and parse command line flags
	tickerPtr := flag.String("ticker", "", "Ticker symbol")
	statsPtr := flag.Bool("stats", false, "Calculate descriptive statistics")

	// Override the default usage function to provide a custom help message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "This is a CLI application for fetching market data.")
		fmt.Fprintln(os.Stderr, "You can specify the ticker symbol and whether to calculate statistics using command line flags.")
		fmt.Fprintln(os.Stderr, "For example, to get the market data for the AAPL ticker and calculate statistics, run:  -ticker=AAPL -stats")
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Create a new reader for user input
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Data Ingestion Microservice Setup")

	// get the root directory
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Construct path to the root directory
	rootDir := filepath.Join(currentDir, "..")
	// Set the path to the cmd directory
	cmdDirPath := filepath.Join(rootDir, "cmd")
	// Create the config file in the cmd directory
	configFilePath := filepath.Join(cmdDirPath, "config.yml")

	// Check if config file already exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fmt.Println("Configuration file not found. Please provide the following information to generate the config file:")

		// Existing code for getting user input and generating the config file...

		fmt.Println("Please provide the following information to generate the config file:")

		// Ask the user for input top generatye the config file
		dbHost, err := getUserInput("Enter Database Host: ", reader)
		dbPort, err := getUserInput("Enter Database Port: ", reader)
		dbUser, err := getUserInput("Enter Database User ", reader)
		dbPassword, err := getUserInput("Enter Database Password: ", reader)
		dbName, err := getUserInput("Enter Database Name: ", reader)
		dbSSLMode, err := getUserInput("Enter Database SSL Mode (disable, require, verify-ca, verify-full): ", reader)

		// Create a data structure to hold the template variables
		data := ConfigTemplate{
			DBHost:     dbHost,
			DBPort:     dbPort,
			DBUser:     dbUser,
			DBPassword: dbPassword,
			DBName:     dbName,
			DBSSLMode:  dbSSLMode,
		}

		// Generate the config file
		err = generateConfigFile(data)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		// If an error occurred for a reason other than the file not existing, log and exit
		log.Fatalf("An error occurred while checking for the configuration file: %v", err)
	}

	err = runServer(rootDir)
	if err != nil {
		log.Fatal(err)
	}

	// Read the server address from the environment variable
	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		fmt.Println("SERVER_ADDRESS environment variable is not set.")
		return
	}
	// Handle ticker symbol based on command line flags
	if *tickerPtr != "" {
		handleTickerSymbol(reader, serverAddress, *tickerPtr, *statsPtr)
	} else {
		handleTickerSymbol(reader, serverAddress, "", false)
	}
}
