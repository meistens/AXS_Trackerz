package client

import (
	"bytes"
	"cmd/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// MoralisClient struct responsible for 'talking' to the Moralis API
type MoralisClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	walletAddr string
}

// NewMoralisClient func creates a new client
func NewMoralisClient(apiKey, baseURL, walletAddr string) *MoralisClient {
	return &MoralisClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

// GetNFTsByWallet gets all NFTs for a wallet
func (c *MoralisClient) GetNFTsByWallet(ctx context.Context, walletAddr string, params models.QueryParams) ([]models.RawNFTData, error) {
	// Build URL
	// Format: baseURL/{address}/nft
	url := fmt.Sprintf("%s/%s/nft", strings.TrimSuffix(c.baseURL, "/"), walletAddr)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request(from client/moralis_client.go/del. later): %w", err)
	}

	// Set header
	req.Header.Set("X-API-Key", c.apiKey)

	// Add query params
	// TODO: handle multiple chains
	query := req.URL.Query()
	query.Add("chain", "ronin")
	if params.Limit > 0 {
		query.Add("limit", fmt.Sprintf("%d", params.Limit))
	}
	if params.ExcludeSpam {
		query.Add("exclude_spam", "true")
	}
	req.URL.RawQuery = query.Encode()

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request(from client/moralis_client): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d(from client/morails_client.go)", resp.StatusCode)
	}

	// Parse response, we need the queries sent
	var apiResp models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("parsing response (from client/moralis_client): %w", err)
	}

	// return the result of the queries
	return apiResp.Result, nil
}

// GetSpecificNFTs func gets specific NFTs by token address and token ID
// Explanation -> Takes the token address and token ID as params, which are in the TokenRequest struct
// Return -> Data from the RawNFTData struct (pick whichever you fancy)
func (c *MoralisClient) GetSpecificNFTs(ctx context.Context, tokens []models.TokenRequest) ([]models.RawNFTData, error) {
	// Format: baseURL/nft/getMultipleNFTs
	url := fmt.Sprintf("%s/nft/getMultipleNFTs", strings.TrimSuffix(c.baseURL, "/"))

	// Create request body
	// since this is different and we have to add the address and id to the body
	// instead of query params
	reqBody := map[string]interface{}{
		"tokens":            tokens,
		"normalizeMetadata": true,  // TODO: change to false to see what it being returned
		"media_items":       false, // TODO: change to true if you plan to use for GUI or want to seee the media url
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling request(client/moralis_client): %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("creating request(client/moralis_client): %w", err)
	}

	// Add headers
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Add query params
	query := req.URL.Query()
	query.Add("chain", "ronin") // TODO: add multiple chain
	req.URL.RawQuery = query.Encode()

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request(client/moralis_client): %w)", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Since this POST does things different - no query params - and it returns
	// the array instead of an object of arrays, return the data directly
	// without parsing
	var nftData []models.RawNFTData
	if err := json.NewDecoder(resp.Body).Decode(&nftData); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	return nftData, nil
}
