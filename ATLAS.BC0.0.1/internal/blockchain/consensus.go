package blockchain

import (
	"atlas-blockchain/pkg/block"
	"atlas-blockchain/pkg/config"
	"atlas-blockchain/pkg/sharding"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/wallet"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"

	"sort"
	"sync"
	"time"
)

// Constants for consensus
const (
	BLOCK_REWARD = 10 // Fixed reward for forging a block
	MIN_STAKE    = 1  // Minimum stake required to be a validator
)

// ErrNoValidators is returned when there are no validators with sufficient stake.
var ErrNoValidators = errors.New("no validators with sufficient stake available")

// ErrInvalidValidator is returned when an invalid validator attempts to forge a block.
var ErrInvalidValidator = errors.New("invalid validator")

// KYCInfo holds KYC/AML information for a validator
type KYCInfo struct {
	FullName string
	Country  string
	IDNumber string
	Verified bool
	// Add more fields as needed
}

// Validator represents a node that can participate in block production
type Validator struct {
	Address     string
	Stake       uint64
	Delegations map[string]int64 // Address -> Amount
	LastBlock   int64
	SlashCount  int
	Active      bool
	// New fields for enhanced validator metrics
	PerformanceScore float64   // Score based on block production and validation
	Uptime           float64   // Percentage of time validator is online
	LastActive       time.Time // Last time validator was active
	BlocksProduced   uint64    // Total blocks produced
	BlocksValidated  uint64    // Total blocks validated
	ReputationScore  float64   // Overall reputation score
	SlashingHistory  []time.Time
	RewardHistory    []uint64
	TotalRewards     uint64 // Sum of all rewards

	KYC KYCInfo // KYC/AML info
}

// ConsensusManager handles validator selection, rotation, and slashing
type ConsensusManager struct {
	validators         map[string]*Validator
	mu                 sync.RWMutex
	config             *config.BlockchainConfig
	blockManager       *BlockManager
	lastRotation       time.Time
	rotationSchedule   map[string]time.Time
	performanceHistory map[string][]float64
	slashingThreshold  int
	rewardMultiplier   float64
	// New fields for block finality
	finalityThreshold  int             // Number of confirmations required for finality
	blockConfirmations map[string]int  // Maps block hash to number of confirmations
	finalizedBlocks    map[string]bool // Tracks finalized blocks
	confirmationWindow time.Duration   // Time window for counting confirmations
	// Sharding support
	shardManager *sharding.ShardManager // Add shard manager
}

// NewConsensusManager creates a new consensus manager with enhanced validator management
func NewConsensusManager(config *config.BlockchainConfig, blockManager *BlockManager) *ConsensusManager {
	// Initialize shard manager
	shardConfig := &sharding.ShardConfig{
		TotalShards:     4,  // Default to 4 shards
		ShardSize:       10, // 10 validators per shard
		CrossShardDelay: time.Second * 5,
		ConsensusType:   "pbft",
	}
	shardManager := sharding.NewShardManager(shardConfig)

	return &ConsensusManager{
		validators:         make(map[string]*Validator),
		config:             config,
		blockManager:       blockManager,
		lastRotation:       time.Now(),
		rotationSchedule:   make(map[string]time.Time),
		performanceHistory: make(map[string][]float64),
		slashingThreshold:  3,
		rewardMultiplier:   1.0,
		// Initialize finality-related fields
		finalityThreshold:  1, // Require 1 confirmation for finality (single-user testing)
		blockConfirmations: make(map[string]int),
		finalizedBlocks:    make(map[string]bool),
		confirmationWindow: time.Minute * 10, // 10-minute window for confirmations
		// Initialize sharding
		shardManager: shardManager,
	}
}

