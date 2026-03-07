package vm

import (
	"fmt"
	"testing"
)

// MockStateAdapter implements StateAdapter for testing
type MockStateAdapter struct {
	balances        map[string]int64
	contractStorage map[string]map[string]int64
	stringStorage   map[string]map[string]string
	totalSupply     int64
	maxSupply       int64
}

func NewMockStateAdapter() *MockStateAdapter {
	return &MockStateAdapter{
		balances:        make(map[string]int64),
		contractStorage: make(map[string]map[string]int64),
		stringStorage:   make(map[string]map[string]string),
		maxSupply:       TokenMaxSupply,
	}
}

func (m *MockStateAdapter) GetBalance(address string) int64 {
	return m.balances[address]
}

func (m *MockStateAdapter) SetBalance(address string, amount int64) {
	m.balances[address] = amount
}

func (m *MockStateAdapter) Transfer(from, to string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if m.balances[from] < amount {
		return fmt.Errorf("insufficient balance: %d < %d", m.balances[from], amount)
	}
	m.balances[from] -= amount
	m.balances[to] += amount
	return nil
}

func (m *MockStateAdapter) Mint(to string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if m.totalSupply+amount > m.maxSupply {
		return fmt.Errorf("exceeds max supply")
	}
	m.totalSupply += amount
	m.balances[to] += amount
	return nil
}

func (m *MockStateAdapter) Burn(from string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	if m.balances[from] < amount {
		return fmt.Errorf("insufficient balance: %d < %d", m.balances[from], amount)
	}
	m.balances[from] -= amount
	m.totalSupply -= amount
	return nil
}

func (m *MockStateAdapter) GetContractStorage(contractAddress, key string) (int64, bool) {
	if cm, ok := m.contractStorage[contractAddress]; ok {
		if v, exists := cm[key]; exists {
			return v, true
		}
	}
	return 0, false
}

func (m *MockStateAdapter) SetContractStorage(contractAddress, key string, value int64) {
	if _, ok := m.contractStorage[contractAddress]; !ok {
		m.contractStorage[contractAddress] = make(map[string]int64)
	}
	m.contractStorage[contractAddress][key] = value
}

func (m *MockStateAdapter) GetStringStorage(contractAddress, key string) (string, bool) {
	if cm, ok := m.stringStorage[contractAddress]; ok {
		if v, exists := cm[key]; exists {
			return v, true
		}
	}
	return "", false
}

func (m *MockStateAdapter) SetStringStorage(contractAddress, key string, value string) {
	if _, ok := m.stringStorage[contractAddress]; !ok {
		m.stringStorage[contractAddress] = make(map[string]string)
	}
	m.stringStorage[contractAddress][key] = value
}

