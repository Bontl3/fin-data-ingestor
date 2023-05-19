# Data Ingestion Microservice Application

# Status: App Being Updated
Application is being updated. Please come back anytime after 13:00 pm 20 May 2023 to use or review the code. In the meantime you are more than welcome to browse around :) 

Again please note, this is an old version and might not work on your machince. 

## Introduction

This is a Command-Line Interface (CLI) application for fetching market data. It is written in Go and allows users to fetch market data for a specified ticker symbol. Optionally, the application can calculate descriptive statistics (mean and standard deviation) for the closing prices of the requested market data.

## The application has been compiled into a binary executable for ease of use, and does not require the Go language to be installed on your machine.

Follow these steps to run the application:

1. Download the latest binary executable from the GitHub releases page.

2. Once downloaded, navigate to the directory containing the binary file using your system's command line interface.
3. Start the application by running:
* `./data_ingestion_microservice` on Unix systems or 
*  `data_ingestion_microservice.exe` on Windows systems.
4. If the configuration file config.yml does not exist in the same directory, the application will ask for database configuration details. These details are necessary for the application to fetch market data.
5. After providing the database configuration details, the application will create a config.yml file in the same directory for future use. Next time you run the application, it will use the details from this file and will not ask for them again.
6. The application will then prompt you to enter a ticker symbol. Type in a valid ticker symbol and press Enter.
7. Optionally, you can specify a data length. If you don't want to specify a length, simply press Enter to use the default value.
8. The application will fetch and display the market data for the specified ticker symbol.

## Advanced Usage

The application accepts two optional command line flags for advanced usage:

* `-ticker=<TICKER_SYMBOL>`
