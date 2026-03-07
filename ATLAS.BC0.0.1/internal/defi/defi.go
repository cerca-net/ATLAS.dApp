package defi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"atlas-blockchain/internal/identity"
)

// DeFiManager handles all DeFi operations
type DeFiManager struct {
	lending      *LendingPool
	trading      *DEX
	staking      *StakingPool
	governance   *GovernanceSystem
	oracles      *OracleNetwork
	risk         *RiskManager
	identity     *identity.IdentityManager
}

// NewDeFiManager creates a new DeFi manager
func NewDeFiManager(identityManager *identity.IdentityManager) *DeFiManager {
	return &DeFiManager{
		lending:      NewLendingPool(),
		trading:      NewDEX(),
		staking:      NewStakingPool(),
		governance:   NewGovernanceSystem(),
		oracles:      NewOracleNetwork(),
		risk:         NewRiskManager(),
		identity:     identityManager,
	}
}

// LendingPool manages lending and borrowing operations
type LendingPool struct {
	assets        map[string]*Asset
	loans         map[string]*Loan
	interestRates map[string]float64
	utilization   map[string]float64
	riskParams    *RiskParameters
}

// Asset represents a lending asset
type Asset struct {
	Symbol        string  `json:"symbol"`
	Address       string  `json:"address"`
	TotalSupply   uint64  `json:"total_supply"`
	TotalBorrowed uint64  `json:"total_borrowed"`
	ReserveFactor float64 `json:"reserve_factor"`
	CollateralFactor float64 `json:"collateral_factor"`
	InterestRate  float64 `json:"interest_rate"`
	LastUpdated   time.Time `json:"last_updated"`
}

// Loan represents a borrowing position
type Loan struct {
	ID            string    `json:"id"`
	Borrower      string    `json:"borrower"`
	Asset         string    `json:"asset"`
	Amount        uint64    `json:"amount"`
	Collateral    []Collateral `json:"collateral"`
	InterestRate  float64   `json:"interest_rate"`
	BorrowedAt    time.Time `json:"borrowed_at"`
	LastPayment   time.Time `json:"last_payment"`
	IsLiquidated  bool      `json:"is_liquidated"`
	LiquidationAt *time.Time `json:"liquidation_at,omitempty"`
}

// Collateral represents collateral for a loan
type Collateral struct {
	Asset   string  `json:"asset"`
	Amount  uint64  `json:"amount"`
	Value   uint64  `json:"value"`
	Factor  float64 `json:"factor"`
}

// RiskParameters defines lending risk parameters
type RiskParameters struct {
	MaxUtilization    float64 `json:"max_utilization"`
	LiquidationThreshold float64 `json:"liquidation_threshold"`
	LiquidationPenalty float64 `json:"liquidation_penalty"`
	MinCollateralRatio float64 `json:"min_collateral_ratio"`
}

// NewLendingPool creates a new lending pool
func NewLendingPool() *LendingPool {
	return &LendingPool{
		assets:        make(map[string]*Asset),
		loans:         make(map[string]*Loan),
		interestRates: make(map[string]float64),
		utilization:   make(map[string]float64),
		riskParams: &RiskParameters{
			MaxUtilization:        0.85,
			LiquidationThreshold:  0.8,
			LiquidationPenalty:    0.1,
			MinCollateralRatio:    1.5,
		},
	}
}

// AddAsset adds a new lending asset
func (lp *LendingPool) AddAsset(symbol, address string, collateralFactor, reserveFactor float64) error {
	if _, exists := lp.assets[symbol]; exists {
		return fmt.Errorf("asset %s already exists", symbol)
	}

	lp.assets[symbol] = &Asset{
		Symbol:          symbol,
		Address:         address,
		TotalSupply:     0,
		TotalBorrowed:   0,
		ReserveFactor:   reserveFactor,
		CollateralFactor: collateralFactor,
		InterestRate:    0.05, // 5% base rate
		LastUpdated:     time.Now(),
	}

	return nil
}

// Supply assets to the lending pool
func (lp *LendingPool) Supply(asset, user string, amount uint64) error {
	assetInfo, exists := lp.assets[asset]
	if !exists {
		return fmt.Errorf("asset %s not found", asset)
	}

	// Update asset supply
	assetInfo.TotalSupply += amount
	assetInfo.LastUpdated = time.Now()

	// Update utilization rate
	lp.updateUtilization(asset)

	// Update interest rate based on utilization
	lp.updateInterestRate(asset)

	return nil
}

// Borrow assets from the lending pool
func (lp *LendingPool) Borrow(asset, user string, amount uint64, collateral []Collateral) (*Loan, error) {
	assetInfo, exists := lp.assets[asset]
	if !exists {
		return nil, fmt.Errorf("asset %s not found", asset)
	}

	// Check if enough liquidity
	if assetInfo.TotalBorrowed+amount > assetInfo.TotalSupply {
		return nil, fmt.Errorf("insufficient liquidity for %s", asset)
	}

	// Validate collateral
	if err := lp.validateCollateral(collateral, amount); err != nil {
		return nil, err
	}

	// Create loan
	loanID := generateLoanID(user, asset)
	loan := &Loan{
		ID:           loanID,
		Borrower:     user,
		Asset:        asset,
		Amount:       amount,
		Collateral:   collateral,
		InterestRate: assetInfo.InterestRate,
		BorrowedAt:   time.Now(),
		LastPayment:  time.Now(),
		IsLiquidated: false,
	}

	// Update asset borrowed amount
	assetInfo.TotalBorrowed += amount
	assetInfo.LastUpdated = time.Now()

	// Store loan
	lp.loans[loanID] = loan

	// Update utilization and interest rate
	lp.updateUtilization(asset)
	lp.updateInterestRate(asset)

	return loan, nil
}

