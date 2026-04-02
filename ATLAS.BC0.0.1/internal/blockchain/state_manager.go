package blockchain

import (
	"atlas-blockchain/pkg/block"
	"atlas-blockchain/pkg/config"
	"atlas-blockchain/pkg/database"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/vm"
	"atlas-blockchain/pkg/wallet"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// StateSnapshot represents a point-in-time snapshot of the blockchain state
type StateSnapshot struct {
	Balances    map[string]int64 `json:"balances"`
	BlockHeight int64            `json:"block_height"`
	Timestamp   time.Time        `json:"timestamp"`
	Checksum    string           `json:"checksum"`
}

// Governance proposal states
const (
	ProposalPending   = "pending"
	ProposalActive    = "active"
	ProposalSucceeded = "succeeded"
	ProposalFailed    = "failed"
	ProposalExecuted  = "executed"
)

type Proposal struct {
	ID           string
	Proposer     string
	Description  string
	Actions      string // For now, a string; could be JSON or structured
	State        string
	VotesFor     int64
	VotesAgainst int64
	StartBlock   int64
	EndBlock     int64
	Voters       map[string]bool // Track who has voted
}

type Vote struct {
	ProposalID string
	Voter      string
	Choice     string // "for" or "against"
	Weight     int64  // Voting power
}

type OracleData struct {
	Key       string
	Value     string
	Timestamp int64
	Source    string
}

// StateManager handles the blockchain state with persistence
type StateManager struct {
	balances         map[string]int64
	mu               sync.RWMutex
	config           *config.BlockchainConfig
	snapshotPath     string
	lastSnapshot     time.Time
	snapshotInterval time.Duration
	checksums        map[string]string // Maps block hash to state checksum
	accounts         map[string]*database.Account

	// Database for persistent storage
	db *database.Database

	// Backup and recovery managers
	backupManager   *database.BackupManager
	recoveryManager *database.RecoveryManager

	// Contract registry: address -> contract
	contracts map[string]*vm.Contract

	// Token contract infrastructure
	stateAdapter      *StateAdapterImpl
	tokenHelper       *vm.TokenContractHelper
	stakingHelper     *vm.StakingContractHelper
	marketplaceHelper *vm.MarketplaceContractHelper
	governanceHelper  *vm.GovernanceContractHelper

	// Governance registries
	proposals map[string]*Proposal
	votes     map[string][]*Vote // proposalID -> votes

	// Oracle data registry: key -> OracleData
	oracleData       map[string]OracleData
	consensusManager *ConsensusManager
}

// NewStateManager creates a new state manager with persistence
func NewStateManager(config *config.BlockchainConfig) *StateManager {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		log.Printf("⚠️  Failed to create data directory: %v", err)
	}

	// Initialize database
	dbPath := filepath.Join(config.DataDir, "blockchain.db")
	db, err := database.NewDatabase(dbPath)
	if err != nil {
		log.Printf("⚠️  Failed to initialize database: %v, falling back to JSON snapshots", err)
		db = nil
	} else {
		log.Printf("✅ Database initialized successfully at %s", dbPath)
	}

	// Initialize backup and recovery managers
	var backupManager *database.BackupManager
	var recoveryManager *database.RecoveryManager

	backupDir := filepath.Join(config.DataDir, "backups")
	fallbackDir := filepath.Join(config.DataDir, "state_snapshots")

	// Always create backup manager, with fallback if database is nil
	backupManager = database.NewBackupManagerWithFallback(db, backupDir, fallbackDir)
	recoveryManager = database.NewRecoveryManager(db, backupDir)
	log.Printf("✅ Backup and recovery managers initialized")

	sm := &StateManager{
		balances:         make(map[string]int64),
		config:           config,
		snapshotPath:     fallbackDir,
		lastSnapshot:     time.Now(),
		snapshotInterval: time.Minute * 5, // Take snapshots every 5 minutes
		checksums:        make(map[string]string),
		accounts:         make(map[string]*database.Account),
		db:               db,
		backupManager:    backupManager,
		recoveryManager:  recoveryManager,
		contracts:        make(map[string]*vm.Contract), // Initialize contract registry
		proposals:        make(map[string]*Proposal),
		votes:            make(map[string][]*Vote),
		oracleData:       make(map[string]OracleData),
	}

	// Create snapshot directory if it doesn't exist
	if err := os.MkdirAll(sm.snapshotPath, 0755); err != nil {
		log.Printf("Failed to create snapshot directory: %v", err)
	}

	// Try to load latest snapshot
	if err := sm.loadLatestSnapshot(); err != nil {
		log.Printf("Failed to load latest snapshot: %v", err)
	}

	// Start periodic snapshot routine
	go sm.startPeriodicSnapshots()

	return sm
}

// GetDatabase returns the underlying database instance
func (sm *StateManager) GetDatabase() *database.Database {
	return sm.db
}

// Add at the top of the file, after imports
func shortAddr(addr string) string {
	if len(addr) > 16 {
		return addr[:16] + "..."
	}
	return addr
}

// Define custom error types
var (
	ErrSnapshotCorrupt = errors.New("state snapshot is corrupt or invalid")
	ErrStateCorrupt    = errors.New("blockchain state is corrupt or inconsistent")
)

