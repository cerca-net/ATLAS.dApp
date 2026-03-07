# PHASE 1 IMPLEMENTATION GUIDE: Transaction Broadcasting
**Network Dynamics - Transaction Propagation**

## Overview
Implement automatic transaction propagation across the network so that when a user submits a transaction to their node, all other nodes receive and process it.

---

## Current State vs Target State

### Current (Multi-Node Foundation)
```
User submits TX to Node A
└─ Transaction added to Node A's mempool
   └─ STOPS HERE ❌
```

### Target (Transaction Broadcasting)
```
User submits TX to Node A
├─ Transaction added to Node A's mempool
└─ Transaction broadcast to all peers
   ├─ Node B receives and validates TX
   │  └─ Adds to Node B's mempool ✅
   └─ Node C receives and validates TX
      └─ Adds to Node C's mempool ✅
```

---

## Implementation Plan

### Step 1: Add Transaction Message Type
**File:** `pkg/network/p2p.go`

Add new message type constant:
```go
const (
    MsgTypeBlock               = "block"
    MsgTypeTransaction         = "transaction"          // ← ADD THIS
    MsgTypeValidatorRegistration = "validator_registration"
    // ... existing types
)
```

### Step 2: Create Transaction Broadcast Function
**File:** `pkg/network/p2p.go`

```go
// BroadcastTransaction sends a transaction to all connected peers
func (node *P2PNode) BroadcastTransaction(tx transaction.Transaction) error {
    txBytes, err := json.Marshal(tx)
    if err != nil {
        return fmt.Errorf("failed to marshal transaction: %v", err)
    }

    msg := NetworkMessage{
        Type:    MsgTypeTransaction,
        Payload: txBytes,
    }

    peers := node.Host.Network().Peers()
    successCount := 0

    for _, peerID := range peers {
        if peerID == node.Host.ID() {
            continue // Skip self
        }

        err := node.SendMessage(peerID, msg)
        if err != nil {
            log.Printf("[P2P] Failed to send transaction to peer %s: %v", peerID, err)
            continue
        }
        successCount++
    }

    log.Printf("[P2P] Broadcast transaction %s to %d peers", tx.Hash, successCount)
    return nil
}
```

### Step 3: Add Transaction Handler
**File:** `cmd/main.go`

In the `main()` function, add transaction handler after validator registration handler:

```go
// Handle incoming transactions from peers
p2pNode.OnTransactionReceived = func(txMsg network.TransactionMessage) {
    log.Printf("[P2P] Received transaction from peer: %s", txMsg.Hash)
    
    // Reconstruct transaction from message
    tx := transaction.Transaction{
        Hash:      txMsg.Hash,
        From:      txMsg.From,
        To:        txMsg.To,
        Amount:    txMsg.Amount,
        Fee:       txMsg.Fee,
        Nonce:     txMsg.Nonce,
        Signature: txMsg.Signature,
        Timestamp: txMsg.Timestamp,
        Data:      txMsg.Data,
    }
    
    // Validate transaction
    if err := transactionManager.ValidateTransaction(&tx, stateManager); err != nil {
        log.Printf("[P2P] Rejected invalid transaction from peer: %v", err)
        return
    }
    
    // Check if we already have this transaction (prevent duplicates)
    if transactionManager.HasTransaction(tx.Hash) {
        log.Printf("[P2P] Ignoring duplicate transaction: %s", tx.Hash)
        return
    }
    
    // Add to local mempool
    if err := transactionManager.AddTransaction(&tx); err != nil {
        log.Printf("[P2P] Failed to add transaction to pool: %v", err)
        return
    }
    
    log.Printf("[P2P] Successfully added transaction %s to mempool", tx.Hash)
}
```

### Step 4: Update P2P Message Handler
**File:** `pkg/network/p2p.go`

In `handleStream()`, add case for transaction messages:

```go
func (node *P2PNode) handleStream(stream network.Stream) {
    // ... existing code ...
    
    var msg NetworkMessage
    if err := json.NewDecoder(stream).Decode(&msg); err != nil {
        // ... error handling ...
        return
    }
    
    switch msg.Type {
    case MsgTypeTransaction:
        if node.OnTransactionReceived != nil {
            var txMsg TransactionMessage
            if err := json.Unmarshal(msg.Payload, &txMsg); err != nil {
                log.Printf("[P2P] Failed to unmarshal transaction: %v", err)
                return
            }
            node.OnTransactionReceived(txMsg)
        }
    case MsgTypeValidatorRegistration:
        // ... existing code ...
    // ... other cases ...
    }
}
```

