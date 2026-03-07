package config

import (
	"errors"
	"time"
)

// BlockchainConfig holds the configuration parameters for the blockchain.
type BlockchainConfig struct {
	// Network parameters
	MaxPeers          int           // Maximum number of peers a node can connect to
	PeerDiscoveryPort int           // Port for peer discovery
	BlockTime         time.Duration // Target time between blocks

	// General configuration
	DataDir string // Directory to store blockchain data (db, snapshots, backups)

	// Block parameters
	MaxBlockSize     int           // Maximum number of transactions per block
	MaxTxPoolSize    int           // Maximum number of transactions in the pool
	TxExpirationTime time.Duration // Time after which unconfirmed transactions expire

	// Consensus parameters
	MinStake          int // Minimum stake required to be a validator
	BlockReward       int // Reward for forging a block
	ValidatorRotation int // Number of blocks before validator rotation

	// Security parameters
	MaxValidators   int // Maximum number of validators in the pool
	SlashingPenalty int // Amount to slash for malicious behavior
}

// DefaultConfig returns the default configuration for the blockchain.
func DefaultConfig() *BlockchainConfig {
	return &BlockchainConfig{
		DataDir:           ".",
		MaxPeers:          10,
		PeerDiscoveryPort: 8000,
		BlockTime:         time.Second * 30,
		MaxBlockSize:      1000,
		MaxTxPoolSize:     5000,
		TxExpirationTime:  time.Hour * 24,
		MinStake:          1000,
		BlockReward:       10,
		ValidatorRotation: 100,
		MaxValidators:     100,
		SlashingPenalty:   50,
	}
}

// Validate checks if the configuration is valid.
func (c *BlockchainConfig) Validate() error {
	if c.MaxPeers <= 0 {
		return errors.New("MaxPeers must be positive")
	}
	if c.PeerDiscoveryPort <= 0 {
		return errors.New("PeerDiscoveryPort must be positive")
	}
	if c.BlockTime <= 0 {
		return errors.New("BlockTime must be positive")
	}
	if c.MaxBlockSize <= 0 {
		return errors.New("MaxBlockSize must be positive")
	}
	if c.MaxTxPoolSize <= 0 {
		return errors.New("MaxTxPoolSize must be positive")
	}
	if c.TxExpirationTime <= 0 {
		return errors.New("TxExpirationTime must be positive")
	}
	if c.MinStake <= 0 {
		return errors.New("MinStake must be positive")
	}
	if c.BlockReward <= 0 {
		return errors.New("BlockReward must be positive")
	}
	if c.ValidatorRotation <= 0 {
		return errors.New("ValidatorRotation must be positive")
	}
	if c.MaxValidators <= 0 {
		return errors.New("MaxValidators must be positive")
	}
	if c.SlashingPenalty <= 0 {
		return errors.New("SlashingPenalty must be positive")
	}
	return nil
}
