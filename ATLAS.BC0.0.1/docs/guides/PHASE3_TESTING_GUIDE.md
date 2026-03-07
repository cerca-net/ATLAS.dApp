# Multi-Node Distributed Consensus Test Guide

## Phase 3: Distributed Block Production - Testing

This guide describes how to test the newly implemented deterministic consensus mechanism for distributed block production across multiple nodes.

## What Was Implemented

### 1. Deterministic Validator Selection
- **Algorithm**: SHA256 hash of `(lastBlockHash + blockHeight)` provides deterministic randomness
- **Weighted Selection**: Validators selected based on:
  - 40% stake weight
  - 30% performance score
  - 20% reputation score
  - 10% uptime
- **Consensus Guarantee**: All nodes will select the same validator for a given block height

### 2. Time Synchronization
- Block production aligned to 30-second intervals
- All nodes synchronize to the same time boundaries
- Prevents timing-based consensus failures

### 3. Automatic Block Broadcasting
- Blocks automatically broadcast to all peers via P2P
- Type-safe `BlockMessage` with full `block.Block` structure
- Validation on receipt before accepting into chain

## Testing Steps

### Unit Test (Already Passed ✅)
```bash
go test -v tests/consensus_test.go
```

Expected output:
- Validates that `ChooseValidator(hash, height)` returns same validator for same inputs
- Confirms deterministic selection works correctly

### Multi-Node Integration Test

#### Step 1: Start Node 1
```bash
# Terminal 1
go run cmd/main.go --datadir=data/node1 --port=8001 --api=8081 --validator=true --validator-key=data/node1/validator.hex
```

#### Step 2: Start Node 2
```bash
# Terminal 2  
go run cmd/main.go --datadir=data/node2 --port=8002 --api=8082 --validator=true --validator-key=data/node2/validator.hex
```

#### Step 3: Connect Nodes
Get Node 1's multiaddress from its startup logs (looks like `/ip4/127.0.0.1/tcp/8001/p2p/...`), then:

```bash
# Connect Node 2 to Node 1
curl -X POST http://localhost:8082/connect-peer \
  -H "Content-Type: application/json" \
  -d '{"multiaddress": "<NODE1_MULTIADDR>"}'
```

Or set environment variable before starting Node 2:
```bash
export NODE1_MULTIADDR="<NODE1_MULTIADDR>"
go run cmd/main.go --datadir=data/node2 --port=8002 --api=8082 --validator=true
```

#### Step 4: Submit a Transaction
```bash
# Use faucet to create a valid transaction
curl http://localhost:8081/faucet?address=0x1234567890abcdef
```

#### Step 5: Observe Consensus

**Watch Node 1 logs** for:
```
🔍 Attempting to choose validator...
🎯 We are the chosen validator! Attempting to forge block...
✅ Block forged successfully: Index=1
📡 Broadcasted new block 1 to peers.
```

**Watch Node 2 logs** for:
```
🔍 Attempting to choose validator...
⏳ Not our turn. Chosen: <validator_id>, Us: <our_id>
[P2P] Received block 1 from network
✅ Successfully added block 1 to local chain
```

### Expected Behavior

1. ✅ **Deterministic Selection**: Both nodes select the same validator for block N
2. ✅ **Only One Producer**: Only the selected validator produces the block
3. ✅ **Automatic Sync**: Non-producing node receives and validates the block
4. ✅ **Rotation**: As blocks progress, different validators are chosen based on deterministic algorithm
5. ✅ **Time Alignment**: Blocks produced at synchronized 30s intervals

### Verification Checklist

- [ ] Both nodes agree on selected validator (check logs)
- [ ] Only selected node produces block
- [ ] Non-selected node receives block via P2P
- [ ] Block height increments on both nodes
- [ ] Validator selection changes between blocks (if multiple validators with different weights)
- [ ] No fork conflicts or duplicate blocks

## Key Log Messages to Watch

### Successful Consensus
```
Node 1: ✅ Chosen validator: aa11223344... (Our address: aa11223344...)
Node 1: 🎯 We are the chosen validator!
Node 2: ⏳ Not our turn. Chosen: aa112233, Us: bb112233
```

### Block Propagation
```
Node 1: 📡 Broadcasted new block 1 to peers.
Node 2: [P2P] Received block 1 from network
Node 2: ✅ Successfully added block 1 to local chain
```

### Potential Issues

If nodes disagree on validator:
- Check that both nodes have synced to the same latest block
- Verify `lastBlockHash` is identical on both nodes
- Ensure both nodes use the same block height for selection

## Advanced: Three-Node Test

For full multi-validator rotation testing:

```bash
# Node 1
go run cmd/main.go --datadir=data/node1 --port=8001 --api=8081 --validator=true

# Node 2  
go run cmd/main.go --datadir=data/node2 --port=8002 --api=8082 --validator=true

# Node 3
go run cmd/main.go --datadir=data/node3 --port=8003 --api=8083 --validator=true
```

Connect all nodes in a mesh topology and observe that:
- Block production rotates between all three validators
- All nodes stay synchronized
- No forks or conflicts occur

## Performance Metrics

Track these metrics during testing:
- **Block Time**: Should be consistent ~30s
- **Propagation Latency**: Time from block production to receipt (<1s on localhost)
- **Consensus Failures**: Should be zero (all nodes agree)
- **Fork Events**: Should be zero in normal operation

## Next Steps (Phase 4)

Once Phase 3 is validated:
1. Implement fast-sync for new nodes joining the network
2. Add snapshot-based state synchronization
3. Implement automatic peer discovery beyond localhost
4. Add Byzantine fault tolerance mechanisms
