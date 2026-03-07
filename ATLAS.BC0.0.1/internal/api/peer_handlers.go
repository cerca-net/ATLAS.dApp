package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// POST /connect-peer - Connects to a remote peer using its multiaddress
func (api *APIServer) handleConnectPeer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PeerAddress string `json:"peer_address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PeerAddress == "" {
		http.Error(w, "peer_address is required", http.StatusBadRequest)
		return
	}

	// Parse the multiaddress
	maddr, err := multiaddr.NewMultiaddr(req.PeerAddress)
	if err != nil {
		http.Error(w, "Invalid multiaddress: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extract peer info from the multiaddress
	peerInfo, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		http.Error(w, "Failed to extract peer info: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Connect to the peer
	ctx := context.Background()
	if err := api.p2pNode.Host.Connect(ctx, *peerInfo); err != nil {
		http.Error(w, "Failed to connect to peer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("[API] Successfully connected to peer: %s", peerInfo.ID.String())

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Successfully connected to peer",
		"peer_id": peerInfo.ID.String(),
	})
}