// updateState updates the blockchain state with a new block's transactions
func (sm *StateManager) updateState(block *block.Block) error {
	log.Printf("🔄 updateState: Starting state update for block %d with %d transactions", block.Index, len(block.Transactions))
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for i, tx := range block.Transactions {
		log.Printf("💸 updateState: Processing transaction %d/%d - Sender: %s, Recipient: %s, Amount: %d, Fee: %d",
			i+1, len(block.Transactions), shortAddr(tx.Sender), shortAddr(tx.Recipient), tx.Amount, tx.Fee)

		sender := tx.Sender
		recipient := tx.Recipient
		if len(sender) > 42 && sender[:2] != "0x" {
			pubKeyBytes, _ := hex.DecodeString(sender)
			sender = wallet.PublicKeyToAddress(pubKeyBytes)
		}
		if len(recipient) > 42 && recipient[:2] != "0x" {
			pubKeyBytes, _ := hex.DecodeString(recipient)
			recipient = wallet.PublicKeyToAddress(pubKeyBytes)
		}

		// Only apply to regular transfers
		if tx.Type == transaction.TxTypeRegular {
			if sender != "network" {
				senderAcct := sm.getAccountUnlocked(sender)
				// Check for sufficient funds (amount + fee)
				if senderAcct.Balance < tx.Amount+tx.Fee {
					log.Printf("❌ updateState: Insufficient funds for sender %s (balance: %d, required: %d)",
						shortAddr(sender), senderAcct.Balance, tx.Amount+tx.Fee)
					return fmt.Errorf("insufficient funds for sender %s: balance %d, required %d", sender, senderAcct.Balance, tx.Amount+tx.Fee)
				}
				senderAcct.Balance -= tx.Amount + tx.Fee
				senderAcct.Nonce++ // Increment nonce for sender
				sm.setAccountUnlocked(senderAcct)
				log.Printf("✅ updateState: Updated sender account - New balance: %d", senderAcct.Balance)
			} else {
				log.Printf("🌐 updateState: Processing network reward transaction")
			}

			log.Printf("💰 updateState: About to update recipient account: %s", shortAddr(recipient))
			recipientAcct := sm.getAccountUnlocked(recipient)
			recipientAcct.Balance += tx.Amount
			sm.setAccountUnlocked(recipientAcct)
			log.Printf("✅ updateState: Updated recipient account - New balance: %d", recipientAcct.Balance)

			// Credit fee to block proposer (validator)
			if tx.Fee > 0 && block.Validator != "" && block.Validator != "GENESIS_VALIDATOR" {
				validatorAddr := block.Validator
				if len(validatorAddr) > 42 && validatorAddr[:2] != "0x" {
					pubKeyBytes, _ := hex.DecodeString(validatorAddr)
					validatorAddr = wallet.PublicKeyToAddress(pubKeyBytes)
				}
				validatorAcct := sm.getAccountUnlocked(validatorAddr)
				validatorAcct.Balance += tx.Fee
				sm.setAccountUnlocked(validatorAcct)
				log.Printf("💎 updateState: Credited fee %d to validator %s", tx.Fee, shortAddr(validatorAddr))
			}
			continue
		}

		// Handle contract transactions
		if tx.Type == transaction.TxTypeDeploy {
			// Parse JSON contract from tx.Data
			var jsonContract vm.JSONContract
			if err := json.Unmarshal([]byte(tx.Data), &jsonContract); err != nil {
				log.Printf("❌ Failed to parse JSON contract: %v", err)
				continue
			}
			// Deploy contract
			contract, err := vm.DeployJSONContract(sender, &jsonContract, true) // Default to upgradable for now
			if err != nil {
				log.Printf("❌ Failed to deploy contract: %v", err)
				continue
			}
			sm.SetContract(contract.Address, contract)
			log.Printf("🚀 Contract '%s' deployed at %s by %s", contract.Name, shortAddr(contract.Address), shortAddr(sender))
			continue
		} else if tx.Type == transaction.TxTypeCall {
			// Parse call data from tx.Data (JSON: {"function":..., "args":...})
			var call struct {
				Function string        `json:"function"`
				Args     []interface{} `json:"args"`
			}
			if err := json.Unmarshal([]byte(tx.Data), &call); err != nil {
				log.Printf("❌ Failed to parse contract call data: %v", err)
				continue
			}

			// SPECIAL CASE: Intercept system marketplace contract
			if recipient == vm.MarketplaceContractAddress && sm.marketplaceHelper != nil {
				execCtx := &vm.ExecutionContext{
					State:     sm.stateAdapter,
					Timestamp: time.Now().Unix(),
					Caller:    sender,
					Value:     tx.Amount,
					GasLimit:  uint64(tx.Fee),
				}
				
				if call.Function == "createOrder" && len(call.Args) >= 4 {
					seller := fmt.Sprint(call.Args[1])
					// Handle floating point parsing from JSON correctly
					amountFloat, ok := call.Args[2].(float64)
					var amount int64
					if ok {
						amount = int64(amountFloat)
					} else {
						amount = tx.Amount // fallback
					}
					orderId := fmt.Sprint(call.Args[3])
					_, err := sm.marketplaceHelper.CreateOrder(execCtx, sender, seller, amount, orderId)
					if err != nil {
						log.Printf("❌ Marketplace createOrder failed: %v", err)
					}
				} else if call.Function == "releaseFunds" && len(call.Args) >= 1 {
					orderId := fmt.Sprint(call.Args[0])
					err := sm.marketplaceHelper.ReleaseFunds(execCtx, sender, orderId)
					if err != nil {
						log.Printf("❌ Marketplace releaseFunds failed: %v", err)
					}
				} else if call.Function == "refundBuyer" && len(call.Args) >= 1 {
					orderId := fmt.Sprint(call.Args[0])
					err := sm.marketplaceHelper.RefundBuyer(execCtx, sender, orderId)
					if err != nil {
						log.Printf("❌ Marketplace refundBuyer failed: %v", err)
					}
				} else if call.Function == "raiseDispute" && len(call.Args) >= 1 {
					orderId := fmt.Sprint(call.Args[0])
					err := sm.marketplaceHelper.RaiseDispute(execCtx, sender, orderId)
					if err != nil {
						log.Printf("❌ Marketplace raiseDispute failed: %v", err)
					}
				} else if call.Function == "resolveDispute" && len(call.Args) >= 2 {
					orderId := fmt.Sprint(call.Args[0])
					payBuyer, ok := call.Args[1].(bool)
					if !ok {
						log.Printf("❌ Marketplace resolveDispute failed: invalid second argument, must be boolean")
					} else {
						err := sm.marketplaceHelper.ResolveDispute(execCtx, sender, orderId, payBuyer)
						if err != nil {
							log.Printf("❌ Marketplace resolveDispute failed: %v", err)
						}
					}
				}
				log.Printf("⚙️ System Contract '%s' executed at %s by %s", call.Function, shortAddr(recipient), shortAddr(sender))
				continue
			}

			// Load contract
			contract, ok := sm.GetContract(recipient)
			if !ok {
				log.Printf("❌ Contract not found at %s", shortAddr(recipient))
				continue
			}
			// Prepare execution context
			execCtx := &vm.ExecutionContext{
				Caller:   sender,
				Value:    tx.Amount,
				GasLimit: uint64(tx.Fee),
			}
			// Create VM instance
			vmInstance := &vm.VM{
				StateManager: sm,
			}
			// Initialize VM memory with contract storage
			vmInstance.Memory = make(map[string]int64)
			for k, v := range contract.Storage {
				if ival, ok := toInt64Safe(v); ok {
					vmInstance.Memory[k] = ival
				}
			}
			// Execute function
			if err := contract.CallFunction(call.Function, call.Args, vmInstance, execCtx); err != nil {
				log.Printf("❌ Contract function '%s' execution failed: %v", call.Function, err)
				continue
			}
			// After execution, write VM memory back to contract storage
			for k, v := range vmInstance.Memory {
				contract.Storage[k] = v
			}
			contract.UpdatedAt = time.Now().Unix()
			sm.SetContract(contract.Address, contract)
			log.Printf("⚙️ Contract '%s' function '%s' executed at %s by %s (gas used: %d)",
				contract.Name, call.Function, shortAddr(recipient), shortAddr(sender), vmInstance.GetGasUsed())
			continue
		}

		// Handle governance transactions
		if tx.Type == transaction.TxTypeProposal {
			// Parse proposal data from tx.Data (JSON: {"description":..., "actions":..., "duration":...})
			var proposalData struct {
				Description string `json:"description"`
				Actions     string `json:"actions"`
				Duration    int64  `json:"duration"`
			}
			if err := json.Unmarshal([]byte(tx.Data), &proposalData); err != nil {
				log.Printf("❌ Failed to parse proposal data: %v", err)
				continue
			}
			startBlock := int64(block.Index)
			endBlock := startBlock + proposalData.Duration
			proposal := sm.SubmitProposal(sender, proposalData.Description, proposalData.Actions, startBlock, endBlock)
			proposal.State = ProposalActive
			log.Printf("🗳️ Proposal submitted by %s: %s (ID: %s)", shortAddr(sender), proposal.Description, proposal.ID)
			continue
		} else if tx.Type == transaction.TxTypeVote {
			// Parse vote data from tx.Data (JSON: {"proposalID":..., "choice":..., "weight":...})
			var voteData struct {
				ProposalID string `json:"proposalID"`
				Choice     string `json:"choice"`
				Weight     int64  `json:"weight"`
			}
			if err := json.Unmarshal([]byte(tx.Data), &voteData); err != nil {
				log.Printf("❌ Failed to parse vote data: %v", err)
				continue
			}
			if err := sm.CastVote(voteData.ProposalID, sender, voteData.Choice, voteData.Weight); err != nil {
				log.Printf("❌ Failed to cast vote: %v", err)
				continue
			}
			log.Printf("🗳️ Vote cast by %s on proposal %s: %s (%d)", shortAddr(sender), voteData.ProposalID, voteData.Choice, voteData.Weight)
			// Tally proposal if voting period ended
			sm.TallyProposal(voteData.ProposalID, int64(block.Index))
			continue
		}

		// Handle staking transactions
		if tx.Type == transaction.TxTypeStake {
			if sender == "network" {
				log.Printf("❌ updateState: Network cannot stake")
				return fmt.Errorf("network cannot stake")
			}
			minStake := int64(1)
			if sm.config != nil && sm.config.MinStake > 0 {
				minStake = int64(sm.config.MinStake)
			}
			if tx.Amount < minStake {
				log.Printf("❌ updateState: Stake amount too low: %d (min: %d)", tx.Amount, minStake)
				return fmt.Errorf("stake amount too low: %d (min: %d)", tx.Amount, minStake)
			}
			senderAcct := sm.getAccountUnlocked(sender)
			if senderAcct.Balance < tx.Amount+tx.Fee {
				log.Printf("❌ updateState: Insufficient funds for staking by %s (balance: %d, required: %d)", shortAddr(sender), senderAcct.Balance, tx.Amount+tx.Fee)
				return fmt.Errorf("insufficient funds for staking by %s: balance %d, required %d", sender, senderAcct.Balance, tx.Amount+tx.Fee)
			}
			senderAcct.Balance -= tx.Amount + tx.Fee
			senderAcct.Nonce++
			sm.setAccountUnlocked(senderAcct)
			log.Printf("✅ updateState: Deducted stake and fee from %s, new balance: %d", shortAddr(sender), senderAcct.Balance)
			// Register or update validator in consensus manager
			if sm.consensusManager != nil {
				err := sm.consensusManager.OnChainStake(sender, uint64(tx.Amount))
				if err != nil {
					log.Printf("❌ updateState: ConsensusManager.OnChainStake failed: %v", err)
					return fmt.Errorf("consensus manager staking failed: %v", err)
				}
			}
			// Credit fee to block proposer (validator)
			if tx.Fee > 0 && block.Validator != "" && block.Validator != "GENESIS_VALIDATOR" {
				validatorAddr := block.Validator
				if len(validatorAddr) > 42 && validatorAddr[:2] != "0x" {
					pubKeyBytes, _ := hex.DecodeString(validatorAddr)
					validatorAddr = wallet.PublicKeyToAddress(pubKeyBytes)
				}
				validatorAcct := sm.getAccountUnlocked(validatorAddr)
				validatorAcct.Balance += tx.Fee
				sm.setAccountUnlocked(validatorAcct)
				log.Printf("💎 updateState: Credited fee %d to validator %s", tx.Fee, shortAddr(validatorAddr))
			}
			continue
		}
	}

	log.Printf("📸 updateState: Checking if snapshot should be created...")
	// Check if we should create a new snapshot
	if time.Since(sm.lastSnapshot) >= time.Hour {
		log.Printf("📸 updateState: Creating new snapshot...")
		if err := sm.createSnapshot(int64(block.Index)); err != nil {
			log.Printf("❌ updateState: Failed to create state snapshot: %v", err)
		} else {
			log.Printf("✅ updateState: Snapshot created successfully")
		}
	} else {
		log.Printf("⏰ updateState: No snapshot needed yet (last: %v ago)", time.Since(sm.lastSnapshot))
	}

	log.Printf("✅ updateState: State update completed successfully for block %d", block.Index)
	return nil
}

