package spark

import (
	"context"
	"time"
)

// ---------------------------------------------------------------------------
// Clock
// ---------------------------------------------------------------------------

// Clock abstracts time for testing.
type Clock interface {
	Now() time.Time
}

// SystemClock returns the real wall clock.
type SystemClock struct{}

// Now returns the current time.
func (SystemClock) Now() time.Time { return time.Now() }

// ---------------------------------------------------------------------------
// Currency formatting
// ---------------------------------------------------------------------------

// CurrencyFormatter formats monetary amounts for display.
type CurrencyFormatter interface {
	FormatAmount(amountMinor int64, currency string, locale string) string
}

// ---------------------------------------------------------------------------
// Storage interfaces
// ---------------------------------------------------------------------------

// CustomerStore abstracts persistence of customer records.
type CustomerStore interface {
	FindByProviderID(ctx context.Context, providerID string) (*Customer, error)
	FindByBillable(ctx context.Context, billableType string, billableID int64) (*Customer, error)
	Create(ctx context.Context, c *Customer) error
	Save(ctx context.Context, c *Customer) error
}

// SubscriptionStore abstracts persistence of subscription records.
type SubscriptionStore interface {
	FindByID(ctx context.Context, id int64) (*Subscription, error)
	FindByProviderID(ctx context.Context, providerID string) (*Subscription, error)
	CurrentForBillable(ctx context.Context, billableType string, billableID int64) (*Subscription, error)
	LatestForPlan(ctx context.Context, billableType string, billableID int64, plan string) (*Subscription, error)
	Create(ctx context.Context, s *Subscription) error
	Save(ctx context.Context, s *Subscription) error
	Delete(ctx context.Context, id int64) error
	FindExpirable(ctx context.Context, before time.Time) ([]*Subscription, error)
	FindPendingWithoutProvider(ctx context.Context) ([]*Subscription, error)
	ExistsAccessibleNonStarter(ctx context.Context, billableType string, billableID int64) (bool, error)
	ActiveForBillable(ctx context.Context, billableType string, billableID int64) ([]*Subscription, error)
}

// SubscriptionItemStore abstracts persistence of subscription items.
type SubscriptionItemStore interface {
	FindBySubscription(ctx context.Context, subscriptionID int64) ([]SubscriptionItem, error)
	Sync(ctx context.Context, subscriptionID int64, items []SubscriptionItem) error
}

// TransactionStore abstracts persistence of transaction records.
type TransactionStore interface {
	FindByProviderID(ctx context.Context, providerID string) (*Transaction, error)
	FindByBillable(ctx context.Context, billableType string, billableID int64, limit int) ([]Transaction, error)
	Create(ctx context.Context, t *Transaction) error
	Save(ctx context.Context, t *Transaction) error
}

// PlanStore abstracts persistence of the plan catalog.
type PlanStore interface {
	ActivePlans(ctx context.Context) ([]Plan, error)
	FindBySlug(ctx context.Context, slug string) (*Plan, error)
	FindByCode(ctx context.Context, code string) (*Plan, error)
	UpsertPlans(ctx context.Context, plans []Plan) error
	UpsertFeatures(ctx context.Context, features []Feature) error
	SyncPlanFeatures(ctx context.Context, planID int64, features []PlanFeature) error
	FindPlanPeriod(ctx context.Context, planID int64, period BillingPeriod) (*PlanPeriod, error)
	UpsertPlanPeriods(ctx context.Context, periods []PlanPeriod) error
}

// PriceStore abstracts persistence of plan period prices.
type PriceStore interface {
	FindByProviderPriceID(ctx context.Context, providerPriceID string) (*PlanPeriodPrice, error)
	ActiveForPeriod(ctx context.Context, planPeriodID int64) (*PlanPeriodPrice, error)
	Create(ctx context.Context, p *PlanPeriodPrice) error
	Save(ctx context.Context, p *PlanPeriodPrice) error
	DeactivateForPeriod(ctx context.Context, planPeriodID int64) error
}