// RegisterValidator registers a new validator with enhanced metrics
// The stake amount is moved from the wallet's balance to staked amount
func (cm *ConsensusManager) RegisterValidator(w *wallet.Wallet, stake uint64, kyc KYCInfo) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	address := wallet.PublicKeyToAddress(w.PublicKey)

	if _, exists := cm.validators[address]; exists {
		return errors.New("validator already registered")
	}

	// Check if wallet has enough balance for stake
	if w == nil || !w.CanStake(stake) {
		return fmt.Errorf("insufficient balance for stake: have %d, need %d", w.GetBalance(), stake)
	}

	// Require KYC verification
	if !kyc.Verified || kyc.FullName == "" || kyc.Country == "" || kyc.IDNumber == "" {
		return errors.New("KYC verification required for validator registration")
	}

	// Move tokens from balance to staked amount
	if err := w.Stake(stake); err != nil {
		return fmt.Errorf("failed to stake tokens: %v", err)
	}

	cm.validators[address] = &Validator{
		Address:          address,
		Stake:            stake,
		Delegations:      make(map[string]int64),
		LastBlock:        0,
		SlashCount:       0,
		Active:           true,
		PerformanceScore: 1.0,
		Uptime:           1.0,
		ReputationScore:  1.0,
		SlashingHistory:  make([]time.Time, 0),
		RewardHistory:    make([]uint64, 0),
		KYC:              kyc,
	}

	// Add validator to appropriate shard
	shardID := cm.shardManager.GetShardForAddress(address)
	if err := cm.shardManager.AddValidatorToShard(shardID, address); err != nil {
		log.Printf("⚠️ Failed to add validator to shard %d: %v", shardID, err)
	} else {
		log.Printf("✅ Added validator %s to shard %d", address[:16]+"...", shardID)
	}

	cm.performanceHistory[address] = make([]float64, 0)
	return nil
}

// RegisterValidatorByAddress registers a validator using an address string instead of a wallet object
func (cm *ConsensusManager) RegisterValidatorByAddress(address string, stake uint64, kyc KYCInfo) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.validators[address]; exists {
		return errors.New("validator already registered")
	}

	// Require KYC verification
	if !kyc.Verified || kyc.FullName == "" || kyc.Country == "" || kyc.IDNumber == "" {
		return errors.New("KYC verification required for validator registration")
	}

	cm.validators[address] = &Validator{
		Address:          address,
		Stake:            stake,
		Delegations:      make(map[string]int64),
		LastBlock:        0,
		SlashCount:       0,
		Active:           true,
		PerformanceScore: 1.0,
		Uptime:           1.0,
		ReputationScore:  1.0,
		SlashingHistory:  make([]time.Time, 0),
		RewardHistory:    make([]uint64, 0),
		KYC:              kyc,
	}

	// Add validator to appropriate shard
	if cm.shardManager != nil {
		shardID := cm.shardManager.GetShardForAddress(address)
		if err := cm.shardManager.AddValidatorToShard(shardID, address); err != nil {
			log.Printf("⚠️ Failed to add validator to shard %d: %v", shardID, err)
		} else {
			log.Printf("✅ Added validator %s to shard %d", address[:16]+"...", shardID)
		}
	}

	cm.performanceHistory[address] = make([]float64, 0)
	return nil
}

// DelegateStake allows a user to delegate their stake to a validator
func (cm *ConsensusManager) DelegateStake(delegator, validator string, amount int64) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Check if validator exists
	val, exists := cm.validators[validator]
	if !exists {
		return fmt.Errorf("validator not found: %s", validator)
	}

	// Check if delegator has enough balance
	balance := cm.blockManager.state.GetBalance(delegator)
	if balance < amount {
		return fmt.Errorf("insufficient balance for delegation")
	}

	// Update delegations
	val.Delegations[delegator] = amount
	val.Stake += uint64(amount)

	// Deduct from delegator's balance
	cm.blockManager.state.SetBalance(delegator, balance-amount)

	return nil
}

// UndelegateStake allows a user to remove their delegation
func (cm *ConsensusManager) UndelegateStake(delegator, validator string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	val, exists := cm.validators[validator]
	if !exists {
		return fmt.Errorf("validator not found: %s", validator)
	}

	// Get delegation amount
	amount, exists := val.Delegations[delegator]
	if !exists {
		return fmt.Errorf("no delegation found for %s", delegator)
	}

	// Update validator's stake
	val.Stake -= uint64(amount)
	delete(val.Delegations, delegator)

	// Return stake to delegator
	balance := cm.blockManager.state.GetBalance(delegator)
	cm.blockManager.state.SetBalance(delegator, balance+amount)

	return nil
}

