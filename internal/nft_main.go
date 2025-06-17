package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) getNftsByAddress(walletAddr string, params QueryParams) ([]Transactions, error) {
	// GET baseURL/{address}/nft
	//
	baseURL := strings.TrimSuffix(c.BaseUrl, "/")
	nftURL := fmt.Sprintf("%s/%s/nft", baseURL, walletAddr)

	// Create HTTP Request
	//
	req, err := http.NewRequest("GET", nftURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request (NFT, delete once done with all): %w", err)
	}

	// Set Headers and Whatnot
	//
	req.Header.Set("X-API-KEY", c.APIKey)
	query := req.URL.Query()
	query.Add("chain", "ronin")

	// Add Query Parameters
	// see swagger docs for references
	//
	if params.Format != "" {
		query.Add("format", params.Format)
	}
	if params.Limit > 0 {
		query.Add("limit", fmt.Sprintf("%d", params.Limit))
	}
	if params.ExcludeSpam {
		query.Add("exclude_spam", "true")
	}
	if params.NormalizMetadata {
		query.Add("normalize_metadata", "true")
	}
	if params.MediaItems {
		query.Add("media_items", "false")
	}
	if params.IncludePrices {
		query.Add("include_prices", "true")
	}

	//
	req.URL.RawQuery = query.Encode()

	// Make HTTP Request
	//
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make a request: %w", err)
	}
	defer res.Body.Close()

	//
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API failed with status %d: %s", res.StatusCode, res.Status)
	}

	// Read HTTP Response
	//
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON Data
	//
	var nftResponse APIResponse
	if err := json.Unmarshal(body, &nftResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	// return NFT API response, cherry-pick later
	return nftResponse.Result, nil
}

// Function extracts NFT Data From the Full NFT Response
// Pass NFT Data, Returns With Data Strictly Needed
func extractNFTData(nfts []Transactions) []NFTExtract {
	var extractedNFTData []NFTExtract

	// loop through the NFTData struct
	//
	for _, nft := range nfts {
		// loop and convert attributes to a map for easy access (nested objects taken care of?)
		//
		attribs := make(map[string]any)
		if nft.NormalizedMetadata != nil && nft.NormalizedMetadata.Attributes != nil {
			for _, attrib := range *nft.NormalizedMetadata.Attributes {
				attribs[attrib.TraitType] = attrib.Value
			}
		}
		// Get image URL
		//
		var imgURL string
		if nft.NormalizedMetadata != nil {
			imgURL = nft.NormalizedMetadata.Image
		}

		// Get description
		//
		var description string
		if nft.NormalizedMetadata != nil && nft.NormalizedMetadata.Description != nil {
			description = *nft.NormalizedMetadata.Description
		}

		// handle unsafe ptr derefs
		var floorPrice, floorPriceCurrency string
		if nft.FloorPrice != nil {
			floorPrice = *nft.FloorPrice
		}

		if nft.FloorPriceCurrency != nil {
			floorPriceCurrency = *nft.FloorPriceCurrency
		}
		// if nft.LastSale != nil && nft.LastSale.Price != nil {
		// 	lastSalePrice = *nft.LastSale.Price
		// }

		extractedNFTData = append(extractedNFTData, NFTExtract{
			TokenID:            nft.TokenID,
			Name:               nft.Name,
			Owner:              nft.OwnerOf,
			TokenAddress:       nft.TokenAddress,
			ContractType:       nft.ContractType,
			FloorPrice:         floorPrice,
			FloorPriceCurrency: floorPriceCurrency,
			Image:              imgURL,
			Description:        description,
			Attributes:         attribs,
			CollectionName:     nft.Name, // or use a different field if available
			IsVerified:         nft.VerifiedCollection,
			PossibleSpam:       nft.PossibleSpam,
			RarityRank:         nft.RarityRank,
			BlockNumber:        nft.BlockNumber,
		})
	}
	return extractedNFTData
}
