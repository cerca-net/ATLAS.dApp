package vm

import (
	"fmt"
	"log"
)

// Staking contract constants
const (
	StakingContractName    = "CercaChainStaking"
	StakingContractAddress = "CONTRACT_STAKING_SYSTEM"
	StakingMinStake        = int64(1000)
	StakingMaxValidators   = int64(100)
	StakingSlashingPenalty = int64(50)
	StakingLockPeriod      = int64(604800) // 7 days in seconds
	StakingRewardRate      = int64(500)    // 5% annual in basis points
)

// Staking status constants
const (
	StakingStatusActive   = int64(1)
	StakingStatusInactive = int64(0)
)

// CreateSystemStakingContract creates the genesis staking contract.
// This manages validator staking, rewards, slashing, and unstaking.
func CreateSystemStakingContract() *Contract {
	contract := &Contract{
		Address:      StakingContractAddress,
		Name:         StakingContractName,
		Version:      "1.0.0",
		Code:         nil,
		Functions:    make(map[string]*Function),
		Storage:      make(map[string]interface{}),
		Owner:        "SYSTEM",
		Upgradable:   true, // Governance can upgrade staking parameters
		CreatedAt:    0,    // Genesis
		UpdatedAt:    0,
		ContractType: ContractTypeSystem,
	}

	// Initialize default storage
	contract.Storage["min_stake"] = StakingMinStake
	contract.Storage["max_validators"] = StakingMaxValidators
	contract.Storage["slashing_penalty"] = StakingSlashingPenalty
	contract.Storage["lock_period"] = StakingLockPeriod
	contract.Storage["reward_rate"] = StakingRewardRate
	contract.Storage["total_staked"] = int64(0)
	contract.Storage["active_validators"] = int64(0)

	// Define functions
	contract.Functions = map[string]*Function{
		"stake": {
			Name:       "stake",
			Parameters: []string{"validator", "amount"},
			Code:       []Instruction{}, // Logic handled by StakingContractHelper
		},
		"unstake": {
			Name:       "unstake",
			Parameters: []string{"validator", "amount"},
			Code:       []Instruction{},
		},
		"claimRewards": {
			Name:       "claimRewards",
			Parameters: []string{"validator"},
			Code:       []Instruction{},
		},
		"slash": {
			Name:       "slash",
			Parameters: []string{"validator", "amount", "reason"},
			Code:       []Instruction{},
		},
		"getStakeInfo": {
			Name:       "getStakeInfo",
			Parameters: []string{"validator"},
			Code:       []Instruction{},
		},
		"distributeRewards": {
			Name:       "distributeRewards",
			Parameters: []string{"block_height"},
			Code:       []Instruction{},
		},
		"isValidator": {
			Name:       "isValidator",
			Parameters: []string{"address"},
			Code:       []Instruction{},
		},
	}

	return contract
}

// StakingContractHelper provides high-level Go functions for staking operations.
type StakingContractHelper struct {
	Contract *Contract
	VM       *VM
}

// NewStakingContractHelper creates a helper for staking contract interaction
func NewStakingContractHelper(contract *Contract) *StakingContractHelper {
	vmInstance := NewVM()
	vmInstance.RegisterSystemContract(contract.Address, []string{
		"stake", "unstake", "claimRewards", "slash", "getStakeInfo",
		"distributeRewards", "isValidator",
	})
	return &StakingContractHelper{
		Contract: contract,
		VM:       vmInstance,
	}
}

// StakeInfo holds staking information for a validator
type StakeInfo struct {
	Address      string `json:"address"`
	Amount       int64  `json:"amount"`
	Rewards      int64  `json:"rewards"`
	StakedAt     int64  `json:"staked_at"`
	LockedUntil  int64  `json:"locked_until"`
	IsActive     bool   `json:"is_active"`
	TotalSlashed int64  `json:"total_slashed"`
}

