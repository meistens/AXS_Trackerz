package config

import (
	"log"
	"os"
)

// holds all app settings
type Config struct {
	MoralisAPIKey  string
	MoralisBaseURL string
	Port           string
	LogLevel       string
	WalletAddress  string
	// CacheEnabled   bool
}

// Load func. reads config from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		MoralisAPIKey:  getEnv("MORALIS_API_KEY", ""),
		MoralisBaseURL: getEnv("MORALIS_BASE_URL", "https://deep-index.moralis.io/api/v2.2"),
		WalletAddress:  getEnv("WALLET_ADDRESS", ""),
		Port:           getEnv("PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		//CacheEnabled:   getEnvAsBool("CACHE_ENABLED", false),
	}

	// check if required fields are set
	if cfg.MoralisAPIKey == "" {
		log.Fatal("MORALIS_API_KEY is required")
	}

	return cfg, nil
}

// getEnv func. fetches environment variables or returns default set in Load()
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
