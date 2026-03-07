package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/block"
)

// PeerAnnounce represents a peer announcement message
type PeerAnnounce struct {
	Address string `json:"address"`
}

// startNetwork starts the network services for the node
func (n *Node) startNetwork() error {
	// Start listening for incoming connections
	listener, err := net.Listen("tcp", n.Address)
	if err != nil {
		return fmt.Errorf("failed to start listener: %v", err)
	}

	// Start peer discovery
	go n.startPeerDiscovery()

	// Start accepting connections
	go n.acceptConnections(listener)

	return nil
}

func (n *Node) startPeerDiscovery() {
	// Create UDP connection for peer discovery
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 8000, // Default peer discovery port
	})
	if err != nil {
		log.Printf("Failed to start peer discovery: %v", err)
		return
	}
	defer conn.Close()

	// Start broadcasting our presence
	go n.broadcastPresence(conn)

	// Start listening for peer announcements
	n.listenForPeers(conn)
}

func (n *Node) broadcastPresence(conn *net.UDPConn) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	announce := PeerAnnounce{Address: n.Address}
	msg := NetworkMessage{
		Type:    "peer_announce",
		Payload: mustMarshal(announce),
	}

	for {
		select {
		case <-ticker.C:
					// Broadcast to the network
		addr := &net.UDPAddr{
			IP:   net.IPv4(255, 255, 255, 255),
			Port: 8000, // Default peer discovery port
		}
			if _, err := conn.WriteTo(mustMarshal(msg), addr); err != nil {
				log.Printf("Failed to broadcast presence: %v", err)
			}
		}
	}
}

func (n *Node) listenForPeers(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		numBytes, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Failed to read from UDP: %v", err)
			continue
		}

		var msg NetworkMessage
		if err := json.Unmarshal(buffer[:numBytes], &msg); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		if msg.Type == "peer_announce" {
			var announce PeerAnnounce
			if err := json.Unmarshal(msg.Payload, &announce); err != nil {
				log.Printf("Failed to unmarshal peer announcement: %v", err)
				continue
			}

			// Add new peer if we haven't reached the maximum
			n.mu.Lock()
			if len(n.Peers) < 10 { // Default max peers
				// Check if we already know this peer
				known := false
				for _, peer := range n.Peers {
					if peer.Address == announce.Address {
						known = true
						break
					}
				}
				if !known {
					newPeer := NewNode(announce.Address, "testtoken")
					n.Peers = append(n.Peers, newPeer)
					log.Printf("Added new peer: %s", announce.Address)
				}
			}
			n.mu.Unlock()
		}
	}
}

func (n *Node) acceptConnections(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go n.handleConnection(conn)
	}
}

func (n *Node) handleNewBlock(block *block.Block) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Validate the block
	if ValidateChain(append(n.LocalBlockchain, block)) {
		n.LocalBlockchain = append(n.LocalBlockchain, block)
		// Update transaction pool
		n.LocalTxPool = removeProcessedTransactionsFromPool(n.LocalTxPool, block.Transactions)
		log.Printf("Added new block %d to chain", block.Index)
	} else {
		log.Printf("Rejected invalid block %d", block.Index)
	}
}

func (n *Node) handleNewTransaction(tx transaction.Transaction) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Validate transaction
	if err := tx.Validate(); err != nil {
		log.Printf("Rejected invalid transaction: %v", err)
		return
	}

	// Add to local pool if not at capacity
	if len(n.LocalTxPool) < 5000 { // Default max tx pool size
		n.LocalTxPool = append(n.LocalTxPool, tx)
		log.Printf("Added new transaction to pool")
	}
}

// Helper function to marshal JSON
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal: %v", err))
	}
	return data
} 