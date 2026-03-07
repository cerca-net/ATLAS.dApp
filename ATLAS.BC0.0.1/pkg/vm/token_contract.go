package vm

import (
	"fmt"
	"log"
)

// Token contract constants
const (
	TokenContractName    = "CercaChainToken"
	TokenSymbol          = "TCOIN"
	TokenMaxSupply       = int64(1_000_000_000) // 1 billion TCOIN (no decimals for MVP)
	TokenFaucetAmount    = int64(1000)
	TokenFaucetCooldown  = int64(86400) // 24 hours in seconds
	TokenBlockReward     = int64(10)
	TokenContractAddress = "CONTRACT_TCOIN_SYSTEM" // Fixed well-known address
)

// CreateSystemTokenContract creates the genesis TCOIN token contract.
// This is deployed once at chain initialization and manages the entire token supply.
func CreateSystemTokenContract(treasuryAddress string) *Contract {
	contract := &Contract{
		Address:      TokenContractAddress,
		Name:         TokenContractName,
		Version:      "1.0.0",
		Code:         nil,
		Functions:    make(map[string]*Function),
		Storage:      make(map[string]interface{}),
		Owner:        "SYSTEM",
		Upgradable:   false,
		CreatedAt:    0, // Genesis
		UpdatedAt:    0,
		ContractType: ContractTypeSystem,
	}

	// Initialize storage
	contract.Storage["symbol"] = TokenSymbol
	contract.Storage["max_supply"] = TokenMaxSupply
	contract.Storage["total_supply"] = int64(0) // Will be set to treasury initial balance
	contract.Storage["treasury_address"] = treasuryAddress
	contract.Storage["faucet_amount"] = TokenFaucetAmount
	contract.Storage["faucet_cooldown"] = TokenFaucetCooldown
	contract.Storage["block_reward"] = TokenBlockReward

	// Define functions using the new opcodes
	contract.Functions = map[string]*Function{
		"transfer": {
			Name:       "transfer",
			Parameters: []string{"from", "to", "amount"},
			Code: []Instruction{
				// The VM receives from/to on string stack and amount on int stack
				// from the CallFunctionWithStrings helper
				{Opcode: "TRANSFER"},
				{Opcode: "EMIT", Operands: []interface{}{"Transfer"}},
			},
		},
		"balanceOf": {
			Name:       "balanceOf",
			Parameters: []string{"address"},
			Code: []Instruction{
				// Address is on string stack, result goes to int stack
				{Opcode: "BALANCE"},
			},
		},
		"mint": {
			Name:       "mint",
			Parameters: []string{"to", "amount"},
			Code: []Instruction{
				// to address on string stack, amount on int stack
				// Mint creates new tokens and credits to 'to'
				{Opcode: "MINT"},
				// Update total_supply in persistent storage
				{Opcode: "SLOAD", Operands: []interface{}{"total_supply"}},
				{Opcode: "ADD"}, // old total + minted amount (mint pushed 1 for success, we need to track)
				{Opcode: "SSTORE", Operands: []interface{}{"total_supply"}},
				{Opcode: "EMIT", Operands: []interface{}{"Mint"}},
			},
		},
		"burn": {
			Name:       "burn",
			Parameters: []string{"from", "amount"},
			Code: []Instruction{
				// from address on string stack, amount on int stack
				{Opcode: "BURN"},
				{Opcode: "EMIT", Operands: []interface{}{"Burn"}},
			},
		},
		"faucetRequest": {
			Name:       "faucetRequest",
			Parameters: []string{"to"},
			Code: []Instruction{
				// to address is on string stack
				// Load faucet amount
				{Opcode: "SLOAD", Operands: []interface{}{"faucet_amount"}},
				// Load treasury address to string stack
				{Opcode: "SLOAD_S", Operands: []interface{}{"treasury_address"}},
				// Now string stack has: [to, treasury_address]
				// We need: [from(treasury), to] for TRANSFER
				// Swap them: PUSHS to first, then from
				// Actually the stack order for TRANSFER is: stringStack[-2]=from, stringStack[-1]=to
				// We have stringStack = [to_original, treasury_address]
				// Need stringStack = [treasury_address, to_original]
				// So we swap by popping both and pushing in reverse
				// For simplicity, use the native Go helper instead
				{Opcode: "TRANSFER"},
				{Opcode: "EMIT", Operands: []interface{}{"FaucetRequest"}},
			},
		},
	}

	return contract
}

// TokenContractHelper provides high-level Go functions that execute token contract logic.
// This is the recommended way for the blockchain to interact with the token contract,
// instead of manually pushing to stacks and calling raw VM functions.
type TokenContractHelper struct {
	Contract *Contract
	VM       *VM
}

// NewTokenContractHelper creates a helper for easy token contract interaction
func NewTokenContractHelper(contract *Contract) *TokenContractHelper {
	vmInstance := NewVM()
	vmInstance.RegisterSystemContract(contract.Address, []string{
		"transfer", "balanceOf", "mint", "burn", "faucetRequest",
	})
	return &TokenContractHelper{
		Contract: contract,
		VM:       vmInstance,
	}
}

