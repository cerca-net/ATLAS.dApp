package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// handleCreate Snapshot creates a new blockchain state snapshot
func (api *APIServer) handleCreateSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if api.fastSyncManager == nil {
		http.Error(w, "Fast sync manager not initialized", http.StatusServiceUnavailable)
		return
	}

	metadata, err := api.fastSyncManager.CreateSnapshot()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"metadata": metadata,
	})
}

// handleGetLatestSnapshot returns information about the latest snapshot
func (api *APIServer) handleGetLatestSnapshot(w http.ResponseWriter, r *http.Request) {
	if api.fastSyncManager == nil {
		http.Error(w, "Fast sync manager not initialized", http.StatusServiceUnavailable)
		return
	}

	snapshotPath, metadata, err := api.fastSyncManager.GetLatestSnapshot()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"path":     snapshotPath,
		"metadata": metadata,
	})
}

// handleLoadSnapshot loads a snapshot from file
func (api *APIServer) handleLoadSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if api.fastSyncManager == nil {
		http.Error(w, "Fast sync manager not initialized", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		SnapshotPath string `json:"snapshot_path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.SnapshotPath == "" {
		// Load latest snapshot if no path specified
		snapshotPath, _, err := api.fastSyncManager.GetLatestSnapshot()
		if err != nil {
			http.Error(w, "No snapshots available", http.StatusNotFound)
			return
		}
		req.SnapshotPath = snapshotPath
	}

	if err := api.fastSyncManager.LoadSnapshot(req.SnapshotPath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Snapshot loaded successfully",
		"path":    req.SnapshotPath,
	})
}

// handleGetPeerStatus returns the status of all known peers
func (api *APIServer) handleGetPeerStatus(w http.ResponseWriter, r *http.Request) {
	if api.peerManager == nil {
		http.Error(w, "Peer manager not initialized", http.StatusServiceUnavailable)
		return
	}

	peers := api.peerManager.GetKnownPeers()

	// Format peer information
	peerStatuses := make([]map[string]interface{}, len(peers))
	for i, peer := range peers {
		peerStatuses[i] = map[string]interface{}{
			"peer_id":      peer.PeerID[:16] + "...",
			"is_validator": peer.IsValidator,
			"stake":        peer.Stake,
			"last_seen":    peer.LastSeen.Format(time.RFC3339),
			"multiaddrs":   peer.Multiaddrs,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"peer_count": len(peers),
		"peers":      peerStatuses,
	})
}

// handleReconnectPeers triggers reconnection to all known peers
func (api *APIServer) handleReconnectPeers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if api.peerManager == nil {
		http.Error(w, "Peer manager not initialized", http.StatusServiceUnavailable)
		return
	}

	if err := api.peerManager.ReconnectToPeers(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Reconnection attempt completed",
	})
}

// handleGetValidatorPeers returns list of validator peers
func (api *APIServer) handleGetValidatorPeers(w http.ResponseWriter, r *http.Request) {
	if api.peerManager == nil {
		http.Error(w, "Peer manager not initialized", http.StatusServiceUnavailable)
		return
	}

	validators := api.peerManager.GetValidators()

	// Format validator information
	validatorInfos := make([]map[string]interface{}, len(validators))
	for i, validator := range validators {
		validatorInfos[i] = map[string]interface{}{
			"peer_id":    validator.PeerID[:16] + "...",
			"stake":      validator.Stake,
			"last_seen":  validator.LastSeen.Format(time.RFC3339),
			"multiaddrs": validator.Multiaddrs,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":         true,
		"validator_count": len(validators),
		"validators":      validatorInfos,
	})
}
