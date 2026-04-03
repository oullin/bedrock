package money

import (
	"github.com/gollin/packages/support/money/exception"
)

// Aggregator performs aggregate operations on Money objects using a shared Manager.
type Aggregator struct {
	manager *Manager
}

// NewAggregator creates a new Aggregator with the provided Manager.
func NewAggregator(manager *Manager) *Aggregator {
	return &Aggregator{manager: manager}
}

// Sum adds multiple Money objects together and returns the result.
// All Money objects must have the same currency.
//
// Preconditions:
//   - All Money arguments must be non-nil (panics if nil)
func (a *Aggregator) Sum(moneys ...*Money) (*Money, error) {
	if a == nil || a.manager == nil {
		return nil, exception.ErrInvalidAggregatorProvider
	}

	if len(moneys) == 0 {
		return nil, exception.ErrNoMoneyProvided
	}

	if len(moneys) == 1 {
		return moneys[0], nil
	}

	// Use the first money as the base and add all others
	result := moneys[0]
	remaining := moneys[1:]

	return a.manager.Add(result, remaining...)
}

// Min returns the Money object with the minimum amount from the given slice.
// All Money objects must have the same currency.
//
// Preconditions:
//   - All Money arguments must be non-nil (panics if nil)
func (a *Aggregator) Min(moneys ...*Money) (*Money, error) {
	if a == nil || a.manager == nil {
		return nil, exception.ErrInvalidAggregatorProvider
	}

	if len(moneys) == 0 {
		return nil, exception.ErrNoMoneyProvided
	}

	if len(moneys) == 1 {
		return moneys[0], nil
	}

	money := moneys[0]
	for _, m := range moneys[1:] {
		if err := money.AssertSameCurrency(m); err != nil {
			return nil, err
		}

		if m.amount < money.amount {
			money = m
		}
	}

	return money, nil
}

// Max returns the Money object with the maximum amount from the given slice.
// All Money objects must have the same currency.
//
// Preconditions:
//   - All Money arguments must be non-nil (panics if nil)
func (a *Aggregator) Max(moneys ...*Money) (*Money, error) {
	if a == nil || a.manager == nil {
		return nil, exception.ErrInvalidAggregatorProvider
	}

	if len(moneys) == 0 {
		return nil, exception.ErrNoMoneyProvided
	}

	if len(moneys) == 1 {
		return moneys[0], nil
	}

	money := moneys[0]
	for _, m := range moneys[1:] {
		if err := money.AssertSameCurrency(m); err != nil {
			return nil, err
		}

		if m.amount > money.amount {
			money = m
		}
	}

	return money, nil
}

// Avg returns the average of multiple Money objects.
// All Money objects must have the same currency.
//
// The result is computed as sum.amount / len(moneys) using integer division,
// which truncates any fractional cents toward zero rather than rounding.
// For example, the average of $1.00 and $2.00 is $1.50, but the average
// of $1.00, $2.00, and $3.00 is $2.00 (not $2.00 rounded from $2.00).
// The average of $0.01, $0.01, and $0.01 remains $0.01.
// The average of $0.02 and $0.01 is $0.01 (truncated from $0.015).
//
// Preconditions:
//   - All Money arguments must be non-nil (panics if nil)
func (a *Aggregator) Avg(moneys ...*Money) (*Money, error) {
	if a == nil || a.manager == nil {
		return nil, exception.ErrInvalidAggregatorProvider
	}

	if len(moneys) == 0 {
		return nil, exception.ErrNoMoneyProvided
	}

	if len(moneys) == 1 {
		return moneys[0], nil
	}

	sum, err := a.Sum(moneys...)
	if err != nil {
		return nil, err
	}

	// Divide by the count
	avgAmount := sum.amount / int64(len(moneys))

	return &Money{
		amount:   avgAmount,
		currency: sum.currency,
	}, nil
}

//var defaultAggregator = NewAggregator(NewManager())
//
//// Sum adds multiple Money objects together and returns the result.
//// All Money objects must have the same currency.
////
//// Preconditions:
////   - All Money arguments must be non-nil (panics if nil)
//func Sum(moneys ...*Money) (*Money, error) {
//	return defaultAggregator.Sum(moneys...)
//}
//
//// Min returns the Money object with the minimum amount from the given slice.
//// All Money objects must have the same currency.
////
//// Preconditions:
////   - All Money arguments must be non-nil (panics if nil)
//func Min(moneys ...*Money) (*Money, error) {
//	return defaultAggregator.Min(moneys...)
//}
//
//// Max returns the Money object with the maximum amount from the given slice.
//// All Money objects must have the same currency.
////
//// Preconditions:
////   - All Money arguments must be non-nil (panics if nil)
//func Max(moneys ...*Money) (*Money, error) {
//	return defaultAggregator.Max(moneys...)
//}
//
//// Avg returns the average of multiple Money objects.
//// All Money objects must have the same currency.
////
//// The result is computed as sum.amount / len(moneys) using integer division,
//// which truncates any fractional cents toward zero rather than rounding.
//// For example, the average of $1.00 and $2.00 is $1.50, but the average
//// of $1.00, $2.00, and $3.00 is $2.00 (not $2.00 rounded from $2.00).
//// The average of $0.01, $0.01, and $0.01 remains $0.01.
//// The average of $0.02 and $0.01 is $0.01 (truncated from $0.015).
////
//// Preconditions:
////   - All Money arguments must be non-nil (panics if nil)
//func Avg(moneys ...*Money) (*Money, error) {
//	return defaultAggregator.Avg(moneys...)
//}
