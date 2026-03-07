package tests

import (
	"atlas-blockchain/internal/blockchain"
	"atlas-blockchain/internal/identity"
	"atlas-blockchain/internal/social"
	"atlas-blockchain/pkg/config"
	"atlas-blockchain/pkg/database"
	"testing"
)

func TestSocialFeatures(t *testing.T) {
	// Initialize components
	cfg := config.DefaultConfig()
	db, _ := database.NewDatabase(cfg.DataDir)
	stateManager := blockchain.NewStateManager(cfg)
	identityManager := identity.NewIdentityManager()

	sm := social.NewSocialManager(identityManager, db, stateManager)

	// Register user
	addr := "0xabc123"
	identityManager.CreateIdentity(addr, nil, "TestUser", "test@test.com")

	// Create a post
	post, err := sm.CreatePost(addr, "Hello World", nil, "public", "post")
	if err != nil {
		t.Fatalf("Failed to create post: %v", err)
	}

	if post.Author != addr {
		t.Errorf("Expected post author %s, got %s", addr, post.Author)
	}

	// Test energize
	// Setup user balance
	stateManager.SetBalance(addr, 100)

	energyState, err := sm.EnergizeObject(post.ID, addr, 10)
	if err != nil {
		t.Fatalf("Failed to energize object: %v", err)
	}

	if energyState.TipBalance != 110 {
		t.Errorf("Expected tip balance 110, got %d", energyState.TipBalance)
	}

	// Verify balance deduction
	newBal := stateManager.GetBalance(addr)
	if newBal != 90 {
		t.Errorf("Expected balance 90, got %d", newBal)
	}
}
