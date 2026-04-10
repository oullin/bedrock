package spark

// SubscriptionStateSnapshot holds the resolved state of a subscription for
// internal consumption and serialisation to the frontend.
type SubscriptionStateSnapshot struct {
	Status           string
	PlanCode         string
	PlanName         string
	BillingPeriod    string
	PayNow           bool
	PendingExpiresAt *string
	PaymentReadyAt   *string
	PortalURL        *string
}

// EmptySubscriptionState returns a zero-value snapshot.

// SubscriptionCTASnapshot holds the call-to-action state for the subscription
// banner/widget.
type SubscriptionCTASnapshot struct {
	Visible          bool
	Href             *string
	Label            string
	PlanName         string
	Status           *string
	PendingExpiresAt *string
	RemainingDays    *int
}

// EmptySubscriptionCTA returns a zero-value CTA snapshot.

// BillingStateSnapshot combines the subscription state and CTA.
type BillingStateSnapshot struct {
	Subscription SubscriptionStateSnapshot
	CTA          SubscriptionCTASnapshot
}

// EmptyBillingState returns a zero-value billing state.

// TransactionSnapshot holds a serialisable view of a transaction.
type TransactionSnapshot struct {
	ID            int64
	TotalMinor    *int64
	TaxMinor      *int64
	Currency      *string
	DisplayTotal  *string
	DisplayTax    *string
	BilledAt      string
	InvoiceNumber string
	Status        string
}

// EmptyTransactionSnapshot returns a zero-value transaction snapshot.

// PlanConfig holds the provider price IDs for a plan's monthly and yearly
// variants.
type PlanConfig struct {
	MonthlyID string
	YearlyID  string
}

func EmptySubscriptionState() SubscriptionStateSnapshot {
	return SubscriptionStateSnapshot{}
}

func EmptySubscriptionCTA() SubscriptionCTASnapshot {
	return SubscriptionCTASnapshot{}
}

func EmptyBillingState() BillingStateSnapshot {
	return BillingStateSnapshot{
		Subscription: EmptySubscriptionState(),
		CTA:          EmptySubscriptionCTA(),
	}
}

func EmptyTransactionSnapshot() TransactionSnapshot {
	return TransactionSnapshot{}
}

// HasConfiguredProviderIDs reports whether at least one provider price ID is set.
func (c PlanConfig) HasConfiguredProviderIDs() bool {
	return c.MonthlyID != "" || c.YearlyID != ""
}

// MatchesPriceID reports whether the given price ID matches either the
// monthly or yearly provider price ID.
func (c PlanConfig) MatchesPriceID(priceID string) bool {
	if priceID == "" {
		return false
	}

	return priceID == c.MonthlyID || priceID == c.YearlyID
}