// ChooseValidator selects the next validator based on deterministic weighted selection
func (cm *ConsensusManager) ChooseValidator(lastBlockHash string, blockHeight int64) (*Validator, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// log.Printf("🔍 ChooseValidator: Checking available validators...")
	// log.Printf("📊 Total validators in pool: %d", len(cm.validators))

	if len(cm.validators) == 0 {
		log.Printf("❌ No validators available in pool")
		return nil, errors.New("no validators available")
	}

	// Deterministic ordering of validators
	var addresses []string
	for addr := range cm.validators {
		addresses = append(addresses, addr)
	}
	sort.Strings(addresses)

	// Calculate total weight and build weighted list
	var totalWeight float64
	weights := make(map[string]float64)

	for _, addr := range addresses {
		validator := cm.validators[addr]

		// Stake weight (40%)
		stakeWeight := float64(validator.Stake) / float64(cm.getTotalStake())
		// Performance weight (30%)
		perfWeight := validator.PerformanceScore
		// Reputation weight (20%)
		repWeight := validator.ReputationScore
		// Uptime weight (10%)
		uptimeWeight := validator.Uptime

		// Calculate total weight
		weight := (stakeWeight * 0.4) + (perfWeight * 0.3) + (repWeight * 0.2) + (uptimeWeight * 0.1)

		// Ensure weight is positive
		if weight < 0 {
			weight = 0
		}

		weights[addr] = weight
		totalWeight += weight
	}

	// Generate deterministic seed from last block hash and height
	seedData := []byte(lastBlockHash)
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(blockHeight))
	seedData = append(seedData, heightBytes...)
	hash := sha256.Sum256(seedData)

	// Convert hash to float64 in [0, 1)
	// Use first 8 bytes for randomness
	seedInt := binary.BigEndian.Uint64(hash[:8])
	// Max uint64 is ~1.8e19. We normalize to 0.0-1.0
	randomFactor := float64(seedInt) / float64(math.MaxUint64)

	targetValue := randomFactor * totalWeight

	// Select validator
	var currentWeight float64
	for _, addr := range addresses {
		currentWeight += weights[addr]
		if targetValue <= currentWeight {
			return cm.validators[addr], nil
		}
	}

	// Fallback to last validator if rounding errors (should be rare)
	if len(addresses) > 0 {
		return cm.validators[addresses[len(addresses)-1]], nil
	}

	return nil, errors.New("failed to select validator")
}

// UpdateValidatorMetrics updates validator performance metrics
func (cm *ConsensusManager) UpdateValidatorMetrics(validator *Validator, success bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	now := time.Now()

	// Update uptime
	if success {
		validator.Uptime = (validator.Uptime*0.9 + 1.0*0.1) // Moving average
	} else {
		validator.Uptime = validator.Uptime * 0.9 // Decrease uptime on failure
	}

	// Update performance score
	if success {
		validator.PerformanceScore = (validator.PerformanceScore*0.9 + 1.0*0.1)
	} else {
		validator.PerformanceScore = validator.PerformanceScore * 0.9
	}

	// Update last active time
	validator.LastActive = now

	// Store performance history
	cm.performanceHistory[validator.Address] = append(
		cm.performanceHistory[validator.Address],
		validator.PerformanceScore,
	)

	// Trim history if too long
	if len(cm.performanceHistory[validator.Address]) > 100 {
		cm.performanceHistory[validator.Address] = cm.performanceHistory[validator.Address][1:]
	}
}

// SlashValidator handles validator slashing
func (cm *ConsensusManager) SlashValidator(address string, reason string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	log.Printf("⚡ SlashValidator called for address: %s, reason: %s", address[:16]+"...", reason)

	validator, exists := cm.validators[address]
	if !exists {
		log.Printf("❌ Validator not found for slashing: %s", address[:16]+"...")
		return fmt.Errorf("validator not found: %s", address)
	}

	log.Printf("📊 Validator before slashing - SlashCount: %d, ReputationScore: %.3f",
		validator.SlashCount, validator.ReputationScore)

	validator.SlashingHistory = append(validator.SlashingHistory, time.Now())
	validator.ReputationScore *= 0.5

	log.Printf("📉 Validator after slashing - SlashCount: %d, ReputationScore: %.3f",
		validator.SlashCount, validator.ReputationScore)

	if len(validator.SlashingHistory) >= cm.slashingThreshold {
		log.Printf("🗑️ Removing validator due to slashing threshold reached: %s", address[:16]+"...")
		delete(cm.validators, address)
		delete(cm.performanceHistory, address)
		log.Printf("📊 Remaining validators: %d", len(cm.validators))
	}

	log.Printf("✅ Validator %s slashed for reason: %s", address[:16]+"...", reason)
	return nil
}

