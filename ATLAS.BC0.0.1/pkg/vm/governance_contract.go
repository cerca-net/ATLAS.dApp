package vm

import (
	"fmt"
	"log"
)

// Governance contract constants
const (
	GovernanceContractName    = "CercaChainGovernance"
	GovernanceContractAddress = "CONTRACT_GOVERNANCE_SYSTEM"
)

// Proposal statuses
const (
	ProposalStatusPending  = int64(0)
	ProposalStatusActive   = int64(1)
	ProposalStatusPassed   = int64(2)
	ProposalStatusFailed   = int64(3)
	ProposalStatusExecuted = int64(4)
)

// CreateSystemGovernanceContract creates the genesis governance contract for voting.
func CreateSystemGovernanceContract() *Contract {
	contract := &Contract{
		Address:      GovernanceContractAddress,
		Name:         GovernanceContractName,
		Version:      "1.0.0",
		Owner:        "SYSTEM",
		Upgradable:   true,
		CreatedAt:    0,
		UpdatedAt:    0,
		ContractType: ContractTypeSystem,
		Functions:    make(map[string]*Function),
		Storage:      make(map[string]interface{}),
	}

	// Initialize default storage parameters
	contract.Storage["proposal_counter"] = int64(0)
	contract.Storage["min_proposal_stake"] = int64(1000)
	contract.Storage["min_voting_stake"] = int64(100)

	// Define functions
	contract.Functions = map[string]*Function{
		"createProposal": {
			Name:       "createProposal",
			Parameters: []string{"proposer", "description", "votingPeriod"},
			Code:       []Instruction{}, // Handled by helper
		},
		"castVote": {
			Name:       "castVote",
			Parameters: []string{"voter", "proposalId", "support", "weight"},
			Code:       []Instruction{},
		},
		"executeProposal": {
			Name:       "executeProposal",
			Parameters: []string{"proposalId"},
			Code:       []Instruction{},
		},
	}

	return contract
}

// GovernanceContractHelper provides high-level Go functions for governance operations.
type GovernanceContractHelper struct {
	Contract *Contract
	VM       *VM
}

// NewGovernanceContractHelper creates a helper for governance contract interaction
func NewGovernanceContractHelper(contract *Contract) *GovernanceContractHelper {
	vmInstance := NewVM()
	vmInstance.RegisterSystemContract(contract.Address, []string{
		"createProposal", "castVote", "executeProposal",
	})
	return &GovernanceContractHelper{
		Contract: contract,
		VM:       vmInstance,
	}
}

