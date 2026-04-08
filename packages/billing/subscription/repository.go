package subscription

import (
	"context"
	"time"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/catalog"
)

// Repository manages the creation and retrieval of subscriptions.
type Repository struct {
	subscriptions billing.SubscriptionStore
	catalog       *catalog.PlanCatalog
	checkout      billing.ProviderCheckoutGenerator
	customers     billing.CustomerStore
	clock         billing.Clock
}

// NewRepository creates a Repository.
func NewRepository(
	subscriptions billing.SubscriptionStore,
	cat *catalog.PlanCatalog,
	checkout billing.ProviderCheckoutGenerator,
	customers billing.CustomerStore,
	clock billing.Clock,
) *Repository {
	return &Repository{
		subscriptions: subscriptions,
		catalog:       cat,
		checkout:      checkout,
		customers:     customers,
		clock:         clock,
	}
}

// StartPendingSubscription creates a new pending subscription for the billable.
func (r *Repository) StartPendingSubscription(
	ctx context.Context,
	billable billing.Billable,
	plan billing.SubscriptionPlan,
	period billing.BillingPeriod,
) (*billing.Subscription, error) {
	existing, _ := r.subscriptions.LatestForPlan(ctx, billable.BillableType(), billable.BillableID(), string(plan))
	if existing != nil && !existing.Status.IsTerminal() {
		return existing, nil
	}

	expiresAt := r.clock.Now().Add(time.Duration(billing.SubscriptionHoldDays) * 24 * time.Hour)

	sub := &billing.Subscription{
		BillableType:     billable.BillableType(),
		BillableID:       billable.BillableID(),
		Type:             billing.DefaultSubscriptionType,
		Status:           billing.StatusPending,
		Plan:             string(plan),
		BillingPeriod:    period,
		PendingExpiresAt: &expiresAt,
		CreatedAt:        r.clock.Now(),
		UpdatedAt:        r.clock.Now(),
	}

	if err := r.subscriptions.Create(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// StartStarterTrial creates a trialing subscription for the starter plan.
func (r *Repository) StartStarterTrial(ctx context.Context, billable billing.Billable) (*billing.Subscription, error) {
	exists, _ := r.subscriptions.ExistsAccessibleNonStarter(ctx, billable.BillableType(), billable.BillableID())
	if exists {
		return nil, billing.ErrAlreadySubscribed
	}

	existing, _ := r.subscriptions.LatestForPlan(ctx, billable.BillableType(), billable.BillableID(), string(billing.PlanStarter))
	if existing != nil && !existing.Status.IsTerminal() {
		return existing, nil
	}

	trialEnd := r.clock.Now().Add(time.Duration(billing.SubscriptionHoldDays) * 24 * time.Hour)

	sub := &billing.Subscription{
		BillableType:  billable.BillableType(),
		BillableID:    billable.BillableID(),
		Type:          billing.DefaultSubscriptionType,
		Status:        billing.StatusTrialing,
		Plan:          string(billing.PlanStarter),
		BillingPeriod: billing.PeriodFree,
		TrialEndsAt:   &trialEnd,
		CreatedAt:     r.clock.Now(),
		UpdatedAt:     r.clock.Now(),
	}

	if err := r.subscriptions.Create(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// BeginCheckout creates a pending subscription and generates a checkout URL.
func (r *Repository) BeginCheckout(
	ctx context.Context,
	billable billing.Billable,
	plan billing.SubscriptionPlan,
	period billing.BillingPeriod,
	price *billing.PlanPeriodPrice,
) (string, error) {
	sub, err := r.StartPendingSubscription(ctx, billable, plan, period)
	if err != nil {
		return "", err
	}

	_ = sub // The subscription is tracked locally for reconciliation.

	customer, err := r.customers.FindByBillable(ctx, billable.BillableType(), billable.BillableID())
	if err != nil {
		return "", err
	}

	session, err := r.checkout.Generate(ctx, customer, []billing.CheckoutItem{
		{PriceID: price.ProviderPriceID, Quantity: 1},
	}, nil)
	if err != nil {
		return "", err
	}

	if session.ReturnURL != "" {
		return session.ReturnURL, nil
	}

	return "", nil
}

// CurrentSubscription returns the latest subscription for the billable.
func (r *Repository) CurrentSubscription(ctx context.Context, billableType string, billableID int64) (*billing.Subscription, error) {
	return r.subscriptions.CurrentForBillable(ctx, billableType, billableID)
}

// PendingExpiryDays returns the number of days a pending subscription is held.
func (r *Repository) PendingExpiryDays() int {
	return billing.SubscriptionHoldDays
}