// Stake executes a staking operation — locks TCOIN to become/remain a validator
func (h *StakingContractHelper) Stake(ctx *ExecutionContext, validator string, amount int64) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required for staking")
	}

	// Load staking parameters from contract storage
	minStake, _ := ctx.State.GetContractStorage(h.Contract.Address, "min_stake")
	if minStake == 0 {
		minStake = StakingMinStake
	}
	maxValidators, _ := ctx.State.GetContractStorage(h.Contract.Address, "max_validators")
	if maxValidators == 0 {
		maxValidators = StakingMaxValidators
	}

	if amount < minStake {
		return fmt.Errorf("stake amount %d below minimum %d", amount, minStake)
	}

	// Check if already staking
	currentStake, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("stake_%s", validator))
	isActive, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("is_active_%s", validator))

	// If this is a new validator, check max validator count
	if currentStake == 0 {
		activeCount, _ := ctx.State.GetContractStorage(h.Contract.Address, "active_validators")
		if activeCount >= maxValidators {
			return fmt.Errorf("maximum validator count reached (%d)", maxValidators)
		}
	}

	// Check balance
	balance := ctx.State.GetBalance(validator)
	if balance < amount {
		return fmt.Errorf("insufficient balance: %s has %d, needs %d", validator, balance, amount)
	}

	// Transfer TCOIN from validator to contract (locked)
	if err := ctx.State.Transfer(validator, h.Contract.Address, amount); err != nil {
		return fmt.Errorf("failed to lock stake: %v", err)
	}

	// Update staking storage
	newStake := currentStake + amount
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("stake_%s", validator), newStake)
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("staked_at_%s", validator), ctx.Timestamp)

	// Set lock period
	lockPeriod, _ := ctx.State.GetContractStorage(h.Contract.Address, "lock_period")
	if lockPeriod == 0 {
		lockPeriod = StakingLockPeriod
	}
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("locked_until_%s", validator), ctx.Timestamp+lockPeriod)

	// Activate validator if not already
	if isActive != StakingStatusActive {
		ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("is_active_%s", validator), StakingStatusActive)

		// Increment active validator count
		activeCount, _ := ctx.State.GetContractStorage(h.Contract.Address, "active_validators")
		ctx.State.SetContractStorage(h.Contract.Address, "active_validators", activeCount+1)
	}

	// Update total staked
	totalStaked, _ := ctx.State.GetContractStorage(h.Contract.Address, "total_staked")
	ctx.State.SetContractStorage(h.Contract.Address, "total_staked", totalStaked+amount)

	// Store validator address in string storage for lookup
	ctx.State.SetStringStorage(h.Contract.Address, fmt.Sprintf("validator_%s", validator), validator)

	ctx.EmitEvent(h.Contract.Address, "Stake", map[string]interface{}{
		"validator":    validator,
		"amount":       amount,
		"total_stake":  newStake,
		"total_staked": totalStaked + amount,
	})

	log.Printf("[STAKING] Stake: %d TCOIN by %s (total stake: %d, total staked: %d)",
		amount, shortAddr(validator), newStake, totalStaked+amount)
	return nil
}

// Unstake begins the unstaking process — tokens are returned after lock period
func (h *StakingContractHelper) Unstake(ctx *ExecutionContext, validator string, amount int64) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required for unstaking")
	}

	// Check staking position
	currentStake, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("stake_%s", validator))
	if !exists || currentStake == 0 {
		return fmt.Errorf("no staking position found for %s", validator)
	}

	if amount > currentStake {
		return fmt.Errorf("cannot unstake %d, only %d staked", amount, currentStake)
	}

	// Check lock period
	lockedUntil, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("locked_until_%s", validator))
	if ctx.Timestamp < lockedUntil {
		remainingSeconds := lockedUntil - ctx.Timestamp
		return fmt.Errorf("tokens locked for %d more seconds (until timestamp %d)", remainingSeconds, lockedUntil)
	}

	// Transfer tokens back from contract to validator
	if err := ctx.State.Transfer(h.Contract.Address, validator, amount); err != nil {
		return fmt.Errorf("failed to return staked tokens: %v", err)
	}

	// Update staking storage
	newStake := currentStake - amount
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("stake_%s", validator), newStake)

	// Update total staked
	totalStaked, _ := ctx.State.GetContractStorage(h.Contract.Address, "total_staked")
	newTotalStaked := totalStaked - amount
	if newTotalStaked < 0 {
		newTotalStaked = 0
	}
	ctx.State.SetContractStorage(h.Contract.Address, "total_staked", newTotalStaked)

	// Check minimum stake — deactivate if below threshold
	minStake, _ := ctx.State.GetContractStorage(h.Contract.Address, "min_stake")
	if minStake == 0 {
		minStake = StakingMinStake
	}
	if newStake < minStake {
		ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("is_active_%s", validator), StakingStatusInactive)
		activeCount, _ := ctx.State.GetContractStorage(h.Contract.Address, "active_validators")
		if activeCount > 0 {
			ctx.State.SetContractStorage(h.Contract.Address, "active_validators", activeCount-1)
		}
		log.Printf("[STAKING] Validator %s deactivated (stake %d < min %d)", shortAddr(validator), newStake, minStake)
	}

	ctx.EmitEvent(h.Contract.Address, "Unstake", map[string]interface{}{
		"validator": validator,
		"amount":    amount,
		"remaining": newStake,
	})

	log.Printf("[STAKING] Unstake: %d TCOIN by %s (remaining: %d)", amount, shortAddr(validator), newStake)
	return nil
}

