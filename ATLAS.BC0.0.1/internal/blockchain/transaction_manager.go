package blockchain

import (
	"atlas-blockchain/pkg/config"
	"atlas-blockchain/pkg/crypto/zk"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/wallet"
	"container/heap"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

// TransactionPriority represents a transaction with its priority score
type TransactionPriority struct {
	Transaction transaction.Transaction
	Priority    float64
	Timestamp   int64
	Fee         int64
	index       int
	Hash        string // Cache the transaction hash
}

// TransactionHeap implements heap.Interface for transaction prioritization
type TransactionHeap []*TransactionPriority

func (h TransactionHeap) Len() int { return len(h) }

func (h TransactionHeap) Less(i, j int) bool {
	return h[i].Priority > h[j].Priority
}

func (h TransactionHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *TransactionHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*TransactionPriority)
	item.index = n
	*h = append(*h, item)
}

func (h *TransactionHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*h = old[0 : n-1]
	return item
}

// TransactionManager handles transaction processing and pool management
type TransactionManager struct {
	pool                  TransactionHeap
	mu                    sync.RWMutex
	byHash                map[string]*TransactionPriority
	config                *config.BlockchainConfig
	stateManager          *StateManager
	senderReputation      map[string]float64
	transactionComplexity map[string]float64
	historicalSuccessRate map[string]float64

	dynamicFeeMultiplier float64
	zkVerifier           *zk.ProofVerifier
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(config *config.BlockchainConfig, stateManager *StateManager) *TransactionManager {
	tm := &TransactionManager{
		pool:                  make(TransactionHeap, 0),
		byHash:                make(map[string]*TransactionPriority),
		config:                config,
		stateManager:          stateManager,
		senderReputation:      make(map[string]float64),
		transactionComplexity: make(map[string]float64),
		historicalSuccessRate: make(map[string]float64),
		dynamicFeeMultiplier:  1.0,
		zkVerifier:            zk.NewProofVerifier(true),
	}
	heap.Init(&tm.pool)
	return tm
}

// Helper to calculate hash
func (tm *TransactionManager) calculateHash(tx transaction.Transaction) string {
	h := wallet.CalculateTxHash(tx)
	return hex.EncodeToString(h)
}

// GetNextNonce calculates the next valid nonce for a sender
func (tm *TransactionManager) GetNextNonce(sender string) uint64 {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	nonce := tm.stateManager.GetNonce(sender)
	maxPoolNonce := nonce
	foundInPool := false

	for _, item := range tm.pool {
		if item.Transaction.Sender == sender {
			if item.Transaction.Nonce >= maxPoolNonce {
				maxPoolNonce = item.Transaction.Nonce
				foundInPool = true
			}
		}
	}

	if foundInPool && maxPoolNonce >= nonce {
		return maxPoolNonce + 1
	}

	return nonce
}

// AddTransaction adds a transaction to the pool
func (tm *TransactionManager) AddTransaction(tx transaction.Transaction) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.UpdateDynamicFeeMultiplier()

	txHash := tm.calculateHash(tx)

	if _, exists := tm.byHash[txHash]; exists {
		return fmt.Errorf("transaction already in pool")
	}

	if tx.Sender != "network" && tm.stateManager != nil {
		stateNonce := tm.stateManager.GetNonce(tx.Sender)
		if tx.Nonce < stateNonce {
			return fmt.Errorf("nonce too low: got %d, state %d", tx.Nonce, stateNonce)
		}
	}

	if tx.Sender != "network" {
		valid, err := wallet.VerifyTransactionSignature(tx)
		if err != nil {
			return fmt.Errorf("signature verification error: %v", err)
		}
		if !valid {
			return errors.New("invalid transaction signature")
		}
	}

	if tx.Type == transaction.TxTypeZK && tx.ZKProof != nil {
		valid, err := tm.zkVerifier.VerifyProof(tx.ZKProof)
		if err != nil {
			return fmt.Errorf("ZK proof verification error: %v", err)
		}
		if !valid {
			return errors.New("invalid ZK proof")
		}
	}

	if len(tm.pool) >= tm.config.MaxTxPoolSize {
		if len(tm.pool) > 0 {
			lowestPriority := tm.pool[0]
			newPriority := tm.calculatePriority(tx)
			if newPriority > lowestPriority.Priority {
				heap.Pop(&tm.pool)
				delete(tm.byHash, lowestPriority.Hash)
			} else {
				return errors.New("mempool full and transaction priority too low")
			}
		}
	}

	priority := tm.calculatePriority(tx)
	item := &TransactionPriority{
		Transaction: tx,
		Priority:    priority,
		Timestamp:   time.Now().Unix(),
		Fee:         int64(tx.Fee),
		Hash:        txHash,
	}

	heap.Push(&tm.pool, item)
	tm.byHash[txHash] = item
	log.Printf("[Mempool] Transaction added: %s (Nonce: %d)", txHash, tx.Nonce)

	return nil
}

