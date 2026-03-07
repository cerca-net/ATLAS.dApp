package vm

import (
	"time"
)

// ExecutionContext provides context for contract execution (caller, value, gas, etc.)
type ExecutionContext struct {
	Caller          string
	Value           int64
	GasLimit        uint64
	GasUsed         uint64
	ContractAddress string        // Address of the contract being executed
	FunctionName    string        // Name of the function being called
	Parameters      []interface{} // Function parameters
	Timestamp       int64         // Block timestamp
	BlockHeight     int64         // Current block height

	// State adapter for interacting with blockchain state
	State StateAdapter

	// Event logs emitted during execution
	Events []*EventLog
}

// NewExecutionContext creates a new execution context
func NewExecutionContext(caller string, gasLimit uint64) *ExecutionContext {
	return &ExecutionContext{
		Caller:      caller,
		GasLimit:    gasLimit,
		GasUsed:     0,
		Timestamp:   time.Now().Unix(),
		BlockHeight: 0, // Will be set by the blockchain
		Parameters:  make([]interface{}, 0),
		Events:      make([]*EventLog, 0),
	}
}

// NewExecutionContextWithState creates a new execution context with a state adapter
func NewExecutionContextWithState(caller string, gasLimit uint64, state StateAdapter, blockHeight int64, timestamp int64) *ExecutionContext {
	return &ExecutionContext{
		Caller:      caller,
		GasLimit:    gasLimit,
		GasUsed:     0,
		Timestamp:   timestamp,
		BlockHeight: blockHeight,
		Parameters:  make([]interface{}, 0),
		Events:      make([]*EventLog, 0),
		State:       state,
	}
}

// SetContractContext sets the contract context for execution
func (ctx *ExecutionContext) SetContractContext(contractAddress, functionName string) {
	ctx.ContractAddress = contractAddress
	ctx.FunctionName = functionName
}

// AddParameter adds a parameter to the function call
func (ctx *ExecutionContext) AddParameter(param interface{}) {
	ctx.Parameters = append(ctx.Parameters, param)
}

// ChargeGas deducts gas for an operation.
func (ctx *ExecutionContext) ChargeGas(amount uint64) bool {
	if ctx.GasUsed+amount > ctx.GasLimit {
		return false // Out of gas
	}
	ctx.GasUsed += amount
	return true
}

// GetRemainingGas returns the remaining gas
func (ctx *ExecutionContext) GetRemainingGas() uint64 {
	if ctx.GasUsed >= ctx.GasLimit {
		return 0
	}
	return ctx.GasLimit - ctx.GasUsed
}

// EmitEvent records an event log during contract execution
func (ctx *ExecutionContext) EmitEvent(contractAddress, eventName string, data interface{}) {
	ctx.Events = append(ctx.Events, &EventLog{
		ContractAddress: contractAddress,
		EventName:       eventName,
		Data:            data,
		BlockHeight:     ctx.BlockHeight,
		Timestamp:       ctx.Timestamp,
	})
}
