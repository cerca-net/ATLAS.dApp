package blockchain

import (
	"atlas-blockchain/pkg/block"
	"atlas-blockchain/pkg/config"
	"atlas-blockchain/pkg/network"
	"atlas-blockchain/pkg/transaction"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

// SyncStatus represents the current synchronization status
type SyncStatus string

const (
	SyncStatusIdle     SyncStatus = "idle"
	SyncStatusSyncing  SyncStatus = "syncing"
	SyncStatusComplete SyncStatus = "complete"
	SyncStatusFailed   SyncStatus = "failed"
	SyncStatusForked   SyncStatus = "forked"
)

// PendingRequest tracks a pending request waiting for a response
type PendingRequest struct {
	RequestID   string
	RequestType network.MessageType
	PeerID      peer.ID
	Timestamp   time.Time
	Response    chan interface{}
	Timeout     time.Duration
}

// ChainSyncManager handles blockchain synchronization
type ChainSyncManager struct {
	blockManager *BlockManager
	stateManager *StateManager
	p2pNode      *network.P2PNode
	config       *config.BlockchainConfig

	// Sync state
	status        SyncStatus
	mu            sync.RWMutex
	syncStartTime time.Time
	syncEndTime   time.Time
	blocksSynced  int64
	totalBlocks   int64

	// Chain validation
	genesisHash    string
	canonicalChain []string // Hashes of canonical chain

	// Message correlation
	pendingRequests map[string]*PendingRequest
	requestMutex    sync.RWMutex

	// Callbacks
	onSyncStart    func()
	onSyncProgress func(blocksSynced, totalBlocks int64)
	onSyncComplete func()
	onSyncFailed   func(error)
	onForkDetected func(forkHeight int64, canonical, forked []string)
}

// NewChainSyncManager creates a new chain synchronization manager
func NewChainSyncManager(blockManager *BlockManager, stateManager *StateManager, p2pNode *network.P2PNode, config *config.BlockchainConfig) *ChainSyncManager {
	csm := &ChainSyncManager{
		blockManager:    blockManager,
		stateManager:    stateManager,
		p2pNode:         p2pNode,
		config:          config,
		status:          SyncStatusIdle,
		genesisHash:     blockManager.GetLatestBlock().Hash, // Genesis block hash
		pendingRequests: make(map[string]*PendingRequest),
	}

	// Set up message handlers
	csm.setupMessageHandlers()

	// Start cleanup routine
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				csm.cleanupExpiredRequests()
			}
		}
	}()

	return csm
}

// setupMessageHandlers registers handlers for sync-related messages
func (csm *ChainSyncManager) setupMessageHandlers() {
	// Override the default message handler
	csm.p2pNode.HandleIncomingMessage = func(msg network.NetworkMessage) {
		switch msg.Type {
		case network.MsgTypeChainStatusRequest:
			csm.handleChainStatusRequest(msg)
		case network.MsgTypeBlockRequest:
			csm.handleBlockRequest(msg)
		case network.MsgTypeSyncRequest:
			csm.handleSyncRequest(msg)
		case network.MsgTypeChainStatusResponse:
			csm.handleChainStatusResponse(msg)
		case network.MsgTypeBlockResponse:
			csm.handleBlockResponse(msg)
		case network.MsgTypeSyncResponse:
			csm.handleSyncResponse(msg)
		case network.MsgTypeForkResolution:
			csm.handleForkResolution(msg)
		default:
			// Let other handlers process non-sync messages
			log.Printf("[SYNC] Unknown message type: %s", msg.Type)
		}
	}
}

