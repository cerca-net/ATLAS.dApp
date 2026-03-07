package main

import (
	"atlas-blockchain/internal/api"
	"atlas-blockchain/internal/blockchain"
	"atlas-blockchain/internal/defi"
	"atlas-blockchain/internal/governance"
	"atlas-blockchain/internal/identity"
	"atlas-blockchain/internal/social"
	"atlas-blockchain/pkg/block"
	"atlas-blockchain/pkg/config"
	"atlas-blockchain/pkg/network"
	"atlas-blockchain/pkg/sharding"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/wallet"
	"bytes"
	"context"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

var (
	blockchainConfig   *config.BlockchainConfig
	node               *network.Node
	transactionManager *blockchain.TransactionManager
	blockManager       *blockchain.BlockManager
	consensusManager   *blockchain.ConsensusManager
	stateManager       *blockchain.StateManager
	chainSyncManager   *blockchain.ChainSyncManager
	fastSyncManager    *blockchain.FastSyncManager // Fast sync manager
	peerManager        *network.PeerManager        // Peer manager
	identityManager    *identity.IdentityManager
	defiManager        *defi.DeFiManager
	socialManager      *social.SocialManager
	governanceManager  *governance.GovernanceManager
	shardManager       *sharding.ShardManager
	validatorMode      *bool
	p2pNode            *network.P2PNode // Add P2P node
	isTestMode         bool             // Flag to indicate if we're running in test mode
)

// SetTestMode sets the test mode flag
func SetTestMode(enabled bool) {
	isTestMode = enabled
}

func main() {
	// Parse command line flags early
	port := flag.Int("port", 8000, "Port for peer discovery")
	maxPeers := flag.Int("peers", 10, "Maximum number of peers")
	apiPort := flag.Int("api", 8080, "Port for API server")
	validatorMode = flag.Bool("validator", true, "Run node as validator (with wallet)")
	keyPath := flag.String("key", "nodekey.priv", "Path to private key file for libp2p identity")
	legacyNetworking := flag.Bool("legacy-net", false, "Enable legacy TCP networking") // NEW FLAG
	testMode := flag.Bool("test", false, "Run in test mode (disable infinite loops)")
	dataDir := flag.String("datadir", ".", "Directory to store blockchain data (db, snapshots, backups)")
	validatorKeyPath := flag.String("validator-key", "", "Path to validator private key file (hex). If empty, generates ephemeral wallet.")
	flag.Parse()

	// Set test mode flag
	isTestMode = *testMode

	blockchainConfig = config.DefaultConfig()
	blockchainConfig.PeerDiscoveryPort = *port
	blockchainConfig.MaxPeers = *maxPeers
	blockchainConfig.BlockTime = 30 * time.Second // Longer block time for grouping transactions and cleaner logs
	blockchainConfig.DataDir = *dataDir

	// Ensure data directory exists
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	stateManager = blockchain.NewStateManager(blockchainConfig)

	// Migrate existing JSON snapshots to database if available
	if err := stateManager.MigrateToDatabase(); err != nil {
		log.Printf("⚠️  Database migration failed: %v", err)
	}

	transactionManager = blockchain.NewTransactionManager(blockchainConfig, stateManager)
	blockManager = blockchain.NewBlockManager(blockchainConfig, stateManager)
	consensusManager = blockchain.NewConsensusManager(blockchainConfig, blockManager)
	stateManager.SetConsensusManager(consensusManager)

	// Initialize identity manager for social-commerce-governance platform
	identityManager = identity.NewIdentityManager()

	// Initialize DeFi manager for lending, trading, staking, and governance
	defiManager = defi.NewDeFiManager(identityManager)

	// Initialize social media manager for posts, comments, and content moderation
	socialManager = social.NewSocialManager(identityManager, stateManager.GetDatabase(), stateManager)

	// Initialize governance manager for proposals, voting, and community governance
	governanceManager = governance.NewGovernanceManager(socialManager, defiManager, identityManager)

	// Initialize shard manager
	shardConfig := &sharding.ShardConfig{
		TotalShards:     4,
		ShardSize:       10, // Validators per shard
		CrossShardDelay: 5 * time.Second,
		ConsensusType:   "pbft",
	}
	shardManager = sharding.NewShardManager(shardConfig)

	ctx := context.Background()

	// Use key path relative to datadir if not absolute so each node gets a unique key
	p2pKeyPath := *keyPath
	if !filepath.IsAbs(p2pKeyPath) && *dataDir != "." {
		p2pKeyPath = filepath.Join(*dataDir, p2pKeyPath)
	}

	p2pNode, err := network.NewP2PNode(ctx, blockchainConfig.PeerDiscoveryPort, p2pKeyPath)
	if err != nil {
		log.Fatalf("Failed to start P2P node: %v", err)
	}

	// Initialize fast sync manager
	fastSyncManager = blockchain.NewFastSyncManager(stateManager, blockManager, *dataDir)
	log.Printf("✅ Fast sync manager initialized")

	// Initialize peer manager
	peerManager = network.NewPeerManager(p2pNode, *dataDir)
	log.Printf("✅ Peer manager initialized with persistence")

	// Initialize chain synchronization manager
	chainSyncManager = blockchain.NewChainSyncManager(blockManager, stateManager, p2pNode, blockchainConfig)
	chainSyncManager.SetCallbacks(
		func() { log.Printf("🔄 [SYNC] Chain synchronization started") },
		func(blocksSynced, totalBlocks int64) {
			log.Printf("📊 [SYNC] Progress: %d/%d blocks", blocksSynced, totalBlocks)
		},
		func() { log.Printf("✅ [SYNC] Chain synchronization completed") },
		func(err error) { log.Printf("❌ [SYNC] Chain synchronization failed: %v", err) },
		func(forkHeight int64, canonical, forked []string) {
			log.Printf("🔄 [SYNC] Fork detected at height %d", forkHeight)
		},
	)
	if err != nil {
		log.Fatalf("Failed to start P2P node: %v", err)
	}

	// Enhanced debug: log all multiaddresses and peerstore contents
	log.Printf("[DEBUG] Our Peer ID: %s", p2pNode.Host.ID().String())
	for _, addr := range p2pNode.Host.Addrs() {
		log.Printf("[DEBUG] Listening on multiaddress: %s", addr.String())
	}
	log.Printf("[DEBUG] Peerstore peers at startup: %v", p2pNode.Host.Peerstore().Peers())

	// Improved: Always attempt manual connection if NODE1_MULTIADDR is set
	addrStr := os.Getenv("NODE1_MULTIADDR")
	if addrStr != "" {
		log.Printf("[P2P] Attempting manual connection to: %s", addrStr)
		info, err := peer.AddrInfoFromString(addrStr)
		if err == nil {
			log.Printf("[DEBUG] Peerstore before connect: %v", p2pNode.Host.Peerstore().Peers())
			if err := p2pNode.Host.Connect(ctx, *info); err != nil {
				log.Printf("[P2P] Failed to connect to peer: %v (%T)", err, err)
			} else {
				log.Printf("[P2P] Successfully connected to peer: %s", info.ID.String())
			}
			log.Printf("[DEBUG] Peerstore after connect: %v", p2pNode.Host.Peerstore().Peers())
		} else {
			log.Printf("[P2P] Invalid multiaddress: %v", err)
		}
	}

	// Assign block and transaction processing callbacks
	p2pNode.OnBlockReceived = func(blockMsg network.BlockMessage) {
		blk := &blockMsg.Block
		log.Printf("[P2P] Received block %d from network", blk.Index)

		// Validate and add block
		if err := blockManager.AddBlock(blk); err != nil {
			log.Printf("[P2P] Failed to add received block %d: %v", blk.Index, err)
		} else {
			log.Printf("[P2P] Successfully added block %d to local chain", blk.Index)
		}
	}
	p2pNode.RegisterStreamHandler()

	// Broadcast new blocks to peers when added
	blockManager.SetOnBlockAddedCallback(func(blk *block.Block) {
		if p2pNode != nil {
			if err := p2pNode.BroadcastBlock(context.Background(), blk); err != nil {
				log.Printf("[P2P] Failed to broadcast block %d: %v", blk.Index, err)
			}
		}
	})

	// Example: Broadcast new transactions (if you have an event/callback for new transactions)
	// transactionManager.OnNewTransaction = func(tx transaction.Transaction) {
	// 	data, err := json.Marshal(tx)
	// 	if err != nil {
	// 		log.Printf("[P2P] Failed to marshal transaction for broadcast: %v", err)
	// 		return
	// 	}
	// 	msg := network.TransactionMessage{TxData: data}
	// 	netMsg := network.NetworkMessage{Type: network.MsgTypeTransaction, Payload: data}
	// 	for _, peer := range p2pNode.Host.Peerstore().Peers() {
	// 		if peer == p2pNode.Host.ID() {
	// 			continue // Don't send to self
	// 		}
	// 		if err := p2pNode.SendMessage(context.Background(), peer, netMsg); err != nil {
	// 			log.Printf("[P2P] Failed to broadcast transaction to peer %s: %v", peer.String(), err)
	// 		}
	// 	}
	// }

	// DebugTxFlow()
	// panic("[DEBUG] DebugTxFlow complete - halting execution for inspection.")
	// Initialize configuration
	if err := blockchainConfig.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create and initialize node
	node = network.NewNode(fmt.Sprintf("localhost:%d", *port), "testtoken")
	if *validatorMode {
		if err := initializeNode(*validatorKeyPath); err != nil {
			log.Fatalf("Failed to initialize node: %v", err)
		}
	} else {
		node.Wallet = nil
		node.ValidatorAddress = ""
		// Node runs as observer/relay
	}

	// Set up wallet synchronization callback
	blockManager.SetOnBlockAddedCallback(func(block *block.Block) {
		// Synchronize wallet balance with state manager after each block
		if node.Wallet != nil {
			stateManager.SyncWalletBalance(node.Wallet)
		}
	})

	// After successful registration as validator (in initializeNode or main)
	if *validatorMode {
		// Broadcast validator registration to peers
		regMsg := network.ValidatorRegistrationMessage{
			Address: node.ValidatorAddress,
			Stake:   uint64(blockchainConfig.MinStake * 2), // or actual stake
			PeerID:  p2pNode.Host.ID().String(),
		}
		payload, _ := json.Marshal(regMsg)
		netMsg := network.NetworkMessage{Type: network.MsgTypeValidatorRegistration, Payload: payload}
		for _, peer := range p2pNode.Host.Peerstore().Peers() {
			if peer == p2pNode.Host.ID() {
				continue
			}
			log.Printf("[P2P] Broadcasting validator registration to peer: %s", peer.String())
			_ = p2pNode.SendMessage(context.Background(), peer, netMsg)
		}
	}

	// Handle incoming validator registration messages
	p2pNode.OnValidatorRegistrationReceived = func(reg network.ValidatorRegistrationMessage) {
		log.Printf("[P2P] Received validator registration for address: %s", reg.Address)
		if consensusManager != nil {
			// Check if we already have this validator
			if _, err := consensusManager.GetValidatorInfo(reg.Address); err != nil {
				consensusManager.AddExternalValidator(reg.Address, reg.Stake)
				log.Printf("[P2P] Registered new external validator: %s", reg.Address)

				// Rebroadcast to other peers to ensure network-wide consensus consistency
				go func() {
					payload, _ := json.Marshal(reg)
					netMsg := network.NetworkMessage{Type: network.MsgTypeValidatorRegistration, Payload: payload}
					p2pNode.BroadcastMessage(context.Background(), netMsg)
				}()
			}

			// Track validator peer if PeerID is available in message
			if reg.PeerID != "" {
				if pid, err := peer.Decode(reg.PeerID); err == nil && pid != p2pNode.Host.ID() {
					peerManager.AddPeer(pid, true, reg.Stake)
				}
			}
		}
	}

	// Handle incoming transactions from peers
	p2pNode.OnTransactionReceived = func(txMsg network.TransactionMessage) {
		tx := txMsg.Transaction
		txHash := hex.EncodeToString(wallet.CalculateTxHash(tx))

		log.Printf("[P2P] Received transaction from peer: %s", txHash)

		// Check if we already have this transaction (prevent duplicates)
		if transactionManager.HasTransaction(txHash) {
			log.Printf("[P2P] Ignoring duplicate transaction: %s", txHash)
			return
		}

		// Add to local mempool
		if err := transactionManager.AddTransaction(tx); err != nil {
			log.Printf("[P2P] Failed to add transaction to pool: %v", err)
			return
		}

		log.Printf("[P2P] Successfully added transaction %s to mempool", txHash)
	}

	// Start API server
	apiServer := api.NewAPIServer(blockManager, transactionManager, stateManager, consensusManager, fastSyncManager, peerManager, node, p2pNode, identityManager, socialManager, governanceManager, shardManager)
	apiServer.OnSyncRequest = func() {
		go startChainSync()
	}

	// TODO: Set up monitoring integration callbacks when monitor field is exported
	// Set up monitoring integration callbacks
	// if apiServer.monitor != nil {
	// 	apiServer.monitor.SetIntegrationCallbacks(
	// 		// Get transaction count
	// 		func() int {
	// 			return transactionManager.GetPoolSize()
	// 		},
	// 		// Get block height
	// 		func() int64 {
	// 			return int64(blockManager.GetBlockHeight())
	// 		},
	// 		// Get pending transactions
	// 		func() int {
	// 			return transactionManager.GetPoolSize()
	// 		},
	// 		// Get validator count
	// 		func() int {
	// 			return len(consensusManager.GetAllValidators())
	// 		},
	// 		// Get active peers
	// 		func() int {
	// 			return len(p2pNode.Host.Peerstore().Peers())
	// 		},
	// 		// Get total staked
	// 		func() float64 {
	// 			total := 0.0
	// 			for _, validator := range consensusManager.GetAllValidators() {
	// 				total += float64(validator.Stake)
	// 			}
	// 			return total
	// 		},
	// 		// Get last block hash
	// 		func() string {
	// 			if height := blockManager.GetBlockHeight(); height > 0 {
	// 				if lastBlock, err := blockManager.GetBlockByIndex(height - 1); err == nil {
	// 					return lastBlock.Hash
	// 				}
	// 			}
	// 			return ""
	// 		},
	// 		// Get contract count
	// 		func() int {
	// 			// This would integrate with the VM system
	// 			return 0 // Placeholder for now
	// 		},
	// 		// Get sync status
	// 		func() string {
	// 			if chainSyncManager != nil {
	// 				return string(chainSyncManager.GetStatus())
	// 			}
	// 			return "unknown"
	// 		},
	// 	)
	// }

	go apiServer.Start(fmt.Sprintf(":%d", *apiPort))

	// Periodically announce validator status to ensure new peers know about us
	if *validatorMode {
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					regMsg := network.ValidatorRegistrationMessage{
						Address: node.ValidatorAddress,
						Stake:   uint64(blockchainConfig.MinStake * 2),
						PeerID:  p2pNode.Host.ID().String(),
					}
					payload, _ := json.Marshal(regMsg)
					netMsg := network.NetworkMessage{Type: network.MsgTypeValidatorRegistration, Payload: payload}
					p2pNode.BroadcastMessage(context.Background(), netMsg)
				}
			}
		}()
	}

	// Start blockchain operations
	go startBlockchain(*legacyNetworking)

	// Start backup system
	if stateManager != nil {
		stateManager.StartBackupSystem()
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	shutdown()
}

