package config

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
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

func Encrypt(plainData []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	plainData = pad(plainData)

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plainData))
	mode.CryptBlocks(ciphertext, plainData)

	ciphertext = append(iv, ciphertext...)

	return []byte(hex.EncodeToString(ciphertext)), nil
}

func pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// LoadConfif, reads the config file and unmarshals it into a Config struct.
// It takes a string aurgument confifPath representing the path to the configuration file.
func LoadEncryptedConfig(configFile, key string) (*Config, error) {
	// Read the contents of the config file
	encryptedData, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("Checking if the error is here")
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

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// Unpad data after decryption.
	ciphertext, err = unpad(ciphertext)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// unpad removes padding from data that was added as per PKCS#7 standard.
func unpad(data []byte) ([]byte, error) {
	length := len(data)
	unpadding := int(data[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when wrong encryption key is used")
	}

	return data[:(length - unpadding)], nil
}
