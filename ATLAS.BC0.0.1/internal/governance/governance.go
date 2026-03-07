package governance

import (
	"atlas-blockchain/internal/defi"
	"atlas-blockchain/internal/identity"
	"atlas-blockchain/internal/social"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// GovernanceManager handles all governance operations
type GovernanceManager struct {
	proposals   map[string]*GovernanceProposal
	votes       map[string]*GovernanceVote
	committees  map[string]*Committee
	referendums map[string]*Referendum
	parameters  *GovernanceParameters
	social      *social.SocialManager
	defi        *defi.DeFiManager
	identity    *identity.IdentityManager
}

// GovernanceProposal represents a governance proposal
type GovernanceProposal struct {
	ID            string                 `json:"id"`
	Proposer      string                 `json:"proposer"`
	Title         string                 `json:"title"`
	Description   string                 `json:"description"`
	Category      string                 `json:"category"` // "platform", "defi", "social", "technical", "economic"
	Actions       []GovernanceAction     `json:"actions"`
	State         string                 `json:"state"` // "draft", "active", "passed", "failed", "executed", "cancelled"
	VotesFor      uint64                 `json:"votes_for"`
	VotesAgainst  uint64                 `json:"votes_against"`
	VotesAbstain  uint64                 `json:"votes_abstain"`
	TotalVotes    uint64                 `json:"total_votes"`
	QuorumMet     bool                   `json:"quorum_met"`
	StartBlock    int64                  `json:"start_block"`
	EndBlock      int64                  `json:"end_block"`
	ExecutedAt    *time.Time             `json:"executed_at,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	Metadata      map[string]interface{} `json:"metadata"`
	SocialMetrics *SocialMetrics         `json:"social_metrics,omitempty"`
	Discussion    []*DiscussionComment   `json:"discussion,omitempty"`
}

// GovernanceAction represents an action to be executed if proposal passes
type GovernanceAction struct {
	Type        string                 `json:"type"`
	Target      string                 `json:"target"`
	Data        map[string]interface{} `json:"data"`
	Description string                 `json:"description"`
	Impact      string                 `json:"impact"` // "low", "medium", "high", "critical"
}

// GovernanceVote represents a governance vote
type GovernanceVote struct {
	ID         string    `json:"id"`
	ProposalID string    `json:"proposal_id"`
	Voter      string    `json:"voter"`
	Choice     string    `json:"choice"` // "for", "against", "abstain"
	Weight     uint64    `json:"weight"`
	Reason     string    `json:"reason,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
	Delegated  bool      `json:"delegated"`
	Delegate   string    `json:"delegate,omitempty"`
}

// Committee represents a governance committee
type Committee struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Members     []string  `json:"members"`
	Chair       string    `json:"chair"`
	Type        string    `json:"type"` // "technical", "economic", "social", "security"
	CreatedAt   time.Time `json:"created_at"`
	IsActive    bool      `json:"is_active"`
}

// Referendum represents a community referendum
type Referendum struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Question    string            `json:"question"`
	Options     []string          `json:"options"`
	State       string            `json:"state"` // "active", "closed", "passed", "failed"
	Votes       map[string]uint64 `json:"votes"`
	TotalVotes  uint64            `json:"total_votes"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
	CreatedAt   time.Time         `json:"created_at"`
}

// SocialMetrics tracks social engagement with proposals
type SocialMetrics struct {
	DiscussionPosts int64     `json:"discussion_posts"`
	Comments        int64     `json:"comments"`
	Likes           int64     `json:"likes"`
	Shares          int64     `json:"shares"`
	Views           int64     `json:"views"`
	Sentiment       string    `json:"sentiment"` // "positive", "neutral", "negative"
	LastUpdated     time.Time `json:"last_updated"`
}

// DiscussionComment represents a comment in proposal discussion
type DiscussionComment struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	ParentID  string    `json:"parent_id,omitempty"`
	Likes     int64     `json:"likes"`
	Replies   int64     `json:"replies"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GovernanceParameters defines governance parameters
