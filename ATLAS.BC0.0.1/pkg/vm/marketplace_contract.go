package vm

import (
	"fmt"
	"log"
)

// Marketplace contract constants
const (
	MarketplaceContractName    = "CercaChainMarketplace"
	MarketplaceContractAddress = "CONTRACT_MARKETPLACE_SYSTEM"
	MarketplaceFeeRate         = int64(20) // 0.2% fee in basis points
)

// Escrow status physical goods / digital goods
const (
	EscrowStatusPending   = int64(0)
	EscrowStatusFunded    = int64(1) // Buyer deposited funds
	EscrowStatusCompleted = int64(2) // Funds released to seller
	EscrowStatusRefunded  = int64(3) // Funds returned to buyer
	EscrowStatusDisputed  = int64(4) // Under arbitration
)

// CreateSystemMarketplaceContract creates the genesis marketplace contract for e-commerce.
// This manages product listings, escrow, and seller payments.
func CreateSystemMarketplaceContract() *Contract {
	contract := &Contract{
		Address:      MarketplaceContractAddress,
		Name:         MarketplaceContractName,
		Version:      "1.0.0",
		Owner:        "SYSTEM",
		Upgradable:   true,
		CreatedAt:    0,
		UpdatedAt:    0,
		ContractType: ContractTypeSystem,
		Functions:    make(map[string]*Function),
		Storage:      make(map[string]interface{}),
	}

	// Initialize default storage
	contract.Storage["fee_rate"] = MarketplaceFeeRate
	contract.Storage["total_volume"] = int64(0)
	contract.Storage["total_fees"] = int64(0)
	contract.Storage["order_counter"] = int64(0)

	// Define functions
	contract.Functions = map[string]*Function{
		"createOrder": {
			Name:       "createOrder",
			Parameters: []string{"buyer", "seller", "amount"},
			Code:       []Instruction{}, // Handled by helper
		},
		"releaseFunds": {
			Name:       "releaseFunds",
			Parameters: []string{"buyer", "orderId"},
			Code:       []Instruction{},
		},
		"refundBuyer": {
			Name:       "refundBuyer",
			Parameters: []string{"seller", "orderId"},
			Code:       []Instruction{},
		},
		"raiseDispute": {
			Name:       "raiseDispute",
			Parameters: []string{"user", "orderId"},
			Code:       []Instruction{},
		},
		"resolveDispute": { // Arbitration by governance/system
			Name:       "resolveDispute",
			Parameters: []string{"arbitrator", "orderId", "resolution"}, // 0=refund, 1=release
			Code:       []Instruction{},
		},
	}

	return contract
}

// MarketplaceContractHelper provides high-level Go functions for commerce operations.
type MarketplaceContractHelper struct {
	Contract *Contract
	VM       *VM
}

// NewMarketplaceContractHelper creates a helper for marketplace contract interaction
func NewMarketplaceContractHelper(contract *Contract) *MarketplaceContractHelper {
	vmInstance := NewVM()
	vmInstance.RegisterSystemContract(contract.Address, []string{
		"createOrder", "releaseFunds", "refundBuyer", "raiseDispute", "resolveDispute",
	})
	return &MarketplaceContractHelper{
		Contract: contract,
		VM:       vmInstance,
	}
}

