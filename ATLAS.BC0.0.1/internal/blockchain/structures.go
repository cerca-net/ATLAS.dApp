package blockchain

import (
	"fmt"
	"time"
	"atlas-blockchain/pkg/wallet"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/block"
)

// Blockchain is the full chain of validated blocks.
// It is stored as a slice of block pointers for efficient memory usage.
var Blockchain []*block.Block

// createNewBlock creates a new block to be added to the chain.
// It validates all transactions before creating the block.
func createNewBlock(transactions []transaction.Transaction, prevBlock *block.Block, validatorWallet *wallet.Wallet) (*block.Block, error) {
	// Validate all transactions
	for _, tx := range transactions {
		if err := tx.Validate(); err != nil {
			return nil, fmt.Errorf("transaction validation failed: %v", err)
		}
		if tx.Sender != "network" {
			valid, err := wallet.VerifyTransactionSignature(tx)
			if err != nil {
				return nil, fmt.Errorf("signature verification error: %v", err)
			}
			if !valid {
				return nil, fmt.Errorf("invalid signature for transaction from %s", tx.Sender)
			}
		}
	}

	validatorPubKeyHex := validatorWallet.PublicKeyStr()
	newBlock := &block.Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevBlock.Hash,
		Validator:    validatorPubKeyHex,
		Signature:    "", // Will be set after signing
	}
	// Sign the block
	sig, err := block.SignBlock(newBlock, validatorWallet)
	if err != nil {
		return nil, fmt.Errorf("failed to sign block: %v", err)
	}
	newBlock.Signature = sig
	newBlock.Hash = block.CalculateHash(*newBlock)
	return newBlock, nil
}

// createGenesisBlock creates the first block in the chain.
// The genesis block has no transactions and a special previous hash.
func createGenesisBlock() *block.Block {
	genesisBlock := &block.Block{
		Index:        0,
		Timestamp:    1640995200, // Fixed timestamp: January 1, 2022 00:00:00 UTC
		Transactions: []transaction.Transaction{},
		PrevHash:     "0",
		Validator:    "GENESIS_VALIDATOR",
		Signature:    "GENESIS_SIGNATURE",
	}
	genesisBlock.Hash = block.CalculateHash(*genesisBlock)
	return genesisBlock
}
