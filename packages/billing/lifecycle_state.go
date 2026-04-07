package billing

import (
	"context"
	"math"
	"time"
)

// BillingStateResolver builds the billing state snapshot for a billable entity.
type BillingStateResolver struct {
	catalog       *PlanCatalog
	subscriptions *SubscriptionRepository
	urls          URLResolver
	clock         Clock
}

// NewBillingStateResolver creates a BillingStateResolver.
func NewBillingStateResolver(
	catalog *PlanCatalog,
	subscriptions *SubscriptionRepository,
	urls URLResolver,
	clock Clock,
) *BillingStateResolver {
	return &BillingStateResolver{
		catalog:       catalog,
		subscriptions: subscriptions,
		urls:          urls,
		clock:         clock,
	}
}

// Read returns the current billing state for the billable.
func (r *BillingStateResolver) Read(ctx context.Context, billableType string, billableID int64) BillingStateSnapshot {
	sub, err := r.subscriptions.CurrentSubscription(ctx, billableType, billableID)
	if err != nil || sub == nil {
		return EmptyBillingState()
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

	payNow := sub.Status == StatusPending || sub.Status == StatusAwaitingPayment

	subState := SubscriptionStateSnapshot{
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

	return BillingStateSnapshot{Subscription: subState, CTA: cta}
}

func (r *BillingStateResolver) buildCTA(sub *Subscription, planName string, pendingExpiresAt *string) SubscriptionCTASnapshot {
	if sub.Status.IsTerminal() {
		return EmptySubscriptionCTA()
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

	return SubscriptionCTASnapshot{
		Visible:          true,
		Href:             &href,
		Label:            "Manage subscription",
		PlanName:         planName,
		Status:           &status,
		PendingExpiresAt: pendingExpiresAt,
		RemainingDays:    remaining,
	}
}
