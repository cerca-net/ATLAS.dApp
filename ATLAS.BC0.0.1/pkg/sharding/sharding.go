package sharding

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// ShardID represents a unique identifier for a shard
type ShardID uint32

// Shard represents a single shard in the network
type Shard struct {
	ID           ShardID
	Validators   map[string]bool // Address -> bool
	Transactions []string         // Transaction hashes
	State        map[string]interface{}
	mu           sync.RWMutex
}

// ShardManager handles shard coordination and cross-shard communication
type ShardManager struct {
	shards       map[ShardID]*Shard
	totalShards  uint32
	coordinator  *ShardCoordinator
	mu           sync.RWMutex
}

// ShardCoordinator manages cross-shard transactions and coordination
type ShardCoordinator struct {
	shardManager *ShardManager
	crossShardTx map[string]*CrossShardTransaction
	mu           sync.RWMutex
}

// CrossShardTransaction represents a transaction that spans multiple shards
type CrossShardTransaction struct {
	ID           string
	SourceShard  ShardID
	TargetShard  ShardID
	Transaction  interface{}
	Status       string // "pending", "committed", "failed"
	Timestamp    int64
	Participants []ShardID
}

// ShardConfig holds configuration for sharding
type ShardConfig struct {
	TotalShards     uint32
	ShardSize       uint32 // Number of validators per shard
	CrossShardDelay time.Duration
	ConsensusType   string // "pbft", "raft", etc.
}

// NewShardManager creates a new shard manager
func NewShardManager(config *ShardConfig) *ShardManager {
	sm := &ShardManager{
		shards:      make(map[ShardID]*Shard),
		totalShards: config.TotalShards,
	}
	
	// Initialize shards
	for i := uint32(0); i < config.TotalShards; i++ {
		shardID := ShardID(i)
		sm.shards[shardID] = &Shard{
			ID:           shardID,
			Validators:   make(map[string]bool),
			Transactions: make([]string, 0),
			State:        make(map[string]interface{}),
		}
	}
	
	// Initialize coordinator
	sm.coordinator = &ShardCoordinator{
		shardManager: sm,
		crossShardTx: make(map[string]*CrossShardTransaction),
	}
	
	return sm
}

// GetShardForAddress determines which shard an address belongs to
func (sm *ShardManager) GetShardForAddress(address string) ShardID {
	hash := sha256.Sum256([]byte(address))
	shardID := ShardID(hash[0]) % ShardID(sm.totalShards)
	return shardID
}

// GetShardForTransaction determines which shard a transaction belongs to
func (sm *ShardManager) GetShardForTransaction(txHash string) ShardID {
	hash := sha256.Sum256([]byte(txHash))
	shardID := ShardID(hash[0]) % ShardID(sm.totalShards)
	return shardID
}

// GetTotalShards returns the total number of shards
func (sm *ShardManager) GetTotalShards() uint32 {
	return sm.totalShards
}

// AddValidatorToShard adds a validator to a specific shard
func (sm *ShardManager) AddValidatorToShard(shardID ShardID, validatorAddress string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	shard, exists := sm.shards[shardID]
	if !exists {
		return fmt.Errorf("shard %d does not exist", shardID)
	}
	
	shard.mu.Lock()
	shard.Validators[validatorAddress] = true
	shard.mu.Unlock()
	
	log.Printf("✅ Added validator %s to shard %d", validatorAddress, shardID)
	return nil
}

// RemoveValidatorFromShard removes a validator from a specific shard
func (sm *ShardManager) RemoveValidatorFromShard(shardID ShardID, validatorAddress string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	shard, exists := sm.shards[shardID]
	if !exists {
		return fmt.Errorf("shard %d does not exist", shardID)
	}
	
	shard.mu.Lock()
	delete(shard.Validators, validatorAddress)
	shard.mu.Unlock()
	
	log.Printf("❌ Removed validator %s from shard %d", validatorAddress, shardID)
	return nil
}

// AddTransactionToShard adds a transaction to the appropriate shard
func (sm *ShardManager) AddTransactionToShard(txHash string, transaction interface{}) error {
	shardID := sm.GetShardForTransaction(txHash)
	
	sm.mu.RLock()
	shard, exists := sm.shards[shardID]
	sm.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("shard %d does not exist", shardID)
	}
	
	shard.mu.Lock()
	shard.Transactions = append(shard.Transactions, txHash)
	shard.State[txHash] = transaction
	shard.mu.Unlock()
	
	log.Printf("📝 Added transaction %s to shard %d", txHash[:16]+"...", shardID)
	return nil
}

