package network

import (
	"atlas-blockchain/pkg/block"
	"atlas-blockchain/pkg/transaction"
	"encoding/json"
)

// MessageType defines the type of message sent over the network
type MessageType string

const (
	MsgTypeBlock                 MessageType = "block"
	MsgTypeTransaction           MessageType = "transaction"
	MsgTypePeerInfo              MessageType = "peer_info"
	MsgTypeChainSync             MessageType = "chain_sync"
	MsgTypeValidatorRegistration MessageType = "validator_registration"
	MsgTypeBlockRequest          MessageType = "block_request"
	MsgTypeBlockResponse         MessageType = "block_response"
	MsgTypeChainStatus           MessageType = "chain_status"
	MsgTypeChainStatusRequest    MessageType = "chain_status_request"
	MsgTypeChainStatusResponse   MessageType = "chain_status_response"
	MsgTypeSyncRequest           MessageType = "sync_request"
	MsgTypeSyncResponse          MessageType = "sync_response"
	MsgTypeForkResolution        MessageType = "fork_resolution"
)

// NetworkMessage is the generic message wrapper
type NetworkMessage struct {
	Type     MessageType     `json:"type"`
	Payload  json.RawMessage `json:"payload"`
	FromPeer string          `json:"from_peer,omitempty"` // Peer ID of sender
}

// BlockMessage, TransactionMessage, PeerInfoMessage, ChainSyncMessage
// can be defined as needed, e.g.:
type BlockMessage struct {
	Block block.Block `json:"block"`
}

type TransactionMessage struct {
	Transaction transaction.Transaction `json:"transaction"`
}

type PeerInfoMessage struct {
	PeerID  string
	Address string
}

type ChainSyncMessage struct {
	Requesting bool
	// Add more fields as needed
}

// ValidatorRegistrationMessage is used to broadcast validator registration
type ValidatorRegistrationMessage struct {
	Address string `json:"address"`
	Stake   uint64 `json:"stake"`
	PeerID  string `json:"peer_id,omitempty"` // libp2p Peer ID
}

// Chain synchronization messages
type BlockRequestMessage struct {
	FromIndex int64  `json:"from_index"`
	ToIndex   int64  `json:"to_index"`
	RequestID string `json:"request_id"`
}

type BlockResponseMessage struct {
	Blocks    [][]byte `json:"blocks"`
	RequestID string   `json:"request_id"`
	Error     string   `json:"error,omitempty"`
}

type ChainStatusMessage struct {
	Height      int64  `json:"height"`
	LatestHash  string `json:"latest_hash"`
	GenesisHash string `json:"genesis_hash"`
	TotalBlocks int64  `json:"total_blocks"`
	IsSyncing   bool   `json:"is_syncing"`
	PeerID      string `json:"peer_id"`
}

type ChainStatusRequestMessage struct {
	RequestID string `json:"request_id"`
}

type ChainStatusResponseMessage struct {
	Status    ChainStatusMessage `json:"status"`
	RequestID string             `json:"request_id"`
	Error     string             `json:"error,omitempty"`
}

type SyncRequestMessage struct {
	FromHeight int64  `json:"from_height"`
	ToHeight   int64  `json:"to_height"`
	RequestID  string `json:"request_id"`
	FastSync   bool   `json:"fast_sync"` // Whether to use fast sync (state snapshots)
}

type SyncResponseMessage struct {
	Blocks     [][]byte `json:"blocks"`
	StateHash  string   `json:"state_hash,omitempty"`
	RequestID  string   `json:"request_id"`
	Error      string   `json:"error,omitempty"`
	IsComplete bool     `json:"is_complete"`
}

type ForkResolutionMessage struct {
	ForkHeight int64    `json:"fork_height"`
	Canonical  []string `json:"canonical_chain"` // Hashes of canonical chain
	Forked     []string `json:"forked_chain"`    // Hashes of forked chain
	Reason     string   `json:"reason"`
}
