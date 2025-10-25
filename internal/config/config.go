package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	DB_URL                string
	FIPE_EXTERNAL_API_URL string
	FIPE_API_TOKEN        string
	SecretKey             []byte
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	return &Config{
		Port:                  getEnv("PORT", "8080"),
		DB_URL:                getEnv("DB_URL", "postgres://user:pass@localhost:5432/dbname?sslmode=disable"),
		FIPE_EXTERNAL_API_URL: getEnv("FIPE_EXTERNAL_API_URL", ""),
		FIPE_API_TOKEN:        getEnv("FIPE_API_TOKEN", ""),
		SecretKey:             []byte(getEnv("SECRET_KEY", "")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
