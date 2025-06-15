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

		// Action flags
		// Subject to changes
		//
		fetchTx  = flag.Bool("tx", false, "Fetch transactions (default action if no other specified)")
		fetchNft = flag.Bool("nft", false, "Fetch NFTs")

		// Query parameter flags (used by Wallet and NFT)
		//
		limit           = flag.Int("limit", 10, "Limit number of transactions (default: 10)")
		cursor          = flag.String("cursor", "", "Cursor for pagination")
		order           = flag.String("order", "DESC", "Order: ASC or DESC (default: DESC)")
		fromDate        = flag.String("from-date", "", "From date (format: seconds or datestring)")
		toDate          = flag.String("to-date", "", "To date (format: seconds or datestring)")
		includeInternal = flag.Bool("include-internal", false, "Include internal transactions")
		nftMetadata     = flag.Bool("nft-metadata", false, "Include NFT metadata")

		// NFT-specific flags
		format            = flag.String("format", "", "Response format")
		excludeSpam       = flag.Bool("exclude-spam", false, "Exclude spam NFTs")
		normalizeMetadata = flag.Bool("normalize-metadata", false, "Normalize metadata")
		mediaItems        = flag.Bool("media-items", false, "Include media items")
		includePrices     = flag.Bool("include-prices", false, "Include price information")
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

	// determine what action to take
	// if no specific action is specified, default to tx
	//
	if !*fetchNft && !*fetchTx {
		*fetchTx = true
	}

	// build query params
	//
	queryParams := internal.QueryParams{
		Limit:                       *limit,
		Cursor:                      getCursorPtr(*cursor),
		Order:                       *order,
		FromDate:                    *fromDate,
		ToDate:                      *toDate,
		IncludeInternalTransactions: *includeInternal,
		NftMetadata:                 *nftMetadata,
		Format:                      *format,
		ExcludeSpam:                 *excludeSpam,
		NormalizMetadata:            *normalizeMetadata,
		MediaItems:                  *mediaItems,
		IncludePrices:               *includePrices,
	}

	// execute actions
	//
	if *fetchTx {
		log.Printf("Fetching transaction stats for wallet: %s", address)
		if err := internal.GetTxByWallet(baseURL, key, address, queryParams); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	if *fetchNft {
		log.Printf("Fetching NFTs for wallet: %s", address)
		if err := internal.GetNftsByWallet(baseURL, key, address, queryParams); err != nil {
			log.Fatalf("Error fetching NFTs: %v", err)
		}
	}
}

// helper func to convert flag value to a ptr
func getCursorPtr(cursor string) *string {
	if cursor == "" {
		return nil // not provided
	}
	return &cursor // provided even if empty after trimming
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
