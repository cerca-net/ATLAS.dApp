package main

import (
	"atlas-blockchain/internal/blockchain"
	"atlas-blockchain/pkg/config"
	"testing"
)

func TestDeterministicValidatorSelection(t *testing.T) {
	// Setup ConsensusManager
	cfg := config.DefaultConfig()
	// blockManager is nil, which is fine as we are using AddExternalValidator and not accessing state
	cm := blockchain.NewConsensusManager(cfg, nil)

	// Add 3 validators with equal stake
	val1 := "00112233445566778899aabbccddeeff"
	val2 := "aa112233445566778899aabbccddeeff"
	val3 := "bb112233445566778899aabbccddeeff"

	cm.AddExternalValidator(val1, 1000)
	cm.AddExternalValidator(val2, 1000)
	cm.AddExternalValidator(val3, 1000)

	// Test Case 1: Same inputs -> Same output
	lastHash := "000000000000000000000000000000000000000000000000000000abc123"
	height := int64(10)

	v1, err := cm.ChooseValidator(lastHash, height)
	if err != nil {
		t.Fatalf("Failed to choose validator: %v", err)
	}

	v2, err := cm.ChooseValidator(lastHash, height)
	if err != nil {
		t.Fatalf("Failed to choose validator: %v", err)
	}

	if v1.Address != v2.Address {
		t.Errorf("Non-deterministic selection! Got %s then %s", v1.Address, v2.Address)
	} else {
		t.Logf("Deterministic Check Passed: Selected %s for hash %s", v1.Address, lastHash)
	}

	// Test Case 2: Different inputs
	lastHash2 := "000000000000000000000000000000000000000000000000000000abc124"
	v3, err := cm.ChooseValidator(lastHash2, height)
	if err != nil {
		t.Fatalf("Failed to choose validator 2: %v", err)
	}
	t.Logf("Selected %s for hash %s", v3.Address, lastHash2)
}
