package state

import (
	"errors"
)

type Account struct {
	Address      string
	Balance      int64
	Nonce        uint64
	IsValidator  bool
	StakedAmount int64
}

// ErrInsufficientFunds is returned when a sender doesn't have enough balance for a transaction.
var ErrInsufficientFunds = errors.New("insufficient funds")

// ErrInvalidTransaction is returned when a transaction is invalid.
var ErrInvalidTransaction = errors.New("invalid transaction")
