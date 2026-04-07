package spark

import (
	"context"
	"time"
)

// SubscriptionRepository manages the creation and retrieval of subscriptions.
type SubscriptionRepository struct {
	subscriptions SubscriptionStore
	catalog       *PlanCatalog
	checkout      ProviderCheckoutGenerator
	customers     CustomerStore
	clock         Clock
}

// NewSubscriptionRepository creates a SubscriptionRepository.
func NewSubscriptionRepository(
	subscriptions SubscriptionStore,
	catalog *PlanCatalog,
	checkout ProviderCheckoutGenerator,
	customers CustomerStore,
	clock Clock,
) *SubscriptionRepository {
	return &SubscriptionRepository{
		subscriptions: subscriptions,
		catalog:       catalog,
		checkout:      checkout,
		customers:     customers,
		clock:         clock,
	}
}

// StartPendingSubscription creates a new pending subscription for the billable.
func (r *SubscriptionRepository) StartPendingSubscription(
	ctx context.Context,
	billable Billable,
	plan SubscriptionPlan,
	period BillingPeriod,
) (*Subscription, error) {
	existing, _ := r.subscriptions.LatestForPlan(ctx, billable.BillableType(), billable.BillableID(), string(plan))
	if existing != nil && !existing.Status.IsTerminal() {
		return existing, nil
	}

	expiresAt := r.clock.Now().Add(time.Duration(SubscriptionHoldDays) * 24 * time.Hour)

	sub := &Subscription{
		BillableType:     billable.BillableType(),
		BillableID:       billable.BillableID(),
		Type:             DefaultSubscriptionType,
		Status:           StatusPending,
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
func (r *SubscriptionRepository) StartStarterTrial(ctx context.Context, billable Billable) (*Subscription, error) {
	exists, _ := r.subscriptions.ExistsAccessibleNonStarter(ctx, billable.BillableType(), billable.BillableID())
	if exists {
		return nil, ErrAlreadySubscribed
	}

	existing, _ := r.subscriptions.LatestForPlan(ctx, billable.BillableType(), billable.BillableID(), string(PlanStarter))
	if existing != nil && !existing.Status.IsTerminal() {
		return existing, nil
	}

	trialEnd := r.clock.Now().Add(time.Duration(SubscriptionHoldDays) * 24 * time.Hour)

	sub := &Subscription{
		BillableType:  billable.BillableType(),
		BillableID:    billable.BillableID(),
		Type:          DefaultSubscriptionType,
		Status:        StatusTrialing,
		Plan:          string(PlanStarter),
		BillingPeriod: PeriodFree,
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
func (r *SubscriptionRepository) BeginCheckout(
	ctx context.Context,
	billable Billable,
	plan SubscriptionPlan,
	period BillingPeriod,
	price *PlanPeriodPrice,
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

	session, err := r.checkout.Generate(ctx, customer, []CheckoutItem{
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
func (r *SubscriptionRepository) CurrentSubscription(ctx context.Context, billableType string, billableID int64) (*Subscription, error) {
	return r.subscriptions.CurrentForBillable(ctx, billableType, billableID)
}

// PendingExpiryDays returns the number of days a pending subscription is held.
func (r *SubscriptionRepository) PendingExpiryDays() int {
	return SubscriptionHoldDays
}
