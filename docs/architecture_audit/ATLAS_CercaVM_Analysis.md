# ATLAS Smart Contract Architecture: The CercaVM

This document details the architecture and operational mechanics of the ATLAS Virtual Machine (`CercaVM`), the execution environment responsible for smart contracts, decentralized commerce, and state mutations within the network.

---

## 1. The CercaVM Architecture
Unlike the Ethereum Virtual Machine (EVM) which executes compiled bytecode, the `CercaVM` is a lightweight, stack-based virtual machine written natively in Go. It is optimized specifically for social-commerce and high-throughput operations.

### Stack & Memory Model
The VM utilizes a dual-stack design to handle data efficiently without excessive casting:
*   **Int64 Stack:** Handles all mathematical operations, boolean logic, and numeric storage.
*   **String Stack:** Specifically designed to handle blockchain addresses (e.g., wallet public keys) and text identifiers (e.g., `ORDER_123`).
*   **Persistent Storage:** Accessed via `SSTORE` (int) and `SSTORE_S` (string), writing directly to the `StateManager`'s SQLite database.

### Gas Metering
To prevent infinite loops and spam, the VM implements a rigid Gas Metering system. Every operation is hardcoded with a specific gas cost.
*   `PUSH`: 3 Gas
*   `MUL`: 5 Gas
*   `CALL`: 10 Gas
*   `TRANSFER`: 21 Gas
*   `MINT`: 30 Gas

If an execution context (`ExecutionContext`) exceeds its `GasLimit`, the VM instantly halts with an "Out of Gas" error, reverting the transaction but keeping the spent fee.

---

## 2. The Permissioned Contract Tiering System
A major vulnerability in public blockchains is the deployment of malicious or unverified smart contracts. ATLAS solves this via a rigid Permissioned Tiering System.

*   **System Contracts:** (e.g., Marketplace, Token, Staking). These are pre-approved by the protocol and have unrestricted access to system opcodes like `MINT`.
*   **Governance Contracts:** Pre-approved contracts that handle treasury and parameter votes.
*   **Voting Contracts:** Standardized templates for decentralized voting.
*   **Custom Contracts:** Users *can* deploy custom JSON-based smart contracts (`TxTypeDeploy`), but they are flagged as `ContractTypeCustom`. A custom contract **cannot execute** until it is explicitly approved (`ApproveCustomContract`), acting as a safety net against exploits.

---

## 3. The "Helper" Pattern & System Contracts
To achieve maximum scalability, the core features of the ATLAS network are built as **System Contracts** using a "Helper" pattern in Go (e.g., `MarketplaceContractHelper`). 

Instead of interpreting bytecode line-by-line, when the network detects a call to the `MarketplaceContractAddress`, it intercepts it and runs pre-compiled, native Go code. 

### The E-Commerce Escrow (Marketplace)
The `marketplace_contract.go` file reveals a highly robust, on-chain escrow system designed specifically to bridge crypto with physical/digital goods.
1.  **Fund (`createOrder`):** A buyer locks `TCOIN` into the contract's escrow. The state is updated to `EscrowStatusFunded`.
2.  **Delivery:** The seller ships the item or delivers the digital good.
3.  **Completion (`releaseFunds`):** The buyer confirms delivery. The contract deducts a `0.2%` (20 basis points) network fee and releases the funds to the seller.
4.  **Dispute (`raiseDispute`):** If the item is not delivered, either party can freeze the funds. It enters `EscrowStatusDisputed`.
5.  **Arbitration (`resolveDispute`):** The Admin Panel (via a privileged key) evaluates the real-world evidence and forcefully refunds the buyer or pays the seller. 

---

## 4. Zero-Knowledge Readiness (V2)
The `vm.go` file contains an active placeholder for Zero-Knowledge Proofs: `VerifyZKProof`.

Currently, the code explicitly logs: `WARNING: VerifyZKProof is MOCKED for V1. Returns true unconditionally.` 

This confirms that the fundamental VM architecture is fully prepared to integrate `gnark` or `libsnark` circuits in the Phase 2 rollout. Once enabled, the VM will be capable of processing entirely private transactions and obfuscated smart contract executions while maintaining cryptographic integrity.
