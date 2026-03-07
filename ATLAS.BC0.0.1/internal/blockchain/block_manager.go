package blockchain

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"atlas-blockchain/pkg/wallet"
	"encoding/hex"
	"atlas-blockchain/pkg/block"
	"atlas-blockchain/pkg/config"
)

// BlockManager handles block operations and state management
type BlockManager struct {
	chain     []*block.Block
	mu        sync.RWMutex
	config    *config.BlockchainConfig
	state     *StateManager
	onBlockAdded func(*block.Block) // Callback function called after block is added
}

// NewBlockManager creates a new block manager
func NewBlockManager(config *config.BlockchainConfig, state *StateManager) *BlockManager {
	return &BlockManager{
		chain:  []*block.Block{createGenesisBlock()},
		config: config,
		state:  state,
	}
}

// SetOnBlockAddedCallback sets the callback function to be called after a block is added
func (bm *BlockManager) SetOnBlockAddedCallback(callback func(*block.Block)) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.onBlockAdded = callback
}

// AddBlock adds a new block to the chain
func (bm *BlockManager) AddBlock(blk *block.Block) error {
	log.Printf("🔧 AddBlock: Starting to add block %d", blk.Index)
	bm.mu.Lock()
	defer bm.mu.Unlock()

	// Validate block
	log.Printf("🔍 AddBlock: Validating block...")
	if err := bm.validateBlock(blk); err != nil {
		log.Printf("❌ AddBlock: Block validation failed: %v", err)
		return fmt.Errorf("invalid block: %v", err)
	}
	log.Printf("✅ AddBlock: Block validation passed")

	// Update state with block transactions
	log.Printf("📊 AddBlock: Updating state with %d transactions...", len(blk.Transactions))
	if err := bm.state.updateState(blk); err != nil {
		log.Printf("❌ AddBlock: State update failed: %v", err)
		return fmt.Errorf("failed to update state: %v", err)
	}
	log.Printf("✅ AddBlock: State update completed")

	// Add block to chain
	log.Printf("⛓️ AddBlock: Adding block to chain...")
	bm.chain = append(bm.chain, blk)
	log.Printf("✅ AddBlock: Block added to chain. Chain length: %d", len(bm.chain))

	// Check if pruning is needed
	if len(bm.chain) > bm.config.MaxBlockSize*2 {
		log.Printf("🧹 AddBlock: Pruning old blocks...")
		bm.pruneOldBlocks()
	}

	if bm.onBlockAdded != nil {
		log.Printf("🔄 AddBlock: Calling onBlockAdded callback...")
		bm.onBlockAdded(blk)
		log.Printf("✅ AddBlock: onBlockAdded callback completed")
	}

	log.Printf("✅ AddBlock: Successfully added block %d to chain", blk.Index)
	return nil
}

// validateBlock validates a block before adding it to the chain
func (bm *BlockManager) validateBlock(blk *block.Block) error {
	if blk == nil {
		return errors.New("block cannot be nil")
	}

	// Check block size
	if len(blk.Transactions) > bm.config.MaxBlockSize {
		return fmt.Errorf("block exceeds maximum size of %d transactions", bm.config.MaxBlockSize)
	}

	// Verify block index
	lastBlock := bm.chain[len(bm.chain)-1]
	if blk.Index != lastBlock.Index+1 {
		return fmt.Errorf("invalid block index: expected %d, got %d", lastBlock.Index+1, blk.Index)
	}

	// Verify previous hash
	if blk.PrevHash != lastBlock.Hash {
		return errors.New("invalid previous hash")
	}

	// Verify block hash
	if blk.Hash != block.CalculateHash(*blk) {
		return errors.New("invalid block hash")
	}

	// Verify block signature
	validatorPubKeyBytes, err := hex.DecodeString(blk.Validator)
	if err != nil {
		return fmt.Errorf("invalid validator public key encoding: %v", err)
	}
	valid, err := block.VerifyBlockSignature(blk, validatorPubKeyBytes)
	if err != nil {
		return fmt.Errorf("block signature verification error: %v", err)
	}
	if !valid {
		return errors.New("invalid block signature")
	}

	// Verify all transactions
	for _, tx := range blk.Transactions {
		if tx.Sender != "network" { // Skip network reward transactions
			valid, err := wallet.VerifyTransactionSignature(tx)
			if err != nil || !valid {
				return fmt.Errorf("invalid transaction: %v", err)
			}
		}
	}

	return nil
}

// pruneOldBlocks removes old blocks while maintaining state
func (bm *BlockManager) pruneOldBlocks() {
	// Keep the last N blocks
	keepBlocks := bm.config.MaxBlockSize
	if len(bm.chain) <= keepBlocks {
		return
	}

	// Create new chain with only the last N blocks
	newChain := make([]*block.Block, keepBlocks)
	copy(newChain, bm.chain[len(bm.chain)-keepBlocks:])
	bm.chain = newChain
}

// GetBlockByIndex retrieves a block by its index
func (bm *BlockManager) GetBlockByIndex(index int) (*block.Block, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if index < 0 || index >= len(bm.chain) {
		return nil, fmt.Errorf("block index %d out of range", index)
	}

	return bm.chain[index], nil
}

// GetLatestBlock returns the most recent block
func (bm *BlockManager) GetLatestBlock() *block.Block {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	if len(bm.chain) == 0 {
		return nil
	}

	return bm.chain[len(bm.chain)-1]
}

// GetChainLength returns the current length of the blockchain
func (bm *BlockManager) GetChainLength() int {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	return len(bm.chain)
}

// GetBlockByHash retrieves a block by its hash
func (bm *BlockManager) GetBlockByHash(hash string) *block.Block {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	for _, blk := range bm.chain {
		if blk.Hash == hash {
			return blk
		}
	}
	return nil
}

// GetBlockHeight returns the current block height (index of the latest block)
func (bm *BlockManager) GetBlockHeight() int {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	if len(bm.chain) == 0 {
		return 0
	}
	return bm.chain[len(bm.chain)-1].Index
}

// GetBlocks returns a slice of blocks with limit and offset
func (bm *BlockManager) GetBlocks(limit, offset int) []*block.Block {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	if offset > len(bm.chain) {
		return []*block.Block{}
	}
	end := offset + limit
	if end > len(bm.chain) {
		end = len(bm.chain)
	}
	return bm.chain[offset:end]
} 