// createSnapshot creates a new state snapshot
func (sm *StateManager) createSnapshot(blockHeight int64) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	snapshot := StateSnapshot{
		Balances:    make(map[string]int64),
		BlockHeight: blockHeight,
		Timestamp:   time.Now(),
	}

	// Copy current balances
	for addr, acct := range sm.accounts {
		snapshot.Balances[addr] = acct.Balance
	}

	// Calculate checksum
	checksum := calculateStateChecksum(snapshot)
	snapshot.Checksum = checksum

	// Marshal to JSON
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %v", err)
	}

	// Write to file
	filename := fmt.Sprintf("snapshot_%d_%s.json", blockHeight, time.Now().Format("20060102_150405"))
	path := filepath.Join(sm.snapshotPath, filename)
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot: %v", err)
	}

	// Clean up old snapshots
	sm.cleanupOldSnapshots()

	sm.lastSnapshot = time.Now()
	return nil
}

// loadLatestSnapshot loads the most recent state snapshot
func (sm *StateManager) loadLatestSnapshot() error {
	files, err := ioutil.ReadDir(sm.snapshotPath)
	if err != nil {
		return fmt.Errorf("failed to read snapshot directory: %v", err)
	}

	if len(files) == 0 {
		// No snapshots to load; proceed with fresh state
		return nil
	}

	// Find latest snapshot
	var latestFile os.FileInfo
	for _, file := range files {
		if latestFile == nil || file.ModTime().After(latestFile.ModTime()) {
			latestFile = file
		}
	}

	// Read and parse snapshot
	data, err := ioutil.ReadFile(filepath.Join(sm.snapshotPath, latestFile.Name()))
	if err != nil {
		return fmt.Errorf("failed to read snapshot file: %v", err)
	}

	var snapshot StateSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return fmt.Errorf("failed to unmarshal snapshot: %v", err)
	}

	// Verify checksum
	if snapshot.Checksum != calculateStateChecksum(snapshot) {
		return ErrSnapshotCorrupt
	}

	// Update state
	sm.mu.Lock()
	for addr, bal := range snapshot.Balances {
		if acct, exists := sm.accounts[addr]; exists {
			acct.Balance = bal
		} else {
			sm.accounts[addr] = &database.Account{Address: addr, Balance: bal}
		}
	}
	sm.lastSnapshot = snapshot.Timestamp
	sm.mu.Unlock()

	return nil
}

