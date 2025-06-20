package service

import (
	"cmd/internal/client"
	"cmd/internal/models"
	"cmd/pkg/logger"
	"context"
	"fmt"
	"time"
)

// NFTService struct handles NFT operations
type NFTService struct {
	moralisClient *client.MoralisClient
	logger        *logger.Logger
}

// NewNFTService func creates a new service
func NewNFTService(client *client.MoralisClient) *NFTService {
	log := logger.New().WithGroup("nft_service")

	return &NFTService{
		moralisClient: client,
		logger:        log,
	}
}

// GetNFTsByWallet (see client/moralis_client for func.)
// Explanation -> func gets NFTs for a wallet, cleans up the data
// Return -> NFT data
func (c *NFTService) GetNFTsByWallet(ctx context.Context, walletAddr string, params models.QueryParams) ([]models.NFT, error) {
	start := time.Now()

	c.logger.Info("Processing NFT wallet request",
		"wallet_address", walletAddr,
		"params", params,
	)

	// get raw data from API
	rawNFTs, err := c.moralisClient.GetNFTsByWallet(ctx, walletAddr, params)
	if err != nil {
		c.logger.Error("Failed to fetch NFTs from API",
			"error", err,
			"wallet_address", walletAddr,
		)
		return nil, fmt.Errorf("fetching NFTs from API: %w", err)
	}

	// Convert raw data to clean data
	cleanNFTs := c.convertRawNFTs(rawNFTs)

	// log results
	duration := time.Since(start)
	c.logger.Info("NFT wallet request processed",
		"wallet_address", walletAddr,
		"raw_nfts", len(rawNFTs),
		"clean_nfts", len(cleanNFTs),
		"duration", duration,
	)

	return cleanNFTs, nil
}

// GetSpecificNFTs (see client/moralis_client for func.)
// Explanation -> func gets specific NFT based on token ID and token address provided
// Return -> NFT data
func (c *NFTService) GetSpecficNFTs(ctx context.Context, tokens []models.TokenRequest) ([]models.NFT, error) {
	start := time.Now()

	// log service call
	c.logger.Info("Getting NFT",
		"tokens", tokens,
	)

	rawNFTs, err := c.moralisClient.GetSpecificNFTs(ctx, tokens)
	if err != nil {
		c.logger.Error("Failed to fetch NFT from API",
			"error", err,
			"tokens", tokens,
		)
		return nil, fmt.Errorf("fetching specific NFTs from API: %w", err)
	}

	specificNFT := c.convertRawNFTs(rawNFTs)

	// Log results
	duration := time.Since(start)
	c.logger.Info("NFT wallet request processed",
		"tokens", tokens,
		"specific_nfts", len(specificNFT),
		"duration", duration,
	)
	return specificNFT, nil
}

// convertRawNFTs
// Explanation -> func cleans up the raw data received from the API to clean data
// Return -> cleaned NFT data
func (c *NFTService) convertRawNFTs(rawNFTs []models.RawNFTData) []models.NFT {
	var cleanNFTs []models.NFT
	spamCount := 0

	for _, raw := range rawNFTs {
		// skip spam count
		if raw.PossibleSpam {
			spamCount++
			continue
		}

		nft := models.NFT{
			TokenID:      raw.TokenID,
			TokenAddress: raw.TokenAddress,
			Name:         raw.Name,
			IsVerified:   raw.VerifiedCollection,
			PossibleSpam: raw.PossibleSpam,
			RarityRank:   raw.RarityRank,
		}

		// floor price, in case it returns nil/null
		if raw.FloorPrice != nil {
			nft.FloorPrice = *raw.FloorPrice
		}

		// metadata, in case it returns nil
		if raw.NormalizedMetadata != nil {
			nft.Image = raw.NormalizedMetadata.Image

			if raw.NormalizedMetadata.Description != nil {
				nft.Description = *raw.NormalizedMetadata.Description
			}

			// convert attriutes to map
			nft.Attributes = make(map[string]interface{})
			if raw.NormalizedMetadata.Attributes != nil {
				for _, attr := range *raw.NormalizedMetadata.Attributes {
					nft.Attributes[attr.TraitType] = attr.Value
				}
			}
		}
		cleanNFTs = append(cleanNFTs, nft)
	}

	// Log conversion stats
	if spamCount > 0 {
		c.logger.Info("Filtered spam NFTs",
			"spam_fltered", spamCount,
			"clean_nfts", len(cleanNFTs),
		)
	}
	return cleanNFTs
}
