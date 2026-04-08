package spark

import "net/http"

// Billable is implemented by any model that can hold subscriptions and be
// billed through a payment provider (e.g. a Team).
type Billable interface {
	BillableID() int64
	BillableType() string
	BillableUUID() string
	BillableName() string
	BillableEmail() string
}

// BillableResolver resolves the current billable entity from an HTTP request.
type BillableResolver interface {
	Resolve(r *http.Request) (Billable, error)
}

// ResolverFunc is a function that resolves a Billable from an HTTP request.
type ResolverFunc func(r *http.Request) (Billable, error)

// AuthorizerFunc is a function that checks whether a billable is authorized
// to view the billing portal.
type AuthorizerFunc func(billable Billable, r *http.Request) bool

// EligibilityFunc is a function that checks whether a billable is eligible
// for a specific plan.
type EligibilityFunc func(billable Billable, plan SparkPlan) error

// SeatCountFunc returns the current seat count for a billable entity.
type SeatCountFunc func(billable Billable) int