// cleanupOldSnapshots removes old snapshots, keeping only the last N
func (sm *StateManager) cleanupOldSnapshots() {
	files, err := ioutil.ReadDir(sm.snapshotPath)
	if err != nil {
		log.Printf("Failed to read snapshot directory: %v", err)
		return
	}

	// Keep only the last 5 snapshots
	maxSnapshots := 5
	if len(files) <= maxSnapshots {
		return
	}

	// Sort by modification time
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().After(files[j].ModTime())
	})

	// Remove old snapshots
	for _, file := range files[maxSnapshots:] {
		path := filepath.Join(sm.snapshotPath, file.Name())
		if err := os.Remove(path); err != nil {
			log.Printf("Failed to remove old snapshot %s: %v", file.Name(), err)
		}
	}
}

// calculateStateChecksum calculates a checksum for the state snapshot
func calculateStateChecksum(snapshot StateSnapshot) string {
	// Create a deterministic string representation
	data := fmt.Sprintf("%d_%s", snapshot.BlockHeight, snapshot.Timestamp)
	// Sort addresses for deterministic order
	addresses := make([]string, 0, len(snapshot.Balances))
	for addr := range snapshot.Balances {
		addresses = append(addresses, addr)
	}
	sort.Strings(addresses)
	for _, addr := range addresses {
		balance := snapshot.Balances[addr]
		data += fmt.Sprintf("_%s_%d", addr, balance)
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
}

// GetBalance returns the balance for a given address
func (sm *StateManager) GetBalance(address string) int64 {
	acct := sm.GetAccount(address)
	return acct.Balance
}

// SetBalance sets the balance for a given address
func (sm *StateManager) SetBalance(address string, amount int64) {
	acct := sm.GetAccount(address)
	acct.Balance = amount
	sm.SetAccount(acct)
}

// GetNonce returns the nonce for a given address
func (sm *StateManager) GetNonce(address string) uint64 {
	acct := sm.GetAccount(address)
	return acct.Nonce
}

// SetNonce sets the nonce for a given address
func (sm *StateManager) SetNonce(address string, nonce uint64) {
	acct := sm.GetAccount(address)
	acct.Nonce = nonce
	sm.SetAccount(acct)
}

// SyncWalletBalance synchronizes a wallet's balance with the state manager
func (sm *StateManager) SyncWalletBalance(w *wallet.Wallet) {
	if w == nil {
		return
	}
	address := wallet.PublicKeyToAddress(w.PublicKey)
	stateBalance := sm.GetBalance(address)
	w.SetBalance(stateBalance)
}

// GetAllBalances returns a copy of all balances
func (sm *StateManager) GetAllBalances() map[string]int64 {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	balances := make(map[string]int64)
	for addr, acct := range sm.accounts {
		balances[addr] = acct.Balance
	}
	return balances
}

// RecoverState attempts to recover the state from the latest valid snapshot
func (sm *StateManager) RecoverState() error {
	// Try to load the latest snapshot
	if err := sm.loadLatestSnapshot(); err != nil {
		if err == ErrSnapshotCorrupt {
			log.Printf("RecoverState: Snapshot is corrupt. Attempting to recover from last valid state.")
			// If snapshot is corrupt, try to load the last valid state from the last snapshot
			files, err := ioutil.ReadDir(sm.snapshotPath)
			if err != nil {
				return fmt.Errorf("failed to read snapshot directory for recovery: %v", err)
			}
			if len(files) == 0 {
				return fmt.Errorf("no valid snapshots found to recover from: %v", ErrStateCorrupt)
			}
			var latestFile os.FileInfo
			for _, file := range files {
				if latestFile == nil || file.ModTime().After(latestFile.ModTime()) {
					latestFile = file
				}
			}
			data, err := ioutil.ReadFile(filepath.Join(sm.snapshotPath, latestFile.Name()))
			if err != nil {
				return fmt.Errorf("failed to read latest valid snapshot for recovery: %v", err)
			}
			var snapshot StateSnapshot
			if err := json.Unmarshal(data, &snapshot); err != nil {
				return fmt.Errorf("failed to unmarshal latest valid snapshot for recovery: %v", err)
			}
			if snapshot.Checksum != calculateStateChecksum(snapshot) {
				return ErrSnapshotCorrupt
			}
			sm.mu.Lock()
			sm.balances = snapshot.Balances
			sm.lastSnapshot = snapshot.Timestamp
			sm.mu.Unlock()
			log.Printf("Recovered state from latest valid snapshot: %s", latestFile.Name())
			return nil
		}
		return fmt.Errorf("failed to recover state: %v", err)
	}
	return nil
}

// ExportState exports the current state to a file
func (sm *StateManager) ExportState(path string) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	balances := make(map[string]int64)
	for addr, acct := range sm.accounts {
		balances[addr] = acct.Balance
	}

	data, err := json.MarshalIndent(balances, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %v", err)
	}

	return nil
}

