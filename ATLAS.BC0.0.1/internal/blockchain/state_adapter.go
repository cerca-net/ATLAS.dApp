package blockchain

import (
	"atlas-blockchain/pkg/vm"
	"fmt"
	"log"
	"sync"
)

// StateAdapterImpl implements vm.StateAdapter, bridging the VM to the StateManager.
// This allows smart contracts to read/write balances and persistent storage
// without direct access to internal blockchain types.
type StateAdapterImpl struct {
	stateManager *StateManager
	mu           sync.RWMutex

	// Contract persistent storage (int64 values)
	contractStorage map[string]map[string]int64 // contractAddress -> key -> value

	// Contract persistent string storage (addresses, etc.)
	contractStringStorage map[string]map[string]string // contractAddress -> key -> value

	// Token supply tracking
	totalSupply int64
	maxSupply   int64
}

// NewStateAdapter creates a new StateAdapter that wraps the StateManager
func NewStateAdapter(sm *StateManager) *StateAdapterImpl {
	return &StateAdapterImpl{
		stateManager:          sm,
		contractStorage:       make(map[string]map[string]int64),
		contractStringStorage: make(map[string]map[string]string),
		maxSupply:             vm.TokenMaxSupply,
	}
}

// Ensure StateAdapterImpl implements vm.StateAdapter
var _ vm.StateAdapter = (*StateAdapterImpl)(nil)

// GetBalance returns the TCOIN balance for an address
func (sa *StateAdapterImpl) GetBalance(address string) int64 {
	return sa.stateManager.GetBalance(address)
}

// SetBalance sets the TCOIN balance for an address
func (sa *StateAdapterImpl) SetBalance(address string, amount int64) {
	sa.stateManager.SetBalance(address, amount)
}

// Transfer moves tokens from one address to another
func (sa *StateAdapterImpl) Transfer(from, to string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	fromBalance := sa.stateManager.GetBalance(from)
	if fromBalance < amount {
		return fmt.Errorf("insufficient balance: %s has %d, needs %d", from, fromBalance, amount)
	}

	// Debit sender
	sa.stateManager.SetBalance(from, fromBalance-amount)

	// Credit recipient
	toBalance := sa.stateManager.GetBalance(to)
	sa.stateManager.SetBalance(to, toBalance+amount)

	log.Printf("[STATE-ADAPTER] Transfer: %d from %s (bal: %d→%d) to %s (bal: %d→%d)",
		amount, shortAddr(from), fromBalance, fromBalance-amount,
		shortAddr(to), toBalance, toBalance+amount)
	return nil
}

// Mint creates new tokens and credits them to an address
func (sa *StateAdapterImpl) Mint(to string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("mint amount must be positive")
	}

	sa.mu.Lock()
	if sa.totalSupply+amount > sa.maxSupply {
		sa.mu.Unlock()
		return fmt.Errorf("minting %d would exceed max supply of %d (current: %d)", amount, sa.maxSupply, sa.totalSupply)
	}
	sa.totalSupply += amount
	sa.mu.Unlock()

	// Credit the recipient
	currentBalance := sa.stateManager.GetBalance(to)
	sa.stateManager.SetBalance(to, currentBalance+amount)

	log.Printf("[STATE-ADAPTER] Mint: %d TCOIN to %s (total supply: %d/%d)", amount, shortAddr(to), sa.totalSupply, sa.maxSupply)
	return nil
}

// Burn destroys tokens from an address
func (sa *StateAdapterImpl) Burn(from string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("burn amount must be positive")
	}

	fromBalance := sa.stateManager.GetBalance(from)
	if fromBalance < amount {
		return fmt.Errorf("insufficient balance to burn: %s has %d, needs %d", from, fromBalance, amount)
	}

	// Debit the sender
	sa.stateManager.SetBalance(from, fromBalance-amount)

	// Decrease total supply
	sa.mu.Lock()
	sa.totalSupply -= amount
	if sa.totalSupply < 0 {
		sa.totalSupply = 0
	}
	sa.mu.Unlock()

	log.Printf("[STATE-ADAPTER] Burn: %d TCOIN from %s (total supply: %d/%d)", amount, shortAddr(from), sa.totalSupply, sa.maxSupply)
	return nil
}

// GetContractStorage reads an int64 value from a contract's persistent storage
func (sa *StateAdapterImpl) GetContractStorage(contractAddress, key string) (int64, bool) {
	sa.mu.RLock()
	defer sa.mu.RUnlock()

	if contractMap, ok := sa.contractStorage[contractAddress]; ok {
		if value, exists := contractMap[key]; exists {
			return value, true
		}
	}
	return 0, false
}

// SetContractStorage writes an int64 value to a contract's persistent storage
func (sa *StateAdapterImpl) SetContractStorage(contractAddress, key string, value int64) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	if _, ok := sa.contractStorage[contractAddress]; !ok {
		sa.contractStorage[contractAddress] = make(map[string]int64)
	}
	sa.contractStorage[contractAddress][key] = value
}

// GetStringStorage reads a string value from a contract's persistent storage
func (sa *StateAdapterImpl) GetStringStorage(contractAddress, key string) (string, bool) {
	sa.mu.RLock()
	defer sa.mu.RUnlock()

	if contractMap, ok := sa.contractStringStorage[contractAddress]; ok {
		if value, exists := contractMap[key]; exists {
			return value, true
		}
	}
	return "", false
}

// SetStringStorage writes a string value to a contract's persistent storage
func (sa *StateAdapterImpl) SetStringStorage(contractAddress, key string, value string) {
	sa.mu.Lock()
	defer sa.mu.Unlock()

	if _, ok := sa.contractStringStorage[contractAddress]; !ok {
		sa.contractStringStorage[contractAddress] = make(map[string]string)
	}
	sa.contractStringStorage[contractAddress][key] = value
}

// GetTotalSupply returns the current total TCOIN supply
func (sa *StateAdapterImpl) GetTotalSupply() int64 {
	sa.mu.RLock()
	defer sa.mu.RUnlock()
	return sa.totalSupply
}

// SetTotalSupply sets the current total TCOIN supply (used during initialization)
func (sa *StateAdapterImpl) SetTotalSupply(supply int64) {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	sa.totalSupply = supply
}

// InitializeTokenContract sets up the initial contract storage for the TCOIN contract
func (sa *StateAdapterImpl) InitializeTokenContract(treasuryAddress string, initialTreasuryBalance int64) {
	// Set int64 storage
	sa.SetContractStorage(vm.TokenContractAddress, "total_supply", initialTreasuryBalance)
	sa.SetContractStorage(vm.TokenContractAddress, "max_supply", vm.TokenMaxSupply)
	sa.SetContractStorage(vm.TokenContractAddress, "faucet_amount", vm.TokenFaucetAmount)
	sa.SetContractStorage(vm.TokenContractAddress, "faucet_cooldown", vm.TokenFaucetCooldown)
	sa.SetContractStorage(vm.TokenContractAddress, "block_reward", vm.TokenBlockReward)

	// Set string storage
	sa.SetStringStorage(vm.TokenContractAddress, "treasury_address", treasuryAddress)
	sa.SetStringStorage(vm.TokenContractAddress, "symbol", vm.TokenSymbol)
	sa.SetStringStorage(vm.TokenContractAddress, "name", vm.TokenContractName)

	// Track total supply
	sa.SetTotalSupply(initialTreasuryBalance)

	log.Printf("[STATE-ADAPTER] Token contract initialized: treasury=%s, initial_supply=%d, max_supply=%d",
		shortAddr(treasuryAddress), initialTreasuryBalance, vm.TokenMaxSupply)
}