// Repay loan
func (lp *LendingPool) Repay(loanID, user string, amount uint64) error {
	loan, exists := lp.loans[loanID]
	if !exists {
		return fmt.Errorf("loan %s not found", loanID)
	}

	if loan.Borrower != user {
		return fmt.Errorf("user %s is not the borrower of loan %s", user, loanID)
	}

	if loan.IsLiquidated {
		return fmt.Errorf("loan %s has been liquidated", loanID)
	}

	// Calculate interest
	interest := lp.calculateInterest(loan)
	totalOwed := loan.Amount + uint64(interest)

	if amount > totalOwed {
		amount = totalOwed
	}

	// Update loan amount
	loan.Amount = totalOwed - amount
	loan.LastPayment = time.Now()

	// Update asset borrowed amount
	assetInfo := lp.assets[loan.Asset]
	assetInfo.TotalBorrowed -= amount
	assetInfo.LastUpdated = time.Now()

	// Update utilization and interest rate
	lp.updateUtilization(loan.Asset)
	lp.updateInterestRate(loan.Asset)

	return nil
}

// Liquidate undercollateralized loans
func (lp *LendingPool) Liquidate(loanID, liquidator string) error {
	loan, exists := lp.loans[loanID]
	if !exists {
		return fmt.Errorf("loan %s not found", loanID)
	}

	if loan.IsLiquidated {
		return fmt.Errorf("loan %s already liquidated", loanID)
	}

	// Check if loan is undercollateralized
	collateralValue := lp.calculateCollateralValue(loan.Collateral)
	borrowedValue := float64(loan.Amount)
	collateralRatio := collateralValue / borrowedValue

	if collateralRatio > lp.riskParams.LiquidationThreshold {
		return fmt.Errorf("loan %s is not undercollateralized", loanID)
	}

	// Execute liquidation
	now := time.Now()
	loan.IsLiquidated = true
	loan.LiquidationAt = &now

	// Calculate liquidation penalty (for future use)
	_ = uint64(float64(loan.Amount) * lp.riskParams.LiquidationPenalty)

	// Update asset borrowed amount
	assetInfo := lp.assets[loan.Asset]
	assetInfo.TotalBorrowed -= loan.Amount
	assetInfo.LastUpdated = time.Now()

	// Update utilization and interest rate
	lp.updateUtilization(loan.Asset)
	lp.updateInterestRate(loan.Asset)

	return nil
}

// Helper methods for lending pool
func (lp *LendingPool) validateCollateral(collateral []Collateral, borrowedAmount uint64) error {
	totalCollateralValue := uint64(0)
	for _, col := range collateral {
		assetInfo, exists := lp.assets[col.Asset]
		if !exists {
			return fmt.Errorf("collateral asset %s not found", col.Asset)
		}
		totalCollateralValue += uint64(float64(col.Value) * assetInfo.CollateralFactor)
	}

	if float64(totalCollateralValue) < float64(borrowedAmount)*lp.riskParams.MinCollateralRatio {
		return fmt.Errorf("insufficient collateral for loan amount")
	}

	return nil
}

func (lp *LendingPool) updateUtilization(asset string) {
	assetInfo := lp.assets[asset]
	if assetInfo.TotalSupply > 0 {
		lp.utilization[asset] = float64(assetInfo.TotalBorrowed) / float64(assetInfo.TotalSupply)
	} else {
		lp.utilization[asset] = 0
	}
}

func (lp *LendingPool) updateInterestRate(asset string) {
	utilization := lp.utilization[asset]
	baseRate := 0.05 // 5% base rate

	// Simple interest rate model based on utilization
	if utilization > 0.8 {
		// High utilization - increase rate
		lp.assets[asset].InterestRate = baseRate + (utilization-0.8)*0.5
	} else {
		// Low utilization - decrease rate
		lp.assets[asset].InterestRate = baseRate * utilization
	}

	lp.assets[asset].LastUpdated = time.Now()
}

func (lp *LendingPool) calculateInterest(loan *Loan) float64 {
	daysSincePayment := time.Since(loan.LastPayment).Hours() / 24
	return float64(loan.Amount) * loan.InterestRate * daysSincePayment / 365
}

func (lp *LendingPool) calculateCollateralValue(collateral []Collateral) float64 {
	total := 0.0
	for _, col := range collateral {
		total += float64(col.Value)
	}
	return total
}

// Helper functions
func generateLoanID(user, asset string) string {
	data := fmt.Sprintf("%s_%s_%d", user, asset, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
} 