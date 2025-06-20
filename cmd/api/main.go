package main

import (
	"cmd/internal/client"
	"cmd/internal/commands"
	"cmd/internal/config"
	"cmd/internal/models"
	"cmd/internal/service"
	"cmd/pkg/logger"
	"cmd/pkg/utils"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dotenv-org/godotenvvault"
)

func main() {
	// Init. logger based on environment
	logLevel := os.Getenv("LOG_LEVEL")
	var log *logger.Logger

	switch logLevel {
	case "debug":
		log = logger.NewWithLevel(slog.LevelDebug)
	case "warn":
		log = logger.NewWithLevel(slog.LevelWarn)
	case "error":
		log = logger.NewWithLevel(slog.LevelError)
	default:
		log = logger.NewWithLevel(slog.LevelInfo)
	}

	// Log application start
	log.Info("Starting NFT CLI application",
		"version", "1.0.0",
		"log_level", logLevel,
	)

	// load env file
	if err := godotenvvault.Load(); err != nil {
		log.Warn("No env file found, using system environment variables")
	}

	// load configs
	// TODO: change log to slog (DONE)
	// TODO: update logger with relevant changes
	cfg, err := config.Load()
	if err != nil {
		log.Error(
			"Failed to load configuration",
			"error", err,
		)
		os.Exit(1)
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

	// Log parsed args
	log.Debug("CLI flags parsed",
		"wallet_address", *walletAddr,
		"fetch_nft", *fetchNFT,
		"fetch_specific", *fetchSpecific,
		"limit", *limit,
	)

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
			log.Info("Executing NFT wallet command",
				"wallet_address", finalWalletAddr,
			)
		}
		if err := nftCommand.GetByWallet(ctx, finalWalletAddr, params); err != nil {
			log.Error("NFT wallet command failed",
				"error", err,
				"wallet_address", finalWalletAddr,
			)
			os.Exit(1)
		}
		log.Info("NFT wallet command completed successfully")
	} else if *fetchSpecific {
		var tokens []models.TokenRequest

		if *tokensFile != "" {
			// Load from file
			tokens, err = utils.LoadTokensFromFile(*tokensFile)
			if err != nil {
				log.Error(
					"Failed to load tokens file/missing JSON file",
					"error", err,
				)
			}
		} else if *tokenAddr != "" && *tokenID != "" {
			// Use single token from flags
			tokens = []models.TokenRequest{
				{
					TokenAddress: *tokenAddr,
					TokenID:      *tokenID,
				},
			}
		} else {
			log.Error("Need either -tokens-file or both -token-address and -token-id")
		}

		if err := nftCommand.GetSpecific(ctx, tokens); err != nil {
			log.Error(
				"Failed to get NFT, check if the tokens are correct",
				"Error:", err)
		}
	} else {
		flag.Usage()
	}
}