// StartSync initiates chain synchronization with peers
func (csm *ChainSyncManager) StartSync(ctx context.Context) error {
	csm.mu.Lock()
	if csm.status == SyncStatusSyncing {
		csm.mu.Unlock()
		return fmt.Errorf("sync already in progress")
	}

	csm.status = SyncStatusSyncing
	csm.syncStartTime = time.Now()
	csm.blocksSynced = 0
	csm.mu.Unlock()

	log.Printf("🔄 [SYNC] Starting chain synchronization...")

	if csm.onSyncStart != nil {
		csm.onSyncStart()
	}

	// Get connected peers
	peers := csm.p2pNode.Host.Peerstore().Peers()
	if len(peers) == 0 {
		csm.setStatus(SyncStatusFailed)
		return fmt.Errorf("no peers available for synchronization")
	}

	// Request chain status from all peers
	chainStatuses := make([]network.ChainStatusMessage, 0)
	for _, peerID := range peers {
		if peerID == csm.p2pNode.Host.ID() {
			continue // Skip self
		}

		status, err := csm.requestChainStatus(ctx, peerID)
		if err != nil {
			log.Printf("⚠️ [SYNC] Failed to get chain status from peer %s: %v", peerID.String(), err)
			continue
		}
		chainStatuses = append(chainStatuses, status)
	}

	if len(chainStatuses) == 0 {
		csm.setStatus(SyncStatusFailed)
		return fmt.Errorf("no valid chain statuses received from peers")
	}

	// Find the highest chain
	highestChain := csm.findHighestChain(chainStatuses)
	log.Printf("📊 [SYNC] Highest chain found: height=%d, hash=%s", highestChain.Height, highestChain.LatestHash)

	// Check if we need to sync
	currentHeight := int64(csm.blockManager.GetBlockHeight())
	if currentHeight >= highestChain.Height {
		log.Printf("✅ [SYNC] Already up to date (local: %d, network: %d)", currentHeight, highestChain.Height)
		csm.setStatus(SyncStatusComplete)
		return nil
	}

	// Start downloading blocks
	return csm.downloadBlocks(ctx, currentHeight+1, highestChain.Height)
}

// requestChainStatus requests chain status from a specific peer
func (csm *ChainSyncManager) requestChainStatus(ctx context.Context, peerID peer.ID) (network.ChainStatusMessage, error) {
	requestID := fmt.Sprintf("status_%d", time.Now().UnixNano())
	request := network.ChainStatusRequestMessage{
		RequestID: requestID,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return network.ChainStatusMessage{}, err
	}

	msg := network.NetworkMessage{
		Type:    network.MsgTypeChainStatusRequest,
		Payload: payload,
	}

	// Add pending request with timeout
	responseChan := csm.addPendingRequest(requestID, network.MsgTypeChainStatusRequest, peerID, 10*time.Second)

	// Send request
	err = csm.p2pNode.SendMessage(ctx, peerID, msg)
	if err != nil {
		csm.resolvePendingRequest(requestID, err)
		return network.ChainStatusMessage{}, err
	}

	// Wait for response
	select {
	case response := <-responseChan:
		if err, ok := response.(error); ok {
			return network.ChainStatusMessage{}, err
		}
		if status, ok := response.(network.ChainStatusMessage); ok {
			return status, nil
		}
		return network.ChainStatusMessage{}, fmt.Errorf("invalid response type")
	case <-ctx.Done():
		csm.resolvePendingRequest(requestID, ctx.Err())
		return network.ChainStatusMessage{}, ctx.Err()
	}
}

// findHighestChain finds the chain with the highest height among peers
func (csm *ChainSyncManager) findHighestChain(statuses []network.ChainStatusMessage) network.ChainStatusMessage {
	if len(statuses) == 0 {
		return network.ChainStatusMessage{}
	}

	highest := statuses[0]
	for _, status := range statuses {
		if status.Height > highest.Height {
			highest = status
		}
	}

	return highest
}

// downloadBlocks downloads blocks from the specified range
func (csm *ChainSyncManager) downloadBlocks(ctx context.Context, fromHeight, toHeight int64) error {
	log.Printf("📥 [SYNC] Downloading blocks from %d to %d", fromHeight, toHeight)

	csm.totalBlocks = toHeight - fromHeight + 1

	// Get connected peers
	peers := csm.p2pNode.Host.Peerstore().Peers()
	if len(peers) == 0 {
		return fmt.Errorf("no peers available")
	}

	// Try to download from each peer
	for _, peerID := range peers {
		if peerID == csm.p2pNode.Host.ID() {
			continue
		}

		err := csm.downloadBlocksFromPeer(ctx, peerID, fromHeight, toHeight)
		if err == nil {
			log.Printf("✅ [SYNC] Successfully downloaded blocks from peer %s", peerID.String())
			return nil
		}

		log.Printf("⚠️ [SYNC] Failed to download from peer %s: %v", peerID.String(), err)
	}

	return fmt.Errorf("failed to download blocks from any peer")
}

