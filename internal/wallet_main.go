package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// function fetches wallet transaction history
func (c *Client) GetTokensByWallet(walletAddr string, params QueryParams) ([]*TxDetails, error) {

	// Build the correct URL for wallet history
	baseURL := strings.TrimSuffix(c.BaseUrl, "/")
	historyURL := fmt.Sprintf("%s/wallets/%s/history", baseURL, walletAddr)

	req, err := http.NewRequest("GET", historyURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.APIKey)
	query := req.URL.Query()
	query.Add("chain", "ronin")

	// Add all query parameters
	if params.Limit > 0 {
		query.Add("limit", fmt.Sprintf("%d", params.Limit))
	}
	if params.Cursor != "" {
		query.Add("cursor", params.Cursor)
	}
	if params.Order != "" {
		query.Add("order", params.Order)
	}
	if params.FromDate != "" {
		query.Add("from_date", params.FromDate)
	}
	if params.ToDate != "" {
		query.Add("to_date", params.ToDate)
	}
	if params.IncludeInternalTransactions {
		query.Add("include_internal_transactions", "true")
	}
	if params.NftMetadata {
		query.Add("nft_metadata", "true")
	}

	req.URL.RawQuery = query.Encode()

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make a request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API failed with status %d: %s", res.StatusCode, res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Debug: Print raw response
	// fmt.Println("Raw API Response:")
	// fmt.Println(string(body))
	// fmt.Println("---")

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert all results to TxDetails
	var allTxs []*TxDetails
	for _, tx := range apiResponse.Result {
		txDetail := &TxDetails{
			TransactionHash: tx.Hash,
			FromAddress:     tx.FromAddress,
			ToAddress:       tx.ToAddress,
			Value:           tx.Value,
			BlockTimestamp:  tx.BlockTimestamp,
		}
		allTxs = append(allTxs, txDetail)
	}

	return allTxs, nil
}

// display all transaction history
func DisplayStats(allTxs []*TxDetails) {
	if len(allTxs) == 0 {
		fmt.Println("No transactions found")
		return
	}

	fmt.Printf("=== Wallet Transaction History ===\n")
	for i, tx := range allTxs {
		fmt.Printf("--- Transaction %d ---\n", i+1)
		fmt.Printf("Transaction Hash: %s\n", tx.TransactionHash)
		fmt.Printf("From Address:     %s\n", tx.FromAddress)
		fmt.Printf("To Address:       %s\n", tx.ToAddress)
		fmt.Printf("Value:            %s\n", tx.Value)
		fmt.Printf("Block Timestamp:  %s\n", tx.BlockTimestamp)
		fmt.Printf("\n")
	}
	fmt.Printf("==================================\n")
}
