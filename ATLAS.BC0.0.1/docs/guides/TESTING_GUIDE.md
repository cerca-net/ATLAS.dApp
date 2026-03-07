# 🚀 Blockchain Testing Guide

Welcome to your blockchain testing environment! This guide will help you run and test your multi-node blockchain system.

## 🎯 What You'll Test

- **Multi-node blockchain network** (3 nodes)
- **Proof-of-Stake consensus** with validators
- **Transaction processing** and block creation
- **Real-time monitoring** through the UI
- **Wallet creation** and token transfers
- **Network synchronization** between nodes

## 🚀 Quick Start

### Option 1: Windows PowerShell (Recommended)
```powershell
# Run the test environment
.\run_test_nodes.ps1
```

### Option 2: Bash/Linux/Mac
```bash
# Run the test environment
./run_test_nodes.sh
```

### Option 3: Manual Start
```bash
# Build the blockchain
go build -o blockchain.exe .

# Start Node 1
./blockchain.exe -port 8001 -api 8081 -peers 5

# Start Node 2 (in new terminal)
./blockchain.exe -port 8002 -api 8082 -peers 5

# Start Node 3 (in new terminal)
./blockchain.exe -port 8003 -api 8083 -peers 5
```

## 📱 Testing Interface

Once your nodes are running, open these URLs in your browser:

### 🏠 Main Hub
```
file://[your-project-path]/frontend/index.html
```

### 📊 Test Dashboard (Real-time monitoring)
```
file://[your-project-path]/frontend/test-dashboard.html
```

### 🖥️ Node Dashboard
```
file://[your-project-path]/frontend/node-dashboard.html
```

### 🔍 Blockchain Explorer
```
file://[your-project-path]/frontend/explorer.html
```

### 💰 Wallet Interface
```
file://[your-project-path]/frontend/wallet.html
```

## 🌐 API Endpoints

Your blockchain nodes expose these APIs:

### Node 1: http://localhost:8081
### Node 2: http://localhost:8082
### Node 3: http://localhost:8083

**Available endpoints:**
- `GET /status` - Node status and metrics
- `GET /blocks` - List recent blocks
- `GET /validators` - List active validators
- `GET /mempool` - Pending transactions
- `POST /submit-transaction` - Submit new transaction
- `POST /faucet` - Create test wallet with tokens
- `GET /balance?address=...` - Check wallet balance

## 🧪 What to Test

### 1. **Node Status Check**
- Open the test dashboard
- Verify all 3 nodes show as "🟢 Online"
- Check block heights and transaction pool sizes

### 2. **Block Production**
- Watch blocks being created every few seconds
- Verify block numbers increase over time
- Check that blocks contain transactions

### 3. **Transaction Testing**
- Use the "📤 Send Test TX" button to create random transactions
- Watch transactions appear in the mempool
- Verify transactions get included in blocks

### 4. **Wallet Creation**
- Use the "💰 Create Wallet" button to create test wallets
- Check wallet balances using the faucet
- Send transactions between wallets

### 5. **Network Synchronization**
- Compare block heights across all nodes
- Verify they stay synchronized (within 1-2 blocks)
- Check that transactions propagate between nodes

### 6. **Validator Testing**
- Monitor validator activity
- Check staking amounts and active status
- Verify consensus is working properly

## 🔧 Manual API Testing

You can also test the APIs directly using curl or Postman:

### Create a test wallet:
```bash
curl -X POST http://localhost:8081/faucet \
  -H "Content-Type: application/json" \
  -d '{"address": "my_test_wallet"}'
```

### Send a transaction:
```bash
curl -X POST http://localhost:8081/submit-transaction \
  -H "Content-Type: application/json" \
  -d '{
    "from": "my_test_wallet",
    "to": "another_wallet",
    "amount": 50,
    "timestamp": 1234567890,
    "fee": 1
  }'
```

### Check node status:
```bash
curl http://localhost:8081/status
```

### Get recent blocks:
```bash
curl http://localhost:8081/blocks?limit=5
```

## 📊 Monitoring

### Real-time Dashboard
The test dashboard automatically refreshes every 5 seconds and shows:
- Node status (online/offline)
- Block heights and transaction counts
- Recent blocks with transaction details
- Active validators and their stakes
- Pending transactions in mempool

### Logs
Check the logs for detailed information:
```bash
# View Node 1 logs
tail -f test_logs/node1.log

# View Node 2 logs  
tail -f test_logs/node2.log

# View Node 3 logs
tail -f test_logs/node3.log
```

## 🛑 Stopping the Test Environment

### Option 1: PowerShell
Press `Ctrl+C` in the PowerShell window running the test script.

### Option 2: Manual
```bash
# Stop all blockchain processes
pkill -f blockchain.exe
```

### Option 3: Windows Task Manager
- Open Task Manager
- Find "blockchain.exe" processes
- End them manually

## 🔍 Troubleshooting

### Nodes not starting?
1. Check if ports 8001-8003 and 8081-8083 are available
2. Ensure Go is installed and in your PATH
3. Check the logs in `test_logs/` directory

### UI not loading?
1. Make sure you're using the correct file paths
2. Try opening the HTML files directly in your browser
3. Check browser console for any JavaScript errors

### Transactions not working?
1. Verify nodes are running and online
2. Check that wallets have sufficient balance
3. Look at the mempool to see if transactions are pending

### Network not syncing?
1. Ensure all nodes can communicate on their peer ports
2. Check that consensus is working properly
3. Verify validators are active and staking

## 🎉 Success Indicators

Your blockchain is working correctly when you see:

✅ **All 3 nodes show as "🟢 Online"**  
✅ **Block heights increase over time**  
✅ **Transactions get processed and included in blocks**  
✅ **Network stays synchronized**  
✅ **Validators are active and participating in consensus**  
✅ **Wallets can send and receive tokens**  

## 🚀 Next Steps

Once you've verified everything is working:

1. **Explore the codebase** - Understand how the blockchain works
2. **Modify parameters** - Adjust block time, stake amounts, etc.
3. **Add new features** - Implement additional functionality
4. **Scale up** - Add more nodes to test larger networks
5. **Deploy** - Set up for production use

Happy testing! 🎯 