package spark

// SubscriptionStatus represents the state of a subscription.
type SubscriptionStatus string

// GrantsAccess reports whether this status allows the billable to use
// the product. Active, trialing, and past-due all grant access.

// SubscriptionInterval describes how often a subscription renews.
type SubscriptionInterval string

// TransactionStatus represents the state of a transaction.
type TransactionStatus string

// ProrationBehavior controls how mid-cycle plan changes are billed.
type ProrationBehavior string

const (
	StatusPending         SubscriptionStatus = "pending"
	StatusAwaitingPayment SubscriptionStatus = "awaiting_payment"
	StatusActive          SubscriptionStatus = "active"
	StatusTrialing        SubscriptionStatus = "trialing"
	StatusPastDue         SubscriptionStatus = "past_due"
	StatusPaused          SubscriptionStatus = "paused"
	StatusCanceled        SubscriptionStatus = "canceled"
	StatusExpired         SubscriptionStatus = "expired"
)

func (s SubscriptionStatus) GrantsAccess() bool {
	switch s {
	case StatusActive, StatusTrialing, StatusPastDue:
		return true
	default:
		return false
	}
}

const (
	IntervalDay   SubscriptionInterval = "day"
	IntervalWeek  SubscriptionInterval = "week"
	IntervalMonth SubscriptionInterval = "month"
	IntervalYear  SubscriptionInterval = "year"
)

const (
	TransactionDraft     TransactionStatus = "draft"
	TransactionReady     TransactionStatus = "ready"
	TransactionBilled    TransactionStatus = "billed"
	TransactionPaid      TransactionStatus = "paid"
	TransactionCompleted TransactionStatus = "completed"
	TransactionCanceled  TransactionStatus = "canceled"
	TransactionPastDue   TransactionStatus = "past_due"
)

const (
	ProrateNextBilling ProrationBehavior = "prorated_next_billing_period"
	FullNextBilling    ProrationBehavior = "full_next_billing_period"
	ProrateImmediately ProrationBehavior = "prorated_immediately"
	FullImmediately    ProrationBehavior = "full_immediately"
	DoNotBill          ProrationBehavior = "do_not_bill"
)

// DefaultSubscriptionType is the type assigned to subscriptions when no
// explicit type is specified.
const DefaultSubscriptionType = "default"
