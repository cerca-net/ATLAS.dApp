#!/bin/bash

echo "🚀 Starting Blockchain Test Environment..."

# Kill any existing processes
echo "🧹 Cleaning up existing processes..."
pkill -f "blockchain.exe" 2>/dev/null || true

# Build the blockchain
echo "🔨 Building blockchain..."
go build -o blockchain.exe .

# Create test directories
mkdir -p test_logs

echo "🌐 Starting blockchain nodes..."

# Start Node 1
echo "🟢 Starting Node 1 (Peer: 8001, API: 8081)"
./blockchain.exe -port 8001 -api 8081 -peers 5 > test_logs/node1.log 2>&1 &
NODE1_PID=$!

# Start Node 2
echo "🟢 Starting Node 2 (Peer: 8002, API: 8082)"
./blockchain.exe -port 8002 -api 8082 -peers 5 > test_logs/node2.log 2>&1 &
NODE2_PID=$!

# Start Node 3
echo "🟢 Starting Node 3 (Peer: 8003, API: 8083)"
./blockchain.exe -port 8003 -api 8083 -peers 5 > test_logs/node3.log 2>&1 &
NODE3_PID=$!

echo "⏳ Waiting for nodes to initialize..."
sleep 10

echo "💰 Creating test wallets..."

# Create test wallets using faucet
curl -s -X POST http://localhost:8081/faucet \
    -H "Content-Type: application/json" \
    -d '{"address": "test_wallet_1"}' > /dev/null

curl -s -X POST http://localhost:8082/faucet \
    -H "Content-Type: application/json" \
    -d '{"address": "test_wallet_2"}' > /dev/null

curl -s -X POST http://localhost:8083/faucet \
    -H "Content-Type: application/json" \
    -d '{"address": "test_wallet_3"}' > /dev/null

echo "⏳ Waiting for initial block production..."
sleep 15

echo "🔄 Running test transactions..."

# Send test transactions
curl -s -X POST http://localhost:8081/submit-transaction \
    -H "Content-Type: application/json" \
    -d '{
        "from": "test_wallet_1",
        "to": "test_wallet_2",
        "amount": 50,
        "timestamp": 1234567890,
        "fee": 1
    }' > /dev/null

sleep 5

curl -s -X POST http://localhost:8082/submit-transaction \
    -H "Content-Type: application/json" \
    -d '{
        "from": "test_wallet_2",
        "to": "test_wallet_3",
        "amount": 30,
        "timestamp": 1234567895,
        "fee": 1
    }' > /dev/null

sleep 5

curl -s -X POST http://localhost:8083/submit-transaction \
    -H "Content-Type: application/json" \
    -d '{
        "from": "test_wallet_3",
        "to": "test_wallet_1",
        "amount": 20,
        "timestamp": 1234567900,
        "fee": 1
    }' > /dev/null

echo "⏳ Waiting for transactions to be processed..."
sleep 10

echo ""
echo "🎉 Test environment is ready!"
echo ""
echo "🌐 Blockchain nodes are running on:"
echo "   Node 1: API http://localhost:8081"
echo "   Node 2: API http://localhost:8082"
echo "   Node 3: API http://localhost:8083"
echo ""
echo "📱 You can now open the UI files in your browser:"
echo "   Main Hub: file://$(pwd)/frontend/index.html"
echo "   Node Dashboard: file://$(pwd)/frontend/node-dashboard.html"
echo "   Blockchain Explorer: file://$(pwd)/frontend/explorer.html"
echo "   Wallet Interface: file://$(pwd)/frontend/wallet.html"
echo ""
echo "📋 To view logs:"
echo "   tail -f test_logs/node1.log"
echo "   tail -f test_logs/node2.log"
echo "   tail -f test_logs/node3.log"
echo ""
echo "🔧 To stop all nodes, press Ctrl+C"
echo ""

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Stopping blockchain nodes..."
    kill $NODE1_PID 2>/dev/null || true
    kill $NODE2_PID 2>/dev/null || true
    kill $NODE3_PID 2>/dev/null || true
    echo "✅ Test environment stopped!"
    exit 0
}

# Set up signal handler
trap cleanup INT TERM

# Keep the script running
echo "🔄 Test environment is running. Press Ctrl+C to stop all nodes."
echo ""

while true; do
    sleep 1
done 