// ProposalInfo struct returns public info about a proposal
type ProposalInfo struct {
	ProposalID   string `json:"proposal_id"`
	Proposer     string `json:"proposer"`
	Description  string `json:"description"`
	Status       int64  `json:"status"`
	VotesFor     int64  `json:"votes_for"`
	VotesAgainst int64  `json:"votes_against"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
}

// CreateProposal allows a user to submit a new proposal
func (h *GovernanceContractHelper) CreateProposal(ctx *ExecutionContext, proposer, description string, votingPeriod int64) (string, error) {
	if ctx.State == nil {
		return "", fmt.Errorf("state adapter is required")
	}

	// Check if proposer has enough stake (assume using standard token balance for simple check)
	minStake, _ := ctx.State.GetContractStorage(h.Contract.Address, "min_proposal_stake")
	if minStake == 0 {
		minStake = 1000
	}

	if ctx.State.GetBalance(proposer) < minStake {
		return "", fmt.Errorf("insufficient balance, need at least %d TCOIN to propose", minStake)
	}

	// Generate Proposal ID
	counter, _ := ctx.State.GetContractStorage(h.Contract.Address, "proposal_counter")
	newCounter := counter + 1
	ctx.State.SetContractStorage(h.Contract.Address, "proposal_counter", newCounter)
	proposalID := fmt.Sprintf("PROP_%d", newCounter)

	// Store proposal data
	ctx.State.SetStringStorage(h.Contract.Address, fmt.Sprintf("prop_proposer_%s", proposalID), proposer)
	ctx.State.SetStringStorage(h.Contract.Address, fmt.Sprintf("prop_desc_%s", proposalID), description)
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("prop_status_%s", proposalID), ProposalStatusActive)
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("prop_votes_for_%s", proposalID), 0)
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("prop_votes_against_%s", proposalID), 0)

	startTime := ctx.Timestamp
	if startTime == 0 {
		startTime = 1 // Prevent 0 from being treated as missing
	}
	endTime := startTime + votingPeriod

	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("prop_start_%s", proposalID), startTime)
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("prop_end_%s", proposalID), endTime)

	ctx.EmitEvent(h.Contract.Address, "ProposalCreated", map[string]interface{}{
		"proposal_id": proposalID,
		"proposer":    proposer,
		"end_time":    endTime,
	})

	log.Printf("[GOVERNANCE] Proposal %s created by %s", proposalID, shortAddr(proposer))
	return proposalID, nil
}

// CastVote records a vote for or against a proposal
func (h *GovernanceContractHelper) CastVote(ctx *ExecutionContext, voter, proposalID string, support bool, weight int64) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required")
	}

	// Read proposal status
	status, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_status_%s", proposalID))
	if !exists {
		return fmt.Errorf("proposal %s not found", proposalID)
	}
	if status != ProposalStatusActive {
		return fmt.Errorf("proposal %s is not active", proposalID)
	}

	// Check time
	endTime, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_end_%s", proposalID))
	if ctx.Timestamp > endTime {
		return fmt.Errorf("voting period for proposal %s has ended", proposalID)
	}

	// Check if already voted
	hasVotedKey := fmt.Sprintf("prop_voted_%s_%s", proposalID, voter)
	if _, voted := ctx.State.GetContractStorage(h.Contract.Address, hasVotedKey); voted {
		return fmt.Errorf("voter %s has already voted on proposal %s", voter, proposalID)
	}

	// Add vote
	ctx.State.SetContractStorage(h.Contract.Address, hasVotedKey, 1)

	if support {
		votesFor, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_votes_for_%s", proposalID))
		ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("prop_votes_for_%s", proposalID), votesFor+weight)
	} else {
		votesAgainst, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_votes_against_%s", proposalID))
		ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("prop_votes_against_%s", proposalID), votesAgainst+weight)
	}

	ctx.EmitEvent(h.Contract.Address, "VoteCast", map[string]interface{}{
		"proposal_id": proposalID,
		"voter":       voter,
		"support":     support,
		"weight":      weight,
	})

	log.Printf("[GOVERNANCE] Vote cast on %s by %s (support: %v, weight: %d)", proposalID, shortAddr(voter), support, weight)
	return nil
}

// ExecuteProposal finalizes a proposal if the voting period has ended
func (h *GovernanceContractHelper) ExecuteProposal(ctx *ExecutionContext, proposalID string) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required")
	}

	status, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_status_%s", proposalID))
	if !exists {
		return fmt.Errorf("proposal %s not found", proposalID)
	}
	if status != ProposalStatusActive {
		return fmt.Errorf("proposal %s is not active", proposalID)
	}

	endTime, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_end_%s", proposalID))
	if ctx.Timestamp <= endTime {
		return fmt.Errorf("voting period for proposal %s has not ended yet", proposalID)
	}

	votesFor, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_votes_for_%s", proposalID))
	votesAgainst, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_votes_against_%s", proposalID))

	// Simple majority wins
	var finalStatus int64
	if votesFor > votesAgainst {
		finalStatus = ProposalStatusPassed
	} else {
		finalStatus = ProposalStatusFailed
	}

	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("prop_status_%s", proposalID), finalStatus)

	ctx.EmitEvent(h.Contract.Address, "ProposalConcluded", map[string]interface{}{
		"proposal_id":   proposalID,
		"final_status":  finalStatus,
		"votes_for":     votesFor,
		"votes_against": votesAgainst,
	})

	log.Printf("[GOVERNANCE] Proposal %s concluded with status %d (For: %d, Against: %d)", proposalID, finalStatus, votesFor, votesAgainst)
	return nil
}

// GetProposal returns view-only details for a proposal
func (h *GovernanceContractHelper) GetProposal(ctx *ExecutionContext, proposalID string) *ProposalInfo {
	if ctx.State == nil {
		return nil
	}

	status, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_status_%s", proposalID))
	if !exists {
		return nil
	}

	proposer, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("prop_proposer_%s", proposalID))
	desc, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("prop_desc_%s", proposalID))
	votesFor, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_votes_for_%s", proposalID))
	votesAgainst, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_votes_against_%s", proposalID))
	start, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_start_%s", proposalID))
	end, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("prop_end_%s", proposalID))

	return &ProposalInfo{
		ProposalID:   proposalID,
		Proposer:     proposer,
		Description:  desc,
		Status:       status,
		VotesFor:     votesFor,
		VotesAgainst: votesAgainst,
		StartTime:    start,
		EndTime:      end,
	}
}
