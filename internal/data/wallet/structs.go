package wallet

import "net/http"

// Moralis Wallet API Response Structures here (could be moved once I get other API responses as I expplore)
// Separated for ease of access to changes
// If this is wrong in terms of project structure, I do not care
// What I care about is making sure I can easily track what what is doing than
// ruin my mental state later down the line as the project grows in complexity
//

// requests params to be sent to the API
type QueryParams struct {
	Limit                       int    `url:"limit,omitempty"`
	Cursor                      string `url:"cursor,omitempty"`
	Order                       string `url:"order,omitempty"`
	FromDate                    string `url:"from_date,omitempty"`
	ToDate                      string `url:"to_date,omitempty"`
	IncludeInternalTransactions bool   `url:"include_internal_transactions,omitempty"`
	NftMetadata                 bool   `url:"nft_metadata,omitempty"`
}

// expected responses from the API, split
// quite time-consuming to do, use AI, review and compare from actual
// response saves time than typing it all
type APIResponse struct {
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
	Cursor   string         `json:"-"`
	Result   []Transactions `json:"result"` // sake of readability
}

type Transactions struct {
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	TransactionIndex  string `json:"transaction_index"`
	FromAddress       string `json:"from_address"`
	GasUsed           string `json:"gas_used"`
	CumulativeGasUsed string `json:"cumulative_gas_used"`
	InputData         string `json:"input"`
	ContractType      string `json:"contract_type"`
	TransactionHash   string `json:"transaction_hash"`
	TransactionType   string `json:"transaction_type"`
	TokenAddress      string `json:"token_address"`
	TokenID           string `json:"token_id"`
	//Amount                   string           `json:"string"`
	//Verified                 int              `json:"verified"`
	//Operator                 string           `json:"operator"`
	//VerifiedCollection       bool             `json:"verified_collection"`
	FromAddressEntity        *string               `json:"from_address_entity,omitempty"`
	FromAddressEntityLogo    *string               `json:"from_address_entity_logo,omitempty"`
	FromAddressLabel         *string               `json:"from_address_label,omitempty"`
	ToAddressEntity          *string               `json:"to_address_entity,omitempty"`
	ToAddressEntityLogo      *string               `json:"to_address_entity_logo,omitempty"`
	ToAddress                string                `json:"to_address"`
	ToAddressLabel           *string               `json:"to_address_label,omitempty"`
	Value                    string                `json:"value"`
	Gas                      string                `json:"gas"`
	GasPrice                 string                `json:"gas_price"`
	ReceiptCumulativeGasUsed string                `json:"receipt_cumulative_gas_used"`
	ReceiptGasUsed           string                `json:"receipt_gas_used"`
	ReceiptContractAddress   *string               `json:"receipt_contract_address,omitempty"`
	ReceiptStatus            string                `json:"receipt_status"`
	BlockTimestamp           string                `json:"block_timestamp"`
	BlockNumber              string                `json:"block_number"`
	BlockHash                string                `json:"block_hash"`
	TransactionFee           string                `json:"transaction_fee"`
	MethodLabel              *string               `json:"method_label,omitempty"`
	NFTTransfers             []NFTTransfer         `json:"nft_transfers"`
	ERC20Transfers           []ERC20Transfer       `json:"erc20_transfers"`
	NativeTransfers          []NativeTransfer      `json:"native_transfers"`
	Summary                  string                `json:"summary"`
	PossibleSpam             bool                  `json:"possible_spam"`
	Category                 string                `json:"category"`
	InternalTransactions     []InternalTransaction `json:"internal_transactions,omitempty"`
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

type NormalizedMetadata struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	AnimationURL *string     `json:"animation_url,omitempty"`
	ExternalLink *string     `json:"external_link,omitempty"`
	Image        string      `json:"image"`
	Attributes   []Attribute `json:"attributes"`
}

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

type NativeTransfer struct {
	// Add fields as needed when native transfers are present
}

// Client struct
type Client struct {
	HTTPClient *http.Client
	BaseUrl    string
	APIKey     string
}

// walletTxHistory.go
type TxDetails struct {
	TransactionHash string `json:"transaction_hash"`
	FromAddress     string `json:"from_address"`
	ToAddress       string `json:"to_address"`
	Value           string `json:"value"`
	BlockTimestamp  string `json:"block_timestamp"`
}
