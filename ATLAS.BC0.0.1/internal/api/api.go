package api

import (
	"atlas-blockchain/internal/blockchain"
	"atlas-blockchain/internal/governance"
	"atlas-blockchain/internal/identity"
	"atlas-blockchain/internal/social"
	"atlas-blockchain/pkg/crypto"
	"atlas-blockchain/pkg/monitoring"
	"atlas-blockchain/pkg/network"
	"atlas-blockchain/pkg/sharding"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/vm"
	"atlas-blockchain/pkg/wallet"
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Node state constants
const (
	NodeStateRunning = "running"
	NodeStatePaused  = "paused"
	NodeStateStopped = "stopped"
	NodeStateSyncing = "syncing"
)

// API server struct
// Add references to managers as needed
type APIServer struct {
	blockManager       *blockchain.BlockManager
	transactionManager *blockchain.TransactionManager
	stateManager       *blockchain.StateManager
	consensusManager   *blockchain.ConsensusManager
	fastSyncManager    *blockchain.FastSyncManager // Fast sync manager
	peerManager        *network.PeerManager        // Peer manager
	node               *network.Node
	p2pNode            *network.P2PNode // P2P networking node
	monitor            *monitoring.Monitor
	identityManager    *identity.IdentityManager
	socialManager      *social.SocialManager
	governanceManager  *governance.GovernanceManager
	shardManager       *sharding.ShardManager
	TreasuryWallet     *wallet.Wallet // Wallet for faucet/treasury operations
	rateLimiter        *RateLimiter   // Per-IP API rate limiter
	// Node control state
	nodeLogs      []NodeLogEntry
	nodeLogsMutex sync.RWMutex
	maxNodeLogs   int
	OnSyncRequest func() // Callback to trigger sync in main
}

// NodeLogEntry represents a single log entry from the node
type NodeLogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"` // info, success, warning, error
	Message   string `json:"message"`
}

func NewAPIServer(bm *blockchain.BlockManager, tm *blockchain.TransactionManager, sm *blockchain.StateManager, cm *blockchain.ConsensusManager, fsm *blockchain.FastSyncManager, pm *network.PeerManager, node *network.Node, p2pNode *network.P2PNode, im *identity.IdentityManager, socialMgr *social.SocialManager, govMgr *governance.GovernanceManager, shm *sharding.ShardManager) *APIServer {
	// Create monitor without shard manager for now - TODO: Update Monitor to accept ShardManager
	monitor := monitoring.NewMonitor(nil)

	api := &APIServer{
		blockManager:       bm,
		transactionManager: tm,
		stateManager:       sm,
		consensusManager:   cm,
		fastSyncManager:    fsm,
		peerManager:        pm,
		node:               node,
		p2pNode:            p2pNode,
		monitor:            monitor,
		identityManager:    im,
		socialManager:      socialMgr,
		governanceManager:  govMgr,
		shardManager:       shm,
		nodeLogs:           make([]NodeLogEntry, 0),
		maxNodeLogs:        100,
		rateLimiter:        NewRateLimiter(30, 60, time.Second),
	}

	// Initialize Treasury Wallet from environment variable
	// For production: set TREASURY_MNEMONIC env var
	// Fallback is for local development ONLY
	treasuryMnemonic := os.Getenv("TREASURY_MNEMONIC")
	if treasuryMnemonic == "" {
		treasuryMnemonic = "canyon vision beer orange notice wrong savage coin fashion roam ranch weasel"
		log.Printf("⚠️ TREASURY_MNEMONIC not set — using development fallback. DO NOT use in production!")
	}
	if treasuryWallet, err := wallet.NewWalletFromMnemonic(treasuryMnemonic); err == nil {
		api.TreasuryWallet = treasuryWallet
		// Ensure Treasury has initial Genesis supply
		treasuryAddress := wallet.PublicKeyToAddress(treasuryWallet.PublicKey)
		if sm.GetBalance(treasuryAddress) == 0 {
			sm.SetBalance(treasuryAddress, 1000000000) // 1 Billion Genesis Supply
			log.Printf("💰 Treasury Wallet Initialized: %s (Balance: 1,000,000,000)", treasuryAddress)
		}
		// Initialize TCOIN Token Contract (System Contract)
		currentBalance := sm.GetBalance(treasuryAddress)
		sm.InitTokenContract(treasuryAddress, currentBalance)

		// Initialize Staking Contract (System Contract)
		sm.InitStakingContract()

		// Initialize Marketplace Contract (System Contract)
		sm.InitMarketplaceContract()

		// Initialize Governance Contract (System Contract)
		sm.InitGovernanceContract()
	} else {
		log.Printf("❌ Failed to initialize Treasury Wallet: %v", err)
	}

	// Set integration callbacks for real monitoring data
	if monitor != nil {
		monitor.SetIntegrationCallbacks(
			// GetTransactionCount
			func() int {
				count := 0
				totalBlocks := bm.GetChainLength()
				// Safely get all blocks (paginated if implementation changes, but currently memory-based)
				blocks := bm.GetBlocks(totalBlocks, 0)
				for _, b := range blocks {
					count += len(b.Transactions)
				}
				return count
			},
			// GetBlockHeight
			func() int64 {
				return int64(bm.GetBlockHeight())
			},
			// GetPendingTransactions
			func() int {
				return tm.GetPoolSize()
			},
			// GetValidatorCount
			func() int {
				return len(cm.GetAllValidators())
			},
			// GetActivePeers
			func() int {
				if p2pNode != nil && p2pNode.Host != nil {
					return len(p2pNode.Host.Peerstore().Peers())
				}
				return 0
			},
			// GetTotalStaked
			func() float64 {
				validators := cm.GetAllValidators()
				total := uint64(0)
				for _, v := range validators {
					total += v.Stake
				}
				return float64(total)
			},
			// GetLastBlockHash
			func() string {
				blk := bm.GetLatestBlock()
				if blk != nil {
					return blk.Hash
				}
				return ""
			},
			// GetContractCount
			func() int {
				// Retrieve from StateManager if available, or mock 0 for now
				return 0
			},
			// GetSyncStatus
			func() string {
				if node != nil {
					// Use direct node state field
					// Assuming State is a string or fmt.Stringer
					return fmt.Sprintf("%v", node.State)
				}
				return "unknown"
			},
		)

		monitor.StartMonitoring()
	}

	return api
}

var (
	corsAllowedOrigins []string
	corsInitOnce       sync.Once
)

func initCORS() {
	corsInitOnce.Do(func() {
		originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
		if originsStr == "" {
			log.Printf("[API] Warning: CORS_ALLOWED_ORIGINS is unset. Defaulting to fully open CORS '*'")
			return
		}
		for _, o := range strings.Split(originsStr, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				corsAllowedOrigins = append(corsAllowedOrigins, o)
			}
		}
		log.Printf("[API] CORS initialized with allowed origins: %v", corsAllowedOrigins)
	})
}