// RewardValidator handles validator rewards
func (cm *ConsensusManager) RewardValidator(address string, amount uint64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	validator, exists := cm.validators[address]
	if !exists {
		return
	}

	// Record reward in metrics
	validator.RewardHistory = append(validator.RewardHistory, amount)
	validator.TotalRewards += amount
	validator.BlocksProduced++
	validator.LastActive = time.Now()

	// Increase reputation score
	validator.ReputationScore = math.Min(validator.ReputationScore*1.1, 1.0)

	// ACTUAL FIX: Credit the balance in the StateManager
	if cm.blockManager != nil && cm.blockManager.state != nil {
		currentBalance := cm.blockManager.state.GetBalance(address)
		cm.blockManager.state.SetBalance(address, currentBalance+int64(amount))
		log.Printf("💸 RewardValidator: Credited %d to %s (New Balance: %d)", amount, address, currentBalance+int64(amount))
	}
}

// RotateValidators performs periodic validator rotation
func (cm *ConsensusManager) RotateValidators() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Sort validators by performance
	type validatorScore struct {
		address string
		score   float64
	}

	scores := make([]validatorScore, 0)
	for addr, validator := range cm.validators {
		score := (validator.PerformanceScore * 0.4) +
			(validator.ReputationScore * 0.3) +
			(validator.Uptime * 0.3)
		scores = append(scores, validatorScore{addr, score})
	}

	// Sort by score
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Remove bottom 20% of validators
	removeCount := int(float64(len(scores)) * 0.2)
	for i := 0; i < removeCount && i < len(scores); i++ {
		addr := scores[len(scores)-1-i].address
		delete(cm.validators, addr)
		delete(cm.performanceHistory, addr)
	}
}

// getTotalStake returns the total stake across all validators
func (cm *ConsensusManager) getTotalStake() uint64 {
	var total uint64
	for _, validator := range cm.validators {
		total += validator.Stake
	}
	return total
}

// GetValidatorInfo returns information about a validator
func (cm *ConsensusManager) GetValidatorInfo(address string) (*Validator, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	val, exists := cm.validators[address]
	if !exists {
		return nil, fmt.Errorf("validator not found: %s", address)
	}

	return val, nil
}

// GetAllValidators returns all registered validators
func (cm *ConsensusManager) GetAllValidators() []*Validator {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	validators := make([]*Validator, 0, len(cm.validators))
	for _, v := range cm.validators {
		validators = append(validators, v)
	}

	return validators
}

