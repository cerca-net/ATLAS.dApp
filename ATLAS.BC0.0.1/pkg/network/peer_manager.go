package network

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

// PeerInfo holds information about a known peer
type PeerInfo struct {
	PeerID      string    `json:"peer_id"`
	Multiaddrs  []string  `json:"multiaddrs"`
	LastSeen    time.Time `json:"last_seen"`
	IsValidator bool      `json:"is_validator"`
	Stake       uint64    `json:"stake,omitempty"`
}

// PeerManager handles peer discovery, persistence, and reconnection
type PeerManager struct {
	p2pNode           *P2PNode
	knownPeers        map[string]*PeerInfo
	mu                sync.RWMutex
	persistPath       string
	autoReconnect     bool
	reconnectInterval time.Duration
}

// NewPeerManager creates a new peer manager
func NewPeerManager(p2pNode *P2PNode, dataDir string) *PeerManager {
	persistPath := filepath.Join(dataDir, "peers.json")

	pm := &PeerManager{
		p2pNode:           p2pNode,
		knownPeers:        make(map[string]*PeerInfo),
		persistPath:       persistPath,
		autoReconnect:     true,
		reconnectInterval: 30 * time.Second,
	}

	// Load persisted peers
	if err := pm.LoadPeers(); err != nil {
		log.Printf("⚠️  Failed to load persisted peers: %v", err)
	}

	// Start automatic peer persistence
	go pm.startPersistenceRoutine()

	// Start automatic reconnection if enabled
	if pm.autoReconnect {
		go pm.startReconnectionRoutine()
	}

	return pm
}

// AddPeer adds or updates a peer in the known peers list
func (pm *PeerManager) AddPeer(peerID peer.ID, isValidator bool, stake uint64) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	peerInfo := pm.p2pNode.Host.Peerstore().PeerInfo(peerID)

	multiaddrs := make([]string, len(peerInfo.Addrs))
	for i, addr := range peerInfo.Addrs {
		multiaddrs[i] = addr.String()
	}

	pm.knownPeers[peerID.String()] = &PeerInfo{
		PeerID:      peerID.String(),
		Multiaddrs:  multiaddrs,
		LastSeen:    time.Now(),
		IsValidator: isValidator,
		Stake:       stake,
	}

	log.Printf("📝 [PEER-MGR] Added peer: %s (validator: %v)", peerID.String()[:16]+"...", isValidator)
}

// UpdatePeerLastSeen updates the last seen timestamp for a peer
func (pm *PeerManager) UpdatePeerLastSeen(peerID peer.ID) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if info, exists := pm.knownPeers[peerID.String()]; exists {
		info.LastSeen = time.Now()
	}
}

// GetKnownPeers returns a list of all known peers
func (pm *PeerManager) GetKnownPeers() []*PeerInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	peers := make([]*PeerInfo, 0, len(pm.knownPeers))
	for _, peer := range pm.knownPeers {
		peers = append(peers, peer)
	}
	return peers
}

// GetValidators returns a list of known validator peers
func (pm *PeerManager) GetValidators() []*PeerInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	validators := make([]*PeerInfo, 0)
	for _, peer := range pm.knownPeers {
		if peer.IsValidator {
			validators = append(validators, peer)
		}
	}
	return validators
}

// SavePeers persists the known peers to disk
func (pm *PeerManager) SavePeers() error {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	data, err := json.MarshalIndent(pm.knownPeers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal peers: %v", err)
	}

	if err := ioutil.WriteFile(pm.persistPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write peers file: %v", err)
	}

	log.Printf("💾 [PEER-MGR] Saved %d peers to %s", len(pm.knownPeers), pm.persistPath)
	return nil
}

