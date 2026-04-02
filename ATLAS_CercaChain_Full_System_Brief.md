# ATLAS CercaChain — Full-Scope System Architecture Brief
**Version BC0.0.1 · March 2026 · Compiled from source**

---

## 1. Project Identity & Vision

**ATLAS CercaChain** is a purpose-built Layer-1 Proof-of-Stake blockchain designed to power a **social-commerce-governance** platform. It is not a fork of Ethereum, Cosmos, or any existing chain — it is an original Go implementation with a native stack-based VM, custom secp256k1/ECDSA cryptography, a hybrid SQLite + JSON-snapshot persistence layer, and a libp2p-powered peer-to-peer networking stack.

The ecosystem consists of three distinct, interconnected applications:

| Application | Role | Language / Stack |
|---|---|---|
| `ATLAS.BC0.0.1` | Layer-1 blockchain node daemon | Go 1.24 |
| `cercaend` | End-user mobile/web app | Flutter / Dart |
| `cerca-admin-panel` | Network operator control plane | React + Vite + TypeScript |

All three communicate through a single REST API served by the Go node on `:8080`.

---

## 2. High-Level Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                   NETWORK LAYER (libp2p)                     │
│  P2P discovery · mDNS · QUIC/TCP · pubsub gossip            │
└───────────────────────┬──────────────────────────────────────┘
                        │ blocks / txs / validators
┌───────────────────────▼──────────────────────────────────────┐
│               ATLAS BLOCKCHAIN ENGINE (Go)                   │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │ BlockManager │  │ TxManager    │  │ ConsensusManager │  │
│  │ (chain head) │  │ (mempool)    │  │ (PBFT / PoS)     │  │
│  └──────────────┘  └──────────────┘  └──────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              StateManager (accounts + contracts)      │   │
│  │  ┌──────┐ ┌─────────┐ ┌──────────────┐ ┌─────────┐  │   │
│  │  │TCOIN │ │Staking  │ │ Marketplace  │ │Governance│  │   │
│  │  │Token │ │Contract │ │  Contract    │ │Contract  │  │   │
│  │  └──────┘ └─────────┘ └──────────────┘ └─────────┘  │   │
│  └──────────────────────────────────────────────────────┘   │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐   │
│  │   CercaVM│ │ShardMgr  │ │IdentityM │ │SocialManager │   │
│  │  (stack) │ │(4 shards)│ │(DID/KYC) │ │(posts/energy)│   │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────┘   │
│  ┌───────────────────────────────────────────────────────┐  │
│  │   REST API Server  (100+ endpoints, CORS, :8080)      │  │
│  └───────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  Persistence: SQLite (blockchain.db) + JSON snapshots │  │
│  └───────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────┘
         │               │                    │
         ▼               ▼                    ▼
   CercaEnd          Admin Panel          Firebase /
   (Flutter)       (React+Vite)           Supabase
                                      (off-chain social DB)
