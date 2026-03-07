package block

import (
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/wallet"
	"testing"
)

func TestBlockSigningAndVerification(t *testing.T) {
	// Test 1: Create wallet and verify it has a private key
	t.Run("Wallet Creation", func(t *testing.T) {
		wallet, err := wallet.NewWallet()
		if err != nil {
			t.Fatalf("Failed to create wallet: %v", err)
		}

		if wallet.PrivateKey == nil {
			t.Fatal("Wallet private key is nil")
		}

		if len(wallet.PublicKey) == 0 {
			t.Fatal("Wallet public key is empty")
		}
	})

	// Test 2: Create genesis block
	t.Run("Genesis Block Creation", func(t *testing.T) {
		genesis := CreateGenesisBlock()

		if genesis.Index != 0 {
			t.Errorf("Genesis block index should be 0, got %d", genesis.Index)
		}

		if genesis.PrevHash != "0" {
			t.Errorf("Genesis block prev hash should be '0', got %s", genesis.PrevHash)
		}

		if genesis.Validator != "GENESIS_VALIDATOR" {
			t.Errorf("Genesis block validator should be 'GENESIS_VALIDATOR', got %s", genesis.Validator)
		}

		if genesis.Signature != "GENESIS_SIGNATURE" {
			t.Errorf("Genesis block signature should be 'GENESIS_SIGNATURE', got %s", genesis.Signature)
		}
	})

	// Test 3: Create new block with real signing
	t.Run("Block Creation and Signing", func(t *testing.T) {
		wallet, err := wallet.NewWallet()
		if err != nil {
			t.Fatalf("Failed to create wallet: %v", err)
		}

		genesis := CreateGenesisBlock()

		// Create new block with empty transaction list to avoid validation issues
		newBlock, err := CreateNewBlock([]transaction.Transaction{}, genesis, wallet)
		if err != nil {
			t.Fatalf("Failed to create new block: %v", err)
		}

		// Verify block properties
		if newBlock.Index != 1 {
			t.Errorf("New block index should be 1, got %d", newBlock.Index)
		}

		if newBlock.PrevHash != genesis.Hash {
			t.Errorf("New block prev hash should match genesis hash")
		}

		if newBlock.Validator != wallet.PublicKeyStr() {
			t.Errorf("New block validator should match wallet public key")
		}

		// Verify block has a signature
		if newBlock.Signature == "" {
			t.Fatal("New block signature is empty")
		}

		// Verify block signature
		valid, err := VerifyBlockSignature(newBlock, wallet.PublicKey)
		if err != nil {
			t.Fatalf("Block signature verification error: %v", err)
		}

		if !valid {
			t.Fatal("Block signature verification failed")
		}
	})

	// Test 4: Test signature tampering detection
	t.Run("Signature Tampering Detection", func(t *testing.T) {
		wallet, err := wallet.NewWallet()
		if err != nil {
			t.Fatalf("Failed to create wallet: %v", err)
		}

		genesis := CreateGenesisBlock()

		// Create new block
		newBlock, err := CreateNewBlock([]transaction.Transaction{}, genesis, wallet)
		if err != nil {
			t.Fatalf("Failed to create new block: %v", err)
		}

		// Tamper with the signature
		originalSignature := newBlock.Signature
		newBlock.Signature = "TAMPERED_SIGNATURE"

		// Verify that tampered signature fails
		valid, err := VerifyBlockSignature(newBlock, wallet.PublicKey)
		if err == nil && valid {
			t.Fatal("Tampered signature should fail verification")
		}

		// Restore original signature
		newBlock.Signature = originalSignature

		// Verify that original signature still works
		valid, err = VerifyBlockSignature(newBlock, wallet.PublicKey)
		if err != nil {
			t.Fatalf("Original signature verification error: %v", err)
		}

		if !valid {
			t.Fatal("Original signature verification failed after restoration")
		}
	})

	// Test 5: Test chain validation
	t.Run("Chain Validation", func(t *testing.T) {
		wallet, err := wallet.NewWallet()
		if err != nil {
			t.Fatalf("Failed to create wallet: %v", err)
		}

		genesis := CreateGenesisBlock()

		// Create a few blocks
		block1, err := CreateNewBlock([]transaction.Transaction{}, genesis, wallet)
		if err != nil {
			t.Fatalf("Failed to create block 1: %v", err)
		}

		block2, err := CreateNewBlock([]transaction.Transaction{}, block1, wallet)
		if err != nil {
			t.Fatalf("Failed to create block 2: %v", err)
		}

		// Validate the chain
		chain := []*Block{genesis, block1, block2}
		if !ValidateChain(chain) {
			t.Fatal("Valid chain validation failed")
		}

		// Test invalid chain (tampered block)
		block1.Hash = "TAMPERED_HASH"
		if ValidateChain(chain) {
			t.Fatal("Invalid chain validation should fail")
		}
	})
}
