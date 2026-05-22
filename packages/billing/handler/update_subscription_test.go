package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/action"
	"github.com/bedrock/packages/billing/handler"
)

type updateSubStore struct {
	sub     *billing.Subscription
	findErr error
	saveErr error
	saved   *billing.Subscription
}

func (s *updateSubStore) FindByID(context.Context, int64) (*billing.Subscription, error) {
	return nil, nil
}
func (s *updateSubStore) FindByProviderID(context.Context, string) (*billing.Subscription, error) {
	return nil, nil
}
func (s *updateSubStore) CurrentForBillable(context.Context, string, int64) (*billing.Subscription, error) {
	return s.sub, s.findErr
}
func (s *updateSubStore) ActiveForBillable(context.Context, string, int64) ([]*billing.Subscription, error) {
	return nil, nil
}
func (s *updateSubStore) Create(context.Context, *billing.Subscription) error { return nil }
func (s *updateSubStore) Save(_ context.Context, sub *billing.Subscription) error {
	s.saved = sub

	return s.saveErr
}
func (s *updateSubStore) Delete(context.Context, int64) error { return nil }

func TestUpdateSubscriptionHandler_ResolverError(t *testing.T) {
	mgr := billing.NewManager()
	store := &updateSubStore{}
	updater := action.NewSubscriptionUpdater(store, mgr)

	h := handler.NewUpdateSubscriptionHandler(updater, mgr, store, func(*http.Request) (billing.Billable, error) {
		return nil, errors.New("no")
	})

	req := httptest.NewRequest(http.MethodPut, "/billing/subscription", bytes.NewReader([]byte(`{}`)))
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}
}

func TestUpdateSubscriptionHandler_InvalidJSON(t *testing.T) {
	mgr := billing.NewManager()
	store := &updateSubStore{sub: &billing.Subscription{}}
	updater := action.NewSubscriptionUpdater(store, mgr)

	h := handler.NewUpdateSubscriptionHandler(updater, mgr, store, func(*http.Request) (billing.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	req := httptest.NewRequest(http.MethodPut, "/billing/subscription", bytes.NewReader([]byte("{not json")))
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}
}

func TestUpdateSubscriptionHandler_NotSubscribed(t *testing.T) {
	mgr := billing.NewManager()
	store := &updateSubStore{}
	updater := action.NewSubscriptionUpdater(store, mgr)

	h := handler.NewUpdateSubscriptionHandler(updater, mgr, store, func(*http.Request) (billing.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	body, _ := json.Marshal(map[string]string{"plan": "pri"})
	req := httptest.NewRequest(http.MethodPut, "/billing/subscription", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d, want 400", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), billing.ErrNotSubscribed.Error()) {
		t.Fatalf("body = %q", rec.Body.String())
	}
}

func TestUpdateSubscriptionHandler_InvalidPlan(t *testing.T) {
	mgr := billing.NewManager()
	mgr.Plan("team", "Pro", "pri_pro").Monthly()

	store := &updateSubStore{sub: &billing.Subscription{ID: 1}}
	updater := action.NewSubscriptionUpdater(store, mgr)

	h := handler.NewUpdateSubscriptionHandler(updater, mgr, store, func(*http.Request) (billing.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	body, _ := json.Marshal(map[string]string{"plan": "unknown"})
	req := httptest.NewRequest(http.MethodPut, "/billing/subscription", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("code = %d, want 422", rec.Code)
	}
}

func TestUpdateSubscriptionHandler_UpdaterError(t *testing.T) {
	mgr := billing.NewManager()
	mgr.Plan("team", "Pro", "pri_pro").Monthly()

	store := &updateSubStore{sub: &billing.Subscription{ID: 1}, saveErr: errors.New("save fail")}
	updater := action.NewSubscriptionUpdater(store, mgr)

	h := handler.NewUpdateSubscriptionHandler(updater, mgr, store, func(*http.Request) (billing.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	body, _ := json.Marshal(map[string]string{"plan": "pri_pro"})
	req := httptest.NewRequest(http.MethodPut, "/billing/subscription", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("code = %d, want 500", rec.Code)
	}
}

func TestUpdateSubscriptionHandler_HappyPath(t *testing.T) {
	mgr := billing.NewManager()
	mgr.Plan("team", "Pro", "pri_pro").Monthly()

	store := &updateSubStore{sub: &billing.Subscription{ID: 1}}
	updater := action.NewSubscriptionUpdater(store, mgr)

	h := handler.NewUpdateSubscriptionHandler(updater, mgr, store, func(*http.Request) (billing.Billable, error) {
		return &stubBillable{id: 1, btype: "team"}, nil
	})

	body, _ := json.Marshal(map[string]string{"plan": "pri_pro"})
	req := httptest.NewRequest(http.MethodPut, "/billing/subscription", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	h.Update(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("code = %d, want 204 (body=%s)", rec.Code, rec.Body.String())
	}

	if store.saved == nil || len(store.saved.Items) != 1 || store.saved.Items[0].PriceID != "pri_pro" {
		t.Fatalf("saved = %#v", store.saved)
	}
}
