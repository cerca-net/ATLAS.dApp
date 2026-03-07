package defi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"atlas-blockchain/pkg/vm"
	"atlas-blockchain/internal/identity"
)

// OracleNetwork manages external data feeds
type OracleNetwork struct {
	oracles      map[string]*Oracle
	aggregator   *PriceAggregator
	validator    *OracleValidator
}

// Oracle represents an external data source
type Oracle struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	URL      string    `json:"url"`
	DataType string    `json:"data_type"`
	IsActive bool      `json:"is_active"`
	LastUpdate time.Time `json:"last_update"`
}

// PriceAggregator aggregates prices from multiple oracles
type PriceAggregator struct {
	prices map[string][]float64
}

// OracleValidator validates oracle data
type OracleValidator struct {
	threshold float64
}

// NewOracleNetwork creates a new oracle network
func NewOracleNetwork() *OracleNetwork {
	return &OracleNetwork{
		oracles:    make(map[string]*Oracle),
		aggregator: &PriceAggregator{prices: make(map[string][]float64)},
		validator:  &OracleValidator{threshold: 0.1}, // 10% deviation threshold
	}
}

// AddOracle adds a new oracle
func (on *OracleNetwork) AddOracle(id, name, url, dataType string) error {
	if _, exists := on.oracles[id]; exists {
		return fmt.Errorf("oracle %s already exists", id)
	}

	on.oracles[id] = &Oracle{
		ID:         id,
		Name:       name,
		URL:        url,
		DataType:   dataType,
		IsActive:   true,
		LastUpdate: time.Now(),
	}

	return nil
}

// GetPrice gets aggregated price for an asset
func (on *OracleNetwork) GetPrice(asset string) (float64, error) {
	prices, exists := on.aggregator.prices[asset]
	if !exists || len(prices) == 0 {
		return 0, fmt.Errorf("no price data available for %s", asset)
	}

	// Simple median calculation
	if len(prices) == 1 {
		return prices[0], nil
	}

	// Calculate median price
	total := 0.0
	for _, price := range prices {
		total += price
	}
	return total / float64(len(prices)), nil
}

// UpdatePrice updates price from an oracle
func (on *OracleNetwork) UpdatePrice(oracleID, asset string, price float64) error {
	oracle, exists := on.oracles[oracleID]
	if !exists {
		return fmt.Errorf("oracle %s not found", oracleID)
	}

	if !oracle.IsActive {
		return fmt.Errorf("oracle %s is not active", oracleID)
	}

	// Validate price deviation
	if !on.validatePrice(asset, price) {
		return fmt.Errorf("price deviation too high for %s", asset)
	}

	// Update price in aggregator
	on.aggregator.prices[asset] = append(on.aggregator.prices[asset], price)
	oracle.LastUpdate = time.Now()

	return nil
}

// validatePrice validates price against existing prices
func (on *OracleNetwork) validatePrice(asset string, newPrice float64) bool {
	existingPrices, exists := on.aggregator.prices[asset]
	if !exists || len(existingPrices) == 0 {
		return true // First price is always valid
	}

	// Calculate average existing price
	total := 0.0
	for _, price := range existingPrices {
		total += price
	}
	avgPrice := total / float64(len(existingPrices))

	// Check deviation
	deviation := (newPrice - avgPrice) / avgPrice
	if deviation < 0 {
		deviation = -deviation
	}

	return deviation <= on.validator.threshold
}

// RiskManager manages DeFi risk
type RiskManager struct {
	riskMetrics map[string]*RiskMetric
	alerts      []*RiskAlert
}

// RiskMetric represents a risk metric
type RiskMetric struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Threshold float64   `json:"threshold"`
	Severity  string    `json:"severity"` // "low", "medium", "high", "critical"
	LastUpdate time.Time `json:"last_update"`
}

// RiskAlert represents a risk alert
type RiskAlert struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
	Acknowledged bool   `json:"acknowledged"`
}

// NewRiskManager creates a new risk manager
func NewRiskManager() *RiskManager {
	return &RiskManager{
		riskMetrics: make(map[string]*RiskMetric),
		alerts:      make([]*RiskAlert, 0),
	}
}

// AddRiskMetric adds a new risk metric
func (rm *RiskManager) AddRiskMetric(name string, threshold float64) {
	rm.riskMetrics[name] = &RiskMetric{
		Name:       name,
		Value:      0,
		Threshold:  threshold,
		Severity:   "low",
		LastUpdate: time.Now(),
	}
}

