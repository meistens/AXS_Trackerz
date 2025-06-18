package models

// NFT represents a single NFT with the data we care about
type NFT struct {
	TokenID      string                 `json:"token_id"`
	TokenAddress string                 `json:"token_address"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Image        string                 `json:"image"`
	FloorPrice   string                 `json:"floor_price"`
	IsVerified   bool                   `json:"is_verified"`
	PossibleSpam bool                   `json:"possible_spam"`
	Attributes   map[string]interface{} `json:"attributes"`
	RarityRank   *int                   `json:"rarity_rank,omitempty"`
}

// TokenRequest is what we send to get specific NFTs
type TokenRequest struct {
	TokenAddress string `json:"token_address"`
	TokenID      string `json:"token_id"`
}

// QueryParams are the filters we can use when getting NFTs
type QueryParams struct {
	Limit         int     `json:"limit"`
	Cursor        *string `json:"cursor"`
	ExcludeSpam   bool    `json:"exclude_spam"`
	IncludePrices bool    `json:"include_prices"`
}

// APIResponse is what Moralis sends back for wallet queries
type APIResponse struct {
	Status   string       `json:"status"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
	Result   []RawNFTData `json:"result"`
}

// RawNFTData struct is the raw data from the Moralis API
// what you feel is necessary
type RawNFTData struct {
	TokenID            string              `json:"token_id"`
	TokenAddress       string              `json:"token_address"`
	Name               string              `json:"name"`
	OwnerOf            string              `json:"owner_of"`
	FloorPrice         *string             `json:"floor_price,omitempty"`
	VerifiedCollection bool                `json:"verified_collection"`
	PossibleSpam       bool                `json:"possible_spam"`
	NormalizedMetadata *NormalizedMetadata `json:"normalized_metadata"`
	RarityRank         *int                `json:"rarity_rank,omitempty"`
	Symbol             string              `json:"symbol"`
	// ... other fields from the Moralis API data
}

type NormalizedMetadata struct {
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	Image       string       `json:"image"`
	Attributes  *[]Attribute `json:"attributes,omitempty"`
}

type Attribute struct {
	TraitType string `json:"trait_type"`
	Value     any    `json:"value"`
}
