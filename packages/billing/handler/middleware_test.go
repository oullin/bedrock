package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/handler"
	"github.com/bedrock/packages/billing/service"
)

func newMiddleware(resolver billing.ResolverFunc, subs []*billing.Subscription) func(http.Handler) http.Handler {
	mgr := billing.NewManager()
	mgr.Billable("team").Resolve(resolver)
	mgr.RegisterBillable(billing.BillableConfig{Model: "team"})

	billing := service.NewBillingService(
		&testSubStore{subs: subs},
		&testOrderStore{},
		&testProductStore{},
	)

	return handler.VerifyBillableIsSubscribed(mgr, billing)
}

// EnsureTeamSubscribedTest::test_unsubscribed_user_is_redirected_to_billing_gateway
func TestVerifyBillableIsSubscribed_Unsubscribed_Redirects(t *testing.T) {
	billable := &stubBillable{id: 1, btype: "team"}
	resolver := func(r *http.Request) (billing.Billable, error) {
		return billable, nil
	}

	mw := newMiddleware(resolver, nil) // no subscriptions

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test-subscribed", nil)
	rec := httptest.NewRecorder()
	mw(inner).ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("status = %d, want %d (redirect)", rec.Code, http.StatusFound)
	}

	location := rec.Header().Get("Location")

	if location != "/billing/choose-provider" {
		t.Errorf("Location = %q, want /billing/choose-provider", location)
	}
}

// EnsureTeamSubscribedTest::test_unsubscribed_xhr_request_returns_402
func TestVerifyBillableIsSubscribed_XHR_Returns402(t *testing.T) {
	billable := &stubBillable{id: 1, btype: "team"}
	resolver := func(r *http.Request) (billing.Billable, error) {
		return billable, nil
	}

	mw := newMiddleware(resolver, nil)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test-subscribed", nil)
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()
	mw(inner).ServeHTTP(rec, req)

	if rec.Code != http.StatusPaymentRequired {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusPaymentRequired)
	}
}

// EnsureTeamSubscribedTest::test_unauthenticated_user_is_redirected
func TestVerifyBillableIsSubscribed_Unauthenticated_Redirects(t *testing.T) {
	resolver := func(r *http.Request) (billing.Billable, error) {
		return nil, errors.New("unauthenticated")
	}

	mw := newMiddleware(resolver, nil)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test-subscribed", nil)
	rec := httptest.NewRecorder()
	mw(inner).ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Errorf("status = %d, want %d (redirect)", rec.Code, http.StatusFound)
	}
}