// ImportState imports state from a file
func (sm *StateManager) ImportState(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read state file: %v", err)
	}

	var balances map[string]int64
	if err := json.Unmarshal(data, &balances); err != nil {
		return fmt.Errorf("failed to unmarshal state: %v", err)
	}

	sm.mu.Lock()
	for addr, bal := range balances {
		if acct, exists := sm.accounts[addr]; exists {
			acct.Balance = bal
		} else {
			sm.accounts[addr] = &database.Account{Address: addr, Balance: bal}
		}
	}
	sm.mu.Unlock()

	return nil
}

// startPeriodicSnapshots starts a goroutine for periodic state snapshots
func (sm *StateManager) startPeriodicSnapshots() {
	ticker := time.NewTicker(sm.snapshotInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := sm.createSnapshot(0); err != nil {
			log.Printf("Failed to create snapshot: %v", err)
		}
	}
}

// GetStateChecksum returns the current state checksum
func (sm *StateManager) GetStateChecksum() string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	balances := make(map[string]int64)
	for addr, acct := range sm.accounts {
		balances[addr] = acct.Balance
	}
	snapshot := StateSnapshot{
		Balances:  balances,
		Timestamp: time.Now(),
	}
	return calculateStateChecksum(snapshot)
}

// VerifyStateIntegrity verifies the integrity of the current state
func (sm *StateManager) VerifyStateIntegrity() bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	balances := make(map[string]int64)
	for addr, acct := range sm.accounts {
		balances[addr] = acct.Balance
	}
	snapshot := StateSnapshot{
		Balances:  balances,
		Timestamp: time.Now(),
	}
	checksum := calculateStateChecksum(snapshot)

	// Verify against stored checksums
	for _, storedChecksum := range sm.checksums {
		if checksum == storedChecksum {
			return true
		}
	}
	return false
}

