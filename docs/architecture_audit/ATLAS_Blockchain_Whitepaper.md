# ATLAS CercaChain: Protocol Architecture & Mechanics

This document serves as a technical brief on the underlying consensus, networking, and tokenomic layers of the ATLAS CercaChain network. Unlike traditional EVM forks, ATLAS is a bespoke Go-based Layer-1 protocol designed specifically for social-commerce integration.

---

## 1. Consensus Mechanism: Multi-Variate Proof-of-Stake (mv-PoS)
ATLAS utilizes a deterministic, weighted Proof-of-Stake algorithm (with PBFT characteristics) to achieve consensus without the energy expenditure of Proof-of-Work.

### 1.1 Deterministic Validator Selection
Traditional PoS systems often use naive randomization or rely purely on stake volume to select the next block producer. ATLAS implements a proprietary **Multi-Variate Selection Formula** that evaluates a validator across four distinct dimensions:

```go
weight = (StakeWeight * 0.4) + (PerformanceScore * 0.3) + (ReputationScore * 0.2) + (Uptime * 0.1)
```

1.  **Stake (40%):** The raw amount of `TCOIN` locked in the bond.
2.  **Performance (30%):** A moving average tracking the validator's success in validating and forging blocks.
3.  **Reputation (20%):** Tied to the validator's Decentralized Identifier (DID) and slashed upon bad behavior.
4.  **Uptime (10%):** Calculated based on P2P heartbeat tracking.

To select a validator, the `ConsensusManager` hashes the *previous block hash* combined with the *current block height* using `SHA-256`. This generates a deterministic pseudorandom seed, ensuring all nodes on the network independently calculate and agree on the exact same next validator without requiring an extra communication round.

### 1.2 Slashing & Finality
*   **Slashing:** Malicious or failing nodes are penalized. A single failure reduces the `ReputationScore` by 50%. Accumulating 3 strikes results in immediate ejection from the active validator set.
*   **Finality:** Implements tracking for block confirmations. Once a block hits the `finalityThreshold`, it is marked as irreversibly finalized.

---

## 2. Tokenomics & Economics (TCOIN)
The native utility token, `TCOIN`, powers the network's gas, security, and social mechanics.

### 2.1 The Native Economy
1.  **Genesis Supply:** A fixed allocation instantiated at block 0, held by the Network Treasury.
2.  **Block Rewards:** Validators receive a hardcoded reward of **10 TCOIN** per forged block. This is injected natively by the `ForgeBlock` function as a `TxTypeRegular` transaction from the `network` sender address.
3.  **Transaction Fees:** Users pay gas fees to submit transactions, which are credited directly to the validator forging the block.
4.  **Social Tipping:** Transactions are used to "Energize" social posts.
5.  **Staking Requirements:** The network enforces a minimum stake (`MIN_STAKE` = 1 TCOIN) to participate in consensus.

### 2.2 System Smart Contracts & The CercaVM
ATLAS natively implements a stack-based Virtual Machine (`CercaVM`). While users can deploy generic JSON-based contracts (`TxTypeDeploy`), the network's core features are hardcoded as **System Contracts** to ensure speed and security.

For example, the **Marketplace Contract** is natively intercepted by the `StateManager` during a `TxTypeCall`. If a transaction targets the Marketplace address, the Go engine natively routes the payload to functions like `createOrder`, `releaseFunds`, or `raiseDispute`. This provides the security of a smart contract with the execution speed of native compiled Go code.

---

## 3. Node Architecture & Networking Layer

### 3.1 Dual-State Persistence
ATLAS nodes utilize a unique, highly resilient storage architecture managed by the `StateManager`.
*   **SQLite Primary DB:** The primary ledger state (Account balances, nonces, and smart contract memory) is tracked utilizing a lightweight SQLite database using Write-Ahead Logging (WAL).
*   **JSON Fast-Sync Snapshots:** Every 5 minutes (or 100 blocks), the network takes a complete checksum-verified JSON snapshot of all balances and contracts. When a new node joins the network, it does not need to replay every block from genesis; it requests the latest snapshot via the `/snapshot/load` API, allowing synchronization in under 1 second.

### 3.2 Peer-to-Peer Networking (`libp2p`)
ATLAS relies on Protocol Labs' `libp2p` networking stack (the same engine powering IPFS and Filecoin).
*   **Transports:** Defaults to QUIC for fast multiplexing, with TCP fallback.
*   **Discovery:** Nodes use mDNS to discover local peers instantly without centralized DNS.
*   **Gossip Protocol:** Uses `go-libp2p-pubsub` to propagate transactions (`TxPool`) and new blocks to all connected peers in milliseconds.

### 3.3 Node Operator Types
*   **Validators:** Bind to a local ECDSA private key. They run the `ConsensusManager`, participate in the PBFT lottery, and write SQLite data.
*   **Relay/Observer Nodes:** Light clients that verify block signatures and gossip them to the network, without participating in consensus block generation. These are crucial for handling read-heavy API requests (like those coming from the Flutter App).
