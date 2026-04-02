# ATLAS CercaChain - V1 MVP

The **ATLAS CercaChain** is a Layer 1 Proof-of-Stake (PBFT) blockchain written in Go with an integrated native Stack-based Virtual Machine (CercaVM). It powers the **CercaEnd** Flutter Web Application natively, handling localized Token logic, Validator Staking, Marketplace Data Unit commerce, and fully on-chain Governance voting. 

![CercaEnd Prototype](https://github.com/CercaChain/CercaChain)

---

## ⚡ Deployment Instructions

As of V1 MVP, the network relies on deterministic bootstrap validators (`genesis_validator.key`) and predefined consensus configurations (`genesis.json`).

The easiest way to boot the ecosystem is using **Docker Compose**, which securely encapsulates the Host OS anomalies (e.g. SQLite CGO thread handling) inside an Alpine Linux instance.

### Option 1: Docker Compose (Recommended)
This method spins up both the **Node Validator Backend** and the **Flutter Web Server Frontend** instantly. Make sure you have Docker Desktop running.

```bash
# Start the full ecosystem in the background
docker-compose up -d

# View logs dynamically
docker-compose logs -f
```
**Access The Network:**
- **Frontend App**: [http://localhost/](http://localhost/)
- **Backend API**: [http://localhost:8080/](http://localhost:8080/)

---

### Option 2: Local Windows Native Devnet

If you are modifying the raw Go code or Flutter services locally and want active recompilation outputs:

**1. Boot the Blockchain Nodes (Powershell Required)**
This spawns 2 local instances dynamically synced to one another using the `start_devnet.ps1` script inside the workspace.
```powershell
cd ATLAS.BC0.0.1
# Compile the L1 Daemon
go build -o build\atlas-node.exe cmd\main.go
# Boot Node 1 (Validator) and Node 2 (Observer)
powershell.exe -ExecutionPolicy Bypass -File .\start_devnet.ps1
```

**2. Boot the Frontend**
```powershell
cd cercaend
# Fetch Dependencies
flutter pub get
# Preview in a Chrome debugger wrapper 
flutter run -d chrome --web-port=3000
```
Then navigate to `http://localhost:3000` to interact with the Local DB.

---

## 🏗 Submodules

- **`ATLAS.BC0.0.1`**: Core Go blockchain architecture (`sqlite`, `gnark` bypassed for V1, native `secp256k1` crypto logic). Has the `generate_genesis_keys` builder script.
- **`cercaend`**: The native frontend UI written with FlutterFlow and Dart Services. Includes all UI endpoints safely verifying `mounted` asynchronous UI callbacks.

---

## 🛡 Network Architecture Notes (V1)
- Privacy verification scaling overhead using `gnark` Groth16 Trusted Setups has been deferred to **V2**. Currently, the verification logic skips circuit constraints and utilizes basic data assertions. 
- The Treasury account pre-allocates `1,000,000,000 TCOIN` securely dictated by the PBFT `genesis.json` configuration logic payload.
- State is natively mapped and snapshotted via the `blockchain.db` SQLite container volume. Shutting down nodes will pause states without deleting the cache via WAL bindings.

Enjoy the Cerca Ecosystem!
