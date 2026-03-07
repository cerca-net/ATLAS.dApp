package defi

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// DEX (Decentralized Exchange) for trading
type DEX struct {
	pairs         map[string]*TradingPair
	orders        map[string]*Order
	transactions  []*Trade
	feeRate       float64
}

// TradingPair represents a trading pair
type TradingPair struct {
	Symbol      string  `json:"symbol"`
	BaseAsset   string  `json:"base_asset"`
	QuoteAsset  string  `json:"quote_asset"`
	Price       float64 `json:"price"`
	Volume24h   uint64  `json:"volume_24h"`
	LastUpdated time.Time `json:"last_updated"`
}

// Order represents a trading order
type Order struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	Pair      string    `json:"pair"`
	Type      string    `json:"type"` // "buy" or "sell"
	Amount    uint64    `json:"amount"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"` // "pending", "filled", "cancelled"
	CreatedAt time.Time `json:"created_at"`
	FilledAt  *time.Time `json:"filled_at,omitempty"`
}

// Trade represents a completed trade
type Trade struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Pair      string    `json:"pair"`
	Buyer     string    `json:"buyer"`
	Seller    string    `json:"seller"`
	Amount    uint64    `json:"amount"`
	Price     float64   `json:"price"`
	Fee       uint64    `json:"fee"`
	Timestamp time.Time `json:"timestamp"`
}

// NewDEX creates a new decentralized exchange
func NewDEX() *DEX {
	return &DEX{
		pairs:        make(map[string]*TradingPair),
		orders:       make(map[string]*Order),
		transactions: make([]*Trade, 0),
		feeRate:      0.003, // 0.3% trading fee
	}
}

// AddTradingPair adds a new trading pair
func (dex *DEX) AddTradingPair(symbol, baseAsset, quoteAsset string, initialPrice float64) error {
	if _, exists := dex.pairs[symbol]; exists {
		return fmt.Errorf("trading pair %s already exists", symbol)
	}

	dex.pairs[symbol] = &TradingPair{
		Symbol:      symbol,
		BaseAsset:   baseAsset,
		QuoteAsset:  quoteAsset,
		Price:       initialPrice,
		Volume24h:   0,
		LastUpdated: time.Now(),
	}

	return nil
}

// PlaceOrder places a trading order
func (dex *DEX) PlaceOrder(user, pair, orderType string, amount uint64, price float64) (*Order, error) {
	if _, exists := dex.pairs[pair]; !exists {
		return nil, fmt.Errorf("trading pair %s not found", pair)
	}

	if orderType != "buy" && orderType != "sell" {
		return nil, fmt.Errorf("invalid order type: %s", orderType)
	}

	orderID := generateOrderID(user, pair)
	order := &Order{
		ID:        orderID,
		User:      user,
		Pair:      pair,
		Type:      orderType,
		Amount:    amount,
		Price:     price,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	dex.orders[orderID] = order

	// Try to match orders
	dex.matchOrders(pair)

	return order, nil
}

// CancelOrder cancels a pending order
func (dex *DEX) CancelOrder(orderID, user string) error {
	order, exists := dex.orders[orderID]
	if !exists {
		return fmt.Errorf("order %s not found", orderID)
	}

	if order.User != user {
		return fmt.Errorf("user %s is not the owner of order %s", user, orderID)
	}

	if order.Status != "pending" {
		return fmt.Errorf("order %s cannot be cancelled", orderID)
	}

	order.Status = "cancelled"
	return nil
}

// GetOrderBook returns the order book for a trading pair
func (dex *DEX) GetOrderBook(pair string) ([]*Order, []*Order, error) {
	if _, exists := dex.pairs[pair]; !exists {
		return nil, nil, fmt.Errorf("trading pair %s not found", pair)
	}

	var buyOrders, sellOrders []*Order
	for _, order := range dex.orders {
		if order.Pair == pair && order.Status == "pending" {
			if order.Type == "buy" {
				buyOrders = append(buyOrders, order)
			} else {
				sellOrders = append(sellOrders, order)
			}
		}
	}

	return buyOrders, sellOrders, nil
}

// GetTradeHistory returns trade history for a pair
func (dex *DEX) GetTradeHistory(pair string, limit int) []*Trade {
	var trades []*Trade
	count := 0

	// Get recent trades in reverse order
	for i := len(dex.transactions) - 1; i >= 0 && count < limit; i-- {
		if dex.transactions[i].Pair == pair {
			trades = append(trades, dex.transactions[i])
			count++
		}
	}

	return trades
}

// matchOrders attempts to match buy and sell orders
func (dex *DEX) matchOrders(pair string) {
	// Get all pending orders for this pair
	var buyOrders, sellOrders []*Order
	for _, order := range dex.orders {
		if order.Pair == pair && order.Status == "pending" {
			if order.Type == "buy" {
				buyOrders = append(buyOrders, order)
			} else {
				sellOrders = append(sellOrders, order)
			}
		}
	}

	// Simple matching: match orders at market price
	for _, buyOrder := range buyOrders {
		for _, sellOrder := range sellOrders {
			if buyOrder.Price >= sellOrder.Price {
				// Match found
				dex.executeTrade(buyOrder, sellOrder)
			}
		}
	}
}

// executeTrade executes a trade between two orders
func (dex *DEX) executeTrade(buyOrder, sellOrder *Order) {
	// Determine trade amount (minimum of both orders)
	tradeAmount := buyOrder.Amount
	if sellOrder.Amount < tradeAmount {
		tradeAmount = sellOrder.Amount
	}

	// Calculate trade price (average of both orders)
	tradePrice := (buyOrder.Price + sellOrder.Price) / 2

	// Calculate fee
	fee := uint64(float64(tradeAmount) * tradePrice * dex.feeRate)

	// Create trade record
	tradeID := generateTradeID(buyOrder.ID, sellOrder.ID)
	trade := &Trade{
		ID:        tradeID,
		OrderID:   buyOrder.ID,
		Pair:      buyOrder.Pair,
		Buyer:     buyOrder.User,
		Seller:    sellOrder.User,
		Amount:    tradeAmount,
		Price:     tradePrice,
		Fee:       fee,
		Timestamp: time.Now(),
	}

	dex.transactions = append(dex.transactions, trade)

	// Update orders
	now := time.Now()
	buyOrder.Amount -= tradeAmount
	sellOrder.Amount -= tradeAmount

	if buyOrder.Amount == 0 {
		buyOrder.Status = "filled"
		buyOrder.FilledAt = &now
	}

	if sellOrder.Amount == 0 {
		sellOrder.Status = "filled"
		sellOrder.FilledAt = &now
	}

	// Update trading pair price and volume
	pair := dex.pairs[buyOrder.Pair]
	pair.Price = tradePrice
	pair.Volume24h += tradeAmount
	pair.LastUpdated = time.Now()
}

// GetMarketData returns market data for all pairs
func (dex *DEX) GetMarketData() map[string]*TradingPair {
	return dex.pairs
}

// Helper functions
func generateOrderID(user, pair string) string {
	data := fmt.Sprintf("%s_%s_%d", user, pair, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
}

func generateTradeID(buyOrderID, sellOrderID string) string {
	data := fmt.Sprintf("%s_%s_%d", buyOrderID, sellOrderID, time.Now().Unix())
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:16])
} 