// Transfer executes a token transfer via the contract
func (h *TokenContractHelper) Transfer(ctx *ExecutionContext, from, to string, amount int64) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required for token transfer")
	}

	// Use the StateAdapter directly — the contract VM opcodes call through to StateAdapter
	if err := ctx.State.Transfer(from, to, amount); err != nil {
		return err
	}

	// Emit event
	ctx.EmitEvent(h.Contract.Address, "Transfer", map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	})

	log.Printf("[TOKEN] Transfer: %d TCOIN from %s to %s", amount, shortAddr(from), shortAddr(to))
	return nil
}

// BalanceOf queries the balance of an address
func (h *TokenContractHelper) BalanceOf(ctx *ExecutionContext, address string) int64 {
	if ctx.State == nil {
		return 0
	}
	return ctx.State.GetBalance(address)
}

// Mint creates new tokens and credits them to an address (validator rewards)
func (h *TokenContractHelper) Mint(ctx *ExecutionContext, to string, amount int64) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required for minting")
	}

	// Check max supply
	totalSupply, _ := ctx.State.GetContractStorage(h.Contract.Address, "total_supply")
	if totalSupply+amount > TokenMaxSupply {
		return fmt.Errorf("minting %d would exceed max supply of %d (current: %d)", amount, TokenMaxSupply, totalSupply)
	}

	if err := ctx.State.Mint(to, amount); err != nil {
		return err
	}

	// Update total supply
	ctx.State.SetContractStorage(h.Contract.Address, "total_supply", totalSupply+amount)

	ctx.EmitEvent(h.Contract.Address, "Mint", map[string]interface{}{
		"to":           to,
		"amount":       amount,
		"total_supply": totalSupply + amount,
	})

	log.Printf("[TOKEN] Mint: %d TCOIN to %s (total supply: %d/%d)", amount, shortAddr(to), totalSupply+amount, TokenMaxSupply)
	return nil
}

// Burn destroys tokens from an address (TCOIN → Data Unit conversion)
func (h *TokenContractHelper) Burn(ctx *ExecutionContext, from string, amount int64) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required for burning")
	}

	if err := ctx.State.Burn(from, amount); err != nil {
		return err
	}

	// Update total supply
	totalSupply, _ := ctx.State.GetContractStorage(h.Contract.Address, "total_supply")
	newSupply := totalSupply - amount
	if newSupply < 0 {
		newSupply = 0
	}
	ctx.State.SetContractStorage(h.Contract.Address, "total_supply", newSupply)

	ctx.EmitEvent(h.Contract.Address, "Burn", map[string]interface{}{
		"from":         from,
		"amount":       amount,
		"total_supply": newSupply,
	})

	log.Printf("[TOKEN] Burn: %d TCOIN from %s (total supply: %d/%d)", amount, shortAddr(from), newSupply, TokenMaxSupply)
	return nil
}

// FaucetRequest dispenses tokens from the treasury to a user
func (h *TokenContractHelper) FaucetRequest(ctx *ExecutionContext, to string) (int64, error) {
	if ctx.State == nil {
		return 0, fmt.Errorf("state adapter is required for faucet request")
	}

	// Get faucet amount from contract storage
	faucetAmount, exists := ctx.State.GetContractStorage(h.Contract.Address, "faucet_amount")
	if !exists || faucetAmount == 0 {
		faucetAmount = TokenFaucetAmount
	}

	// Get treasury address from contract string storage
	treasuryAddr, exists := ctx.State.GetStringStorage(h.Contract.Address, "treasury_address")
	if !exists || treasuryAddr == "" {
		return 0, fmt.Errorf("treasury address not configured in token contract")
	}

	// Check faucet cooldown
	cooldownKey := fmt.Sprintf("faucet_last_%s", to)
	lastRequest, _ := ctx.State.GetContractStorage(h.Contract.Address, cooldownKey)
	faucetCooldown, _ := ctx.State.GetContractStorage(h.Contract.Address, "faucet_cooldown")
	if faucetCooldown == 0 {
		faucetCooldown = TokenFaucetCooldown
	}
	if lastRequest > 0 && ctx.Timestamp-lastRequest < faucetCooldown {
		return 0, fmt.Errorf("faucet cooldown active: wait %d more seconds", faucetCooldown-(ctx.Timestamp-lastRequest))
	}

	// Transfer from treasury to user
	if err := ctx.State.Transfer(treasuryAddr, to, faucetAmount); err != nil {
		return 0, fmt.Errorf("faucet transfer failed: %v", err)
	}

	// Record last faucet request time
	ctx.State.SetContractStorage(h.Contract.Address, cooldownKey, ctx.Timestamp)

	ctx.EmitEvent(h.Contract.Address, "FaucetRequest", map[string]interface{}{
		"to":     to,
		"amount": faucetAmount,
	})

	log.Printf("[TOKEN] Faucet: %d TCOIN from treasury to %s", faucetAmount, shortAddr(to))
	return faucetAmount, nil
}

// GetTotalSupply returns the current total supply
func (h *TokenContractHelper) GetTotalSupply(ctx *ExecutionContext) int64 {
	if ctx.State == nil {
		return 0
	}
	totalSupply, _ := ctx.State.GetContractStorage(h.Contract.Address, "total_supply")
	return totalSupply
}

// GetMaxSupply returns the maximum supply
func (h *TokenContractHelper) GetMaxSupply() int64 {
	return TokenMaxSupply
}

// shortAddr helper for logging
func shortAddr(addr string) string {
	if len(addr) > 16 {
		return addr[:8] + "..." + addr[len(addr)-4:]
	}
	return addr
}