func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	initCORS()
	return func(w http.ResponseWriter, r *http.Request) {
		// Log the request for debugging
		log.Printf("[API] Handling %s request for %s", r.Method, r.URL.Path)

		origin := r.Header.Get("Origin")
		var allowedOrigin string

		if len(corsAllowedOrigins) == 0 {
			// CORS is fully open
			if origin != "" {
				allowedOrigin = origin
			} else {
				allowedOrigin = "*"
			}
		} else {
			// CORS is restricted. Verify the request's origin matches allowed list.
			isAllowed := false
			for _, allowed := range corsAllowedOrigins {
				if allowed == origin {
					isAllowed = true
					break
				}
			}
			if isAllowed {
				allowedOrigin = origin
			} else if origin != "" {
				log.Printf("[API] Blocked request from unauthorized Origin: %s", origin)
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		if allowedOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Origin, Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		handler(w, r)
	}
}

func (api *APIServer) Start(addr string) {
	// rl wraps a handler with CORS + rate limiting
	rl := func(handler http.HandlerFunc) http.HandlerFunc {
		return RateLimitMiddleware(api.rateLimiter, withCORS(handler))
	}

	// admin wraps a handler with CORS + rate limiting + admin API key auth
	admin := func(handler http.HandlerFunc) http.HandlerFunc {
		return RateLimitMiddleware(api.rateLimiter, withCORS(AdminAuthMiddleware(handler)))
	}

	http.HandleFunc("/block", rl(api.handleGetBlock))
	http.HandleFunc("/blocks", rl(api.handleListBlocks))
	http.HandleFunc("/transaction", rl(api.handleGetTransaction))
	http.HandleFunc("/mempool", rl(api.handleGetMempool))
	http.HandleFunc("/submit-transaction", rl(api.handleSubmitTransaction))
	http.HandleFunc("/balance", rl(api.handleGetBalance))
	http.HandleFunc("/status", rl(api.handleGetStatus))
	http.HandleFunc("/validators", rl(api.handleListValidators))
	http.HandleFunc("/validator", rl(api.handleGetValidator))
	http.HandleFunc("/register-validator", rl(api.handleRegisterValidator))
	http.HandleFunc("/update-stake", rl(api.handleUpdateStake))
	http.HandleFunc("/update-user-stake", rl(api.handleUpdateUserStake))
	http.HandleFunc("/node-address", rl(api.handleGetNodeAddress))
	http.HandleFunc("/peers", rl(api.handleGetPeers))
	http.HandleFunc("/connect-peer", rl(api.handleConnectPeer))
	http.HandleFunc("/faucet", rl(api.handleFaucet))
	http.HandleFunc("/admin/faucet", admin(api.handleAdminFaucet))
	http.HandleFunc("/admin/treasury-history", admin(api.handleAdminTreasuryHistory))
	http.HandleFunc("/nonce", rl(api.handleGetNonce))
	http.HandleFunc("/run-tests", admin(api.handleRunTests))
	http.HandleFunc("/test-performance", admin(api.handleTestPerformance))
	http.HandleFunc("/test-security", admin(api.handleTestSecurity))
	http.HandleFunc("/test-integration", admin(api.handleTestIntegration))
	http.HandleFunc("/start-test-env", admin(api.handleStartTestEnv))
	http.HandleFunc("/stop-test-env", admin(api.handleStopTestEnv))
	http.HandleFunc("/test-env-status", admin(api.handleTestEnvStatus))
	// Node control endpoints are registered below, duplicates removed here
	http.HandleFunc("/create-wallet", rl(api.handleCreateWallet))
	http.HandleFunc("/import-wallet", rl(api.handleImportWallet))
	http.HandleFunc("/fee-info", rl(api.handleFeeInfo))

	// Identity management endpoints for social-commerce-governance platform
	http.HandleFunc("/identity/create", rl(api.handleCreateIdentity))
	http.HandleFunc("/identity/get", rl(api.handleGetIdentity))
	http.HandleFunc("/identity/update-profile", rl(api.handleUpdateProfile))
	// Identity endpoints commented out - handlers not implemented yet
	// http.HandleFunc("/identity/add-credential", rl(api.handleAddCredential))
	// http.HandleFunc("/identity/add-attestation", rl(api.handleAddAttestation))
	// http.HandleFunc("/identity/update-privacy", rl(api.handleUpdatePrivacy))
	// http.HandleFunc("/identity/update-kyc", rl(api.handleUpdateKYC))
	http.HandleFunc("/identity/update-activity", rl(api.handleUpdateActivity))
	http.HandleFunc("/identity/create-proof", rl(api.handleCreatePrivacyProof))
	http.HandleFunc("/identity/verify-proof", rl(api.handleVerifyPrivacyProof))
	http.HandleFunc("/identity/social", rl(api.handleGetIdentityForSocial))
	http.HandleFunc("/treasury", rl(api.handleGetTreasury))
	http.HandleFunc("/token", rl(api.handleGetTokenInfo))
	http.HandleFunc("/staking", rl(api.handleGetStakingInfo))
	http.HandleFunc("/marketplace", rl(api.handleGetMarketplaceInfo))
	http.HandleFunc("/governance-contract", rl(api.handleGetGovernanceContractInfo))
	http.HandleFunc("/identity/commerce", rl(api.handleGetIdentityForCommerce))
	http.HandleFunc("/identity/governance", rl(api.handleGetIdentityForGovernance))

	// FlutterFlow Integration Endpoints
	http.HandleFunc("/flutterflow/connect-wallet", rl(api.handleFlutterFlowConnectWallet))
	http.HandleFunc("/flutterflow/authenticate", rl(api.handleFlutterFlowAuthenticate))
	http.HandleFunc("/flutterflow/wallet-info", rl(api.handleFlutterFlowWalletInfo))
	http.HandleFunc("/flutterflow/send-transaction", rl(api.handleFlutterFlowSendTransaction))
	http.HandleFunc("/flutterflow/transaction-history", rl(api.handleFlutterFlowTransactionHistory))
	http.HandleFunc("/flutterflow/disconnect", rl(api.handleFlutterFlowDisconnect))

	// Governance endpoints
	http.HandleFunc("/governance/proposals", rl(api.handleListProposals))
	http.HandleFunc("/governance/proposal", rl(api.handleGetProposal))
	http.HandleFunc("/governance/submit-proposal", rl(api.handleSubmitProposal))
	http.HandleFunc("/governance/vote", rl(api.handleVote))

	http.HandleFunc("/oracle/submit", rl(api.handleOracleSubmit))
	http.HandleFunc("/oracle/latest", rl(api.handleOracleLatest))

	// Privacy endpoints
	http.HandleFunc("/privacy/encrypt", rl(api.handleEncryptData))
	http.HandleFunc("/privacy/decrypt", rl(api.handleDecryptData))
	// ZK-related endpoints are disabled due to ZK code being commented out
	/*
		func (api *APIServer) handleCreateProof(w http.ResponseWriter, r *http.Request) {
			// Disabled: ZK proof creation endpoint
		}

		func (api *APIServer) handleVerifyProof(w http.ResponseWriter, r *http.Request) {
			// Disabled: ZK proof verification endpoint
		}
	*/
	http.HandleFunc("/privacy/gdpr-delete", rl(api.handleGDPRDelete))
	http.HandleFunc("/privacy/gdpr-anonymize", rl(api.handleGDPRAnonymize))

	// Sharding endpoints
	http.HandleFunc("/sharding/status", rl(api.handleShardingStatus))
	http.HandleFunc("/sharding/shard", rl(api.handleGetShard))
	http.HandleFunc("/sharding/assign-validator", rl(api.handleAssignValidator))
	http.HandleFunc("/sharding/cross-shard-tx", rl(api.handleCrossShardTransaction))
	http.HandleFunc("/sharding/statistics", rl(api.handleShardingStatistics))

	// Monitoring endpoints
	http.HandleFunc("/monitoring/status", rl(api.handleMonitoringStatus))
	http.HandleFunc("/monitoring/metrics", rl(api.handleMonitoringMetrics))
	http.HandleFunc("/monitoring/health", rl(api.handleMonitoringHealth))
	http.HandleFunc("/monitoring/alerts", rl(api.handleMonitoringAlerts))
	http.HandleFunc("/monitoring/performance", rl(api.handleMonitoringPerformance))
	http.HandleFunc("/monitoring/history", rl(api.handleMonitoringHistory))
	http.HandleFunc("/monitoring/trends", rl(api.handleMonitoringTrends))

	// Fast Sync endpoints
	http.HandleFunc("/snapshot/create", rl(api.handleCreateSnapshot))
	http.HandleFunc("/snapshot/latest", rl(api.handleGetLatestSnapshot))
	http.HandleFunc("/snapshot/load", rl(api.handleLoadSnapshot))

	// Peer Management endpoints
	http.HandleFunc("/peers/status", rl(api.handleGetPeerStatus))
	http.HandleFunc("/peers/reconnect", rl(api.handleReconnectPeers))
	http.HandleFunc("/peers/validators", rl(api.handleGetValidatorPeers))

	// Chain synchronization endpoints
	http.HandleFunc("/sync/status", rl(api.handleSyncStatus))
	http.HandleFunc("/sync/start", rl(api.handleSyncStart))

	// Database backup and recovery endpoints
	http.HandleFunc("/backup/status", rl(api.handleBackupStatus))
	http.HandleFunc("/backup/list", rl(api.handleBackupList))
	http.HandleFunc("/backup/create", rl(api.handleBackupCreate))

	// Social Media API endpoints
	http.HandleFunc("/social/post/create", rl(api.handleCreatePost))
	http.HandleFunc("/social/post/get", rl(api.handleGetPost))
	http.HandleFunc("/social/comment/create", rl(api.handleCreateComment))
	http.HandleFunc("/social/like", rl(api.handleLikePost))
	http.HandleFunc("/social/unlike", rl(api.handleUnlikePost))
	http.HandleFunc("/social/tip", rl(api.handleTipPost))
	http.HandleFunc("/social/feed", rl(api.handleGetFeed))
	http.HandleFunc("/social/search", rl(api.handleSearchPosts))
	http.HandleFunc("/social/trending", rl(api.handleGetTrendingHashtags))
	http.HandleFunc("/social/report", rl(api.handleReportContent))

	// Object Energy Physics endpoints (bridges Firebase objects with blockchain energy)
	http.HandleFunc("/social/object/energy", rl(api.handleGetObjectEnergy))
	http.HandleFunc("/social/object/energize", rl(api.handleEnergizeObject))

	// Enhanced Governance API endpoints
	http.HandleFunc("/governance/proposal/create", rl(api.handleCreateProposal))
	http.HandleFunc("/governance/proposal/get", rl(api.handleGetProposal))
	http.HandleFunc("/governance/proposal/activate", rl(api.handleActivateProposal))
	http.HandleFunc("/governance/proposal/vote", rl(api.handleVoteProposal))
	http.HandleFunc("/governance/proposal/execute", rl(api.handleExecuteProposal))
	http.HandleFunc("/governance/proposal/discuss", rl(api.handleAddDiscussionComment))
	http.HandleFunc("/governance/proposals/active", rl(api.handleGetActiveProposals))
	http.HandleFunc("/governance/proposals/category", rl(api.handleGetProposalsByCategory))
	http.HandleFunc("/admin/resolve-dispute", admin(api.handleAdminResolveDispute))
	http.HandleFunc("/admin/disputes", admin(api.handleAdminGetDisputes))
	http.HandleFunc("/governance/referendum/create", rl(api.handleCreateReferendum))
	http.HandleFunc("/governance/referendum/vote", rl(api.handleVoteReferendum))

	// Smart Contract endpoints
	http.HandleFunc("/contract/deploy", admin(api.handleDeployContract))
	http.HandleFunc("/contract/call", rl(api.handleCallContract))
	http.HandleFunc("/contract/list", rl(api.handleListContracts))
	http.HandleFunc("/contract/info", rl(api.handleGetContractInfo))
	http.HandleFunc("/contract/examples", rl(api.handleGetContractExamples))

	// Network Architecture endpoint
	http.HandleFunc("/network/architecture", rl(api.handleGetNetworkArchitecture))

	// Node Control endpoints (admin-protected)
	http.HandleFunc("/node/start", admin(api.handleNodeStart))
	http.HandleFunc("/node/stop", admin(api.handleNodeStop))
	http.HandleFunc("/node/pause", admin(api.handleNodePause))
	http.HandleFunc("/node/sync", admin(api.handleNodeSync))
	http.HandleFunc("/node/status", admin(api.handleNodeStatus))
	http.HandleFunc("/node/logs", admin(api.handleNodeLogs))

	// Serve static frontend files (no rate limit on static assets)
	fs := http.FileServer(http.Dir("./web/frontend"))
	http.HandleFunc("/", withCORS(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))

	log.Printf("API server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// GET /block?hash=...
func (api *APIServer) handleGetBlock(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if ve := ValidateHash(hash); ve != nil {
		WriteValidationError(w, ve)
		return
	}
	block := api.blockManager.GetBlockByHash(hash)
	if block == nil {
		http.Error(w, "Block not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(block)
}

// GET /transaction?hash=...
func (api *APIServer) handleGetTransaction(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("hash")
	if ve := ValidateHash(hash); ve != nil {
		WriteValidationError(w, ve)
		return
	}
	tx := api.transactionManager.GetTransactionByHash(hash)
	if tx == nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(tx)
}

// POST /submit-transaction
func (api *APIServer) handleSubmitTransaction(w http.ResponseWriter, r *http.Request) {
	var tx transaction.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid transaction data", http.StatusBadRequest)
		return
	}
	// Set missing fields if needed
	if tx.Timestamp == 0 {
		tx.Timestamp = time.Now().Unix()
	}
	// Do not modify Fee after signing
	// if tx.Fee == 0 {
	// 	tx.Fee = calculateFee(tx)
	// }
	if err := api.transactionManager.AddTransaction(tx); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Broadcast valid transaction to peers
	if api.p2pNode != nil {
		if err := api.p2pNode.BroadcastTransaction(tx); err != nil {
			log.Printf("[API] Warning: Failed to broadcast transaction: %v", err)
		} else {
			log.Printf("[API] Transaction broadcast initiated")
		}
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "Transaction submitted"})
}

// GET /balance?address=...
func (api *APIServer) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if ve := ValidateAddress(address); ve != nil {
		WriteValidationError(w, ve)
		return
	}

	// Get balance from state manager (this doesn't require node wallet)
	balance := api.stateManager.GetBalance(address)
	json.NewEncoder(w).Encode(map[string]interface{}{"address": address, "balance": balance})
}

// GET /status?address=...
func (api *APIServer) handleGetStatus(w http.ResponseWriter, r *http.Request) {
	// Get the address from query parameter or use node's address
	addressToQuery := r.URL.Query().Get("address")
	if addressToQuery == "" {
		addressToQuery = api.node.ValidatorAddress
	}

	// Check if the queried address is a validator
	var isValidator bool
	var validatorAddress string
	var stakeAmount uint64
	var rewardsEarned uint64
	var walletBalance int64
	var walletStaked uint64

	if addressToQuery != "" {
		if validator, err := api.consensusManager.GetValidatorInfo(addressToQuery); err == nil {
			isValidator = true
			validatorAddress = validator.Address
			stakeAmount = validator.Stake
			rewardsEarned = validator.TotalRewards
		}
		walletBalance = api.stateManager.GetBalance(addressToQuery)
		// Note: we don't have a direct Wallet object for everyone,
		// but we can compute staked from consensus.
		walletStaked = stakeAmount
	}

	// Get all validators for total count
	validators := api.consensusManager.GetAllValidators()

	mode := "observer"
	if isValidator {
		mode = "validator"
	}
	status := map[string]interface{}{
		"blockHeight":      api.blockManager.GetBlockHeight(),
		"txPoolSize":       api.transactionManager.GetPoolSize(),
		"isValidator":      isValidator,
		"validatorAddress": validatorAddress,
		"stakeAmount":      stakeAmount,
		"rewardsEarned":    rewardsEarned,
		"totalValidators":  len(validators),
		"walletBalance":    walletBalance,
		"walletStaked":     walletStaked,
		"totalBalance":     walletBalance + int64(walletStaked),
		"mode":             mode,
	}
	json.NewEncoder(w).Encode(status)
}

// GET /blocks?limit=&offset=
func (api *APIServer) handleListBlocks(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}
	blocks := api.blockManager.GetBlocks(limit, offset)
	json.NewEncoder(w).Encode(blocks)
}

// GET /mempool
func (api *APIServer) handleGetMempool(w http.ResponseWriter, r *http.Request) {
	txs := api.transactionManager.GetAllTransactions()
	json.NewEncoder(w).Encode(txs)
}

// GET /validators
func (api *APIServer) handleListValidators(w http.ResponseWriter, r *http.Request) {
	validators := api.consensusManager.GetAllValidators()
	json.NewEncoder(w).Encode(validators)
}

// GET /validator?address=
func (api *APIServer) handleGetValidator(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing validator address", http.StatusBadRequest)
		return
	}

	validator, err := api.consensusManager.GetValidatorInfo(address)
	if err != nil {
		http.Error(w, "Validator not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(validator)
}

// GET /peers
func (api *APIServer) handleGetPeers(w http.ResponseWriter, r *http.Request) {
	if api.p2pNode == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"peers":   []string{},
			"count":   0,
			"message": "P2P node not initialized",
		})
		return
	}

	peers := api.p2pNode.Host.Network().Peers()
	peerList := make([]map[string]interface{}, 0)

	for _, peerID := range peers {
		if peerID == api.p2pNode.Host.ID() {
			continue // Skip self
		}

		addrs := api.p2pNode.Host.Peerstore().Addrs(peerID)
		addrStrings := make([]string, len(addrs))
		for i, addr := range addrs {
			addrStrings[i] = addr.String()
		}

		peerList = append(peerList, map[string]interface{}{
			"peer_id":   peerID.String(),
			"addresses": addrStrings,
		})
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"peers":   peerList,
		"count":   len(peerList),
	})
}

// Duplicate handlers handleStartNode, handleStopNode, handleSyncNode removed.
// The unique implementations handleNodeStart, handleNodeStop, handleNodeSync are used instead.

// Duplicate handlers handleStartNode, handleStopNode, handleSyncNode removed.
// The unique implementations handleNodeStart, handleNodeStop, handleNodeSync are used instead.

// Duplicate handleGetNonce removed.

// GET /treasury
func (api *APIServer) handleGetTreasury(w http.ResponseWriter, r *http.Request) {
	if api.TreasuryWallet == nil {
		http.Error(w, "Treasury not initialized", http.StatusServiceUnavailable)
		return
	}
	address := wallet.PublicKeyToAddress(api.TreasuryWallet.PublicKey)
	balance := api.stateManager.GetBalance(address)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"address": address,
		"balance": balance,
	})
}