// AddExternalValidator adds a validator received from the network if not already present
func (cm *ConsensusManager) AddExternalValidator(address string, stake uint64) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if _, exists := cm.validators[address]; !exists {
		cm.validators[address] = &Validator{
			Address: address,
			Stake:   stake,
			Active:  true,
		}
	}
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// forgeBlock creates a new block for the given validator.
// 1. Validates that the validator has sufficient stake
// 2. Creates a reward transaction for the validator
// 3. Creates a new block with pending transactions and the reward
// Returns the new block and any error that occurred.
func ForgeBlock(validatorWallet *wallet.Wallet, lastBlock *block.Block, stateManager *StateManager, transactionManager *TransactionManager, consensusManager *ConsensusManager) (*block.Block, error) {
	if lastBlock == nil {
		return nil, errors.New("last block cannot be nil")
	}

	// Check if validator exists and has sufficient stake
	// Use the wallet's address format, not the public key string
	validatorAddress := wallet.PublicKeyToAddress(validatorWallet.PublicKey)
	validatorInfo, err := consensusManager.GetValidatorInfo(validatorAddress)
	if err != nil {
		return nil, fmt.Errorf("%w: validator not found", ErrInvalidValidator)
	}

	if validatorInfo.Stake < MIN_STAKE {
		return nil, fmt.Errorf("%w: insufficient stake", ErrInvalidValidator)
	}

	// Create reward transaction
	shortValidator := wallet.PublicKeyToAddress(validatorWallet.PublicKey)
	rewardTx := transaction.Transaction{
		Type:      transaction.TxTypeRegular,
		Sender:    "network",
		Recipient: shortValidator,
		Amount:    BLOCK_REWARD,
		Fee:       0,
		Timestamp: time.Now().Unix(),
		Nonce:     0, // Network reward transactions use nonce 0
		Data:      "block reward",
		Signature: "NETWORK_REWARD_SIGNATURE", // Add signature for network rewards
	}
	// Get transactions from TransactionManager and filter valid ones
	pendingTxs := transactionManager.GetTransactionsForBlock()
	var transactionsToInclude []transaction.Transaction
	var invalidTxs []transaction.Transaction

	for _, tx := range pendingTxs {
		// 1. Basic format validation
		if err := tx.Validate(); err != nil {
			log.Printf("⚠️ [Forge] Dropping invalid tx (validation failed): %v", err)
			invalidTxs = append(invalidTxs, tx)
			continue
		}

		// 2. State-dependent validation (Balance & Nonce)
		if tx.Sender != "network" {
			senderBalance := stateManager.GetBalance(tx.Sender)
			if senderBalance < tx.Amount+tx.Fee {
				log.Printf("⚠️ [Forge] Dropping tx from %s: Insufficient balance (%d < %d)", tx.Sender[:16], senderBalance, tx.Amount+tx.Fee)
				invalidTxs = append(invalidTxs, tx)
				continue
			}

			stateNonce := stateManager.GetNonce(tx.Sender)
			if tx.Nonce != stateNonce {
				// If nonce is in the future, we keep it in pool but don't include in THIS block
				// If nonce is in the past, we drop it
				if tx.Nonce < stateNonce {
					log.Printf("⚠️ [Forge] Dropping stale tx from %s: Nonce too low (%d < %d)", tx.Sender[:16], tx.Nonce, stateNonce)
					invalidTxs = append(invalidTxs, tx)
				}
				continue
			}
		}

		transactionsToInclude = append(transactionsToInclude, tx)
	}

	// Prune strictly invalid or stale transactions from mempool
	if len(invalidTxs) > 0 {
		transactionManager.RemoveTransactions(invalidTxs)
	}

	transactionsToInclude = append(transactionsToInclude, rewardTx)
	newBlock, err := block.CreateNewBlock(transactionsToInclude, lastBlock, validatorWallet)
	if err != nil {
		return nil, fmt.Errorf("failed to create new block: %v", err)
	}
	fmt.Printf("Validator %s forged block %d and received %d coin(s).\n", shortValidator, newBlock.Index, BLOCK_REWARD)
	return newBlock, nil
}

// getTransactionsFromPool returns a subset of transactions from the pool.
// This helps prevent blocks from becoming too large.
func getTransactionsFromPool(transactionManager *TransactionManager) []transaction.Transaction {
	// For now, we'll take all transactions from the pool
	// In a real implementation, you might want to limit the number of transactions
	// or prioritize certain transactions based on fees
	return transactionManager.GetTransactionsForBlock()
}

// TrackBlockConfirmation tracks a block's confirmation status
func (cm *ConsensusManager) TrackBlockConfirmation(blockHash string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Initialize confirmation count if not exists
	if _, exists := cm.blockConfirmations[blockHash]; !exists {
		cm.blockConfirmations[blockHash] = 0
	}

	// Increment confirmation count
	cm.blockConfirmations[blockHash]++

	// Check if block has reached finality
	if cm.blockConfirmations[blockHash] >= cm.finalityThreshold {
		cm.finalizedBlocks[blockHash] = true
		log.Printf("Block %s has reached finality with %d confirmations",
			blockHash, cm.blockConfirmations[blockHash])
	}
}

// IsBlockFinalized checks if a block has reached finality
func (cm *ConsensusManager) IsBlockFinalized(blockHash string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.finalizedBlocks[blockHash]
}

// GetBlockConfirmations returns the number of confirmations for a block
func (cm *ConsensusManager) GetBlockConfirmations(blockHash string) int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.blockConfirmations[blockHash]
}

// VerifyBlockFinality verifies if a block can be considered final
func (cm *ConsensusManager) VerifyBlockFinality(block *block.Block) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Check if block is already finalized
	if cm.finalizedBlocks[block.Hash] {
		return true
	}

	// Get current confirmation count
	confirmations := cm.blockConfirmations[block.Hash]

	// Check if block has enough confirmations
	if confirmations >= cm.finalityThreshold {
		// Mark block as finalized
		cm.finalizedBlocks[block.Hash] = true
		return true
	}

	return false
}

// CleanupOldConfirmations removes old confirmation data
func (cm *ConsensusManager) CleanupOldConfirmations() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	now := time.Now()
	for blockHash, confirmations := range cm.blockConfirmations {
		// If block is finalized and confirmation window has passed, remove from tracking
		if cm.finalizedBlocks[blockHash] &&
			now.Sub(time.Unix(int64(confirmations), 0)) > cm.confirmationWindow {
			delete(cm.blockConfirmations, blockHash)
			delete(cm.finalizedBlocks, blockHash)
		}
	}
}

