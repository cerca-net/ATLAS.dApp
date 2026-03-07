# 📡 Phase 5: Transaction Broadcasting & Network Consensus Fixes

**Date:** January 30, 2026  
**Status:** ✅ Completed

---

## **1. Executive Summary**

This phase focused on resolving critical issues preventing successful transaction broadcasting and block consensus in a multi-node environment. The primary blockers were identified as non-deterministic treasury wallet generation and incorrect transaction address formatting. These issues caused peer nodes to reject valid transactions and blocks, preventing network synchronization.

We successfully implemented fixes that ensure consistent wallet derivation across all nodes and validated the full lifecycle of a transaction from Faucet Request -> Block Production -> Network Propagation -> State Update.

---

## **2. Root Cause Analysis**

### **A. Treasury Wallet Mismatch**
*   **Symptom:** Peer nodes rejected blocks containing Faucet transactions with the error `insufficient funds`.
*   **Cause:** The Treasury Wallet was being initialized using `ecdsa.GenerateKey` with a `bytes.Reader` seeded from the mnemonic. while intended to be deterministic, the internal consumption of entropy by `GenerateKey` varied or was not strictly deterministic across different execution environments or runs, resulting in different nodes deriving different Treasury Addresses (e.g., Node 1 derived `0x850e...` vs Node 3 derived `0xd28f...`).
*   **Impact:** Since the Genesis state assigned 1 Billion tokens to *one* specific address, nodes with a different derived address saw the Treasury as having 0 balance.

### **B. Address Formatting**
*   **Symptom:** Faucet transactions failed validation with `invalid recipient address length`.
*   **Cause:** The system expects a full 40-character hex string (excluding `0x` prefix) for addresses. The initial test requests used truncated or incorrect formats.
*   **Impact:** Transactions were rejected before entering the mempool.

---

## **3. Technical Implementation**

### **3.1. Deterministic Wallet Derivation**
We updated `pkg/wallet/wallet.go` to use a strictly deterministic method for key generation. Instead of relying on `ecdsa.GenerateKey`'s entropy consumption, we now hash the BIP39 seed using SHA-256 to produce a fixed 32-byte scalar, which is used as the private key `D` value.

```go
// NewWalletFromMnemonic derives a wallet from a BIP39 mnemonic.
func NewWalletFromMnemonic(mnemonic string) (*Wallet, error) {
    seed := bip39.NewSeed(mnemonic, "") 

    // Hash the seed to get exactly 32 bytes
    hasher := sha256.New()
    hasher.Write(seed)
    keyBytes := hasher.Sum(nil)

    // Create private key from the deterministic bytes
    curve := elliptic.P256()
    privKey := new(ecdsa.PrivateKey)
    privKey.PublicKey.Curve = curve
    privKey.D = new(big.Int).SetBytes(keyBytes)
    
    // ... derive public key ...
}
```

### **3.2. Faucet Request Validation**
We verified the correct payload format for the Faucet API:
```json
POST /faucet
{
    "address": "0x1234567890abcdef1234567890abcdef12345678"
}
```
*Note: The address must be a valid 40-character hex string (20 bytes).*

---

## **4. Verification & Testing**

We conducted a live 3-node network test:

| Node | Role | Status | Outcome |
| :--- | :--- | :--- | :--- |
| **Node 3** | Creator | Validator | Forged Block 1 containing Faucet Tx (Hash: `c970...`) |
| **Node 2** | Peer | Validator | Received Block 1 via broadcast; Address Balance updated to `1000` |
| **Node 1** | Peer | Validator | Joined after broadcast; Pending sync (Requires restart/catchup) |

**Success Criteria Met:**
- [x] All nodes derive the exact same Treasury Address (`0xd28f2b1294f15d29229016bff098fe1a9cfede16`).
- [x] Faucet transaction accepted and mined.
- [x] Block broadcasted to connected peers.
- [x] Peer state updated correctly.

---

## **5. Next Steps**

1.  **Node Synchronization:** Ensure `Node 1` (or any late-joining node) successfully runs the `Chain Sync` process to fetch missed blocks.
2.  **Automated Testing:** Incorporate this multi-node faucet test into the CI/CD pipeline.
3.  **Documentation:** Update the "Running a Network" guide with these findings.
