# ATLAS CercaChain: Multi-Device & P2P LAN Testing Strategy
**Version 1.0 · May 2026**

This document details the strategy for deploying the host computer as the central server-seed node (Treasury Node) and testing ecosystem operations from external client devices (mobile/tablet/secondary PCs) over a Local Area Network (LAN).

---

## 1. Network Topology Overview

We will test two primary deployment topologies on the local network:

1. **Remote Client Mode (Light Client)**: External devices run only the Flutter frontend (compiled for Web/Android/iOS) and route API requests directly to the host machine's Go REST API.
2. **Full P2P Node Mode**: A secondary PC runs a separate ATLAS Go node daemon, which discovers the host's seed node via libp2p (mDNS or direct dial). The client frontend on the secondary PC routes to its own local node.

```mermaid
graph TD
    subgraph Host Computer [Host Machine: Server-Seed Node]
        NodeA[ATLAS Go Node A<br/>API: :8080 | P2P: :8000]
        WebA[Flutter Web Server<br/>Port: :8082]
        AdminA[React Admin Panel<br/>Port: :3001]
        DB_A[(SQLite State DB)]
        NodeA <--> DB_A
    end

    subgraph WiFiRouter [Local Wi-Fi Router / LAN]
        R[Wi-Fi Access Point]
    end

    subgraph ClientDevice [External Client Device: Phone/Tablet]
        AppB[CercaEnd Flutter Client<br/>Browser / App]
    end

    subgraph PeerComputer [Secondary PC: Peer Node]
        NodeB[ATLAS Go Node B<br/>API: :8081 | P2P: :8001]
        AppC[CercaEnd Flutter Client<br/>Browser]
        DB_B[(SQLite State DB)]
        NodeB <--> DB_B
        AppC <--> NodeB
    end

    %% Network Connections
    NodeA <--> R
    WebA <--> R
    AdminA <--> R
    AppB <--> R
    NodeB <--> R

    %% Interactions
    AppB -- HTTP API Requests --> NodeA
    NodeB -- libp2p mDNS / GossipSub --> NodeA
```

---

## 2. Environment Configuration & Setup

### Prerequisite 1: Retrieve Host Computer LAN IP
To allow external devices to locate the host machine, identify its IPv4 address on the local network:
1. Open PowerShell on the host computer.
2. Run: `ipconfig`
3. Locate the active network adapter (e.g., *Wireless LAN adapter Wi-Fi* or *Ethernet adapter*).
4. Note the IPv4 Address (e.g., `192.168.1.125`). We will refer to this as `<HOST_LAN_IP>`.

### Prerequisite 2: Configure Host Firewall
Windows Firewall blocks inbound traffic on custom ports by default. We must allow incoming connections for the required services:
1. Open PowerShell as Administrator on the host machine.
2. Run the following commands to open ports for the REST API, P2P network, and Web App host:
   ```powershell
   # Open API port for external client connections
   New-NetFirewallRule -DisplayName "ATLAS Go API" -Direction Inbound -Action Allow -Protocol TCP -LocalPort 8080,8081
   
   # Open P2P ports for secondary Go node syncing
   New-NetFirewallRule -DisplayName "ATLAS P2P Host" -Direction Inbound -Action Allow -Protocol TCP -LocalPort 8000,8001
   
   # Open ports for serving the Frontend client and Admin Panel
   New-NetFirewallRule -DisplayName "ATLAS HTTP Serving" -Direction Inbound -Action Allow -Protocol TCP -LocalPort 8082,3001
   ```

### Prerequisite 3: Check Router Access Point (AP) Isolation
Ensure that **AP Isolation / Client Isolation** is disabled in your Wi-Fi router settings. If enabled, the router prevents wireless clients from communicating with each other and the host computer.

---

## 3. Server-Seed Node Deployment (Host Machine)

We will use the primary validator configuration for the host. 

### Step 1: Start the Primary Go Node & Admin Panel
To run the server-seed node natively (using SQLite persistence) and boot the React Admin Panel, execute the existing script on the host:
```powershell
.\start_admin.bat
```
*   **P2P port**: `8000` (launches as a validator using keys at `.data_admin/nodekey.priv`).
*   **API port**: `8080`.
*   **Admin Panel**: `http://localhost:3001` (communicates with the node via local loopback).

### Step 2: Serve the Flutter Web Client to the LAN
To allow external devices (phones/tablets) to load the web interface, we must run the Flutter web server listening on all network interfaces (`0.0.0.0`) and pointing to our external LAN API:
1. Open a new PowerShell terminal.
2. Navigate to the `cercaend` directory:
   ```powershell
   cd cercaend
   ```
3. Run the development server with host bindings and the target API URL:
   ```powershell
   flutter run -d web-server --web-hostname 0.0.0.0 --web-port 8082 --dart-define=BLOCKCHAIN_API_URL=http://<HOST_LAN_IP>:8080
   ```
   *Replace `<HOST_LAN_IP>` with your actual IPv4 address (e.g., `192.168.1.125`).*

---

## 4. Launching the App on External Devices

### Web Browser (Easiest Method)
1. Ensure the external device (phone, tablet, or secondary PC) is connected to the same Wi-Fi network.
2. Open a browser on the device and navigate to:
   ```http
   http://<HOST_LAN_IP>:8082
   ```
3. The browser will download the compiled Flutter application. When actions (like wallet creation, post creation, or order placements) are triggered in the UI, the browser will send HTTP requests to the host's REST API at `http://<HOST_LAN_IP>:8080`.

