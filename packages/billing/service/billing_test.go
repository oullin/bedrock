package service_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/service"
)

// stubSubscriptionStore implements billing.SubscriptionStore for testing.
type stubSubscriptionStore struct {
	subs []*billing.Subscription
}

type stubOrderStore struct{}

type stubProductStore struct{}

func (s *stubSubscriptionStore) FindByID(_ context.Context, _ int64) (*billing.Subscription, error) {
	return nil, nil
}

func (s *stubSubscriptionStore) FindByProviderID(_ context.Context, _ string) (*billing.Subscription, error) {
	return nil, nil
}

func (s *stubSubscriptionStore) CurrentForBillable(_ context.Context, _ string, _ int64) (*billing.Subscription, error) {
	if len(s.subs) > 0 {
		return s.subs[0], nil
	}

	return nil, nil
}

func (s *stubSubscriptionStore) ActiveForBillable(_ context.Context, _ string, _ int64) ([]*billing.Subscription, error) {
	return s.subs, nil
}

func (s *stubSubscriptionStore) Create(_ context.Context, _ *billing.Subscription) error { return nil }
func (s *stubSubscriptionStore) Save(_ context.Context, _ *billing.Subscription) error   { return nil }
func (s *stubSubscriptionStore) Delete(_ context.Context, _ int64) error               { return nil }

func (s *stubOrderStore) FindByID(_ context.Context, _ int64) (*billing.Order, error) { return nil, nil }
func (s *stubOrderStore) FindByBillable(_ context.Context, _ int64, _ int) ([]billing.Order, error) {
	return nil, nil
}
func (s *stubOrderStore) Create(_ context.Context, _ *billing.Order) error { return nil }
func (s *stubOrderStore) Save(_ context.Context, _ *billing.Order) error   { return nil }
func (s *stubOrderStore) HasCompletedForProduct(_ context.Context, _ int64, _ int64) (bool, error) {
	return false, nil
}

func (s *stubProductStore) FindByID(_ context.Context, _ int64) (*billing.Product, error) {
	return nil, nil
}
func (s *stubProductStore) Active(_ context.Context) ([]billing.Product, error) { return nil, nil }
func (s *stubProductStore) ActiveSubscriptions(_ context.Context) ([]billing.Product, error) {
	return nil, nil
}
func (s *stubProductStore) ActiveOneTime(_ context.Context) ([]billing.Product, error) { return nil, nil }

// Mirrors BillingServiceTest::test_get_active_subscription_returns_null_when_no_subscription
func TestBillingService_GetActiveSubscription_ReturnsNilWhenNone(t *testing.T) {
	svc := service.NewBillingService(
		&stubSubscriptionStore{subs: nil},
		&stubOrderStore{},
		&stubProductStore{},
	)

	result, err := svc.GetActiveSubscription(context.Background(), "team", 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != nil {
		t.Error("expected nil active subscription, got non-nil")
	}
}

// Mirrors BillingServiceTest::test_team_is_not_subscribed_to_any_provider_by_default
func TestBillingService_IsNotSubscribedByDefault(t *testing.T) {
	svc := service.NewBillingService(
		&stubSubscriptionStore{subs: nil},
		&stubOrderStore{},
		&stubProductStore{},
	)

	subscribed, err := svc.IsSubscribedToAnyProvider(context.Background(), "team", 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if subscribed {
		t.Error("expected not subscribed, got subscribed")
	}
}