func TestTokenContract(t *testing.T) {
	t.Run("CreateTokenContract", func(t *testing.T) {
		contract := CreateSystemTokenContract("0xTREASURY")
		if contract == nil {
			t.Fatal("Failed to create token contract")
		}
		if contract.Address != TokenContractAddress {
			t.Errorf("Expected address %s, got %s", TokenContractAddress, contract.Address)
		}
		if contract.Name != TokenContractName {
			t.Errorf("Expected name %s, got %s", TokenContractName, contract.Name)
		}
		if contract.ContractType != ContractTypeSystem {
			t.Errorf("Expected system contract type")
		}
		if contract.Upgradable {
			t.Errorf("Token contract should not be upgradable")
		}

		// Verify functions exist
		expectedFunctions := []string{"transfer", "balanceOf", "mint", "burn", "faucetRequest"}
		for _, fn := range expectedFunctions {
			if _, ok := contract.Functions[fn]; !ok {
				t.Errorf("Missing function: %s", fn)
			}
		}
	})

	t.Run("TokenContractHelper_Transfer", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xALICE"] = 1000
		state.balances["0xBOB"] = 500

		contract := CreateSystemTokenContract("0xTREASURY")
		helper := NewTokenContractHelper(contract)
		ctx := NewExecutionContextWithState("0xALICE", 100000, state, 1, 1000)

		// Transfer 200 from Alice to Bob
		err := helper.Transfer(ctx, "0xALICE", "0xBOB", 200)
		if err != nil {
			t.Fatalf("Transfer failed: %v", err)
		}

		if state.balances["0xALICE"] != 800 {
			t.Errorf("Expected Alice balance 800, got %d", state.balances["0xALICE"])
		}
		if state.balances["0xBOB"] != 700 {
			t.Errorf("Expected Bob balance 700, got %d", state.balances["0xBOB"])
		}

		// Verify event was emitted
		if len(ctx.Events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(ctx.Events))
		}
		if ctx.Events[0].EventName != "Transfer" {
			t.Errorf("Expected Transfer event, got %s", ctx.Events[0].EventName)
		}
	})

	t.Run("TokenContractHelper_Transfer_InsufficientFunds", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xALICE"] = 100

		contract := CreateSystemTokenContract("0xTREASURY")
		helper := NewTokenContractHelper(contract)
		ctx := NewExecutionContextWithState("0xALICE", 100000, state, 1, 1000)

		err := helper.Transfer(ctx, "0xALICE", "0xBOB", 200)
		if err == nil {
			t.Fatal("Expected insufficient funds error")
		}
	})

	t.Run("TokenContractHelper_Mint", func(t *testing.T) {
		state := NewMockStateAdapter()
		// Initialize contract storage for total_supply tracking
		state.SetContractStorage(TokenContractAddress, "total_supply", 0)

		contract := CreateSystemTokenContract("0xTREASURY")
		helper := NewTokenContractHelper(contract)
		ctx := NewExecutionContextWithState("SYSTEM", 100000, state, 1, 1000)

		// Mint 500 to Alice
		err := helper.Mint(ctx, "0xALICE", 500)
		if err != nil {
			t.Fatalf("Mint failed: %v", err)
		}

		if state.balances["0xALICE"] != 500 {
			t.Errorf("Expected Alice balance 500, got %d", state.balances["0xALICE"])
		}

		// Check total supply updated
		totalSupply, _ := state.GetContractStorage(TokenContractAddress, "total_supply")
		if totalSupply != 500 {
			t.Errorf("Expected total supply 500, got %d", totalSupply)
		}
	})

	t.Run("TokenContractHelper_Mint_ExceedsMaxSupply", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.SetContractStorage(TokenContractAddress, "total_supply", TokenMaxSupply-100)

		contract := CreateSystemTokenContract("0xTREASURY")
		helper := NewTokenContractHelper(contract)
		ctx := NewExecutionContextWithState("SYSTEM", 100000, state, 1, 1000)

		// Try to mint more than remaining supply
		err := helper.Mint(ctx, "0xALICE", 200)
		if err == nil {
			t.Fatal("Expected max supply exceeded error")
		}
	})

	t.Run("TokenContractHelper_Burn", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xALICE"] = 1000
		state.SetContractStorage(TokenContractAddress, "total_supply", 5000)

		contract := CreateSystemTokenContract("0xTREASURY")
		helper := NewTokenContractHelper(contract)
		ctx := NewExecutionContextWithState("0xALICE", 100000, state, 1, 1000)

		err := helper.Burn(ctx, "0xALICE", 300)
		if err != nil {
			t.Fatalf("Burn failed: %v", err)
		}

		if state.balances["0xALICE"] != 700 {
			t.Errorf("Expected Alice balance 700, got %d", state.balances["0xALICE"])
		}

		// Check total supply decreased
		totalSupply, _ := state.GetContractStorage(TokenContractAddress, "total_supply")
		if totalSupply != 4700 {
			t.Errorf("Expected total supply 4700, got %d", totalSupply)
		}
	})

	t.Run("TokenContractHelper_FaucetRequest", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xTREASURY"] = 1_000_000

		// Initialize contract storage
		state.SetContractStorage(TokenContractAddress, "faucet_amount", TokenFaucetAmount)
		state.SetContractStorage(TokenContractAddress, "faucet_cooldown", TokenFaucetCooldown)
		state.SetStringStorage(TokenContractAddress, "treasury_address", "0xTREASURY")

		contract := CreateSystemTokenContract("0xTREASURY")
		helper := NewTokenContractHelper(contract)
		ctx := NewExecutionContextWithState("0xUSER1", 100000, state, 1, 1000)

		// Request faucet
		amount, err := helper.FaucetRequest(ctx, "0xUSER1")
		if err != nil {
			t.Fatalf("FaucetRequest failed: %v", err)
		}
		if amount != TokenFaucetAmount {
			t.Errorf("Expected faucet amount %d, got %d", TokenFaucetAmount, amount)
		}

		if state.balances["0xUSER1"] != TokenFaucetAmount {
			t.Errorf("Expected user balance %d, got %d", TokenFaucetAmount, state.balances["0xUSER1"])
		}
		if state.balances["0xTREASURY"] != 1_000_000-TokenFaucetAmount {
			t.Errorf("Expected treasury balance %d, got %d", 1_000_000-TokenFaucetAmount, state.balances["0xTREASURY"])
		}
	})

	t.Run("TokenContractHelper_FaucetRequest_Cooldown", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.balances["0xTREASURY"] = 1_000_000

		state.SetContractStorage(TokenContractAddress, "faucet_amount", TokenFaucetAmount)
		state.SetContractStorage(TokenContractAddress, "faucet_cooldown", TokenFaucetCooldown)
		state.SetStringStorage(TokenContractAddress, "treasury_address", "0xTREASURY")

		contract := CreateSystemTokenContract("0xTREASURY")
		helper := NewTokenContractHelper(contract)

		// First request at timestamp 1000
		ctx1 := NewExecutionContextWithState("0xUSER1", 100000, state, 1, 1000)
		_, err := helper.FaucetRequest(ctx1, "0xUSER1")
		if err != nil {
			t.Fatalf("First FaucetRequest failed: %v", err)
		}

		// Second request too soon (within cooldown)
		ctx2 := NewExecutionContextWithState("0xUSER1", 100000, state, 2, 2000)
		_, err = helper.FaucetRequest(ctx2, "0xUSER1")
		if err == nil {
			t.Fatal("Expected cooldown error on second request")
		}

		// Third request after cooldown expires
		ctx3 := NewExecutionContextWithState("0xUSER1", 100000, state, 100, 1000+TokenFaucetCooldown+1)
		_, err = helper.FaucetRequest(ctx3, "0xUSER1")
		if err != nil {
			t.Fatalf("Third FaucetRequest should succeed after cooldown: %v", err)
		}
	})

	t.Run("TokenContractHelper_GetTotalSupply", func(t *testing.T) {
		state := NewMockStateAdapter()
		state.SetContractStorage(TokenContractAddress, "total_supply", 42000)

		contract := CreateSystemTokenContract("0xTREASURY")
		helper := NewTokenContractHelper(contract)
		ctx := NewExecutionContextWithState("test", 100000, state, 1, 1000)

		supply := helper.GetTotalSupply(ctx)
		if supply != 42000 {
			t.Errorf("Expected total supply 42000, got %d", supply)
		}
	})
}