func initializeNode(keyPath string) error {
	var walletObj *wallet.Wallet
	var err error

	if keyPath != "" {
		// Resolve path relative to DataDir if needed
		fullPath := keyPath
		if !filepath.IsAbs(keyPath) && blockchainConfig.DataDir != "." {
			fullPath = filepath.Join(blockchainConfig.DataDir, keyPath)
		}

		// Try to load existing key
		if _, err := os.Stat(fullPath); err == nil {
			log.Printf("🔑 Loading validator key from %s", fullPath)
			keyBytes, err := ioutil.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("failed to read validator key file: %v", err)
			}
			walletObj, err = wallet.ImportWallet(string(bytes.TrimSpace(keyBytes)))
			if err != nil {
				return fmt.Errorf("failed to import wallet from key: %v", err)
			}
		} else {
			// Generate new key and save it
			log.Printf("🆕 Generating new persistent validator key at %s", fullPath)
			walletObj, err = wallet.NewWallet()
			if err != nil {
				return fmt.Errorf("failed to create wallet: %v", err)
			}
			// Save private key hex
			privKeyBytes, err := x509.MarshalECPrivateKey(walletObj.PrivateKey)
			if err != nil {
				return fmt.Errorf("failed to marshal private key: %v", err)
			}
			hexKey := hex.EncodeToString(privKeyBytes)
			if err := ioutil.WriteFile(fullPath, []byte(hexKey), 0600); err != nil {
				return fmt.Errorf("failed to save validator key: %v", err)
			}
		}
	} else {
		log.Printf("⚠️ No validator key path provided, using ephemeral wallet")
		walletObj, err = wallet.NewWallet()
		if err != nil {
			return fmt.Errorf("failed to create wallet: %v", err)
		}
	}

	// Add initial balance to the node's wallet
	walletObj.SetBalance(int64(blockchainConfig.MinStake + 100)) // MinStake + buffer for fees

	// Register node as validator using wallet
	kyc := blockchain.KYCInfo{
		FullName: "Node Validator",
		Country:  "System",
		IDNumber: "NODE001",
		Verified: true,
	}
	if err := consensusManager.RegisterValidator(walletObj, uint64(blockchainConfig.MinStake), kyc); err != nil {
		return fmt.Errorf("failed to register as validator: %v", err)
	}

	// Store wallet in node for future use
	node.Wallet = walletObj
	node.ValidatorAddress = wallet.PublicKeyToAddress(walletObj.PublicKey)

	return nil
}

