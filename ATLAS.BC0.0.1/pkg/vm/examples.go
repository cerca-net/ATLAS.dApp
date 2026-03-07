package vm

// ExampleContracts provides sample contracts for testing and demonstration

// GetSimpleTokenContract returns a simple token contract example
func GetSimpleTokenContract() *JSONContract {
	return &JSONContract{
		Name:    "SimpleToken",
		Version: "1.0",
		Owner:   "0x1234567890abcdef",
		Functions: map[string]*JSONFunction{
			"transfer": {
				Parameters: []string{"to", "amount"},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "balance"},
					{Op: "PUSH", Value: "amount"},
					{Op: "SUB"},
					{Op: "STORE", Key: "balance"},
				},
			},
			"balance": {
				Parameters: []string{},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "balance"},
					{Op: "RETURN"},
				},
			},
			"mint": {
				Parameters: []string{"amount"},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "balance"},
					{Op: "PUSH", Value: "amount"},
					{Op: "ADD"},
					{Op: "STORE", Key: "balance"},
				},
			},
		},
		Storage: map[string]interface{}{
			"balance": 1000,
			"owner":   "0x1234567890abcdef",
		},
	}
}

// GetVotingContract returns a simple voting contract example
func GetVotingContract() *JSONContract {
	return &JSONContract{
		Name:    "SimpleVoting",
		Version: "1.0",
		Owner:   "0x1234567890abcdef",
		Functions: map[string]*JSONFunction{
			"vote": {
				Parameters: []string{"proposal", "choice"},
				Code: []JSONInstruction{
					{Op: "PUSH", Value: "choice"},
					{Op: "STORE", Key: "vote"},
					{Op: "LOAD", Key: "votes"},
					{Op: "PUSH", Value: 1},
					{Op: "ADD"},
					{Op: "STORE", Key: "votes"},
				},
			},
			"getVotes": {
				Parameters: []string{},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "votes"},
					{Op: "RETURN"},
				},
			},
			"getVote": {
				Parameters: []string{},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "vote"},
					{Op: "RETURN"},
				},
			},
		},
		Storage: map[string]interface{}{
			"votes": 0,
			"vote":  0,
		},
	}
}

// GetEscrowContract returns a simple escrow contract example
func GetEscrowContract() *JSONContract {
	return &JSONContract{
		Name:    "SimpleEscrow",
		Version: "1.0",
		Owner:   "0x1234567890abcdef",
		Functions: map[string]*JSONFunction{
			"deposit": {
				Parameters: []string{"amount"},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "balance"},
					{Op: "PUSH", Value: "amount"},
					{Op: "ADD"},
					{Op: "STORE", Key: "balance"},
				},
			},
			"withdraw": {
				Parameters: []string{"amount"},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "balance"},
					{Op: "PUSH", Value: "amount"},
					{Op: "SUB"},
					{Op: "STORE", Key: "balance"},
				},
			},
			"getBalance": {
				Parameters: []string{},
				Code: []JSONInstruction{
					{Op: "LOAD", Key: "balance"},
					{Op: "RETURN"},
				},
			},
		},
		Storage: map[string]interface{}{
			"balance": 0,
		},
	}
} 