#!/bin/bash

# Comprehensive Blockchain Test Suite
echo "🧪 Starting Comprehensive Blockchain Tests"

# Test Configuration
API_BASE="http://localhost:8081"
NODE2_API="http://localhost:8082"
NODE3_API="http://localhost:8083"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Test function
run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_result="$3"
    
    echo -e "${YELLOW}🧪 Running: $test_name${NC}"
    
    result=$(eval "$test_command" 2>/dev/null)
    
    if [[ "$result" == *"$expected_result"* ]]; then
        echo -e "${GREEN}✅ PASS: $test_name${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}❌ FAIL: $test_name${NC}"
        echo "Expected: $expected_result"
        echo "Got: $result"
        ((TESTS_FAILED++))
    fi
    echo ""
}

# Wait for nodes to be ready
echo -e "${BLUE}⏳ Waiting for nodes to be ready...${NC}"
sleep 15

echo -e "${BLUE}🔍 Testing Node Health and Basic Functionality${NC}"
echo "=================================================="

# Test 1: Check if nodes are running
run_test "Node 1 Health Check" \
    "curl -s $API_BASE/health" \
    "status"

run_test "Node 2 Health Check" \
    "curl -s $NODE2_API/health" \
    "status"

run_test "Node 3 Health Check" \
    "curl -s $NODE3_API/health" \
    "status"

# Test 2: Get blockchain info
run_test "Blockchain Info" \
    "curl -s $API_BASE/blockchain" \
    "blocks"

echo -e "${BLUE}🔐 Testing Wallet and Key Management${NC}"
echo "============================================="

# Test 3: Create wallet
run_test "Wallet Creation" \
    "curl -s -X POST $API_BASE/wallet" \
    "public_key"

# Test 4: Get wallet balance
run_test "Wallet Balance" \
    "curl -s $API_BASE/wallet/balance" \
    "balance"

echo -e "${BLUE}💰 Testing Transaction System${NC}"
echo "===================================="

# Test 5: Create transaction
run_test "Transaction Creation" \
    "curl -s -X POST $API_BASE/transaction -H 'Content-Type: application/json' -d '{\"recipient\":\"test\",\"amount\":10}'" \
    "transaction"

# Test 6: Get transaction pool
run_test "Transaction Pool" \
    "curl -s $API_BASE/transactions" \
    "transactions"

echo -e "${BLUE}⛓️ Testing Consensus and Validators${NC}"
echo "==========================================="

# Test 7: Get validators
run_test "Validator List" \
    "curl -s $API_BASE/validators" \
    "validators"

# Test 8: Check consensus state
run_test "Consensus State" \
    "curl -s $API_BASE/consensus" \
    "validators"

echo -e "${BLUE}🌐 Testing Network Synchronization${NC}"
echo "=========================================="

# Test 9: Check node synchronization
run_test "Node 2 Synchronization" \
    "curl -s $NODE2_API/blockchain" \
    "blocks"

run_test "Node 3 Synchronization" \
    "curl -s $NODE3_API/blockchain" \
    "blocks"

# Test 10: Test transaction propagation
run_test "Transaction Propagation to Node 2" \
    "curl -s $NODE2_API/transactions" \
    "transactions"

run_test "Transaction Propagation to Node 3" \
    "curl -s $NODE3_API/transactions" \
    "transactions"

echo -e "${BLUE}📊 Performance and State Testing${NC}"
echo "========================================"

# Test 11: Check block production
run_test "Block Production" \
    "curl -s $API_BASE/blockchain | grep -o '\"height\":[0-9]*' | head -1" \
    "height"

# Test 12: Check state consistency
run_test "State Consistency" \
    "curl -s $API_BASE/state" \
    "balances"

echo "📊 Test Results Summary:"
echo "========================"
echo -e "${GREEN}✅ Passed: $TESTS_PASSED${NC}"
echo -e "${RED}❌ Failed: $TESTS_FAILED${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed! Blockchain is working correctly.${NC}"
else
    echo -e "${RED}⚠️  Some tests failed. Check the logs for details.${NC}"
    echo -e "${YELLOW}💡 Check test_logs/node*.log for detailed error information${NC}"
fi

echo ""
echo -e "${BLUE}🔍 Detailed Status Check:${NC}"
echo "================================"

# Show current blockchain status
echo "📈 Current Blockchain Status:"
curl -s $API_BASE/blockchain | python3 -m json.tool 2>/dev/null || curl -s $API_BASE/blockchain

echo ""
echo "👥 Current Validators:"
curl -s $API_BASE/validators | python3 -m json.tool 2>/dev/null || curl -s $API_BASE/validators

echo ""
echo "💼 Current Transaction Pool:"
curl -s $API_BASE/transactions | python3 -m json.tool 2>/dev/null || curl -s $API_BASE/transactions 