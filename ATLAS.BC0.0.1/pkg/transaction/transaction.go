package transaction

import (
	"atlas-blockchain/pkg/crypto/zk"
	"encoding/hex"
	"errors"
	"fmt"
)

// TransactionType defines the type of transaction (regular, contract deployment, contract call)
type TransactionType string

const (
	TxTypeRegular  TransactionType = "regular"
	TxTypeDeploy   TransactionType = "deploy_contract"
	TxTypeCall     TransactionType = "call_contract"
	TxTypeProposal TransactionType = "proposal"
	TxTypeVote     TransactionType = "vote"
	TxTypeStake    TransactionType = "stake"    // Stake DUT to become validator
	TxTypeUnstake  TransactionType = "unstake"  // Unstake DUT (future)
	TxTypeZK       TransactionType = "zk_proof" // New ZK transaction type
	TxTypeEnergize TransactionType = "energize" // Provide energy/tips to objects
)

// Transaction represents a transfer of value or a contract operation.
type Transaction struct {
	Type            TransactionType `json:"type"`
	Sender          string          `json:"sender"`
	SenderPublicKey string          `json:"senderPublicKey"`
	Recipient       string          `json:"recipient"`
	Amount          int64           `json:"amount"`
	Fee             int64           `json:"fee"`
	Timestamp       int64           `json:"timestamp"`
	Nonce           uint64          `json:"nonce"`
	Data            string          `json:"data"`
	Signature       string          `json:"signature"`
	// Privacy features
	IsEncrypted     bool        `json:"is_encrypted,omitempty"`      // Whether Data field is encrypted
	EncryptionKeyID string      `json:"encryption_key_id,omitempty"` // ID of encryption key used
	ZKProof         *zk.ZKProof `json:"zk_proof,omitempty"`          // Zero-Knowledge Proof data
}

// Validate checks if a transaction is valid.
func (t *Transaction) Validate() error {
	if t.Sender == "" {
		return errors.New("sender cannot be empty")
	}
	if t.Recipient == "" {
		return errors.New("recipient cannot be empty")
	}
	if t.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if t.Fee < 0 {
		return errors.New("fee cannot be negative")
	}
	if t.Signature == "" {
		return errors.New("signature cannot be empty")
	}
	if t.Sender != "network" && t.SenderPublicKey == "" {
		return errors.New("sender public key cannot be empty")
	}

	// Validate ZK Proof if present or if type requires it
	if t.Type == TxTypeZK {
		if t.ZKProof == nil {
			return errors.New("ZK transaction missing ZK proof")
		}
		// Basic structural check only; cryptographic verification happens in TransactionManager
		if t.ZKProof.Proof == "" {
			return errors.New("ZK proof data cannot be empty")
		}
	}

	// Normalize addresses (remove 0x prefix if present)
	sender := t.Sender
	recipient := t.Recipient
	if len(sender) > 2 && sender[:2] == "0x" {
		sender = sender[2:]
	}
	if len(recipient) > 2 && recipient[:2] == "0x" {
		recipient = recipient[2:]
	}
	if t.Sender != "network" {
		if len(sender) != 40 {
			return fmt.Errorf("invalid sender address length: %d (expected 40)", len(sender))
		}
		if _, err := hex.DecodeString(sender); err != nil {
			return fmt.Errorf("invalid sender address format (not valid hex): %s", sender)
		}
	}
	if len(recipient) != 40 {
		return fmt.Errorf("invalid recipient address length: %d (expected 40)", len(recipient))
	}
	if _, err := hex.DecodeString(recipient); err != nil {
		return fmt.Errorf("invalid recipient address format (not valid hex): %s", recipient)
	}
	return nil
}
