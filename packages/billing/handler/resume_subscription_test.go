package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/handler"
	"github.com/bedrock/packages/billing/service"
)

func TestResumeSubscriptionHandler_ResolverError(t *testing.T) {
	billing := service.NewBillingService(&cancelSubStore{}, &testOrderStore{}, &testProductStore{})

	h := handler.NewResumeSubscriptionHandler(billing, func(*http.Request) (billing.Billable, error) {
		return nil, errors.New("no")
	})

	rec := httptest.NewRecorder()
	h.Resume(rec, httptest.NewRequest(http.MethodPut, "/billing/subscription/resume", nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}
}

func TestResumeSubscriptionHandler_BillingError(t *testing.T) {
	billing := service.NewBillingService(&cancelSubStore{}, &testOrderStore{}, &testProductStore{})

	h := handler.NewResumeSubscriptionHandler(billing, func(*http.Request) (billing.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	rec := httptest.NewRecorder()
	h.Resume(rec, httptest.NewRequest(http.MethodPut, "/billing/subscription/resume", nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400 (no active sub)", rec.Code)
	}
}

func TestResumeSubscriptionHandler_HappyPath(t *testing.T) {
	future := time.Now().Add(48 * time.Hour)
	active := &billing.Subscription{
		ID:       1,
		Status:   billing.StatusCanceled,
		EndsAt:   &future,
		PaddleID: "sub_1",
	}

	store := &cancelSubStore{subs: []*billing.Subscription{active}}
	billing := service.NewBillingService(store, &testOrderStore{}, &testProductStore{})

	h := handler.NewResumeSubscriptionHandler(billing, func(*http.Request) (billing.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	rec := httptest.NewRecorder()
	h.Resume(rec, httptest.NewRequest(http.MethodPut, "/billing/subscription/resume", nil))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("code = %d, want 204", rec.Code)
	}

	if active.Status != billing.StatusActive {
		t.Fatalf("status = %q, want active", active.Status)
	}
}
