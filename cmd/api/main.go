package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
		fetchTx          = flag.Bool("tx", false, "Fetch transactions (default action if no other specified)")
		fetchNft         = flag.Bool("nft", false, "Fetch NFTs")
		fetchSpecificNft = flag.Bool("specific-nft", false, "Fetch specific NFTs by token address and ID")

		// specific NFT flags
		tokenAddress = flag.String("token-address", "", "Token contract address (required with -specific-nft)")
		tokenID      = flag.String("token-id", "", "Token ID (required with -specific-nft)")
		tokensFile   = flag.String("tokens-file", "", "JSON file containing array of tokens (alternative to single token)")

		// Query parameter flags (used by Wallet and NFT)
		//
		limit  = flag.Int("limit", 10, "Limit number of transactions (default: 10)")
		cursor = flag.String("cursor", "", "Cursor for pagination")
		//	order           = flag.String("order", "DESC", "Order: ASC or DESC (default: DESC)")
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

	// validate based on action
	if *fetchSpecificNft {
		if *tokensFile == "" && (*tokenAddress == "" || *tokenID == "") {
			log.Fatal("For specific NFT fetching, provide token address and token id")
		}
		if *tokensFile != "" && (*tokenAddress != "" || *tokenID != "") {
			log.Fatal("Use either tokensFile or individual token flags")
		}
	} else {
		// For wallet-based queries, address is required
		if address == "" {
			log.Fatal("Wallet address is required. Set WALLET_ADDRESS environment variable or use -wallet flag")
		}
	}

	// determine what action to take
	// if no specific action is specified, default to tx
	//
	if !*fetchNft && !*fetchTx && !*fetchSpecificNft {
		*fetchTx = true
	}

	// build query params
	//
	queryParams := internal.QueryParams{
		Limit:  *limit,
		Cursor: getCursorPtr(*cursor),
		//Order:                       *order,
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

	if *fetchSpecificNft {
		log.Printf("Fetching specific NFTs...")

		var tokens []internal.TokenRequest
		var err error

		if *tokensFile != "" {
			// load tokens from JSON file
			tokens, err = loadTokensFromFile(*tokensFile)
			if err != nil {
				log.Fatalf("Error loading from file: %v", err)
			}
		} else {
			// use from flags if no JSON file found
			tokens = []internal.TokenRequest{
				{
					TokenAddr: *tokenAddress,
					TokenID:   *tokenID,
				},
			}
		}

		if err := internal.GetSpecificNFTs(baseURL, key, tokens, *normalizeMetadata, *mediaItems); err != nil {
			log.Fatalf("Error fetching specific NFTs: %v", err)
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

// loadTokensFromFile loads tokens from a JSON file
func loadTokensFromFile(filename string) ([]internal.TokenRequest, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var tokens []internal.TokenRequest
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tokens); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return tokens, nil
}