// OrderInfo struct returns public info about an escrowed transaction
type OrderInfo struct {
	OrderID   string `json:"order_id"`
	Buyer     string `json:"buyer"`
	Seller    string `json:"seller"`
	Amount    int64  `json:"amount"`
	Fee       int64  `json:"fee"`
	Status    int64  `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

// CreateOrder places the buyer's funds in escrow for a particular seller.
func (h *MarketplaceContractHelper) CreateOrder(ctx *ExecutionContext, buyer, seller string, amount int64, orderID string) (string, error) {
	if ctx.State == nil {
		return "", fmt.Errorf("state adapter is required")
	}

	if amount <= 0 {
		return "", fmt.Errorf("order amount must be greater than 0")
	}

	if buyer == seller {
		return "", fmt.Errorf("buyer and seller cannot be the same")
	}

	// Check buyer balance
	if ctx.State.GetBalance(buyer) < amount {
		return "", fmt.Errorf("insufficient balance for order")
	}

	if orderID == "" {
		// Generate Order ID if not provided
		orderCounter, _ := ctx.State.GetContractStorage(h.Contract.Address, "order_counter")
		newCounter := orderCounter + 1
		ctx.State.SetContractStorage(h.Contract.Address, "order_counter", newCounter)
		orderID = fmt.Sprintf("ORDER_%d", newCounter)
	}

	// Move funds to contract (escrow)
	if err := ctx.State.Transfer(buyer, h.Contract.Address, amount); err != nil {
		return "", fmt.Errorf("failed to fund escrow: %v", err)
	}

	// Store order state
	ctx.State.SetStringStorage(h.Contract.Address, fmt.Sprintf("order_buyer_%s", orderID), buyer)
	ctx.State.SetStringStorage(h.Contract.Address, fmt.Sprintf("order_seller_%s", orderID), seller)
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("order_amount_%s", orderID), amount)
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID), EscrowStatusFunded)
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("order_created_%s", orderID), ctx.Timestamp)

	ctx.EmitEvent(h.Contract.Address, "OrderCreated", map[string]interface{}{
		"order_id": orderID,
		"buyer":    buyer,
		"seller":   seller,
		"amount":   amount,
	})

	log.Printf("[MARKETPLACE] Order %s created: %d TCOIN from %s to %s", orderID, amount, shortAddr(buyer), shortAddr(seller))
	return orderID, nil
}

// ReleaseFunds releases funds from escrow to the seller upon confirmation of delivery by the buyer.
func (h *MarketplaceContractHelper) ReleaseFunds(ctx *ExecutionContext, buyer, orderID string) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required")
	}

	// Validate order
	status, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID))
	if !exists {
		return fmt.Errorf("order %s not found", orderID)
	}
	if status != EscrowStatusFunded {
		return fmt.Errorf("order %s is not in funded status (status=%d)", orderID, status)
	}

	orderBuyer, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_buyer_%s", orderID))
	if buyer != orderBuyer {
		return fmt.Errorf("only the buyer (%s) can release funds, got %s", orderBuyer, buyer)
	}

	orderSeller, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_seller_%s", orderID))
	amount, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("order_amount_%s", orderID))

	// Calculate fee
	feeRate, _ := ctx.State.GetContractStorage(h.Contract.Address, "fee_rate")
	if feeRate == 0 {
		feeRate = MarketplaceFeeRate
	}

	fee := (amount * feeRate) / 10000
	sellerAmount := amount - fee

	// Transfer funds to seller
	if err := ctx.State.Transfer(h.Contract.Address, orderSeller, sellerAmount); err != nil {
		return fmt.Errorf("failed to release funds to seller: %v", err)
	}

	// If there's a fee, burn it or send to treasury. Let's burn it as a network sink.
	if fee > 0 {
		ctx.State.Burn(h.Contract.Address, fee)
		totalFees, _ := ctx.State.GetContractStorage(h.Contract.Address, "total_fees")
		ctx.State.SetContractStorage(h.Contract.Address, "total_fees", totalFees+fee)
	}

	// Update status
	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID), EscrowStatusCompleted)

	// Update analytics
	totalVolume, _ := ctx.State.GetContractStorage(h.Contract.Address, "total_volume")
	ctx.State.SetContractStorage(h.Contract.Address, "total_volume", totalVolume+amount)

	ctx.EmitEvent(h.Contract.Address, "FundsReleased", map[string]interface{}{
		"order_id": orderID,
		"seller":   orderSeller,
		"amount":   sellerAmount,
		"fee":      fee,
	})

	log.Printf("[MARKETPLACE] Order %s completed: %d to %s (fee: %d)", orderID, sellerAmount, shortAddr(orderSeller), fee)
	return nil
}

// RefundBuyer allows the seller to cancel the order and refund the buyer.
func (h *MarketplaceContractHelper) RefundBuyer(ctx *ExecutionContext, seller, orderID string) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required")
	}

	status, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID))
	if !exists {
		return fmt.Errorf("order %s not found", orderID)
	}
	if status != EscrowStatusFunded {
		return fmt.Errorf("order %s cannot be refunded (status=%d)", orderID, status)
	}

	orderSeller, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_seller_%s", orderID))
	if seller != orderSeller {
		return fmt.Errorf("only the seller (%s) can issue a refund willingly, got %s", orderSeller, seller)
	}

	orderBuyer, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_buyer_%s", orderID))
	amount, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("order_amount_%s", orderID))

	// Return funds to buyer in full
	if err := ctx.State.Transfer(h.Contract.Address, orderBuyer, amount); err != nil {
		return fmt.Errorf("failed to return funds to buyer: %v", err)
	}

	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID), EscrowStatusRefunded)

	ctx.EmitEvent(h.Contract.Address, "OrderRefunded", map[string]interface{}{
		"order_id": orderID,
		"buyer":    orderBuyer,
		"amount":   amount,
	})

	log.Printf("[MARKETPLACE] Order %s refunded: %d returned to %s", orderID, amount, shortAddr(orderBuyer))
	return nil
}

// RaiseDispute freezes the order. Neither party can release/refund until resolved by arbiter.
func (h *MarketplaceContractHelper) RaiseDispute(ctx *ExecutionContext, user, orderID string) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required")
	}

	status, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID))
	if !exists {
		return fmt.Errorf("order %s not found", orderID)
	}
	if status != EscrowStatusFunded {
		return fmt.Errorf("order %s cannot be disputed, not active", orderID)
	}

	buyer, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_buyer_%s", orderID))
	seller, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_seller_%s", orderID))

	if user != buyer && user != seller {
		return fmt.Errorf("only buyer or seller can dispute")
	}

	ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID), EscrowStatusDisputed)

	ctx.EmitEvent(h.Contract.Address, "OrderDisputed", map[string]interface{}{
		"order_id":  orderID,
		"raised_by": user,
	})

	return nil
}

// ResolveDispute resolves a disputed order, paying either the buyer or the seller.
func (h *MarketplaceContractHelper) ResolveDispute(ctx *ExecutionContext, admin, orderID string, payBuyer bool) error {
	if ctx.State == nil {
		return fmt.Errorf("state adapter is required")
	}

	status, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID))
	if !exists {
		return fmt.Errorf("order %s not found", orderID)
	}
	if status != EscrowStatusDisputed {
		return fmt.Errorf("order %s is not disputed", orderID)
	}

	buyer, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_buyer_%s", orderID))
	seller, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_seller_%s", orderID))
	amountStr, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_amount_%s", orderID))

	var amount uint64
	fmt.Sscanf(amountStr, "%d", &amount)

	var payee string
	if payBuyer {
		payee = buyer
		ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID), EscrowStatusRefunded)
	} else {
		payee = seller
		ctx.State.SetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID), EscrowStatusCompleted)
	}

	err := ctx.State.Transfer(h.Contract.Address, payee, int64(amount))
	if err != nil {
		return fmt.Errorf("failed to transfer funds to %s: %v", payee, err)
	}

	ctx.EmitEvent(h.Contract.Address, "OrderDisputeResolved", map[string]interface{}{
		"order_id":    orderID,
		"resolved_by": admin,
		"payee":       payee,
		"amount":      amount,
	})

	log.Printf("[MARKETPLACE] Order %s dispute resolved: %d paid to %s", orderID, amount, shortAddr(payee))
	return nil
}

// GetOrder returns view-only details for an order
func (h *MarketplaceContractHelper) GetOrder(ctx *ExecutionContext, orderID string) *OrderInfo {
	if ctx.State == nil {
		return nil
	}

	status, exists := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("order_status_%s", orderID))
	if !exists {
		return nil
	}

	buyer, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_buyer_%s", orderID))
	seller, _ := ctx.State.GetStringStorage(h.Contract.Address, fmt.Sprintf("order_seller_%s", orderID))
	amount, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("order_amount_%s", orderID))
	created, _ := ctx.State.GetContractStorage(h.Contract.Address, fmt.Sprintf("order_created_%s", orderID))

	feeRate, _ := ctx.State.GetContractStorage(h.Contract.Address, "fee_rate")
	if feeRate == 0 {
		feeRate = MarketplaceFeeRate
	}

	fee := (amount * feeRate) / 10000

	return &OrderInfo{
		OrderID:   orderID,
		Buyer:     buyer,
		Seller:    seller,
		Amount:    amount,
		Fee:       fee,
		Status:    status,
		CreatedAt: created,
	}
}
