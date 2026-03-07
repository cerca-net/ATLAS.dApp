package vm

// StateAdapter is the interface the VM uses to interact with blockchain state.
// This decouples the VM from the StateManager, allowing the VM to read/write
// balances, nonces, and contract storage without importing internal packages.
type StateAdapter interface {
	// GetBalance returns the TCOIN balance for an address
	GetBalance(address string) int64

	// SetBalance sets the TCOIN balance for an address
	SetBalance(address string, amount int64)

	// Transfer moves tokens from one address to another.
	// Returns an error if the sender has insufficient funds.
	Transfer(from, to string, amount int64) error

	// Mint creates new tokens and credits them to an address.
	// Returns an error if minting would exceed max supply.
	Mint(to string, amount int64) error

	// Burn destroys tokens from an address.
	// Returns an error if the address has insufficient funds.
	Burn(from string, amount int64) error

	// GetContractStorage reads a value from a contract's persistent storage
	GetContractStorage(contractAddress, key string) (int64, bool)

	// SetContractStorage writes a value to a contract's persistent storage
	SetContractStorage(contractAddress, key string, value int64)

	// GetStringStorage reads a string value from a contract's persistent storage
	GetStringStorage(contractAddress, key string) (string, bool)

	// SetStringStorage writes a string value to a contract's persistent storage
	SetStringStorage(contractAddress, key string, value string)
}

// EventLog represents an emitted event from contract execution
type EventLog struct {
	ContractAddress string      `json:"contract_address"`
	EventName       string      `json:"event_name"`
	Data            interface{} `json:"data"`
	BlockHeight     int64       `json:"block_height"`
	Timestamp       int64       `json:"timestamp"`
}
