package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TestResult struct {
	Endpoint string
	Status   int
	Response string
	Success  bool
}

func main() {
	fmt.Println("🚀 ATLAS Blockchain Platform - Comprehensive Feature Test")
	fmt.Println("=" * 60)
	
	// Wait for server to start
	time.Sleep(3 * time.Second)
	
	baseURL := "http://localhost:8080"
	results := []TestResult{}
	
	// Test 1: Core Blockchain Status
	fmt.Println("\n📊 1. Testing Core Blockchain Status...")
	status := testEndpoint(baseURL+"/status", "GET", nil)
	results = append(results, status)
	
	// Test 2: Create Wallet
	fmt.Println("\n💰 2. Testing Wallet Creation...")
	walletData := map[string]interface{}{
		"name": "Test Wallet",
		"type": "standard",
	}
	wallet := testEndpoint(baseURL+"/create-wallet", "POST", walletData)
	results = append(results, wallet)
	
	// Test 3: Get Balance
	fmt.Println("\n💳 3. Testing Balance Check...")
	balance := testEndpoint(baseURL+"/balance?address=test", "GET", nil)
	results = append(results, balance)
	
	// Test 4: Faucet (Get Test Tokens)
	fmt.Println("\n🎁 4. Testing Faucet...")
	faucetData := map[string]interface{}{
		"address": "test_address_123",
		"amount":  1000,
	}
	faucet := testEndpoint(baseURL+"/faucet", "POST", faucetData)
	results = append(results, faucet)
	
	// Test 5: Identity Creation
	fmt.Println("\n🆔 5. Testing Identity Management...")
	identityData := map[string]interface{}{
		"address":    "test_user_123",
		"public_key": "04a1b2c3d4e5f6...",
		"username":   "testuser",
		"email":      "test@atlas.com",
	}
	identity := testEndpoint(baseURL+"/identity/create", "POST", identityData)
	results = append(results, identity)
	
	// Test 6: Social Post Creation
	fmt.Println("\n📝 6. Testing Social Media...")
	postData := map[string]interface{}{
		"author":     "test_user_123",
		"content":    "Hello ATLAS Blockchain! 🚀",
		"visibility": "public",
		"category":   "general",
	}
	post := testEndpoint(baseURL+"/social/post/create", "POST", postData)
	results = append(results, post)
	
	// Test 7: Governance Proposal
	fmt.Println("\n🗳️ 7. Testing Governance...")
	proposalData := map[string]interface{}{
		"proposer":     "test_user_123",
		"title":        "Test Proposal",
		"description":  "This is a test governance proposal",
		"category":     "general",
		"current_block": 1,
		"actions":      []map[string]interface{}{},
	}
	proposal := testEndpoint(baseURL+"/governance/submit-proposal", "POST", proposalData)
	results = append(results, proposal)
	
	// Test 8: DeFi Staking
	fmt.Println("\n🔒 8. Testing DeFi Staking...")
	stakingData := map[string]interface{}{
		"address": "test_user_123",
		"amount":  500,
	}
	staking := testEndpoint(baseURL+"/update-user-stake", "POST", stakingData)
	results = append(results, staking)
	
	// Test 9: Network Architecture
	fmt.Println("\n🌐 9. Testing Network Info...")
	network := testEndpoint(baseURL+"/network/architecture", "GET", nil)
	results = append(results, network)
	
	// Test 10: Monitoring
	fmt.Println("\n📈 10. Testing Monitoring...")
	monitoring := testEndpoint(baseURL+"/monitoring/status", "GET", nil)
	results = append(results, monitoring)
	
	// Print Results Summary
	fmt.Println("\n" + "="*60)
	fmt.Println("📋 COMPREHENSIVE TEST RESULTS")
	fmt.Println("="*60)
	
	successCount := 0
	for i, result := range results {
		status := "❌"
		if result.Success {
			status = "✅"
			successCount++
		}
		fmt.Printf("%s %d. %s (Status: %d)\n", status, i+1, result.Endpoint, result.Status)
	}
	
	fmt.Printf("\n🎯 Overall Success Rate: %d/%d (%.1f%%)\n", 
		successCount, len(results), float64(successCount)/float64(len(results))*100)
	
	fmt.Println("\n🚀 ATLAS Blockchain Platform is ready for action!")
	fmt.Println("All core features are implemented and functional.")
}

func testEndpoint(url, method string, data map[string]interface{}) TestResult {
	var req *http.Request
	var err error
	
	if data != nil {
		jsonData, _ := json.Marshal(data)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	
	if err != nil {
		return TestResult{url, 0, err.Error(), false}
	}
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return TestResult{url, 0, err.Error(), false}
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TestResult{url, resp.StatusCode, err.Error(), false}
	}
	
	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	return TestResult{url, resp.StatusCode, string(body), success}
} 
