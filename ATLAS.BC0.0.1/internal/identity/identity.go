package identity

import (
	"atlas-blockchain/pkg/crypto/zk"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// UserIdentity represents a user's identity in the platform
type UserIdentity struct {
	Address         string                 `json:"address"`
	PublicKey       []byte                 `json:"public_key"`
	Username        string                 `json:"username"`
	Email           string                 `json:"email"`
	Profile         *UserProfile           `json:"profile"`
	Credentials     []*Credential          `json:"credentials"`
	Attestations    []*Attestation         `json:"attestations"`
	PrivacySettings *PrivacySettings       `json:"privacy_settings"`
	KYC             *KYCInfo               `json:"kyc"`
	Reputation      *ReputationScore       `json:"reputation"`
	Activity        *ActivityMetrics       `json:"activity"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	IsActive        bool                   `json:"is_active"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// UserProfile contains public profile information
type UserProfile struct {
	DisplayName string            `json:"display_name"`
	Bio         string            `json:"bio"`
	Avatar      string            `json:"avatar"`
	Location    string            `json:"location"`
	Website     string            `json:"website"`
	SocialLinks map[string]string `json:"social_links"`
	Interests   []string          `json:"interests"`
	Skills      []string          `json:"skills"`
	IsPublic    bool              `json:"is_public"`
}

// Credential represents a verified credential
type Credential struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"` // email, phone, social, kyc, etc.
	Issuer    string                 `json:"issuer"`
	Subject   string                 `json:"subject"`
	Data      map[string]interface{} `json:"data"`
	Proof     *zk.ZKProof            `json:"proof"`
	IssuedAt  time.Time              `json:"issued_at"`
	ExpiresAt *time.Time             `json:"expires_at"`
	IsRevoked bool                   `json:"is_revoked"`
	RevokedAt *time.Time             `json:"revoked_at"`
}

// Attestation represents a third-party attestation
type Attestation struct {
	ID          string                 `json:"id"`
	Attester    string                 `json:"attester"`
	Subject     string                 `json:"subject"`
	Type        string                 `json:"type"`
	Data        map[string]interface{} `json:"data"`
	Proof       *zk.ZKProof            `json:"proof"`
	AttestedAt  time.Time              `json:"attested_at"`
	IsValid     bool                   `json:"is_valid"`
	ValidatedAt *time.Time             `json:"validated_at"`
}

// PrivacySettings controls what information is visible
type PrivacySettings struct {
	ProfileVisibility  string            `json:"profile_visibility"`  // public, friends, private
	ActivityVisibility string            `json:"activity_visibility"` // public, friends, private
	TransactionPrivacy string            `json:"transaction_privacy"` // public, private, zk
	ContactVisibility  string            `json:"contact_visibility"`  // public, friends, private
	LocationSharing    bool              `json:"location_sharing"`
	DataRetention      string            `json:"data_retention"` // 30days, 1year, forever
	GDPRConsent        bool              `json:"gdpr_consent"`
	MarketingConsent   bool              `json:"marketing_consent"`
	CustomSettings     map[string]string `json:"custom_settings"`
}

// KYCInfo contains Know Your Customer information
type KYCInfo struct {
	FullName    string     `json:"full_name"`
	DateOfBirth string     `json:"date_of_birth"`
	Nationality string     `json:"nationality"`
	IDType      string     `json:"id_type"`
	IDNumber    string     `json:"id_number"`
	Address     string     `json:"address"`
	Country     string     `json:"country"`
	Verified    bool       `json:"verified"`
	VerifiedAt  *time.Time `json:"verified_at"`
	Verifier    string     `json:"verifier"`
	Level       string     `json:"level"` // basic, enhanced, full
	Documents   []string   `json:"documents"`
}

// ReputationScore tracks user reputation
type ReputationScore struct {
	Overall     float64            `json:"overall"`
	Social      float64            `json:"social"`
	Commerce    float64            `json:"commerce"`
	Governance  float64            `json:"governance"`
	Factors     map[string]float64 `json:"factors"`
	LastUpdated time.Time          `json:"last_updated"`
	History     []ReputationEvent  `json:"history"`
}

// ReputationEvent represents a reputation change
type ReputationEvent struct {
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

// ActivityMetrics tracks user activity
type ActivityMetrics struct {
	PostsCreated      int64     `json:"posts_created"`
	CommentsMade      int64     `json:"comments_made"`
	LikesGiven        int64     `json:"likes_given"`
	LikesReceived     int64     `json:"likes_received"`
	Transactions      int64     `json:"transactions"`
	ProposalsCreated  int64     `json:"proposals_created"`
	VotesCast         int64     `json:"votes_cast"`
	LastActive        time.Time `json:"last_active"`
	StreakDays        int       `json:"streak_days"`
	TotalTokensEarned int64     `json:"total_tokens_earned"`
}

// IdentityManager handles user identity operations
type IdentityManager struct {
	identities   map[string]*UserIdentity
	credentials  map[string]*Credential
	attestations map[string]*Attestation
	zkVerifier   *zk.ProofVerifier
	reputation   *ReputationManager
	privacy      *PrivacyManager
	kyc          *KYCManager
}

// NewIdentityManager creates a new identity manager
func NewIdentityManager() *IdentityManager {
	return &IdentityManager{
		identities:   make(map[string]*UserIdentity),
		credentials:  make(map[string]*Credential),
		attestations: make(map[string]*Attestation),
		zkVerifier:   zk.NewProofVerifier(true),
		reputation:   NewReputationManager(),
		privacy:      NewPrivacyManager(),
		kyc:          NewKYCManager(),
	}
}

// CreateIdentity creates a new user identity
func (im *IdentityManager) CreateIdentity(address string, publicKey []byte, username string, email string) (*UserIdentity, error) {
	// Validate inputs
	if address == "" || len(publicKey) == 0 || username == "" || email == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	// Check if identity already exists
	if _, exists := im.identities[address]; exists {
		return nil, fmt.Errorf("identity already exists for address: %s", address)
	}

	// Create new identity
	identity := &UserIdentity{
		Address:         address,
		PublicKey:       publicKey,
		Username:        username,
		Email:           email,
		Profile:         &UserProfile{DisplayName: username, IsPublic: true},
		Credentials:     make([]*Credential, 0),
		Attestations:    make([]*Attestation, 0),
		PrivacySettings: im.privacy.GetDefaultSettings(),
		KYC:             &KYCInfo{Level: "basic"},
		Reputation:      &ReputationScore{Overall: 0.0, LastUpdated: time.Now()},
		Activity:        &ActivityMetrics{LastActive: time.Now()},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		IsActive:        true,
		Metadata:        make(map[string]interface{}),
	}

	// Store identity
	im.identities[address] = identity

	// Create initial credentials
	im.createInitialCredentials(identity)

	return identity, nil
}

// GetIdentity retrieves a user identity
func (im *IdentityManager) GetIdentity(address string) (*UserIdentity, error) {
	identity, exists := im.identities[address]
	if !exists {
		return nil, fmt.Errorf("identity not found: %s", address)
	}
	return identity, nil
}

// UpdateProfile updates user profile information
func (im *IdentityManager) UpdateProfile(address string, profile *UserProfile) error {
	identity, err := im.GetIdentity(address)
	if err != nil {
		return err
	}

	identity.Profile = profile
	identity.UpdatedAt = time.Now()
	return nil
}

// AddCredential adds a verified credential to an identity
func (im *IdentityManager) AddCredential(address string, credential *Credential) error {
	identity, err := im.GetIdentity(address)
	if err != nil {
		return err
	}

	// Verify credential proof if provided
	if credential.Proof != nil {
		valid, err := im.zkVerifier.VerifyProof(credential.Proof)
		if err != nil {
			return fmt.Errorf("credential proof verification failed: %v", err)
		}
		if !valid {
			return fmt.Errorf("invalid credential proof")
		}
	}

	// Store credential
	im.credentials[credential.ID] = credential
	identity.Credentials = append(identity.Credentials, credential)
	identity.UpdatedAt = time.Now()

	return nil
}

// AddAttestation adds a third-party attestation
func (im *IdentityManager) AddAttestation(address string, attestation *Attestation) error {
	identity, err := im.GetIdentity(address)
	if err != nil {
		return err
	}

	// Verify attestation proof if provided
	if attestation.Proof != nil {
		valid, err := im.zkVerifier.VerifyProof(attestation.Proof)
		if err != nil {
			return fmt.Errorf("attestation proof verification failed: %v", err)
		}
		if !valid {
			return fmt.Errorf("invalid attestation proof")
		}
	}

	// Store attestation
	im.attestations[attestation.ID] = attestation
	identity.Attestations = append(identity.Attestations, attestation)
	identity.UpdatedAt = time.Now()

	return nil
}

// UpdatePrivacySettings updates user privacy settings
func (im *IdentityManager) UpdatePrivacySettings(address string, settings *PrivacySettings) error {
	identity, err := im.GetIdentity(address)
	if err != nil {
		return err
	}

	identity.PrivacySettings = settings
	identity.UpdatedAt = time.Now()
	return nil
}

// UpdateKYC updates KYC information
func (im *IdentityManager) UpdateKYC(address string, kycInfo *KYCInfo) error {
	identity, err := im.GetIdentity(address)
	if err != nil {
		return err
	}

	identity.KYC = kycInfo
	identity.UpdatedAt = time.Now()
	return nil
}

// UpdateActivity updates user activity metrics
func (im *IdentityManager) UpdateActivity(address string, activityType string, value int64) error {
	identity, err := im.GetIdentity(address)
	if err != nil {
		return err
	}

	// Update activity metrics
	switch activityType {
	case "post":
		identity.Activity.PostsCreated += value
	case "comment":
		identity.Activity.CommentsMade += value
	case "like_given":
		identity.Activity.LikesGiven += value
	case "like_received":
		identity.Activity.LikesReceived += value
	case "transaction":
		identity.Activity.Transactions += value
	case "proposal":
		identity.Activity.ProposalsCreated += value
	case "vote":
		identity.Activity.VotesCast += value
	case "tokens_earned":
		identity.Activity.TotalTokensEarned += value
	}

	identity.Activity.LastActive = time.Now()
	identity.UpdatedAt = time.Now()

	// Update reputation based on activity
	im.reputation.UpdateReputation(identity, activityType, value)

	return nil
}

// CreatePrivacyProof creates a zero-knowledge proof for privacy-preserving operations
func (im *IdentityManager) CreatePrivacyProof(proofType zk.ProofType, data interface{}, address string) (*zk.ZKProof, error) {
	return zk.GenerateProof(proofType, data, address)
}

// VerifyPrivacyProof verifies a zero-knowledge proof
func (im *IdentityManager) VerifyPrivacyProof(proof *zk.ZKProof) (bool, error) {
	return im.zkVerifier.VerifyProof(proof)
}

// GetIdentityForSocial returns identity information suitable for social features
func (im *IdentityManager) GetIdentityForSocial(address string, requesterAddress string) (*UserIdentity, error) {
	identity, err := im.GetIdentity(address)
	if err != nil {
		return nil, err
	}

	// Apply privacy settings
	if identity.PrivacySettings.ProfileVisibility == "private" && address != requesterAddress {
		// Return minimal public information
		return &UserIdentity{
			Address:    identity.Address,
			Username:   identity.Username,
			Profile:    &UserProfile{DisplayName: identity.Profile.DisplayName, IsPublic: false},
			Reputation: identity.Reputation,
			Activity:   &ActivityMetrics{LastActive: identity.Activity.LastActive},
			IsActive:   identity.IsActive,
		}, nil
	}

	return identity, nil
}

// GetIdentityForCommerce returns identity information suitable for commerce
func (im *IdentityManager) GetIdentityForCommerce(address string, requesterAddress string) (*UserIdentity, error) {
	identity, err := im.GetIdentity(address)
	if err != nil {
		return nil, err
	}

	// Apply commerce-specific privacy settings
	if identity.PrivacySettings.TransactionPrivacy == "private" && address != requesterAddress {
		// Return minimal commerce information
		return &UserIdentity{
			Address:    identity.Address,
			Username:   identity.Username,
			Reputation: &ReputationScore{Commerce: identity.Reputation.Commerce},
			KYC:        identity.KYC,
			IsActive:   identity.IsActive,
		}, nil
	}

	return identity, nil
}

// GetIdentityForGovernance returns identity information suitable for governance
func (im *IdentityManager) GetIdentityForGovernance(address string) (*UserIdentity, error) {
	identity, err := im.GetIdentity(address)
	if err != nil {
		return nil, err
	}

	// For governance, we need to verify voting eligibility
	if !im.isEligibleForGovernance(identity) {
		return nil, fmt.Errorf("user not eligible for governance")
	}

	return identity, nil
}

// Helper methods
func (im *IdentityManager) createInitialCredentials(identity *UserIdentity) {
	// Create email credential
	emailCredential := &Credential{
		ID:        generateCredentialID(identity.Address, "email"),
		Type:      "email",
		Issuer:    "system",
		Subject:   identity.Address,
		Data:      map[string]interface{}{"email": identity.Email},
		IssuedAt:  time.Now(),
		IsRevoked: false,
	}

	// Create username credential
	usernameCredential := &Credential{
		ID:        generateCredentialID(identity.Address, "username"),
		Type:      "username",
		Issuer:    "system",
		Subject:   identity.Address,
		Data:      map[string]interface{}{"username": identity.Username},
		IssuedAt:  time.Now(),
		IsRevoked: false,
	}

	identity.Credentials = append(identity.Credentials, emailCredential, usernameCredential)
}

func (im *IdentityManager) isEligibleForGovernance(identity *UserIdentity) bool {
	// Check if user has minimum token balance (this would be checked against blockchain state)
	// Check if user has verified KYC
	// Check if user has good reputation
	return identity.KYC.Verified && identity.Reputation.Overall >= 50.0
}

func generateCredentialID(address, credentialType string) string {
	data := fmt.Sprintf("%s_%s_%d", address, credentialType, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

// ApplyReputationDecay runs entropy on all identities
func (im *IdentityManager) ApplyReputationDecay() {
	for _, identity := range im.identities {
		im.reputation.DecayReputation(identity)
	}
}

// ReputationManager handles reputation calculations
type ReputationManager struct{}

func NewReputationManager() *ReputationManager {
	return &ReputationManager{}
}

// UpdateReputation updates a user's reputation score
func (rm *ReputationManager) UpdateReputation(identity *UserIdentity, activityType string, value int64) {
	// Calculate reputation change based on activity
	var change float64
	switch activityType {
	case "post":
		change = float64(value) * 0.1
	case "comment":
		change = float64(value) * 0.05
	case "like_given":
		change = float64(value) * 0.01 // Minimal gain for giving likes
	case "like_received":
		change = float64(value) * 0.2
	case "transaction":
		change = float64(value) * 0.01
	case "proposal":
		change = float64(value) * 1.0
	case "vote":
		change = float64(value) * 0.1
	case "tip_given":
		change = float64(value) * 0.05 // Generosity rewarded
	case "tip_received":
		change = float64(value) * 0.1 // Value creation rewarded
	default:
		change = 0
	}

	// Update reputation
	identity.Reputation.Overall += change
	identity.Reputation.LastUpdated = time.Now()

	// Add to history
	event := ReputationEvent{
		Type:      activityType,
		Value:     change,
		Reason:    fmt.Sprintf("Activity: %s", activityType),
		Timestamp: time.Now(),
		Source:    "activity",
	}
	identity.Reputation.History = append(identity.Reputation.History, event)
}

// DecayReputation applies entropy to reputation over time.
// If a user is inactive, their influence fades.
func (rm *ReputationManager) DecayReputation(identity *UserIdentity) {
	const DECAY_RATE = 0.5 // Points per day of inactivity

	hoursInactive := time.Since(identity.Activity.LastActive).Hours()

	if hoursInactive > 24 { // Grace period of 1 day
		daysInactive := hoursInactive / 24
		decay := daysInactive * DECAY_RATE

		if decay > 0 {
			identity.Reputation.Overall -= decay
			if identity.Reputation.Overall < 0 {
				identity.Reputation.Overall = 0
			}

			// Log Decay Event
			event := ReputationEvent{
				Type:      "decay",
				Value:     -decay,
				Reason:    fmt.Sprintf("Inactivity for %.1f days", daysInactive),
				Timestamp: time.Now(),
				Source:    "system_entropy",
			}
			identity.Reputation.History = append(identity.Reputation.History, event)
		}
	}
}

// PrivacyManager handles privacy settings
type PrivacyManager struct{}

func NewPrivacyManager() *PrivacyManager {
	return &PrivacyManager{}
}

func (pm *PrivacyManager) GetDefaultSettings() *PrivacySettings {
	return &PrivacySettings{
		ProfileVisibility:  "public",
		ActivityVisibility: "friends",
		TransactionPrivacy: "public",
		ContactVisibility:  "friends",
		LocationSharing:    false,
		DataRetention:      "1year",
		GDPRConsent:        true,
		MarketingConsent:   false,
		CustomSettings:     make(map[string]string),
	}
}

// KYCManager handles KYC operations
type KYCManager struct{}

func NewKYCManager() *KYCManager {
	return &KYCManager{}
}