```

---

## 3. Layer-1 Blockchain Engine — `ATLAS.BC0.0.1`

### 3.1 Language & Dependencies
- **Language**: Go 1.24 (`module atlas-blockchain`)
- **Key runtime deps**:
  - `github.com/libp2p/go-libp2p v0.47.0` — P2P host
  - `github.com/libp2p/go-libp2p-pubsub` — gossip protocol
  - `github.com/mattn/go-sqlite3` — CGO SQLite driver
  - `github.com/consensys/gnark + gnark-crypto` — ZKP framework (deferred to V2)
  - `github.com/tyler-smith/go-bip39` — BIP-39 mnemonic support
  - `github.com/decred/dcrd/dcrec/secp256k1/v4` — secp256k1 curve
  - `github.com/shirou/gopsutil` — system monitoring
  - `github.com/multiformats/go-multiaddr` — libp2p multiaddresses
  - `github.com/quic-go/quic-go` — QUIC transport
  - `github.com/pion/webrtc/v4` — WebRTC transport (via libp2p)

### 3.2 Package Layout

```
ATLAS.BC0.0.1/
├── cmd/main.go                # Entry point — boots all subsystems
├── internal/
│   ├── api/                   # HTTP REST API (api.go, ~4,300 lines)
│   ├── blockchain/            # Core chain logic
│   │   ├── block_manager.go
│   │   ├── consensus.go       # PBFT PoS validator selection & slashing
│   │   ├── state_manager.go   # Account balances, contracts, oracle data
│   │   ├── transaction_manager.go
│   │   ├── chain_sync.go
│   │   ├── fast_sync.go
│   │   └── state_adapter.go   # Bridge: StateManager <-> CercaVM
│   ├── defi/                  # DeFi: lending, DEX, staking logic
│   ├── governance/            # Proposals, votes, referendums, committees
│   ├── identity/              # DID/KYC identity management
│   └── social/                # Posts, comments, tips, fossilization
├── pkg/
│   ├── block/
│   ├── config/                # Genesis config, blockchain params
│   ├── crypto/                # ECDSA/secp256k1 key utilities + ZK stub
│   ├── database/              # SQLite ORM, backup/recovery manager
│   ├── monitoring/
│   ├── network/               # P2P node, peer manager, message types
│   ├── sharding/              # 4-shard manager + cross-shard TXs
│   ├── transaction/
│   ├── vm/                    # CercaVM + 4 system smart contracts
│   └── wallet/                # HD wallet, BIP-39, address derivation
```

### 3.3 Node Startup Sequence

1. Parse CLI flags (`--port`, `--api`, `--validator`, `--key`, `--datadir`, `--genesis`)
2. Load `genesis.json` → block time, max block size, min stake, initial allocations
3. Initialize `StateManager` (opens SQLite at `<datadir>/blockchain.db`)
4. Apply genesis allocations (idempotent)
5. Run DB migration (JSON → SQLite)
6. Initialize `TransactionManager`, `BlockManager`, `ConsensusManager`
7. Register genesis validators
8. Initialize `IdentityManager`, `DeFiManager`, `SocialManager`, `GovernanceManager`
9. Initialize `ShardManager` (4 shards, 10 validators/shard, PBFT)
10. Start libp2p P2P node (load/generate key at `<datadir>/nodekey.priv`)
11. Initialize `FastSyncManager`, `PeerManager`, `ChainSyncManager`
12. Register P2P callbacks (block received, tx received, validator registration)
13. Start `APIServer` on `:8080`
14. Start block production loop (`produceBlocks()`)
15. Start backup system
16. Await `SIGINT/SIGTERM` → graceful shutdown

### 3.4 Consensus Engine

**Algorithm**: Deterministic weighted Proof-of-Stake (PBFT-flavored, single-node finality for V1)

**Validator selection formula**:
```
weight = (stake_weight × 0.4) + (performance × 0.3) + (reputation × 0.2) + (uptime × 0.1)
```
where `stake_weight = validator.Stake / totalStake`.

A SHA-256 seed derived from `lastBlockHash + blockHeight` drives a deterministic lottery over the sorted validator set.

**Block production conditions**:
- Node state must be `running`
- Mempool must be non-empty (activator mode — no empty blocks)
- DevNet override: `isLocalValidator = true` (any node produces blocks)

**Constants**:
- Block reward: 10 TCOIN/block
- Minimum stake: 1 TCOIN (genesis param)
- Slashing threshold: 3 events → validator ejected
- Finality threshold: 1 confirmation

**Validator struct**:
- `Address`, `Stake`, `Delegations`, `LastBlock`, `SlashCount`, `Active`
- `PerformanceScore`, `Uptime`, `LastActive`, `BlocksProduced`, `BlocksValidated`
- `ReputationScore`, `SlashingHistory`, `RewardHistory`, `TotalRewards`
- `KYC` (FullName, Country, IDNumber, Verified)

### 3.5 State Manager

The canonical source of truth for all account state:

- **Primary storage**: SQLite (`blockchain.db`)
- **Fallback**: SHA-256 checksum-verified JSON snapshots
- **Snapshot interval**: every 5 minutes + every block if >1 hour since last
- **Backup system**: `BackupManager` + `RecoveryManager` (WAL mode)

Account model: `{Address, Balance int64, Nonce uint64}`

Transaction types in `updateState()`:
- `TxTypeRegular` — transfer + fee deduction + nonce increment
- `TxTypeDeploy` — deploy JSON-encoded smart contract into registry
- `TxTypeCall` — call contract function (system contracts intercepted directly)
- `TxTypeProposal` — submit on-chain governance proposal
- `TxTypeVote` — cast governance vote
- `TxTypeStake` — register/increase validator stake

---

## 4. CercaVM — Native Stack-Based Virtual Machine

`pkg/vm/vm.go` (767 lines) implements a minimal register-less stack machine.

### 4.1 Architecture

**Two stacks**: `stack []int64` + `stringStack []string`  
**Memory**: `Memory map[string]int64` + `StringMem map[string]string` (volatile)  
**Persistent storage**: via `StateAdapter` (SSTORE/SLOAD opcodes backed by SQLite)  
**Gas model**: per-opcode fixed costs  
**Max call depth**: 10

### 4.2 Opcode Set

| Category | Opcodes |
|---|---|
| Stack ops | PUSH, POP, DUP, SWAP, PUSHS, POPS |
| Arithmetic | ADD, SUB, MUL, DIV |
| Logic | GT, LT, EQ, NEQ, AND, OR, NOT |
| Control flow | JUMP, JUMPIF, CALL, RETURN, REQUIRE |
| Memory | STORE, LOAD, SSTORE, SLOAD, SSTORE_S, SLOAD_S |
| Token ops | TRANSFER, BALANCE, MINT, BURN |
| Context | CALLER, TIMESTAMP, BLOCKNUM |
| Events | EMIT |

### 4.3 Gas Costs (selected)

| Opcode | Cost |
|---|---|
| PUSH | 3 |
| ADD/SUB | 3 |
| MUL/DIV | 5 |
| TRANSFER | 21 |
| MINT | 30 |
| BURN | 20 |
| SSTORE | 20 |
| SLOAD | 5 |
| CALL | 10 |

### 4.4 Contract Permission Model

- **System, governance, voting contracts**: auto-approved, unrestricted function access
- **Custom contracts**: governance-approved; `AllowedFunctions` whitelist enforced

### 4.5 ZK Proofs (V1 Stub)

gnark Groth16 scaffolding present. `VerifyZKProof()` returns `true` unconditionally in V1. Full ZK circuit integration deferred to V2.

---

## 5. System Smart Contracts (4 Core Contracts)

Initialized at node startup in `NewAPIServer()`.

### 5.1 TCOIN Token Contract (`CONTRACT_TCOIN_SYSTEM`)
- **Genesis supply**: 1,000,000,000 TCOIN (mint to Treasury)
- **Treasury address**: derived from fixed BIP-39 mnemonic
- **Functions**: transfer, mint, burn, balanceOf, approve, transferFrom

### 5.2 Staking Contract (`CONTRACT_STAKING_SYSTEM`)
- Min stake, max validators, lock period, slashing penalty (all configurable via `vm.*` constants)
- **Functions**: stake, unstake, getStakeInfo, getActiveValidatorCount, getTotalStaked

### 5.3 Marketplace Contract (`vm.MarketplaceContractAddress`)
Escrow-based P2P commerce:
- **Order lifecycle**: `createOrder` → `releaseFunds` OR `raiseDispute` → `resolveDispute`/`refundBuyer`
- Fee rate in basis points; all tracked in contract storage

### 5.4 Governance Contract (`vm.GovernanceContractAddress`)
- Min proposal stake: 1,000 | Min voting stake: 100
- Voting period: 1,000 blocks | Quorum: 10% | Pass: 60%
- Execution delay: 100 blocks after vote end

---

## 6. REST API — Complete Endpoint Catalogue (100+)

**Base URL**: `http://localhost:8080` | **CORS**: Fully permissive (testnet)

