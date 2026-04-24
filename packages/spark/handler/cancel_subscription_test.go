package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/handler"
	"github.com/bedrock/packages/spark/service"
)

type cancelSubStore struct {
	subs    []*spark.Subscription
	findErr error
	saveErr error
}

func (s *cancelSubStore) FindByID(context.Context, int64) (*spark.Subscription, error) {
	return nil, nil
}
func (s *cancelSubStore) FindByProviderID(context.Context, string) (*spark.Subscription, error) {
	return nil, nil
}
func (s *cancelSubStore) CurrentForBillable(context.Context, string, int64) (*spark.Subscription, error) {
	return nil, nil
}
func (s *cancelSubStore) ActiveForBillable(context.Context, string, int64) ([]*spark.Subscription, error) {
	return s.subs, s.findErr
}
func (s *cancelSubStore) Create(context.Context, *spark.Subscription) error { return nil }
func (s *cancelSubStore) Save(context.Context, *spark.Subscription) error   { return s.saveErr }
func (s *cancelSubStore) Delete(context.Context, int64) error               { return nil }

func TestCancelSubscriptionHandler_ResolverError(t *testing.T) {
	billing := service.NewBillingService(&cancelSubStore{}, &testOrderStore{}, &testProductStore{})

	h := handler.NewCancelSubscriptionHandler(billing, func(*http.Request) (spark.Billable, error) {
		return nil, errors.New("no")
	})

	rec := httptest.NewRecorder()
	h.Cancel(rec, httptest.NewRequest(http.MethodPut, "/spark/subscription/cancel", nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}
}

func TestCancelSubscriptionHandler_BillingError(t *testing.T) {
	billing := service.NewBillingService(&cancelSubStore{}, &testOrderStore{}, &testProductStore{})

	h := handler.NewCancelSubscriptionHandler(billing, func(*http.Request) (spark.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	rec := httptest.NewRecorder()
	h.Cancel(rec, httptest.NewRequest(http.MethodPut, "/spark/subscription/cancel", nil))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400 (no active sub)", rec.Code)
	}
}

func TestCancelSubscriptionHandler_HappyPath(t *testing.T) {
	active := &spark.Subscription{ID: 1, Status: spark.StatusActive, PaddleID: "sub_1"}
	store := &cancelSubStore{subs: []*spark.Subscription{active}}
	billing := service.NewBillingService(store, &testOrderStore{}, &testProductStore{})

	h := handler.NewCancelSubscriptionHandler(billing, func(*http.Request) (spark.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	rec := httptest.NewRecorder()
	h.Cancel(rec, httptest.NewRequest(http.MethodPut, "/spark/subscription/cancel", nil))

	if rec.Code != http.StatusNoContent {
		t.Fatalf("code = %d, want 204", rec.Code)
	}

	if active.EndsAt == nil || active.EndsAt.After(time.Now().Add(time.Hour)) {
		t.Fatalf("subscription EndsAt = %v, want near now", active.EndsAt)
	}
}
