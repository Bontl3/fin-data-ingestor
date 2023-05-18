# fin-data-ingestor

# Status: App Being Updated
Application is being updated. Please come back anytime after 23:00 pm 18 May 2023 to use or review the code. In the meantime you are more than welcome to browse around :) 

Again please note, this is an old version and might not work on your machince. 

## App Overview

This is data ingestion microservice. It is designed to fetch market data for a specific ticker symbol from an external API and store it in a PostgreSQL database. Users can interact with the microservice through a command-line interface (CLI) to retrieve market data for a particular ticker symbol.

Here's an overview of the application's functionality:

Configuration: The application loads its configuration from a YAML file. The configuration includes settings for the server (such as the port to listen on) and the database connection.

Database Connection: The application establishes a connection to a PostgreSQL database using the provided configuration settings. It sets connection pool parameters to optimize performance.

Database Migration: The application ensures the database schema is up to date by running database migrations. Migrations allow for the creation and modification of database tables and schema over time.

Data Ingestion: The application fetches market data for a specific ticker symbol from an external API. It makes HTTP requests to retrieve the data and stores it in the PostgreSQL database. The fetched data includes information such as the ticker symbol, date, open price, high price, low price, closing price, and volume.

CLI Interaction: The application provides a command-line interface (CLI) for users to interact with. Users can enter a ticker symbol and retrieve the corresponding market data from the database. 