### Core Chain
`GET /block` · `GET /blocks` · `GET /transaction` · `GET /mempool` · `POST /submit-transaction` · `GET /balance` · `GET /status` · `GET /nonce` · `GET /fee-info`

### Validators
`GET /validators` · `GET /validator` · `POST /register-validator` · `POST /update-stake` · `POST /update-user-stake` · `GET /node-address`

### Wallet
`POST /create-wallet` · `POST /import-wallet`

### Treasury / Faucet
`POST /faucet` (1,000 TCOIN) · `POST /admin/faucet` · `GET /treasury` · `GET /admin/treasury-history`

### System Contracts
`GET /token` · `GET /staking` · `GET /marketplace` · `GET /governance-contract`

### FlutterFlow Integration
`POST /flutterflow/connect-wallet` · `POST /flutterflow/authenticate` · `GET /flutterflow/wallet-info` · `POST /flutterflow/send-transaction` · `GET /flutterflow/transaction-history` · `POST /flutterflow/disconnect`

### Identity
`POST /identity/create` · `GET /identity/get` · `POST /identity/update-profile` · `POST /identity/update-activity` · `POST /identity/create-proof` · `POST /identity/verify-proof` · `GET /identity/social` · `GET /identity/commerce` · `GET /identity/governance`

### Social Platform
`POST /social/post/create` · `GET /social/post/get` · `POST /social/comment/create` · `POST /social/like` · `POST /social/unlike` · `POST /social/tip` · `GET /social/feed` · `GET /social/search` · `GET /social/trending` · `POST /social/report` · `GET /social/object/energy` · `POST /social/object/energize`

