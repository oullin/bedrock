package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/handler"
	"github.com/bedrock/packages/spark/service"
)

func TestResumeSubscriptionHandler_ResolverError(t *testing.T) {
	billing := service.NewBillingService(&cancelSubStore{}, &testOrderStore{}, &testProductStore{})

	h := handler.NewResumeSubscriptionHandler(billing, func(*http.Request) (spark.Billable, error) {
		return nil, errors.New("no")
	})

	rec := httptest.NewRecorder()
	h.Resume(rec, httptest.NewRequest(http.MethodPut, "/spark/subscription/resume", nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}
}

func TestResumeSubscriptionHandler_BillingError(t *testing.T) {
	billing := service.NewBillingService(&cancelSubStore{}, &testOrderStore{}, &testProductStore{})

	h := handler.NewResumeSubscriptionHandler(billing, func(*http.Request) (spark.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	rec := httptest.NewRecorder()
	h.Resume(rec, httptest.NewRequest(http.MethodPut, "/spark/subscription/resume", nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400 (no active sub)", rec.Code)
	}
}

func TestResumeSubscriptionHandler_HappyPath(t *testing.T) {
	future := time.Now().Add(48 * time.Hour)
	active := &spark.Subscription{
		ID:       1,
		Status:   spark.StatusCanceled,
		EndsAt:   &future,
		PaddleID: "sub_1",
	}

	store := &cancelSubStore{subs: []*spark.Subscription{active}}
	billing := service.NewBillingService(store, &testOrderStore{}, &testProductStore{})

	h := handler.NewResumeSubscriptionHandler(billing, func(*http.Request) (spark.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	rec := httptest.NewRecorder()
	h.Resume(rec, httptest.NewRequest(http.MethodPut, "/spark/subscription/resume", nil))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("code = %d, want 204", rec.Code)
	}

	if active.Status != spark.StatusActive {
		t.Fatalf("status = %q, want active", active.Status)
	}
}
