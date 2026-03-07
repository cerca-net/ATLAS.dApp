package vm

import (
	"testing"
	"time"
)

func TestPermissionedSmartContractSystem(t *testing.T) {
	// Create a new VM instance
	vm := NewVM()
	
	// Test 1: Register system contracts
	t.Run("Register System Contracts", func(t *testing.T) {
		// Register voting contract
		votingContract := CreateSystemVotingContract()
		vm.RegisterVotingContract(votingContract.Address, []string{"createProposal", "vote", "getVoteCount"})
		
		// Register governance contract
		governanceContract := CreateSystemGovernanceContract()
		vm.RegisterGovernanceContract(governanceContract.Address, []string{"submitProposal", "executeProposal"}, "SYSTEM")
		
		// Verify contracts are approved
		if !vm.IsContractApproved(votingContract.Address) {
			t.Errorf("Voting contract should be approved")
		}
		if !vm.IsContractApproved(governanceContract.Address) {
			t.Errorf("Governance contract should be approved")
		}
		
		// Verify functions are allowed
		if !vm.IsFunctionAllowed(votingContract.Address, "vote") {
			t.Errorf("Vote function should be allowed for voting contract")
		}
		if !vm.IsFunctionAllowed(governanceContract.Address, "submitProposal") {
			t.Errorf("SubmitProposal function should be allowed for governance contract")
		}
	})
	
	// Test 2: Custom contract approval
	t.Run("Custom Contract Approval", func(t *testing.T) {
		// Create a custom contract
		customFunctions := map[string]*JSONFunction{
			"increment": {
				Parameters: []string{},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "counter"},
					{Op: "PUSH", Value: 1},
					{Op: "ADD"},
					{Op: "STORE", Key: "counter"},
				},
			},
		}
		
		customContract := CreateCustomContract("user123", "TestContract", customFunctions)
		
		// Initially, custom contract should not be approved
		if vm.IsContractApproved(customContract.Address) {
			t.Errorf("Custom contract should not be approved initially")
		}
		
		// Approve the custom contract
		err := vm.ApproveCustomContract(customContract.Address, []string{"increment"}, "admin", 1234567890)
		if err != nil {
			t.Errorf("Failed to approve custom contract: %v", err)
		}
		
		// Now it should be approved
		if !vm.IsContractApproved(customContract.Address) {
			t.Errorf("Custom contract should be approved after approval")
		}
		
		// Check function permissions
		if !vm.IsFunctionAllowed(customContract.Address, "increment") {
			t.Errorf("Increment function should be allowed for approved custom contract")
		}
		
		// Unapproved function should not be allowed
		if vm.IsFunctionAllowed(customContract.Address, "unauthorized_function") {
			t.Errorf("Unauthorized function should not be allowed")
		}
	})
	
	// Test 3: Contract execution
	t.Run("Contract Execution", func(t *testing.T) {
		// Create a simple test contract
		testContract := &Contract{
			Address:      "CONTRACT_TEST123",
			Name:         "TestContract",
			ContractType: ContractTypeSystem,
			Functions: map[string]*Function{
				"add": {
					Parameters: []string{"a", "b"},
					Code: []Instruction{
						{Opcode: "ADD"},
					},
				},
				"store": {
					Parameters: []string{"value"},
					Code: []Instruction{
						{Opcode: "STORE", Operands: []interface{}{"test_value"}},
					},
				},
				"load": {
					Parameters: []string{},
					Code: []Instruction{
						{Opcode: "LOAD", Operands: []interface{}{"test_value"}},
					},
				},
			},
		}
		
		// Register the contract
		vm.RegisterSystemContract(testContract.Address, []string{"add", "store", "load"})
		
		// Test function execution
		context := NewExecutionContext("caller123", 1000)
		context.SetContractContext(testContract.Address, "add")
		
		// Execute add function with parameters
		err := testContract.CallFunction("add", []interface{}{5, 3}, vm, context)
		if err != nil {
			t.Errorf("Failed to execute add function: %v", err)
		}
		
		// Check result
		if len(vm.stack) != 1 || vm.stack[0] != 8 {
			t.Errorf("Expected result 8, got %v", vm.stack)
		}
	})
	
	// Test 4: Gas metering
	t.Run("Gas Metering", func(t *testing.T) {
		vm := NewVM()
		context := NewExecutionContext("caller123", 10) // Low gas limit
		
		// Create a contract that uses more gas than available
		testContract := &Contract{
			Address:      "CONTRACT_GAS_TEST",
			ContractType: ContractTypeSystem,
			Functions: map[string]*Function{
				"expensive": {
					Parameters: []string{},
					Code: []Instruction{
						{Opcode: "PUSH", Operands: []interface{}{1}},
						{Opcode: "PUSH", Operands: []interface{}{2}},
						{Opcode: "PUSH", Operands: []interface{}{3}},
						{Opcode: "PUSH", Operands: []interface{}{4}},
						{Opcode: "PUSH", Operands: []interface{}{5}},
					},
				},
			},
		}
		
		vm.RegisterSystemContract(testContract.Address, []string{"expensive"})
		context.SetContractContext(testContract.Address, "expensive")
		
		// This should fail due to gas limit
		err := testContract.CallFunction("expensive", []interface{}{}, vm, context)
		if err == nil {
			t.Errorf("Expected gas limit error, but execution succeeded")
		}
	})
	
	// Test 5: Call depth protection
	t.Run("Call Depth Protection", func(t *testing.T) {
		vm := NewVM()
		vm.maxCallDepth = 2 // Set low call depth limit
		
		// Create a contract with recursive function
		testContract := &Contract{
			Address:      "CONTRACT_RECURSIVE",
			ContractType: ContractTypeSystem,
			Functions: map[string]*Function{
				"recursive": {
					Parameters: []string{},
					Code: []Instruction{
						{Opcode: "CALL", Operands: []interface{}{"recursive"}},
					},
				},
			},
		}
		
		vm.RegisterSystemContract(testContract.Address, []string{"recursive"})
		context := NewExecutionContext("caller123", 1000)
		context.SetContractContext(testContract.Address, "recursive")
		
		// This should fail due to call depth limit
		err := testContract.CallFunction("recursive", []interface{}{}, vm, context)
		if err == nil {
			t.Errorf("Expected call depth error, but execution succeeded")
		}
	})
	
	// Test 6: Contract address generation
	t.Run("Contract Address Generation", func(t *testing.T) {
		owner1 := "owner1"
		
		// Generate addresses for same owner
		addr1 := generateContractAddress(owner1, nil)
		
		// Add a small delay to ensure different timestamps
		time.Sleep(1 * time.Millisecond)
		
		addr2 := generateContractAddress(owner1, nil)
		
		// Addresses should be different due to timestamp
		if addr1 == addr2 {
			t.Errorf("Contract addresses should be unique: %s == %s", addr1, addr2)
		}
		
		// Addresses should start with CONTRACT_
		if addr1[:9] != "CONTRACT_" {
			t.Errorf("Contract address should start with CONTRACT_: %s", addr1)
		}
	})
	
	// Test 7: System contract permissions
	t.Run("System Contract Permissions", func(t *testing.T) {
		vm := NewVM()
		
		// System contracts should have full access
		systemContract := CreateSystemVotingContract()
		vm.RegisterVotingContract(systemContract.Address, []string{})
		
		// Any function should be allowed for system contracts
		if !vm.IsFunctionAllowed(systemContract.Address, "any_function_name") {
			t.Errorf("System contracts should have full function access")
		}
		
		// Governance contracts should also have full access
		governanceContract := CreateSystemGovernanceContract()
		vm.RegisterGovernanceContract(governanceContract.Address, []string{}, "SYSTEM")
		
		if !vm.IsFunctionAllowed(governanceContract.Address, "any_function_name") {
			t.Errorf("Governance contracts should have full function access")
		}
	})
}

