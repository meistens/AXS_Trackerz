package main

import (
	"flag"
	"log"
	"os"

	"github.com/dotenv-org/godotenvvault"

	"cmd/internal"
)

func main() {
	// load env file
	if err := godotenvvault.Load(); err != nil {
		log.Println("no env file found")
	}

	// cli flags
	var (
		apiURL        = flag.String("url", "", "API base URL")
		apiKey        = flag.String("key", "", "API key")
		walletAddress = flag.String("wallet", "", "Wallet address to query")
		// Query parameter flags
		limit           = flag.Int("limit", 10, "Limit number of transactions (default: 10)")
		cursor          = flag.String("cursor", "", "Cursor for pagination")
		order           = flag.String("order", "DESC", "Order: ASC or DESC (default: DESC)")
		fromDate        = flag.String("from-date", "", "From date (format: seconds or datestring)")
		toDate          = flag.String("to-date", "", "To date (format: seconds or datestring)")
		includeInternal = flag.Bool("include-internal", false, "Include internal transactions")
		nftMetadata     = flag.Bool("nft-metadata", false, "Include NFT metadata")
	)

	//
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
	queryParams := internal.QueryParams{
		Limit:                       *limit,
		Cursor:                      *cursor,
		Order:                       *order,
		FromDate:                    *fromDate,
		ToDate:                      *toDate,
		IncludeInternalTransactions: *includeInternal,
		NftMetadata:                 *nftMetadata,
	}

	if err := internal.GetTxByWallet(baseURL, key, address, queryParams); err != nil {
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
