# 🔐 Theory: Deterministic Wallet Generation in ATLAS

**Version:** 1.0  
**Date:** January 30, 2026

---

## **1. The Problem Space**

In a decentralized network, specific system-level wallets (like the **Treasury Wallet** or **Genesis Validators**) must be identical across all nodes. If Node A and Node B derive different addresses for the Treasury, they will have different Genesis states (Node A credits Treasury A, Node B credits Treasury B).

When Node A receives a transaction signed by Treasury B, it checks its state for Treasury B's balance. Since Node A only knows about Treasury A, it sees 0 balance and rejects the transaction.

## **2. The Challenge with Standard ECDSA**

The standard Go `crypto/ecdsa` library's `GenerateKey` function is designed for security, not cross-device determinism. It consumes entropy from a `io.Reader`.

```go
// Conventional (Non-Deterministic) Approach
seed := bip39.NewSeed(mnemonic, "")
reader := bytes.NewReader(seed) // Hoping entropy consumption is linear
key, _ := ecdsa.GenerateKey(curve, reader)
```

While `bytes.NewReader` provides deterministic output stream, the internal implementation of `GenerateKey` (specifically rejection sampling for the private scalar `d`) is not guaranteed to consume a fixed amount of bytes across different platforms or versions. If `d` is rejected (because it's >= N), it reads more bytes. This small variability breaks determinism.

## **3. The ATLAS Solution: Hash-To-Scalar**

To guarantee 100% determinism independent of the underlying ECDSA implementation, we moved to a **Hash-To-Scalar** approach.

### **Mechanism:**
1.  **Seed Generation:** `BIP39 Mnemonic` -> `512-bit Seed` (Standard).
2.  **Entropy Compression:** We hash the 512-bit seed using `SHA-256` to produce a fixed `256-bit (32-byte)` digest.
3.  **Direct Scalar Construction:** We use this 32-byte digest *directly* as the private key scalar `D`.
4.  **Curve Validation:** We ensure `D < N` (handling the rare edge case by modulo reduction if necessary).
5.  **Public Key Derivation:** We derive `(X, Y) = G * D`.

### **Implementation Details:**

```go
func NewWalletFromMnemonic(mnemonic string) (*Wallet, error) {
    seed := bip39.NewSeed(mnemonic, "")
    
    // 1. Hash seed to get fixed 32 bytes
    hasher := sha256.New()
    hasher.Write(seed)
    keyBytes := hasher.Sum(nil) // Strictly 32 bytes
    
    // 2. Construct Private Key directly
    privKey := new(ecdsa.PrivateKey)
    privKey.PublicKey.Curve = elliptic.P256()
    privKey.D = new(big.Int).SetBytes(keyBytes)
    
    // 3. Ensure D is valid
    if privKey.D.Cmp(curve.Params().N) >= 0 {
        privKey.D.Mod(privKey.D, curve.Params().N)
    }
    
    // ... derive public key
}
```

## **4. Benefits**

1.  **Platform Independence:** Works identically on Windows, Linux, macOS, and WASM.
2.  **Version Independence:** Immune to changes in Go's `crypto/rand` or `ecdsa` internals.
3.  **Verifiability:** Any developer can manually verify the derivation path using standard tools (SHA256).

This change ensures the **Treasury Address** is consistently `0xd28f2b1294f15d29229016bff098fe1a9cfede16` across the entire ATLAS network.