type GovernanceParameters struct {
	MinProposalStake    uint64  `json:"min_proposal_stake"`
	MinVotingStake      uint64  `json:"min_voting_stake"`
	VotingPeriod        int64   `json:"voting_period"`
	QuorumThreshold     float64 `json:"quorum_threshold"`
	PassThreshold       float64 `json:"pass_threshold"`
	ExecutionDelay      int64   `json:"execution_delay"`
	MaxProposalLength   int     `json:"max_proposal_length"`
	MinDiscussionPeriod int64   `json:"min_discussion_period"`
	DelegationEnabled   bool    `json:"delegation_enabled"`
	CommitteeRequired   bool    `json:"committee_required"`
}

// NewGovernanceManager creates a new governance manager
func NewGovernanceManager(socialManager *social.SocialManager, defiManager *defi.DeFiManager, identityManager *identity.IdentityManager) *GovernanceManager {
	return &GovernanceManager{
		proposals:   make(map[string]*GovernanceProposal),
		votes:       make(map[string]*GovernanceVote),
		committees:  make(map[string]*Committee),
		referendums: make(map[string]*Referendum),
		parameters: &GovernanceParameters{
			MinProposalStake:    1000,
			MinVotingStake:      100,
			VotingPeriod:        1000, // blocks
			QuorumThreshold:     0.1,  // 10%
			PassThreshold:       0.6,  // 60%
			ExecutionDelay:      100,  // blocks
			MaxProposalLength:   5000,
			MinDiscussionPeriod: 100, // blocks
			DelegationEnabled:   true,
			CommitteeRequired:   false,
		},
		social:   socialManager,
		defi:     defiManager,
		identity: identityManager,
	}
}

// CreateProposal creates a new governance proposal
func (gm *GovernanceManager) CreateProposal(proposer, title, description, category string, actions []GovernanceAction, currentBlock int64) (*GovernanceProposal, error) {
	// Validate proposal length
	if len(description) > gm.parameters.MaxProposalLength {
		return nil, fmt.Errorf("proposal description exceeds maximum length of %d characters", gm.parameters.MaxProposalLength)
	}

	// Check if proposer has minimum stake
	if gm.identity != nil {
		identity, err := gm.identity.GetIdentity(proposer)
		if err != nil {
			return nil, fmt.Errorf("proposer identity not found")
		}

		// Check if user has enough activity/reputation to propose
		if identity.Activity.TotalTokensEarned < int64(gm.parameters.MinProposalStake) {
			return nil, fmt.Errorf("insufficient activity to create proposal")
		}
	}

	proposalID := generateProposalID(proposer, title)

	proposal := &GovernanceProposal{
		ID:           proposalID,
		Proposer:     proposer,
		Title:        title,
		Description:  description,
		Category:     category,
		Actions:      actions,
		State:        "draft",
		VotesFor:     0,
		VotesAgainst: 0,
		VotesAbstain: 0,
		TotalVotes:   0,
		QuorumMet:    false,
		StartBlock:   currentBlock + gm.parameters.MinDiscussionPeriod,
		EndBlock:     currentBlock + gm.parameters.MinDiscussionPeriod + gm.parameters.VotingPeriod,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Metadata:     make(map[string]interface{}),
		SocialMetrics: &SocialMetrics{
			DiscussionPosts: 0,
			Comments:        0,
			Likes:           0,
			Shares:          0,
			Views:           0,
			Sentiment:       "neutral",
			LastUpdated:     time.Now(),
		},
		Discussion: make([]*DiscussionComment, 0),
	}

	gm.proposals[proposalID] = proposal

	// Create social media post for discussion
	if gm.social != nil {
		socialContent := fmt.Sprintf("New Governance Proposal: %s\n\n%s\n\nCategory: %s\n\n#governance #proposal",
			title, description, category)

		_, err := gm.social.CreatePost(proposer, socialContent, []string{}, "public", "governance")
		if err != nil {
			// Log error but don't fail proposal creation
			fmt.Printf("Failed to create social post for proposal: %v\n", err)
		}
	}

	// Update user activity
	if gm.identity != nil {
		gm.identity.UpdateActivity(proposer, "proposal", 1)
	}

	return proposal, nil
}

