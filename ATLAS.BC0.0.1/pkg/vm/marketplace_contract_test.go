package vm

import (
	"testing"
)

func TestMarketplaceContract(t *testing.T) {
	mpContract := CreateSystemMarketplaceContract()
	helper := NewMarketplaceContractHelper(mpContract)

	stateAdapter := NewMockStateAdapter()
	ctx := &ExecutionContext{
		State:     stateAdapter,
		Timestamp: 1672531200,
	}

	stateAdapter.balances["0xBUYER"] = 50000

	t.Run("CreateOrder_Success", func(t *testing.T) {
		orderID, err := helper.CreateOrder(ctx, "0xBUYER", "0xSELLER", 10000)
		if err != nil {
			t.Fatalf("Failed to create order: %v", err)
		}

		if orderID != "ORDER_1" {
			t.Errorf("Expected ORDER_1, got %s", orderID)
		}

		// Contract should hold the funds
		ctBalance := stateAdapter.GetBalance(MarketplaceContractAddress)
		if ctBalance != 10000 {
			t.Errorf("Expected contract balance 10000, got %d", ctBalance)
		}

		// Buyer should have 40000 left
		buyerBalance := stateAdapter.GetBalance("0xBUYER")
		if buyerBalance != 40000 {
			t.Errorf("Expected buyer balance 40000, got %d", buyerBalance)
		}

		// Check order info
		order := helper.GetOrder(ctx, orderID)
		if order == nil {
			t.Fatalf("Order not found")
		}
		if order.Status != EscrowStatusFunded {
			t.Errorf("Expected status %d, got %d", EscrowStatusFunded, order.Status)
		}
	})

	t.Run("CreateOrder_InsufficientFunds", func(t *testing.T) {
		_, err := helper.CreateOrder(ctx, "0xPOORBUYER", "0xSELLER", 10000)
		if err == nil {
			t.Fatalf("Expected error for insufficient funds")
		}
	})

	t.Run("ReleaseFunds_Success", func(t *testing.T) {
		// Buyer releases funds to Seller
		err := helper.ReleaseFunds(ctx, "0xBUYER", "ORDER_1")
		if err != nil {
			t.Fatalf("Failed to release funds: %v", err)
		}

		// Fee is 0.2% = 20
		// Seller should receive 9980
		sellerBalance := stateAdapter.GetBalance("0xSELLER")
		if sellerBalance != 9980 {
			t.Errorf("Expected seller to have 9980, got %d", sellerBalance)
		}

		// Fee should be burned (removed from system)
		// Initially total supply in adapter is 50000
		// We expect it to drop by 20. But MockAdapter Transfer doesn't affect total supply directly unless explicit burn.
		// So we won't assert total supply here unless we check Burn history, but let's just check the order status.

		order := helper.GetOrder(ctx, "ORDER_1")
		if order.Status != EscrowStatusCompleted {
			t.Errorf("Expected status %d, got %d", EscrowStatusCompleted, order.Status)
		}
	})

	t.Run("RefundBuyer", func(t *testing.T) {
		orderID, _ := helper.CreateOrder(ctx, "0xBUYER", "0xSELLER", 5000)

		err := helper.RefundBuyer(ctx, "0xSELLER", orderID)
		if err != nil {
			t.Fatalf("Failed to refund: %v", err)
		}

		order := helper.GetOrder(ctx, orderID)
		if order.Status != EscrowStatusRefunded {
			t.Errorf("Expected status %d, got %d", EscrowStatusRefunded, order.Status)
		}
	})

	t.Run("RaiseDispute", func(t *testing.T) {
		orderID, _ := helper.CreateOrder(ctx, "0xBUYER", "0xSELLER", 5000)

		err := helper.RaiseDispute(ctx, "0xBUYER", orderID)
		if err != nil {
			t.Fatalf("Failed to dispute: %v", err)
		}

		order := helper.GetOrder(ctx, orderID)
		if order.Status != EscrowStatusDisputed {
			t.Errorf("Expected status %d, got %d", EscrowStatusDisputed, order.Status)
		}
	})
}
