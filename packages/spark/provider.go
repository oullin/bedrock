package spark

import (
	"context"
	"io"
)

// CustomerCreateOptions carries provider metadata for customer creation.
type CustomerCreateOptions struct {
	TrialEndsAt any
}

// PaymentMethodUpdateTransaction describes a provider transaction used to
// collect updated payment method details.
type PaymentMethodUpdateTransaction struct {
	ID   string
	Data map[string]any
}

// InvoiceDownload contains the provider invoice stream and response metadata.
type InvoiceDownload struct {
	FileName    string
	ContentType string
	Body        io.ReadCloser
}

// ProviderOperations defines the Cashier/Paddle behavior Spark delegates to
// provider SDKs in Laravel.
type ProviderOperations interface {
	CreateCustomer(ctx context.Context, billable Billable, options CustomerCreateOptions) (*Customer, error)
	PreviewPrices(ctx context.Context, priceIDs []string, options map[string]any) ([]PricePreview, error)
	CreateCheckoutSession(ctx context.Context, billable Billable, checkout *Checkout, options map[string]any) (*Checkout, error)
	CreatePaymentMethodUpdateTransaction(ctx context.Context, billable Billable, subscription *Subscription, options map[string]any) (*PaymentMethodUpdateTransaction, error)
	DownloadInvoice(ctx context.Context, transaction *Transaction) (*InvoiceDownload, error)
	LastPayment(ctx context.Context, subscription *Subscription) (*Payment, error)
	NextPayment(ctx context.Context, subscription *Subscription) (*Payment, error)
	UpdateSubscriptionQuantity(ctx context.Context, subscription *Subscription, quantity int, behavior ProrationBehavior) error
}

// PaymentMethodUpdater is the narrow provider contract needed by Spark's
// payment-method update route.
type PaymentMethodUpdater interface {
	CreatePaymentMethodUpdateTransaction(ctx context.Context, billable Billable, subscription *Subscription, options map[string]any) (*PaymentMethodUpdateTransaction, error)
}

// InvoiceDownloader is the narrow provider contract needed by invoice download
// routes.
type InvoiceDownloader interface {
	DownloadInvoice(ctx context.Context, transaction *Transaction) (*InvoiceDownload, error)
}

// PaymentReporter provides payment summaries for billing portal state.
type PaymentReporter interface {
	LastPayment(ctx context.Context, subscription *Subscription) (*Payment, error)
	NextPayment(ctx context.Context, subscription *Subscription) (*Payment, error)
}

// QuantityUpdater applies subscription seat quantity mutations.
type QuantityUpdater interface {
	UpdateSubscriptionQuantity(ctx context.Context, subscription *Subscription, quantity int, behavior ProrationBehavior) error
}
