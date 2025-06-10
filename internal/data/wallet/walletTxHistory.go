package wallet

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Moralis Wallet History API Response Structure
type APIResponse struct {
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
	Cursor   string   `json:"cursor"`
	Result   []Result `json:"result"`
}

type Result struct {
	Hash                     string `json:"hash"`
	Nonce                    string `json:"nonce"`
	TransactionIndex         string `json:"transaction_index"`
	FromAddress              string `json:"from_address"`
	ToAddress                string `json:"to_address"`
	Value                    string `json:"value"`
	Gas                      string `json:"gas"`
	GasPrice                 string `json:"gas_price"`
	GasUsed                  string `json:"gas_used"`
	CumulativeGasUsed        string `json:"cumulative_gas_used"`
	InputData                string `json:"input"`
	ReceiptContractAddress   string `json:"receipt_contract_address"`
	ReceiptCumulativeGasUsed string `json:"receipt_cumulative_gas_used"`
	ReceiptGasUsed           string `json:"receipt_gas_used"`
	ReceiptStatus            string `json:"receipt_status"`
	BlockTimestamp           string `json:"block_timestamp"`
	BlockNumber              string `json:"block_number"`
	BlockHash                string `json:"block_hash"`
}

// What you want to display
type TxDetails struct {
	TransactionHash string `json:"transaction_hash"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	Value           string `json:"value"`
	BlockTimestamp  string `json:"block_timestamp"`
}

// Client struct
type Client struct {
	HTTPClient *http.Client
	BaseUrl    string
	APIKey     string
}

// function creates a new wallet API client
func NewClient(baseurl string, apikey string) *Client {
	return &Client{
		HTTPClient: &http.Client{},
		BaseUrl:    baseurl,
		APIKey:     apikey,
	}
}

// function fetches wallet transaction history
func (c *Client) getTokensByWallet(walletAddr string) ([]*TxDetails, error) {
	// Build the correct URL for wallet history
	baseURL := strings.TrimSuffix(c.BaseUrl, "/tokens")
	historyURL := fmt.Sprintf("%s/%s/history", baseURL, walletAddr)

	req, err := http.NewRequest("GET", historyURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.APIKey)
	query := req.URL.Query()
	query.Add("chain", "ronin")
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
func displayStats(allTxs []*TxDetails) {
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

// export zone
func GetTxByWallet(baseURL string, apiKey string, walletAddr string) error {
	client := NewClient(baseURL, apiKey)
	txs, err := client.getTokensByWallet(walletAddr)
	if err != nil {
		return fmt.Errorf("failed to fetch transaction history: %w", err)
	}
	displayStats(txs)
	return nil
}
