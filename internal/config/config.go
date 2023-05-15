package config

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"log"
	"os"

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

// LoadConfif, reads the config file and unmarshals it into a Config struct.
// It takes a string aurgument confifPath representing the path to the configuration file.
func LoadEncryptedConfig(configFile, key string) (*Config, error) {
	// Read the contents of the config file
	encryptedData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}

	decryptedData, err := decrypt(encryptedData, key)
	if err != nil {
		return nil, err
	}

	// Create a new Config struct to hold the parsed configuration settings.
	config := &Config{}
	// Unmarshal the file contents (YAML) into the config struct.
	err = yaml.Unmarshal(decryptedData, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func decrypt(encrypted []byte, key string) ([]byte, error) {
	ciphertext, _ := hex.DecodeString(string(encrypted))

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}