// downloadBlocksFromPeer downloads blocks from a specific peer
func (csm *ChainSyncManager) downloadBlocksFromPeer(ctx context.Context, peerID peer.ID, fromHeight, toHeight int64) error {
	// Split download into chunks to avoid overwhelming the peer
	chunkSize := int64(100) // Download 100 blocks at a time

	for currentHeight := fromHeight; currentHeight <= toHeight; currentHeight += chunkSize {
		endHeight := currentHeight + chunkSize - 1
		if endHeight > toHeight {
			endHeight = toHeight
		}

		err := csm.downloadBlockChunk(ctx, peerID, currentHeight, endHeight)
		if err != nil {
			return fmt.Errorf("failed to download chunk %d-%d: %v", currentHeight, endHeight, err)
		}

		csm.blocksSynced += endHeight - currentHeight + 1
		if csm.onSyncProgress != nil {
			csm.onSyncProgress(csm.blocksSynced, csm.totalBlocks)
		}

		log.Printf("📊 [SYNC] Progress: %d/%d blocks", csm.blocksSynced, csm.totalBlocks)
	}

	return nil
}

// downloadBlockChunk downloads a chunk of blocks
func (csm *ChainSyncManager) downloadBlockChunk(ctx context.Context, peerID peer.ID, fromHeight, toHeight int64) error {
	requestID := fmt.Sprintf("blocks_%d_%d_%d", fromHeight, toHeight, time.Now().UnixNano())
	request := network.BlockRequestMessage{
		FromIndex: fromHeight,
		ToIndex:   toHeight,
		RequestID: requestID,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	msg := network.NetworkMessage{
		Type:    network.MsgTypeBlockRequest,
		Payload: payload,
	}

	// Add pending request with timeout
	responseChan := csm.addPendingRequest(requestID, network.MsgTypeBlockRequest, peerID, 30*time.Second)

	// Send request
	err = csm.p2pNode.SendMessage(ctx, peerID, msg)
	if err != nil {
		csm.resolvePendingRequest(requestID, err)
		return err
	}

	log.Printf("📥 [SYNC] Requested blocks %d-%d from peer %s", fromHeight, toHeight, peerID.String())

	// Wait for response
	select {
	case response := <-responseChan:
		if err, ok := response.(error); ok {
			return err
		}
		if blockResponse, ok := response.(network.BlockResponseMessage); ok {
			if blockResponse.Error != "" {
				return fmt.Errorf("peer error: %s", blockResponse.Error)
			}
			return csm.processDownloadedBlocks(blockResponse.Blocks)
		}
		return fmt.Errorf("invalid response type")
	case <-ctx.Done():
		csm.resolvePendingRequest(requestID, ctx.Err())
		return ctx.Err()
	}
}

// processDownloadedBlocks processes and validates downloaded blocks
func (csm *ChainSyncManager) processDownloadedBlocks(blockData [][]byte) error {
	for i, data := range blockData {
		var blk block.Block
		if err := json.Unmarshal(data, &blk); err != nil {
			return fmt.Errorf("failed to unmarshal block %d: %v", i, err)
		}

		// Validate block
		if err := csm.validateDownloadedBlock(&blk); err != nil {
			return fmt.Errorf("block %d validation failed: %v", blk.Index, err)
		}

		// Add block to chain
		if err := csm.blockManager.AddBlock(&blk); err != nil {
			return fmt.Errorf("failed to add block %d: %v", blk.Index, err)
		}

		log.Printf("✅ [SYNC] Added block %d to chain", blk.Index)
	}

	return nil
}

// validateDownloadedBlock validates a downloaded block
func (csm *ChainSyncManager) validateDownloadedBlock(blk *block.Block) error {
	// Check if block index is sequential
	expectedIndex := csm.blockManager.GetBlockHeight() + 1
	if blk.Index != expectedIndex {
		return fmt.Errorf("invalid block index: expected %d, got %d", expectedIndex, blk.Index)
	}

	// Check previous hash
	lastBlock := csm.blockManager.GetLatestBlock()
	if blk.PrevHash != lastBlock.Hash {
		return fmt.Errorf("invalid previous hash: expected %s, got %s", lastBlock.Hash, blk.PrevHash)
	}

	// Verify block hash
	expectedHash := block.CalculateHash(*blk)
	if blk.Hash != expectedHash {
		return fmt.Errorf("invalid block hash: expected %s, got %s", expectedHash, blk.Hash)
	}

	return nil
}

// createMockBlock creates a mock block for testing purposes
func (csm *ChainSyncManager) createMockBlock(index int) *block.Block {
	prevBlock := csm.blockManager.GetLatestBlock()

	// Create a simple mock block
	mockBlock := &block.Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: []transaction.Transaction{}, // Empty transactions for now
		PrevHash:     prevBlock.Hash,
		Validator:    "mock_validator",
		Signature:    "mock_signature",
	}

	// Calculate hash
	mockBlock.Hash = block.CalculateHash(*mockBlock)

	return mockBlock
}