// ActivateProposal activates a proposal for voting
func (gm *GovernanceManager) ActivateProposal(proposalID string, currentBlock int64) error {
	proposal, exists := gm.proposals[proposalID]
	if !exists {
		return fmt.Errorf("proposal %s not found", proposalID)
	}

	if proposal.State != "draft" {
		return fmt.Errorf("proposal %s is not in draft state", proposalID)
	}

	if currentBlock < proposal.StartBlock {
		return fmt.Errorf("proposal %s cannot be activated before start block", proposalID)
	}

	proposal.State = "active"
	proposal.UpdatedAt = time.Now()

	return nil
}

// Vote on a proposal
func (gm *GovernanceManager) Vote(proposalID, voter, choice, reason string, weight uint64) error {
	proposal, exists := gm.proposals[proposalID]
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
	if _, exists := gm.votes[voteID]; exists {
		return fmt.Errorf("user %s already voted on proposal %s", voter, proposalID)
	}

	// Validate voting power
	if gm.identity != nil {
		identity, err := gm.identity.GetIdentity(voter)
		if err != nil {
			return fmt.Errorf("voter identity not found")
		}

		// Calculate voting power based on activity and reputation
		votingPower := gm.calculateVotingPower(identity)
		if votingPower < float64(gm.parameters.MinVotingStake) {
			return fmt.Errorf("insufficient voting power")
		}

		weight = uint64(votingPower)
	}

	vote := &GovernanceVote{
		ID:         voteID,
		ProposalID: proposalID,
		Voter:      voter,
		Choice:     choice,
		Weight:     weight,
		Reason:     reason,
		Timestamp:  time.Now(),
		Delegated:  false,
	}

	gm.votes[voteID] = vote

	// Update proposal vote counts
	switch choice {
	case "for":
		proposal.VotesFor += weight
	case "against":
		proposal.VotesAgainst += weight
	case "abstain":
		proposal.VotesAbstain += weight
	}

	proposal.TotalVotes += weight
	proposal.UpdatedAt = time.Now()

	// Check quorum
	gm.checkQuorum(proposal)

	// Update user activity
	if gm.identity != nil {
		gm.identity.UpdateActivity(voter, "vote", 1)
	}

	return nil
}

// DelegateVote delegates voting power to another user
func (gm *GovernanceManager) DelegateVote(delegator, delegate string) error {
	if !gm.parameters.DelegationEnabled {
		return fmt.Errorf("vote delegation is not enabled")
	}

	// Validate both users exist
	if gm.identity != nil {
		if _, err := gm.identity.GetIdentity(delegator); err != nil {
			return fmt.Errorf("delegator identity not found")
		}
		if _, err := gm.identity.GetIdentity(delegate); err != nil {
			return fmt.Errorf("delegate identity not found")
		}
	}

	// Store delegation (simplified implementation)
	// In a full implementation, this would be stored in a separate delegation table

	return nil
}

