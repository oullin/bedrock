package subscription

import (
	"context"
	"math"
	"time"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/catalog"
)

// StateResolver builds the billing state snapshot for a billable entity.
type StateResolver struct {
	catalog       *catalog.PlanCatalog
	subscriptions *Repository
	urls          spark.URLResolver
	clock         spark.Clock
}

// NewStateResolver creates a StateResolver.
func NewStateResolver(
	cat *catalog.PlanCatalog,
	subscriptions *Repository,
	urls spark.URLResolver,
	clock spark.Clock,
) *StateResolver {
	return &StateResolver{
		catalog:       cat,
		subscriptions: subscriptions,
		urls:          urls,
		clock:         clock,
	}
}

// Read returns the current billing state for the billable.
func (r *StateResolver) Read(ctx context.Context, billableType string, billableID int64) spark.BillingStateSnapshot {
	sub, err := r.subscriptions.CurrentSubscription(ctx, billableType, billableID)

	if err != nil || sub == nil {
		return spark.EmptyBillingState()
	}

	planName := r.catalog.MarketingName(ctx, sub.Plan)
	portalURL, _ := r.urls.PortalURL(billableType, "")

	var pendingExpiresAt *string

	if sub.PendingExpiresAt != nil {
		s := sub.PendingExpiresAt.Format(time.RFC3339)
		pendingExpiresAt = &s
	}

	var paymentReadyAt *string

	if sub.PaymentReadyAt != nil {
		s := sub.PaymentReadyAt.Format(time.RFC3339)
		paymentReadyAt = &s
	}

	payNow := sub.Status == spark.StatusPending || sub.Status == spark.StatusAwaitingPayment

	subState := spark.SubscriptionStateSnapshot{
		Status:           string(sub.Status),
		PlanCode:         sub.Plan,
		PlanName:         planName,
		BillingPeriod:    string(sub.BillingPeriod),
		PayNow:           payNow,
		PendingExpiresAt: pendingExpiresAt,
		PaymentReadyAt:   paymentReadyAt,
		PortalURL:        &portalURL,
	}

	cta := r.buildCTA(sub, planName, pendingExpiresAt)

	return spark.BillingStateSnapshot{Subscription: subState, CTA: cta}
}

func (r *StateResolver) buildCTA(sub *spark.Subscription, planName string, pendingExpiresAt *string) spark.SubscriptionCTASnapshot {
	if sub.Status.IsTerminal() {
		return spark.EmptySubscriptionCTA()
	}

	status := string(sub.Status)

	var remaining *int

	if sub.PendingExpiresAt != nil {
		days := int(math.Ceil(time.Until(*sub.PendingExpiresAt).Hours() / 24))

		if days < 0 {
			days = 0
		}

		remaining = &days
	}

	href := r.urls.SubscriptionShowURL(sub.Plan, sub.BillingPeriod)

	return spark.SubscriptionCTASnapshot{
		Visible:          true,
		Href:             &href,
		Label:            "Manage subscription",
		PlanName:         planName,
		Status:           &status,
		PendingExpiresAt: pendingExpiresAt,
		RemainingDays:    remaining,
	}
}
