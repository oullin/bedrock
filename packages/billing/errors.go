package billing

import (
	"errors"
	"fmt"
	"strings"
)

// Sentinel errors used across the package.
var (
	ErrNotFound               = errors.New("billing: resource not found")
	ErrAlreadySubscribed      = errors.New("billing: billable is already subscribed")
	ErrNotSubscribed          = errors.New("billing: billable is not subscribed")
	ErrPlanNotEligible        = errors.New("billing: billable is not eligible for this plan")
	ErrSubscriptionNotOnGrace = errors.New("billing: subscription is not on a grace period")
	ErrSubscriptionTerminal   = errors.New("billing: subscription is in a terminal state")
	ErrNoActivePrice          = errors.New("billing: no active price available for the plan period")
	ErrPendingCheckoutExists  = errors.New("billing: a pending checkout already exists")
	ErrRecoveryFailed         = errors.New("billing: explicit subscription recovery failed")
	ErrUnauthorized           = errors.New("billing: not authorized to view billing portal")
	ErrBillableRequired       = errors.New("billing: billable context is required")
)

// ProviderError represents an error returned by the payment provider.
type ProviderError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface.
func (e *ProviderError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("billing: provider error [%s]: %s", e.Code, e.Message)
	}

	return fmt.Sprintf("billing: provider error: %s", e.Message)
}

// Unwrap returns the underlying error.
func (e *ProviderError) Unwrap() error {
	return e.Err
}

// InvalidPlanPriceError indicates that a plan price is missing required fields.
type InvalidPlanPriceError struct {
	Identifier  string
	PricingMode string
	Missing     []string
}

// Error implements the error interface.
func (e *InvalidPlanPriceError) Error() string {
	return fmt.Sprintf(
		"billing: plan price [%s] with pricing mode [%s] is missing required money fields [%s]",
		e.Identifier,
		e.PricingMode,
		strings.Join(e.Missing, ", "),
	)
}

// NewInvalidPlanPriceError builds an InvalidPlanPriceError for a price with
// incomplete money fields.
func NewInvalidPlanPriceError(identifier, pricingMode string, amountMinor *int64, currency string) *InvalidPlanPriceError {
	var missing []string

	if amountMinor == nil {
		missing = append(missing, "amount_minor")
	}

	if strings.TrimSpace(currency) == "" {
		missing = append(missing, "currency")
	}

	return &InvalidPlanPriceError{
		Identifier:  identifier,
		PricingMode: pricingMode,
		Missing:     missing,
	}
}

// ValidationError holds field-level validation errors.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("billing: validation error on %s: %s", e.Field, e.Message)
}

// ValidationErrors collects multiple field validation failures.
type ValidationErrors []ValidationError

// Error implements the error interface.
func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return "billing: validation failed"
	}

	msgs := make([]string, len(e))

	for i, v := range e {
		msgs[i] = v.Error()
	}

	return strings.Join(msgs, "; ")
}

// HasField reports whether the collection contains an error for the given field.
func (e ValidationErrors) HasField(field string) bool {
	for _, v := range e {
		if v.Field == field {
			return true
		}
	}

	return false
}
