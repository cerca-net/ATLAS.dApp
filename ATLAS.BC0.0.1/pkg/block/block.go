package block

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/wallet"
	"time"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"math/big"
	"fmt"
)

// Block is the fundamental component of the blockchain.
// Each block contains a list of transactions and links to the previous block.
type Block struct {
	Index        int
	Timestamp    int64
	Transactions []transaction.Transaction
	PrevHash     string
	Hash         string
	Validator    string // Store validator public key (hex)
	Signature    string
}

// CalculateHash generates the hash for a given block.
// The hash is computed using SHA256 and includes all block data.
func CalculateHash(block Block) string {
	var txDetails string
	for _, tx := range block.Transactions {
		txDetails += tx.Sender + tx.Recipient + strconv.FormatInt(tx.Amount, 10)
	}
	record := strconv.Itoa(block.Index) + strconv.FormatInt(block.Timestamp, 10) + txDetails + block.PrevHash + block.Validator + block.Signature
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// HashBlockForSigning returns the hash of the block excluding the signature field.
func HashBlockForSigning(block *Block) []byte {
	var txDetails string
	for _, tx := range block.Transactions {
		txDetails += tx.Sender + tx.Recipient + strconv.FormatInt(tx.Amount, 10)
	}
	record := strconv.Itoa(block.Index) + strconv.FormatInt(block.Timestamp, 10) + txDetails + block.PrevHash + block.Validator
	h := sha256.New()
	h.Write([]byte(record))
	return h.Sum(nil)
}

// CreateGenesisBlock creates the first block in the chain.
// The genesis block has no transactions and a special previous hash.
func CreateGenesisBlock() *Block {
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    1640995200, // Fixed timestamp: January 1, 2022 00:00:00 UTC
		Transactions: []transaction.Transaction{},
		PrevHash:     "0",
		Validator:    "GENESIS_VALIDATOR",
		Signature:    "GENESIS_SIGNATURE",
	}
	genesisBlock.Hash = CalculateHash(*genesisBlock)
	return genesisBlock
}

// CreateNewBlock creates a new block to be added to the chain.
// It validates all transactions before creating the block.
func CreateNewBlock(transactions []transaction.Transaction, prevBlock *Block, validatorWallet *wallet.Wallet) (*Block, error) {
	// Validate all transactions
	for _, tx := range transactions {
		if err := tx.Validate(); err != nil {
			return nil, err
		}
		if tx.Sender != "network" {
			valid, err := wallet.VerifyTransactionSignature(tx)
			if err != nil {
				return nil, err
			}
			if !valid {
				return nil, err
			}
		}
	}

	validatorPubKeyHex := validatorWallet.PublicKeyStr()
	newBlock := &Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    int64(time.Now().Unix()),
		Transactions: transactions,
		PrevHash:     prevBlock.Hash,
		Validator:    validatorPubKeyHex,
		Signature:    "",
	}
	sig, err := SignBlock(newBlock, validatorWallet)
	if err != nil {
		return nil, err
	}
	newBlock.Signature = sig
	newBlock.Hash = CalculateHash(*newBlock)
	return newBlock, nil
}

// SignBlock signs the hash of the block (excluding the signature field) with the validator's private key.
func SignBlock(b *Block, w *wallet.Wallet) (string, error) {
	if w.PrivateKey == nil {
		fmt.Println("[DEBUG] SignBlock: wallet.PrivateKey is nil!")
		return "", fmt.Errorf("wallet has no private key")
	}
	fmt.Printf("[DEBUG] SignBlock: PrivateKey type: %T\n", w.PrivateKey)
	hash := HashBlockForSigning(b)
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, hash)
	if err != nil {
		fmt.Printf("[DEBUG] SignBlock: failed to sign block: %v\n", err)
		return "", fmt.Errorf("failed to sign block: %v", err)
	}
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	sig := append(rBytes, sBytes...)
	sigHex := hex.EncodeToString(sig)
	fmt.Printf("[DEBUG] SignBlock: signature (hex): %s\n", sigHex)
	return sigHex, nil
}

// VerifyBlockSignature verifies the block's signature using the validator's public key (hex-encoded DER format).
func VerifyBlockSignature(b *Block, validatorPubKey []byte) (bool, error) {
	pubKeyInterface, err := x509.ParsePKIXPublicKey(validatorPubKey)
	if err != nil {
		return false, fmt.Errorf("failed to parse validator public key: %v", err)
	}
	pubKey, ok := pubKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("validator public key is not ECDSA")
	}
	hash := HashBlockForSigning(b)
	sigBytes, err := hex.DecodeString(b.Signature)
	if err != nil {
		return false, fmt.Errorf("invalid block signature encoding: %v", err)
	}
	if len(sigBytes)%2 != 0 {
		return false, fmt.Errorf("invalid signature length")
	}
	r := new(big.Int).SetBytes(sigBytes[:len(sigBytes)/2])
	s := new(big.Int).SetBytes(sigBytes[len(sigBytes)/2:])
	return ecdsa.Verify(pubKey, hash, r, s), nil
}

// ValidateChain checks the integrity of the entire blockchain.
// It verifies:
// 1. Each block's previous hash matches the hash of the preceding block
// 2. Each block's hash is correctly calculated
// 3. The chain is properly ordered by index
// 4. Each block's signature is valid
func ValidateChain(chain []*Block) bool {
	if len(chain) == 0 {
		return false
	}
	for i := 1; i < len(chain); i++ {
		prevBlock := chain[i-1]
		currentBlock := chain[i]
		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
		if currentBlock.Hash != CalculateHash(*currentBlock) {
			return false
		}
		if currentBlock.Index != prevBlock.Index+1 {
			return false
		}
		// Validate signature if not genesis
		if currentBlock.Index != 0 && currentBlock.Validator != "GENESIS_VALIDATOR" {
			validatorPubKeyBytes, err := hex.DecodeString(currentBlock.Validator)
			if err != nil {
				return false
			}
			valid, err := VerifyBlockSignature(currentBlock, validatorPubKeyBytes)
			if err != nil || !valid {
				return false
			}
		}
	}
	return true
} 