package defi

import (
	"sync"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/internal/blockchain"
)

// Tokenomics encapsulates all logic related to the native token
// including metadata, supply, minting, burning, fees, and rewards

type Tokenomics struct {
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply uint64
	mu          sync.RWMutex
}

func NewTokenomics(name, symbol string, decimals uint8, initialSupply uint64) *Tokenomics {
	return &Tokenomics{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: initialSupply,
	}
}

// Mint new tokens
func (t *Tokenomics) Mint(amount uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.TotalSupply += amount
}

// Burn tokens
func (t *Tokenomics) Burn(amount uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.TotalSupply >= amount {
		t.TotalSupply -= amount
	}
}

// Calculate transaction fee (example: base fee + size-based fee)
func (t *Tokenomics) CalculateFee(tx transaction.Transaction) uint64 {
	return 10 + uint64(len(tx.Data))/100 // Adjust as needed
}

// Distribute rewards to validators (example: equal distribution)
func (t *Tokenomics) DistributeRewards(validators []*blockchain.Validator, totalReward uint64) {
	if len(validators) == 0 {
		return
	}
	reward := totalReward / uint64(len(validators))
	for _, v := range validators {
		v.Stake += reward
	}
} 