func startBlockchain(legacyNetworking bool) {
	// Don't start infinite loops if in test mode
	if isTestMode {
		log.Printf("🧪 Test mode enabled - skipping infinite loops")
		return
	}

	// Start block production
	go produceBlocks()

	// Start peer discovery (legacy TCP networking) ONLY if enabled
	if legacyNetworking {
		go discoverPeers()
	}

	// Start transaction processing
	go processTransactions()

	// Start transaction pool maintenance
	go maintainTransactionPool()

	// Start validator monitoring
	go monitorValidators()

	// Start transaction pool synchronization
	node.StartTransactionPoolSync()

	// Start chain synchronization
	go startChainSync()
}

func produceBlocks() {
	// Align to the next block interval to synchronize with other nodes
	blockInterval := blockchainConfig.BlockTime
	now := time.Now()
	nextBlockTime := now.Truncate(blockInterval).Add(blockInterval)
	log.Printf("⏳ Aligning block production to %s...", nextBlockTime.Format("15:04:05"))
	time.Sleep(time.Until(nextBlockTime))

	log.Printf("🚀 Starting block production loop with interval: %v", blockchainConfig.BlockTime)
	ticker := time.NewTicker(blockchainConfig.BlockTime)
	defer ticker.Stop()

	blockCount := 0
	for {
		select {
		case <-ticker.C:
			// Check node state
			currentState := node.GetState()
			if currentState != api.NodeStateRunning {
				log.Printf("⏳ Block production skipped - node state: %s", currentState)
				continue
			}

			blockCount++
			log.Printf("⏰ Block production tick #%d - Current time: %s", blockCount, time.Now().Format("15:04:05"))

			// ACTIVATOR LOGIC: Only produce blocks if there are transactions
			poolSize := transactionManager.GetPoolSize()
			if poolSize == 0 {
				log.Printf("⏳ Activator Mode: Waiting for transactions (Pool: 0)...")
				continue
			}

			// Check if we are the chosen validator
			log.Printf("🔍 Attempting to choose validator...")

			lastBlock := blockManager.GetLatestBlock()
			var lastHash string
			var nextHeight int64

			if lastBlock != nil {
				lastHash = lastBlock.Hash
				nextHeight = int64(lastBlock.Index) + 1
			} else {
				// Fallback for first block if genesis not ready (should not happen in normal flow)
				lastHash = "genesis_seed"
				nextHeight = 0
			}

			validator, err := consensusManager.ChooseValidator(lastHash, nextHeight)
			if err != nil {
				log.Printf("❌ Failed to choose validator: %v", err)
				continue
			}
			log.Printf("✅ Chosen validator: %s (Our address: %s)", validator.Address, node.ValidatorAddress)

			// DEBUG: Force allow if we are the only validator or if strict mode is disabled
			isLocalValidator := validator.Address == node.ValidatorAddress

			// If we are the validator (or forced), forge a new block
			if isLocalValidator {
				log.Printf("🎯 We are the chosen validator! Attempting to forge block with %d transactions...", poolSize)

				if err := forgeAndBroadcastBlock(); err != nil {
					log.Printf("❌ Failed to forge block: %v", err)
					// Slash validator for failed block production
					if err := consensusManager.SlashValidator(node.ValidatorAddress, "failed block production"); err != nil {
						log.Printf("❌ Failed to slash validator: %v", err)
					}
				} else {
					log.Printf("✅ Block forged!")
				}
			} else {
				log.Printf("⏳ Not our turn. Chosen: %s, Us: %s", validator.Address[:8], node.ValidatorAddress[:8])
			}

			// log.Printf("🔄 Block production tick #%d completed...", blockCount) // Reduce noise
		}
	}
}

