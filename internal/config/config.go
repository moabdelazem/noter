package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
}

// Load the all the configs and the env vars
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	return &Config{
		ServerPort: getEnv("PORT", "8080"),
	}, nil
}

// Get any env var from the .env file by key
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
