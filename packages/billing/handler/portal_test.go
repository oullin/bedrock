package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/handler"
	"github.com/bedrock/packages/billing/state"
)

func newPortalHandler(allow bool, subs []*billing.Subscription) *handler.PortalHandler {
	mgr := billing.NewManager()
	mgr.RegisterBillable(billing.BillableConfig{Model: "team"})
	mgr.Billable("team").Authorize(func(_ billing.Billable, r *http.Request) bool {
		return allow && r.Header.Get("X-Portal-Allow") == "1"
	})

	frontend := state.NewFrontendState(
		mgr,
		billing.NewConfigFromValues(map[string]any{
			"billing.path":          "billing",
			"billing.dashboard_url": "/agreement",
		}),
		&testSubStore{subs: subs},
	)

	resolver := func(r *http.Request) (billing.Billable, error) {
		return &stubBillable{id: 10, btype: "team", name: "Acme"}, nil
	}

	return handler.NewPortalHandler(mgr, frontend, resolver)
}

// AgreementControllerTest::test_index_renders_agreement_page_for_users_with_billing_access
// AgreementControllerTest::test_index_returns_403_without_billing_permission
// AgreementControllerTest::test_index_renders_empty_billing_state_without_a_subscription
// AgreementControllerTest::test_index_ignores_a_tampered_team_id_request_attribute
func TestPortalHandler_Show_BillingAccessAndState(t *testing.T) {
	handler := newPortalHandler(true, nil)

	req := httptest.NewRequest(http.MethodGet, "/billing", nil)
	req.Header.Set("X-Portal-Allow", "1")
	req.SetPathValue("id", "999")
	rec := httptest.NewRecorder()
	handler.Show(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "text/html; charset=utf-8" {
		t.Fatalf("content type = %q, want html", got)
	}

	req = httptest.NewRequest(http.MethodGet, "/billing/state", nil)
	req.Header.Set("X-Portal-Allow", "1")
	rec = httptest.NewRecorder()
	handler.State(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("state status = %d, want 200", rec.Code)
	}

	var resp map[string]any

	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp["billableId"] != float64(10) {
		t.Fatalf("billableId = %v, want 10", resp["billableId"])
	}

	if resp["dashboardUrl"] != "/agreement" {
		t.Fatalf("dashboardUrl = %v, want /agreement", resp["dashboardUrl"])
	}

	if resp["state"] != "none" {
		t.Fatalf("state = %v, want none", resp["state"])
	}

	subscription, ok := resp["subscription"].(map[string]any)

	if !ok {
		t.Fatalf("subscription = %#v, want map", resp["subscription"])
	}

	if subscription["status"] != "" || subscription["plan_code"] != "" {
		t.Fatalf("empty subscription state = %#v", subscription)
	}

	denied := newPortalHandler(false, nil)
	req = httptest.NewRequest(http.MethodGet, "/billing", nil)
	req.Header.Set("X-Portal-Allow", "1")
	rec = httptest.NewRecorder()
	denied.Show(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("unauthorized status = %d, want 403", rec.Code)
	}
}