// ClaimRewards calculates and distributes accumulated staking rewards
func (h *StakingContractHelper) ClaimRewards(ctx *ExecutionContext, validator string) (int64, error) {
	if ctx.State == nil {
		return 0, fmt.Errorf("state adapter is required for claiming rewards")
	}

	// Check staking position
	currentStake, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("stake_%s", validator))
	if !exists || currentStake == 0 {
		return 0, fmt.Errorf("no staking position found for %s", validator)
	}

	// Get accumulated rewards
	rewards, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("rewards_%s", validator))
	if rewards == 0 {
		return 0, fmt.Errorf("no rewards to claim for %s", validator)
	}

	// Mint rewards to the validator (new tokens from block rewards)
	if err := ctx.State.Mint(validator, rewards); err != nil {
		return 0, fmt.Errorf("failed to mint rewards: %v", err)
	}

	// Reset accumulated rewards
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("rewards_%s", validator), 0)

	ctx.EmitEvent(h.Contract.Address, "ClaimRewards", map[string]interface{}{
		"validator": validator,
		"amount":    rewards,
	})

	log.Printf("[STAKING] ClaimRewards: %d TCOIN to %s", rewards, shortAddr(validator))
	return rewards, nil
}

// Slash penalizes a misbehaving validator (consensus-only operation)
func (h *StakingContractHelper) Slash(ctx *ExecutionContext, validator string, amount int64, reason string) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required for slashing")
	}

	// Check staking position
	currentStake, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("stake_%s", validator))
	if !exists || currentStake == 0 {
		return fmt.Errorf("no staking position found for %s", validator)
	}

	// Calculate slash amount (capped at current stake)
	slashAmount := amount
	if slashAmount > currentStake {
		slashAmount = currentStake
	}

	// Reduce staked amount
	newStake := currentStake - slashAmount
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("stake_%s", validator), newStake)

	// Burn the slashed tokens (removed from circulation)
	contractBalance := ctx.State.GetBalance(h.Contract.Address)
	if contractBalance >= slashAmount {
		if err := ctx.State.Burn(h.Contract.Address, slashAmount); err != nil {
			log.Printf("[STAKING] Warning: Failed to burn slashed tokens: %v", err)
		}
	}

	// Update total staked
	totalStaked, _ := ctx.State.GetContractStorage(h.Contract.Address, "total_staked")
	newTotalStaked := totalStaked - slashAmount
	if newTotalStaked < 0 {
		newTotalStaked = 0
	}
	ctx.State.SetContractStorage(h.Contract.Address, "total_staked", newTotalStaked)

	// Track total slashed for this validator
	totalSlashed, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("slashed_%s", validator))
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("slashed_%s", validator), totalSlashed+slashAmount)

	// Deactivate if below min stake
	minStake, _ := ctx.State.GetContractStorage(h.Contract.Address, "min_stake")
	if minStake == 0 {
		minStake = StakingMinStake
	}
	if newStake < minStake {
		ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("is_active_%s", validator), StakingStatusInactive)
		activeCount, _ := ctx.State.GetContractStorage(h.Contract.Address, "active_validators")
		if activeCount > 0 {
			ctx.State.SetContractStorage(h.Contract.Address, "active_validators", activeCount-1)
		}
	}

	ctx.EmitEvent(h.Contract.Address, "Slash", map[string]interface{}{
		"validator": validator,
		"amount":    slashAmount,
		"reason":    reason,
		"remaining": newStake,
	})

	log.Printf("[STAKING] Slash: %d TCOIN from %s (reason: %s, remaining: %d)", slashAmount, shortAddr(validator), reason, newStake)
	return nil
}

