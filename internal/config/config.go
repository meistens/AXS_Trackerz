package config

import (
	"cmd/pkg/logger"
	"os"
)

type Config struct {
	// API Configuration
	MoralisAPIKey  string
	MoralisBaseURL string

	// Server Configuration
	Port     string
	LogLevel string

	// Database
	DatabaseURL string

	// Discord
	DiscordToken    string
	DiscordClientID string

	// Payment
	ETHWalletAddress string
	BTCWalletAddress string

	// Security
	JWTSecret string

	// Optional defaults
	WalletAddress string
	TokenAddress  string
}

func Load() (*Config, error) {
	log := logger.New()

	cfg := &Config{
		MoralisAPIKey:    getEnv("MORALIS_API_KEY", ""),
		MoralisBaseURL:   getEnv("MORALIS_BASE_URL", "https://deep-index.moralis.io/api/v2.2"),
		Port:             getEnv("PORT", "8080"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		DiscordToken:     getEnv("DISCORD_BOT_TOKEN", ""),
		DiscordClientID:  getEnv("DISCORD_CLIENT_ID", ""),
		ETHWalletAddress: getEnv("ETH_WALLET_ADDRESS", ""),
		BTCWalletAddress: getEnv("BTC_WALLET_ADDRESS", ""),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		WalletAddress:    getEnv("WALLET_ADDRESS", ""),
		TokenAddress:     getEnv("TOKEN_ADDRESS", ""),
	}

	// Log configuration loading with structured data
	log.Info("configuration loaded",
		"moralis_base_url", cfg.MoralisBaseURL,
		"port", cfg.Port,
		"log_level", cfg.LogLevel,
	)

	// Check if requiired fields are set
	if cfg.MoralisAPIKey == "" {
		log.Error("Missing required configuration",
			"field", "MORALIS_API_KEY",
			"message", "MORALIS_API_KEY environment variable is required",
		)
		os.Exit(1)
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
