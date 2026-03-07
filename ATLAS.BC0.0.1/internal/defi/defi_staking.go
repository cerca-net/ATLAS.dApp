package defi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// StakingPool manages staking operations
type StakingPool struct {
	stakers       map[string]*Staker
	totalStaked   uint64
	rewardRate    float64
	lockPeriod    time.Duration
	lastReward    time.Time
}

// Staker represents a staking position
type Staker struct {
	Address     string    `json:"address"`
	Amount      uint64    `json:"amount"`
	Rewards     uint64    `json:"rewards"`
	StakedAt    time.Time `json:"staked_at"`
	LastClaim   time.Time `json:"last_claim"`
	LockedUntil time.Time `json:"locked_until"`
}

// NewStakingPool creates a new staking pool
func NewStakingPool() *StakingPool {
	return &StakingPool{
		stakers:     make(map[string]*Staker),
		totalStaked: 0,
		rewardRate:  0.1, // 10% annual reward rate
		lockPeriod:  time.Hour * 24 * 7, // 7 days lock period
		lastReward:  time.Now(),
	}
}

// Stake tokens
func (sp *StakingPool) Stake(address string, amount uint64) error {
	if amount == 0 {
		return fmt.Errorf("stake amount must be greater than 0")
	}

	staker, exists := sp.stakers[address]
	if !exists {
		staker = &Staker{
			Address:     address,
			Amount:      0,
			Rewards:     0,
			StakedAt:    time.Now(),
			LastClaim:   time.Now(),
			LockedUntil: time.Now().Add(sp.lockPeriod),
		}
		sp.stakers[address] = staker
	}

	staker.Amount += amount
	sp.totalStaked += amount

	return nil
}

// Unstake tokens
func (sp *StakingPool) Unstake(address string, amount uint64) error {
	staker, exists := sp.stakers[address]
	if !exists {
		return fmt.Errorf("no staking position found for %s", address)
	}

	if time.Now().Before(staker.LockedUntil) {
		return fmt.Errorf("tokens are still locked until %s", staker.LockedUntil)
	}

	if amount > staker.Amount {
		return fmt.Errorf("insufficient staked amount")
	}

	staker.Amount -= amount
	sp.totalStaked -= amount

	if staker.Amount == 0 {
		delete(sp.stakers, address)
	}

	return nil
}

// Claim rewards
func (sp *StakingPool) ClaimRewards(address string) (uint64, error) {
	staker, exists := sp.stakers[address]
	if !exists {
		return 0, fmt.Errorf("no staking position found for %s", address)
	}

	// Calculate rewards
	rewards := sp.calculateRewards(staker)
	if rewards == 0 {
		return 0, fmt.Errorf("no rewards to claim")
	}

	staker.Rewards = 0
	staker.LastClaim = time.Now()

	return rewards, nil
}

// GetStakingInfo returns staking information for an address
func (sp *StakingPool) GetStakingInfo(address string) (*Staker, error) {
	staker, exists := sp.stakers[address]
	if !exists {
		return nil, fmt.Errorf("no staking position found for %s", address)
	}

	// Calculate pending rewards
	pendingRewards := sp.calculateRewards(staker)
	staker.Rewards = pendingRewards

	return staker, nil
}

// GetTotalStaked returns total staked amount
func (sp *StakingPool) GetTotalStaked() uint64 {
	return sp.totalStaked
}

// GetStakers returns all stakers
func (sp *StakingPool) GetStakers() map[string]*Staker {
	return sp.stakers
}

// calculateRewards calculates staking rewards
func (sp *StakingPool) calculateRewards(staker *Staker) uint64 {
	daysSinceClaim := time.Since(staker.LastClaim).Hours() / 24
	return uint64(float64(staker.Amount) * sp.rewardRate * daysSinceClaim / 365)
}

// GovernanceSystem manages on-chain governance
type GovernanceSystem struct {
	proposals    map[string]*Proposal
	votes        map[string]*Vote
	treasury     *TreasuryManager
	parameters   *GovernanceParameters
}