// LoadPeers loads persisted peers from disk
func (pm *PeerManager) LoadPeers() error {
	if _, err := os.Stat(pm.persistPath); os.IsNotExist(err) {
		log.Printf("ℹ️  [PEER-MGR] No persisted peers file found")
		return nil
	}

	data, err := ioutil.ReadFile(pm.persistPath)
	if err != nil {
		return fmt.Errorf("failed to read peers file: %v", err)
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	if err := json.Unmarshal(data, &pm.knownPeers); err != nil {
		return fmt.Errorf("failed to unmarshal peers: %v", err)
	}

	log.Printf("📂 [PEER-MGR] Loaded %d persisted peers", len(pm.knownPeers))
	return nil
}

// startPersistenceRoutine periodically saves peers to disk
func (pm *PeerManager) startPersistenceRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := pm.SavePeers(); err != nil {
			log.Printf("⚠️  [PEER-MGR] Failed to save peers: %v", err)
		}
	}
}

// ReconnectToPeers attempts to reconnect to all known peers
func (pm *PeerManager) ReconnectToPeers() error {
	pm.mu.RLock()
	knownPeers := make([]*PeerInfo, 0, len(pm.knownPeers))
	for _, peer := range pm.knownPeers {
		knownPeers = append(knownPeers, peer)
	}
	pm.mu.RUnlock()

	if len(knownPeers) == 0 {
		log.Printf("ℹ️  [PEER-MGR] No known peers to reconnect to")
		return nil
	}

	log.Printf("🔄 [PEER-MGR] Attempting to reconnect to %d peers...", len(knownPeers))

	successCount := 0
	for _, peerInfo := range knownPeers {
		// Check if already connected
		peerID, err := peer.Decode(peerInfo.PeerID)
		if err != nil {
			log.Printf("⚠️  [PEER-MGR] Invalid peer ID: %v", err)
			continue
		}

		// Skip if already connected
		if pm.p2pNode.Host.Network().Connectedness(peerID) == 1 { // Connected
			successCount++
			pm.UpdatePeerLastSeen(peerID)
			continue
		}

		// Try to reconnect using stored multiaddresses
		for _, addrStr := range peerInfo.Multiaddrs {
			addrInfo, err := peer.AddrInfoFromString(addrStr + "/p2p/" + peerInfo.PeerID)
			if err != nil {
				continue
			}

			if err := pm.p2pNode.Host.Connect(context.Background(), *addrInfo); err != nil {
				log.Printf("⚠️  [PEER-MGR] Failed to reconnect to %s: %v", peerInfo.PeerID[:16]+"...", err)
				continue
			}

			log.Printf("✅ [PEER-MGR] Reconnected to peer: %s", peerInfo.PeerID[:16]+"...")
			pm.UpdatePeerLastSeen(peerID)
			successCount++
			break
		}
	}

	log.Printf("📊 [PEER-MGR] Reconnection complete: %d/%d peers connected", successCount, len(knownPeers))
	return nil
}

// startReconnectionRoutine periodically attempts to reconnect to known peers
func (pm *PeerManager) startReconnectionRoutine() {
	ticker := time.NewTicker(pm.reconnectInterval)
	defer ticker.Stop()

	for range ticker.C {
		if err := pm.ReconnectToPeers(); err != nil {
			log.Printf("⚠️  [PEER-MGR] Reconnection failed: %v", err)
		}
	}
}

// SetAutoReconnect enables or disables automatic reconnection
func (pm *PeerManager) SetAutoReconnect(enabled bool) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.autoReconnect = enabled
}

// CleanupStalePeers removes peers that haven't been seen in a while
func (pm *PeerManager) CleanupStalePeers(maxAge time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	now := time.Now()
	removedCount := 0

	for peerID, peerInfo := range pm.knownPeers {
		if now.Sub(peerInfo.LastSeen) > maxAge {
			delete(pm.knownPeers, peerID)
			removedCount++
		}
	}

	if removedCount > 0 {
		log.Printf("🗑️  [PEER-MGR] Cleaned up %d stale peers", removedCount)
	}
}