// UpdateRiskMetric updates a risk metric
func (rm *RiskManager) UpdateRiskMetric(name string, value float64) error {
	metric, exists := rm.riskMetrics[name]
	if !exists {
		return fmt.Errorf("risk metric %s not found", name)
	}

	metric.Value = value
	metric.LastUpdate = time.Now()

	// Update severity based on threshold
	if value >= metric.Threshold {
		metric.Severity = "critical"
		rm.createAlert(name, "critical", fmt.Sprintf("Risk metric %s exceeded threshold: %.2f", name, value))
	} else if value >= metric.Threshold*0.8 {
		metric.Severity = "high"
		rm.createAlert(name, "high", fmt.Sprintf("Risk metric %s approaching threshold: %.2f", name, value))
	} else if value >= metric.Threshold*0.6 {
		metric.Severity = "medium"
	} else {
		metric.Severity = "low"
	}

	return nil
}

// createAlert creates a new risk alert
func (rm *RiskManager) createAlert(metricName, severity, message string) {
	alert := &RiskAlert{
		ID:          generateAlertID(metricName),
		Type:        "risk_metric",
		Message:     message,
		Severity:    severity,
		Timestamp:   time.Now(),
		Acknowledged: false,
	}

	rm.alerts = append(rm.alerts, alert)
}

// GetRiskMetrics returns all risk metrics
func (rm *RiskManager) GetRiskMetrics() map[string]*RiskMetric {
	return rm.riskMetrics
}

// GetAlerts returns all risk alerts
func (rm *RiskManager) GetAlerts() []*RiskAlert {
	return rm.alerts
}

// AcknowledgeAlert acknowledges a risk alert
func (rm *RiskManager) AcknowledgeAlert(alertID string) error {
	for _, alert := range rm.alerts {
		if alert.ID == alertID {
			alert.Acknowledged = true
			return nil
		}
	}
	return fmt.Errorf("alert %s not found", alertID)
}



