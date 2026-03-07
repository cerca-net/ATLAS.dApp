#!/bin/bash

# FlutterFlow Integration Test Script
# This script tests all the FlutterFlow-specific API endpoints

echo "🧪 Testing FlutterFlow Blockchain Integration"
echo "=============================================="

# Base URL for the blockchain API
BASE_URL="http://localhost:8080"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to test API endpoint
test_endpoint() {
    local name="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    
    echo -e "\n${YELLOW}Testing: $name${NC}"
    echo "URL: $method $url"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
    fi
    
    # Extract status code (last line)
    status_code=$(echo "$response" | tail -n1)
    # Extract response body (all lines except last)
    body=$(echo "$response" | head -n -1)
    
    if [ "$status_code" = "200" ]; then
        echo -e "${GREEN}✅ SUCCESS${NC}"
        echo "Response: $body" | jq '.' 2>/dev/null || echo "Response: $body"
    else
        echo -e "${RED}❌ FAILED (Status: $status_code)${NC}"
        echo "Response: $body"
    fi
}

# Check if blockchain is running
echo "🔍 Checking if blockchain is running..."
if curl -s "$BASE_URL/status" > /dev/null; then
    echo -e "${GREEN}✅ Blockchain is running${NC}"
else
    echo -e "${RED}❌ Blockchain is not running. Please start it first:${NC}"
    echo "go run main.go --port=8000 --api=8080"
    exit 1
fi

# Test 1: Connect Wallet (Create)
echo -e "\n${YELLOW}=== Test 1: Connect Wallet (Create) ===${NC}"
test_endpoint "Connect Wallet" "POST" "$BASE_URL/flutterflow/connect-wallet" '{"action": "create"}'

# Extract wallet address from response for subsequent tests
WALLET_RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d '{"action": "create"}' "$BASE_URL/flutterflow/connect-wallet")
WALLET_ADDRESS=$(echo "$WALLET_RESPONSE" | jq -r '.data.address' 2>/dev/null)
SESSION_TOKEN=$(echo "$WALLET_RESPONSE" | jq -r '.data.sessionToken' 2>/dev/null)

if [ "$WALLET_ADDRESS" != "null" ] && [ "$WALLET_ADDRESS" != "" ]; then
    echo -e "${GREEN}✅ Wallet created: ${WALLET_ADDRESS:0:20}...${NC}"
else
    echo -e "${RED}❌ Failed to create wallet${NC}"
    exit 1
fi

# Test 2: Authenticate Session
echo -e "\n${YELLOW}=== Test 2: Authenticate Session ===${NC}"
test_endpoint "Authenticate" "POST" "$BASE_URL/flutterflow/authenticate" "{\"sessionToken\": \"$SESSION_TOKEN\", \"address\": \"$WALLET_ADDRESS\"}"

# Test 3: Get Wallet Info
echo -e "\n${YELLOW}=== Test 3: Get Wallet Info ===${NC}"
test_endpoint "Wallet Info" "GET" "$BASE_URL/flutterflow/wallet-info?address=$WALLET_ADDRESS"

# Test 4: Send Transaction
echo -e "\n${YELLOW}=== Test 4: Send Transaction ===${NC}"
test_endpoint "Send Transaction" "POST" "$BASE_URL/flutterflow/send-transaction" "{\"from\": \"$WALLET_ADDRESS\", \"to\": \"0x1234567890abcdef1234567890abcdef12345678\", \"amount\": 50, \"fee\": 1, \"signature\": \"TODO: In test setup, generate a real signature for the transaction using the test wallet's private key\", \"sessionToken\": \"$SESSION_TOKEN\"}"

# Test 5: Get Transaction History
echo -e "\n${YELLOW}=== Test 5: Transaction History ===${NC}"
test_endpoint "Transaction History" "GET" "$BASE_URL/flutterflow/transaction-history?address=$WALLET_ADDRESS"

# Test 6: Disconnect Wallet
echo -e "\n${YELLOW}=== Test 6: Disconnect Wallet ===${NC}"
test_endpoint "Disconnect Wallet" "POST" "$BASE_URL/flutterflow/disconnect" "{\"sessionToken\": \"$SESSION_TOKEN\", \"address\": \"$WALLET_ADDRESS\"}"

# Test 7: Faucet (Get Test Tokens)
echo -e "\n${YELLOW}=== Test 7: Faucet (Get Test Tokens) ===${NC}"
test_endpoint "Faucet" "POST" "$BASE_URL/faucet" "{\"address\": \"$WALLET_ADDRESS\"}"

# Test 8: Import Wallet
echo -e "\n${YELLOW}=== Test 8: Import Wallet ===${NC}"
test_endpoint "Import Wallet" "POST" "$BASE_URL/flutterflow/connect-wallet" "{\"action\": \"import\", \"privateKey\": \"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef\"}"

echo -e "\n${GREEN}🎉 FlutterFlow Integration Tests Complete!${NC}"
echo -e "\n${YELLOW}Next Steps:${NC}"
echo "1. Use these endpoints in your FlutterFlow app"
echo "2. Replace 'localhost:8080' with your production URL"
echo "3. Implement proper error handling in FlutterFlow"
echo "4. Add input validation for user inputs"
echo -e "\n${GREEN}Happy integrating! 🚀${NC}" 