// GetAccount returns the Account for a given address, or creates it if not present
func (sm *StateManager) GetAccount(address string) *database.Account {
	// Try database first if available
	if sm.db != nil {
		dbAccount, err := sm.db.GetAccount(address)
		if err != nil {
			log.Printf("⚠️  Failed to get account from database: %v", err)
		} else if dbAccount != nil {
			log.Printf("✅ GetAccount: Found account in database for %s with balance %d", shortAddr(address), dbAccount.Balance)
			return dbAccount
		}
	}

	// Fallback to in-memory storage
	sm.mu.Lock()
	defer sm.mu.Unlock()

	log.Printf("🔍 GetAccount: Getting account for %s", shortAddr(address))
	acct, exists := sm.accounts[address]
	if !exists {
		log.Printf("🆕 GetAccount: Creating new account for %s", shortAddr(address))
		acct = &database.Account{Address: address, Balance: 0, Nonce: 0}
		sm.accounts[address] = acct
	} else {
		log.Printf("✅ GetAccount: Found existing account for %s with balance %d", shortAddr(address), acct.Balance)
	}
	return acct
}

// getAccountUnlocked returns the Account for a given address without acquiring locks
// Use this when you already have the lock (e.g., from updateState)
func (sm *StateManager) getAccountUnlocked(address string) *database.Account {
	log.Printf("🔍 getAccountUnlocked: Getting account for %s", shortAddr(address))
	acct, exists := sm.accounts[address]
	if !exists {
		log.Printf("🆕 getAccountUnlocked: Creating new account for %s", shortAddr(address))
		acct = &database.Account{Address: address, Balance: 0, Nonce: 0}
		sm.accounts[address] = acct
	} else {
		log.Printf("✅ getAccountUnlocked: Found existing account for %s with balance %d", shortAddr(address), acct.Balance)
	}
	return acct
}

// SetAccount sets the Account for a given address
func (sm *StateManager) SetAccount(acct *database.Account) {
	// Try database first if available
	if sm.db != nil {
		if err := sm.db.SetAccount(acct); err != nil {
			log.Printf("⚠️  Failed to set account in database: %v", err)
		} else {
			log.Printf("💾 SetAccount: Saved account to database for %s with balance %d", shortAddr(acct.Address), acct.Balance)
		}
	}

	// Also update in-memory storage
	sm.mu.Lock()
	defer sm.mu.Unlock()

	log.Printf("💾 SetAccount: Setting account for %s with balance %d", shortAddr(acct.Address), acct.Balance)
	sm.accounts[acct.Address] = acct
}

// setAccountUnlocked sets the Account for a given address without acquiring locks
// Use this when you already have the lock (e.g., from updateState)
func (sm *StateManager) setAccountUnlocked(acct *database.Account) {
	log.Printf("💾 setAccountUnlocked: Setting account for %s with balance %d", shortAddr(acct.Address), acct.Balance)
	sm.accounts[acct.Address] = acct
}

// GetContract retrieves a contract by address
func (sm *StateManager) GetContract(address string) (*vm.Contract, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	c, ok := sm.contracts[address]
	return c, ok
}

// SetContract stores or updates a contract
func (sm *StateManager) SetContract(address string, contract *vm.Contract) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.contracts[address] = contract
}

// InitTokenContract initializes the TCOIN token contract at chain startup.
// This deploys the system token contract, creates the state adapter, and
// initializes the token contract helper for subsequent token operations.
func (sm *StateManager) InitTokenContract(treasuryAddress string, initialTreasuryBalance int64) {
	// Create the state adapter
	sa := NewStateAdapter(sm)
	sm.stateAdapter = sa

	// Create and register the token contract
	tokenContract := vm.CreateSystemTokenContract(treasuryAddress)
	sm.SetContract(tokenContract.Address, tokenContract)

	// Initialize contract storage via the adapter
	sa.InitializeTokenContract(treasuryAddress, initialTreasuryBalance)

	// Create the token helper
	sm.tokenHelper = vm.NewTokenContractHelper(tokenContract)

	log.Printf("🪙 TCOIN Token Contract initialized at %s (treasury: %s, supply: %d/%d)",
		vm.TokenContractAddress, shortAddr(treasuryAddress), initialTreasuryBalance, vm.TokenMaxSupply)
}

// GetStateAdapter returns the state adapter for VM interactions
func (sm *StateManager) GetStateAdapter() *StateAdapterImpl {
	return sm.stateAdapter
}