### Step 5: Add Transaction Broadcast to API
**File:** `internal/api/api.go`

Update `handleSubmitTransaction()`:

```go
func (api *APIServer) handleSubmitTransaction(w http.ResponseWriter, r *http.Request) {
    // ... existing validation and transaction creation code ...
    
    // Add transaction to local pool
    if err := api.transactionManager.AddTransaction(&tx); err != nil {
        http.Error(w, "Failed to add transaction: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    // ← ADD THIS: Broadcast to network
    if api.p2pNode != nil {
        if err := api.p2pNode.BroadcastTransaction(tx); err != nil {
            log.Printf("[API] Warning: Failed to broadcast transaction: %v", err)
            // Don't fail the request - transaction is still in local pool
        } else {
            log.Printf("[API] Transaction %s broadcast to network", tx.Hash)
        }
    }
    
    // ... existing response code ...
}
```

### Step 6: Add Duplicate Detection
**File:** `internal/blockchain/transaction_manager.go`

```go
// HasTransaction checks if a transaction exists in the pool by hash
func (tm *TransactionManager) HasTransaction(hash string) bool {
    tm.mu.RLock()
    defer tm.mu.RUnlock()
    
    for _, tx := range tm.pool {
        if tx.Hash == hash {
            return true
        }
    }
    return false
}
```

---

## Testing Procedure

### Setup
1. Start 2 nodes:
   ```bash
   scripts\start_network.bat
   ```

2. Connect the nodes:
   ```bash
   scripts\connect_nodes.bat "/ip4/127.0.0.1/tcp/8000/p2p/12D3Koo..."
   ```

3. Verify connection:
   ```bash
   curl http://localhost:8081/peers
   # Should show 1 peer
   ```

### Test 1: Transaction Propagation
```bash
# Submit transaction to Node 1
curl -X POST http://localhost:8080/submit-transaction \
  -H "Content-Type: application/json" \
  -d '{
    "from": "0xALICE_ADDRESS",
    "to": "0xBOB_ADDRESS",
    "amount": 100,
    "privateKey": "ALICE_PRIVATE_KEY"
  }'

# Check Node 1's mempool
curl http://localhost:8080/transactions

# Check Node 2's mempool (should contain the same transaction)
curl http://localhost:8081/transactions
```

### Expected Results
✅ Node 1 console shows:
```
[API] Transaction abc123... broadcast to network
[P2P] Broadcast transaction abc123... to 1 peers
```

✅ Node 2 console shows:
```
[P2P] Received transaction from peer: abc123...
[P2P] Successfully added transaction abc123... to mempool
```

✅ Both `/transactions` endpoints return the same transaction

### Test 2: Duplicate Prevention
```bash
# Submit same transaction again to Node 1
curl -X POST http://localhost:8080/submit-transaction (same data)
```

✅ Node 2 console shows:
```
[P2P] Ignoring duplicate transaction: abc123...
```

### Test 3: Invalid Transaction Rejection
```bash
# Submit transaction with invalid signature to Node 1
```

✅ Node 2 console shows:
```
[P2P] Rejected invalid transaction from peer: invalid signature
```

---

## Success Criteria Checklist

- [ ] Transaction submitted to Node 1 appears in Node 2's mempool
- [ ] Node 1 logs show successful broadcast
- [ ] Node 2 logs show successful reception
- [ ] Duplicate transactions are not re-added
- [ ] Invalid transactions from peers are rejected
- [ ] Transaction appears in both nodes' `/transactions` endpoint
- [ ] Block production includes transactions from network-wide mempool

---

## Rollback Plan

If issues occur:
1. The changes are isolated to transaction handling
2. Commenting out the `BroadcastTransaction()` call reverts to single-node behavior
3. No existing functionality is broken
4. Flutter app continues to work unchanged

---

## Next Phase Preview

After Phase 1 completion, Phase 2 will add:
- Block broadcasting when a block is produced
- Block validation before accepting from peers
- Chain synchronization for identical blockchain state across all nodes

---

**Estimated Implementation Time:** 2-3 hours
**Testing Time:** 1-2 hours
**Total:** Half a working day

Ready to implement! 🚀
