package billing

import (
	"context"
	"time"
)

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
