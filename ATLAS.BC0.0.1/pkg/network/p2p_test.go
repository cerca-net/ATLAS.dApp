package network

import (
	"context"
	"fmt"
	"testing"
	"time"
	"github.com/libp2p/go-libp2p/core/peer"
	"encoding/json"
)

func TestP2PBasicMessageExchange(t *testing.T) {
	t.Log("[Test] Starting P2PBasicMessageExchange test...")
	time.Sleep(10 * time.Second)
	ctx := context.Background()

	// Start Node A on port 9001
	nodeA, err := NewP2PNode(ctx, 9001)
	if err != nil {
		t.Fatalf("Failed to start Node A: %v", err)
	}
	nodeA.RegisterStreamHandler()

	// Start Node B on port 9002
	nodeB, err := NewP2PNode(ctx, 9002)
	if err != nil {
		t.Fatalf("Failed to start Node B: %v", err)
	}
	nodeB.RegisterStreamHandler()

	// Node A's multiaddress
	var nodeAAddr string
	for _, addr := range nodeA.Host.Addrs() {
		nodeAAddr = fmt.Sprintf("%s/p2p/%s", addr, nodeA.Host.ID().String())
		break
	}

	// Node B connects to Node A
	info, err := peer.AddrInfoFromString(nodeAAddr)
	if err != nil {
		t.Fatalf("Failed to parse Node A address: %v", err)
	}
	if err := nodeB.Host.Connect(ctx, *info); err != nil {
		t.Fatalf("Node B failed to connect to Node A: %v", err)
	}

	// Set up a channel to confirm receipt
	received := make(chan struct{})
	nodeA.HandleIncomingMessage = func(msg NetworkMessage) {
		if msg.Type == MsgTypePeerInfo {
			fmt.Println("[Test] Node A received PeerInfoMessage!")
			close(received)
		}
	}

	// Node B sends a PeerInfoMessage to Node A
	peerMsg := PeerInfoMessage{
		PeerID:  nodeB.Host.ID().String(),
		Address: nodeB.Host.Addrs()[0].String(),
	}
	payload, _ := EncodePayload(peerMsg)
	msg := NetworkMessage{
		Type:    MsgTypePeerInfo,
		Payload: payload,
	}
	if err := nodeB.SendMessage(ctx, nodeA.Host.ID(), msg); err != nil {
		t.Fatalf("Node B failed to send message: %v", err)
	}

	// Wait for confirmation or timeout
	select {
	case <-received:
		fmt.Println("[Test] Message exchange successful.")
	case <-time.After(3 * time.Second):
		t.Fatal("Timed out waiting for message receipt")
	}
}

// EncodePayload marshals a struct to json.RawMessage
func EncodePayload(v interface{}) ([]byte, error) {
	return json.Marshal(v)
} 