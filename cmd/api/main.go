package main

import (
	"flag"
	"log"
	"os"

	"github.com/dotenv-org/godotenvvault"

	"cmd/internal/data/wallet"
)

func main() {
	// load env file
	if err := godotenvvault.Load(); err != nil {
		log.Println("no env file found")
	}

	// cli flags
	var (
		// Change this line in your main.go:
		apiURL        = flag.String("url", "", "API base URL")
		apiKey        = flag.String("key", "", "API key")
		walletAddress = flag.String("wallet", "", "Wallet address to query")
	)
	flag.Parse()

	// get configs from env or flags
	// empty string reads automagically env
	baseURL := getConfig(*apiURL, "URI_RON", "")
	key := getConfig(*apiKey, "API", "")
	address := getConfig(*walletAddress, "WALLET_ADDRESS", "")

	// validate
	if key == "" {
		log.Fatal("API key is required. Set API_KEY environment variable or use -key flag")
	}

	if address == "" {
		log.Fatal("Wallet address is required. Set WALLET_ADDRESS environment variable or use -wallet flag")
	}

	// Fetch and display transaction stats
	log.Printf("Fetching transaction stats for wallet: %s", address)
	if err := wallet.GetTxByWallet(baseURL, key, address); err != nil {
		log.Fatalf("Error: %v", err)
	}

}

// getConfig returns the first non-empty value from flag, environment variable, or default
func getConfig(flagValue, envKey, defaultValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}
	return defaultValue
}
