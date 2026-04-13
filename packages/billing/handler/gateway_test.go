package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/handler"
	"github.com/bedrock/packages/billing/service"
)

// --- test doubles ---

type stubBillable struct {
	id       int64
	btype    string
	name     string
	email    string
	provider string
}

type testSubStore struct {
	subs []*billing.Subscription
}

type testOrderStore struct{}

type testProductStore struct {
	subProducts     []billing.Product
	oneTimeProducts []billing.Product
}

func (b *stubBillable) BillableID() int64       { return b.id }
func (b *stubBillable) BillableType() string    { return b.btype }
func (b *stubBillable) BillableName() string    { return b.name }
func (b *stubBillable) BillableEmail() string   { return b.email }
func (b *stubBillable) PaymentProvider() string { return b.provider }
func (b *stubBillable) SetPaymentProvider(p string) error {
	b.provider = p

	return nil
}

func (s *testSubStore) FindByID(_ context.Context, _ int64) (*billing.Subscription, error) {
	return nil, nil
}
func (s *testSubStore) FindByProviderID(_ context.Context, _ string) (*billing.Subscription, error) {
	return nil, nil
}
func (s *testSubStore) CurrentForBillable(_ context.Context, _ string, _ int64) (*billing.Subscription, error) {
	return nil, nil
}
func (s *testSubStore) ActiveForBillable(_ context.Context, _ string, _ int64) ([]*billing.Subscription, error) {
	return s.subs, nil
}
func (s *testSubStore) Create(_ context.Context, _ *billing.Subscription) error { return nil }
func (s *testSubStore) Save(_ context.Context, _ *billing.Subscription) error   { return nil }
func (s *testSubStore) Delete(_ context.Context, _ int64) error               { return nil }

func (s *testOrderStore) FindByID(_ context.Context, _ int64) (*billing.Order, error) { return nil, nil }
func (s *testOrderStore) FindByBillable(_ context.Context, _ int64, _ int) ([]billing.Order, error) {
	return nil, nil
}
func (s *testOrderStore) Create(_ context.Context, _ *billing.Order) error { return nil }
func (s *testOrderStore) Save(_ context.Context, _ *billing.Order) error   { return nil }
func (s *testOrderStore) HasCompletedForProduct(_ context.Context, _ int64, _ int64) (bool, error) {
	return false, nil
}

func (s *testProductStore) FindByID(_ context.Context, _ int64) (*billing.Product, error) {
	return nil, nil
}
func (s *testProductStore) Active(_ context.Context) ([]billing.Product, error) {
	return append(s.subProducts, s.oneTimeProducts...), nil
}
func (s *testProductStore) ActiveSubscriptions(_ context.Context) ([]billing.Product, error) {
	return s.subProducts, nil
}
func (s *testProductStore) ActiveOneTime(_ context.Context) ([]billing.Product, error) {
	return s.oneTimeProducts, nil
}

// --- helper ---

func newGateway(resolver billing.ResolverFunc, products *testProductStore) *handler.GatewayHandler {
	mgr := billing.NewManager()
	billing := service.NewBillingService(&testSubStore{}, &testOrderStore{}, products)

	return handler.NewGatewayHandler(mgr, billing, products, resolver)
}

// --- tests ---