// GetTransactions returns top priority transactions up to limit
// Uses a copy to sort and retrieve without modifying the heap
func (tm *TransactionManager) GetTransactions(limit int) []transaction.Transaction {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	n := len(tm.pool)
	if limit > n {
		limit = n
	}

	result := make([]transaction.Transaction, 0, limit)

	// Copy pointers to avoid modifying original heap structure
	items := make([]*TransactionPriority, len(tm.pool))
	copy(items, tm.pool)

	// Sort by priority (descending)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Priority > items[j].Priority
	})

	for i := 0; i < limit; i++ {
		result = append(result, items[i].Transaction)
	}

	return result
}

// GetTransactionsForBlock returns transactions for block creation
// Defaults to retrieving up to 2000 transactions (or configured block size logic)
func (tm *TransactionManager) GetTransactionsForBlock() []transaction.Transaction {
	return tm.GetTransactions(2000)
}

// GetTransactionByHash returns a transaction by its hash
func (tm *TransactionManager) GetTransactionByHash(hash string) *transaction.Transaction {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if item, ok := tm.byHash[hash]; ok {
		tx := item.Transaction
		return &tx
	}
	return nil
}

// HasTransaction checks if a transaction is in the pool
func (tm *TransactionManager) HasTransaction(hash string) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	_, exists := tm.byHash[hash]
	return exists
}

// RemoveExpiredTransactions removes transactions older than a certain threshold
func (tm *TransactionManager) RemoveExpiredTransactions() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	threshold := time.Now().Add(-24 * time.Hour).Unix()

	idsToRemove := make(map[string]bool)
	for _, item := range tm.pool {
		if item.Timestamp < threshold {
			idsToRemove[item.Hash] = true
			delete(tm.byHash, item.Hash)
		}
	}

	if len(idsToRemove) == 0 {
		return
	}

	newPool := make(TransactionHeap, 0, len(tm.pool))
	for _, item := range tm.pool {
		if !idsToRemove[item.Hash] {
			newPool = append(newPool, item)
		}
	}

	tm.pool = newPool
	heap.Init(&tm.pool)
	log.Printf("[Mempool] Removed %d expired transactions", len(idsToRemove))
}

// GetAllTransactions returns all transactions in the pool
func (tm *TransactionManager) GetAllTransactions() []transaction.Transaction {
	return tm.GetTransactions(len(tm.pool))
}

// GetDynamicFeeMultiplier returns the current fee multiplier
func (tm *TransactionManager) GetDynamicFeeMultiplier() float64 {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.dynamicFeeMultiplier
}

// GetPoolSize returns current pool size
func (tm *TransactionManager) GetPoolSize() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.pool)
}

// UpdateDynamicFeeMultiplier updates the fee multiplier based on congestion
func (tm *TransactionManager) UpdateDynamicFeeMultiplier() {
	usage := float64(len(tm.pool)) / float64(tm.config.MaxTxPoolSize)
	if usage > 0.8 {
		tm.dynamicFeeMultiplier *= 1.2
		if tm.dynamicFeeMultiplier > 10.0 {
			tm.dynamicFeeMultiplier = 10.0
		}
	} else if usage < 0.2 {
		tm.dynamicFeeMultiplier *= 0.8
		if tm.dynamicFeeMultiplier < 1.0 {
			tm.dynamicFeeMultiplier = 1.0
		}
	}
}

// RemoveTransactions removes processed transactions from the pool
func (tm *TransactionManager) RemoveTransactions(txs []transaction.Transaction) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	idsToRemove := make(map[string]bool)
	for _, tx := range txs {
		h := tm.calculateHash(tx)
		idsToRemove[h] = true
		delete(tm.byHash, h)
	}

	newPool := make(TransactionHeap, 0, len(tm.pool))
	for _, item := range tm.pool {
		if !idsToRemove[item.Hash] {
			newPool = append(newPool, item)
		}
	}

	tm.pool = newPool
	heap.Init(&tm.pool)
}

func (tm *TransactionManager) calculatePriority(tx transaction.Transaction) float64 {
	feeScore := float64(tx.Fee)
	repScore := tm.senderReputation[tx.Sender]
	if repScore == 0 {
		repScore = 1.0
	}
	return feeScore * repScore
}