// StartFinalityTracking starts the finality tracking routine
func (cm *ConsensusManager) StartFinalityTracking() {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			cm.CleanupOldConfirmations()
		}
	}()
}

// Update block validation to include finality checks
func (cm *ConsensusManager) ValidateBlock(block *block.Block) error {
	// Basic block validation
	if err := cm.validateBlockBasic(block); err != nil {
		return err
	}

	// Check if block's parent is finalized
	if !cm.IsBlockFinalized(block.PrevHash) {
		return fmt.Errorf("parent block %s is not finalized", block.PrevHash)
	}

	// Track block confirmation
	cm.TrackBlockConfirmation(block.Hash)

	return nil
}

// validateBlockBasic performs basic block validation
func (cm *ConsensusManager) validateBlockBasic(block *block.Block) error {
	// Validate block structure
	if block == nil {
		return errors.New("block is nil")
	}

	// Validate block hash
	if block.Hash == "" {
		return errors.New("block hash is empty")
	}

	// Validate timestamp
	if block.Timestamp > time.Now().Unix() {
		return errors.New("block timestamp is in the future")
	}

	return nil
}

// GetFinalityStatus returns the finality status of a block
func (cm *ConsensusManager) GetFinalityStatus(blockHash string) (int, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	confirmations := cm.blockConfirmations[blockHash]
	finalized := cm.finalizedBlocks[blockHash]

	return confirmations, finalized
}

// GetFinalityThreshold returns the number of confirmations required for finality
func (cm *ConsensusManager) GetFinalityThreshold() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.finalityThreshold
}

// UpdateValidatorStake updates the stake amount for a validator
// The difference is moved between the wallet's balance and staked amount
func (cm *ConsensusManager) UpdateValidatorStake(w *wallet.Wallet, newStake uint64) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	address := wallet.PublicKeyToAddress(w.PublicKey)

	validator, exists := cm.validators[address]
	if !exists {
		return fmt.Errorf("validator not found: %s", address)
	}

	currentStake := validator.Stake
	stakeDiff := int64(newStake) - int64(currentStake)

	if stakeDiff > 0 {
		// Increasing stake - check if wallet has enough balance
		if w == nil || !w.CanStake(uint64(stakeDiff)) {
			return fmt.Errorf("insufficient balance for stake increase: have %d, need %d", w.GetBalance(), stakeDiff)
		}
		// Move additional tokens from balance to staked
		if err := w.Stake(uint64(stakeDiff)); err != nil {
			return fmt.Errorf("failed to stake additional tokens: %v", err)
		}
	} else if stakeDiff < 0 {
		// Decreasing stake - move tokens from staked back to balance
		if err := w.Unstake(uint64(-stakeDiff)); err != nil {
			return fmt.Errorf("failed to unstake tokens: %v", err)
		}
	}

	// Update validator's stake
	validator.Stake = newStake
	return nil
}

// OnChainStake registers or increases stake for a validator via on-chain transaction
func (cm *ConsensusManager) OnChainStake(address string, amount uint64) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	minStake := uint64(1)
	if cm.config != nil && cm.config.MinStake > 0 {
		minStake = uint64(cm.config.MinStake)
	}
	if amount < minStake {
		return fmt.Errorf("stake amount too low: %d (min: %d)", amount, minStake)
	}

	validator, exists := cm.validators[address]
	if !exists {
		// Register new validator (no KYC for on-chain)
		cm.validators[address] = &Validator{
			Address:          address,
			Stake:            amount,
			Delegations:      make(map[string]int64),
			LastBlock:        0,
			SlashCount:       0,
			Active:           true,
			PerformanceScore: 1.0,
			Uptime:           1.0,
			ReputationScore:  1.0,
			SlashingHistory:  make([]time.Time, 0),
			RewardHistory:    make([]uint64, 0),
			KYC:              KYCInfo{},
		}
		// Add to shard
		if cm.shardManager != nil {
			shardID := cm.shardManager.GetShardForAddress(address)
			_ = cm.shardManager.AddValidatorToShard(shardID, address)
		}
		cm.performanceHistory[address] = make([]float64, 0)
		return nil
	}
	// Already a validator: increase stake
	validator.Stake += amount
	return nil
}