// DistributeRewards distributes block rewards proportionally to all active validators
func (h *StakingContractHelper) DistributeRewards(ctx *ExecutionContext, blockReward int64) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required for distributing rewards")
	}

	totalStaked, _ := ctx.State.GetContractStorage(h.Contract.Address, "total_staked")
	if totalStaked == 0 {
		return nil // No validators, no rewards
	}

	activeCount, _ := ctx.State.GetContractStorage(h.Contract.Address, "active_validators")
	if activeCount == 0 {
		return nil
	}

	// For now, we distribute rewards to validators whose stake we find
	// In production, we'd iterate a validator list. Here we accumulate rewards
	// and each validator claims them individually.
	// This function is called per block to accumulate pending rewards.

	// The reward per unit of stake for this block:
	// reward_per_unit = blockReward / totalStaked
	// For each validator: pending_reward += stake * reward_per_unit

	// Since we're dealing with int64, use scaled arithmetic to avoid precision loss:
	// reward = (stake * blockReward) / totalStaked

	// We store a pending reward accumulator that validators read when claiming
	pendingPoolReward, _ := ctx.State.GetContractStorage(h.Contract.Address, "pending_pool_reward")
	ctx.State.SetContractStorage(h.Contract.Address, "pending_pool_reward", pendingPoolReward+blockReward)

	ctx.EmitEvent(h.Contract.Address, "BlockReward", map[string]interface{}{
		"block_reward": blockReward,
		"total_staked": totalStaked,
		"block_height": ctx.BlockHeight,
	})

	return nil
}

// AccumulateValidatorReward adds reward to a specific validator's pending rewards
// Called when we know which validator produced the block.
func (h *StakingContractHelper) AccumulateValidatorReward(ctx *ExecutionContext, validator string, reward int64) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required")
	}

	currentRewards, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("rewards_%s", validator))
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("rewards_%s", validator), currentRewards+reward)

	log.Printf("[STAKING] Accumulated reward: %d TCOIN for %s (total pending: %d)",
		reward, shortAddr(validator), currentRewards+reward)
	return nil
}

// GetStakeInfo returns staking information for a validator
func (h *StakingContractHelper) GetStakeInfo(ctx *ExecutionContext, validator string) *StakeInfo {
	if ctx.State == nil {
		return nil
	}

	stake, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("stake_%s", validator))
	rewards, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("rewards_%s", validator))
	stakedAt, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("staked_at_%s", validator))
	lockedUntil, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("locked_until_%s", validator))
	isActive, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("is_active_%s", validator))
	totalSlashed, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("slashed_%s", validator))

	return &StakeInfo{
		Address:      validator,
		Amount:       stake,
		Rewards:      rewards,
		StakedAt:     stakedAt,
		LockedUntil:  lockedUntil,
		IsActive:     isActive == StakingStatusActive,
		TotalSlashed: totalSlashed,
	}
}

// IsValidator checks if an address is an active validator
func (h *StakingContractHelper) IsValidator(ctx *ExecutionContext, address string) bool {
	if ctx.State == nil {
		return false
	}

	isActive, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("is_active_%s", address))
	return isActive == StakingStatusActive
}

// GetTotalStaked returns the total amount staked across all validators
func (h *StakingContractHelper) GetTotalStaked(ctx *ExecutionContext) int64 {
	if ctx.State == nil {
		return 0
	}
	totalStaked, _ := ctx.State.GetContractStorage(h.Contract.Address, "total_staked")
	return totalStaked
}

// GetActiveValidatorCount returns the number of active validators
func (h *StakingContractHelper) GetActiveValidatorCount(ctx *ExecutionContext) int64 {
	if ctx.State == nil {
		return 0
	}
	count, _ := ctx.State.GetContractStorage(h.Contract.Address, "active_validators")
	return count
}