// GetTokenHelper returns the token contract helper for token operations
func (sm *StateManager) GetTokenHelper() *vm.TokenContractHelper {
	return sm.tokenHelper
}

// InitStakingContract initializes the Staking system contract.
func (sm *StateManager) InitStakingContract() {
	if sm.stateAdapter == nil {
		log.Printf("⚠️  Cannot init staking contract: StateAdapter not initialized. Call InitTokenContract first.")
		return
	}

	stakingContract := vm.CreateSystemStakingContract()
	sm.SetContract(stakingContract.Address, stakingContract)

	// Initialize staking parameters in contract storage
	sm.stateAdapter.SetContractStorage(vm.StakingContractAddress, "min_stake", vm.StakingMinStake)
	sm.stateAdapter.SetContractStorage(vm.StakingContractAddress, "max_validators", vm.StakingMaxValidators)
	sm.stateAdapter.SetContractStorage(vm.StakingContractAddress, "slashing_penalty", vm.StakingSlashingPenalty)
	sm.stateAdapter.SetContractStorage(vm.StakingContractAddress, "lock_period", vm.StakingLockPeriod)
	sm.stateAdapter.SetContractStorage(vm.StakingContractAddress, "reward_rate", vm.StakingRewardRate)
	sm.stateAdapter.SetContractStorage(vm.StakingContractAddress, "total_staked", 0)
	sm.stateAdapter.SetContractStorage(vm.StakingContractAddress, "active_validators", 0)

	sm.stakingHelper = vm.NewStakingContractHelper(stakingContract)

	log.Printf("⚡ Staking Contract initialized at %s (minStake: %d, maxValidators: %d, lockPeriod: %ds)",
		vm.StakingContractAddress, vm.StakingMinStake, vm.StakingMaxValidators, vm.StakingLockPeriod)
}

// GetStakingHelper returns the staking contract helper.
func (sm *StateManager) GetStakingHelper() *vm.StakingContractHelper {
	return sm.stakingHelper
}

// InitMarketplaceContract initializes the system Marketplace contract for escrow and commerce
func (sm *StateManager) InitMarketplaceContract() {
	if sm.stateAdapter == nil {
		log.Printf("⚠️  Cannot init marketplace contract: StateAdapter not initialized. Call InitTokenContract first.")
		return
	}

	marketplaceContract := vm.CreateSystemMarketplaceContract()
	sm.SetContract(marketplaceContract.Address, marketplaceContract)

	// We don't overwrite generic storage here because the factory function
	// CreateSystemMarketplaceContract already set the defaults (fee_rate, etc.).
	// We'll migrate those defaults to the real StateAdapter storage.
	for key, value := range marketplaceContract.Storage {
		if valInt, ok := value.(int64); ok {
			sm.stateAdapter.SetContractStorage(marketplaceContract.Address, key, valInt)
		} else if valStr, ok := value.(string); ok {
			sm.stateAdapter.SetStringStorage(marketplaceContract.Address, key, valStr)
		}
	}

	sm.marketplaceHelper = vm.NewMarketplaceContractHelper(marketplaceContract)

	log.Printf("🛒 Marketplace Contract initialized at %s (feeRate: %d bp)",
		vm.MarketplaceContractAddress, vm.MarketplaceFeeRate)
}

// GetMarketplaceHelper returns the marketplace contract helper.
func (sm *StateManager) GetMarketplaceHelper() *vm.MarketplaceContractHelper {
	return sm.marketplaceHelper
}

// InitGovernanceContract initializes the system Governance contract for parameter override tracking
func (sm *StateManager) InitGovernanceContract() {
	if sm.stateAdapter == nil {
		log.Printf("⚠️  Cannot init governance contract: StateAdapter not initialized. Call InitTokenContract first.")
		return
	}

	governanceContract := vm.CreateSystemGovernanceContract()
	sm.SetContract(governanceContract.Address, governanceContract)

	// We don't overwrite generic storage here because the factory function
	// CreateSystemGovernanceContract already set the defaults (min_staking...).
	for key, value := range governanceContract.Storage {
		if valInt, ok := value.(int64); ok {
			sm.stateAdapter.SetContractStorage(governanceContract.Address, key, valInt)
		} else if valStr, ok := value.(string); ok {
			sm.stateAdapter.SetStringStorage(governanceContract.Address, key, valStr)
		}
	}

	sm.governanceHelper = vm.NewGovernanceContractHelper(governanceContract)

	stake, _ := sm.stateAdapter.GetContractStorage(vm.GovernanceContractAddress, "min_proposal_stake")
	log.Printf("🏛️  Governance Contract initialized at %s (minProposalStake: %d TCOIN)",
		vm.GovernanceContractAddress, stake)
}

// GetGovernanceHelper returns the governance contract helper.
func (sm *StateManager) GetGovernanceHelper() *vm.GovernanceContractHelper {
	return sm.governanceHelper
}

// SubmitProposal adds a new proposal
func (sm *StateManager) SubmitProposal(proposer, description, actions string, startBlock, endBlock int64) *Proposal {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	id := fmt.Sprintf("proposal_%d", len(sm.proposals)+1)
	proposal := &Proposal{
		ID:           id,
		Proposer:     proposer,
		Description:  description,
		Actions:      actions,
		State:        ProposalPending,
		VotesFor:     0,
		VotesAgainst: 0,
		StartBlock:   startBlock,
		EndBlock:     endBlock,
		Voters:       make(map[string]bool),
	}
	sm.proposals[id] = proposal
	return proposal
}