// GetShardInfo returns information about a specific shard
func (sm *ShardManager) GetShardInfo(shardID ShardID) (*Shard, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	shard, exists := sm.shards[shardID]
	if !exists {
		return nil, fmt.Errorf("shard %d does not exist", shardID)
	}
	
	return shard, nil
}

// GetAllShards returns information about all shards
func (sm *ShardManager) GetAllShards() map[ShardID]*Shard {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	result := make(map[ShardID]*Shard)
	for id, shard := range sm.shards {
		result[id] = shard
	}
	
	return result
}

// CreateCrossShardTransaction creates a new cross-shard transaction
func (sm *ShardManager) CreateCrossShardTransaction(sourceShard, targetShard ShardID, transaction interface{}) (*CrossShardTransaction, error) {
	txID := fmt.Sprintf("cst_%d_%d_%d", sourceShard, targetShard, time.Now().Unix())
	
	cst := &CrossShardTransaction{
		ID:           txID,
		SourceShard:  sourceShard,
		TargetShard:  targetShard,
		Transaction:  transaction,
		Status:       "pending",
		Timestamp:    time.Now().Unix(),
		Participants: []ShardID{sourceShard, targetShard},
	}
	
	sm.coordinator.mu.Lock()
	sm.coordinator.crossShardTx[txID] = cst
	sm.coordinator.mu.Unlock()
	
	log.Printf("🔗 Created cross-shard transaction %s from shard %d to shard %d", txID, sourceShard, targetShard)
	return cst, nil
}

// ProcessCrossShardTransaction processes a cross-shard transaction
func (sm *ShardManager) ProcessCrossShardTransaction(txID string) error {
	sm.coordinator.mu.RLock()
	cst, exists := sm.coordinator.crossShardTx[txID]
	sm.coordinator.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("cross-shard transaction %s not found", txID)
	}
	
	// Simulate cross-shard processing
	log.Printf("🔄 Processing cross-shard transaction %s", txID)
	
	// Update status to committed
	sm.coordinator.mu.Lock()
	cst.Status = "committed"
	sm.coordinator.mu.Unlock()
	
	log.Printf("✅ Cross-shard transaction %s committed", txID)
	return nil
}

// GetCrossShardTransactions returns all cross-shard transactions
func (sm *ShardManager) GetCrossShardTransactions() []*CrossShardTransaction {
	sm.coordinator.mu.RLock()
	defer sm.coordinator.mu.RUnlock()
	
	var transactions []*CrossShardTransaction
	for _, cst := range sm.coordinator.crossShardTx {
		transactions = append(transactions, cst)
	}
	
	return transactions
}

// ShardConsensus represents consensus state for a shard
type ShardConsensus struct {
	ShardID      ShardID
	BlockHeight  uint64
	Validators   []string
	ConsensusType string
	mu           sync.RWMutex
}

// NewShardConsensus creates a new consensus instance for a shard
func NewShardConsensus(shardID ShardID, consensusType string) *ShardConsensus {
	return &ShardConsensus{
		ShardID:       shardID,
		BlockHeight:   0,
		Validators:    make([]string, 0),
		ConsensusType: consensusType,
	}
}

// ProposeBlock proposes a new block for the shard
func (sc *ShardConsensus) ProposeBlock(block interface{}) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	
	// Simulate block proposal
	log.Printf("📦 Proposing block for shard %d at height %d", sc.ShardID, sc.BlockHeight+1)
	
	// Update block height
	sc.BlockHeight++
	
	log.Printf("✅ Block proposed for shard %d, new height: %d", sc.ShardID, sc.BlockHeight)
	return nil
}

