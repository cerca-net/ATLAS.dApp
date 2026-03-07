# FlutterFlow Integration with Atlas Blockchain

This document explains how to integrate FlutterFlow applications with the Atlas Blockchain.

## Overview

The Atlas Blockchain provides dedicated API endpoints for FlutterFlow integration, making it easy to connect your FlutterFlow frontend with the Atlas blockchain backend. These endpoints handle wallet management, authentication, transactions, and more.

## Prerequisites

1. Running Atlas Blockchain node with API enabled
2. FlutterFlow project
3. Basic understanding of REST APIs

## Available FlutterFlow API Endpoints

The Atlas Blockchain provides the following FlutterFlow-specific endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/flutterflow/connect-wallet` | POST | Create, import, or connect to a wallet |
| `/flutterflow/authenticate` | POST | Authenticate a session |
| `/flutterflow/wallet-info` | GET | Get wallet information |
| `/flutterflow/send-transaction` | POST | Send a transaction |
| `/flutterflow/transaction-history` | GET | Get transaction history |
| `/flutterflow/disconnect` | POST | Disconnect wallet |

## Integration Steps

### 1. Start Atlas Blockchain Node

First, start the Atlas Blockchain node with API enabled:

```bash
cd "ATLAS.BC 0.0.1"
go run cmd/main.go -api=8080
```

The API will be available at `http://localhost:8080`.

### 2. Add API Service to FlutterFlow

Create a new Dart file in your FlutterFlow project (e.g., `atlas_api_service.dart`) with functions to call the Atlas API endpoints. You can use the provided example in the cercaend project as a reference.

### 3. Use API Functions in Your Widgets

Import and use the API functions in your FlutterFlow widgets to interact with the blockchain.

## API Endpoint Details

### Connect Wallet

**Endpoint**: `POST /flutterflow/connect-wallet`

Connect to a wallet by creating a new one, importing an existing one, or connecting to an address.

Request body:
```json
{
  "action": "create|import|connect",
  "privateKey": "private_key_if_importing", // Optional
  "address": "address_if_connecting" // Optional
}
```

Response:
```json
{
  "success": true,
  "message": "Wallet connected successfully",
  "data": {
    "address": "wallet_address",
    "sessionToken": "authentication_token",
    "balance": 1000,
    "isValidator": false
  }
}
```

### Authenticate

**Endpoint**: `POST /flutterflow/authenticate`

Authenticate a session with a session token.

Request body:
```json
{
  "sessionToken": "token_from_connect",
  "address": "wallet_address"
}
```

Response:
```json
{
  "success": true,
  "message": "Authentication successful",
  "data": {
    "address": "wallet_address",
    "balance": 1000,
    "isValidator": false,
    "validatorInfo": {}
  }
}
```

### Get Wallet Info

**Endpoint**: `GET /flutterflow/wallet-info?address=wallet_address`

Get information about a wallet.

Response:
```json
{
  "success": true,
  "data": {
    "address": "wallet_address",
    "balance": 1000,
    "isValidator": false,
    "validatorInfo": {},
    "recentTransactions": [],
    "nonce": 5
  }
}
```

### Send Transaction

**Endpoint**: `POST /flutterflow/send-transaction`

Send a transaction from one wallet to another.

Request body:
```json
{
  "from": "sender_address",
  "to": "recipient_address",
  "amount": 100,
  "fee": 1,
  "data": "optional_data",
  "signature": "transaction_signature",
  "sessionToken": "authentication_token"
}
```

Response:
```json
{
  "success": true,
  "message": "Transaction submitted successfully",
  "data": {
    "transactionHash": "hash_of_transaction",
    "from": "sender_address",
    "to": "recipient_address",
    "amount": 100,
    "status": "pending"
  }
}
```

### Get Transaction History

**Endpoint**: `GET /flutterflow/transaction-history?address=wallet_address`

Get transaction history for a wallet.

Response:
```json
{
  "success": true,
  "data": {
    "address": "wallet_address",
    "transactions": [
      {
        "hash": "transaction_hash",
        "blockIndex": 123,
        "sender": "sender_address",
        "recipient": "recipient_address",
        "amount": 100,
        "fee": 1,
        "timestamp": 1234567890,
        "type": "sent|received"
      }
    ],
    "totalCount": 1
  }
}
```

### Disconnect

**Endpoint**: `POST /flutterflow/disconnect`

Disconnect a wallet session.

Request body:
```json
{
  "sessionToken": "authentication_token",
  "address": "wallet_address"
}
```

Response:
```json
{
  "success": true,
  "message": "Wallet disconnected successfully"
}
```

## Example Implementation

See the `cercaend/lib/backend/api_requests/atlas_api_service.dart` file for a complete Dart implementation of all API endpoints.

## Testing

To test the integration:

1. Start the Atlas Blockchain node
2. Run your FlutterFlow app
3. Use the wallet connection functionality to create or connect to a wallet
4. Perform transactions and check balances

## Troubleshooting

### Common Issues

1. **CORS errors**: Make sure the Atlas Blockchain node is configured to allow requests from your FlutterFlow app's domain.

2. **Connection refused**: Ensure the Atlas Blockchain node is running and accessible from your FlutterFlow app.

3. **Authentication failures**: Make sure you're using valid session tokens and addresses.

### Debugging Tips

1. Check the Atlas Blockchain logs for error messages.
2. Use browser developer tools to inspect network requests.
3. Verify that all required fields are included in API requests.

## Security Considerations

1. Never expose private keys in the frontend.
2. Use HTTPS in production environments.
3. Validate all inputs on both client and server sides.
4. Store session tokens securely.