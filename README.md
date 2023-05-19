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
4.  `./data_ingestion_microservice` on Unix systems or 
5.  `data_ingestion_microservice.exe` on Windows systems.
If the configuration file config.yml does not exist in the same directory, the application will ask for database configuration details. These details are necessary for the application to fetch market data.
After providing the database configuration details, the application will create a config.yml file in the same directory for future use. Next time you run the application, it will use the details from this file and will not ask for them again.
The application will then prompt you to enter a ticker symbol. Type in a valid ticker symbol and press Enter.
Optionally, you can specify a data length. If you don't want to specify a length, simply press Enter to use the default value.
The application will fetch and display the market data for the specified ticker symbol.
The application has been compiled into a binary executable for ease of use, and does not require the Go language to be installed on your machine.

Follow these steps to run the application:

Download the latest binary executable from the GitHub releases page.
Once downloaded, navigate to the directory containing the binary file using your system's command line interface.
Start the application by running ./data_ingestion_microservice on Unix systems or data_ingestion_microservice.exe on Windows systems.
If the configuration file config.yml does not exist in the same directory, the application will ask for database configuration details. These details are necessary for the application to fetch market data.
After providing the database configuration details, the application will create a config.yml file in the same directory for future use. Next time you run the application, it will use the details from this file and will not ask for them again.
The application will then prompt you to enter a ticker symbol. Type in a valid ticker symbol and press Enter.
Optionally, you can specify a data length. If you don't want to specify a length, simply press Enter to use the default value.
The application will fetch and display the market data for the specified ticker symbol.
Here's an overview of the application's functionality:

Configuration: The application loads its configuration from a YAML file. The configuration includes settings for the server (such as the port to listen on) and the database connection.

Database Connection: The application establishes a connection to a PostgreSQL database using the provided configuration settings. It sets connection pool parameters to optimize performance.

Database Migration: The application ensures the database schema is up to date by running database migrations. Migrations allow for the creation and modification of database tables and schema over time.

Data Ingestion: The application fetches market data for a specific ticker symbol from an external API. It makes HTTP requests to retrieve the data and stores it in the PostgreSQL database. The fetched data includes information such as the ticker symbol, date, open price, high price, low price, closing price, and volume.

CLI Interaction: The application provides a command-line interface (CLI) for users to interact with. Users can enter a ticker symbol and retrieve the corresponding market data from the database. 
