// Package service provides the multi-provider billing orchestration layer.
package service

import (
	"context"
	"time"

	"github.com/bedrock/packages/spark"
)

// ProviderCheckout represents a checkout session from any provider.
type ProviderCheckout struct {
	Provider     string         // "stripe" or "paddle"
	CheckoutURL  string         // redirect URL for Stripe, or empty for Paddle inline
	CheckoutData map[string]any // Paddle inline checkout data
}

// ActiveSubscription contains the provider and subscription details
// for a currently active subscription.
type ActiveSubscription struct {
	Provider     string
	Subscription *spark.Subscription
}

// BillingService orchestrates billing operations across multiple
// payment providers. php.
type BillingService struct {
	subscriptions spark.SubscriptionStore
	orders        spark.OrderStore
	products      spark.ProductStore
}

// NewBillingService creates a BillingService with the required stores.
func NewBillingService(
	subs spark.SubscriptionStore,
	orders spark.OrderStore,
	products spark.ProductStore,
) *BillingService {
	return &BillingService{
		subscriptions: subs,
		orders:        orders,
		products:      products,
	}
}

// GetActiveSubscription checks all providers for an active subscription.
// Returns nil if the billable has no active subscription.
func (s *BillingService) GetActiveSubscription(ctx context.Context, billableType string, billableID int64) (*ActiveSubscription, error) {
	subs, err := s.subscriptions.ActiveForBillable(ctx, billableType, billableID)

	if err != nil {
		return nil, err
	}

	for _, sub := range subs {
		if sub.Valid() {
			return &ActiveSubscription{
				Provider:     providerForSubscription(sub),
				Subscription: sub,
			}, nil
		}
	}

	return nil, nil
}

func providerForSubscription(sub *spark.Subscription) string {
	if sub.PaddleID != "" {
		return "paddle"
	}

	return "stripe"
}

// IsSubscribedToAnyProvider reports whether the billable has any valid
// subscription regardless of provider.
func (s *BillingService) IsSubscribedToAnyProvider(ctx context.Context, billableType string, billableID int64) (bool, error) {
	active, err := s.GetActiveSubscription(ctx, billableType, billableID)

	if err != nil {
		return false, err
	}

	return active != nil, nil
}

// CancelSubscription cancels the billable's active subscription.
func (s *BillingService) CancelSubscription(ctx context.Context, billableType string, billableID int64) error {
	active, err := s.GetActiveSubscription(ctx, billableType, billableID)

	if err != nil {
		return err
	}

	if active == nil {
		return spark.ErrNotSubscribed
	}

	now := time.Now()

	if !active.Subscription.Cancel(now) {
		return nil
	}

	return s.subscriptions.Save(ctx, active.Subscription)
}

// ResumeSubscription resumes a canceled subscription that is still on
// its grace period.
func (s *BillingService) ResumeSubscription(ctx context.Context, billableType string, billableID int64) error {
	active, err := s.GetActiveSubscription(ctx, billableType, billableID)

	if err != nil {
		return err
	}

	if active == nil {
		return spark.ErrNotSubscribed
	}

	if !active.Subscription.Resume(time.Now()) {
		return spark.ErrNotSubscribed
	}

	return s.subscriptions.Save(ctx, active.Subscription)
}

// CreateOneTimeOrder creates an Order record for a one-time purchase.
func (s *BillingService) CreateOneTimeOrder(ctx context.Context, teamID int64, product *spark.Product, provider string) (*spark.Order, error) {
	now := time.Now()
	order := &spark.Order{
		TeamID:          teamID,
		ProductID:       product.ID,
		PaymentProvider: provider,
		Amount:          product.PriceAmount,
		Currency:        product.Currency,
		Status:          spark.OrderStatusPending,
		Metadata:        map[string]any{"price_id": product.PriceIDForProvider(provider)},
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.orders.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}
