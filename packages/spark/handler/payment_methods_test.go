package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/handler"
)

type paymentMethodSubStore struct {
	sub *spark.Subscription
	err error
}

type paymentMethodProvider struct {
	options map[string]any
	err     error
}

func (s paymentMethodSubStore) FindByID(context.Context, int64) (*spark.Subscription, error) {
	return nil, nil
}

func (s paymentMethodSubStore) FindByProviderID(context.Context, string) (*spark.Subscription, error) {
	return nil, nil
}

func (s paymentMethodSubStore) CurrentForBillable(context.Context, string, int64) (*spark.Subscription, error) {
	return s.sub, s.err
}

func (s paymentMethodSubStore) ActiveForBillable(context.Context, string, int64) ([]*spark.Subscription, error) {
	return nil, nil
}

func (s paymentMethodSubStore) Create(context.Context, *spark.Subscription) error { return nil }
func (s paymentMethodSubStore) Save(context.Context, *spark.Subscription) error   { return nil }
func (s paymentMethodSubStore) Delete(context.Context, int64) error               { return nil }

func (p *paymentMethodProvider) CreatePaymentMethodUpdateTransaction(_ context.Context, _ spark.Billable, _ *spark.Subscription, options map[string]any) (*spark.PaymentMethodUpdateTransaction, error) {
	p.options = options

	if p.err != nil {
		return nil, p.err
	}

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

func TestPaymentMethodsHandler_Setup_ErrorPaths(t *testing.T) {
	okResolver := func(*http.Request) (spark.Billable, error) {
		return &stubBillable{id: 10, btype: "team"}, nil
	}

	tests := []struct {
		name     string
		resolver spark.ResolverFunc
		store    paymentMethodSubStore
		provider *paymentMethodProvider
		wantCode int
		wantBody string
	}{
		{
			name:     "resolver error",
			resolver: func(*http.Request) (spark.Billable, error) { return nil, errors.New("no auth") },
			store:    paymentMethodSubStore{},
			provider: &paymentMethodProvider{},
			wantCode: http.StatusUnauthorized,
			wantBody: "Unauthorized",
		},
		{
			name:     "store error",
			resolver: okResolver,
			store:    paymentMethodSubStore{err: errors.New("db")},
			provider: &paymentMethodProvider{},
			wantCode: http.StatusInternalServerError,
			wantBody: "db",
		},
		{
			name:     "no subscription",
			resolver: okResolver,
			store:    paymentMethodSubStore{},
			provider: &paymentMethodProvider{},
			wantCode: http.StatusBadRequest,
			wantBody: spark.ErrNotSubscribed.Error(),
		},
		{
			name:     "provider error",
			resolver: okResolver,
			store:    paymentMethodSubStore{sub: &spark.Subscription{}},
			provider: &paymentMethodProvider{err: errors.New("provider down")},
			wantCode: http.StatusInternalServerError,
			wantBody: "provider down",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.NewPaymentMethodsHandler(tt.resolver, tt.store, tt.provider, nil)

			req := httptest.NewRequest(http.MethodPut, "/spark/subscription/payment-method", nil)
			rec := httptest.NewRecorder()
			h.Setup(rec, req)

			if rec.Code != tt.wantCode {
				t.Fatalf("code = %d, want %d", rec.Code, tt.wantCode)
			}

			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Fatalf("body = %q, want contains %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestPaymentMethodsHandler_SetDefault(t *testing.T) {
	h := handler.NewPaymentMethodsHandler(
		func(*http.Request) (spark.Billable, error) {
			return &stubBillable{id: 1, btype: "team"}, nil
		},
		paymentMethodSubStore{}, &paymentMethodProvider{}, nil,
	)

	rec := httptest.NewRecorder()
	h.SetDefault(rec, httptest.NewRequest(http.MethodPut, "/p", nil))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("SetDefault ok code = %d, want 204", rec.Code)
	}

	denied := handler.NewPaymentMethodsHandler(
		func(*http.Request) (spark.Billable, error) { return nil, errors.New("no") },
		paymentMethodSubStore{}, &paymentMethodProvider{}, nil,
	)

	rec = httptest.NewRecorder()
	denied.SetDefault(rec, httptest.NewRequest(http.MethodPut, "/p", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("SetDefault denied code = %d, want 401", rec.Code)
	}
}

func TestPaymentMethodsHandler_Delete(t *testing.T) {
	h := handler.NewPaymentMethodsHandler(
		func(*http.Request) (spark.Billable, error) {
			return &stubBillable{id: 1, btype: "team"}, nil
		},
		paymentMethodSubStore{}, &paymentMethodProvider{}, nil,
	)

	rec := httptest.NewRecorder()
	h.Delete(rec, httptest.NewRequest(http.MethodDelete, "/p", nil))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("Delete ok code = %d, want 204", rec.Code)
	}

	denied := handler.NewPaymentMethodsHandler(
		func(*http.Request) (spark.Billable, error) { return nil, errors.New("no") },
		paymentMethodSubStore{}, &paymentMethodProvider{}, nil,
	)

	rec = httptest.NewRecorder()
	denied.Delete(rec, httptest.NewRequest(http.MethodDelete, "/p", nil))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("Delete denied code = %d, want 401", rec.Code)
	}
}
