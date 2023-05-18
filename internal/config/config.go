package config

import (
	"log"
	"os"
    "fmt"
    "io/ioutil"
    "strconv"
	"gopkg.in/yaml.v3"
)

// Define a Config struct to hold the configuration settings
// Serves as container for overall application configuration settings
type Config struct {
	Server ServerConfig   `yaml:"server"`
	DB     DatabaseConfig `yaml:"db"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

// Define a DatabaseConfig struct to hold the database-specific settings.
// Focused specifically on the database configuration settings.
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// LoadConfig loads the configuration from the specified YAML file
func LoadConfig(FilePath string)(*Config, error){
	// read the contents of the config file
	configData, err := ioutil.ReadFile(FilePath)
	if err != nil{
		return nil, err
	}

	// Create a new Config struct to hold parsed configuraion settings
	config :=&Config{}

	// Unmarshal the YAML daya into the config struct
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}


// UpdateConfigWithEnv updates the configuration settings with environment variables.
// It modifies the provided config struct in-place
func UpdateConfigWithEnv(cfg *Config) error {
	if cfg.Server.Port == "" {
		cfg.Server.Port = os.Getenv("SERVER_PORT")
	}

	if cfg.DB.Host == "" {
		cfg.DB.Host = os.Getenv("DB_HOST")
	}

	if cfg.DB.Port == 0 {
		portStr := os.Getenv("DB_PORT")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("failed to parse DB_PORT: %w", err)
		}
		cfg.DB.Port = port
	}

	if cfg.DB.User == "" {
		cfg.DB.User = os.Getenv("DB_USER")
	}

	if cfg.DB.Password == "" {
		cfg.DB.Password = os.Getenv("DB_PASSWORD")
	}

	if cfg.DB.DBName == "" {
		cfg.DB.DBName = os.Getenv("DB_NAME")
	}

	if cfg.DB.SSLMode == "" {
		cfg.DB.SSLMode = os.Getenv("DB_SSL_MODE")
	}

	// Update other configuration values from environment variables

	return nil
}

// parseInt parses the string representation of an integer and returns the corresponding integer value.
// It returns an error if the parsing fails.
func parseInt(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Failed to parse integer: %v", err)
	}
	return value
}
