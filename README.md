# fin-data-ingestor

The Data Ingestion Microservice is a command-line application that allows users to ingest and process market data. It provides functionality for connecting to a PostgreSQL database, ingesting data, and performing various operations on the data.

## Installation
To install the Data Ingestion Microservice CLI application on your desktop, follow these steps:

1. Clone the repository:

`git clone https://github.com/your-username/data_ingestion_microservice.git`
2. Navigate to the cli directory:

`cd data_ingestion_microservice/cli`

3. Build the application:

`go build -o data_ingestion_cli`

4. Add the executable to your system's PATH:

* Linux/Mac:

shell
Copy code
sudo mv data_ingestion_cli /usr/local/bin/
Windows:

Add the data_ingestion_cli executable to a directory included in your system's PATH environment variable.

Verify the installation:

shell
Copy code
data_ingestion_cli --version
You should see the version information displayed if the installation was successful.

Usage
The Data Ingestion Microservice CLI application provides various commands and options for interacting with the market data. Here are some examples:

To ingest market data from a CSV file:

shell
Copy code
data_ingestion_cli ingest --file <path-to-csv-file>
To fetch market data for a specific ticker:

shell
Copy code
data_ingestion_cli fetch --ticker <ticker-symbol>
For a complete list of commands and options, use the --help flag:

shell
Copy code
data_ingestion_cli --help
Configuration
The Data Ingestion Microservice CLI application uses a configuration file (config.yml) to specify database connection details and other settings. Before running the application, make sure to update the configuration file with your own values.

Contributing
Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.