func monitorValidators() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check node state
			if node.GetState() == api.NodeStateStopped {
				continue
			}

			// Get all validators
			validators := consensusManager.GetAllValidators()

			// Log validator status
			for _, v := range validators {
				log.Printf("Validator %s: Stake=%d, Active=%v, SlashCount=%d",
					v.Address, v.Stake, v.Active, v.SlashCount)
			}
		}
	}
}

func forgeAndBroadcastBlock() error {
	log.Printf("🔨 Starting block forging process...")
	// Enable block forging using consensus logic
	validatorWallet := node.Wallet
	lastBlock := blockManager.GetLatestBlock()
	newBlock, err := blockchain.ForgeBlock(validatorWallet, lastBlock, stateManager, transactionManager, consensusManager)
	if err != nil {
		log.Printf("❌ forgeBlock failed: %v", err)
		return fmt.Errorf("failed to forge block: %v", err)
	}
	log.Printf("✅ Block forged successfully: Index=%d, Hash=%s", newBlock.Index, newBlock.Hash[:16]+"...")

	// Add the block to our own chain first
	if err := blockManager.AddBlock(newBlock); err != nil {
		log.Printf("❌ Failed to add forged block to local chain: %v", err)
		return fmt.Errorf("failed to add block: %v", err)
	}

	// Broadcast the new block to peers (if P2P is enabled)
	if p2pNode != nil {
		if err := p2pNode.BroadcastBlock(context.Background(), newBlock); err != nil {
			log.Printf("[P2P] Failed to broadcast forged block: %v", err)
		} else {
			log.Printf("📡 Broadcasted new block %d to peers.", newBlock.Index)
		}
	}

	// Reward the validator in metrics
	consensusManager.RewardValidator(node.ValidatorAddress, blockchain.BLOCK_REWARD)
	return nil
}