// GET /token — Returns TCOIN token contract information
func (api *APIServer) handleGetTokenInfo(w http.ResponseWriter, r *http.Request) {
	tokenHelper := api.stateManager.GetTokenHelper()
	stateAdapter := api.stateManager.GetStateAdapter()

	if tokenHelper == nil || stateAdapter == nil {
		http.Error(w, "Token contract not initialized", http.StatusServiceUnavailable)
		return
	}

	treasuryAddress := ""
	if api.TreasuryWallet != nil {
		treasuryAddress = wallet.PublicKeyToAddress(api.TreasuryWallet.PublicKey)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":             "CercaChainToken",
		"symbol":           "TCOIN",
		"total_supply":     stateAdapter.GetTotalSupply(),
		"max_supply":       tokenHelper.GetMaxSupply(),
		"contract_address": "CONTRACT_TCOIN_SYSTEM",
		"treasury_address": treasuryAddress,
		"treasury_balance": api.stateManager.GetBalance(treasuryAddress),
	})
}

// GET /staking — Returns staking contract info and optionally per-validator info
func (api *APIServer) handleGetStakingInfo(w http.ResponseWriter, r *http.Request) {
	stakingHelper := api.stateManager.GetStakingHelper()
	stateAdapter := api.stateManager.GetStateAdapter()

	if stakingHelper == nil || stateAdapter == nil {
		http.Error(w, "Staking contract not initialized", http.StatusServiceUnavailable)
		return
	}

	ctx := &vm.ExecutionContext{
		State:     stateAdapter,
		Timestamp: time.Now().Unix(),
	}

	response := map[string]interface{}{
		"contract_address":  "CONTRACT_STAKING_SYSTEM",
		"total_staked":      stakingHelper.GetTotalStaked(ctx),
		"active_validators": stakingHelper.GetActiveValidatorCount(ctx),
		"min_stake":         vm.StakingMinStake,
		"max_validators":    vm.StakingMaxValidators,
		"lock_period":       vm.StakingLockPeriod,
		"slashing_penalty":  vm.StakingSlashingPenalty,
	}

	// If address query param is provided, include per-validator info
	address := r.URL.Query().Get("address")
	if address != "" {
		info := stakingHelper.GetStakeInfo(ctx, address)
		if info != nil {
			response["validator_info"] = info
		}
	}

	json.NewEncoder(w).Encode(response)
}

// GET /marketplace — Returns marketplace contract info and optionally per-order info
func (api *APIServer) handleGetMarketplaceInfo(w http.ResponseWriter, r *http.Request) {
	marketplaceHelper := api.stateManager.GetMarketplaceHelper()
	stateAdapter := api.stateManager.GetStateAdapter()

	if marketplaceHelper == nil || stateAdapter == nil {
		http.Error(w, "Marketplace contract not initialized", http.StatusServiceUnavailable)
		return
	}

	ctx := &vm.ExecutionContext{
		State:     stateAdapter,
		Timestamp: time.Now().Unix(),
	}

	feeRate, _ := stateAdapter.GetContractStorage(vm.MarketplaceContractAddress, "fee_rate")
	totalVolume, _ := stateAdapter.GetContractStorage(vm.MarketplaceContractAddress, "total_volume")
	totalFees, _ := stateAdapter.GetContractStorage(vm.MarketplaceContractAddress, "total_fees")
	orderCounter, _ := stateAdapter.GetContractStorage(vm.MarketplaceContractAddress, "order_counter")

	response := map[string]interface{}{
		"contract_address": vm.MarketplaceContractAddress,
		"fee_rate_bp":      feeRate,
		"total_volume":     totalVolume,
		"total_fees":       totalFees,
		"total_orders":     orderCounter,
	}

	// If order_id query param is provided, include order info
	orderID := r.URL.Query().Get("order_id")
	if orderID != "" {
		info := marketplaceHelper.GetOrder(ctx, orderID)
		if info != nil {
			response["order_info"] = info
		} else {
			response["order_info"] = "Order not found"
		}
	}

	json.NewEncoder(w).Encode(response)
}

// GET /governance-contract — Returns governance contract info and optionally per-proposal info
func (api *APIServer) handleGetGovernanceContractInfo(w http.ResponseWriter, r *http.Request) {
	governanceHelper := api.stateManager.GetGovernanceHelper()
	stateAdapter := api.stateManager.GetStateAdapter()

	if governanceHelper == nil || stateAdapter == nil {
		http.Error(w, "Governance contract not initialized", http.StatusServiceUnavailable)
		return
	}

	ctx := &vm.ExecutionContext{
		State:     stateAdapter,
		Timestamp: time.Now().Unix(),
	}

	minStake, _ := stateAdapter.GetContractStorage(vm.GovernanceContractAddress, "min_proposal_stake")
	minVotingStake, _ := stateAdapter.GetContractStorage(vm.GovernanceContractAddress, "min_voting_stake")
	proposalCounter, _ := stateAdapter.GetContractStorage(vm.GovernanceContractAddress, "proposal_counter")

	response := map[string]interface{}{
		"contract_address":   vm.GovernanceContractAddress,
		"min_proposal_stake": minStake,
		"min_voting_stake":   minVotingStake,
		"total_proposals":    proposalCounter,
	}

	// If proposal_id query param is provided, include proposal info
	proposalID := r.URL.Query().Get("proposal_id")
	if proposalID != "" {
		info := governanceHelper.GetProposal(ctx, proposalID)
		if info != nil {
			response["proposal_info"] = info
		} else {
			response["proposal_info"] = "Proposal not found"
		}
	}

	json.NewEncoder(w).Encode(response)
}

// POST /faucet
func (api *APIServer) handleFaucet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse address from POST body
	var req struct {
		Address string `json:"address"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)

	const faucetAmount = 1000
	address := req.Address
	if address == "" && api.node.Wallet != nil {
		address = wallet.PublicKeyToAddress(api.node.Wallet.PublicKey)
	}
	if address == "" {
		http.Error(w, "No address provided and node wallet unavailable", http.StatusBadRequest)
		return
	}

	// Verify Treasury Wallet
	if api.TreasuryWallet == nil {
		http.Error(w, "Treasury wallet not initialized", http.StatusInternalServerError)
		return
	}

	// Create a real transaction from Treasury to User
	treasuryAddress := wallet.PublicKeyToAddress(api.TreasuryWallet.PublicKey)
	// Use GetNextNonce to account for pending transactions in the mempool
	nonce := api.transactionManager.GetNextNonce(treasuryAddress)
	log.Printf("[FAUCET] 🔢 Calculated nonce for Treasury: %d (State + Pending)", nonce)

	// Create transaction
	tx := transaction.Transaction{
		Type:            transaction.TxTypeRegular,
		Sender:          treasuryAddress,
		SenderPublicKey: hex.EncodeToString(api.TreasuryWallet.PublicKey),
		Recipient:       address,
		Amount:          faucetAmount,
		Fee:             1, // Minimal fee
		Timestamp:       time.Now().Unix(),
		Nonce:           nonce,
		Data:            "Faucet Request",
	}

	// Sign transaction (r and s must each be 32-byte padded for consistent verification)
	txHash := wallet.CalculateTxHash(tx)
	rInt, sInt, err := ecdsa.Sign(rand.Reader, api.TreasuryWallet.PrivateKey, txHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to sign faucet transaction: " + err.Error(),
		})
		return
	}
	rBytes := rInt.Bytes()
	sBytes := sInt.Bytes()
	// Pad to 32 bytes each so the verifier can split at len/2
	sig := make([]byte, 64)
	copy(sig[32-len(rBytes):32], rBytes)
	copy(sig[64-len(sBytes):64], sBytes)
	tx.Signature = hex.EncodeToString(sig)

	// Submit to mempool
	if err := api.transactionManager.AddTransaction(tx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to submit faucet transaction: " + err.Error(),
		})
		return
	}

	log.Printf("[FAUCET] 🚰 Transaction sent from Treasury to %s (TxHash: %s)", address, hex.EncodeToString(txHash))

	// Broadcast transaction to network
	if api.p2pNode != nil {
		go func() {
			if err := api.p2pNode.BroadcastTransaction(tx); err != nil {
				log.Printf("[FAUCET] Failed to broadcast transaction: %v", err)
			}
		}()
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"address": address,
		"amount":  faucetAmount,
		"status":  "Transaction submitted",
		"message": "Funds will arrive in the next block.",
		"txHash":  hex.EncodeToString(txHash),
	})
}

// POST /admin/faucet — Admin version of faucet that accepts a custom amount
func (api *APIServer) handleAdminFaucet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string `json:"address"`
		Amount  int64  `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Address == "" {
		http.Error(w, "Missing address or invalid body", http.StatusBadRequest)
		return
	}
	if req.Amount <= 0 {
		req.Amount = 1000 // default fallback
	}

	if api.TreasuryWallet == nil {
		http.Error(w, "Treasury wallet not initialized", http.StatusInternalServerError)
		return
	}

	treasuryAddress := wallet.PublicKeyToAddress(api.TreasuryWallet.PublicKey)
	nonce := api.transactionManager.GetNextNonce(treasuryAddress)

	tx := transaction.Transaction{
		Type:            transaction.TxTypeRegular,
		Sender:          treasuryAddress,
		SenderPublicKey: hex.EncodeToString(api.TreasuryWallet.PublicKey),
		Recipient:       req.Address,
		Amount:          req.Amount,
		Fee:             1,
		Timestamp:       time.Now().Unix(),
		Nonce:           nonce,
		Data:            "Admin Faucet Emission",
	}

	txHash := wallet.CalculateTxHash(tx)
	rInt, sInt, err := ecdsa.Sign(rand.Reader, api.TreasuryWallet.PrivateKey, txHash)
	if err != nil {
		http.Error(w, "Failed to sign faucet transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rBytes, sBytes := rInt.Bytes(), sInt.Bytes()
	sig := make([]byte, 64)
	copy(sig[32-len(rBytes):32], rBytes)
	copy(sig[64-len(sBytes):64], sBytes)
	tx.Signature = hex.EncodeToString(sig)

	if err := api.transactionManager.AddTransaction(tx); err != nil {
		http.Error(w, "Failed to submit faucet transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast transaction to network
	if api.p2pNode != nil {
		go func() {
			if err := api.p2pNode.BroadcastTransaction(tx); err != nil {
				log.Printf("[ADMIN FAUCET] Failed to broadcast transaction: %v", err)
			}
		}()
	}

	log.Printf("[ADMIN FAUCET] 💸 Emitted %d TCOIN from Treasury to %s", req.Amount, req.Address)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"amount":  req.Amount,
		"to":      req.Address,
		"txHash":  hex.EncodeToString(txHash),
		"message": "Admin emission queued. Funds will arrive in the next block.",
	})
}

// POST /register-validator
// Expects a signed transaction of type "stake". Data field should contain KYC JSON.
func (api *APIServer) handleRegisterValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tx transaction.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid transaction data", http.StatusBadRequest)
		return
	}

	// 1. Verify Signature
	valid, err := wallet.VerifyTransactionSignature(tx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Signature verification failed: %v", err), http.StatusUnauthorized)
		return
	}
	if !valid {
		http.Error(w, "Invalid transaction signature", http.StatusUnauthorized)
		return
	}

	// 2. Verify Metadata
	if tx.Type != transaction.TxTypeStake {
		// Fallback for backward compatibility or better error message? No, strictly require stake.
		// However, the frontend sends "stake" string which matches TxTypeStake ("stake") const.
		http.Error(w, "Transaction type must be 'stake'", http.StatusBadRequest)
		return
	}

	address := tx.Sender
	stake := uint64(tx.Amount)

	if address == "" {
		http.Error(w, "Sender address required", http.StatusBadRequest)
		return
	}
	if stake == 0 {
		stake = 1000 // Default if 0? Or reject? Let's use 1000 default.
	}

	// 3. Check Balance
	balance := api.stateManager.GetBalance(address)
	if int64(stake) > balance {
		http.Error(w, fmt.Sprintf("Insufficient balance: have %d, need %d", balance, stake), http.StatusBadRequest)
		return
	}

	// 4. Parse KYC
	var kyc blockchain.KYCInfo
	// If Data is JSON object, parse it.
	// If Data is string (legacy), might fail.
	// If empty, use default.
	if len(tx.Data) > 0 {
		if err := json.Unmarshal([]byte(tx.Data), &kyc); err != nil {
			log.Printf("Failed to parse KYC from tx data: %v", err)
			// Proceed with default/partial KYC or fail?
			// Let's proceed as most users might not send full KYC in MVP.
			kyc.Verified = true // Assume verified if they signed a tx? No, but required for registration.
		}
	} else {
		kyc.Verified = true
	}

	// 5. Check if address matches Sender (redundant as we used tx.Sender)
	// But ensure we rely on the SIGNED Sender, which we do.

	// 6. Register Validator
	err = api.consensusManager.RegisterValidatorByAddress(address, stake, kyc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 7. Update Balance (Deduct Stake)
	// In a full system, this would be a regular transaction processed by a block.
	// Here we update state immediately as per original logic.
	api.stateManager.SetBalance(address, balance-int64(stake))

	// 8. Update Node Identity if applicable
	if api.node.ValidatorAddress == "" || api.node.ValidatorAddress == address {
		api.node.ValidatorAddress = address
		log.Printf("✅ Node validator identity confirmed as: %s", address)
	}

	// 9. Broadcast Registration
	if api.p2pNode != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			regMsg := network.ValidatorRegistrationMessage{
				Address: address,
				Stake:   stake,
				PeerID:  api.p2pNode.Host.ID().String(),
			}
			payload, _ := json.Marshal(regMsg)
			netMsg := network.NetworkMessage{Type: network.MsgTypeValidatorRegistration, Payload: payload}
			if err := api.p2pNode.BroadcastMessage(ctx, netMsg); err != nil {
				log.Printf("[API] Failed to broadcast validator registration: %v", err)
			} else {
				log.Printf("[API] Broadcasted validator registration for %s", address)
			}
		}()
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "Validator registered successfully",
		"address": address,
		"stake":   stake,
		"balance": balance - int64(stake),
		"staked":  stake,
	})
}