// handleChainStatusRequest handles incoming chain status requests
func (csm *ChainSyncManager) handleChainStatusRequest(msg network.NetworkMessage) {
	var request network.ChainStatusRequestMessage
	if err := json.Unmarshal(msg.Payload, &request); err != nil {
		log.Printf("❌ [SYNC] Failed to unmarshal chain status request: %v", err)
		return
	}

	// Create response
	latestBlock := csm.blockManager.GetLatestBlock()
	response := network.ChainStatusResponseMessage{
		Status: network.ChainStatusMessage{
			Height:      int64(csm.blockManager.GetBlockHeight()),
			LatestHash:  latestBlock.Hash,
			GenesisHash: csm.genesisHash,
			TotalBlocks: int64(csm.blockManager.GetChainLength()),
			IsSyncing:   csm.status == SyncStatusSyncing,
			PeerID:      csm.p2pNode.Host.ID().String(),
		},
		RequestID: request.RequestID,
	}

	// Send response back to the requesting peer
	payload, err := json.Marshal(response)
	if err != nil {
		log.Printf("❌ [SYNC] Failed to marshal chain status response: %v", err)
		return
	}

	responseMsg := network.NetworkMessage{
		Type:    network.MsgTypeChainStatusResponse,
		Payload: payload,
	}

	// Send response to the requesting peer
	if msg.FromPeer != "" {
		peerID, err := peer.Decode(msg.FromPeer)
		if err != nil {
			log.Printf("❌ [SYNC] Failed to decode peer ID: %v", err)
			return
		}

		ctx := context.Background()
		err = csm.p2pNode.SendMessage(ctx, peerID, responseMsg)
		if err != nil {
			log.Printf("❌ [SYNC] Failed to send chain status response: %v", err)
		} else {
			log.Printf("📤 [SYNC] Sent chain status response to peer %s: height=%d", msg.FromPeer, response.Status.Height)
		}
	}

	// Also resolve any pending request if it exists
	if csm.resolvePendingRequest(request.RequestID, response.Status) {
		log.Printf("✅ [SYNC] Resolved pending chain status request: %s", request.RequestID)
	}
}

// handleBlockRequest handles incoming block requests
func (csm *ChainSyncManager) handleBlockRequest(msg network.NetworkMessage) {
	var request network.BlockRequestMessage
	if err := json.Unmarshal(msg.Payload, &request); err != nil {
		log.Printf("❌ [SYNC] Failed to unmarshal block request: %v", err)
		return
	}

	log.Printf("📤 [SYNC] Received block request: %d-%d", request.FromIndex, request.ToIndex)

	// Get requested blocks
	blocks := make([][]byte, 0)
	for i := request.FromIndex; i <= request.ToIndex; i++ {
		blk, err := csm.blockManager.GetBlockByIndex(int(i))
		if err != nil {
			log.Printf("❌ [SYNC] Block %d not found: %v", i, err)
			continue
		}

		blockData, err := json.Marshal(blk)
		if err != nil {
			log.Printf("❌ [SYNC] Failed to marshal block %d: %v", i, err)
			continue
		}

		blocks = append(blocks, blockData)
	}

	// Create response
	response := network.BlockResponseMessage{
		Blocks:    blocks,
		RequestID: request.RequestID,
	}

	if len(blocks) == 0 {
		response.Error = "no blocks found in requested range"
	}

	// Send response back to the requesting peer
	payload, err := json.Marshal(response)
	if err != nil {
		log.Printf("❌ [SYNC] Failed to marshal block response: %v", err)
		return
	}

	responseMsg := network.NetworkMessage{
		Type:    network.MsgTypeBlockResponse,
		Payload: payload,
	}

	// Send response to the requesting peer
	if msg.FromPeer != "" {
		peerID, err := peer.Decode(msg.FromPeer)
		if err != nil {
			log.Printf("❌ [SYNC] Failed to decode peer ID: %v", err)
			return
		}

		ctx := context.Background()
		err = csm.p2pNode.SendMessage(ctx, peerID, responseMsg)
		if err != nil {
			log.Printf("❌ [SYNC] Failed to send block response: %v", err)
		} else {
			log.Printf("📤 [SYNC] Sent %d blocks to peer %s for request %d-%d", len(blocks), msg.FromPeer, request.FromIndex, request.ToIndex)
		}
	}

	// Also resolve any pending request if it exists
	if csm.resolvePendingRequest(request.RequestID, response) {
		log.Printf("✅ [SYNC] Resolved pending block request: %s", request.RequestID)
	}
}

