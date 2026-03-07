# Multi-Node Network Simulation

This setup allows you to run two independent blockchain nodes on your local machine to simulate a decentralized network.

## Prerequisites
- Windows OS
- Go installed
- `curl` or any HTTP client (for manual peer connection)

## How to Start
1. Run `scripts\start_network.bat`
2. Two command windows will open:
   - **Node 1**: Genesis/Faucet Node (Ports 8000/8080)
   - **Node 2**: User Node (Ports 8001/8081)

## Connecting the Nodes

### Method 1: Using the Connection Script
After both nodes are running:
1. Copy the **full multiaddress** from Node 1's console. It looks like:
   ```
   /ip4/127.0.0.1/tcp/8000/p2p/12D3KooWLYq8QetDBFhFppY9A47V81ZH5B6cwLmcJtG3NFdcmF1C
   ```
2. Run:
   ```cmd
   scripts\connect_nodes.bat "/ip4/127.0.0.1/tcp/8000/p2p/12D3KooWLYq8..."
   ```

### Method 2: Using curl (or any HTTP client)
```powershell
curl -X POST http://localhost:8081/connect-peer -H "Content-Type: application/json" -d '{\"peer_address\": \"/ip4/127.0.0.1/tcp/8000/p2p/12D3KooWLYq8...\"}'
```

### Verifying Connection
Check for success messages in the node consoles:
- `[P2P] Successfully connected to peer`
- `[P2P] Stream handler triggered from peer`

Or query the API:
```powershell
curl http://localhost:8081/peers
```

## API Endpoints
- **Node 1**: http://localhost:8080/status
- **Node 2**: http://localhost:8081/status
- **List Peers**: GET `/peers`
- **Connect Peer**: POST `/connect-peer` with body `{"peer_address": "<multiaddress>"}`

## Troubleshooting
- If nodes don't connect, ensure both are running and using different ports
- Check that the multiaddress is copied exactly (including the `/p2p/` part)
- Verify no firewall is blocking local connections