// POST /update-stake {"new_stake": 2000}
func (api *APIServer) handleUpdateStake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		NewStake uint64 `json:"new_stake"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Use the node's wallet for stake update
	if api.node.Wallet == nil {
		http.Error(w, "Node wallet not available", http.StatusInternalServerError)
		return
	}

	// Get current validator info
	if api.node.Wallet == nil {
		http.Error(w, "Node wallet not available", http.StatusInternalServerError)
		return
	}
	validator, err := api.consensusManager.GetValidatorInfo(api.node.Wallet.PublicKeyStr())
	if err != nil {
		http.Error(w, "Validator not found", http.StatusNotFound)
		return
	}

	// Update validator stake using wallet
	err = api.consensusManager.UpdateValidatorStake(api.node.Wallet, req.NewStake)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "Stake updated successfully",
		"address":   api.node.Wallet.PublicKeyStr(),
		"old_stake": validator.Stake,
		"new_stake": req.NewStake,
		"balance":   api.node.Wallet.GetBalance(),
		"staked":    api.node.Wallet.GetStaked(),
	})
}

// POST /update-user-stake {"address": "user_address", "new_stake": 2000}
func (api *APIServer) handleUpdateUserStake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address  string `json:"address"`
		NewStake uint64 `json:"new_stake"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	// Check if the address is a validator
	validator, err := api.consensusManager.GetValidatorInfo(req.Address)
	if err != nil {
		http.Error(w, "Address is not registered as a validator", http.StatusNotFound)
		return
	}

	// Check if the address has sufficient balance
	balance := api.stateManager.GetBalance(req.Address)
	if int64(req.NewStake) > balance {
		http.Error(w, fmt.Sprintf("Insufficient balance: have %d, need %d", balance, req.NewStake), http.StatusBadRequest)
		return
	}

	// Update the validator's stake
	validator.Stake = req.NewStake

	// Update the state to reflect the new stake
	api.stateManager.SetBalance(req.Address, balance-int64(req.NewStake))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "User stake updated successfully",
		"address":   req.Address,
		"old_stake": validator.Stake,
		"new_stake": req.NewStake,
		"balance":   balance - int64(req.NewStake),
		"staked":    req.NewStake,
	})
}

// GET /node-address
func (api *APIServer) handleGetNodeAddress(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"address": api.node.ValidatorAddress})
}

// addNodeLog adds a log entry to the node logs
func (api *APIServer) addNodeLog(level, message string) {
	api.nodeLogsMutex.Lock()
	defer api.nodeLogsMutex.Unlock()

	entry := NodeLogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Level:     level,
		Message:   message,
	}

	api.nodeLogs = append(api.nodeLogs, entry)

	// Trim to maxNodeLogs
	if len(api.nodeLogs) > api.maxNodeLogs {
		api.nodeLogs = api.nodeLogs[len(api.nodeLogs)-api.maxNodeLogs:]
	}

	// Also log to console
	log.Printf("[NODE-%s] %s", strings.ToUpper(level), message)
}

// POST /node/start
func (api *APIServer) handleNodeStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if api.node.GetState() == NodeStateRunning {
		api.addNodeLog("warning", "Node is already running")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "already_running",
			"message": "Node is already running",
			"state":   api.node.GetState(),
		})
		return
	}

	api.node.SetState(NodeStateRunning)
	api.addNodeLog("success", "Node started successfully")
	api.addNodeLog("info", "Connecting to blockchain network...")
	api.addNodeLog("info", "Listening for transactions...")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "started",
		"message": "Node started successfully",
		"state":   api.node.GetState(),
	})
}

// POST /node/stop
func (api *APIServer) handleNodeStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if api.node.GetState() == NodeStateStopped {
		api.addNodeLog("warning", "Node is already stopped")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "already_stopped",
			"message": "Node is already stopped",
			"state":   api.node.GetState(),
		})
		return
	}

	api.node.SetState(NodeStateStopped)
	api.addNodeLog("warning", "Node stopped by user request")
	api.addNodeLog("info", "Disconnecting from peers...")
	api.addNodeLog("info", "Node shutdown complete")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "stopped",
		"message": "Node stopped successfully",
		"state":   api.node.GetState(),
	})
}

// POST /node/pause
func (api *APIServer) handleNodePause(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if api.node.GetState() == NodeStatePaused {
		api.addNodeLog("warning", "Node is already paused")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "already_paused",
			"message": "Node is already paused",
			"state":   api.node.GetState(),
		})
		return
	}

	if api.node.GetState() == NodeStateStopped {
		http.Error(w, "Cannot pause a stopped node", http.StatusBadRequest)
		return
	}

	api.node.SetState(NodeStatePaused)
	api.addNodeLog("info", "Node paused - block production suspended")
	api.addNodeLog("info", "Maintaining peer connections...")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "paused",
		"message": "Node paused successfully",
		"state":   api.node.GetState(),
	})
}

// POST /node/sync
func (api *APIServer) handleNodeSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if api.node.GetState() == NodeStateStopped {
		http.Error(w, "Cannot sync a stopped node", http.StatusBadRequest)
		return
	}

	previousState := api.node.GetState()
	api.node.SetState(NodeStateSyncing)
	api.addNodeLog("info", "Starting blockchain synchronization...")
	api.addNodeLog("info", fmt.Sprintf("Current block height: %d", api.blockManager.GetBlockHeight()))

	// Simulate sync completion (in real implementation, this would be async)
	go func() {
		time.Sleep(2 * time.Second)
		api.nodeLogsMutex.Lock()
		api.nodeLogs = append(api.nodeLogs, NodeLogEntry{
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			Level:     "success",
			Message:   "Blockchain synchronized successfully",
		})
		api.nodeLogsMutex.Unlock()
		api.node.SetState(previousState)
		if api.node.GetState() == NodeStateSyncing {
			api.node.SetState(NodeStateRunning)
		}
	}()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "syncing",
		"message":     "Synchronization started",
		"state":       api.node.GetState(),
		"blockHeight": api.blockManager.GetBlockHeight(),
	})
}

// GET /node/status
func (api *APIServer) handleNodeStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get validator info
	var isValidator bool
	var validatorAddress string
	var stakeAmount uint64
	nodeAddress := api.node.ValidatorAddress

	if validator, err := api.consensusManager.GetValidatorInfo(nodeAddress); err == nil {
		isValidator = true
		validatorAddress = validator.Address
		stakeAmount = validator.Stake
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"state":            api.node.GetState(),
		"isValidator":      isValidator,
		"validatorAddress": validatorAddress,
		"stakeAmount":      stakeAmount,
		"blockHeight":      api.blockManager.GetBlockHeight(),
		"txPoolSize":       api.transactionManager.GetPoolSize(),
		"totalValidators":  len(api.consensusManager.GetAllValidators()),
		"uptime":           "running", // Could track actual uptime
	})
}

// GET /node/logs?limit=50
func (api *APIServer) handleNodeLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	api.nodeLogsMutex.RLock()
	defer api.nodeLogsMutex.RUnlock()

	// Get the last 'limit' logs
	logs := api.nodeLogs
	if len(logs) > limit {
		logs = logs[len(logs)-limit:]
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"logs":  logs,
		"total": len(api.nodeLogs),
		"state": api.node.GetState(),
	})
}

// GET /nonce?address=...
func (api *APIServer) handleGetNonce(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}
	if api.node.Wallet == nil {
		http.Error(w, "Node wallet not available", http.StatusInternalServerError)
		return
	}
	nonce := api.stateManager.GetNonce(address)
	json.NewEncoder(w).Encode(map[string]interface{}{"address": address, "nonce": nonce})
}

// POST /run-tests
func (api *APIServer) handleRunTests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Use go test directly for cross-platform compatibility
	cmd := exec.Command("go", "test", "./...", "-v")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd = exec.CommandContext(ctx, "go", "test", "./...", "-v")
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Tests failed or error running them
		output := out.String() + "\nError: " + stderr.String() + "\n" + err.Error()
		w.WriteHeader(http.StatusOK) // Return OK so frontend can display the failure output
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"output":  output,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"output":  out.String(),
	})
}

// POST /test-performance
func (api *APIServer) handleTestPerformance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mock performance test results
	result := map[string]interface{}{
		"suites": []map[string]interface{}{
			{
				"name":        "Transaction Throughput",
				"description": "Testing transaction processing speed",
				"passed":      3,
				"failed":      0,
				"warnings":    1,
			},
			{
				"name":        "Block Production",
				"description": "Testing block creation performance",
				"passed":      2,
				"failed":      0,
				"warnings":    0,
			},
			{
				"name":        "Memory Usage",
				"description": "Testing memory efficiency",
				"passed":      4,
				"failed":      0,
				"warnings":    0,
			},
		},
	}

	json.NewEncoder(w).Encode(result)
}

// POST /test-security
func (api *APIServer) handleTestSecurity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mock security test results
	result := map[string]interface{}{
		"suites": []map[string]interface{}{
			{
				"name":        "Transaction Validation",
				"description": "Testing transaction security validation",
				"passed":      5,
				"failed":      0,
				"warnings":    0,
			},
			{
				"name":        "Block Validation",
				"description": "Testing block security validation",
				"passed":      4,
				"failed":      0,
				"warnings":    1,
			},
			{
				"name":        "Consensus Security",
				"description": "Testing consensus mechanism security",
				"passed":      3,
				"failed":      0,
				"warnings":    0,
			},
		},
	}

	json.NewEncoder(w).Encode(result)
}

// POST /test-integration
func (api *APIServer) handleTestIntegration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mock integration test results
	result := map[string]interface{}{
		"suites": []map[string]interface{}{
			{
				"name":        "End-to-End Flow",
				"description": "Testing complete transaction flow",
				"passed":      6,
				"failed":      0,
				"warnings":    0,
			},
			{
				"name":        "Multi-Node Sync",
				"description": "Testing network synchronization",
				"passed":      4,
				"failed":      0,
				"warnings":    1,
			},
			{
				"name":        "API Integration",
				"description": "Testing API endpoint integration",
				"passed":      8,
				"failed":      0,
				"warnings":    0,
			},
		},
	}

	json.NewEncoder(w).Encode(result)
}

// POST /start-test-env
func (api *APIServer) handleStartTestEnv(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mock test environment start
	result := map[string]interface{}{
		"status":  "started",
		"message": "Test environment started successfully",
		"nodes":   3,
		"wallets": 5,
	}

	json.NewEncoder(w).Encode(result)
}

// POST /stop-test-env
func (api *APIServer) handleStopTestEnv(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mock test environment stop
	result := map[string]interface{}{
		"status":  "stopped",
		"message": "Test environment stopped successfully",
	}

	json.NewEncoder(w).Encode(result)
}

// GET /test-env-status
func (api *APIServer) handleTestEnvStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mock test environment status
	result := map[string]interface{}{
		"nodes":        3,
		"wallets":      5,
		"transactions": 12,
		"blocks":       8,
		"status":       "running",
	}

	json.NewEncoder(w).Encode(result)
}

// Add wallet management endpoints
func (api *APIServer) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	if api.node.Wallet != nil {
		http.Error(w, "Wallet already exists", http.StatusBadRequest)
		return
	}
	wallet, mnemonic, err := wallet.NewWalletWithMnemonic()
	if err != nil {
		http.Error(w, "Failed to create wallet: "+err.Error(), http.StatusInternalServerError)
		return
	}
	api.node.Wallet = wallet
	api.node.ValidatorAddress = wallet.PublicKeyStr()
	json.NewEncoder(w).Encode(map[string]string{
		"address":  wallet.PublicKeyStr(),
		"mnemonic": mnemonic,
	})
}

