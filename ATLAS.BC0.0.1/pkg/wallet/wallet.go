package wallet

import (
	"atlas-blockchain/pkg/transaction"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/tyler-smith/go-bip39"
)

// Wallet represents a cryptographic wallet that can sign transactions and manage stake.
// It holds a private key for signing and a public key for verification.
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Balance    int64  // Available balance (not staked)
	Staked     uint64 // Amount currently staked as validator
	mu         sync.RWMutex
}

// NewWallet creates a new wallet with a fresh ECDSA key pair using the P-256 curve.
// Returns an error if key generation fails.
func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %v", err)
	}
	address := PublicKeyToAddress(publicKeyBytes)
	fmt.Printf("[DEBUG] Backend Wallet Address (short): %s\n", address)
	fmt.Printf("[DEBUG] Backend Public Key Bytes: %v\n", publicKeyBytes)
	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKeyBytes,
		Balance:    0,
		Staked:     0,
	}, nil
}

// PublicKeyStr returns the public key as a hex string.
// This is useful for storing or transmitting the public key.
func (w *Wallet) PublicKeyStr() string {
	return hex.EncodeToString(w.PublicKey)
}

// GetPublicKey returns the ECDSA public key in its native format.
// Returns an error if the public key cannot be parsed.
func (w *Wallet) GetPublicKey() (*ecdsa.PublicKey, error) {
	pubKey, err := x509.ParsePKIXPublicKey(w.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}
	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type ECDSA")
	}
	return ecdsaPubKey, nil
}

// GetBalance returns the available balance (not staked)
func (w *Wallet) GetBalance() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Balance
}

// SetBalance sets the available balance
func (w *Wallet) SetBalance(balance int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Balance = balance
}

// GetStaked returns the amount currently staked
func (w *Wallet) GetStaked() uint64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Staked
}

// SetStaked sets the staked amount
func (w *Wallet) SetStaked(staked uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Staked = staked
}

// GetTotalBalance returns the total balance (available + staked)
func (w *Wallet) GetTotalBalance() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Balance + int64(w.Staked)
}

// CanStake checks if the wallet can stake the given amount
func (w *Wallet) CanStake(amount uint64) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Balance >= int64(amount)
}

// Stake moves tokens from balance to staked amount
func (w *Wallet) Stake(amount uint64) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Balance < int64(amount) {
		return fmt.Errorf("insufficient balance for stake: have %d, need %d", w.Balance, amount)
	}

	w.Balance -= int64(amount)
	w.Staked += amount
	return nil
}

// Unstake moves tokens from staked amount back to balance
func (w *Wallet) Unstake(amount uint64) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Staked < amount {
		return fmt.Errorf("insufficient staked amount: have %d, need %d", w.Staked, amount)
	}

	w.Staked -= amount
	w.Balance += int64(amount)
	return nil
}

// Exported version of SignTransaction
func (w *Wallet) SignTransaction(tx *transaction.Transaction) error {
	if w.PrivateKey == nil {
		return fmt.Errorf("wallet has no private key")
	}
	// Calculate the transaction hash (excluding the signature field)
	hash := CalculateTxHash(*tx)
	r, s, err := ecdsa.Sign(rand.Reader, w.PrivateKey, hash)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}
	rBytes := r.Bytes()
	sBytes := s.Bytes()
	sig := append(rBytes, sBytes...)
	tx.Signature = hex.EncodeToString(sig)
	return nil
}

// Exported version of CalculateTxHash
func CalculateTxHash(tx transaction.Transaction) []byte {
	// Exclude the signature field from the hash
	// Format must match Frontend: sender + recipient + amount + fee + timestamp + nonce + data
	record := tx.Sender +
		tx.Recipient +
		fmt.Sprintf("%d", tx.Amount) +
		fmt.Sprintf("%d", tx.Fee) +
		fmt.Sprintf("%d", tx.Timestamp) +
		fmt.Sprintf("%d", tx.Nonce) +
		tx.Data

	h := sha256.New()
	h.Write([]byte(record))
	return h.Sum(nil)
}

// Exported version of VerifyTransactionSignature
func VerifyTransactionSignature(tx transaction.Transaction) (bool, error) {
	if tx.Signature == "" {
		return false, fmt.Errorf("transaction signature is empty")
	}
	if tx.SenderPublicKey == "" {
		return false, fmt.Errorf("transaction sender public key is empty")
	}
	pubKeyBytes, err := hex.DecodeString(tx.SenderPublicKey)
	if err != nil {
		return false, fmt.Errorf("invalid sender public key encoding: %v", err)
	}
	pubKeyInterface, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse sender public key: %v", err)
	}
	pubKey, ok := pubKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("sender public key is not ECDSA")
	}
	hash := CalculateTxHash(tx)
	sigBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return false, fmt.Errorf("invalid transaction signature encoding: %v", err)
	}
	if len(sigBytes)%2 != 0 {
		return false, fmt.Errorf("invalid signature length")
	}
	r := new(big.Int).SetBytes(sigBytes[:len(sigBytes)/2])
	s := new(big.Int).SetBytes(sigBytes[len(sigBytes)/2:])
	return ecdsa.Verify(pubKey, hash, r, s), nil
}

func PublicKeyToAddress(pubKey []byte) string {
	hash := sha256.Sum256(pubKey)
	address := hash[len(hash)-20:]
	return "0x" + hex.EncodeToString(address)
}

// ImportWallet creates a wallet from a given private key (hex string)
func ImportWallet(privateKeyHex string) (*Wallet, error) {
	privBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key hex: %v", err)
	}
	privKey, err := x509.ParseECPrivateKey(privBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %v", err)
	}
	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  pubKeyBytes,
		Balance:    0,
		Staked:     0,
	}, nil
}

// NewWalletWithMnemonic creates a new wallet along with a 12-word BIP39 mnemonic.
func NewWalletWithMnemonic() (*Wallet, string, error) {
	// Generate 128 bits of entropy (produces a 12-word mnemonic)
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate entropy: %v", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate mnemonic: %v", err)
	}

	w, err := NewWalletFromMnemonic(mnemonic)
	if err != nil {
		return nil, "", err
	}

	return w, mnemonic, nil
}

// NewWalletFromMnemonic derives a wallet from a BIP39 mnemonic.
func NewWalletFromMnemonic(mnemonic string) (*Wallet, error) {
	seed := bip39.NewSeed(mnemonic, "") // No password for simplicity

	// For true deterministic derivation, we hash the seed to get exactly 32 bytes
	// and use that as the private key scalar directly
	hasher := sha256.New()
	hasher.Write(seed)
	keyBytes := hasher.Sum(nil)

	// Create private key from the deterministic bytes
	curve := elliptic.P256()
	privKey := new(ecdsa.PrivateKey)
	privKey.PublicKey.Curve = curve
	privKey.D = new(big.Int).SetBytes(keyBytes)

	// Validate that D is within the curve's valid range
	if privKey.D.Cmp(curve.Params().N) >= 0 {
		// If it's too large, reduce it modulo N
		privKey.D.Mod(privKey.D, curve.Params().N)
	}

	// Derive the public key
	privKey.PublicKey.X, privKey.PublicKey.Y = curve.ScalarBaseMult(keyBytes)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %v", err)
	}

	address := PublicKeyToAddress(publicKeyBytes)
	fmt.Printf("[DEBUG] Backend Wallet Derived Address: %s\n", address)

	return &Wallet{
		PrivateKey: privKey,
		PublicKey:  publicKeyBytes,
		Balance:    0,
		Staked:     0,
	}, nil
}
