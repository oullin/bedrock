package spark

import "net/http"

// Billable describes an entity that can be billed (e.g. a team or user).
type Billable interface {
	BillableID() int64
	BillableType() string
	BillableName() string
	BillableEmail() string
}

// ProviderConfigurable is implemented by billables that can switch between
// payment providers (e.g. "stripe" or "paddle").
type ProviderConfigurable interface {
	PaymentProvider() string
	SetPaymentProvider(provider string) error
}

// ResolverFunc resolves the current billable from an HTTP request.
type ResolverFunc func(r *http.Request) (Billable, error)

// AuthorizerFunc determines whether a billable is authorised to view
// the billing portal.
type AuthorizerFunc func(billable Billable, r *http.Request) bool

// EligibilityFunc checks whether a billable is eligible for a plan.
// Return a non-nil error to reject the plan.
type EligibilityFunc func(billable Billable, plan Plan) error

// SeatCountFunc returns how many seats a billable currently occupies.
type SeatCountFunc func(billable Billable) int

// BillableResolver wraps a ResolverFunc so it can be used as a
// dependency without bare function types.
type BillableResolver struct {
	Resolve ResolverFunc
}