// Mirrors BillingGatewayTest::test_guest_cannot_access_billing_gateway
func TestGatewayHandler_Show_UnauthenticatedReturns401(t *testing.T) {
	resolver := func(r *http.Request) (billing.Billable, error) {
		return nil, errors.New("unauthenticated")
	}

	h := newGateway(resolver, &testProductStore{})

	req := httptest.NewRequest("GET", "/billing/choose-provider", nil)
	rec := httptest.NewRecorder()
	h.Show(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// Mirrors BillingGatewayTest::test_authenticated_user_can_access_billing_gateway
func TestGatewayHandler_Show_AuthenticatedReturnsOK(t *testing.T) {
	billable := &stubBillable{id: 1, btype: "team", name: "Test"}
	resolver := func(r *http.Request) (billing.Billable, error) {
		return billable, nil
	}

	h := newGateway(resolver, &testProductStore{})

	req := httptest.NewRequest("GET", "/billing/choose-provider", nil)
	rec := httptest.NewRecorder()
	h.Show(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

// Mirrors BillingGatewayTest::test_user_can_select_stripe_as_provider
func TestGatewayHandler_Store_Stripe(t *testing.T) {
	billable := &stubBillable{id: 1, btype: "team"}
	resolver := func(r *http.Request) (billing.Billable, error) {
		return billable, nil
	}

	h := newGateway(resolver, &testProductStore{})

	body, _ := json.Marshal(map[string]string{"provider": "stripe"})
	req := httptest.NewRequest("POST", "/billing/choose-provider", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Store(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if billable.provider != "stripe" {
		t.Errorf("provider = %q, want %q", billable.provider, "stripe")
	}

	var resp map[string]string

	json.NewDecoder(rec.Body).Decode(&resp)

	if resp["redirect"] != "/billing" {
		t.Errorf("redirect = %q, want %q", resp["redirect"], "/billing")
	}
}

// Mirrors BillingGatewayTest::test_user_can_select_paddle_as_provider
func TestGatewayHandler_Store_Paddle(t *testing.T) {
	billable := &stubBillable{id: 1, btype: "team"}
	resolver := func(r *http.Request) (billing.Billable, error) {
		return billable, nil
	}

	h := newGateway(resolver, &testProductStore{})

	body, _ := json.Marshal(map[string]string{"provider": "paddle"})
	req := httptest.NewRequest("POST", "/billing/choose-provider", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Store(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if billable.provider != "paddle" {
		t.Errorf("provider = %q, want %q", billable.provider, "paddle")
	}

	var resp map[string]string

	json.NewDecoder(rec.Body).Decode(&resp)

	if resp["redirect"] != "/billing/paddle" {
		t.Errorf("redirect = %q, want %q", resp["redirect"], "/billing/paddle")
	}
}

// Mirrors BillingGatewayTest::test_invalid_provider_is_rejected
func TestGatewayHandler_Store_InvalidProviderReturns422(t *testing.T) {
	billable := &stubBillable{id: 1, btype: "team"}
	resolver := func(r *http.Request) (billing.Billable, error) {
		return billable, nil
	}

	h := newGateway(resolver, &testProductStore{})

	body, _ := json.Marshal(map[string]string{"provider": "invalid"})
	req := httptest.NewRequest("POST", "/billing/choose-provider", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Store(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
}

// Mirrors BillingGatewayTest::test_billing_gateway_shows_products
func TestGatewayHandler_Show_IncludesProducts(t *testing.T) {
	billable := &stubBillable{id: 1, btype: "team"}
	resolver := func(r *http.Request) (billing.Billable, error) {
		return billable, nil
	}

	products := &testProductStore{
		subProducts: []billing.Product{
			{ID: 1, Name: "Pro Plan", Type: billing.ProductTypeSubscription, Active: true},
		},
		oneTimeProducts: []billing.Product{
			{ID: 2, Name: "Theme Pack", Type: billing.ProductTypeOneTime, Active: true},
		},
	}

	h := newGateway(resolver, products)

	req := httptest.NewRequest("GET", "/billing/choose-provider", nil)
	rec := httptest.NewRecorder()
	h.Show(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp map[string]any

	json.NewDecoder(rec.Body).Decode(&resp)

	subProds, ok := resp["subscriptionProducts"].([]any)

	if !ok || len(subProds) != 1 {
		t.Errorf("subscriptionProducts count = %d, want 1", len(subProds))
	}

	oneTimeProds, ok := resp["oneTimeProducts"].([]any)

	if !ok || len(oneTimeProds) != 1 {
		t.Errorf("oneTimeProducts count = %d, want 1", len(oneTimeProds))
	}
}
