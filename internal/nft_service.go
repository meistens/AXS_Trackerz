package internal

import (
	"fmt"
	"net/http"
)

// function fetches NFTs for a wallet address
// and displays them
func GetNftsByWallet(baseURL, apiKey, walletAddr string, params QueryParams) error {
	// Create client
	//
	client := &Client{
		HTTPClient: &http.Client{},
		BaseUrl:    baseURL,
		APIKey:     apiKey,
	}

	// Fetch NFT data
	//
	nftData, err := client.getNftsByAddress(walletAddr, params)
	if err != nil {
		return fmt.Errorf("failed to fetch NFTs: %w", err)
	}

	// Extract simplified data
	//
	extractedData := extractNFTData(nftData)

	// Display results (output can be customized)
	//
	fmt.Printf("Found %d NFTs for wallet %s\n", len(extractedData), walletAddr)

	for i, nft := range extractedData {
		fmt.Printf("NFT #%d:\n", i+1)
		fmt.Printf("  Token ID: %s\n", nft.TokenID)
		fmt.Printf("  Name: %s\n", nft.Name)
		fmt.Printf("  Contract: %s\n", nft.TokenAddress)
		fmt.Printf("  Floor Price: %s %s\n", nft.FloorPrice, nft.FloorPriceCurrency)
		if nft.RarityRank != nil {
			fmt.Printf("  Rarity Rank: %d\n", *nft.RarityRank)
		}
		fmt.Printf("  Verified: %t\n", nft.IsVerified)
		fmt.Printf("  Possible Spam: %t\n", nft.PossibleSpam)

		// Show some attributes if available
		if len(nft.Attributes) > 0 {
			fmt.Printf("  Attributes:\n")
			count := 0
			for trait, value := range nft.Attributes {
				if count >= 3 { // Limit to first 3 attributes for readability
					fmt.Printf("    ... and %d more\n", len(nft.Attributes)-3)
					break
				}
				fmt.Printf("    %s: %v\n", trait, value)
				count++
			}
		}
		fmt.Println()
	}
	return nil
}
