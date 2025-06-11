# Yes, generated with AI. If there's a bulldozer around to make short work of digging, why use a shovel???

# Some Blockchain - And More - API

A Go-based command-line application for fetching and displaying blockchain stuff - and more on the way - using the Moralis and Ronin (not there yet) API.

## 🚧 Work in Progress

This project is currently under active development. Features and documentation are subject to change.

## Overview

This tool allows you to query the blockchain on the Ronin network. It fetches transaction data from the Moralis API and displays it in a readable format.

## Features

- Support for Ronin blockchain network
- Environment variable and command-line flag configuration
- Clean, formatted transaction output

## Project Structure (More In The Works)

```
.
├── cmd/
│   └── api/
│       └── main.go              # Main application entry point
└── internal/
    └── data/
        └── wallet/
            └── walletTxHistory.go   # Wallet API client and transaction handling
```

## Current Functionality

### Transaction Data Retrieved
- Transaction hash
- From/To addresses
- Transaction value

### Configuration Options
The application supports configuration through:
- Command-line flags
- Environment variables

## Dependencies

- `github.com/dotenv-org/godotenvvault` - Environment variable management
- Standard Go libraries (`flag`, `log`, `os`, `net/http`, `encoding/json`)

## API Integration

Currently integrated with Moralis API for blockchain data:
- Supports Ronin chain queries
- Includes error handling

## Development Status

TBA

## Notes

This is an experimental project focused on blockchain transaction analysis. The codebase is being actively refined and expanded.

---

*More documentation and usage instructions will be added as development progresses.*