### Governance
`GET /governance/proposals` · `GET /governance/proposal` · `POST /governance/proposal/create` · `POST /governance/proposal/activate` · `POST /governance/proposal/vote` · `POST /governance/proposal/execute` · `POST /governance/proposal/discuss` · `GET /governance/proposals/active` · `GET /governance/proposals/category` · `POST /governance/referendum/create` · `POST /governance/referendum/vote` · `POST /admin/resolve-dispute` · `GET /admin/disputes`

### Custom Smart Contracts
`POST /contract/deploy` · `POST /contract/call` · `GET /contract/list` · `GET /contract/info` · `GET /contract/examples`

### Oracle
`POST /oracle/submit` · `GET /oracle/latest`

### Privacy & GDPR
`POST /privacy/encrypt` · `POST /privacy/decrypt` · `POST /privacy/gdpr-delete` · `POST /privacy/gdpr-anonymize`

### Sharding
`GET /sharding/status` · `GET /sharding/shard` · `POST /sharding/assign-validator` · `POST /sharding/cross-shard-tx` · `GET /sharding/statistics`

### Monitoring
`GET /monitoring/status` · `GET /monitoring/metrics` · `GET /monitoring/health` · `GET /monitoring/alerts` · `GET /monitoring/performance` · `GET /monitoring/history` · `GET /monitoring/trends`

### Snapshots & Sync
`POST /snapshot/create` · `GET /snapshot/latest` · `POST /snapshot/load` · `GET /sync/status` · `POST /sync/start`

### P2P Peers
`GET /peers` · `POST /connect-peer` · `GET /peers/status` · `POST /peers/reconnect` · `GET /peers/validators`

### Database Backup
`GET /backup/status` · `GET /backup/list` · `POST /backup/create`

### Node Control
`POST /node/start` · `POST /node/stop` · `POST /node/pause` · `POST /node/sync` · `GET /node/status` · `GET /node/logs`

### Testing
`POST /run-tests` · `POST /test-performance` · `POST /test-security` · `POST /test-integration` · `POST /start-test-env` · `POST /stop-test-env` · `GET /test-env-status`

### Other
`GET /network/architecture` · `GET /` (serves `./web/frontend` static files)

---

## 7. P2P Networking Layer

**Library**: `go-libp2p v0.47`  
**Transport**: QUIC (primary) + TCP (optional `--legacy-net` flag)  
**Discovery**: mDNS (LAN) + manual `NODE1_MULTIADDR` env var  
**Peer identity**: ED25519 key at `<datadir>/nodekey.priv`  
**Gossip**: `go-libp2p-pubsub`

**Message types**: `BlockMessage`, `TransactionMessage`, `ValidatorRegistrationMessage`, `NetworkMessage`

**P2P callbacks**:
- `OnBlockReceived` → validate + add to local chain
- `OnTransactionReceived` → deduplicate → add to mempool
- `OnValidatorRegistrationReceived` → register external validator → rebroadcast

**Validator heartbeat**: Rebroadcast registration every 30 seconds.  
**`PeerManager`**: Persistent peer store with validator tracking.

---

## 8. Sharding Architecture

