package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/handler"
)

type paymentMethodSubStore struct {
	sub *billing.Subscription
}

type paymentMethodProvider struct {
	options map[string]any
}

func (s paymentMethodSubStore) FindByID(context.Context, int64) (*billing.Subscription, error) {
	return nil, nil
}

func (s paymentMethodSubStore) FindByProviderID(context.Context, string) (*billing.Subscription, error) {
	return nil, nil
}

func (s paymentMethodSubStore) CurrentForBillable(context.Context, string, int64) (*billing.Subscription, error) {
	return s.sub, nil
}

func (s paymentMethodSubStore) ActiveForBillable(context.Context, string, int64) ([]*billing.Subscription, error) {
	return nil, nil
}

func (s paymentMethodSubStore) Create(context.Context, *billing.Subscription) error { return nil }
func (s paymentMethodSubStore) Save(context.Context, *billing.Subscription) error   { return nil }
func (s paymentMethodSubStore) Delete(context.Context, int64) error               { return nil }

func (p *paymentMethodProvider) CreatePaymentMethodUpdateTransaction(_ context.Context, _ billing.Billable, _ *billing.Subscription, options map[string]any) (*billing.PaymentMethodUpdateTransaction, error) {
	p.options = options

	return &billing.PaymentMethodUpdateTransaction{
		ID:   "txn_123",
		Data: map[string]any{"mode": "payment_method_update"},
	}, nil
}

// UpdatePaymentMethodController::__invoke
func TestPaymentMethodsHandlerReturnsProviderTransaction(t *testing.T) {
	billable := &stubBillable{id: 10, btype: "team"}
	provider := &paymentMethodProvider{}
	manager := billing.NewManager()
	manager.SetPaymentMethodSessionOptions("team", func(billing.Billable) map[string]any {
		return map[string]any{"return_url": "/billing"}
	})
	h := handler.NewPaymentMethodsHandler(
		func(*http.Request) (billing.Billable, error) {
			return billable, nil
		},
		paymentMethodSubStore{sub: &billing.Subscription{BillableType: "team", BillableID: 10, PaddleID: "sub_123"}},
		provider,
		manager,
	)

	req := httptest.NewRequest(http.MethodPut, "/billing/subscription/payment-method", nil)
	rec := httptest.NewRecorder()
	h.Setup(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	var resp map[string]any

	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp["transaction_id"] != "txn_123" {
		t.Fatalf("transaction_id = %v, want txn_123", resp["transaction_id"])
	}

	if provider.options["return_url"] != "/billing" {
		t.Fatalf("options = %#v, want return_url", provider.options)
	}
}
