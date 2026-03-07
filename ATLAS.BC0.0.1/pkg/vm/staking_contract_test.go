package vm

import (
	"testing"
)

func TestStakingContract(t *testing.T) {
	t.Run("CreateStakingContract", func(t *testing.T) {
		contract := CreateSystemStakingContract()
		if contract == nil {
			t.Fatal("Failed to create staking contract")
		}
		if contract.Address != StakingContractAddress {
			t.Errorf("Expected address %s, got %s", StakingContractAddress, contract.Address)
		}
		if contract.ContractType != ContractTypeSystem {
			t.Errorf("Expected system contract type")
		}
		if !contract.Upgradable {
			t.Errorf("Staking contract should be upgradable via governance")
		}

		expectedFunctions := []string{"stake", "unstake", "claimRewards", "slash", "getStakeInfo", "distributeRewards", "isValidator"}
		for _, fn := range expectedFunctions {
			if _, ok := contract.Functions[fn]; !ok {
				t.Errorf("Missing function: %s", fn)
			}
		}
	})

	t.Run("Stake_Success", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 10000

		// Initialize contract storage
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)
		state.SetContractStorage(StakingContractAddress, "max_validators", StakingMaxValidators)
		state.SetContractStorage(StakingContractAddress, "lock_period", StakingLockPeriod)
		state.SetContractStorage(StakingContractAddress, "total_staked", 0)
		state.SetContractStorage(StakingContractAddress, "active_validators", 0)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)
		ctx := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 1, 1000)

		// Stake 2000
		err := helper.Stake(ctx, "0xVALIDATOR1", 2000)
		if err != nil {
			t.Fatalf("Stake failed: %v", err)
		}

		// Verify balance decreased
		if state.balances["0xVALIDATOR1"] != 8000 {
			t.Errorf("Expected validator balance 8000, got %d", state.balances["0xVALIDATOR1"])
		}

		// Verify stake recorded
		stakeAmount, _ := state.GetContractStorage(StakingContractAddress, "stake_0xVALIDATOR1")
		if stakeAmount != 2000 {
			t.Errorf("Expected stake 2000, got %d", stakeAmount)
		}

		// Verify validator is active
		isActive, _ := state.GetContractStorage(StakingContractAddress, "is_active_0xVALIDATOR1")
		if isActive != StakingStatusActive {
			t.Errorf("Expected validator to be active")
		}

		// Verify total staked
		totalStaked, _ := state.GetContractStorage(StakingContractAddress, "total_staked")
		if totalStaked != 2000 {
			t.Errorf("Expected total staked 2000, got %d", totalStaked)
		}

		// Verify active validator count
		activeCount, _ := state.GetContractStorage(StakingContractAddress, "active_validators")
		if activeCount != 1 {
			t.Errorf("Expected 1 active validator, got %d", activeCount)
		}

		// Verify lock period set
		lockedUntil, _ := state.GetContractStorage(StakingContractAddress, "locked_until_0xVALIDATOR1")
		if lockedUntil != 1000+StakingLockPeriod {
			t.Errorf("Expected locked until %d, got %d", 1000+StakingLockPeriod, lockedUntil)
		}
	})

	t.Run("Stake_BelowMinimum", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 10000
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)
		ctx := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 1, 1000)

		err := helper.Stake(ctx, "0xVALIDATOR1", 500) // Below min stake of 1000
		if err == nil {
			t.Fatal("Expected error for stake below minimum")
		}
	})

	t.Run("Stake_InsufficientBalance", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 500
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)
		ctx := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 1, 1000)

		err := helper.Stake(ctx, "0xVALIDATOR1", 2000)
		if err == nil {
			t.Fatal("Expected error for insufficient balance")
		}
	})

	t.Run("Stake_AdditionalStake", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 10000
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)
		state.SetContractStorage(StakingContractAddress, "max_validators", StakingMaxValidators)
		state.SetContractStorage(StakingContractAddress, "lock_period", StakingLockPeriod)
		state.SetContractStorage(StakingContractAddress, "total_staked", 0)
		state.SetContractStorage(StakingContractAddress, "active_validators", 0)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)
		ctx := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 1, 1000)

		// First stake
		err := helper.Stake(ctx, "0xVALIDATOR1", 2000)
		if err != nil {
			t.Fatalf("First stake failed: %v", err)
		}

		// Second stake (additional)
		err = helper.Stake(ctx, "0xVALIDATOR1", 3000)
		if err != nil {
			t.Fatalf("Second stake failed: %v", err)
		}

		// Total stake should be 5000
		stakeAmount, _ := state.GetContractStorage(StakingContractAddress, "stake_0xVALIDATOR1")
		if stakeAmount != 5000 {
			t.Errorf("Expected stake 5000, got %d", stakeAmount)
		}

		// Active validators should still be 1
		activeCount, _ := state.GetContractStorage(StakingContractAddress, "active_validators")
		if activeCount != 1 {
			t.Errorf("Expected 1 active validator, got %d", activeCount)
		}
	})

	t.Run("Unstake_Success", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 10000
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)
		state.SetContractStorage(StakingContractAddress, "max_validators", StakingMaxValidators)
		state.SetContractStorage(StakingContractAddress, "lock_period", StakingLockPeriod)
		state.SetContractStorage(StakingContractAddress, "total_staked", 0)
		state.SetContractStorage(StakingContractAddress, "active_validators", 0)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)

		// Stake at time 1000
		ctx1 := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 1, 1000)
		err := helper.Stake(ctx1, "0xVALIDATOR1", 3000)
		if err != nil {
			t.Fatalf("Stake failed: %v", err)
		}

		// Try unstake before lock period (should fail)
		ctx2 := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 10, 2000)
		err = helper.Unstake(ctx2, "0xVALIDATOR1", 1000)
		if err == nil {
			t.Fatal("Expected lock period error")
		}

		// Unstake after lock period
		ctx3 := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 100, 1000+StakingLockPeriod+1)
		err = helper.Unstake(ctx3, "0xVALIDATOR1", 1000)
		if err != nil {
			t.Fatalf("Unstake failed: %v", err)
		}

		// Verify balance returned
		if state.balances["0xVALIDATOR1"] != 8000 { // 10000 - 3000 + 1000
			t.Errorf("Expected balance 8000, got %d", state.balances["0xVALIDATOR1"])
		}

		// Verify remaining stake
		stakeAmount, _ := state.GetContractStorage(StakingContractAddress, "stake_0xVALIDATOR1")
		if stakeAmount != 2000 {
			t.Errorf("Expected remaining stake 2000, got %d", stakeAmount)
		}
	})

	t.Run("Unstake_BelowMinimum_Deactivates", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 10000
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)
		state.SetContractStorage(StakingContractAddress, "max_validators", StakingMaxValidators)
		state.SetContractStorage(StakingContractAddress, "lock_period", int64(1)) // 1 second lock
		state.SetContractStorage(StakingContractAddress, "total_staked", 0)
		state.SetContractStorage(StakingContractAddress, "active_validators", 0)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)

		// Stake at timestamp 1000
		ctx1 := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 1, 1000)
		err := helper.Stake(ctx1, "0xVALIDATOR1", 1500)
		if err != nil {
			t.Fatalf("Stake failed: %v", err)
		}

		// Unstake after lock (timestamp 1002 > locked_until 1001)
		ctx2 := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 10, 1002)
		err = helper.Unstake(ctx2, "0xVALIDATOR1", 1000)
		if err != nil {
			t.Fatalf("Unstake failed: %v", err)
		}

		isActive, _ := state.GetContractStorage(StakingContractAddress, "is_active_0xVALIDATOR1")
		if isActive != StakingStatusInactive {
			t.Errorf("Expected validator to be deactivated")
		}

		activeCount, _ := state.GetContractStorage(StakingContractAddress, "active_validators")
		if activeCount != 0 {
			t.Errorf("Expected 0 active validators, got %d", activeCount)
		}
	})

	t.Run("ClaimRewards", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 10000
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)
		state.SetContractStorage(StakingContractAddress, "max_validators", StakingMaxValidators)
		state.SetContractStorage(StakingContractAddress, "lock_period", StakingLockPeriod)
		state.SetContractStorage(StakingContractAddress, "total_staked", 0)
		state.SetContractStorage(StakingContractAddress, "active_validators", 0)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)
		ctx := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 1, 1000)

		// Stake
		err := helper.Stake(ctx, "0xVALIDATOR1", 2000)
		if err != nil {
			t.Fatalf("Stake failed: %v", err)
		}

		// Accumulate rewards
		err = helper.AccumulateValidatorReward(ctx, "0xVALIDATOR1", 100)
		if err != nil {
			t.Fatalf("AccumulateReward failed: %v", err)
		}

		// Claim rewards
		claimed, err := helper.ClaimRewards(ctx, "0xVALIDATOR1")
		if err != nil {
			t.Fatalf("ClaimRewards failed: %v", err)
		}
		if claimed != 100 {
			t.Errorf("Expected 100 rewards, got %d", claimed)
		}

		// Verify rewards reset to 0
		rewards, _ := state.GetContractStorage(StakingContractAddress, "rewards_0xVALIDATOR1")
		if rewards != 0 {
			t.Errorf("Expected rewards reset to 0, got %d", rewards)
		}

		// Verify balance includes minted rewards
		if state.balances["0xVALIDATOR1"] != 8100 { // 10000 - 2000 staked + 100 minted
			t.Errorf("Expected balance 8100, got %d", state.balances["0xVALIDATOR1"])
		}
	})

	t.Run("Slash", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 10000
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)
		state.SetContractStorage(StakingContractAddress, "max_validators", StakingMaxValidators)
		state.SetContractStorage(StakingContractAddress, "lock_period", StakingLockPeriod)
		state.SetContractStorage(StakingContractAddress, "total_staked", 0)
		state.SetContractStorage(StakingContractAddress, "active_validators", 0)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)
		ctx := NewExecutionContextWithState("CONSENSUS", 100000, state, 1, 1000)

		// Stake 2000
		err := helper.Stake(ctx, "0xVALIDATOR1", 2000)
		if err != nil {
			t.Fatalf("Stake failed: %v", err)
		}

		// Slash 500
		err = helper.Slash(ctx, "0xVALIDATOR1", 500, "double_signing")
		if err != nil {
			t.Fatalf("Slash failed: %v", err)
		}

		// Verify remaining stake
		stakeAmount, _ := state.GetContractStorage(StakingContractAddress, "stake_0xVALIDATOR1")
		if stakeAmount != 1500 {
			t.Errorf("Expected stake 1500 after slash, got %d", stakeAmount)
		}

		// Verify total slashed tracked
		totalSlashed, _ := state.GetContractStorage(StakingContractAddress, "slashed_0xVALIDATOR1")
		if totalSlashed != 500 {
			t.Errorf("Expected total slashed 500, got %d", totalSlashed)
		}
	})

	t.Run("Slash_DeactivatesValidator", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 10000
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)
		state.SetContractStorage(StakingContractAddress, "max_validators", StakingMaxValidators)
		state.SetContractStorage(StakingContractAddress, "lock_period", StakingLockPeriod)
		state.SetContractStorage(StakingContractAddress, "total_staked", 0)
		state.SetContractStorage(StakingContractAddress, "active_validators", 0)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)
		ctx := NewExecutionContextWithState("CONSENSUS", 100000, state, 1, 1000)

		// Stake exactly min stake
		err := helper.Stake(ctx, "0xVALIDATOR1", StakingMinStake)
		if err != nil {
			t.Fatalf("Stake failed: %v", err)
		}

		// Slash 500 → below min stake → deactivation
		err = helper.Slash(ctx, "0xVALIDATOR1", 500, "downtime")
		if err != nil {
			t.Fatalf("Slash failed: %v", err)
		}

		isActive, _ := state.GetContractStorage(StakingContractAddress, "is_active_0xVALIDATOR1")
		if isActive != StakingStatusInactive {
			t.Errorf("Expected validator to be deactivated after slash below min stake")
		}
	})

	t.Run("GetStakeInfo", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 10000
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)
		state.SetContractStorage(StakingContractAddress, "max_validators", StakingMaxValidators)
		state.SetContractStorage(StakingContractAddress, "lock_period", StakingLockPeriod)
		state.SetContractStorage(StakingContractAddress, "total_staked", 0)
		state.SetContractStorage(StakingContractAddress, "active_validators", 0)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)
		ctx := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 1, 1000)

		helper.Stake(ctx, "0xVALIDATOR1", 2000)

		info := helper.GetStakeInfo(ctx, "0xVALIDATOR1")
		if info == nil {
			t.Fatal("Expected stake info")
		}
		if info.Amount != 2000 {
			t.Errorf("Expected stake amount 2000, got %d", info.Amount)
		}
		if !info.IsActive {
			t.Errorf("Expected validator to be active")
		}
		if info.Address != "0xVALIDATOR1" {
			t.Errorf("Expected address 0xVALIDATOR1, got %s", info.Address)
		}
	})

	t.Run("IsValidator", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVALIDATOR1"] = 10000
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)
		state.SetContractStorage(StakingContractAddress, "max_validators", StakingMaxValidators)
		state.SetContractStorage(StakingContractAddress, "lock_period", StakingLockPeriod)
		state.SetContractStorage(StakingContractAddress, "total_staked", 0)
		state.SetContractStorage(StakingContractAddress, "active_validators", 0)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)
		ctx := NewExecutionContextWithState("0xVALIDATOR1", 100000, state, 1, 1000)

		// Not a validator yet
		if helper.IsValidator(ctx, "0xVALIDATOR1") {
			t.Errorf("Should not be a validator before staking")
		}

		// Stake → becomes validator
		helper.Stake(ctx, "0xVALIDATOR1", 2000)

		if !helper.IsValidator(ctx, "0xVALIDATOR1") {
			t.Errorf("Should be a validator after staking")
		}

		// Unknown address
		if helper.IsValidator(ctx, "0xRANDOM") {
			t.Errorf("Unknown address should not be a validator")
		}
	})

	t.Run("MultipleValidators", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xVAL1"] = 10000
		state.balances["0xVAL2"] = 10000
		state.balances["0xVAL3"] = 10000
		state.SetContractStorage(StakingContractAddress, "min_stake", StakingMinStake)
		state.SetContractStorage(StakingContractAddress, "max_validators", StakingMaxValidators)
		state.SetContractStorage(StakingContractAddress, "lock_period", StakingLockPeriod)
		state.SetContractStorage(StakingContractAddress, "total_staked", 0)
		state.SetContractStorage(StakingContractAddress, "active_validators", 0)

		contract := CreateSystemStakingContract()
		helper := NewStakingContractHelper(contract)
		ctx := NewExecutionContextWithState("test", 100000, state, 1, 1000)

		helper.Stake(ctx, "0xVAL1", 2000)
		helper.Stake(ctx, "0xVAL2", 3000)
		helper.Stake(ctx, "0xVAL3", 5000)

		totalStaked := helper.GetTotalStaked(ctx)
		if totalStaked != 10000 {
			t.Errorf("Expected total staked 10000, got %d", totalStaked)
		}

		activeCount := helper.GetActiveValidatorCount(ctx)
		if activeCount != 3 {
			t.Errorf("Expected 3 active validators, got %d", activeCount)
		}
	})
}
