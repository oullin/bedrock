package billing

import "time"

// SubscriptionPlan identifies a plan tier.
type SubscriptionPlan string

const (
	PlanStarter    SubscriptionPlan = "starter"
	PlanPro        SubscriptionPlan = "pro"
	PlanEnterprise SubscriptionPlan = "enterprise"
)

// Valid reports whether p is a recognised plan.
func (p SubscriptionPlan) Valid() bool {
	switch p {
	case PlanStarter, PlanPro, PlanEnterprise:
		return true
	}

	return false
}

// SubscriptionStatus represents the lifecycle state of a subscription.
type SubscriptionStatus string

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

// Valid reports whether s is a recognised status.
func (s SubscriptionStatus) Valid() bool {
	switch s {
	case StatusPending, StatusAwaitingPayment, StatusActive,
		StatusTrialing, StatusPastDue, StatusPaused,
		StatusCanceled, StatusExpired:
		return true
	}

	return false
}

// GrantsAccess reports whether a subscription in this status should grant
// access to gated features.
func (s SubscriptionStatus) GrantsAccess() bool {
	switch s {
	case StatusActive, StatusTrialing, StatusPastDue:
		return true
	}

	return false
}

// IsTerminal reports whether the subscription has reached an end state.
func (s SubscriptionStatus) IsTerminal() bool {
	return s == StatusCanceled || s == StatusExpired
}

// BillingPeriod represents the billing cycle length.
type BillingPeriod string

const (
	PeriodMonthly BillingPeriod = "monthly"
	PeriodYearly  BillingPeriod = "yearly"
	PeriodFree    BillingPeriod = "free"
	PeriodCustom  BillingPeriod = "custom"
)

// Valid reports whether p is a recognised billing period.
func (p BillingPeriod) Valid() bool {
	switch p {
	case PeriodMonthly, PeriodYearly, PeriodFree, PeriodCustom:
		return true
	}

	return false
}

// DisplayLabel returns a human-readable label for the billing period.
func (p BillingPeriod) DisplayLabel() string {
	switch p {
	case PeriodMonthly:
		return "Billed monthly"
	case PeriodYearly:
		return "Billed annually"
	case PeriodFree:
		return "Free"
	case PeriodCustom:
		return "Custom"
	default:
		return string(p)
	}
}

// SubscriptionFeatureCode identifies an entitlement granted by a plan.
type SubscriptionFeatureCode string

const (
	FeatureActiveRoleCapacity   SubscriptionFeatureCode = "active_role_capacity"
	FeatureBookingDispatch      SubscriptionFeatureCode = "booking_dispatch"
	FeatureCalendarIntegrations SubscriptionFeatureCode = "calendar_integrations"
	FeatureViewLLMAnalysis      SubscriptionFeatureCode = "view_llm_analysis"
)

// Valid reports whether c is a recognised feature code.
func (c SubscriptionFeatureCode) Valid() bool {
	switch c {
	case FeatureActiveRoleCapacity, FeatureBookingDispatch,
		FeatureCalendarIntegrations, FeatureViewLLMAnalysis:
		return true
	}

	return false
}

// PlanPricingMode describes how a plan period price is calculated.
type PlanPricingMode string

const (
	PricingModeMoney  PlanPricingMode = "money"
	PricingModeFree   PlanPricingMode = "free"
	PricingModeCustom PlanPricingMode = "custom"
)

// Valid reports whether m is a recognised pricing mode.
func (m PlanPricingMode) Valid() bool {
	switch m {
	case PricingModeMoney, PricingModeFree, PricingModeCustom:
		return true
	}

	return false
}

// DisplayLabel returns a human-readable label for the pricing mode.
func (m PlanPricingMode) DisplayLabel() string {
	switch m {
	case PricingModeMoney:
		return "Money"
	case PricingModeFree:
		return "Free"
	case PricingModeCustom:
		return "Custom"
	default:
		return string(m)
	}
}

// TransactionStatus represents the state of a payment transaction.
type TransactionStatus string

const (
	TransactionDraft     TransactionStatus = "draft"
	TransactionReady     TransactionStatus = "ready"
	TransactionBilled    TransactionStatus = "billed"
	TransactionPaid      TransactionStatus = "paid"
	TransactionCompleted TransactionStatus = "completed"
	TransactionCanceled  TransactionStatus = "canceled"
	TransactionPastDue   TransactionStatus = "past_due"
)

// Valid reports whether s is a recognised transaction status.
func (s TransactionStatus) Valid() bool {
	switch s {
	case TransactionDraft, TransactionReady, TransactionBilled,
		TransactionPaid, TransactionCompleted, TransactionCanceled,
		TransactionPastDue:
		return true
	}

	return false
}

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
