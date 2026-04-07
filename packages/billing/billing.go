package billing

import "time"

// SubscriptionHoldDays is the number of days a pending subscription is kept
// before it is considered stale and eligible for expiry.
const SubscriptionHoldDays = 14

// DefaultSubscriptionType is the default subscription type used when none is
// specified explicitly.
const DefaultSubscriptionType = "default"

// ProrationBehavior controls how charges are adjusted when a subscription
// plan or quantity changes mid-cycle.
type ProrationBehavior string

const (
	ProratedNextBillingPeriod ProrationBehavior = "prorated_next_billing_period"
	FullNextBillingPeriod     ProrationBehavior = "full_next_billing_period"
	ProratedImmediately       ProrationBehavior = "prorated_immediately"
	FullImmediately           ProrationBehavior = "full_immediately"
	DoNotBill                 ProrationBehavior = "do_not_bill"
)

// Valid reports whether b is a recognised proration behavior.
func (b ProrationBehavior) Valid() bool {
	switch b {
	case ProratedNextBillingPeriod, FullNextBillingPeriod,
		ProratedImmediately, FullImmediately, DoNotBill:
		return true
	}

	return false
}

// SubscriptionInterval represents a billing frequency on the provider side.
type SubscriptionInterval string

const (
	IntervalDay   SubscriptionInterval = "day"
	IntervalWeek  SubscriptionInterval = "week"
	IntervalMonth SubscriptionInterval = "month"
	IntervalYear  SubscriptionInterval = "year"
)

// WebhookMaxTimeDrift is the default maximum acceptable difference between the
// webhook timestamp and the server clock, used for signature verification.
const WebhookMaxTimeDrift = 5 * time.Second
