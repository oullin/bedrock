package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/handler"
)

type paymentMethodSubStore struct {
	sub *spark.Subscription
}

type paymentMethodProvider struct {
	options map[string]any
}

func (s paymentMethodSubStore) FindByID(context.Context, int64) (*spark.Subscription, error) {
	return nil, nil
}

func (s paymentMethodSubStore) FindByProviderID(context.Context, string) (*spark.Subscription, error) {
	return nil, nil
}

func (s paymentMethodSubStore) CurrentForBillable(context.Context, string, int64) (*spark.Subscription, error) {
	return s.sub, nil
}

func (s paymentMethodSubStore) ActiveForBillable(context.Context, string, int64) ([]*spark.Subscription, error) {
	return nil, nil
}

func (s paymentMethodSubStore) Create(context.Context, *spark.Subscription) error { return nil }
func (s paymentMethodSubStore) Save(context.Context, *spark.Subscription) error   { return nil }
func (s paymentMethodSubStore) Delete(context.Context, int64) error               { return nil }

func (p *paymentMethodProvider) CreatePaymentMethodUpdateTransaction(_ context.Context, _ spark.Billable, _ *spark.Subscription, options map[string]any) (*spark.PaymentMethodUpdateTransaction, error) {
	p.options = options

	return &spark.PaymentMethodUpdateTransaction{
		ID:   "txn_123",
		Data: map[string]any{"mode": "payment_method_update"},
	}, nil
}

// UpdatePaymentMethodController::__invoke
func TestPaymentMethodsHandlerReturnsProviderTransaction(t *testing.T) {
	billable := &stubBillable{id: 10, btype: "team"}
	provider := &paymentMethodProvider{}
	manager := spark.NewManager()
	manager.SetPaymentMethodSessionOptions("team", func(spark.Billable) map[string]any {
		return map[string]any{"return_url": "/billing"}
	})
	h := handler.NewPaymentMethodsHandler(
		func(*http.Request) (spark.Billable, error) {
			return billable, nil
		},
		paymentMethodSubStore{sub: &spark.Subscription{BillableType: "team", BillableID: 10, PaddleID: "sub_123"}},
		provider,
		manager,
	)

	req := httptest.NewRequest(http.MethodPut, "/spark/subscription/payment-method", nil)
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