// handleSyncRequest handles incoming sync requests
func (csm *ChainSyncManager) handleSyncRequest(msg network.NetworkMessage) {
	var request network.SyncRequestMessage
	if err := json.Unmarshal(msg.Payload, &request); err != nil {
		log.Printf("❌ [SYNC] Failed to unmarshal sync request: %v", err)
		return
	}

	log.Printf("📤 [SYNC] Received sync request: %d-%d (fast_sync=%v)", request.FromHeight, request.ToHeight, request.FastSync)

	// TODO: Send sync response with requested blocks
}

// handleChainStatusResponse handles incoming chain status responses
func (csm *ChainSyncManager) handleChainStatusResponse(msg network.NetworkMessage) {
	var response network.ChainStatusResponseMessage
	if err := json.Unmarshal(msg.Payload, &response); err != nil {
		log.Printf("❌ [SYNC] Failed to unmarshal chain status response: %v", err)
		return
	}

	log.Printf("📥 [SYNC] Received chain status: height=%d, hash=%s", response.Status.Height, response.Status.LatestHash)

	// Resolve the pending request
	if csm.resolvePendingRequest(response.RequestID, response.Status) {
		log.Printf("✅ [SYNC] Resolved chain status request: %s", response.RequestID)
	} else {
		log.Printf("⚠️ [SYNC] No pending request found for chain status: %s", response.RequestID)
	}
}

// handleBlockResponse handles incoming block responses
func (csm *ChainSyncManager) handleBlockResponse(msg network.NetworkMessage) {
	var response network.BlockResponseMessage
	if err := json.Unmarshal(msg.Payload, &response); err != nil {
		log.Printf("❌ [SYNC] Failed to unmarshal block response: %v", err)
		return
	}

	log.Printf("📥 [SYNC] Received %d blocks", len(response.Blocks))

	// Resolve the pending request
	if csm.resolvePendingRequest(response.RequestID, response) {
		log.Printf("✅ [SYNC] Resolved block request: %s", response.RequestID)
	} else {
		log.Printf("⚠️ [SYNC] No pending request found for block response: %s", response.RequestID)
	}
}

// handleSyncResponse handles incoming sync responses
func (csm *ChainSyncManager) handleSyncResponse(msg network.NetworkMessage) {
	var response network.SyncResponseMessage
	if err := json.Unmarshal(msg.Payload, &response); err != nil {
		log.Printf("❌ [SYNC] Failed to unmarshal sync response: %v", err)
		return
	}

	log.Printf("📥 [SYNC] Received sync response: %d blocks, complete=%v", len(response.Blocks), response.IsComplete)

	// TODO: Process received blocks and state
}

// handleForkResolution handles incoming fork resolution messages
func (csm *ChainSyncManager) handleForkResolution(msg network.NetworkMessage) {
	var forkMsg network.ForkResolutionMessage
	if err := json.Unmarshal(msg.Payload, &forkMsg); err != nil {
		log.Printf("❌ [SYNC] Failed to unmarshal fork resolution: %v", err)
		return
	}

	log.Printf("🔄 [SYNC] Fork detected at height %d", forkMsg.ForkHeight)

	if csm.onForkDetected != nil {
		csm.onForkDetected(forkMsg.ForkHeight, forkMsg.Canonical, forkMsg.Forked)
	}
}

