package spark

import (
	"context"
	"time"
)

// Clock abstracts time for testing.
type Clock interface {
	Now() time.Time
}

// SystemClock returns the real wall clock.
type SystemClock struct{}

// Now returns the current time.
func (SystemClock) Now() time.Time { return time.Now() }

// CurrencyFormatter formats monetary amounts for display.
type CurrencyFormatter interface {
	FormatAmount(amountMinor int64, currency string, locale string) string
}

// TransactionManager wraps database transactions.
type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// EventDispatcher publishes domain events.
type EventDispatcher interface {
	Dispatch(ctx context.Context, event any) error
}

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