func TestNewOpcodes(t *testing.T) {
	t.Run("TRANSFER_Opcode", func(t *testing.T) {
		vmInst := NewVM()
		state := NewMockStateAdapter()
		state.balances["0xALICE"] = 1000
		state.balances["0xBOB"] = 500

		ctx := NewExecutionContextWithState("0xALICE", 100000, state, 1, 1000)

		instructions := []Instruction{
			{Opcode: "PUSHS", Operands: []interface{}{"0xALICE"}}, // from
			{Opcode: "PUSHS", Operands: []interface{}{"0xBOB"}},   // to
			{Opcode: "PUSH", Operands: []interface{}{200}},        // amount
			{Opcode: "TRANSFER"},
		}

		err := vmInst.Execute(instructions, ctx)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if state.balances["0xALICE"] != 800 {
			t.Errorf("Expected Alice 800, got %d", state.balances["0xALICE"])
		}
		if state.balances["0xBOB"] != 700 {
			t.Errorf("Expected Bob 700, got %d", state.balances["0xBOB"])
		}

		// TRANSFER pushes 1 on success
		if len(vmInst.stack) != 1 || vmInst.stack[0] != 1 {
			t.Errorf("Expected [1] on stack, got %v", vmInst.stack)
		}
	})

	t.Run("BALANCE_Opcode", func(t *testing.T) {
		vmInst := NewVM()
		state := NewMockStateAdapter()
		state.balances["0xALICE"] = 12345

		ctx := NewExecutionContextWithState("test", 100000, state, 1, 1000)

		instructions := []Instruction{
			{Opcode: "PUSHS", Operands: []interface{}{"0xALICE"}},
			{Opcode: "BALANCE"},
		}

		err := vmInst.Execute(instructions, ctx)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if len(vmInst.stack) != 1 || vmInst.stack[0] != 12345 {
			t.Errorf("Expected [12345] on stack, got %v", vmInst.stack)
		}
	})

	t.Run("MINT_Opcode", func(t *testing.T) {
		vmInst := NewVM()
		state := NewMockStateAdapter()

		ctx := NewExecutionContextWithState("SYSTEM", 100000, state, 1, 1000)

		instructions := []Instruction{
			{Opcode: "PUSHS", Operands: []interface{}{"0xALICE"}},
			{Opcode: "PUSH", Operands: []interface{}{500}},
			{Opcode: "MINT"},
		}

		err := vmInst.Execute(instructions, ctx)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if state.balances["0xALICE"] != 500 {
			t.Errorf("Expected Alice 500, got %d", state.balances["0xALICE"])
		}
	})

	t.Run("BURN_Opcode", func(t *testing.T) {
		vmInst := NewVM()
		state := NewMockStateAdapter()
		state.balances["0xALICE"] = 1000
		state.totalSupply = 1000

		ctx := NewExecutionContextWithState("0xALICE", 100000, state, 1, 1000)

		instructions := []Instruction{
			{Opcode: "PUSHS", Operands: []interface{}{"0xALICE"}},
			{Opcode: "PUSH", Operands: []interface{}{300}},
			{Opcode: "BURN"},
		}

		err := vmInst.Execute(instructions, ctx)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if state.balances["0xALICE"] != 700 {
			t.Errorf("Expected Alice 700, got %d", state.balances["0xALICE"])
		}
	})

	t.Run("CALLER_TIMESTAMP_BLOCKNUM", func(t *testing.T) {
		vmInst := NewVM()
		state := NewMockStateAdapter()

		ctx := NewExecutionContextWithState("0xCALLER", 100000, state, 42, 9999)

		instructions := []Instruction{
			{Opcode: "CALLER"},
			{Opcode: "TIMESTAMP"},
			{Opcode: "BLOCKNUM"},
		}

		err := vmInst.Execute(instructions, ctx)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		// String stack should have caller
		if len(vmInst.stringStack) != 1 || vmInst.stringStack[0] != "0xCALLER" {
			t.Errorf("Expected caller on string stack, got %v", vmInst.stringStack)
		}

		// Int stack should have timestamp and blocknum
		if len(vmInst.stack) != 2 {
			t.Fatalf("Expected 2 values on stack, got %d", len(vmInst.stack))
		}
		if vmInst.stack[0] != 9999 {
			t.Errorf("Expected timestamp 9999, got %d", vmInst.stack[0])
		}
		if vmInst.stack[1] != 42 {
			t.Errorf("Expected blocknum 42, got %d", vmInst.stack[1])
		}
	})

	t.Run("REQUIRE_Success", func(t *testing.T) {
		vmInst := NewVM()
		ctx := NewExecutionContext("test", 100000)

		instructions := []Instruction{
			{Opcode: "PUSH", Operands: []interface{}{1}}, // true
			{Opcode: "REQUIRE", Operands: []interface{}{"should pass"}},
		}

		err := vmInst.Execute(instructions, ctx)
		if err != nil {
			t.Fatalf("REQUIRE should pass: %v", err)
		}
	})

	t.Run("REQUIRE_Failure", func(t *testing.T) {
		vmInst := NewVM()
		ctx := NewExecutionContext("test", 100000)

		instructions := []Instruction{
			{Opcode: "PUSH", Operands: []interface{}{0}}, // false
			{Opcode: "REQUIRE", Operands: []interface{}{"insufficient balance"}},
		}

		err := vmInst.Execute(instructions, ctx)
		if err == nil {
			t.Fatal("REQUIRE should fail")
		}
	})

	t.Run("EMIT_Opcode", func(t *testing.T) {
		vmInst := NewVM()
		ctx := NewExecutionContextWithState("test", 100000, NewMockStateAdapter(), 5, 2000)

		instructions := []Instruction{
			{Opcode: "PUSH", Operands: []interface{}{42}},
			{Opcode: "EMIT", Operands: []interface{}{"Transfer"}},
		}

		err := vmInst.Execute(instructions, ctx)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if len(ctx.Events) != 1 {
			t.Fatalf("Expected 1 event, got %d", len(ctx.Events))
		}
		if ctx.Events[0].EventName != "Transfer" {
			t.Errorf("Expected event name 'Transfer', got '%s'", ctx.Events[0].EventName)
		}
	})

	t.Run("SSTORE_SLOAD", func(t *testing.T) {
		vmInst := NewVM()
		state := NewMockStateAdapter()
		ctx := NewExecutionContextWithState("test", 100000, state, 1, 1000)

		instructions := []Instruction{
			{Opcode: "PUSH", Operands: []interface{}{42}},
			{Opcode: "SSTORE", Operands: []interface{}{"my_key"}},
			{Opcode: "SLOAD", Operands: []interface{}{"my_key"}},
		}

		err := vmInst.Execute(instructions, ctx)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if len(vmInst.stack) != 1 || vmInst.stack[0] != 42 {
			t.Errorf("Expected [42] on stack, got %v", vmInst.stack)
		}

		// Verify it's persisted in state
		val, ok := state.GetContractStorage("", "my_key")
		if !ok || val != 42 {
			t.Errorf("Expected persisted value 42, got %d (exists: %v)", val, ok)
		}
	})

	t.Run("SSTORE_S_SLOAD_S", func(t *testing.T) {
		vmInst := NewVM()
		state := NewMockStateAdapter()
		ctx := NewExecutionContextWithState("test", 100000, state, 1, 1000)

		instructions := []Instruction{
			{Opcode: "PUSHS", Operands: []interface{}{"0xADDRESS123"}},
			{Opcode: "SSTORE_S", Operands: []interface{}{"owner"}},
			{Opcode: "SLOAD_S", Operands: []interface{}{"owner"}},
		}

		err := vmInst.Execute(instructions, ctx)
		if err != nil {
			t.Fatalf("Execute failed: %v", err)
		}

		if len(vmInst.stringStack) != 1 || vmInst.stringStack[0] != "0xADDRESS123" {
			t.Errorf("Expected ['0xADDRESS123'] on string stack, got %v", vmInst.stringStack)
		}
	})
}