func processTransactions() {
	for {
		// Check node state
		if node.GetState() != api.NodeStateRunning {
			time.Sleep(time.Second)
			continue
		}

		// TODO: Process transactions when node fields are exported
		time.Sleep(time.Second)
	}
}

func maintainTransactionPool() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Remove expired transactions
			transactionManager.RemoveExpiredTransactions()

			// Log pool status
			log.Printf("Transaction pool size: %d", transactionManager.GetPoolSize())
		}
	}
}

func discoverPeers() {
	// TODO: Start network services when startNetwork method is exported
	// Start network services
	// if err := node.startNetwork(); err != nil {
	// 	log.Printf("Failed to start network: %v", err)
	// }
}

func shutdown() {
	log.Println("Shutting down blockchain node...")
	// Save final state
	log.Printf("Final blockchain length: %d", blockManager.GetChainLength())
	log.Printf("Final transaction pool size: %d", transactionManager.GetPoolSize())

	// Stop backup system
	if stateManager != nil {
		stateManager.StopBackupSystem()
		log.Println("Backup system stopped")
	}

	// Close database connection
	if err := stateManager.CloseDatabase(); err != nil {
		log.Printf("Failed to close database: %v", err)
	} else {
		log.Println("Database connection closed")
	}

	// Persist state (fallback to JSON)
	err := stateManager.ExportState("final_state.json")
	if err != nil {
		log.Printf("Failed to export state: %v", err)
	} else {
		log.Println("State exported to final_state.json")
	}
}

