package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// wrap to Tokens as per the request body of the API
type TokenRequest struct {
	TokenAddr string `json:"token_address"`
	TokenID   string `json:"token_id"`
}

// request body complete
// WATCH FOR SPELLING!!!!!
type MultipleNFTRequest struct {
	Tokens            []TokenRequest `json:"tokens"`
	NormalizeMetadata bool           `json:"normalizeMetadata"`
	MediaItems        bool           `json:"media_items"`
}

// most of what will be returned already covered in the Transactions struct
func (c *Client) nftByTokenAddr(tokens []TokenRequest, normalizeMetadata, mediaItems bool) ([]Transactions, error) {
	// POST /nft/getMultipleNFTs
	baseURL := strings.TrimSuffix(c.BaseUrl, "/")
	nftURL := fmt.Sprintf("%s/nft/getMultipleNFTs", baseURL)

	// create request body
	reqBody := MultipleNFTRequest{
		Tokens:            tokens,
		NormalizeMetadata: normalizeMetadata,
		MediaItems:        mediaItems,
	}

	// marshal JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// create HTTP request
	req, err := http.NewRequest("POST", nftURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	// set headers
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// add chain query params
	// TODO: proper array of chain (you know it)
	query := req.URL.Query()
	query.Add("chain", "ronin")
	req.URL.RawQuery = query.Encode()

	// make HTTP request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API failed with  status %d: %s", res.StatusCode, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// endpoint returns an array as opposed to an object with an array inside
	// so return the data directly as thus
	var nftData []Transactions
	if err := json.Unmarshal(body, &nftData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nftData, nil
}