// setStatus updates the sync status
func (csm *ChainSyncManager) setStatus(status SyncStatus) {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	csm.status = status
	if status == SyncStatusComplete || status == SyncStatusFailed {
		csm.syncEndTime = time.Now()
	}

	log.Printf("📊 [SYNC] Status changed to: %s", status)

	switch status {
	case SyncStatusComplete:
		if csm.onSyncComplete != nil {
			csm.onSyncComplete()
		}
	case SyncStatusFailed:
		if csm.onSyncFailed != nil {
			csm.onSyncFailed(fmt.Errorf("sync failed"))
		}
	}
}

// GetStatus returns the current sync status
func (csm *ChainSyncManager) GetStatus() SyncStatus {
	csm.mu.RLock()
	defer csm.mu.RUnlock()
	return csm.status
}

// GetSyncProgress returns the current sync progress
func (csm *ChainSyncManager) GetSyncProgress() (blocksSynced, totalBlocks int64) {
	csm.mu.RLock()
	defer csm.mu.RUnlock()
	return csm.blocksSynced, csm.totalBlocks
}

// GetSyncDuration returns the duration of the current sync
func (csm *ChainSyncManager) GetSyncDuration() time.Duration {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	if csm.status == SyncStatusIdle {
		return 0
	}

	endTime := csm.syncEndTime
	if csm.status == SyncStatusSyncing {
		endTime = time.Now()
	}

	return endTime.Sub(csm.syncStartTime)
}

// SetCallbacks sets the callback functions for sync events
func (csm *ChainSyncManager) SetCallbacks(
	onSyncStart func(),
	onSyncProgress func(blocksSynced, totalBlocks int64),
	onSyncComplete func(),
	onSyncFailed func(error),
	onForkDetected func(forkHeight int64, canonical, forked []string),
) {
	csm.onSyncStart = onSyncStart
	csm.onSyncProgress = onSyncProgress
	csm.onSyncComplete = onSyncComplete
	csm.onSyncFailed = onSyncFailed
	csm.onForkDetected = onForkDetected
}

// addPendingRequest adds a pending request to track
func (csm *ChainSyncManager) addPendingRequest(requestID string, requestType network.MessageType, peerID peer.ID, timeout time.Duration) chan interface{} {
	csm.requestMutex.Lock()
	defer csm.requestMutex.Unlock()

	responseChan := make(chan interface{}, 1)
	pendingReq := &PendingRequest{
		RequestID:   requestID,
		RequestType: requestType,
		PeerID:      peerID,
		Timestamp:   time.Now(),
		Response:    responseChan,
		Timeout:     timeout,
	}

	csm.pendingRequests[requestID] = pendingReq

	// Start timeout goroutine
	go func() {
		time.Sleep(timeout)
		csm.requestMutex.Lock()
		if req, exists := csm.pendingRequests[requestID]; exists {
			select {
			case req.Response <- fmt.Errorf("request timeout"):
			default:
			}
			delete(csm.pendingRequests, requestID)
		}
		csm.requestMutex.Unlock()
	}()

	return responseChan
}

// resolvePendingRequest resolves a pending request with a response
func (csm *ChainSyncManager) resolvePendingRequest(requestID string, response interface{}) bool {
	csm.requestMutex.Lock()
	defer csm.requestMutex.Unlock()

	if req, exists := csm.pendingRequests[requestID]; exists {
		select {
		case req.Response <- response:
		default:
		}
		delete(csm.pendingRequests, requestID)
		return true
	}
	return false
}

// cleanupExpiredRequests removes expired requests
func (csm *ChainSyncManager) cleanupExpiredRequests() {
	csm.requestMutex.Lock()
	defer csm.requestMutex.Unlock()

	now := time.Now()
	for requestID, req := range csm.pendingRequests {
		if now.Sub(req.Timestamp) > req.Timeout {
			select {
			case req.Response <- fmt.Errorf("request expired"):
			default:
			}
			delete(csm.pendingRequests, requestID)
		}
	}
}
