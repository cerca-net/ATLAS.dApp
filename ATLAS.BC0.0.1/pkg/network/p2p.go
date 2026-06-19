package network

import (
	"atlas-blockchain/pkg/block"
	"atlas-blockchain/pkg/transaction"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const ProtocolID = "/blockchain/1.0.0"

// P2PNode represents a minimal libp2p node
type P2PNode struct {
	Host                  host.Host
	HandleIncomingMessage func(msg NetworkMessage) // Function field for message handling
	// Callbacks for integration
	OnBlockReceived                 func(block BlockMessage)
	OnTransactionReceived           func(tx TransactionMessage)
	OnValidatorRegistrationReceived func(reg ValidatorRegistrationMessage) // New callback

	// GossipSub state for efficient propagation
	PubSub     *pubsub.PubSub
	BlockTopic *pubsub.Topic
	TxTopic    *pubsub.Topic
}

// loadOrCreatePrivKey loads a private key from file or generates and saves a new one
func loadOrCreatePrivKey(path string) (crypto.PrivKey, error) {
	if _, err := os.Stat(path); err == nil {
		// File exists, load it
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		privKey, err := crypto.UnmarshalPrivateKey(data)
		if err != nil {
			return nil, err
		}
		return privKey, nil
	}
	// File does not exist, generate new key
	privKey, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}
	data, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		return nil, err
	}
	return privKey, nil
}

// NewP2PNode creates and starts a new libp2p node listening on the given port, using a persistent private key
func NewP2PNode(ctx context.Context, listenPort int, keyPath ...string) (*P2PNode, error) {
	log.Printf("[DEBUG] NewP2PNode called with listenPort=%d, keyPath=%v", listenPort, keyPath)
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/libp2p/go-libp2p" {
				log.Printf("[DEBUG] Using libp2p version: %s", dep.Version)
			}
		}
	}
	var privKey crypto.PrivKey
	var err error
	if len(keyPath) > 0 && keyPath[0] != "" {
		privKey, err = loadOrCreatePrivKey(keyPath[0])
		if err != nil {
			return nil, err
		}
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
	}
	if privKey != nil {
		opts = append(opts, libp2p.Identity(privKey))
	}
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}
	fmt.Println("[libp2p] Node ID:", h.ID())
	for _, addr := range h.Addrs() {
		fmt.Println("[libp2p] Listening on:", addr)
	}

	node := &P2PNode{Host: h}
	// Default handler does nothing
	node.HandleIncomingMessage = func(msg NetworkMessage) {}

	// Initialize GossipSub
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		log.Printf("[P2P] Failed to initialize GossipSub: %v", err)
	} else {
		node.PubSub = ps
		node.BlockTopic, _ = ps.Join("blocks")
		node.TxTopic, _ = ps.Join("transactions")
	}

	// Enable MDNS for local peer discovery
	// This allows nodes on the same local network to discover each other automatically
	mdnsService := mdns.NewMdnsService(h, "cercachain-mdns", &discoveryNotifee{h: h})
	if err := mdnsService.Start(); err != nil {
		log.Printf("[P2P] Failed to start MDNS: %v", err)
	} else {
		log.Printf("[P2P] MDNS peer discovery enabled for local network")
	}

	// Start background listeners for GossipSub
	_ = node.StartGossipSubscriptions(ctx)

	return node, nil
}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if pi.ID == n.h.ID() {
		return
	}
	log.Printf("[P2P] MDNS discovered local peer: %s", pi.ID.String())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := n.h.Connect(ctx, pi)
	if err != nil {
		log.Printf("[P2P] Failed to connect to discovered peer %s: %v", pi.ID.String(), err)
		// Remove stale/phantom peer from peerstore to prevent repeated dial attempts
		n.h.Peerstore().RemovePeer(pi.ID)
		n.h.Peerstore().ClearAddrs(pi.ID)
		log.Printf("[P2P] Evicted stale peer %s from peerstore", pi.ID.String())
	} else {
		log.Printf("[P2P] Successfully connected to discovered peer: %s", pi.ID.String())
	}
}

// StartGossipSubscriptions begins listening for blocks and transactions on the public pubsub topics.
func (node *P2PNode) StartGossipSubscriptions(ctx context.Context) error {
	if node.PubSub == nil {
		return fmt.Errorf("GossipSub not initialized")
	}

	if node.BlockTopic != nil {
		blockSub, err := node.BlockTopic.Subscribe()
		if err == nil {
			go func() {
				for {
					msg, err := blockSub.Next(ctx)
					if err != nil {
						return
					}
					if msg.ReceivedFrom == node.Host.ID() {
						continue
					}
					var blockMsg BlockMessage
					if err := json.Unmarshal(msg.Data, &blockMsg); err == nil && node.OnBlockReceived != nil {
						node.OnBlockReceived(blockMsg)
					}
				}
			}()
		}
	}

	if node.TxTopic != nil {
		txSub, err := node.TxTopic.Subscribe()
		if err == nil {
			go func() {
				for {
					msg, err := txSub.Next(ctx)
					if err != nil {
						return
					}
					if msg.ReceivedFrom == node.Host.ID() {
						continue
					}
					var txMsg TransactionMessage
					if err := json.Unmarshal(msg.Data, &txMsg); err == nil && node.OnTransactionReceived != nil {
						node.OnTransactionReceived(txMsg)
					}
				}
			}()
		}
	}

	return nil
}

