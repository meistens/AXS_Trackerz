package internal

import (
	"fmt"
	"net/http"
)

// GetTxByWallet fetches and displays wallet transaction history
func GetTxByWallet(baseURL string, apiKey string, walletAddr string, params QueryParams) error {
	//	Create client directly here to avoid import cycle
	client := &Client{
		HTTPClient: &http.Client{},
		BaseUrl:    baseURL,
		APIKey:     apiKey,
	}

	txs, err := client.GetTokensByWallet(walletAddr, params)
	if err != nil {
		return fmt.Errorf("failed to fetch transaction history: %w", err)
	}

	DisplayStats(txs)
	return nil
}