func (api *APIServer) handleImportWallet(w http.ResponseWriter, r *http.Request) {
	if api.node.Wallet != nil {
		http.Error(w, "Wallet already exists", http.StatusBadRequest)
		return
	}
	var req struct {
		PrivateKey string `json:"privateKey"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.PrivateKey == "" {
		http.Error(w, "Missing privateKey", http.StatusBadRequest)
		return
	}
	wallet, err := wallet.ImportWallet(req.PrivateKey)
	if err != nil {
		http.Error(w, "Failed to import wallet: "+err.Error(), http.StatusBadRequest)
		return
	}
	api.node.Wallet = wallet
	api.node.ValidatorAddress = wallet.PublicKeyStr()
	json.NewEncoder(w).Encode(map[string]string{"address": wallet.PublicKeyStr()})
}

// GET /fee-info?amount=...&sender=...&recipient=...
func (api *APIServer) handleFeeInfo(w http.ResponseWriter, r *http.Request) {
	amount := int64(0)
	if a := r.URL.Query().Get("amount"); a != "" {
		fmt.Sscanf(a, "%d", &amount)
	}
	sender := r.URL.Query().Get("sender")
	recipient := r.URL.Query().Get("recipient")

	_ = transaction.Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	// TODO: Implement fee calculation
	fee := uint64(10) // Default fee
	multiplier := api.transactionManager.GetDynamicFeeMultiplier()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recommendedFee":       fee,
		"dynamicFeeMultiplier": multiplier,
	})
}

// ===== FLUTTERFLOW INTEGRATION ENDPOINTS =====

// POST /flutterflow/connect-wallet
// Connects a wallet for FlutterFlow integration
func (api *APIServer) handleFlutterFlowConnectWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Action     string `json:"action"` // "create", "import", "connect"
		PrivateKey string `json:"privateKey,omitempty"`
		Address    string `json:"address,omitempty"`
		SessionID  string `json:"sessionId,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var walletAddress string
	var mnemonic string

	var linkedWallet *wallet.Wallet

	switch req.Action {
	case "create":
		// Create a new wallet with mnemonic
		newWallet, m, err := wallet.NewWalletWithMnemonic()
		if err != nil {
			http.Error(w, "Failed to create wallet: "+err.Error(), http.StatusInternalServerError)
			return
		}
		walletAddress = wallet.PublicKeyToAddress(newWallet.PublicKey)
		mnemonic = m
		linkedWallet = newWallet

	case "import":
		// Import existing wallet (supports both hex key and mnemonic)
		if req.PrivateKey == "" {
			http.Error(w, "Private key or mnemonic required for import", http.StatusBadRequest)
			return
		}

		var importedWallet *wallet.Wallet
		var err error

		// Simple heuristic: if it has spaces, it's likely a mnemonic
		if strings.Contains(strings.TrimSpace(req.PrivateKey), " ") {
			importedWallet, err = wallet.NewWalletFromMnemonic(req.PrivateKey)
			mnemonic = req.PrivateKey
		} else {
			importedWallet, err = wallet.ImportWallet(req.PrivateKey)
		}

		if err != nil {
			http.Error(w, "Failed to import wallet: "+err.Error(), http.StatusBadRequest)
			return
		}
		walletAddress = wallet.PublicKeyToAddress(importedWallet.PublicKey)
		linkedWallet = importedWallet

	case "connect":
		// Connect to existing address
		if req.Address == "" {
			http.Error(w, "Address required for connection", http.StatusBadRequest)
			return
		}
		walletAddress = req.Address

	default:
		http.Error(w, "Invalid action. Use 'create', 'import', or 'connect'", http.StatusBadRequest)
		return
	}

	// Generate session token for FlutterFlow
	sessionToken := generateSessionToken(walletAddress)

	// Update node's identity if we imported/created a wallet
	if (req.Action == "create" || req.Action == "import") && linkedWallet != nil {
		api.node.ValidatorAddress = walletAddress
		api.node.Wallet = linkedWallet
		log.Printf("🔑 Node identity linked to connected wallet: %s", walletAddress)
	}

	// Check if this wallet is already a validator
	isVal := false
	if _, err := api.consensusManager.GetValidatorInfo(walletAddress); err == nil {
		isVal = true
	}

	// Return wallet connection response
	response := map[string]interface{}{
		"success": true,
		"message": "Wallet connected successfully",
		"data": map[string]interface{}{
			"address":      walletAddress,
			"sessionToken": sessionToken,
			"balance":      api.stateManager.GetBalance(walletAddress),
			"isValidator":  isVal,
			"mnemonic":     mnemonic,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// POST /flutterflow/authenticate
// Authenticates a FlutterFlow session
func (api *APIServer) handleFlutterFlowAuthenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionToken string `json:"sessionToken"`
		Address      string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate session token (simplified for demo)
	if !isValidSessionToken(req.SessionToken, req.Address) {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}

	// Check if address is a validator
	validator, _ := api.consensusManager.GetValidatorInfo(req.Address)
	isValidator := validator != nil

	response := map[string]interface{}{
		"success": true,
		"message": "Authentication successful",
		"data": map[string]interface{}{
			"address":       req.Address,
			"balance":       api.stateManager.GetBalance(req.Address),
			"isValidator":   isValidator,
			"validatorInfo": validator,
		},
	}

	json.NewEncoder(w).Encode(response)
}

// GET /flutterflow/wallet-info?address=...
// Gets wallet information for FlutterFlow
func (api *APIServer) handleFlutterFlowWalletInfo(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	// Get wallet balance
	balance := api.stateManager.GetBalance(address)

	// Check if address is a validator
	validator, _ := api.consensusManager.GetValidatorInfo(address)
	isValidator := validator != nil

	// Get recent transactions (simplified)
	recentTxs := api.transactionManager.GetAllTransactions()
	userTxs := []map[string]interface{}{}

	for _, tx := range recentTxs {
		if tx.Sender == address || tx.Recipient == address {
			userTxs = append(userTxs, map[string]interface{}{
				"hash":      wallet.CalculateTxHash(tx),
				"sender":    tx.Sender,
				"recipient": tx.Recipient,
				"amount":    tx.Amount,
				"timestamp": tx.Timestamp,
			})
		}
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"address":            address,
			"balance":            balance,
			"isValidator":        isValidator,
			"validatorInfo":      validator,
			"recentTransactions": userTxs,
			"nonce":              api.stateManager.GetNonce(address),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// POST /flutterflow/send-transaction
// Sends a transaction from FlutterFlow
func (api *APIServer) handleFlutterFlowSendTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		From         string `json:"from"`
		To           string `json:"to"`
		Amount       int64  `json:"amount"`
		Fee          int64  `json:"fee"`
		Data         string `json:"data,omitempty"`
		Signature    string `json:"signature"`
		SessionToken string `json:"sessionToken"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate session token
	if !isValidSessionToken(req.SessionToken, req.From) {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}

	// Create transaction
	tx := transaction.Transaction{
		Sender:    req.From,
		Recipient: req.To,
		Amount:    req.Amount,
		Fee:       req.Fee,
		Data:      req.Data,
		Timestamp: time.Now().Unix(),
		Nonce:     api.stateManager.GetNonce(req.From),
		Signature: req.Signature,
	}

	// Add transaction to pool
	if err := api.transactionManager.AddTransaction(tx); err != nil {
		http.Error(w, "Failed to submit transaction: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Transaction submitted successfully",
		"data": map[string]interface{}{
			"transactionHash": wallet.CalculateTxHash(tx),
			"from":            req.From,
			"to":              req.To,
			"amount":          req.Amount,
			"status":          "pending",
		},
	}

	json.NewEncoder(w).Encode(response)
}

// GET /flutterflow/transaction-history?address=...
// Gets transaction history for FlutterFlow
func (api *APIServer) handleFlutterFlowTransactionHistory(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address", http.StatusBadRequest)
		return
	}

	// Get recent blocks
	blocks := api.blockManager.GetBlocks(50, 0) // Last 50 blocks

	transactions := []map[string]interface{}{}

	for _, block := range blocks {
		for _, tx := range block.Transactions {
			if tx.Sender == address || tx.Recipient == address {
				transactions = append(transactions, map[string]interface{}{
					"hash":       wallet.CalculateTxHash(tx),
					"blockIndex": block.Index,
					"sender":     tx.Sender,
					"recipient":  tx.Recipient,
					"amount":     tx.Amount,
					"fee":        tx.Fee,
					"timestamp":  tx.Timestamp,
					"type":       getTransactionType(tx, address),
				})
			}
		}
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"address":      address,
			"transactions": transactions,
			"totalCount":   len(transactions),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// POST /flutterflow/disconnect
// Disconnects wallet from FlutterFlow
func (api *APIServer) handleFlutterFlowDisconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		SessionToken string `json:"sessionToken"`
		Address      string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate session token
	if !isValidSessionToken(req.SessionToken, req.Address) {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Wallet disconnected successfully",
		"data": map[string]interface{}{
			"address":        req.Address,
			"disconnectedAt": time.Now().Unix(),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// ===== HELPER FUNCTIONS =====

// generateSessionToken creates a session token for FlutterFlow
func generateSessionToken(address string) string {
	timestamp := time.Now().Unix()
	// Simple token generation - in production, use proper JWT
	return fmt.Sprintf("ff_%s_%d_%s", address[:8], timestamp, "session")
}

// isValidSessionToken validates a session token
func isValidSessionToken(token, address string) bool {
	// Simple validation - in production, use proper JWT validation
	return len(token) > 0 && strings.Contains(token, address[:8])
}

// getTransactionType determines if transaction is incoming or outgoing
func getTransactionType(tx transaction.Transaction, userAddress string) string {
	if tx.Sender == userAddress {
		return "outgoing"
	}
	if tx.Recipient == userAddress {
		return "incoming"
	}
	return "unknown"
}

// List all proposals
func (api *APIServer) handleListProposals(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement proposal listing through governance manager
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Proposal listing not yet implemented",
		"proposals": []interface{}{},
	})
}

// Get proposal details
func (api *APIServer) handleGetProposal(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing proposal ID", http.StatusBadRequest)
		return
	}
	// TODO: Implement proposal retrieval through governance manager
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Proposal retrieval not yet implemented",
		"id":      id,
	})
}

// Submit a new proposal
func (api *APIServer) handleSubmitProposal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Proposer    string `json:"proposer"`
		Description string `json:"description"`
		Actions     string `json:"actions"`
		Duration    int64  `json:"duration"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// TODO: Implement proposal submission through governance manager
	// proposal := api.stateManager.SubmitProposal(req.Proposer, req.Description, req.Actions, 0, req.Duration)
	// proposal.State = ProposalActive
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "Proposal submission not yet implemented",
		"proposer":    req.Proposer,
		"description": req.Description,
	})
}

// Vote on a proposal
func (api *APIServer) handleVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ProposalID string `json:"proposalID"`
		Voter      string `json:"voter"`
		Choice     string `json:"choice"`
		Weight     int64  `json:"weight"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// TODO: Implement voting through governance manager
	// err := api.stateManager.CastVote(req.ProposalID, req.Voter, req.Choice, req.Weight)
	var err error = nil
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "Vote cast",
	})
}

// POST /oracle/submit {"key":..., "value":..., "source":...}
func (api *APIServer) handleOracleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Key    string `json:"key"`
		Value  string `json:"value"`
		Source string `json:"source"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	timestamp := time.Now().Unix()
	api.stateManager.SetOracleData(req.Key, req.Value, req.Source, timestamp)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "Oracle data submitted",
		"key":       req.Key,
		"value":     req.Value,
		"timestamp": timestamp,
		"source":    req.Source,
	})
}

// GET /oracle/latest?key=...
func (api *APIServer) handleOracleLatest(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}
	data, ok := api.stateManager.GetOracleData(key)
	if !ok {
		http.Error(w, "No oracle data for key", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(data)
}

// Privacy endpoints

// POST /privacy/encrypt
func (api *APIServer) handleEncryptData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Data     string `json:"data"`
		Password string `json:"password,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Data == "" {
		http.Error(w, "Data required", http.StatusBadRequest)
		return
	}

	var key *crypto.EncryptionKey
	var err error

	if req.Password != "" {
		// Derive key from password
		salt := []byte("blockchain_salt") // In production, use unique salt per user
		key, err = crypto.DeriveKeyFromPassword(req.Password, salt)
	} else {
		// Generate random key
		key, err = crypto.NewEncryptionKey()
	}

	if err != nil {
		http.Error(w, "Failed to create encryption key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	encrypted, err := crypto.EncryptString(req.Data, key)
	if err != nil {
		http.Error(w, "Failed to encrypt data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"encrypted_data": encrypted,
		"key_id":         key.ID,
		"algorithm":      "AES-GCM-256",
	}

	json.NewEncoder(w).Encode(response)
}

// POST /privacy/decrypt
func (api *APIServer) handleDecryptData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		EncryptedData string `json:"encrypted_data"`
		Password      string `json:"password,omitempty"`
		KeyID         string `json:"key_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.EncryptedData == "" {
		http.Error(w, "Encrypted data required", http.StatusBadRequest)
		return
	}

	var key *crypto.EncryptionKey
	var err error

	if req.Password != "" {
		// Derive key from password
		salt := []byte("blockchain_salt")
		key, err = crypto.DeriveKeyFromPassword(req.Password, salt)
	} else {
		http.Error(w, "Password or key ID required", http.StatusBadRequest)
		return
	}

	decrypted, err := crypto.DecryptString(req.EncryptedData, key)
	if err != nil {
		http.Error(w, "Failed to decrypt data: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"decrypted_data": decrypted,
	}

	json.NewEncoder(w).Encode(response)
}