// RegisterStreamHandler sets up the handler for incoming libp2p streams
func (node *P2PNode) RegisterStreamHandler() {
	node.Host.SetStreamHandler(ProtocolID, func(s network.Stream) {
		remotePeer := s.Conn().RemotePeer()
		log.Printf("[P2P] Stream handler triggered from peer: %s", remotePeer.String())
		defer s.Close()
		var msg NetworkMessage
		decoder := json.NewDecoder(s)
		if err := decoder.Decode(&msg); err != nil {
			log.Printf("[P2P] Error decoding message: %v", err)
			return
		}

		// Add peer information to the message
		msg.FromPeer = remotePeer.String()
		log.Printf("[P2P] Received message of type: %s from peer: %s", msg.Type, remotePeer.String())

		// Integration: handle block and transaction messages
		switch msg.Type {
		case MsgTypeBlock:
			var blockMsg BlockMessage
			if err := json.Unmarshal(msg.Payload, &blockMsg); err == nil && node.OnBlockReceived != nil {
				node.OnBlockReceived(blockMsg)
			}
		case MsgTypeTransaction:
			var txMsg TransactionMessage
			if err := json.Unmarshal(msg.Payload, &txMsg); err == nil && node.OnTransactionReceived != nil {
				node.OnTransactionReceived(txMsg)
			}
		case MsgTypeValidatorRegistration:
			var regMsg ValidatorRegistrationMessage
			if err := json.Unmarshal(msg.Payload, &regMsg); err == nil && node.OnValidatorRegistrationReceived != nil {
				node.OnValidatorRegistrationReceived(regMsg)
			}
		default:
			node.HandleIncomingMessage(msg)
		}
	})
}

// SendMessage sends a NetworkMessage to the given peer over a new libp2p stream
func (node *P2PNode) SendMessage(ctx context.Context, peerID peer.ID, msg NetworkMessage) error {
	log.Printf("[P2P] Attempting to send message of type: %s to peer: %s", msg.Type, peerID.String())
	s, err := node.Host.NewStream(ctx, peerID, ProtocolID)
	if err != nil {
		log.Printf("[P2P] Error opening stream to peer: %v", err)
		return err
	}
	defer s.Close()
	encoder := json.NewEncoder(s)
	err = encoder.Encode(msg)
	if err != nil {
		log.Printf("[P2P] Error sending message: %v", err)
	} else {
		log.Printf("[P2P] Message sent successfully to peer: %s", peerID.String())
	}
	return err
}

// BroadcastBlock sends a block to all connected peers
func (node *P2PNode) BroadcastBlock(ctx context.Context, b *block.Block) error {
	peers := node.Host.Peerstore().Peers()
	log.Printf("[DEBUG] BroadcastBlock: Broadcasting block %d to %d known peers", b.Index, len(peers))
	sentCount := 0

	blockMsg := BlockMessage{Block: *b}
	payload, err := json.Marshal(blockMsg)
	if err != nil {
		log.Printf("[P2P] Failed to marshal BlockMessage: %v", err)
		return err
	}

	if node.BlockTopic != nil {
		err = node.BlockTopic.Publish(ctx, payload)
		if err != nil {
			log.Printf("[P2P] Failed to publish block to GossipSub: %v", err)
		} else {
			log.Printf("[P2P] Block %d published to GossipSub", b.Index)
		}
	}

	for _, peerID := range peers {
		msg := NetworkMessage{
			Type:    MsgTypeBlock,
			Payload: payload,
		}
		err = node.SendMessage(ctx, peerID, msg)
		if err != nil {
			log.Printf("[P2P] Failed to send block to peer %s: %v", peerID.String(), err)
		} else {
			log.Printf("[P2P] Block %d sent to peer %s", b.Index, peerID.String())
			sentCount++
		}
	}
	log.Printf("[P2P] Broadcast block %d to %d peers", b.Index, sentCount)
	return nil
}

// BroadcastTransaction sends a transaction to all connected peers
func (node *P2PNode) BroadcastTransaction(tx transaction.Transaction) error {
	payload, err := json.Marshal(TransactionMessage{Transaction: tx})
	if err != nil {
		return fmt.Errorf("failed to marshal TransactionMessage: %v", err)
	}

	if node.TxTopic != nil {
		err = node.TxTopic.Publish(context.Background(), payload)
		if err != nil {
			log.Printf("[P2P] Failed to publish transaction to GossipSub: %v", err)
		} else {
			log.Printf("[P2P] Transaction published to GossipSub")
		}
	}

	msg := NetworkMessage{
		Type:    MsgTypeTransaction,
		Payload: payload,
	}

	return node.BroadcastMessage(context.Background(), msg)
}

