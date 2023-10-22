package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBHost          string
	DBPort          int
	DBUser          string
	DBPassword      string
	DBName          string
	HttpPort        string
	HttpHost        string
	AgifyHost       string
	GenderizeHost   string
	NationalizeHost string
}

func LoadConfig() (Config, error) {

	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Error loading .env file")
			panic(err)
		}
	}

	DbPort, err := strconv.Atoi(os.Getenv("EM_DB_PORT"))
	if err != nil {
		return Config{}, fmt.Errorf("error loading db port: %w", err)
	}

	config := Config{
		DBHost:          os.Getenv("EM_DB_HOST"),
		DBPort:          DbPort,
		DBUser:          os.Getenv("EM_DB_USER"),
		DBPassword:      os.Getenv("EM_DB_PASSWORD"),
		DBName:          os.Getenv("EM_DB_NAME"),
		HttpPort:        os.Getenv("EM_HTTP_PORT"),
		HttpHost:        os.Getenv("EM_HTTP_HOST"),
		AgifyHost:       os.Getenv("EM_AGIFY_API_HOST"),
		GenderizeHost:   os.Getenv("EM_GENDERIZE_API_HOST"),
		NationalizeHost: os.Getenv("EM_NATIONALIZE_API_HOST"),
	}

	log.Printf("config: %#v\n", config)
	return config, nil
}