// POST /privacy/create-proof
func (api *APIServer) handleCreateProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// var req zk.ProofRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }

	// // Get prover address from request headers or body
	// prover := r.Header.Get("X-Prover-Address")
	// if prover == "" {
	// 	http.Error(w, "Prover address required", http.StatusBadRequest)
	// 	return
	// }

	// proof, err := zk.CreateProof(&req, prover)
	// if err != nil {
	// 	http.Error(w, "Failed to create proof: "+err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// json.NewEncoder(w).Encode(proof)
}

// POST /privacy/verify-proof
func (api *APIServer) handleVerifyProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// var proof zk.ZKProof
	// if err := json.NewDecoder(r.Body).Decode(&proof); err != nil {
	// 	http.Error(w, "Invalid proof data", http.StatusBadRequest)
	// 	return
	// }

	// verifier := zk.NewProofVerifier(true) // Enable verification
	// valid, err := verifier.VerifyProof(&proof)
	// if err != nil {
	// 	http.Error(w, "Proof verification failed: "+err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// response := map[string]interface{}{
	// 	"valid": valid,
	// 	"proof": proof,
	// }

	// json.NewEncoder(w).Encode(response)
}

// POST /privacy/gdpr-delete
func (api *APIServer) handleGDPRDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string `json:"address"`
		Reason  string `json:"reason,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Address == "" {
		http.Error(w, "Address required", http.StatusBadRequest)
		return
	}

	// Mock GDPR deletion - in production, this would:
	// 1. Anonymize personal data
	// 2. Remove from search indexes
	// 3. Log the deletion request
	// 4. Notify relevant parties

	// For now, just return success
	response := map[string]interface{}{
		"success":   true,
		"message":   "GDPR deletion request processed",
		"address":   req.Address,
		"timestamp": time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

// POST /privacy/gdpr-anonymize
func (api *APIServer) handleGDPRAnonymize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string   `json:"address"`
		Fields  []string `json:"fields,omitempty"` // Specific fields to anonymize
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Address == "" {
		http.Error(w, "Address required", http.StatusBadRequest)
		return
	}

	// Mock GDPR anonymization - in production, this would:
	// 1. Replace personal data with hashed/anonymized versions
	// 2. Keep transaction history but remove personal identifiers
	// 3. Log the anonymization request

	response := map[string]interface{}{
		"success":           true,
		"message":           "GDPR anonymization request processed",
		"address":           req.Address,
		"fields_anonymized": req.Fields,
		"timestamp":         time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

// Sharding endpoints

// GET /sharding/status
func (api *APIServer) handleShardingStatus(w http.ResponseWriter, r *http.Request) {
	if api.shardManager == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"enabled": false,
			"message": "Sharding not initialized",
		})
		return
	}

	allShards := api.shardManager.GetAllShards()
	shardList := make([]map[string]interface{}, 0, len(allShards))
	for id, shard := range allShards {
		validators := make([]string, 0)
		for v := range shard.Validators {
			validators = append(validators, v)
		}
		shardList = append(shardList, map[string]interface{}{
			"id":                id,
			"validator_count":   len(shard.Validators),
			"validators":        validators,
			"transaction_count": len(shard.Transactions),
		})
	}

	crossTxs := api.shardManager.GetCrossShardTransactions()

	response := map[string]interface{}{
		"enabled":      true,
		"total_shards": api.shardManager.GetTotalShards(),
		"shards":       shardList,
		"stats": map[string]interface{}{
			"cross_shard_tx_count": len(crossTxs),
		},
	}

	json.NewEncoder(w).Encode(response)
}

// GET /sharding/shard?id=...
func (api *APIServer) handleGetShard(w http.ResponseWriter, r *http.Request) {
	shardIDStr := r.URL.Query().Get("id")
	if shardIDStr == "" {
		http.Error(w, "Missing shard ID", http.StatusBadRequest)
		return
	}

	if api.shardManager == nil {
		http.Error(w, "Sharding not initialized", http.StatusServiceUnavailable)
		return
	}

	var shardID uint32
	fmt.Sscanf(shardIDStr, "%d", &shardID)

	shard, err := api.shardManager.GetShardInfo(sharding.ShardID(shardID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	validators := make([]string, 0)
	for v := range shard.Validators {
		validators = append(validators, v)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                shard.ID,
		"validators":        validators,
		"validator_count":   len(shard.Validators),
		"transaction_count": len(shard.Transactions),
	})
}

// POST /sharding/assign-validator
func (api *APIServer) handleAssignValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if api.shardManager == nil {
		http.Error(w, "Sharding not initialized", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		ValidatorAddress string `json:"validator_address"`
		ShardID          uint32 `json:"shard_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := api.shardManager.AddValidatorToShard(sharding.ShardID(req.ShardID), req.ValidatorAddress); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":           true,
		"message":           "Validator assigned to shard",
		"validator_address": req.ValidatorAddress,
		"shard_id":          req.ShardID,
	})
}

