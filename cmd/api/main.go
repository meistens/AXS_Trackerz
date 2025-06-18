package main

import (
	"cmd/internal/client"
	"cmd/internal/commands"
	"cmd/internal/config"
	"cmd/internal/models"
	"cmd/internal/service"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dotenv-org/godotenvvault"
)

func main() {
	// load env file
	if err := godotenvvault.Load(); err != nil {
		log.Println("No env file found (using system env)")
	}

	// load configs
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Loading config:", err)
	}

	// parse CLI flags
	var (
		walletAddr    = flag.String("wallet", "", "Wallet address")                // wallet addr.
		tokenAddr     = flag.String("token-address", "", "Token contract address") // token addr.
		tokenID       = flag.String("token-id", "", "Token ID")                    // token id
		tokensFile    = flag.String("tokens-file", "", "File with tokens JSON")    // token file, containing token addr. and token id
		limit         = flag.Int("limit", 10, "Limit results")                     // query param
		excludeSpam   = flag.Bool("exclude-spam", false, "Exclude spam")           // query param
		fetchNFT      = flag.Bool("nft", false, "Fetch NFTs by wallet")            // get NFT by wallet addr.
		fetchSpecific = flag.Bool("specific-nft", false, "Fetch specific NFTs")    // get metadata for NFTs
	)
	flag.Parse()

	// set up deps
	moralisClient := client.NewMoralisClient(cfg.MoralisAPIKey, cfg.MoralisBaseURL, cfg.WalletAddress)
	nftService := service.NewNFTService(moralisClient)
	nftCommand := commands.NewNFTCommand(nftService)

	// set up ctx for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// goroutine listens for ctrl + c signal from the terminal
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// build query params
	params := models.QueryParams{
		Limit:       *limit,
		ExcludeSpam: *excludeSpam,
	}

	// determine wallet address: CLI overrides env
	var finalWalletAddr string
	if *walletAddr != "" {
		finalWalletAddr = *walletAddr // CLI flag provided
	} else {
		finalWalletAddr = cfg.WalletAddress // use from env
	}

	// Execute commands
	if *fetchNFT {
		if finalWalletAddr == "" {
			log.Fatal("Wallet address required: use -wallet flag or set WALLET_ADDRESS env")
		}
		if err := nftCommand.GetByWallet(ctx, finalWalletAddr, params); err != nil {
			log.Fatal("Error:", err)
		}
	} else if *fetchSpecific {
		var tokens []models.TokenRequest

		if *tokensFile != "" {
			// Load from file
			tokens, err = LoadTokensFromFile(*tokensFile)
			if err != nil {
				log.Fatal("Loading tokens file:", err)
			}
		} else if *tokenAddr != "" && *tokenID != "" {
			// Use single token from flags
			tokens = []models.TokenRequest{
				{TokenAddress: *tokenAddr, TokenID: *tokenID},
			}
		} else {
			log.Fatal("Need either -tokens-file or both -token-address and -token-id")
		}

		if err := nftCommand.GetSpecific(ctx, tokens); err != nil {
			log.Fatal("Error:", err)
		}
	} else {
		flag.Usage()
	}
} // ← This closing brace was missing!

// LoadTokensFromFile loads tokens from a JSON file
// This is now a standalone function, not a method
func LoadTokensFromFile(filename string) ([]models.TokenRequest, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	var tokens []models.TokenRequest
	if err := json.NewDecoder(file).Decode(&tokens); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	return tokens, nil
}