// BroadcastMessage sends a NetworkMessage to all connected peers
func (node *P2PNode) BroadcastMessage(ctx context.Context, msg NetworkMessage) error {
	peers := node.Host.Peerstore().Peers()
	successCount := 0

	for _, peerID := range peers {
		if peerID == node.Host.ID() {
			continue // Skip self
		}

		// Only send to peers we are actually connected to
		if node.Host.Network().Connectedness(peerID) != network.Connected {
			continue
		}

		err := node.SendMessage(ctx, peerID, msg)
		if err != nil {
			log.Printf("[P2P] Failed to broadcast %s to peer %s: %v", msg.Type, peerID, err)
			// Evict unreachable peer from peerstore
			node.Host.Peerstore().RemovePeer(peerID)
			continue
		}
		successCount++
	}

	log.Printf("[P2P] Broadcasted %s to %d peers", msg.Type, successCount)
	return nil
}

// ConnectToBootstrapPeers connects to a list of bootstrap peer multiaddresses.
// It returns the slice of Peer IDs that were successfully connected to.
func (node *P2PNode) ConnectToBootstrapPeers(ctx context.Context, addresses []string) []peer.ID {
	var connected []peer.ID
	for _, addr := range addresses {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}
		log.Printf("[P2P] Parsing bootstrap address: %s", addr)
		info, err := peer.AddrInfoFromString(addr)
		if err != nil {
			log.Printf("[P2P] Invalid bootstrap multiaddress '%s': %v", addr, err)
			continue
		}

		// Don't connect to ourselves
		if info.ID == node.Host.ID() {
			log.Printf("[P2P] Skipping bootstrap address: self-connection (%s)", info.ID)
			continue
		}

		// Attempt manual connection with retries
		connectedToPeer := false
		maxAttempts := 3
		backoff := 5 * time.Second

		for attempt := 1; attempt <= maxAttempts; attempt++ {
			log.Printf("[P2P] Attempt %d/%d to connect to bootstrap peer: %s", attempt, maxAttempts, info.ID)
			dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err = node.Host.Connect(dialCtx, *info)
			cancel()

			if err == nil {
				log.Printf("[P2P] Successfully connected to bootstrap peer: %s (attempt %d)", info.ID, attempt)
				connected = append(connected, info.ID)
				connectedToPeer = true
				break
			}

			log.Printf("[P2P] Failed to connect to bootstrap peer %s (attempt %d): %v", info.ID, attempt, err)
			if attempt < maxAttempts {
				select {
				case <-ctx.Done():
					log.Printf("[P2P] Connection aborted due to context cancellation")
					return connected
				case <-time.After(backoff):
				}
			}
		}

		if !connectedToPeer {
			log.Printf("[P2P] Exhausted retry attempts to connect to bootstrap peer: %s", info.ID)
		}
	}
	return connected
}

// GetMultiaddress returns a representable multiaddress for this node,
// prioritizing a non-loopback IPv4 address if available.
func (node *P2PNode) GetMultiaddress() string {
	peerID := node.Host.ID().String()
	
	// Try to find a non-loopback address in Host.Addrs() first
	var fallback string
	for _, addr := range node.Host.Addrs() {
		addrStr := addr.String()
		// If it's not a loopback address or wildcard, we can use it
		if !strings.Contains(addrStr, "/127.0.0.1") && !strings.Contains(addrStr, "/0.0.0.0") && !strings.Contains(addrStr, "/::1") {
			return fmt.Sprintf("%s/p2p/%s", addrStr, peerID)
		}
		if strings.Contains(addrStr, "/127.0.0.1") && fallback == "" {
			fallback = fmt.Sprintf("%s/p2p/%s", addrStr, peerID)
		}
	}

	// If we only have wildcard (0.0.0.0) addresses, let's resolve our local IP
	localIP := getLocalIP()
	if localIP != "" {
		for _, addr := range node.Host.Addrs() {
			addrStr := addr.String()
			// Replace 0.0.0.0 with our actual local IP
			if strings.Contains(addrStr, "/0.0.0.0") {
				replaced := strings.Replace(addrStr, "/0.0.0.0", "/"+localIP, 1)
				return fmt.Sprintf("%s/p2p/%s", replaced, peerID)
			}
		}
	}

	if fallback != "" {
		return fallback
	}

	// Ultimate fallback: if no addresses exist, just use loopback with port 8000
	return fmt.Sprintf("/ip4/127.0.0.1/tcp/8000/p2p/%s", peerID)
}

// getLocalIP finds a non-loopback IPv4 address on the host interfaces.
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// WriteMultiaddressToFile writes the node's representable multiaddress to the specified file path.
func (node *P2PNode) WriteMultiaddressToFile(filePath string) error {
	multiaddr := node.GetMultiaddress()
	
	// Ensure parent directory exists
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for multiaddress file: %v", err)
		}
	}

	err := ioutil.WriteFile(filePath, []byte(multiaddr), 0644)
	if err != nil {
		return fmt.Errorf("failed to write multiaddress to %s: %v", filePath, err)
	}
	log.Printf("[P2P] Wrote multiaddress to %s: %s", filePath, multiaddr)
	return nil
}
