package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"atlas-blockchain/internal/blockchain"
	"atlas-blockchain/pkg/config"
	"atlas-blockchain/pkg/network"
)

func TestHandleGetStatus(t *testing.T) {
	// 1. Setup

	// Create temp directory for test data
	tmpDir, err := ioutil.TempDir("", "atlas_api_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize Config
	cfg := config.DefaultConfig()
	cfg.DataDir = tmpDir

	// Initialize Managers
	// Note: Dependencies must match main.go wiring
	sm := blockchain.NewStateManager(cfg)
	bm := blockchain.NewBlockManager(cfg, sm)
	tm := blockchain.NewTransactionManager(cfg, sm)
	cm := blockchain.NewConsensusManager(cfg, bm)

	// Wire up circular dependencies if methods exist (sanity check)
	// sm.SetConsensusManager(cm)

	// Initialize Node
	node := network.NewNode("localhost:8000", "test-token")
	node.ValidatorAddress = "mock_validator_address"

	// Initialize API Server
	// func NewAPIServer(bm, tm, sm, cm, fsm, pm, node, p2pNode, im, socialMgr, govMgr)
	apiServer := NewAPIServer(bm, tm, sm, cm, nil, nil, node, nil, nil, nil, nil, nil)

	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()

	// 2. Execution
	// We use the handler directly. NOTE: handleGetStatus is a method on *APIServer
	handler := http.HandlerFunc(apiServer.handleGetStatus)
	handler.ServeHTTP(rr, req)

	// 3. Verification
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	expectedKeys := []string{
		"blockHeight",
		"txPoolSize",
		"isValidator",
		"totalValidators",
		"mode",
	}

	for _, key := range expectedKeys {
		if _, ok := response[key]; !ok {
			t.Errorf("response missing expected key: %s", key)
		}
	}

	// Check values (fresh state)
	if val, ok := response["blockHeight"].(float64); !ok || val != 0 {
		t.Errorf("blockHeight has wrong value: got %v, want 0", response["blockHeight"])
	}
	if val, ok := response["txPoolSize"].(float64); !ok || val != 0 {
		t.Errorf("txPoolSize has wrong value: got %v, want 0", response["txPoolSize"])
	}
	if val, ok := response["mode"].(string); !ok || val != "observer" {
		// Since "mock_validator_address" is not registered in consensus manager, it should be observer
		t.Errorf("mode has wrong value: got %v, want observer", response["mode"])
	}
}