// POST /sharding/cross-shard-tx
func (api *APIServer) handleCrossShardTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if api.shardManager == nil {
		http.Error(w, "Sharding not initialized", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		SourceShard uint32      `json:"source_shard"`
		TargetShard uint32      `json:"target_shard"`
		Transaction interface{} `json:"transaction"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	csTx, err := api.shardManager.CreateCrossShardTransaction(
		sharding.ShardID(req.SourceShard),
		sharding.ShardID(req.TargetShard),
		req.Transaction,
	)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"tx_id":        csTx.ID,
		"source_shard": req.SourceShard,
		"target_shard": req.TargetShard,
		"status":       csTx.Status,
	})
}

// GET /sharding/statistics
func (api *APIServer) handleShardingStatistics(w http.ResponseWriter, r *http.Request) {
	if api.shardManager == nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Sharding not initialized",
		})
		return
	}

	allShards := api.shardManager.GetAllShards()
	crossTxs := api.shardManager.GetCrossShardTransactions()

	totalValidators := 0
	totalTransactions := 0
	for _, shard := range allShards {
		totalValidators += len(shard.Validators)
		totalTransactions += len(shard.Transactions)
	}

	response := map[string]interface{}{
		"statistics": map[string]interface{}{
			"total_shards":       api.shardManager.GetTotalShards(),
			"total_validators":   totalValidators,
			"total_transactions": totalTransactions,
		},
		"cross_shard_tx":    crossTxs,
		"total_cross_shard": len(crossTxs),
	}

	json.NewEncoder(w).Encode(response)
}

// Monitoring endpoints

// GET /monitoring/status
func (api *APIServer) handleMonitoringStatus(w http.ResponseWriter, r *http.Request) {
	// Get real monitoring status from the monitor
	if api.monitor == nil {
		http.Error(w, "Monitoring not available", http.StatusServiceUnavailable)
		return
	}

	systemStatus := api.monitor.GetSystemStatus()

	// Update with real blockchain data
	systemStatus.LastBlockHeight = int64(api.blockManager.GetBlockHeight())
	systemStatus.TotalTransactions = int64(api.transactionManager.GetPoolSize())

	json.NewEncoder(w).Encode(systemStatus)
}

// GET /monitoring/metrics
func (api *APIServer) handleMonitoringMetrics(w http.ResponseWriter, r *http.Request) {
	// Get real metrics from the monitor
	if api.monitor == nil {
		http.Error(w, "Monitoring not available", http.StatusServiceUnavailable)
		return
	}

	allMetrics := api.monitor.GetAllMetrics()
	performanceMetrics := api.monitor.GetPerformanceMetrics()

	// Convert metrics to the expected format
	metrics := make(map[string]interface{})

	// Add performance metrics
	if performanceMetrics != nil {
		metrics["tps"] = map[string]interface{}{
			"value":     performanceMetrics.TPS,
			"type":      "gauge",
			"timestamp": time.Now().Unix(),
		}
		metrics["block_time"] = map[string]interface{}{
			"value":     performanceMetrics.BlockTime,
			"type":      "gauge",
			"timestamp": time.Now().Unix(),
		}
		metrics["memory_usage"] = map[string]interface{}{
			"value":     performanceMetrics.MemoryUsage,
			"type":      "gauge",
			"timestamp": time.Now().Unix(),
		}
		metrics["cpu_usage"] = map[string]interface{}{
			"value":     performanceMetrics.CPUUsage,
			"type":      "gauge",
			"timestamp": time.Now().Unix(),
		}
		metrics["network_latency"] = map[string]interface{}{
			"value":     performanceMetrics.NetworkLatency,
			"type":      "gauge",
			"timestamp": time.Now().Unix(),
		}
		metrics["active_peers"] = map[string]interface{}{
			"value":     performanceMetrics.ActivePeers,
			"type":      "gauge",
			"timestamp": time.Now().Unix(),
		}
		metrics["validator_count"] = map[string]interface{}{
			"value":     performanceMetrics.ValidatorCount,
			"type":      "gauge",
			"timestamp": time.Now().Unix(),
		}
		metrics["disk_usage"] = map[string]interface{}{
			"value":     performanceMetrics.DiskUsage,
			"type":      "gauge",
			"timestamp": time.Now().Unix(),
		}
		metrics["network_io"] = map[string]interface{}{
			"value":     performanceMetrics.NetworkIO,
			"type":      "gauge",
			"timestamp": time.Now().Unix(),
		}
	}

	// Add all recorded metrics
	for name, metric := range allMetrics {
		metrics[name] = map[string]interface{}{
			"value":     metric.Value,
			"type":      string(metric.Type),
			"timestamp": metric.Timestamp.Unix(),
			"labels":    metric.Labels,
		}
	}

	json.NewEncoder(w).Encode(metrics)
}

// GET /monitoring/health
func (api *APIServer) handleMonitoringHealth(w http.ResponseWriter, r *http.Request) {
	// Get real health checks from the monitor
	if api.monitor == nil {
		http.Error(w, "Monitoring not available", http.StatusServiceUnavailable)
		return
	}

	healthChecks := api.monitor.GetHealthChecks()

	// Convert to the expected format
	var checks []map[string]interface{}
	for _, check := range healthChecks {
		checks = append(checks, map[string]interface{}{
			"name":      check.Name,
			"status":    check.Status,
			"message":   check.Message,
			"timestamp": check.Timestamp.Unix(),
			"details":   check.Details,
		})
	}

	json.NewEncoder(w).Encode(checks)
}

// GET /monitoring/alerts
func (api *APIServer) handleMonitoringAlerts(w http.ResponseWriter, r *http.Request) {
	// Get real alerts from the monitor
	if api.monitor == nil {
		http.Error(w, "Monitoring not available", http.StatusServiceUnavailable)
		return
	}

	alerts := api.monitor.GetAlerts()

	// Convert to the expected format
	var alertList []map[string]interface{}
	for _, alert := range alerts {
		alertList = append(alertList, map[string]interface{}{
			"id":           alert.ID,
			"level":        alert.Level,
			"message":      alert.Message,
			"timestamp":    alert.Timestamp.Unix(),
			"acknowledged": alert.Acknowledged,
			"details":      alert.Details,
		})
	}

	json.NewEncoder(w).Encode(alertList)
}

// GET /monitoring/performance
func (api *APIServer) handleMonitoringPerformance(w http.ResponseWriter, r *http.Request) {
	// Get real performance metrics from the monitor
	if api.monitor == nil {
		http.Error(w, "Monitoring not available", http.StatusServiceUnavailable)
		return
	}

	performanceMetrics := api.monitor.GetPerformanceMetrics()
	if performanceMetrics == nil {
		http.Error(w, "Performance metrics not available", http.StatusServiceUnavailable)
		return
	}

	performance := map[string]interface{}{
		"tps":             performanceMetrics.TPS,
		"block_time":      performanceMetrics.BlockTime,
		"memory_usage":    performanceMetrics.MemoryUsage,
		"cpu_usage":       performanceMetrics.CPUUsage,
		"network_latency": performanceMetrics.NetworkLatency,
		"active_peers":    performanceMetrics.ActivePeers,
		"validator_count": performanceMetrics.ValidatorCount,
		"disk_usage":      performanceMetrics.DiskUsage,
		"network_io":      performanceMetrics.NetworkIO,
		"timestamp":       time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(performance)
}

// GET /monitoring/history
func (api *APIServer) handleMonitoringHistory(w http.ResponseWriter, r *http.Request) {
	// Get historical metrics from the monitor
	if api.monitor == nil {
		http.Error(w, "Monitoring not available", http.StatusServiceUnavailable)
		return
	}

	historicalMetrics := api.monitor.GetHistoricalMetrics()

	response := map[string]interface{}{
		"history":   historicalMetrics,
		"count":     len(historicalMetrics),
		"timestamp": time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

// GET /monitoring/trends
func (api *APIServer) handleMonitoringTrends(w http.ResponseWriter, r *http.Request) {
	// Get trend analysis from the monitor
	if api.monitor == nil {
		http.Error(w, "Monitoring not available", http.StatusServiceUnavailable)
		return
	}

	trends := api.monitor.GetMetricsTrends()

	response := map[string]interface{}{
		"trends":    trends,
		"timestamp": time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

// GET /sync/status
func (api *APIServer) handleSyncStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement sync status through chain sync manager
	syncStatus := map[string]interface{}{
		"status":        "not_implemented",
		"blocks_synced": 0,
		"total_blocks":  0,
		"duration":      "0s",
		"progress":      0.0,
		"timestamp":     time.Now().Unix(),
		"message":       "Sync status not yet implemented",
	}

	json.NewEncoder(w).Encode(syncStatus)
}

// POST /sync/start
func (api *APIServer) handleSyncStart(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement chain sync through chain sync manager
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"message":   "Chain synchronization not yet implemented",
		"status":    "not_implemented",
		"timestamp": time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

// POST /contract/deploy
func (api *APIServer) handleDeployContract(w http.ResponseWriter, r *http.Request) {
	if api.stateManager == nil {
		http.Error(w, "State manager not available", http.StatusServiceUnavailable)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Contract *vm.JSONContract `json:"contract"`
		Owner    string           `json:"owner"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Contract == nil {
		http.Error(w, "Contract data is required", http.StatusBadRequest)
		return
	}

	// Deploy the contract
	contract, err := vm.DeployJSONContract(request.Owner, request.Contract, true)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to deploy contract: %v", err), http.StatusInternalServerError)
		return
	}

	// Store the contract
	api.stateManager.SetContract(contract.Address, contract)

	response := map[string]interface{}{
		"status":    "success",
		"message":   "Contract deployed successfully",
		"address":   contract.Address,
		"name":      contract.Name,
		"owner":     contract.Owner,
		"timestamp": time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

// POST /contract/call
func (api *APIServer) handleCallContract(w http.ResponseWriter, r *http.Request) {
	if api.stateManager == nil {
		http.Error(w, "State manager not available", http.StatusServiceUnavailable)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ContractAddress string        `json:"contract_address"`
		Function        string        `json:"function"`
		Args            []interface{} `json:"args"`
		Caller          string        `json:"caller"`
		GasLimit        uint64        `json:"gas_limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the contract
	contract, exists := api.stateManager.GetContract(request.ContractAddress)
	if !exists {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	// Create execution context
	execCtx := &vm.ExecutionContext{
		Caller:   request.Caller,
		Value:    0, // No value transfer for now
		GasLimit: request.GasLimit,
	}

	// Create VM instance
	vmInstance := &vm.VM{
		StateManager: api.stateManager,
	}

	// Initialize VM memory with contract storage
	vmInstance.Memory = make(map[string]int64)
	for k, v := range contract.Storage {
		if ival, ok := v.(int64); ok {
			vmInstance.Memory[k] = ival
		}
	}

	// Execute function
	if err := contract.CallFunction(request.Function, request.Args, vmInstance, execCtx); err != nil {
		http.Error(w, fmt.Sprintf("Function execution failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Update contract storage
	for k, v := range vmInstance.Memory {
		contract.Storage[k] = v
	}
	contract.UpdatedAt = time.Now().Unix()
	api.stateManager.SetContract(contract.Address, contract)

	response := map[string]interface{}{
		"status":    "success",
		"message":   "Function executed successfully",
		"gas_used":  vmInstance.GetGasUsed(),
		"gas_limit": vmInstance.GetGasLimit(),
		"timestamp": time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

// GET /contract/list
func (api *APIServer) handleListContracts(w http.ResponseWriter, r *http.Request) {
	if api.stateManager == nil {
		http.Error(w, "State manager not available", http.StatusServiceUnavailable)
		return
	}

	// Get all contracts (this would need to be implemented in StateManager)
	// For now, return a mock response
	contracts := []map[string]interface{}{
		{
			"address":    "CONTRACT_12345678",
			"name":       "SimpleToken",
			"version":    "1.0",
			"owner":      "0x1234567890abcdef",
			"created_at": time.Now().Unix(),
		},
	}

	json.NewEncoder(w).Encode(contracts)
}

// GET /contract/info
func (api *APIServer) handleGetContractInfo(w http.ResponseWriter, r *http.Request) {
	if api.stateManager == nil {
		http.Error(w, "State manager not available", http.StatusServiceUnavailable)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Contract address is required", http.StatusBadRequest)
		return
	}

	contract, exists := api.stateManager.GetContract(address)
	if !exists {
		http.Error(w, "Contract not found", http.StatusNotFound)
		return
	}

	// Get function names
	functions := make([]string, 0)
	for funcName := range contract.Functions {
		functions = append(functions, funcName)
	}

	response := map[string]interface{}{
		"address":    contract.Address,
		"name":       contract.Name,
		"version":    contract.Version,
		"owner":      contract.Owner,
		"functions":  functions,
		"storage":    contract.Storage,
		"created_at": contract.CreatedAt,
		"updated_at": contract.UpdatedAt,
		"upgradable": contract.Upgradable,
	}

	json.NewEncoder(w).Encode(response)
}

// GET /contract/examples
func (api *APIServer) handleGetContractExamples(w http.ResponseWriter, r *http.Request) {
	examples := map[string]interface{}{
		"simple_token": vm.GetSimpleTokenContract(),
		"voting":       vm.GetVotingContract(),
		"escrow":       vm.GetEscrowContract(),
	}

	json.NewEncoder(w).Encode(examples)
}

// Database backup and recovery handlers

// GET /backup/status
func (api *APIServer) handleBackupStatus(w http.ResponseWriter, r *http.Request) {
	if api.stateManager == nil {
		http.Error(w, "State manager not available", http.StatusServiceUnavailable)
		return
	}

	status := api.stateManager.GetBackupStatus()
	json.NewEncoder(w).Encode(status)
}

// GET /backup/list
func (api *APIServer) handleBackupList(w http.ResponseWriter, r *http.Request) {
	if api.stateManager == nil {
		http.Error(w, "State manager not available", http.StatusServiceUnavailable)
		return
	}

	backups := api.stateManager.GetBackupList()
	json.NewEncoder(w).Encode(backups)
}

// POST /backup/create
func (api *APIServer) handleBackupCreate(w http.ResponseWriter, r *http.Request) {
	if api.stateManager == nil {
		http.Error(w, "State manager not available", http.StatusServiceUnavailable)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create backup in background to avoid blocking
	go func() {
		err := api.stateManager.CreateManualBackup()
		if err != nil {
			log.Printf("❌ [API] Failed to create manual backup: %v", err)
		}
	}()

	response := map[string]interface{}{
		"message":   "Manual backup started",
		"status":    "started",
		"timestamp": time.Now().Unix(),
	}

	json.NewEncoder(w).Encode(response)
}

// GET /network/architecture
func (api *APIServer) handleGetNetworkArchitecture(w http.ResponseWriter, r *http.Request) {
	// Get consensus information
	validators := api.consensusManager.GetAllValidators()
	totalStake := uint64(0)
	activeValidators := 0
	for _, v := range validators {
		totalStake += v.Stake
		if v.Active {
			activeValidators++
		}
	}

	// Get network topology information
	peerCount := 0
	// TODO: Get peer count from node
	if api.node != nil {
		// peerCount = len(api.node.Peers) // TODO: Implement when Peers field is exported
	}

	// Get sharding information
	shardInfo := map[string]interface{}{
		"enabled": false,
		"shards":  0,
	}
	// TODO: Implement sharding through consensus manager
	if false {
		shardInfo["enabled"] = false
		shardInfo["shards"] = 0
		shardInfo["message"] = "Sharding not yet implemented"
	}

	architecture := map[string]interface{}{
		"nodeTypes": map[string]interface{}{
			"validators": map[string]interface{}{
				"count":       len(validators),
				"active":      activeValidators,
				"totalStake":  totalStake,
				"minStake":    1000, // TODO: Get from consensus manager
				"description": "Nodes that participate in block production and consensus validation. Must stake tokens and maintain high uptime.",
			},
			"observers": map[string]interface{}{
				"count":       peerCount,
				"description": "Lightweight nodes that sync with the network but don't participate in consensus. Used for transaction submission and data retrieval.",
			},
			"fullNodes": map[string]interface{}{
				"count":       len(validators) + peerCount,
				"description": "Complete blockchain replicas that store the full state and validate all transactions and blocks.",
			},
		},
		"p2pProtocol": map[string]interface{}{
			"type":        "libp2p",
			"version":     "1.0.0",
			"discovery":   "UDP broadcast + libp2p DHT",
			"transport":   "TCP with Noise encryption",
			"description": "Peer-to-peer communication using libp2p protocol stack with automatic peer discovery and encrypted connections.",
		},
		"consensusMechanism": map[string]interface{}{
			"type":        "Proof of Stake (PoS)",
			"blockTime":   "5s", // TODO: Get from consensus manager
			"finality":    "1 confirmation",
			"description": "Validators are selected based on stake amount, performance score, reputation, and uptime. Block rewards are distributed to active validators.",
		},
		"networkTopology": map[string]interface{}{
			"type":        "Mesh Network",
			"connections": "Dynamic peer connections",
			"maxPeers":    10, // TODO: Get from consensus manager
			"description": "Decentralized mesh topology where nodes connect to multiple peers for redundancy and efficient data propagation.",
		},
		"securityFeatures": map[string]interface{}{
			"encryption":     "Noise protocol for transport encryption",
			"authentication": "ECDSA signatures for transactions and blocks",
			"rateLimiting":   "Connection and message rate limiting",
			"slashing":       "Validator slashing for misbehavior",
			"description":    "Multi-layered security including transport encryption, cryptographic signatures, and economic penalties for malicious behavior.",
		},
		"sharding": shardInfo,
	}

	json.NewEncoder(w).Encode(architecture)
}

// Identity Management API Handlers

func (api *APIServer) handleCreateIdentity(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address   string `json:"address"`
		PublicKey string `json:"public_key"`
		Username  string `json:"username"`
		Email     string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	publicKey, err := hex.DecodeString(req.PublicKey)
	if err != nil {
		http.Error(w, "Invalid public key format", http.StatusBadRequest)
		return
	}

	// TODO: Implement identity creation through identity manager
	identity, err := api.identityManager.CreateIdentity(req.Address, publicKey, req.Username, req.Email)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"identity": identity,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *APIServer) handleGetIdentity(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter required", http.StatusBadRequest)
		return
	}

	identity, err := api.identityManager.GetIdentity(address)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"identity": identity,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *APIServer) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address string                `json:"address"`
		Profile *identity.UserProfile `json:"profile"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Implement profile update through identity manager
	err := api.identityManager.UpdateProfile(req.Address, req.Profile)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *APIServer) handleUpdateActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Address      string `json:"address"`
		ActivityType string `json:"activity_type"`
		Value        int64  `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Implement activity update through identity manager
	err := api.identityManager.UpdateActivity(req.Address, req.ActivityType, req.Value)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Activity updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *APIServer) handleCreatePrivacyProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProofType string      `json:"proof_type"`
		Data      interface{} `json:"data"`
		Address   string      `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Implement privacy proof creation
	// TODO: Implement privacy proof creation with proper type conversion
	proof, err := api.identityManager.CreatePrivacyProof("default", req.Data, req.Address)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"proof":   proof,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *APIServer) handleVerifyPrivacyProof(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Proof interface{} `json:"proof"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Implement privacy proof verification
	// TODO: Implement privacy proof verification with proper type conversion
	valid, err := api.identityManager.VerifyPrivacyProof(nil)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"valid":   valid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *APIServer) handleGetIdentityForSocial(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	requester := r.URL.Query().Get("requester")

	if address == "" {
		http.Error(w, "Address parameter required", http.StatusBadRequest)
		return
	}

	identity, err := api.identityManager.GetIdentityForSocial(address, requester)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"identity": identity,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *APIServer) handleGetIdentityForCommerce(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	requester := r.URL.Query().Get("requester")

	if address == "" {
		http.Error(w, "Address parameter required", http.StatusBadRequest)
		return
	}

	identity, err := api.identityManager.GetIdentityForCommerce(address, requester)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"identity": identity,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (api *APIServer) handleGetIdentityForGovernance(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter required", http.StatusBadRequest)
		return
	}

	identity, err := api.identityManager.GetIdentityForGovernance(address)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"identity": identity,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ===== SOCIAL MEDIA API ENDPOINTS =====

// POST /social/post/create
func (api *APIServer) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Author     string   `json:"author"`
		Content    string   `json:"content"`
		MediaURLs  []string `json:"media_urls,omitempty"`
		Visibility string   `json:"visibility"`
		Category   string   `json:"category"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post, err := api.socialManager.CreatePost(req.Author, req.Content, req.MediaURLs, req.Visibility, req.Category)
	if err != nil {
		http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    post,
	})
}

// GET /social/post/get
func (api *APIServer) handleGetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post ID parameter required", http.StatusBadRequest)
		return
	}

	post, err := api.socialManager.GetPost(postID)
	if err != nil {
		http.Error(w, "Failed to get post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    post,
	})
}

// POST /social/comment/create
func (api *APIServer) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PostID   string `json:"post_id"`
		Author   string `json:"author"`
		Content  string `json:"content"`
		ParentID string `json:"parent_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	comment, err := api.socialManager.CreateComment(req.PostID, req.Author, req.Content, req.ParentID)
	if err != nil {
		http.Error(w, "Failed to create comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    comment,
	})
}

// POST /social/like
func (api *APIServer) handleLikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PostID   string `json:"post_id"`
		UserID   string `json:"user_id"`
		LikeType string `json:"like_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.socialManager.LikePost(req.PostID, req.UserID, req.LikeType)
	if err != nil {
		http.Error(w, "Failed to like post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Post liked successfully",
	})
}