// AddDiscussionComment adds a comment to proposal discussion
func (gm *GovernanceManager) AddDiscussionComment(proposalID, author, content, parentID string) error {
	proposal, exists := gm.proposals[proposalID]
	if !exists {
		return fmt.Errorf("proposal %s not found", proposalID)
	}

	if proposal.State == "executed" || proposal.State == "cancelled" {
		return fmt.Errorf("cannot comment on %s proposal", proposal.State)
	}

	commentID := generateCommentID(author, proposalID)
	comment := &DiscussionComment{
		ID:        commentID,
		Author:    author,
		Content:   content,
		ParentID:  parentID,
		Likes:     0,
		Replies:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	proposal.Discussion = append(proposal.Discussion, comment)
	proposal.SocialMetrics.Comments++
	proposal.SocialMetrics.LastUpdated = time.Now()

	// Update proposal sentiment based on discussion
	gm.updateProposalSentiment(proposal)

	return nil
}

// GetProposal returns a proposal by ID
func (gm *GovernanceManager) GetProposal(proposalID string) (*GovernanceProposal, error) {
	proposal, exists := gm.proposals[proposalID]
	if !exists {
		return nil, fmt.Errorf("proposal %s not found", proposalID)
	}
	return proposal, nil
}

// GetActiveProposals returns all active proposals
func (gm *GovernanceManager) GetActiveProposals() []*GovernanceProposal {
	var activeProposals []*GovernanceProposal
	for _, proposal := range gm.proposals {
		if proposal.State == "active" {
			activeProposals = append(activeProposals, proposal)
		}
	}
	return activeProposals
}

// GetProposalsByCategory returns proposals filtered by category
func (gm *GovernanceManager) GetProposalsByCategory(category string) []*GovernanceProposal {
	var proposals []*GovernanceProposal
	for _, proposal := range gm.proposals {
		if proposal.Category == category {
			proposals = append(proposals, proposal)
		}
	}
	return proposals
}

// Execute proposal if it passed
func (gm *GovernanceManager) ExecuteProposal(proposalID string, currentBlock int64) error {
	proposal, exists := gm.proposals[proposalID]
	if !exists {
		return fmt.Errorf("proposal %s not found", proposalID)
	}

	if proposal.State != "passed" {
		return fmt.Errorf("proposal %s is not in passed state", proposalID)
	}

	if currentBlock < proposal.EndBlock+gm.parameters.ExecutionDelay {
		return fmt.Errorf("execution delay not met for proposal %s", proposalID)
	}

	// Execute actions
	for _, action := range proposal.Actions {
		if err := gm.executeAction(action); err != nil {
			return fmt.Errorf("failed to execute action: %v", err)
		}
	}

	// Update proposal state
	proposal.State = "executed"
	now := time.Now()
	proposal.ExecutedAt = &now
	proposal.UpdatedAt = now

	return nil
}

// CreateReferendum creates a new community referendum
func (gm *GovernanceManager) CreateReferendum(title, description, question string, options []string, duration time.Duration) (*Referendum, error) {
	if len(options) < 2 {
		return nil, fmt.Errorf("referendum must have at least 2 options")
	}

	referendumID := generateReferendumID(title)
	referendum := &Referendum{
		ID:          referendumID,
		Title:       title,
		Description: description,
		Question:    question,
		Options:     options,
		State:       "active",
		Votes:       make(map[string]uint64),
		TotalVotes:  0,
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(duration),
		CreatedAt:   time.Now(),
	}

	gm.referendums[referendumID] = referendum

	return referendum, nil
}

// VoteReferendum votes on a referendum
func (gm *GovernanceManager) VoteReferendum(referendumID, voter, option string) error {
	referendum, exists := gm.referendums[referendumID]
	if !exists {
		return fmt.Errorf("referendum %s not found", referendumID)
	}

	if referendum.State != "active" {
		return fmt.Errorf("referendum %s is not active", referendumID)
	}

	if time.Now().After(referendum.EndTime) {
		return fmt.Errorf("referendum %s has ended", referendumID)
	}

	// Validate option
	validOption := false
	for _, opt := range referendum.Options {
		if opt == option {
			validOption = true
			break
		}
	}
	if !validOption {
		return fmt.Errorf("invalid option: %s", option)
	}

	// Calculate voting weight
	weight := uint64(1) // Default weight
	if gm.identity != nil {
		if identity, err := gm.identity.GetIdentity(voter); err == nil {
			weight = uint64(gm.calculateVotingPower(identity))
		}
	}

	// Add vote
	referendum.Votes[option] += weight
	referendum.TotalVotes += weight

	return nil
}

// Helper functions
func (gm *GovernanceManager) calculateVotingPower(identity *identity.UserIdentity) float64 {
	// Base voting power from activity
	power := float64(identity.Activity.TotalTokensEarned) * 0.1

	// Bonus from reputation
	power += identity.Reputation.Overall * 10

	// Bonus from governance activity
	power += float64(identity.Activity.ProposalsCreated) * 100
	power += float64(identity.Activity.VotesCast) * 10

	// Bonus from social activity
	power += float64(identity.Activity.PostsCreated) * 5
	power += float64(identity.Activity.CommentsMade) * 2

	return power
}

func (gm *GovernanceManager) checkQuorum(proposal *GovernanceProposal) {
	totalStake := gm.getTotalStake()
	if totalStake > 0 {
		quorumPercentage := float64(proposal.TotalVotes) / float64(totalStake)
		proposal.QuorumMet = quorumPercentage >= gm.parameters.QuorumThreshold
	}
}

func (gm *GovernanceManager) getTotalStake() uint64 {
	// In a real implementation, this would query the total staked tokens
	// For now, return a fixed value
	return 1000000
}

func (gm *GovernanceManager) updateProposalSentiment(proposal *GovernanceProposal) {
	// Simple sentiment analysis based on discussion activity
	// In production, this would use NLP/AI
	if proposal.SocialMetrics.Comments > 100 {
		proposal.SocialMetrics.Sentiment = "positive"
	} else if proposal.SocialMetrics.Comments > 50 {
		proposal.SocialMetrics.Sentiment = "neutral"
	} else {
		proposal.SocialMetrics.Sentiment = "negative"
	}
}

func (gm *GovernanceManager) executeAction(action GovernanceAction) error {
	switch action.Type {
	case "parameter_change":
		return gm.executeParameterChange(action)
	case "defi_parameter":
		return gm.executeDeFiParameterChange(action)
	case "social_parameter":
		return gm.executeSocialParameterChange(action)
	case "treasury_transfer":
		return gm.executeTreasuryTransfer(action)
	case "contract_upgrade":
		return gm.executeContractUpgrade(action)
	case "committee_creation":
		return gm.executeCommitteeCreation(action)
	default:
		return fmt.Errorf("unknown action type: %s", action.Type)
	}
}

func (gm *GovernanceManager) executeParameterChange(action GovernanceAction) error {
	// Implementation for parameter changes
	params := action.Data

	if val, ok := params["min_proposal_stake"]; ok {
		gm.parameters.MinProposalStake = uint64(val.(float64))
	}
	if val, ok := params["min_voting_stake"]; ok {
		gm.parameters.MinVotingStake = uint64(val.(float64))
	}
	if val, ok := params["voting_period"]; ok {
		gm.parameters.VotingPeriod = int64(val.(float64))
	}
	if val, ok := params["quorum_threshold"]; ok {
		gm.parameters.QuorumThreshold = val.(float64)
	}
	if val, ok := params["pass_threshold"]; ok {
		gm.parameters.PassThreshold = val.(float64)
	}
	if val, ok := params["execution_delay"]; ok {
		gm.parameters.ExecutionDelay = int64(val.(float64))
	}

	fmt.Printf("✅ Governance parameters updated: %+v\n", gm.parameters)
	return nil
}

func (gm *GovernanceManager) executeDeFiParameterChange(action GovernanceAction) error {
	// Implementation for DeFi parameter changes
	return nil
}

func (gm *GovernanceManager) executeSocialParameterChange(action GovernanceAction) error {
	// Implementation for social parameter changes
	return nil
}

func (gm *GovernanceManager) executeTreasuryTransfer(action GovernanceAction) error {
	// Implementation for treasury transfers
	// In a real implementation, this would trigger a transaction from the treasury account
	// checking if the treasury has enough funds via StateManager (not currently available here)

	target := action.Target
	amount, ok := action.Data["amount"].(float64)
	if !ok {
		return fmt.Errorf("invalid amount for treasury transfer")
	}

	fmt.Printf("💰 Executing Treasury Transfer: %f DUT to %s\n", amount, target)
	// TODO: Inject TransactionManager/StateManager to actually move funds

	return nil
}

func (gm *GovernanceManager) executeContractUpgrade(action GovernanceAction) error {
	// Implementation for contract upgrades
	return nil
}

func (gm *GovernanceManager) executeCommitteeCreation(action GovernanceAction) error {
	// Implementation for committee creation
	return nil
}

// Helper ID generation functions
func generateProposalID(proposer, title string) string {
	data := fmt.Sprintf("proposal_%s_%s_%d", proposer, title, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func generateVoteID(proposalID, voter string) string {
	data := fmt.Sprintf("vote_%s_%s_%d", proposalID, voter, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func generateCommentID(author, proposalID string) string {
	data := fmt.Sprintf("comment_%s_%s_%d", author, proposalID, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func generateReferendumID(title string) string {
	data := fmt.Sprintf("referendum_%s_%d", title, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}