func TestVMInstructions(t *testing.T) {
	vm := NewVM()
	context := NewExecutionContext("test", 1000)
	
	t.Run("Basic Instructions", func(t *testing.T) {
		// Test PUSH and ADD
		instructions := []Instruction{
			{Opcode: "PUSH", Operands: []interface{}{5}},
			{Opcode: "PUSH", Operands: []interface{}{3}},
			{Opcode: "ADD"},
		}
		
		err := vm.Execute(instructions, context)
		if err != nil {
			t.Errorf("Failed to execute basic instructions: %v", err)
		}
		
		if len(vm.stack) != 1 || vm.stack[0] != 8 {
			t.Errorf("Expected result 8, got %v", vm.stack)
		}
	})
	
	t.Run("Memory Operations", func(t *testing.T) {
		vm := NewVM()
		context := NewExecutionContext("test", 1000)
		
		instructions := []Instruction{
			{Opcode: "PUSH", Operands: []interface{}{42}},
			{Opcode: "STORE", Operands: []interface{}{"test_key"}},
			{Opcode: "LOAD", Operands: []interface{}{"test_key"}},
		}
		
		err := vm.Execute(instructions, context)
		if err != nil {
			t.Errorf("Failed to execute memory operations: %v", err)
		}
		
		if len(vm.stack) != 1 || vm.stack[0] != 42 {
			t.Errorf("Expected result 42, got %v", vm.stack)
		}
	})
	
	t.Run("Control Flow", func(t *testing.T) {
		vm := NewVM()
		context := NewExecutionContext("test", 1000)
		
		// Simple test: if condition is true, jump over the next instruction
		instructions := []Instruction{
			{Opcode: "PUSH", Operands: []interface{}{1}}, // Condition (true)
			{Opcode: "JUMPIF", Operands: []interface{}{2}}, // Jump to instruction 2
			{Opcode: "PUSH", Operands: []interface{}{42}}, // This should be executed
		}
		
		err := vm.Execute(instructions, context)
		if err != nil {
			t.Errorf("Failed to execute control flow: %v", err)
		}
		
		// Should only have 42 on the stack
		if len(vm.stack) != 1 || vm.stack[0] != 42 {
			t.Errorf("Expected result 42, got %v", vm.stack)
		}
	})
} 