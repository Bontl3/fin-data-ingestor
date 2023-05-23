# Data Ingestion Microservice Application

# Status: App Being Updated
Application is being updated. Please come back anytime after 8 am CAT 22 May 2023 to use or review the code. In the meantime you are more than welcome to browse around :) 

Again please note, this is an old version and might not work on your machince. 

# Data Ingestion Microservice

## Overview

The Data Ingestion Microservice is a command-line interface (CLI) application written in Go. This application allows you to ingest financial data and store it in a PostgreSQL database.

## Prerequisites

Before you begin, ensure you have met the following requirements:

* You have a computer system with a recent version of Windows.
* Please note that the application currently supports Windows only. Compatibility with MacOS and Linux will be added in due course
* You have PostgreSQL installed and properly configured. 

**Important**: This application requires a PostgreSQL database to function correctly. The application will prompt you for your database credentials to generate a configuration file. If you don't have PostgreSQL installed, you can download it from the official website: [https://www.postgresql.org/download/](https://www.postgresql.org/download/)

## PostgreSQL Setup

Here is a simple guide on how to set up a PostgreSQL database:
1. Install PostgreSQL from
2. Once installed, create a new database which will be used by the application. This can be done by using the PostgreSQL command line interface (psql) or a graphical interface such as PgAdmin.
3. Ensure that you have the following database details:
 * Database Host (e.g., localhost if running on your machine)
  * Database Port (typically 5432 for PostgreSQL)
  * Database User
  * Database Password
  * Database Name
  * SSL Mode (Set to 'disable' for local setup)

## How to Install the Application

1. Go to the repository page on GitHub: [https://github.com/Bontl3/fin-data-ingestor](https://github.com/Bontl3/fin-data-ingestor)
2. Click on the "Releases" tab.
3. Download the latest release for your operating system.
4. Extract the downloaded file to your desired location.

## How to Use the Application

1. Open a terminal or command prompt.
2. Navigate to the location where you extracted the downloaded file.
3. Run the application using the command: `./data_ingestion_cli_app`
4. Follow the prompts to input your PostgreSQL database credentials and other details as requested by the application.

## Commands

`-home`: Displays the home screen with basic instructions.

`-ticker <ticker symbol>`: Queries and displays the financial data related to the specified ticker symbol. Please replace `<ticker symbol>` with the actual ticker you are interested in (e.g., `AAPL` for Apple Inc.).

`-stats <ticker symbol>`: Displays descriptive statistics about the ticker you have just entered.

`-exit`: Safely shuts down the application.

## Support

If you encounter any problems or have questions, please create an issue on the GitHub repository: [https://github.com/Bontl3/fin-data-ingestor/issues](https://github.com/Bontl3/fin-data-ingestor/issues)

## License

This project is licensed under the terms of the MIT license.