- **4 shards**, 10 validators/shard, PBFT consensus per shard
- Validator-to-shard: deterministic `hash(address) mod 4`
- Cross-shard TXs: routed via `ShardManager` with 5s default delay
- API endpoints for shard inspection and manual assignment

---

## 9. Social Platform — Energy Physics Engine

### 9.1 Object Model

| Property | Physics Analogy | Function |
|---|---|---|
| `TipBalance` | Energy | Keeps object alive |
| `InfluenceScore` | Velocity | Determines feed rank |
| `Upvotes` | Positive gravity | Boosts influence |
| `Downvotes` | Negative gravity | Reduces influence |

**Influence formula**:
```
InfluenceScore = ((upvotes×10) - (downvotes×20) + (tipBalance×5)) - (ageHours × 100)
```

**Fossilization**: `TipBalance < 50` → object marked `"fossilized"`, removed from active feed.  
**Revival cost**: 50 TCOIN minimum.  
**Grace period**: All new objects start with 100 TCOIN energy.  
**Comment cost**: 2 TCOIN (feeds the post's energy).

### 9.2 Causal Time (Lamport Clock)

`SocialManager.logicalClock` implements Lamport timestamps:
- Each post increments the global clock
- Each comment: `logicalTime = max(post.LogicalTime, localClock) + 1`

### 9.3 Firebase Object Bridge

`/social/object/energy` + `/social/object/energize`: Firebase document IDs treated as blockchain objects. TCOIN transferred on-chain when a Firebase object is "energized."

### 9.4 Content Moderation
- Keyword/regex/AI filter types
- `AIModerator` stub (confidence threshold)
- Auto-flagged content → `"hidden"` status
- Report system with priority levels and review workflow

---

## 10. Identity System

**`UserIdentity`** fields:
- `Address`, `PublicKey`, `DID`
- `Profile`: DisplayName, Bio, Avatar, Location, Website, VerifiedBadge
- `KYC`: status, level, country, documents, expiry
- `Privacy`: DataMinimization, ZKProofEnabled, SelectiveDisclosure, GDPRConsent
- `Activity`: PostsCreated, CommentsMade, VotesCast, ProposalsCreated, TotalTokensEarned/Spent, TipsGiven/Received, OrdersCreated/Completed
- `Reputation`: Overall, TrustScore, CommunityScore, CommerceScore, GovernanceScore
- `Credentials`: verifiable attestation list

**Governance voting power**:
```
power = (TokensEarned × 0.1) + (Reputation.Overall × 10)
      + (ProposalsCreated × 100) + (VotesCast × 10)
      + (PostsCreated × 5) + (CommentsMade × 2)
```

---

## 11. Governance System

**Proposal lifecycle**: `draft` → `active` → `passed`/`failed` → `executed`/`cancelled`

**Proposal categories**: platform, defi, social, technical, economic

**Action types**: `parameter_change`, `defi_parameter`, `social_parameter`, `treasury_transfer`, `contract_upgrade`, `committee_creation`

**Default parameters**: Min proposal stake 1,000 | Min voting stake 100 | Voting period 1,000 blocks | Quorum 10% | Pass threshold 60% | Execution delay 100 blocks

**Committees**: Named groups (technical/economic/social/security) with Chair + Members.  
**Referendums**: Multi-option community votes with weighted voting power.  
**Social integration**: Every proposal auto-creates a social post with `#governance #proposal`.  

---

## 12. DeFi Module

`internal/defi/` (5 files, ~41KB):
- `defi.go` — core coordinator
- `defi_components.go` — lending pool, liquidity
- `defi_dex.go` — decentralized exchange swaps
- `defi_staking.go` — advanced staking mechanics
- `tokenomics.go` — token economic model parameters

Feeds data to `GovernanceManager` for parameter change execution and `IdentityManager` for commerce reputation scoring.

---

## 13. CercaEnd — Flutter Web Application

### 13.1 Stack

Flutter web (primary), Android, iOS, Windows.  
Dart SDK `>=3.0.0 <4.0.0`. State: `Provider` pattern.

### 13.2 Dual-Database Architecture

| Data | Storage | Reason |
|---|---|---|
| Auth, profiles, posts, media, orders | Firebase + Supabase | Real-time, scalable, offline-first |
| TCOIN balances, txs, validators, governance | ATLAS blockchain (HTTP) | Immutable, trustless |

### 13.3 Source Layout

```
cercaend/lib/
├── main.dart, app_state.dart
├── auth/                    # Firebase auth flows
├── backend/                 # Supabase + Firebase bridge (~49KB)
├── services/blockchain/     # HTTP clients for ATLAS API
│   ├── blockchain_service.dart  # 1,143 lines
│   ├── wallet_service.dart      # Key management
│   ├── energy_service.dart
│   └── social_service.dart
├── mainpages/
│   ├── block_explorer/
│   ├── feedpage/
│   ├── node_dashboard/
│   ├── orderpage/
│   ├── publicpage/
│   └── userpage/
└── secondarypages/
    ├── hashingpage/
    ├── method_order/
    ├── method_wallet/
    ├── new_catalogue/
    ├── pinned_objects/
    ├── pinned_users/
    ├── settingspages/
    └── userrating/
```

### 13.4 Transaction Signing Flow

1. Fetch nonce from `/nonce?address=...`
2. Build `SendTransactionRequest` struct
3. SHA-256 hash canonical JSON `{Sender, Recipient, Amount, Fee, Timestamp, Nonce, Data}`
4. ECDSA sign with secp256k1 private key (via `elliptic` + `ecdsa` Dart packages)
5. POST to `/submit-transaction`

### 13.5 Key Dependencies

| Package | Purpose |
|---|---|
| `firebase_auth`, `cloud_firestore`, `firebase_storage` | Firebase |
| `supabase_flutter` | Supabase |
| `go_router` | Routing |
| `provider` | State management |
| `http` | ATLAS API HTTP client |
| `elliptic`, `ecdsa` | ECDSA secp256k1 signing |
| `bip39`, `hex` | BIP-39 + hex encoding |
| `flutter_secure_storage` | Secure key storage |
| `fl_chart` | Dashboard charts |
| `flutter_animate`, `lottie` | Animations |
| `google_fonts` | Typography |
| `image_picker`, `file_picker` | Media upload |
| `flutter_map`, `latlong2` | Maps |
| `infinite_scroll_pagination` | Feed paging |

---

## 14. Cerca Admin Panel

### 14.1 Stack
React + Vite + TypeScript. Icons: `lucide-react`. Target: `localhost:8080`.

### 14.2 Navigation

| Section | Pages |
|---|---|
| OVERVIEW | Dashboard |
| BLOCKCHAIN | Block Explorer, Transactions, System Contracts |
| TREASURY | Treasury (Faucet), Tx History, Arbitration |
| NODE | Node Control, Peers & Validators |

### 14.3 Pages

- **Dashboard**: Real-time block height, tx count, validator count, total staked, health indicators
- **Block Explorer**: Paginated blocks with hash, height, tx count; block drill-down
- **Transactions**: Mempool + confirmed TX list; detail view
- **System Contracts**: Live state of all 4 contracts; key/value storage inspection
- **Treasury / Faucet**: Send TCOIN to any address; treasury balance
- **Tx History**: Full paginated treasury TX history with filters
- **Arbitration**: Open dispute list; resolve dispute form (pay seller / refund buyer)
- **Node Control**: Start/Stop/Pause/Sync buttons; live state indicator; scrollable live log stream
- **Peers & Validators**: Connected P2P peer list; validator set with stake/performance/uptime/slash count; reconnect controls

---

## 15. Tokenomics Model

**Token**: TCOIN | **Genesis supply**: 1,000,000,000 (1B) | **V1 inflation**: 10 TCOIN/block (constant)

**Value flow**:
```
Treasury (1B TCOIN)
    │
    ├── Faucet → Users (1,000 TCOIN/request)
    ├── Block rewards → Validators (10 TCOIN/block)
    ├── TX fees → Block's validator
    ├── Social tips → Post TipBalance (energy)
    │   └── Comment cost: 2 TCOIN (feeds post)
    ├── Marketplace: Buyer → Escrow → Seller (minus fee%)
    └── Staking → Locked in validator bond
```

---

## 16. Deployment Architecture

### 16.1 Docker Compose (Production)

```
docker-compose up -d
# Node (port 8080) + Flutter web (port 80, nginx)
# blockchain.db persisted as Docker volume
```

### 16.2 DevNet (Local Windows)

```powershell
# ATLAS.BC0.0.1/start_devnet.ps1
# Node 1: port 8000 + API 8080 (Validator)
# Node 2: port 8001 + API 8081 (Observer, connects to Node 1)
```

### 16.3 Build Commands

```powershell
# Backend
go build -o build/atlas-node.exe cmd/main.go

# Flutter
flutter pub get && flutter build web

# Admin Panel
npm install && npm run build

# Full stack
docker-compose up -d
```

---

## 17. Security Model

| Layer | Mechanism |
|---|---|
| Transaction signing | ECDSA secp256k1 (P-256) |
| Client key storage | `flutter_secure_storage` (hardware-backed) |
| Key derivation | BIP-39 → secp256k1 private key |
| Address format | `0x` + hex(SHA-256(publicKey)[0:20]) |
| Replay protection | Account nonce incremented per TX |
| Double-spend protection | Balance check in `StateManager.updateState()` |
| Validator integrity | Slashing (3 events → ejection) + KYC required |
| Privacy proofs | gnark Groth16 scaffold (V1 mock, V2 real) |
| GDPR | gdpr-delete + gdpr-anonymize endpoints |
| P2P identity | ED25519 persistent node key |

> **WARNING**: Treasury uses a **hardcoded BIP-39 mnemonic** — testnet only. Replace with hardware-managed key for mainnet.
> **WARNING**: API CORS is fully open — must be restricted for mainnet.

---

## 18. Monitoring Subsystem

`pkg/monitoring/Monitor` collects via integration callbacks:

| Metric | Source |
|---|---|
| Block height | `BlockManager.GetBlockHeight()` |
| Transaction count | Scan all blocks |
| Pending transactions | `TxManager.GetPoolSize()` |
| Validator count | `ConsensusManager.GetAllValidators()` |
| Active peers | `P2PNode.Peerstore().Peers()` |
| Total staked | Sum of all validator stakes |
| Last block hash | `BlockManager.GetLatestBlock().Hash` |
| Sync status | `Node.State` string |

---

## 19. Current MVP Status (V1 — March 2026)

### Implemented & Working
- PBFT PoS consensus with deterministic validator selection
- 4 system smart contracts (Token, Staking, Marketplace, Governance)
- CercaVM with 30+ opcodes, gas metering, permissioned contracts
- REST API (100+ endpoints)
- SQLite persistence + JSON snapshot fallback + backup system
- libp2p P2P networking (QUIC + mDNS)
- BIP-39 wallet + ECDSA transaction signing
- Faucet with Treasury wallet
- Social platform (posts, comments, likes, tips, fossilization, Lamport clock)
- Identity system (DID, KYC, activity, reputation)
- Governance (proposals, votes, referendums, committees)
- DeFi module
- 4-shard manager
- Fast sync (snapshot-based)
- Admin panel (all 9 pages operational, live-connected to node)
- Flutter app (wallet, node dashboard, block explorer, feed, userpage, marketplace)
- Docker + DevNet deployment

### Known V1 Limitations
- ZK proofs mocked (always return `true`)
- `isLocalValidator = true` (DevNet override)
- Treasury mnemonic is hardcoded
- GovernanceManager `getTotalStake()` returns hardcoded `1,000,000`
- Governance treasury transfer action has TODO for StateManager injection
- Contract count endpoint returns `0`
- CORS fully open

---

## 20. V2 Roadmap

| Feature | Priority |
|---|---|
| gnark Groth16 ZK proof circuits | High |
| True multi-validator PBFT finality | High |
| Real stake-weighted governance voting | High |
| GovernanceManager → StateManager injection | High |
| DHT-based mainnet peer discovery | Medium |
| Ed25519 identity standardization | Medium |
| Production CORS hardening | Required for mainnet |
| Hardware treasury key management | Required for mainnet |
| Chaos/partition stress testing | Medium |
| CercaChain external developer documentation | Low |
| Mobile (Android/iOS) production builds | Low |

---

*Generated from source: `ATLAS.BC0.0.1/` + `cercaend/` + `cerca-admin-panel/` — March 28, 2026*
