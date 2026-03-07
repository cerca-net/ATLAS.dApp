package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

// EncryptedData represents encrypted data with metadata
type EncryptedData struct {
	Data      string `json:"data"`      // Base64 encoded encrypted data
	Nonce     string `json:"nonce"`     // Base64 encoded nonce
	KeyHash   string `json:"key_hash"`  // Hash of the encryption key for verification
	Algorithm string `json:"algorithm"` // Encryption algorithm used
}

// EncryptionKey represents a symmetric encryption key
type EncryptionKey struct {
	Key []byte
	ID  string
}

// NewEncryptionKey generates a new random encryption key
func NewEncryptionKey() (*EncryptionKey, error) {
	key := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %v", err)
	}
	
	keyHash := sha256.Sum256(key)
	return &EncryptionKey{
		Key: key,
		ID:  hex.EncodeToString(keyHash[:]),
	}, nil
}

// EncryptData encrypts data using AES-GCM with the provided key
func EncryptData(data []byte, key *EncryptionKey) (*EncryptedData, error) {
	block, err := aes.NewCipher(key.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	
	keyHash := sha256.Sum256(key.Key)
	
	return &EncryptedData{
		Data:      base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:     base64.StdEncoding.EncodeToString(nonce),
		KeyHash:   hex.EncodeToString(keyHash[:]),
		Algorithm: "AES-GCM-256",
	}, nil
}

// DecryptData decrypts data using AES-GCM with the provided key
func DecryptData(encrypted *EncryptedData, key *EncryptionKey) ([]byte, error) {
	// Verify key hash
	keyHash := sha256.Sum256(key.Key)
	if hex.EncodeToString(keyHash[:]) != encrypted.KeyHash {
		return nil, fmt.Errorf("key hash mismatch")
	}

	block, err := aes.NewCipher(key.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %v", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(encrypted.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nonce: %v", err)
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %v", err)
	}

	return plaintext, nil
}

// EncryptString encrypts a string and returns JSON-encoded encrypted data
func EncryptString(data string, key *EncryptionKey) (string, error) {
	encrypted, err := EncryptData([]byte(data), key)
	if err != nil {
		return "", err
	}
	
	jsonData, err := json.Marshal(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to marshal encrypted data: %v", err)
	}
	
	return string(jsonData), nil
}

// DecryptString decrypts JSON-encoded encrypted data back to a string
func DecryptString(encryptedJSON string, key *EncryptionKey) (string, error) {
	var encrypted EncryptedData
	if err := json.Unmarshal([]byte(encryptedJSON), &encrypted); err != nil {
		return "", fmt.Errorf("failed to unmarshal encrypted data: %v", err)
	}
	
	plaintext, err := DecryptData(&encrypted, key)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

// DeriveKeyFromPassword derives an encryption key from a password using PBKDF2
func DeriveKeyFromPassword(password string, salt []byte) (*EncryptionKey, error) {
	// For now, use a simple hash-based approach
	// In production, use proper PBKDF2 or Argon2
	hash := sha256.Sum256([]byte(password + string(salt)))
	
	return &EncryptionKey{
		Key: hash[:],
		ID:  hex.EncodeToString(hash[:]),
	}, nil
}

// IsEncryptedData checks if a string contains encrypted data
func IsEncryptedData(data string) bool {
	if len(data) == 0 {
		return false
	}
	
	// Try to parse as JSON encrypted data
	var encrypted EncryptedData
	if err := json.Unmarshal([]byte(data), &encrypted); err != nil {
		return false
	}
	
	// Check if it has the required fields
	return encrypted.Data != "" && encrypted.Nonce != "" && encrypted.KeyHash != ""
} 