// CastVote records a vote for a proposal
func (sm *StateManager) CastVote(proposalID, voter, choice string, weight int64) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	proposal, ok := sm.proposals[proposalID]
	if !ok {
		return fmt.Errorf("proposal not found")
	}
	if proposal.Voters[voter] {
		return fmt.Errorf("voter has already voted")
	}
	if proposal.State != ProposalActive {
		return fmt.Errorf("proposal not active")
	}
	vote := &Vote{
		ProposalID: proposalID,
		Voter:      voter,
		Choice:     choice,
		Weight:     weight,
	}
	if choice == "for" {
		proposal.VotesFor += weight
	} else if choice == "against" {
		proposal.VotesAgainst += weight
	} else {
		return fmt.Errorf("invalid vote choice")
	}
	proposal.Voters[voter] = true
	sm.votes[proposalID] = append(sm.votes[proposalID], vote)
	return nil
}

// TallyProposal updates proposal state based on votes and block height
func (sm *StateManager) TallyProposal(proposalID string, currentBlock int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	proposal, ok := sm.proposals[proposalID]
	if !ok {
		return
	}
	if proposal.State == ProposalActive && currentBlock > proposal.EndBlock {
		if proposal.VotesFor > proposal.VotesAgainst {
			proposal.State = ProposalSucceeded
		} else {
			proposal.State = ProposalFailed
		}
	}
}

// SetOracleData sets the value for a given oracle key
func (sm *StateManager) SetOracleData(key, value, source string, timestamp int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.oracleData[key] = OracleData{
		Key:       key,
		Value:     value,
		Timestamp: timestamp,
		Source:    source,
	}
}

// GetOracleData retrieves the latest value for a given oracle key
func (sm *StateManager) GetOracleData(key string) (OracleData, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	data, ok := sm.oracleData[key]
	return data, ok
}

// SetConsensusManager sets the consensus manager reference for on-chain staking
func (sm *StateManager) SetConsensusManager(cm *ConsensusManager) {
	sm.consensusManager = cm
}

// MigrateToDatabase migrates existing JSON snapshots to the database
func (sm *StateManager) MigrateToDatabase() error {
	if sm.db == nil {
		return fmt.Errorf("database not available")
	}

	log.Printf("🔄 Starting migration from JSON snapshots to database...")

	// Load the latest snapshot
	if err := sm.loadLatestSnapshot(); err != nil {
		log.Printf("⚠️  No existing snapshots to migrate: %v", err)
		return nil
	}

	// Migrate all accounts from in-memory to database
	sm.mu.RLock()
	accountCount := 0
	for address, account := range sm.accounts {
		dbAccount := &database.Account{
			Address:      address,
			Balance:      account.Balance,
			Nonce:        account.Nonce,
			IsValidator:  account.IsValidator,
			StakedAmount: account.StakedAmount,
		}
		if err := sm.db.SetAccount(dbAccount); err != nil {
			log.Printf("⚠️  Failed to migrate account %s: %v", shortAddr(address), err)
		} else {
			accountCount++
		}
	}
	sm.mu.RUnlock()

	log.Printf("✅ Migration completed: %d accounts migrated to database", accountCount)
	return nil
}

// CloseDatabase closes the database connection
func (sm *StateManager) CloseDatabase() error {
	if sm.db != nil {
		return sm.db.Close()
	}
	return nil
}

// Backup and recovery methods

// StartBackupSystem starts the automated backup system
func (sm *StateManager) StartBackupSystem() {
	if sm.backupManager != nil {
		sm.backupManager.StartAutomatedBackups()
		log.Printf("🔄 Automated backup system started")
	} else {
		log.Printf("⚠️ Backup system not available (database not initialized)")
	}
}

// StopBackupSystem stops the automated backup system
func (sm *StateManager) StopBackupSystem() {
	if sm.backupManager != nil {
		sm.backupManager.StopAutomatedBackups()
		log.Printf("🛑 Automated backup system stopped")
	}
}

// GetBackupStatus returns the current backup status
func (sm *StateManager) GetBackupStatus() map[string]interface{} {
	if sm.backupManager != nil {
		return sm.backupManager.GetBackupStatus()
	}
	return map[string]interface{}{
		"status": "unavailable",
		"error":  "backup system not initialized",
	}
}

// GetBackupList returns a list of all backups
func (sm *StateManager) GetBackupList() []*database.BackupInfo {
	if sm.backupManager != nil {
		return sm.backupManager.GetBackupList()
	}
	return []*database.BackupInfo{}
}

// CreateManualBackup creates a manual backup
func (sm *StateManager) CreateManualBackup() error {
	if sm.backupManager != nil {
		return sm.backupManager.CreateBackup()
	}
	return fmt.Errorf("backup system not available")
}

// PerformAutomaticRecovery performs automatic recovery if corruption is detected
func (sm *StateManager) PerformAutomaticRecovery() error {
	if sm.recoveryManager != nil {
		return sm.recoveryManager.AutomaticRecovery()
	}
	return fmt.Errorf("recovery system not available")
}

// Helper to safely convert interface{} to int64
func toInt64Safe(val interface{}) (int64, bool) {
	switch v := val.(type) {
	case int:
		return int64(v), true
	case int64:
		return v, true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}