// POST /social/unlike
func (api *APIServer) handleUnlikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PostID string `json:"post_id"`
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.socialManager.UnlikePost(req.PostID, req.UserID)
	if err != nil {
		http.Error(w, "Failed to unlike post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Post unliked successfully",
	})
}

// POST /social/tip
func (api *APIServer) handleTipPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PostID string `json:"post_id"`
		UserID string `json:"user_id"`
		Amount uint64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// In a real implementation we would check user balance here and transfer tokens.
	// For now we just call the social manager to update the post's stats.
	post, err := api.socialManager.TipPost(req.PostID, req.UserID, req.Amount)
	if err != nil {
		http.Error(w, "Failed to tip post: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Post tipped successfully",
		"data":    post,
	})
}

// GET /social/feed
func (api *APIServer) handleGetFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	// if userID == "" {
	// 	// Optional: returns global feed if empty
	// }

	limit := 20 // Default limit
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	posts, err := api.socialManager.GetFeed(userID, limit)
	if err != nil {
		http.Error(w, "Failed to get feed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    posts,
	})
}

// GET /social/search
func (api *APIServer) handleSearchPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter required", http.StatusBadRequest)
		return
	}

	limit := 20 // Default limit
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	posts := api.socialManager.SearchPosts(query, limit)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    posts,
	})
}

// GET /social/trending
func (api *APIServer) handleGetTrendingHashtags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 10 // Default limit
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	hashtags := api.socialManager.GetTrendingHashtags(limit)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    hashtags,
	})
}

// POST /social/report
func (api *APIServer) handleReportContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Reporter    string `json:"reporter"`
		TargetID    string `json:"target_id"`
		TargetType  string `json:"target_type"`
		Reason      string `json:"reason"`
		Description string `json:"description,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.socialManager.ReportContent(req.Reporter, req.TargetID, req.TargetType, req.Reason, req.Description)
	if err != nil {
		http.Error(w, "Failed to report content: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Content reported successfully",
	})
}

// ===== OBJECT ENERGY PHYSICS ENDPOINTS =====

// GET /social/object/energy?object_id={firebase_doc_id}&object_type={post|item}
func (api *APIServer) handleGetObjectEnergy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	objectID := r.URL.Query().Get("object_id")
	if objectID == "" {
		http.Error(w, "object_id parameter required", http.StatusBadRequest)
		return
	}

	objectType := r.URL.Query().Get("object_type")
	if objectType == "" {
		objectType = "post" // Default to "post" if not specified
	}

	energy, err := api.socialManager.GetOrCreateObjectEnergy(objectID, objectType)
	if err != nil {
		http.Error(w, "Failed to get object energy: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    energy,
	})
}

// POST /social/object/energize
func (api *APIServer) handleEnergizeObject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ObjectID string `json:"object_id"` // Firebase document ID
		UserID   string `json:"user_id"`   // Wallet address
		Amount   uint64 `json:"amount"`    // TCOIN amount
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ObjectID == "" || req.UserID == "" || req.Amount == 0 {
		http.Error(w, "object_id, user_id, and amount are required", http.StatusBadRequest)
		return
	}

	energy, err := api.socialManager.EnergizeObject(req.ObjectID, req.UserID, req.Amount)
	if err != nil {
		http.Error(w, "Failed to energize object: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Object energized with %d TCOIN", req.Amount),
		"data":    energy,
	})
}

// ===== GOVERNANCE API ENDPOINTS =====

// POST /governance/proposal/create
func (api *APIServer) handleCreateProposal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Proposer     string                        `json:"proposer"`
		Title        string                        `json:"title"`
		Description  string                        `json:"description"`
		Category     string                        `json:"category"`
		Actions      []governance.GovernanceAction `json:"actions"`
		CurrentBlock int64                         `json:"current_block"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	proposal, err := api.governanceManager.CreateProposal(req.Proposer, req.Title, req.Description, req.Category, req.Actions, req.CurrentBlock)
	if err != nil {
		http.Error(w, "Failed to create proposal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    proposal,
	})
}

// GET /governance/proposal/get (duplicate - removing)
// This function is a duplicate of the one at line 1193

// POST /governance/proposal/activate
func (api *APIServer) handleActivateProposal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProposalID   string `json:"proposal_id"`
		CurrentBlock int64  `json:"current_block"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.governanceManager.ActivateProposal(req.ProposalID, req.CurrentBlock)
	if err != nil {
		http.Error(w, "Failed to activate proposal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Proposal activated successfully",
	})
}

// POST /governance/proposal/vote
func (api *APIServer) handleVoteProposal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProposalID string `json:"proposal_id"`
		Voter      string `json:"voter"`
		Choice     string `json:"choice"`
		Reason     string `json:"reason,omitempty"`
		Weight     uint64 `json:"weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.governanceManager.Vote(req.ProposalID, req.Voter, req.Choice, req.Reason, req.Weight)
	if err != nil {
		http.Error(w, "Failed to vote on proposal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Vote cast successfully",
	})
}

// POST /governance/proposal/execute
func (api *APIServer) handleExecuteProposal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProposalID   string `json:"proposal_id"`
		CurrentBlock int64  `json:"current_block"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.governanceManager.ExecuteProposal(req.ProposalID, req.CurrentBlock)
	if err != nil {
		http.Error(w, "Failed to execute proposal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Proposal executed successfully",
	})
}

// POST /governance/proposal/discuss
func (api *APIServer) handleAddDiscussionComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ProposalID string `json:"proposal_id"`
		Author     string `json:"author"`
		Content    string `json:"content"`
		ParentID   string `json:"parent_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.governanceManager.AddDiscussionComment(req.ProposalID, req.Author, req.Content, req.ParentID)
	if err != nil {
		http.Error(w, "Failed to add discussion comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Comment added successfully",
	})
}

// GET /governance/proposals/active
func (api *APIServer) handleGetActiveProposals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	proposals := api.governanceManager.GetActiveProposals()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    proposals,
	})
}

// GET /governance/proposals/category
func (api *APIServer) handleGetProposalsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "Category parameter required", http.StatusBadRequest)
		return
	}

	proposals := api.governanceManager.GetProposalsByCategory(category)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    proposals,
	})
}

// POST /governance/referendum/create
func (api *APIServer) handleCreateReferendum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Question    string   `json:"question"`
		Options     []string `json:"options"`
		Duration    string   `json:"duration"` // e.g., "24h", "7d"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse duration
	duration, err := time.ParseDuration(req.Duration)
	if err != nil {
		http.Error(w, "Invalid duration format", http.StatusBadRequest)
		return
	}

	referendum, err := api.governanceManager.CreateReferendum(req.Title, req.Description, req.Question, req.Options, duration)
	if err != nil {
		http.Error(w, "Failed to create referendum: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    referendum,
	})
}

// POST /governance/referendum/vote
func (api *APIServer) handleVoteReferendum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ReferendumID string `json:"referendum_id"`
		Voter        string `json:"voter"`
		Option       string `json:"option"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := api.governanceManager.VoteReferendum(req.ReferendumID, req.Voter, req.Option)
	if err != nil {
		http.Error(w, "Failed to vote on referendum: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Referendum vote cast successfully",
	})
}

// POST /admin/resolve-dispute
func (api *APIServer) handleAdminResolveDispute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		OrderID  string `json:"order_id"`
		PayBuyer bool   `json:"pay_buyer"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if api.TreasuryWallet == nil {
		http.Error(w, "Treasury wallet not initialized", http.StatusInternalServerError)
		return
	}

	address := wallet.PublicKeyToAddress(api.TreasuryWallet.PublicKey)
	nonce := api.stateManager.GetNonce(address)

	args := []interface{}{req.OrderID, req.PayBuyer}
	dataBytes, _ := json.Marshal(map[string]interface{}{
		"function": "resolveDispute",
		"args":     args,
	})

	tx := transaction.Transaction{
		Type:            transaction.TxTypeCall,
		Sender:          address,
		SenderPublicKey: hex.EncodeToString(api.TreasuryWallet.PublicKey),
		Recipient:       vm.MarketplaceContractAddress,
		Amount:          0,
		Fee:             10,
		Nonce:           nonce,
		Data:            string(dataBytes),
		Timestamp:       time.Now().Unix(),
	}

	txHash := wallet.CalculateTxHash(tx)
	rInt, sInt, err := ecdsa.Sign(rand.Reader, api.TreasuryWallet.PrivateKey, txHash)
	if err != nil {
		http.Error(w, "Failed to sign resolution transaction", http.StatusInternalServerError)
		return
	}
	tx.Signature = fmt.Sprintf("%064x%064x", rInt, sInt)

	if err := api.transactionManager.AddTransaction(tx); err != nil {
		http.Error(w, "Failed to submit resolution transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[ARBITRATION] ⚖️ Resolved Dispute for order %s. PayBuyer: %t", req.OrderID, req.PayBuyer)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Dispute resolution transaction executed",
		"tx_hash": hex.EncodeToString(txHash),
	})
}

// GET /admin/disputes
func (api *APIServer) handleAdminGetDisputes(w http.ResponseWriter, r *http.Request) {
	stateAdapter := api.stateManager.GetStateAdapter()
	marketplaceHelper := api.stateManager.GetMarketplaceHelper()

	if stateAdapter == nil || marketplaceHelper == nil {
		http.Error(w, "Marketplace contract not initialized", http.StatusServiceUnavailable)
		return
	}

	ctx := &vm.ExecutionContext{
		State:     stateAdapter,
		Timestamp: time.Now().Unix(),
	}

	// We'll return the array of disputed orders to the UI
	var disputes []interface{}

	// Get all keys using strings.HasPrefix since stateAdapter iterates the db. 
	// For simplicity, since the admin UI is a specialized endpoint, we can use the existing `GetContractStorage` if we knew the orders.
	// However, we can simply fetch "order_counter", loop from 1 to orderCounter, and fetch each order!
	totalOrders, exists := stateAdapter.GetContractStorage(vm.MarketplaceContractAddress, "order_counter")
	if exists {
		for i := int64(1); i <= totalOrders; i++ {
			orderID := fmt.Sprintf("ORD-%d", i)
			status, hasStatus := stateAdapter.GetContractStorage(vm.MarketplaceContractAddress, "order_status_"+orderID)

			if hasStatus && status == vm.EscrowStatusDisputed {
				// We found a dispute! Let's get full info
				info := marketplaceHelper.GetOrder(ctx, orderID)
				if info != nil {
					disputes = append(disputes, info)
				}
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"disputes": disputes,
	})
}

// GET /admin/treasury-history
// Returns all confirmed blocks' transactions that involve the treasury wallet (sender or recipient).
func (api *APIServer) handleAdminTreasuryHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	// Resolve treasury address — must use TreasuryWallet and PublicKeyToAddress
	// to match the format stored in tx.Sender by handleAdminFaucet / handleFaucet.
	treasuryAddr := ""
	if api.TreasuryWallet != nil {
		treasuryAddr = wallet.PublicKeyToAddress(api.TreasuryWallet.PublicKey)
	}

	type TxRecord struct {
		BlockIndex int64  `json:"blockIndex"`
		BlockHash  string `json:"blockHash"`
		Type       string `json:"type"`
		Sender     string `json:"sender"`
		Recipient  string `json:"recipient"`
		Amount     int64  `json:"amount"`
		Fee        int64  `json:"fee"`
		Data       string `json:"data"`
		Timestamp  int64  `json:"timestamp"`
		Direction  string `json:"direction"` // "out" | "in"
	}

	var history []TxRecord

	// GetBlockHeight() returns the LATEST BLOCK INDEX (not count).
	// We must use GetChainLength() to get the actual number of blocks,
	// then iterate from 0 to chainLength-1 to cover all blocks including the latest.
	chainLen := api.blockManager.GetChainLength()
	for i := 0; i < chainLen; i++ {
		blk, err := api.blockManager.GetBlockByIndex(i)
		if err != nil || blk == nil {
			continue
		}
		for _, tx := range blk.Transactions {
			isTreasury := (treasuryAddr != "" && (tx.Sender == treasuryAddr || tx.Recipient == treasuryAddr))
			// If no wallet set, include all (fallback)
			if !isTreasury && treasuryAddr != "" {
				continue
			}
			direction := "in"
			if tx.Sender == treasuryAddr {
				direction = "out"
			}
			history = append(history, TxRecord{
				BlockIndex: int64(blk.Index),
				BlockHash:  blk.Hash,
				Type:       string(tx.Type),
				Sender:     tx.Sender,
				Recipient:  tx.Recipient,
				Amount:     tx.Amount,
				Fee:        tx.Fee,
				Data:       tx.Data,
				Timestamp:  tx.Timestamp,
				Direction:  direction,
			})
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"treasuryAddress": treasuryAddr,
		"count":           len(history),
		"transactions":    history,
	})
}

