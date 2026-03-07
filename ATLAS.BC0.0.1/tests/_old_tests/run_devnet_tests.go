package main

import (
	"atlas-blockchain/pkg/block"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/wallet"
	"fmt"
	"time"
)

// DevnetTestRunner runs all tests in proper order for devnet testing
func DevnetTestRunner() {
	fmt.Println("🚀 Starting Devnet Test Suite")
	fmt.Println("================================")
	fmt.Printf("Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Test Configuration
	config := DefaultConfig()
	fmt.Printf("Test Configuration:\n")
	fmt.Printf("  - Block Time: %v\n", config.BlockTime)
	fmt.Printf("  - Min Stake: %d\n", config.MinStake)
	fmt.Printf("  - Max Block Size: %d\n", config.MaxBlockSize)
	fmt.Printf("  - Max Tx Pool Size: %d\n", config.MaxTxPoolSize)
	fmt.Println()

	// Run test suites in order
	runCoreInfrastructureTests()
	runWalletTests()
	runTransactionTests()
	runBlockTests()
	runConsensusTests()
	runIntegrationTests()

	fmt.Println("✅ Devnet Test Suite Completed!")
}

func runCoreInfrastructureTests() {
	fmt.Println("📋 Core Infrastructure Tests")
	fmt.Println("-----------------------------")

	// Test 1: Configuration
	fmt.Print("  Testing Configuration Validation... ")
	if err := DefaultConfig().Validate(); err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 2: State Manager
	fmt.Print("  Testing State Manager... ")
	stateManager := NewStateManager(DefaultConfig())
	if stateManager == nil {
		fmt.Println("❌ FAILED: State manager is nil")
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 3: Transaction Manager
	fmt.Print("  Testing Transaction Manager... ")
	txManager := NewTransactionManager(DefaultConfig(), stateManager)
	if txManager == nil {
		fmt.Println("❌ FAILED: Transaction manager is nil")
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 4: Block Manager
	fmt.Print("  Testing Block Manager... ")
	blockManager := NewBlockManager(DefaultConfig(), stateManager)
	if blockManager == nil {
		fmt.Println("❌ FAILED: Block manager is nil")
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 5: Consensus Manager
	fmt.Print("  Testing Consensus Manager... ")
	consensusManager := NewConsensusManager(DefaultConfig(), blockManager)
	if consensusManager == nil {
		fmt.Println("❌ FAILED: Consensus manager is nil")
	} else {
		fmt.Println("✅ PASSED")
	}

	fmt.Println()
}

func runWalletTests() {
	fmt.Println("🔐 Wallet & Key Management Tests")
	fmt.Println("--------------------------------")

	// Test 1: Wallet Creation
	fmt.Print("  Testing Wallet Creation... ")
	wallet1, err := wallet.NewWallet()
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 2: Wallet Uniqueness
	fmt.Print("  Testing Wallet Uniqueness... ")
	wallet2, err := wallet.NewWallet()
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else if wallet1.PublicKeyStr() == wallet2.PublicKeyStr() {
		fmt.Println("❌ FAILED: Wallets have same public key")
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 3: Transaction Signing
	fmt.Print("  Testing Transaction Signing... ")
	tx := transaction.Transaction{
		Sender:    wallet1.PublicKeyStr(),
		Recipient: "test_recipient",
		Amount:    100,
		Fee:       1,
		Timestamp: time.Now().Unix(),
		Nonce:     0,
		Data:      "test",
	}
	if err := wallet1.SignTransaction(&tx); err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 4: Signature Verification
	fmt.Print("  Testing Signature Verification... ")
	valid, err := wallet.VerifyTransactionSignature(tx)
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else if !valid {
		fmt.Println("❌ FAILED: Signature verification failed")
	} else {
		fmt.Println("✅ PASSED")
	}

	fmt.Println()
}

func runTransactionTests() {
	fmt.Println("💰 Transaction System Tests")
	fmt.Println("---------------------------")

	// Test 1: Transaction Validation
	fmt.Print("  Testing Transaction Validation... ")
	tx := transaction.Transaction{
		Sender:    "1234567890abcdef",
		Recipient: "abcdef1234567890",
		Amount:    100,
		Fee:       1,
		Timestamp: time.Now().Unix(),
		Nonce:     0,
		Data:      "test",
		// Replace Signature: "valid_signature" with real signing logic or a TODO
		// TODO: In test setup, use wallet.SignBlock to generate a real signature for the block using the test validator's wallet
		Signature: "valid_signature",
	}
	if err := tx.Validate(); err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 2: Transaction Hash
	fmt.Print("  Testing Transaction Hash... ")
	hash1 := wallet.CalculateTxHash(tx)
	hash2 := wallet.CalculateTxHash(tx)
	if string(hash1) != string(hash2) {
		fmt.Println("❌ FAILED: Same transaction should produce same hash")
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 3: Transaction Pool
	fmt.Print("  Testing Transaction Pool... ")
	config := DefaultConfig()
	stateManager := NewStateManager(config)
	txManager := NewTransactionManager(config, stateManager)
	if err := txManager.AddTransaction(tx); err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else if txManager.GetPoolSize() != 1 {
		fmt.Printf("❌ FAILED: Pool size should be 1, got %d\n", txManager.GetPoolSize())
	} else {
		fmt.Println("✅ PASSED")
	}

	fmt.Println()
}

func runBlockTests() {
	fmt.Println("⛓️ Block System Tests")
	fmt.Println("---------------------")

	// Test 1: Genesis Block
	fmt.Print("  Testing Genesis Block... ")
	genesisBlock := block.CreateGenesisBlock()
	if genesisBlock == nil {
		fmt.Println("❌ FAILED: Genesis block is nil")
	} else if genesisBlock.Index != 0 {
		fmt.Printf("❌ FAILED: Genesis index should be 0, got %d\n", genesisBlock.Index)
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 2: Block Creation
	fmt.Print("  Testing Block Creation... ")
	wallet, err := wallet.NewWallet()
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else {
		tx := transaction.Transaction{
			Sender:    wallet.PublicKeyStr(),
			Recipient: "test_recipient",
			Amount:    100,
			Fee:       1,
			Timestamp: time.Now().Unix(),
			Nonce:     0,
			Data:      "test",
		}
		block, err := block.CreateNewBlock([]transaction.Transaction{tx}, genesisBlock, wallet)
		if err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
		} else if block == nil {
			fmt.Println("❌ FAILED: Block is nil")
		} else {
			fmt.Println("✅ PASSED")
		}
	}

	// Test 3: Block Hash
	fmt.Print("  Testing Block Hash... ")
	hash1 := block.CalculateHash(*genesisBlock)
	hash2 := block.CalculateHash(*genesisBlock)
	if hash1 != hash2 {
		fmt.Println("❌ FAILED: Same block should produce same hash")
	} else {
		fmt.Println("✅ PASSED")
	}

	// Test 4: Block Validation
	fmt.Print("  Testing Block Validation... ")
	config := DefaultConfig()
	stateManager := NewStateManager(config)
	blockManager := NewBlockManager(config, stateManager)
	if err := blockManager.AddBlock(genesisBlock); err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else {
		fmt.Println("✅ PASSED")
	}

	fmt.Println()
}

func runConsensusTests() {
	fmt.Println("👥 Consensus & PoS Tests")
	fmt.Println("------------------------")

	// Test 1: Validator Registration
	fmt.Print("  Testing Validator Registration... ")
	config := DefaultConfig()
	stateManager := NewStateManager(config)
	blockManager := NewBlockManager(config, stateManager)
	consensusManager := NewConsensusManager(config, blockManager)

	wallet, err := wallet.NewWallet()
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else {
		// Set initial balance for the wallet
		wallet.SetBalance(1000)

		kyc := KYCInfo{
			FullName: "Test Validator",
			Country:  "Test Country",
			IDNumber: "TEST001",
			Verified: true,
		}
		err = consensusManager.RegisterValidator(wallet, 1000, kyc)
		if err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
		} else {
			validator, err := consensusManager.GetValidatorInfo(wallet.PublicKeyStr())
			if err != nil {
				fmt.Printf("❌ FAILED: %v\n", err)
			} else if validator == nil {
				fmt.Println("❌ FAILED: Validator is nil")
			} else {
				fmt.Println("✅ PASSED")
			}
		}
	}

	// Test 2: Validator Selection
	fmt.Print("  Testing Validator Selection... ")
	validator, err := consensusManager.ChooseValidator()
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else if validator == nil {
		fmt.Println("❌ FAILED: Chosen validator is nil")
	} else {
		fmt.Println("✅ PASSED")
	}

	fmt.Println()
}

func runIntegrationTests() {
	fmt.Println("🔗 Integration Tests")
	fmt.Println("-------------------")

	// Test 1: End-to-End Transaction Flow
	fmt.Print("  Testing End-to-End Transaction Flow... ")
	config := DefaultConfig()
	config.BlockTime = time.Second * 2
	stateManager := NewStateManager(config)
	transactionManager := NewTransactionManager(config, stateManager)
	blockManager := NewBlockManager(config, stateManager)

	walletA, err := wallet.NewWallet()
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else {
		walletB, err := wallet.NewWallet()
		if err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
		} else {
			stateManager.SetBalance(walletA.PublicKeyStr(), 1000)

			tx := transaction.Transaction{
				Sender:    walletA.PublicKeyStr(),
				Recipient: walletB.PublicKeyStr(),
				Amount:    100,
				Fee:       1,
				Timestamp: time.Now().Unix(),
				Nonce:     0,
				Data:      "test transfer",
			}

			if err := walletA.SignTransaction(&tx); err != nil {
				fmt.Printf("❌ FAILED: %v\n", err)
			} else if err := transactionManager.AddTransaction(tx); err != nil {
				fmt.Printf("❌ FAILED: %v\n", err)
			} else {
				transactions := transactionManager.GetTransactionsForBlock()
				lastBlock := blockManager.GetLatestBlock()
				block, err := block.CreateNewBlock(transactions, lastBlock, walletA)
				if err != nil {
					fmt.Printf("❌ FAILED: %v\n", err)
				} else if err := blockManager.AddBlock(block); err != nil {
					fmt.Printf("❌ FAILED: %v\n", err)
				} else {
					fmt.Println("✅ PASSED")
				}
			}
		}
	}

	// Test 2: Multiple Transactions
	fmt.Print("  Testing Multiple Transactions... ")
	walletA, err = wallet.NewWallet()
	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
	} else {
		walletB, err := wallet.NewWallet()
		if err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
		} else {
			stateManager.SetBalance(walletA.PublicKeyStr(), 1000)

			for i := 0; i < 3; i++ {
				tx := transaction.Transaction{
					Sender:    walletA.PublicKeyStr(),
					Recipient: walletB.PublicKeyStr(),
					Amount:    10,
					Fee:       1,
					Timestamp: time.Now().Unix(),
					Nonce:     0,
					Data:      fmt.Sprintf("transfer %d", i),
				}

				if err := walletA.SignTransaction(&tx); err != nil {
					fmt.Printf("❌ FAILED: %v\n", err)
					break
				}

				if err := transactionManager.AddTransaction(tx); err != nil {
					fmt.Printf("❌ FAILED: %v\n", err)
					break
				}
			}

			transactions := transactionManager.GetTransactionsForBlock()
			if len(transactions) == 3 {
				fmt.Println("✅ PASSED")
			} else {
				fmt.Printf("❌ FAILED: Should have 3 transactions, got %d\n", len(transactions))
			}
		}
	}

	fmt.Println()
}