// Proposal represents a governance proposal
type Proposal struct {
	ID          string                 `json:"id"`
	Proposer    string                 `json:"proposer"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Actions     []GovernanceAction     `json:"actions"`
	State       string                 `json:"state"` // "active", "passed", "failed", "executed"
	VotesFor    uint64                 `json:"votes_for"`
	VotesAgainst uint64                `json:"votes_against"`
	StartBlock  int64                  `json:"start_block"`
	EndBlock    int64                  `json:"end_block"`
	ExecutedAt  *time.Time             `json:"executed_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// GovernanceAction represents an action to be executed if proposal passes
type GovernanceAction struct {
	Type    string                 `json:"type"`
	Target  string                 `json:"target"`
	Data    map[string]interface{} `json:"data"`
}

// Vote represents a governance vote
type Vote struct {
	ID         string    `json:"id"`
	ProposalID string    `json:"proposal_id"`
	Voter      string    `json:"voter"`
	Choice     string    `json:"choice"` // "for", "against", "abstain"
	Weight     uint64    `json:"weight"`
	Timestamp  time.Time `json:"timestamp"`
}

// GovernanceParameters defines governance parameters
type GovernanceParameters struct {
	MinProposalStake    uint64  `json:"min_proposal_stake"`
	MinVotingStake      uint64  `json:"min_voting_stake"`
	VotingPeriod        int64   `json:"voting_period"`
	QuorumThreshold     float64 `json:"quorum_threshold"`
	PassThreshold       float64 `json:"pass_threshold"`
	ExecutionDelay      int64   `json:"execution_delay"`
}

// NewGovernanceSystem creates a new governance system
func NewGovernanceSystem() *GovernanceSystem {
	return &GovernanceSystem{
		proposals:  make(map[string]*Proposal),
		votes:      make(map[string]*Vote),
		treasury:   NewTreasuryManager(),
		parameters: &GovernanceParameters{
			MinProposalStake: 1000,
			MinVotingStake:   100,
			VotingPeriod:     1000, // blocks
			QuorumThreshold:  0.1,  // 10%
			PassThreshold:    0.6,  // 60%
			ExecutionDelay:   100,  // blocks
		},
	}
}

// CreateProposal creates a new governance proposal
func (gs *GovernanceSystem) CreateProposal(proposer, title, description string, actions []GovernanceAction, currentBlock int64) (*Proposal, error) {
	proposalID := generateProposalID(proposer, title)
	
	proposal := &Proposal{
		ID:          proposalID,
		Proposer:    proposer,
		Title:       title,
		Description: description,
		Actions:     actions,
		State:       "active",
		VotesFor:    0,
		VotesAgainst: 0,
		StartBlock:  currentBlock,
		EndBlock:    currentBlock + gs.parameters.VotingPeriod,
		CreatedAt:   time.Now(),
		Metadata:    make(map[string]interface{}),
	}

	gs.proposals[proposalID] = proposal
	return proposal, nil
}

// Vote on a proposal
func (gs *GovernanceSystem) Vote(proposalID, voter, choice string, weight uint64) error {
	proposal, exists := gs.proposals[proposalID]
	if !exists {
		return fmt.Errorf("proposal %s not found", proposalID)
	}

	if proposal.State != "active" {
		return fmt.Errorf("proposal %s is not active", proposalID)
	}

	if choice != "for" && choice != "against" && choice != "abstain" {
		return fmt.Errorf("invalid vote choice: %s", choice)
	}

	// Check if user already voted
	voteID := generateVoteID(proposalID, voter)
	if _, exists := gs.votes[voteID]; exists {
		return fmt.Errorf("user %s already voted on proposal %s", voter, proposalID)
	}

	vote := &Vote{
		ID:         voteID,
		ProposalID: proposalID,
		Voter:      voter,
		Choice:     choice,
		Weight:     weight,
		Timestamp:  time.Now(),
	}

	gs.votes[voteID] = vote

	// Update proposal vote counts
	if choice == "for" {
		proposal.VotesFor += weight
	} else if choice == "against" {
		proposal.VotesAgainst += weight
	}

	return nil
}

// GetProposal returns a proposal by ID
func (gs *GovernanceSystem) GetProposal(proposalID string) (*Proposal, error) {
	proposal, exists := gs.proposals[proposalID]
	if !exists {
		return nil, fmt.Errorf("proposal %s not found", proposalID)
	}
	return proposal, nil
}

// GetActiveProposals returns all active proposals
func (gs *GovernanceSystem) GetActiveProposals() []*Proposal {
	var activeProposals []*Proposal
	for _, proposal := range gs.proposals {
		if proposal.State == "active" {
			activeProposals = append(activeProposals, proposal)
		}
	}
	return activeProposals
}

// Execute proposal if it passed
func (gs *GovernanceSystem) ExecuteProposal(proposalID string, currentBlock int64) error {
	proposal, exists := gs.proposals[proposalID]
	if !exists {
		return fmt.Errorf("proposal %s not found", proposalID)
	}

	if proposal.State != "passed" {
		return fmt.Errorf("proposal %s is not in passed state", proposalID)
	}

	if currentBlock < proposal.EndBlock+gs.parameters.ExecutionDelay {
		return fmt.Errorf("execution delay not met for proposal %s", proposalID)
	}

	// Execute actions
	for _, action := range proposal.Actions {
		if err := gs.executeAction(action); err != nil {
			return fmt.Errorf("failed to execute action: %v", err)
		}
	}

	// Update proposal state
	proposal.State = "executed"
	now := time.Now()
	proposal.ExecutedAt = &now

	return nil
}

// executeAction executes a governance action
func (gs *GovernanceSystem) executeAction(action GovernanceAction) error {
	switch action.Type {
	case "parameter_change":
		return gs.executeParameterChange(action)
	case "treasury_transfer":
		return gs.executeTreasuryTransfer(action)
	case "contract_upgrade":
		return gs.executeContractUpgrade(action)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// executeParameterChange executes a parameter change action
func (gs *GovernanceSystem) executeParameterChange(action GovernanceAction) error {
	// Implementation for parameter changes
	return nil
}

// executeTreasuryTransfer executes a treasury transfer action
func (gs *GovernanceSystem) executeTreasuryTransfer(action GovernanceAction) error {
	// Implementation for treasury transfers
	return nil
}

// executeContractUpgrade executes a contract upgrade action
func (gs *GovernanceSystem) executeContractUpgrade(action GovernanceAction) error {
	// Implementation for contract upgrades
	return nil
}

// TreasuryManager manages protocol treasury
type TreasuryManager struct {
	balance     uint64
	transactions []*TreasuryTransaction
}

// TreasuryTransaction represents a treasury transaction
type TreasuryTransaction struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Amount    uint64    `json:"amount"`
	Recipient string    `json:"recipient"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// NewTreasuryManager creates a new treasury manager
func NewTreasuryManager() *TreasuryManager {
	return &TreasuryManager{
		balance:     0,
		transactions: make([]*TreasuryTransaction, 0),
	}
}

// GetBalance returns treasury balance
func (tm *TreasuryManager) GetBalance() uint64 {
	return tm.balance
}

// AddFunds adds funds to treasury
func (tm *TreasuryManager) AddFunds(amount uint64, reason string) {
	tm.balance += amount
	tm.transactions = append(tm.transactions, &TreasuryTransaction{
		ID:        generateTreasuryTxID(),
		Type:      "deposit",
		Amount:    amount,
		Recipient: "treasury",
		Reason:    reason,
		Timestamp: time.Now(),
	})
}

// TransferFunds transfers funds from treasury
func (tm *TreasuryManager) TransferFunds(amount uint64, recipient, reason string) error {
	if amount > tm.balance {
		return fmt.Errorf("insufficient treasury balance")
	}

	tm.balance -= amount
	tm.transactions = append(tm.transactions, &TreasuryTransaction{
		ID:        generateTreasuryTxID(),
		Type:      "transfer",
		Amount:    amount,
		Recipient: recipient,
		Reason:    reason,
		Timestamp: time.Now(),
	})

	return nil
}

// GetTransactions returns treasury transactions
func (tm *TreasuryManager) GetTransactions() []*TreasuryTransaction {
	return tm.transactions
}

// Helper functions
func generateProposalID(proposer, title string) string {
	data := fmt.Sprintf("%s_%s_%d", proposer, title, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func generateVoteID(proposalID, voter string) string {
	data := fmt.Sprintf("%s_%s_%d", proposalID, voter, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func generateTreasuryTxID() string {
	data := fmt.Sprintf("treasury_%d", time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
} 