package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/handler"
	"github.com/bedrock/packages/billing/service"
)

type cancelSubStore struct {
	subs    []*billing.Subscription
	findErr error
	saveErr error
}

func (s *cancelSubStore) FindByID(context.Context, int64) (*billing.Subscription, error) {
	return nil, nil
}
func (s *cancelSubStore) FindByProviderID(context.Context, string) (*billing.Subscription, error) {
	return nil, nil
}
func (s *cancelSubStore) CurrentForBillable(context.Context, string, int64) (*billing.Subscription, error) {
	return nil, nil
}
func (s *cancelSubStore) ActiveForBillable(context.Context, string, int64) ([]*billing.Subscription, error) {
	return s.subs, s.findErr
}
func (s *cancelSubStore) Create(context.Context, *billing.Subscription) error { return nil }
func (s *cancelSubStore) Save(context.Context, *billing.Subscription) error   { return s.saveErr }
func (s *cancelSubStore) Delete(context.Context, int64) error                 { return nil }

func TestCancelSubscriptionHandler_ResolverError(t *testing.T) {
	svc := service.NewBillingService(&cancelSubStore{}, &testOrderStore{}, &testProductStore{})

	h := handler.NewCancelSubscriptionHandler(svc, func(*http.Request) (billing.Billable, error) {
		return nil, errors.New("no")
	})

	rec := httptest.NewRecorder()
	h.Cancel(rec, httptest.NewRequest(http.MethodPut, "/billing/subscription/cancel", nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}
}

func TestCancelSubscriptionHandler_BillingError(t *testing.T) {
	svc := service.NewBillingService(&cancelSubStore{}, &testOrderStore{}, &testProductStore{})

	h := handler.NewCancelSubscriptionHandler(svc, func(*http.Request) (billing.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	rec := httptest.NewRecorder()
	h.Cancel(rec, httptest.NewRequest(http.MethodPut, "/billing/subscription/cancel", nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400 (no active sub)", rec.Code)
	}
}

func TestCancelSubscriptionHandler_HappyPath(t *testing.T) {
	active := &billing.Subscription{ID: 1, Status: billing.StatusActive, PaddleID: "sub_1"}
	store := &cancelSubStore{subs: []*billing.Subscription{active}}
	svc := service.NewBillingService(store, &testOrderStore{}, &testProductStore{})

	h := handler.NewCancelSubscriptionHandler(svc, func(*http.Request) (billing.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	rec := httptest.NewRecorder()
	h.Cancel(rec, httptest.NewRequest(http.MethodPut, "/billing/subscription/cancel", nil))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("code = %d, want 204", rec.Code)
	}

	if active.EndsAt == nil || active.EndsAt.After(time.Now().Add(time.Hour)) {
		t.Fatalf("subscription EndsAt = %v, want near now", active.EndsAt)
	}
}
