package spark

import (
	"context"
	"time"
)

// ReconcileSubscriptionsJob is a background job that reconciles pending
// subscriptions against the payment provider.
type ReconcileSubscriptionsJob struct {
	billing *BillingWorkflow
}

// NewReconcileSubscriptionsJob creates the job.
func NewReconcileSubscriptionsJob(billing *BillingWorkflow) *ReconcileSubscriptionsJob {
	return &ReconcileSubscriptionsJob{billing: billing}
}

// Handle executes the job.
func (j *ReconcileSubscriptionsJob) Handle(ctx context.Context) error {
	return j.billing.ReconcilePendingSubscriptions(ctx)
}

// MaxRetries returns the maximum number of retry attempts.
func (j *ReconcileSubscriptionsJob) MaxRetries() int { return 5 }

// Timeout returns the maximum execution time.
func (j *ReconcileSubscriptionsJob) Timeout() time.Duration { return 120 * time.Second }

// Backoff returns the delay between retries.
func (j *ReconcileSubscriptionsJob) Backoff() []time.Duration {
	return []time.Duration{60 * time.Second, 300 * time.Second, 900 * time.Second}
}

// Ensure ReconcileSubscriptionsJob implements Job.
var _ Job = (*ReconcileSubscriptionsJob)(nil)