// GetShardStatistics returns statistics for all shards
func (sm *ShardManager) GetShardStatistics() map[string]interface{} {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	stats := make(map[string]interface{})
	
	for shardID, shard := range sm.shards {
		shard.mu.RLock()
		validatorCount := len(shard.Validators)
		txCount := len(shard.Transactions)
		shard.mu.RUnlock()
		
		stats[fmt.Sprintf("shard_%d", shardID)] = map[string]interface{}{
			"id":              shardID,
			"validator_count": validatorCount,
			"tx_count":        txCount,
			"active":          validatorCount > 0,
		}
	}
	
	// Add cross-shard transaction stats
	if sm.coordinator != nil {
		sm.coordinator.mu.RLock()
		crossShardCount := len(sm.coordinator.crossShardTx)
		sm.coordinator.mu.RUnlock()
		stats["cross_shard_tx_count"] = crossShardCount
	} else {
		stats["cross_shard_tx_count"] = 0
	}
	stats["total_shards"] = sm.totalShards
	
	return stats
}

// ValidateShardAssignment validates if a validator is properly assigned to a shard
func (sm *ShardManager) ValidateShardAssignment(validatorAddress string, expectedShardID ShardID) bool {
	actualShardID := sm.GetShardForAddress(validatorAddress)
	return actualShardID == expectedShardID
}

// RebalanceShards rebalances validators across shards
func (sm *ShardManager) RebalanceShards() error {
	log.Printf("🔄 Starting shard rebalancing...")
	
	// Collect all validators
	allValidators := make(map[string]ShardID)
	for shardID, shard := range sm.shards {
		shard.mu.RLock()
		for validator := range shard.Validators {
			allValidators[validator] = shardID
		}
		shard.mu.RUnlock()
	}
	
	// Reassign validators based on address hash
	for validator, oldShardID := range allValidators {
		newShardID := sm.GetShardForAddress(validator)
		
		if newShardID != oldShardID {
			// Remove from old shard
			sm.RemoveValidatorFromShard(oldShardID, validator)
			
			// Add to new shard
			sm.AddValidatorToShard(newShardID, validator)
			
			log.Printf("🔄 Moved validator %s from shard %d to shard %d", validator, oldShardID, newShardID)
		}
	}
	
	log.Printf("✅ Shard rebalancing completed")
	return nil
}

// ExportShardState exports the state of all shards
func (sm *ShardManager) ExportShardState() ([]byte, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	state := make(map[string]interface{})
	
	for shardID, shard := range sm.shards {
		shard.mu.RLock()
		shardState := map[string]interface{}{
			"id":           shard.ID,
			"validators":   shard.Validators,
			"transactions": shard.Transactions,
			"state":        shard.State,
		}
		shard.mu.RUnlock()
		
		state[fmt.Sprintf("shard_%d", shardID)] = shardState
	}
	
	// Add coordinator state
	sm.coordinator.mu.RLock()
	coordinatorState := make(map[string]*CrossShardTransaction)
	for id, cst := range sm.coordinator.crossShardTx {
		coordinatorState[id] = cst
	}
	sm.coordinator.mu.RUnlock()
	
	state["coordinator"] = coordinatorState
	state["total_shards"] = sm.totalShards
	
	return json.MarshalIndent(state, "", "  ")
}

// ImportShardState imports the state of all shards
func (sm *ShardManager) ImportShardState(data []byte) error {
	var state map[string]interface{}
	if err := json.Unmarshal(data, &state); err != nil {
		return fmt.Errorf("failed to unmarshal shard state: %v", err)
	}
	
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	// Import shard states
	for key, value := range state {
		if key == "coordinator" || key == "total_shards" {
			continue
		}
		
		// Parse shard ID from key
		var shardID ShardID
		if _, err := fmt.Sscanf(key, "shard_%d", &shardID); err != nil {
			continue
		}
		
		shardData, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		
		shard := sm.shards[shardID]
		if shard == nil {
			continue
		}
		
		shard.mu.Lock()
		
		// Import validators
		if validatorsData, ok := shardData["validators"].(map[string]interface{}); ok {
			shard.Validators = make(map[string]bool)
			for validator := range validatorsData {
				shard.Validators[validator] = true
			}
		}
		
		// Import transactions
		if txData, ok := shardData["transactions"].([]interface{}); ok {
			shard.Transactions = make([]string, 0)
			for _, tx := range txData {
				if txStr, ok := tx.(string); ok {
					shard.Transactions = append(shard.Transactions, txStr)
				}
			}
		}
		
		// Import state
		if stateData, ok := shardData["state"].(map[string]interface{}); ok {
			shard.State = stateData
		}
		
		shard.mu.Unlock()
	}
	
	log.Printf("✅ Shard state imported successfully")
	return nil
} 