### Native Mobile Builds (Android/iOS)
If testing with native apps:
1. Compile the app with the target URL injected:
   ```powershell
   flutter build apk --release --dart-define=BLOCKCHAIN_API_URL=http://<HOST_LAN_IP>:8080
   ```
2. Install the generated APK on your Android device.

---

## 5. Peer-to-Peer Multi-Node Sync Testing
To test block replication, transaction propagation, and consensus across physical machines, set up a secondary node on another PC in the same LAN.

### Step 1: Auto-Discovery (mDNS)
By default, the libp2p engine has mDNS enabled under the service tag `"cercachain-mdns"`. 
1. Build or transfer the `ATLAS.BC0.0.1` folder to the secondary computer.
2. Launch the node on the secondary computer:
   ```powershell
   go run cmd/main.go --port 8000 --api 8080 --datadir ./data_peer1
   ```
3. The secondary node should automatically detect the host node on the local subnet and establish a connection.

### Step 2: Manual Fallback Connection
If your local router blocks multicast packets (preventing mDNS), perform a manual connection:
1. On the host computer, check the logs or call `GET /node-address` to get the host node's libp2p multiaddress:
   *Example:* `/ip4/<HOST_LAN_IP>/tcp/8000/p2p/QmYyQzSZyD5G5k...`
2. On the secondary computer, start the node with the environment variable set:
   ```powershell
   # Windows PowerShell
   $env:NODE1_MULTIADDR="/ip4/<HOST_LAN_IP>/tcp/8000/p2p/QmYyQzSZyD5G5k..."
   go run cmd/main.go --port 8000 --api 8080 --datadir ./data_peer1
   ```
3. Verify connection by querying the peer list endpoint on either machine:
   ```http
   GET http://localhost:8080/peers
   ```

---

## 6. Activity Monitoring Checklist

You can track all network operations in real-time from the host computer using the following interfaces:

1. **React Admin Panel (`http://localhost:3001`)**:
   *   **Dashboard**: Monitor aggregate block height, active peers, and tx pool size.
   *   **Transactions**: View the incoming mempool queue and trace finalized blocks.
   *   **Node Control / Peers**: View detailed latency and address lists of all connected peer nodes.
2. **Command Line Console**:
   *   Observe the terminal stdout of `go run cmd/main.go`. Look for logs prefixed with `[P2P]`, `[API]`, and `[SYNC]`.
3. **API Logging Endpoint**:
   *   Fetch recent JSON logs by calling:
     ```http
     GET http://localhost:8080/node/logs?limit=100
     ```

---

## 7. Step-by-Step Scenario Verification

Use the following step-by-step scenarios to perform structural integration testing:

### Scenario 1: Identity & Faucet
*   **Objective**: Test account generation, secp256k1 key derivation, and faucet distribution.
*   **Steps**:
    1. On the external client device (`http://<HOST_LAN_IP>:8082`), navigate to the **User/Wallet** page.
    2. Click **Create Wallet** (generates mnemonic and derives keypair locally).
    3. Click **Request Faucet Tokens**.
    4. Confirm that 1,000 TCOIN is credited to the wallet balance.
    5. Check the **React Admin Panel** under the *Treasury* tab to verify the transaction was successfully recorded in block storage.

### Scenario 2: Peer-to-Peer Transfers
*   **Objective**: Verify ECDSA transaction signing, nonce tracking, and mempool validation.
*   **Steps**:
    1. Generate two separate wallets on two external devices (Client A and Client B).
    2. Request faucet tokens for Client A.
    3. Copy Client B's address.
    4. In Client A's wallet page, input Client B's address, set the amount to `200 TCOIN`, and click **Send**.
    5. Verify that the transaction enters the host node's mempool, is packaged into the next block, and that balances update correctly on both clients.

### Scenario 3: Social Interactions & Energy Physics
*   **Objective**: Verify Lamport Clock logic, CercaVM state-execution, and object decay rules.
*   **Steps**:
    1. From Client A, create a new Post in the **Social Feed**.
    2. From Client B, navigate to the Feed, locate Client A's post, and click **Like** (verifies profile verification).
    3. Comment on the post (verifies that 2 TCOIN is deducted from Client B and added to the post energy pool).
    4. Tip the post `50 TCOIN` (transfers funds directly into the post's on-chain energy balance).
    5. Verify that the post's `InfluenceScore` rises in the feed ranking.
    6. Let the post age (or simulate decay) to verify that objects with energy `< 50` become `"fossilized"` and require revival payments.

### Scenario 4: E-Commerce Escrow Flow
*   **Objective**: Test the multi-party Escrow System contract (`vm.MarketplaceContractAddress`).
*   **Steps**:
    1. Client A (Seller) registers an item in the Catalog.
    2. Client B (Buyer) places an order (funds are moved on-chain into the Escrow System Contract storage).
    3. Confirm that Client B's balance decreases, but Client A's balance does *not* increase yet.
    4. Client A marks the order as shipped.
    5. Client B clicks **Release Funds** in their orders dashboard.
    6. Verify that the escrow contract executes the transfer, adding the TCOIN to Client A's balance (minus the system transaction fee).

---

## 8. Summary Action Plan

To initiate this testing phase:
1. Ensure your host machine and mobile/external devices are on the same Wi-Fi subnet.
2. Allow ports `8080`, `8000`, `8082`, and `3001` through your host firewall.
3. Start the node and admin dashboard using `.\start_admin.bat`.
4. Serve the web frontend with host-bound parameters (`0.0.0.0` and `--dart-define=BLOCKCHAIN_API_URL=http://<HOST_LAN_IP>:8080`).
5. Open browser access on external devices and execute scenarios.
