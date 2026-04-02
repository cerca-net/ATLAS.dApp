package config

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type GenesisAlloc struct {
	Address     string `json:"address"`
	Balance     int64  `json:"balance"`
	Description string `json:"description"`
}

type GenesisValidator struct {
	Address string `json:"address"`
	PubKey  string `json:"pub_key"`
	Stake   int64  `json:"stake"`
}

type ConsensusParams struct {
	BlockTimeMs       int64 `json:"block_time_ms"`
	MaxBlockSizeTxs   int   `json:"max_block_size_txs"`
	MaxTxPoolSize     int   `json:"max_tx_pool_size"`
	MinValidatorStake int   `json:"min_validator_stake"`
	BlockReward       int   `json:"block_reward"`
}

type SystemContracts struct {
	TcoinContractAddress       string `json:"tcoin_contract_address"`
	StakingContractAddress     string `json:"staking_contract_address"`
	MarketplaceContractAddress string `json:"marketplace_contract_address"`
	GovernanceContractAddress  string `json:"governance_contract_address"`
}

type GenesisConfig struct {
	ChainID          string             `json:"chain_id"`
	GenesisTime      time.Time          `json:"genesis_time"`
	ConsensusParams  ConsensusParams    `json:"consensus_params"`
	Alloc            []GenesisAlloc     `json:"alloc"`
	Validators       []GenesisValidator `json:"validators"`
	SystemContracts  SystemContracts    `json:"system_contracts"`
}

// LoadGenesis loads a genesis configuration from a JSON file.
func LoadGenesis(filePath string) (*GenesisConfig, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var genesis GenesisConfig
	if err := json.Unmarshal(data, &genesis); err != nil {
		return nil, err
	}

	return &genesis, nil
}
