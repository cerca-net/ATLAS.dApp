package indexer

import (
	"atlas-blockchain/pkg/block"
	"atlas-blockchain/pkg/transaction"
	"atlas-blockchain/pkg/wallet"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// SupabaseIndexer bridges the on-chain ATLAS state with the off-chain Supabase database.
type SupabaseIndexer struct {
	URL        string
	ServiceKey string
	HttpClient *http.Client
}

// NewSupabaseIndexer initializes the indexer. If environment variables are missing, it returns nil.
func NewSupabaseIndexer() *SupabaseIndexer {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if url == "" || key == "" {
		log.Printf("[INDEXER] ⚠️ Supabase credentials not found. On-chain indexing disabled.")
		return nil
	}

	return &SupabaseIndexer{
		URL:        url,
		ServiceKey: key,
		HttpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// ProcessBlock intercepts confirmed blocks and pushes relevant events to Supabase
func (i *SupabaseIndexer) ProcessBlock(b *block.Block) {
	if i == nil {
		return
	}

	log.Printf("[INDEXER] 🔄 Processing block %d for Supabase consistency sync...", b.Index)

	for _, tx := range b.Transactions {
		// Index transfers and smart contract interactions to keep the Flutter UI synced
		i.syncTransaction(tx, int64(b.Index))
	}
}

func (i *SupabaseIndexer) syncTransaction(tx transaction.Transaction, blockIndex int64) {
	endpoint := fmt.Sprintf("%s/rest/v1/onchain_transactions", i.URL)

	txHashBytes := wallet.CalculateTxHash(tx)
	txID := hex.EncodeToString(txHashBytes)

	payload := map[string]interface{}{
		"tx_id":        txID,
		"sender":       tx.Sender,
		"recipient":    tx.Recipient,
		"amount":       tx.Amount,
		"tx_type":      tx.Type,
		"block_height": blockIndex,
		"timestamp":    time.Now().Unix(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[INDEXER] Failed to marshal payload: %v", err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("[INDEXER] Failed to create request: %v", err)
		return
	}

	req.Header.Set("apikey", i.ServiceKey)
	req.Header.Set("Authorization", "Bearer "+i.ServiceKey)
	req.Header.Set("Content-Type", "application/json")
	// If the tx already exists, ignore or merge.
	req.Header.Set("Prefer", "resolution=merge-duplicates")

	resp, err := i.HttpClient.Do(req)
	if err != nil {
		log.Printf("[INDEXER] HTTP error syncing tx %s: %v", txID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		log.Printf("[INDEXER] ❌ Supabase API returned status %d for tx %s", resp.StatusCode, txID)
	} else {
		log.Printf("[INDEXER] ✅ Successfully synced tx %s to Supabase", txID)
	}
}
