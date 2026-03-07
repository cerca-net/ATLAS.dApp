package vm

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Contract represents a deployed smart contract.
type Contract struct {
	Address      string
	Name         string
	Version      string
	Code         []Instruction
	Functions    map[string]*Function
	Storage      map[string]interface{}
	Owner        string
	Upgradable   bool
	CreatedAt    int64
	UpdatedAt    int64
	ContractType ContractType // Type of contract (system, governance, custom, voting)
}

// Function represents a contract function
type Function struct {
	Name       string
	Parameters []string
	Code       []Instruction
}

// JSONContract represents a contract in JSON format
type JSONContract struct {
	Name         string                   `json:"name"`
	Version      string                   `json:"version"`
	Owner        string                   `json:"owner"`
	Functions    map[string]*JSONFunction `json:"functions"`
	Storage      map[string]interface{}   `json:"storage"`
	ContractType ContractType             `json:"contract_type,omitempty"`
}

// JSONFunction represents a function in JSON format
type JSONFunction struct {
	Parameters []string          `json:"params"`
	Code       []JSONInstruction `json:"code"`
}

// JSONInstruction represents an instruction in JSON format
type JSONInstruction struct {
	Op    string      `json:"op"`
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// DeployContract deploys a new contract to the blockchain.
func DeployContract(owner string, code []Instruction, upgradable bool, contractType ContractType) *Contract {
	now := time.Now().Unix()
	contract := &Contract{
		Address:      generateContractAddress(owner, code),
		Code:         code,
		Functions:    make(map[string]*Function),
		Storage:      make(map[string]interface{}),
		Owner:        owner,
		Upgradable:   upgradable,
		CreatedAt:    now,
		UpdatedAt:    now,
		ContractType: contractType,
	}
	return contract
}

// DeployJSONContract deploys a contract from JSON format
func DeployJSONContract(owner string, jsonContract *JSONContract, upgradable bool) (*Contract, error) {
	contract := &Contract{
		Address:      generateContractAddress(owner, nil), // Will be updated
		Name:         jsonContract.Name,
		Version:      jsonContract.Version,
		Functions:    make(map[string]*Function),
		Storage:      jsonContract.Storage,
		Owner:        owner,
		Upgradable:   upgradable,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
		ContractType: jsonContract.ContractType,
	}

	// Parse functions
	for funcName, jsonFunc := range jsonContract.Functions {
		function := &Function{
			Name:       funcName,
			Parameters: jsonFunc.Parameters,
			Code:       make([]Instruction, 0),
		}

		// Convert JSON instructions to VM instructions
		for _, jsonInstr := range jsonFunc.Code {
			instr := Instruction{
				Opcode: jsonInstr.Op,
			}

			// Add operands based on instruction type
			switch jsonInstr.Op {
			case "PUSH":
				instr.Operands = []interface{}{jsonInstr.Value}
			case "STORE", "LOAD":
				instr.Operands = []interface{}{jsonInstr.Key}
			case "JUMP", "JUMPIF":
				if jsonInstr.Value != nil {
					instr.Operands = []interface{}{jsonInstr.Value}
				}
			case "CALL":
				if jsonInstr.Key != "" {
					instr.Operands = []interface{}{jsonInstr.Key}
				}
			}

			function.Code = append(function.Code, instr)
		}

		contract.Functions[funcName] = function
	}

	// Generate proper address based on contract content
	contract.Address = generateContractAddress(owner, contract.Code)

	return contract, nil
}

// CallFunction executes a specific function on a contract
func (c *Contract) CallFunction(functionName string, params []interface{}, vm *VM, context *ExecutionContext) error {
	// Set the current contract context in the VM
	vm.currentContract = c

	// Set the contract context in the execution context
	context.SetContractContext(c.Address, functionName)

	function, exists := c.Functions[functionName]
	if !exists {
		return fmt.Errorf("function '%s' not found in contract", functionName)
	}

	// Validate parameter count
	if len(params) != len(function.Parameters) {
		return fmt.Errorf("function '%s' expects %d parameters, got %d", functionName, len(function.Parameters), len(params))
	}

	// Push parameters onto stack
	for _, param := range params {
		if val, ok := toInt64(param); ok {
			vm.stack = append(vm.stack, val)
		} else {
			return fmt.Errorf("invalid parameter type for function '%s'", functionName)
		}
	}

	// Execute function code
	return vm.Execute(function.Code, context)
}

// generateContractAddress creates a unique address for a contract.
func generateContractAddress(owner string, code []Instruction) string {
	// Create a hash of owner + timestamp + code for uniqueness
	data := owner + time.Now().String()
	if code != nil {
		for _, instr := range code {
			data += instr.Opcode
			for _, operand := range instr.Operands {
				data += fmt.Sprintf("%v", operand)
			}
		}
	}

	hash := sha256.Sum256([]byte(data))
	return "CONTRACT_" + hex.EncodeToString(hash[:8])
}

// CreateSystemVotingContract creates a pre-approved voting system contract
func CreateSystemVotingContract() *Contract {
	// Simple voting contract with basic functionality
	votingContract := &JSONContract{
		Name:         "SystemVoting",
		Version:      "1.0.0",
		Owner:        "SYSTEM",
		ContractType: ContractTypeVoting,
		Functions: map[string]*JSONFunction{
			"createProposal": {
				Parameters: []string{"proposalId", "description"},
				Code: []JSONInstruction{
					{Op: "PUSH", Value: 1},
					{Op: "STORE", Key: "proposal_active"},
				},
			},
			"vote": {
				Parameters: []string{"proposalId", "vote"},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "votes_for"},
					{Op: "PUSH", Value: 1},
					{Op: "ADD"},
					{Op: "STORE", Key: "votes_for"},
				},
			},
			"getVoteCount": {
				Parameters: []string{"proposalId"},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "votes_for"},
				},
			},
		},
		Storage: map[string]interface{}{
			"votes_for":       0,
			"votes_against":   0,
			"proposal_active": 0,
		},
	}

	contract, _ := DeployJSONContract("SYSTEM", votingContract, false)
	return contract
}

// CreateCustomContract creates a custom contract that requires approval
func CreateCustomContract(owner string, name string, functions map[string]*JSONFunction) *Contract {
	customContract := &JSONContract{
		Name:         name,
		Version:      "1.0.0",
		Owner:        owner,
		ContractType: ContractTypeCustom,
		Functions:    functions,
		Storage:      make(map[string]interface{}),
	}

	contract, _ := DeployJSONContract(owner, customContract, true)
	return contract
}