// SmartContract represents a DeFi smart contract
type SmartContract struct {
	Address     string                 `json:"address"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"` // "lending", "trading", "staking", "governance"
	Code        []vm.Instruction       `json:"code"`
	State       map[string]interface{} `json:"state"`
	Balance     uint64                 `json:"balance"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	IsActive    bool                   `json:"is_active"`
}

// DeFiContractManager manages DeFi smart contracts
type DeFiContractManager struct {
	contracts map[string]*SmartContract
	vm        *vm.VM
}

// NewDeFiContractManager creates a new DeFi contract manager
func NewDeFiContractManager(vm *vm.VM) *DeFiContractManager {
	return &DeFiContractManager{
		contracts: make(map[string]*SmartContract),
		vm:        vm,
	}
}

// DeployContract deploys a new DeFi contract
func (dcm *DeFiContractManager) DeployContract(name, contractType string, code []vm.Instruction, initialBalance uint64) (*SmartContract, error) {
	address := generateContractAddress(name)
	
	contract := &SmartContract{
		Address:   address,
		Name:      name,
		Type:      contractType,
		Code:      code,
		State:     make(map[string]interface{}),
		Balance:   initialBalance,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsActive:  true,
	}

	dcm.contracts[address] = contract

	// Register contract with VM
	dcm.vm.RegisterSystemContract(address, []string{"execute", "getState", "setState"})

	return contract, nil
}

// ExecuteContract executes a DeFi contract
func (dcm *DeFiContractManager) ExecuteContract(address string, function string, params map[string]interface{}) (interface{}, error) {
	contract, exists := dcm.contracts[address]
	if !exists {
		return nil, fmt.Errorf("contract %s not found", address)
	}

	if !contract.IsActive {
		return nil, fmt.Errorf("contract %s is not active", address)
	}

	// Create execution context
	context := vm.NewExecutionContext("", 1000000)
	context.SetContractContext(address, function)
	
	// Add parameters
	for _, param := range params {
		context.AddParameter(param)
	}

	// Execute contract code
	if err := dcm.vm.Execute(contract.Code, context); err != nil {
		return nil, fmt.Errorf("contract execution failed: %v", err)
	}

	contract.UpdatedAt = time.Now()
	return dcm.vm.GetGasUsed(), nil
}

// GetContract returns a contract by address
func (dcm *DeFiContractManager) GetContract(address string) (*SmartContract, error) {
	contract, exists := dcm.contracts[address]
	if !exists {
		return nil, fmt.Errorf("contract %s not found", address)
	}
	return contract, nil
}

// GetContractsByType returns all contracts of a specific type
func (dcm *DeFiContractManager) GetContractsByType(contractType string) []*SmartContract {
	var contracts []*SmartContract
	for _, contract := range dcm.contracts {
		if contract.Type == contractType {
			contracts = append(contracts, contract)
		}
	}
	return contracts
}

// ActivityRewardManager manages activity-based token rewards
type ActivityRewardManager struct {
	rewards     map[string]*ActivityReward
	rates       map[string]float64
	identity    *identity.IdentityManager
}

// ActivityReward represents a reward for user activity
type ActivityReward struct {
	User       string    `json:"user"`
	Activity   string    `json:"activity"`
	Amount     uint64    `json:"amount"`
	Rate       float64   `json:"rate"`
	Timestamp  time.Time `json:"timestamp"`
	Claimed    bool      `json:"claimed"`
	ClaimedAt  *time.Time `json:"claimed_at,omitempty"`
}

// NewActivityRewardManager creates a new activity reward manager
func NewActivityRewardManager(identityManager *identity.IdentityManager) *ActivityRewardManager {
	return &ActivityRewardManager{
		rewards:  make(map[string]*ActivityReward),
		rates:    make(map[string]float64),
		identity: identityManager,
	}

	// Set default reward rates
	arm := &ActivityRewardManager{
		rewards:  make(map[string]*ActivityReward),
		rates:    make(map[string]float64),
		identity: identityManager,
	}

	// Set default reward rates for different activities
	arm.rates["post"] = 1.0
	arm.rates["comment"] = 0.5
	arm.rates["like_given"] = 0.1
	arm.rates["like_received"] = 0.2
	arm.rates["transaction"] = 0.1
	arm.rates["proposal"] = 10.0
	arm.rates["vote"] = 1.0

	return arm
}

// CalculateReward calculates reward for user activity
func (arm *ActivityRewardManager) CalculateReward(user, activity string, activityValue int64) uint64 {
	rate, exists := arm.rates[activity]
	if !exists {
		return 0
	}

	// Get user identity for reputation multiplier
	identity, err := arm.identity.GetIdentity(user)
	if err != nil {
		// Use base rate if identity not found
		return uint64(float64(activityValue) * rate)
	}

	// Apply reputation multiplier
	reputationMultiplier := 1.0 + (identity.Reputation.Overall / 100.0)
	
	return uint64(float64(activityValue) * rate * reputationMultiplier)
}

// AwardReward awards tokens for user activity
func (arm *ActivityRewardManager) AwardReward(user, activity string, amount uint64) error {
	rewardID := generateRewardID(user, activity)
	
	reward := &ActivityReward{
		User:      user,
		Activity:  activity,
		Amount:    amount,
		Rate:      arm.rates[activity],
		Timestamp: time.Now(),
		Claimed:   false,
	}

	arm.rewards[rewardID] = reward

	// Update user activity metrics
	arm.identity.UpdateActivity(user, "tokens_earned", int64(amount))

	return nil
}

// ClaimReward claims a user's reward
func (arm *ActivityRewardManager) ClaimReward(user, activity string) (uint64, error) {
	rewardID := generateRewardID(user, activity)
	
	reward, exists := arm.rewards[rewardID]
	if !exists {
		return 0, fmt.Errorf("reward not found for user %s and activity %s", user, activity)
	}

	if reward.Claimed {
		return 0, fmt.Errorf("reward already claimed")
	}

	reward.Claimed = true
	now := time.Now()
	reward.ClaimedAt = &now

	return reward.Amount, nil
}

// GetUserRewards returns all rewards for a user
func (arm *ActivityRewardManager) GetUserRewards(user string) []*ActivityReward {
	var userRewards []*ActivityReward
	for _, reward := range arm.rewards {
		if reward.User == user {
			userRewards = append(userRewards, reward)
		}
	}
	return userRewards
}

// GetPendingRewards returns pending rewards for a user
func (arm *ActivityRewardManager) GetPendingRewards(user string) []*ActivityReward {
	var pendingRewards []*ActivityReward
	for _, reward := range arm.rewards {
		if reward.User == user && !reward.Claimed {
			pendingRewards = append(pendingRewards, reward)
		}
	}
	return pendingRewards
}

// Helper functions
func generateAlertID(metricName string) string {
	data := fmt.Sprintf("alert_%s_%d", metricName, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func generateContractAddress(name string) string {
	data := fmt.Sprintf("contract_%s_%d", name, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:20])
}

func generateRewardID(user, activity string) string {
	data := fmt.Sprintf("reward_%s_%s_%d", user, activity, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
} 