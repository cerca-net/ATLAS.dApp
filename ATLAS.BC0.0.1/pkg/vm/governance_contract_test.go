package vm

import (
	"testing"
)

func TestGovernanceContract(t *testing.T) {
	govContract := CreateSystemGovernanceContract()
	helper := NewGovernanceContractHelper(govContract)

	stateAdapter := NewMockStateAdapter()
	ctx := &ExecutionContext{
		State:     stateAdapter,
		Timestamp: 1000,
	}

	stateAdapter.balances["0xPROPOSER"] = 50000
	stateAdapter.balances["0xPOOR"] = 50

	t.Run("CreateProposal_Success", func(t *testing.T) {
		propID, err := helper.CreateProposal(ctx, "0xPROPOSER", "Change min stake to 2000", 3600)
		if err != nil {
			t.Fatalf("Failed to create proposal: %v", err)
		}

		if propID != "PROP_1" {
			t.Errorf("Expected PROP_1, got %s", propID)
		}

		info := helper.GetProposal(ctx, propID)
		if info == nil {
			t.Fatalf("Proposal not found")
		}
		if info.Status != ProposalStatusActive {
			t.Errorf("Expected Active status (%d), got %d", ProposalStatusActive, info.Status)
		}
		if info.EndTime != 4600 {
			t.Errorf("Expected end time 4600, got %d", info.EndTime)
		}
	})

	t.Run("CreateProposal_InsufficientFunds", func(t *testing.T) {
		_, err := helper.CreateProposal(ctx, "0xPOOR", "Should fail due to stake", 3600)
		if err == nil {
			t.Fatalf("Expected error for insufficient funds")
		}
	})

	t.Run("CastVote_Success", func(t *testing.T) {
		// ALICE supports, weight 100
		err := helper.CastVote(ctx, "0xALICE", "PROP_1", true, 100)
		if err != nil {
			t.Fatalf("Failed to cast vote: %v", err)
		}

		// BOB opposes, weight 50
		err = helper.CastVote(ctx, "0xBOB", "PROP_1", false, 50)
		if err != nil {
			t.Fatalf("Failed to cast vote: %v", err)
		}

		info := helper.GetProposal(ctx, "PROP_1")
		if info.VotesFor != 100 || info.VotesAgainst != 50 {
			t.Errorf("Vote tally mismatch. For: %d, Against: %d", info.VotesFor, info.VotesAgainst)
		}
	})

	t.Run("CastVote_DoubleVoting_Fails", func(t *testing.T) {
		err := helper.CastVote(ctx, "0xALICE", "PROP_1", true, 50)
		if err == nil {
			t.Fatalf("Expected error for double voting")
		}
	})

	t.Run("ExecuteProposal_TooEarly", func(t *testing.T) {
		err := helper.ExecuteProposal(ctx, "PROP_1")
		if err == nil {
			t.Fatalf("Expected error trying to execute before voting ends")
		}
	})

	t.Run("ExecuteProposal_Success_Passes", func(t *testing.T) {
		// Fast-forward time
		ctx.Timestamp = 5000

		err := helper.ExecuteProposal(ctx, "PROP_1")
		if err != nil {
			t.Fatalf("Failed to execute proposal: %v", err)
		}

		info := helper.GetProposal(ctx, "PROP_1")
		if info.Status != ProposalStatusPassed {
			t.Errorf("Expected status %d (Passed), got %d", ProposalStatusPassed, info.Status)
		}
	})

	t.Run("ExecuteProposal_Success_Fails", func(t *testing.T) {
		ctx.Timestamp = 1000 // Reset
		propID, _ := helper.CreateProposal(ctx, "0xPROPOSER", "Change fee to 0%", 3600)

		// 0xBOB opposes heavily
		_ = helper.CastVote(ctx, "0xBOB", propID, false, 5000)

		ctx.Timestamp = 5000 // Fast-forward time
		err := helper.ExecuteProposal(ctx, propID)
		if err != nil {
			t.Fatalf("Failed to execute proposal: %v", err)
		}

		info := helper.GetProposal(ctx, propID)
		if info.Status != ProposalStatusFailed {
			t.Errorf("Expected status %d (Failed), got %d", ProposalStatusFailed, info.Status)
		}
	})
}
