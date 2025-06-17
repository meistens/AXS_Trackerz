package internal

import (
	"net/http"
)

type APIResponse struct {
	Status   string         `json:"status"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Cursor   string         `json:"-"`      // leaving as is
	Result   []Transactions `json:"result"` // contains the main stuff
}

// transactions, more to add, and a mix of wallet and NFT
// comment every 7th for readability
type Transactions struct {
	Amount           string `json:"amount"`
	TokenID          string `json:"token_id"`
	TokenAddress     string `json:"token_address"`
	ContractType     string `json:"contract_type"`
	OwnerOf          string `json:"owner_of"`
	LastMetadataSync string `json:"last_metadata_sync"`
	LastTokenUriSync string `json:"last_token_uri_sync"`

	// skipped metadata, will add later after tackling other stuff...

	BlockNumber       string  `json:"block_number"`
	BlockNumberMinted *string `json:"block_number_minted,omitempty"` // string or null
	Name              string  `json:"name"`
	Symbol            string  `json:"symbol"`
	TokenHash         string  `json:"token_hash"`
	TokenURI          string  `json:"token_uri"`
	MinterAddress     *string `json:"minter_address,omitempty"`

	//

	RarityRank         *int                `json:"rarity_rank,omitempty"`
	RarityPercentage   *float64            `json:"rarity_percentage,omitempty"`
	RarityLabel        *string             `json:"rarity_label,omitempty"`
	VerifiedCollection bool                `json:"verified_collection"`
	PossibleSpam       bool                `json:"possible_spam"`
	LastSale           *int                `json:"last_sale,omitempty"`           // *any, come back to this when data is available from other sources
	NormalizedMetadata *NormalizedMetadata `json:"normalized_metadata,omitempty"` // object with data

	//

	Media                 Media  `json:"media"`
	CollectionLogo        string `json:"collection_logo"`
	CollectionBannerImage string `json:"collection_banner_image"`
	CollectionCategory    string `json:"collection_category"`
	ProjectURL            string `json:"project_url"`
	WikiURL               string `json:"wiki_url"`
	DiscordURL            string `json:"discord_url"`

	//

	TelegramURL        string    `json:"telegram_url"`
	TwitterUsername    string    `json:"twitter_username"`
	InstagramUsername  string    `json:"instagram_username"`
	ListPrice          ListPrice `json:"list_price"`
	FloorPrice         *string   `json:"floor_price,omitempty"`
	FloorPriceUSD      *string   `json:"floor_price_usd,omitempty"`
	FloorPriceCurrency *string   `json:"floor_price_currency,omitempty"`
	//==================================================================
	//
	//

	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	TransactionIndex  string `json:"transaction_index"`
	FromAddress       string `json:"from_address"`
	GasUsed           string `json:"gas_used"`
	CumulativeGasUsed string `json:"cumulative_gas_used"`
	InputData         string `json:"input"`

	TransactionHash string `json:"transaction_hash"`
	TransactionType string `json:"transaction_type"`

	//Verified                 int              `json:"verified"`
	//Operator                 string           `json:"operator"`
	//VerifiedCollection       bool             `json:"verified_collection"`
	FromAddressEntity        *string `json:"from_address_entity,omitempty"`
	FromAddressEntityLogo    *string `json:"from_address_entity_logo,omitempty"`
	FromAddressLabel         *string `json:"from_address_label,omitempty"`
	ToAddressEntity          *string `json:"to_address_entity,omitempty"`
	ToAddressEntityLogo      *string `json:"to_address_entity_logo,omitempty"`
	ToAddress                string  `json:"to_address"`
	ToAddressLabel           *string `json:"to_address_label,omitempty"`
	Value                    string  `json:"value"`
	Gas                      string  `json:"gas"`
	GasPrice                 string  `json:"gas_price"`
	ReceiptCumulativeGasUsed string  `json:"receipt_cumulative_gas_used"`
	ReceiptGasUsed           string  `json:"receipt_gas_used"`
	ReceiptContractAddress   *string `json:"receipt_contract_address,omitempty"`
	ReceiptStatus            string  `json:"receipt_status"`
	BlockTimestamp           string  `json:"block_timestamp"`

	BlockHash       string           `json:"block_hash"`
	TransactionFee  string           `json:"transaction_fee"`
	MethodLabel     *string          `json:"method_label,omitempty"`
	NFTTransfers    []NFTTransfer    `json:"nft_transfers"`
	ERC20Transfers  []ERC20Transfer  `json:"erc20_transfers"`
	NativeTransfers []NativeTransfer `json:"native_transfers"`
	Summary         string           `json:"summary"`

	Category             string                `json:"category"`
	InternalTransactions []InternalTransaction `json:"internal_transactions,omitempty"`
}

//

type NormalizedMetadata struct {
	Name         string       `json:"name"`
	Description  *string      `json:"description,omitempty"`
	AnimationURL *string      `json:"animation_url,omitempty"`
	ExternalLink *string      `json:"external_link,omitempty"`
	Image        string       `json:"image"`
	Attributes   *[]Attribute `json:"attributes,omitempty"` // array of data, same keys, different values, can be null
}

//

type Attribute struct {
	TraitType   string   `json:"trait_type"`
	Value       any      `json:"value"` // Can be string or number
	DisplayType *string  `json:"display_type,omitempty"`
	MaxValue    *int     `json:"max_value,omitempty"`
	TraitCount  int      `json:"trait_count"`
	Order       *int     `json:"order,omitempty"`
	RarityLabel *string  `json:"rarity_label,omitempty"`
	Count       *int     `json:"count,omitempty"`
	Percentage  *float64 `json:"percentage,omitempty"`
}

//

type Media struct {
	Status           string           `json:"status"`
	UpdatedAt        string           `json:"updatedAt"`
	Mimetype         *string          `json:"mimetype,omitempty"`
	ParentHash       *string          `json:"parent_hash,omitempty"`
	MediaCollection  *MediaCollection `json:"media_collection,omitempty"`
	OriginalMediaURL string           `json:"original_media_url"`
}

type MediaCollection struct {
	Low    MediaSize `json:"low"`
	Medium MediaSize `json:"medium"`
	High   MediaSize `json:"high"`
}

type MediaSize struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}

type ListPrice struct {
	Listed        bool    `json:"listed"`
	Price         *string `json:"price,omitempty"`
	PriceCurrency *string `json:"price_currency,omitempty"`
	PriceUSD      *string `json:"price_usd,omitempty"`
	Marketplace   *string `json:"marketplace,omitempty"`
}

// NFTExtract - Simplified struct with only the data you need
type NFTExtract struct {
	TokenID            string                 `json:"token_id"`
	Name               string                 `json:"name"`
	Owner              string                 `json:"owner"`
	TokenAddress       string                 `json:"token_address"`
	ContractType       string                 `json:"contract_type"`
	FloorPrice         string                 `json:"floor_price"`
	FloorPriceCurrency string                 `json:"floor_price_currency"`
	Image              string                 `json:"image"`
	Description        string                 `json:"description"`
	Attributes         map[string]interface{} `json:"attributes"`
	CollectionName     string                 `json:"collection_name"`
	IsVerified         bool                   `json:"is_verified"`
	PossibleSpam       bool                   `json:"possible_spam"`
	RarityRank         *int                   `json:"rarity_rank,omitempty"`
	LastSalePrice      *string                `json:"last_sale_price,omitempty"`
	BlockNumber        string                 `json:"block_number"`
}

//
//

type QueryParams struct {
	// Chain  string  `json:"chain"`
	Limit  int     `url:"limit,omitempty"`
	Cursor *string `url:"cursor,omitempty"` // adjusted for null return values
	// Order                       string  `url:"order,omitempty"`
	FromDate                    string `url:"from_date,omitempty"`
	ToDate                      string `url:"to_date,omitempty"`
	IncludeInternalTransactions bool   `url:"include_internal_transactions,omitempty"`
	NftMetadata                 bool   `url:"nft_metadata,omitempty"`
	Format                      string `json:"format"`
	ExcludeSpam                 bool   `json:"exclude_spam"`
	IncludePrices               bool   `json:"include_prices"`
	NormalizMetadata            bool   `json:"nomalize_metadata"`
	MediaItems                  bool   `json:"media_items"`
}

type InternalTransaction struct {
	// if anything new pops up while exploring API, worthwhile or not, add
	// who knows when it will be useful...
}

type ERC20Transfer struct {
	TokenName             string  `json:"token_name"`
	TokenSymbol           string  `json:"token_symbol"`
	TokenLogo             string  `json:"token_logo"`
	TokenDecimals         string  `json:"token_decimals"`
	FromAddressEntity     *string `json:"from_address_entity,omitempty"`
	FromAddressEntityLogo *string `json:"from_address_entity_logo,omitempty"`
	FromAddress           string  `json:"from_address"`
	FromAddressLabel      *string `json:"from_address_label,omitempty"`
	ToAddressEntity       *string `json:"to_address_entity,omitempty"`
	ToAddressEntityLogo   *string `json:"to_address_entity_logo,omitempty"`
	ToAddress             string  `json:"to_address"`
	ToAddressLabel        *string `json:"to_address_label,omitempty"`
	Address               string  `json:"address"`
	LogIndex              int     `json:"log_index"`
	Value                 string  `json:"value"`
	PossibleSpam          bool    `json:"possible_spam"`
	VerifiedContract      bool    `json:"verified_contract"`
	SecurityScore         *int    `json:"security_score,omitempty"`
	Direction             string  `json:"direction"`
	ValueFormatted        string  `json:"value_formatted"`
}

type NFTTransfer struct {
	LogIndex              int                 `json:"log_index"`
	Value                 string              `json:"value"`
	ContractType          string              `json:"contract_type"`
	TransactionType       string              `json:"transaction_type"`
	TokenAddress          string              `json:"token_address"`
	TokenID               string              `json:"token_id"`
	FromAddressEntity     *string             `json:"from_address_entity,omitempty"`
	FromAddressEntityLogo *string             `json:"from_address_entity_logo,omitempty"`
	FromAddress           string              `json:"from_address"`
	FromAddressLabel      *string             `json:"from_address_label,omitempty"`
	ToAddressEntity       *string             `json:"to_address_entity,omitempty"`
	ToAddressEntityLogo   *string             `json:"to_address_entity_logo,omitempty"`
	ToAddress             string              `json:"to_address"`
	ToAddressLabel        *string             `json:"to_address_label,omitempty"`
	Amount                string              `json:"amount"`
	Operator              string              `json:"operator"`
	PossibleSpam          bool                `json:"possible_spam"`
	VerifiedCollection    bool                `json:"verified_collection"`
	Direction             string              `json:"direction"`
	CollectionLogo        string              `json:"collection_logo"`
	CollectionBannerImage string              `json:"collection_banner_image"`
	NormalizedMetadata    *NormalizedMetadata `json:"normalized_metadata,omitempty"`
}

type NativeTransfer struct {
	// Add fields as needed when native transfers are present
}

// Client struct
type Client struct {
	HTTPClient *http.Client
	BaseUrl    string
	APIKey     string
}

// wallet_main.go
type TxDetails struct {
	TransactionHash string `json:"transaction_hash"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	Value           string `json:"value"`
	BlockTimestamp  string `json:"block_timestamp"`
}

// nft_service.go
// type nftByAddr struct {
// 	Amount       string `json:"amount"`
// 	TokenID      string `json:"token_id"`
// 	TokenAddress string `json:"token_address"`
// 	ContractType string `json:"contract_type"`
// }

// ERC20 Token Transfers By Wallet
type Erc20Tokens struct {
	TokenName        string `json:"token_name"`
	TokenSymbol      string `json:"token_symbol"`
	FromAddress      string `json:"from_address"`
	ToAddress        string `json:"to_address"`
	BlockTimestamp   string `json:"block_timestamp"`
	ValueDecimal     string `json:"value_decimal"`
	VerifiedContract bool   `json:"verified_contract"`
}

// Add these new structs referenced in NFTData
// Modifiable?
// type LastSale struct {
// 	Price     *string `json:"price"`
// 	Currency  *string `json:"currency"`
// 	PriceUSD  *string `json:"price_usd"`
// 	Timestamp *string `json:"timestamp"`
// }
