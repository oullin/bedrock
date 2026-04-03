package money

import (
	"github.com/gollin/packages/support/money/calculator"
	"github.com/gollin/packages/support/money/currency"
)

// Amount is a data structure that stores the amount being used for calculations.
type Amount = calculator.Amount

// Money represents a monetary value with an amount and a currency.
type Money struct {
	amount   Amount             `db:"amount"`
	currency *currency.Currency `db:"currency"`
}