// SubscriptionFeatureStore abstracts persistence of subscription features.
type SubscriptionFeatureStore interface {
	FeaturesForSubscription(ctx context.Context, subscriptionID int64) ([]SubscriptionFeature, error)
	Sync(ctx context.Context, subscriptionID int64, features []SubscriptionFeature) error
	SetActivationState(ctx context.Context, subscriptionID int64, active bool, at time.Time) error
	DeleteAll(ctx context.Context, subscriptionID int64) error
}

// ---------------------------------------------------------------------------
// Provider interfaces (payment gateway abstraction)
// ---------------------------------------------------------------------------

// ProviderClient is a low-level HTTP client for the payment provider API.
type ProviderClient interface {
	Call(ctx context.Context, method, uri string, payload any) ([]byte, error)
}

// ProviderSubscriptionManager manages subscriptions on the provider side.
type ProviderSubscriptionManager interface {
	Fetch(ctx context.Context, providerSubID string) (*ProviderSubscription, error)
	Cancel(ctx context.Context, providerSubID string, cancelNow bool) error
	Pause(ctx context.Context, providerSubID string, pauseNow bool, until *time.Time) error
	Resume(ctx context.Context, providerSubID string, resumeAt *time.Time) error
	StopCancellation(ctx context.Context, providerSubID string) error
	Swap(ctx context.Context, providerSubID string, items []CheckoutItem, proration ProrationBehavior) error
	UpdateQuantity(ctx context.Context, providerSubID string, priceID string, quantity int, proration ProrationBehavior) error
	PaymentMethodUpdateTransaction(ctx context.Context, providerSubID string) (string, error)
}

// ProviderSubscription is the provider-neutral representation of a
// subscription fetched from the payment provider.
type ProviderSubscription struct {
	ID     string
	Status string
	Items  []ProviderSubscriptionItem
}

// ProviderSubscriptionItem is a line item from a provider subscription.
type ProviderSubscriptionItem struct {
	ProductID string
	PriceID   string
	Status    string
	Quantity  int
}

// ProviderCheckoutGenerator creates checkout sessions on the provider.
type ProviderCheckoutGenerator interface {
	Generate(ctx context.Context, customer *Customer, items []CheckoutItem, opts map[string]any) (*CheckoutSession, error)
}

// ProviderPricePreviewer fetches price previews (with tax) from the provider.
type ProviderPricePreviewer interface {
	PreviewPrices(ctx context.Context, items []CheckoutItem, opts map[string]any) ([]PricePreview, error)
}

// ProviderTransactionManager handles transaction operations on the provider.
type ProviderTransactionManager interface {
	InvoicePDFURL(ctx context.Context, providerTransactionID string) (string, error)
	Refund(ctx context.Context, providerTransactionID string, reason string, prices []string) error
}

// ProviderSubscriptionFetcher fetches subscriptions from the provider for
// reconciliation purposes.
type ProviderSubscriptionFetcher interface {
	FetchForBillable(ctx context.Context, billable Billable) (*ProviderSubscription, error)
	FetchExplicit(ctx context.Context, providerSubscriptionID string) (*ProviderSubscription, error)
}

// ---------------------------------------------------------------------------
// Domain interfaces
// ---------------------------------------------------------------------------

// TransactionManager wraps database transactions.
type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// EventDispatcher publishes domain events.
type EventDispatcher interface {
	Dispatch(ctx context.Context, event any) error
}

// ---------------------------------------------------------------------------
// Application layer interfaces
// ---------------------------------------------------------------------------

// Mailer sends email messages.
type Mailer interface {
	Send(ctx context.Context, to string, msg MailMessage) error
}

// MailMessage describes an email to be sent.
type MailMessage interface {
	Subject() string
	ReplyTo() string
	Body() string
}

// Logger provides structured logging for commands and services.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// JobDispatcher dispatches background jobs.
type JobDispatcher interface {
	Dispatch(ctx context.Context, job Job) error
	DispatchSync(ctx context.Context, job Job) error
}

// Job represents an asynchronous unit of work.
type Job interface {
	Handle(ctx context.Context) error
	MaxRetries() int
	Timeout() time.Duration
	Backoff() []time.Duration
}

// URLResolver generates application URLs.
type URLResolver interface {
	SubscriptionShowURL(plan string, period BillingPeriod) string
	PortalURL(billableType string, billableUUID string) (string, bool)
	BillingPortalRedirect(billableType string) string
}
