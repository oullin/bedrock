package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/handler"
	"github.com/bedrock/packages/spark/state"
)

func newPortalHandler(allow bool, subs []*spark.Subscription) *handler.PortalHandler {
	mgr := spark.NewManager()
	mgr.RegisterBillable(spark.BillableConfig{Model: "team"})
	mgr.Billable("team").Authorize(func(_ spark.Billable, r *http.Request) bool {
		return allow && r.Header.Get("X-Portal-Allow") == "1"
	})

	frontend := state.NewFrontendState(
		mgr,
		spark.NewConfigFromValues(map[string]any{
			"spark.path":          "billing",
			"spark.dashboard_url": "/agreement",
		}),
		&testSubStore{subs: subs},
	)

	resolver := func(r *http.Request) (spark.Billable, error) {
		return &stubBillable{id: 10, btype: "team", name: "Acme"}, nil
	}

	return handler.NewPortalHandler(mgr, frontend, resolver)
}

func newMultiBillablePortalHandler() *handler.PortalHandler {
	mgr := spark.NewManager()
	mgr.RegisterBillable(spark.BillableConfig{Model: "user"})
	mgr.Plan("user", "Personal", "pri_user_monthly").Monthly()
	mgr.Plan("team", "Team", "pri_team_monthly").Monthly()
	mgr.Billable("team").ChargePerSeat("member", func(_ spark.Billable) int {
		return 7
	})

	frontend := state.NewFrontendState(
		mgr,
		spark.NewConfigFromValues(map[string]any{
			"spark.path": "billing",
			"spark.billables": map[string]spark.BillableConfig{
				"team": {DefaultInterval: "yearly"},
				"user": {DefaultInterval: "monthly"},
			},
		}),
		&testSubStore{},
	)

	resolver := func(r *http.Request) (spark.Billable, error) {
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

	req = httptest.NewRequest(http.MethodGet, "/spark/state", nil)
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

func TestPortalHandler_State_UsesResolvedBillableTypeForUntypedPortal(t *testing.T) {
	handler := newMultiBillablePortalHandler()

	req := httptest.NewRequest(http.MethodGet, "/billing", nil)
	rec := httptest.NewRecorder()
	handler.State(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	resp := decodePortalState(t, rec)

	assertTeamPortalState(t, resp)
}

func TestPortalHandler_State_UsesTypedPortalRoute(t *testing.T) {
	handler := newMultiBillablePortalHandler()

	req := httptest.NewRequest(http.MethodGet, "/billing/team", nil)
	req.SetPathValue("type", "team")
	rec := httptest.NewRecorder()
	handler.State(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	resp := decodePortalState(t, rec)

	assertTeamPortalState(t, resp)
}

func TestPortalHandler_State_RejectsMismatchedTypedPortalRoute(t *testing.T) {
	handler := newMultiBillablePortalHandler()

	tests := []struct {
		name      string
		routeType string
		routeID   string
	}{
		{name: "type mismatch", routeType: "user"},
		{name: "id mismatch", routeType: "team", routeID: "999"},
		{name: "invalid id", routeType: "team", routeID: "not-an-id"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/billing/"+tt.routeType, nil)
			req.SetPathValue("type", tt.routeType)

			if tt.routeID != "" {
				req = httptest.NewRequest(http.MethodGet, "/billing/"+tt.routeType+"/"+tt.routeID, nil)
				req.SetPathValue("type", tt.routeType)
				req.SetPathValue("id", tt.routeID)
			}

			rec := httptest.NewRecorder()
			handler.State(rec, req)

			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404", rec.Code)
			}
		})
	}
}

func TestPortalHandler_State_AcceptsMatchingTypedPortalBillableID(t *testing.T) {
	handler := newMultiBillablePortalHandler()

	req := httptest.NewRequest(http.MethodGet, "/billing/team/10", nil)
	req.SetPathValue("type", "team")
	req.SetPathValue("id", "10")
	rec := httptest.NewRecorder()
	handler.State(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	resp := decodePortalState(t, rec)

	assertTeamPortalState(t, resp)
}

func decodePortalState(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var resp map[string]any

	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	return resp
}

func assertTeamPortalState(t *testing.T, resp map[string]any) {
	t.Helper()

	if resp["billableType"] != "team" {
		t.Fatalf("billableType = %v, want team", resp["billableType"])
	}

	if resp["defaultInterval"] != "yearly" {
		t.Fatalf("defaultInterval = %v, want yearly", resp["defaultInterval"])
	}

	if resp["seatName"] != "member" {
		t.Fatalf("seatName = %v, want member", resp["seatName"])
	}

	monthlyPlans, ok := resp["monthlyPlans"].([]any)

	if !ok || len(monthlyPlans) != 1 {
		t.Fatalf("monthlyPlans = %#v, want one plan", resp["monthlyPlans"])
	}

	plan, ok := monthlyPlans[0].(map[string]any)

	if !ok {
		t.Fatalf("monthly plan = %#v, want map", monthlyPlans[0])
	}

	if plan["id"] != "pri_team_monthly" {
		t.Fatalf("monthly plan id = %v, want pri_team_monthly", plan["id"])
	}

	if _, ok := plan["ID"]; ok {
		t.Fatalf("monthly plan exposed raw struct ID key: %#v", plan)
	}
}