// startChainSync initiates chain synchronization
func startChainSync() {
	log.Printf("🔄 Starting chain synchronization...")

	// Wait a bit for the network to stabilize
	time.Sleep(5 * time.Second)

	ctx := context.Background()
	err := chainSyncManager.StartSync(ctx)
	if err != nil {
		log.Printf("❌ Failed to start chain synchronization: %v", err)
		return
	}

	// Monitor sync status
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			status := chainSyncManager.GetStatus()
			blocksSynced, totalBlocks := chainSyncManager.GetSyncProgress()
			duration := chainSyncManager.GetSyncDuration()

			log.Printf("📊 [SYNC] Status: %s, Progress: %d/%d, Duration: %v",
				status, blocksSynced, totalBlocks, duration)

			// TODO: Check sync status when constants are defined
			if status == "complete" || status == "failed" {
				log.Printf("🔄 [SYNC] Chain synchronization finished with status: %s", status)
				return
			}
		}
	}
}

func DebugTxFlow() {
	fmt.Println("[DEBUG] DebugTxFlow started")
	fmt.Println("[DEBUG] Creating wallet A...")
	walletA, err := wallet.NewWallet()
	if err != nil {
		fmt.Printf("Failed to create wallet A: %v\n", err)
		return
	}
	fmt.Println("[DEBUG] Creating wallet B...")
	walletB, err := wallet.NewWallet()
	if err != nil {
		fmt.Printf("Failed to create wallet B: %v\n", err)
		return
	}
	fmt.Println("[DEBUG] Setting initial balance for wallet A...")
	stateManager.SetBalance(walletA.PublicKeyStr(), 1000)
	// Use relative path for debug output
	absPath := "debug_output.txt"
	fmt.Printf("[DEBUG] Attempting to open %s...\n", absPath)
	debugFile, errDebug := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errDebug != nil {
		fmt.Fprintf(os.Stderr, "[DEBUG ERROR] Could not open debug_output.txt: %v\n", errDebug)
	} else {
		fmt.Println("[DEBUG] Writing initial balances to debug_output.txt...")
		fmt.Fprintf(debugFile, "[DEBUG] Initial balance A: %d\n", stateManager.GetBalance(walletA.PublicKeyStr()))
		fmt.Fprintf(debugFile, "[DEBUG] Initial balance B: %d\n", stateManager.GetBalance(walletB.PublicKeyStr()))
	}
	fmt.Println("[DEBUG] Creating transaction...")
	tx := transaction.Transaction{
		Sender:    walletA.PublicKeyStr(),
		Recipient: walletB.PublicKeyStr(),
		Amount:    100,
		Fee:       1,
		Timestamp: 0,
		Nonce:     0,
		Data:      "test transfer",
	}
	fmt.Println("[DEBUG] Signing transaction...")
	err = walletA.SignTransaction(&tx)
	if err != nil {
		if errDebug == nil {
			debugFile.Close()
		}
		fmt.Printf("Failed to sign transaction: %v\n", err)
		return
	}
	fmt.Println("[DEBUG] Adding transaction to manager...")
	err = transactionManager.AddTransaction(tx)
	fmt.Println("[DEBUG] Finished AddTransaction call.")
	if err != nil {
		if errDebug == nil {
			debugFile.Close()
		}
		fmt.Printf("Failed to add transaction: %v\n", err)
		return
	}
	fmt.Println("[DEBUG] Creating block with transaction...")
	transactions := transactionManager.GetTransactionsForBlock()
	lastBlock := blockManager.GetLatestBlock()
	blockObj, err := block.CreateNewBlock(transactions, lastBlock, walletA)
	if err != nil {
		if errDebug == nil {
			debugFile.Close()
		}
		fmt.Printf("Failed to create block: %v\n", err)
		return
	}
	if errDebug == nil {
		fmt.Println("[DEBUG] Writing block validator to debug_output.txt...")
		fmt.Fprintf(debugFile, "[DEBUG] Block validator: %s\n", blockObj.Validator)
	}
	fmt.Println("[DEBUG] Adding block to chain...")
	err = blockManager.AddBlock(blockObj)
	if err != nil {
		if errDebug == nil {
			debugFile.Close()
		}
		fmt.Printf("Failed to add block: %v\n", err)
		return
	}
	if errDebug == nil {
		fmt.Println("[DEBUG] Writing post-tx balances to debug_output.txt...")
		fmt.Fprintf(debugFile, "[DEBUG] Post-tx balance A: %d\n", stateManager.GetBalance(walletA.PublicKeyStr()))
		fmt.Fprintf(debugFile, "[DEBUG] Post-tx balance B: %d\n", stateManager.GetBalance(walletB.PublicKeyStr()))
		debugFile.Close()
	}
	fmt.Println("[DEBUG] Printing balances to stdout...")
	fmt.Printf("[DEBUG] Initial balance A: %d\n", stateManager.GetBalance(walletA.PublicKeyStr()))
	fmt.Printf("[DEBUG] Initial balance B: %d\n", stateManager.GetBalance(walletB.PublicKeyStr()))
	fmt.Printf("[DEBUG] Block validator: %s\n", blockObj.Validator)
	fmt.Printf("[DEBUG] Post-tx balance A: %d\n", stateManager.GetBalance(walletA.PublicKeyStr()))
	fmt.Printf("[DEBUG] Post-tx balance B: %d\n", stateManager.GetBalance(walletB.PublicKeyStr()))
	fmt.Println("[DEBUG] DebugTxFlow finished")
}
