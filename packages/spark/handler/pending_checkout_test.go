package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/handler"
)

type pendingCheckoutCustomerStore struct {
	customer *spark.Customer
	created  *spark.Customer
	saved    *spark.Customer
}

func (s *pendingCheckoutCustomerStore) FindByProviderID(context.Context, string) (*spark.Customer, error) {
	return nil, nil
}

func (s *pendingCheckoutCustomerStore) FindByBillable(context.Context, string, int64) (*spark.Customer, error) {
	return s.customer, nil
}

func (s *pendingCheckoutCustomerStore) Create(_ context.Context, customer *spark.Customer) error {
	s.created = customer
	s.customer = customer

	return nil
}

func (s *pendingCheckoutCustomerStore) Save(_ context.Context, customer *spark.Customer) error {
	s.saved = customer
	s.customer = customer

	return nil
}

// NewPendingCheckoutController::__invoke
func TestPendingCheckoutHandlerCreatesCustomerWhenMissing(t *testing.T) {
	store := &pendingCheckoutCustomerStore{}
	billable := &stubBillable{id: 10, btype: "team", name: "Acme", email: "billing@example.com"}
	h := handler.NewPendingCheckoutHandler(store, func(*http.Request) (spark.Billable, error) {
		return billable, nil
	})

	body, _ := json.Marshal(map[string]string{"checkout_id": "chk_123"})
	req := httptest.NewRequest(http.MethodPost, "/spark/pending-checkout", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}

	if store.created == nil {
		t.Fatal("expected customer to be created")
	}

	if store.created.BillableType != "team" || store.created.BillableID != 10 {
		t.Fatalf("customer billable = %s/%d, want team/10", store.created.BillableType, store.created.BillableID)
	}

	if store.created.PendingCheckoutID != "chk_123" {
		t.Fatalf("pending checkout = %q, want chk_123", store.created.PendingCheckoutID)
	}
}

// NewPendingCheckoutController::__invoke
func TestPendingCheckoutHandlerUpdatesExistingCustomer(t *testing.T) {
	store := &pendingCheckoutCustomerStore{
		customer: &spark.Customer{BillableType: "team", BillableID: 10, PendingCheckoutID: "old"},
	}
	h := handler.NewPendingCheckoutHandler(store, func(*http.Request) (spark.Billable, error) {
		return &stubBillable{id: 10, btype: "team"}, nil
	})

	body, _ := json.Marshal(map[string]string{"checkout_id": "chk_456"})
	req := httptest.NewRequest(http.MethodPost, "/spark/pending-checkout", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Create(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}

	if store.saved == nil {
		t.Fatal("expected customer to be saved")
	}

	if store.saved.PendingCheckoutID != "chk_456" {
		t.Fatalf("pending checkout = %q, want chk_456", store.saved.PendingCheckoutID)
	}
}
