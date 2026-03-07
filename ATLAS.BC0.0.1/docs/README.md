# Simple Blockchain Project

## Overview
A minimal blockchain implementation in Go with a web-based frontend. Features include block/transaction management, staking, validator rotation, tokenomics, state snapshots, and a testnet faucet.

## Features

### Core Blockchain
- ECDSA wallet generation and transaction signing (client-side)
- Block and transaction validation
- Staking and validator management
- Tokenomics (mint, burn, fees, rewards)
- State snapshots and recovery
- Peer-to-peer networking
- REST API and web dashboard
- Faucet for test tokens

### Advanced Features
- **Smart Contracts**: Custom VM with upgradable and formally verified contracts
- **Privacy Features**: Encryption, zero-knowledge proofs, and GDPR compliance
- **On-chain Governance**: Proposal and voting system with stake-weighted voting
- **Sharding Architecture**: Horizontal scaling with cross-shard transactions
- **Real-time Monitoring**: Comprehensive monitoring dashboard with alerts
- **Oracle Integration**: External data integration for smart contracts
- **Dynamic Fees**: Flexible transaction fee system
- **KYC for Validators**: Identity verification for validator registration

## Setup
### Prerequisites
- Go 1.20+
- Node.js (for frontend development, optional)

### Backend
1. Clone the repo.
2. Run: `go run main.go --port=8000 --api=8080`
3. For multi-node: use different ports on each machine.

### Frontend
1. Open `frontend/index.html` in your browser.
2. The dashboard will connect to the backend API (default: `localhost:8080`).

## Usage
- Generate/import/export wallet in the dashboard.
- Request faucet tokens to your address.
- Send transactions and view blockchain state.
- Connect multiple nodes for network testing.

## Architecture
- **Go Backend:** Handles blockchain logic, consensus, state, and API.
- **Frontend:** Pure JS/HTML, manages wallet and signs transactions in-browser.
- **API:** REST endpoints for blocks, transactions, balances, validators, peers, and faucet.

## Testing with a Partner
- Run backend on both computers (different ports).
- Open frontend on each, use faucet to get tokens.
- Send transactions between addresses.
- Check block/transaction propagation and state sync.

## Security Note
**Private keys are never sent to the backend.** All signing is done in-browser. Never share your private key.

## License
MIT 

## vm/ Directory (Custom Virtual Machine)

This directory contains the implementation of the custom virtual machine (VM) for executing Turing-complete smart contracts on the blockchain. The VM is designed to support contract deployment, execution, upgradability, and formal verification. It will be integrated with the transaction and block processing pipeline, enabling programmable logic and advanced features for the blockchain. 