package state

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/billing"
)

type stateBillable struct {
	id    int64
	btype string
	name  string
	email string
}

type fakeSubStore struct {
	sub *billing.Subscription
	err error
}

type fakeCustomerStore struct {
	customer   *billing.Customer
	findErr    error
	saveErr    error
	createErr  error
	savedCalls int
	savedLast  *billing.Customer
}

type fakeTransactionStore struct {
	txns []billing.Transaction
	err  error
}

type fakePaymentReporter struct {
	last, next *billing.Payment
	lastErr    error
	nextErr    error
}

func (b *stateBillable) BillableID() int64               { return b.id }
func (b *stateBillable) BillableType() string            { return b.btype }
func (b *stateBillable) BillableName() string            { return b.name }
func (b *stateBillable) BillableEmail() string           { return b.email }
func (b *stateBillable) PaymentProvider() string         { return "paddle" }
func (b *stateBillable) SetPaymentProvider(string) error { return nil }

func (f *fakeSubStore) FindByID(context.Context, int64) (*billing.Subscription, error) {
	return nil, nil
}
func (f *fakeSubStore) FindByProviderID(context.Context, string) (*billing.Subscription, error) {
	return nil, nil
}
func (f *fakeSubStore) CurrentForBillable(context.Context, string, int64) (*billing.Subscription, error) {
	return f.sub, f.err
}
func (f *fakeSubStore) ActiveForBillable(context.Context, string, int64) ([]*billing.Subscription, error) {
	return nil, nil
}
func (f *fakeSubStore) Create(context.Context, *billing.Subscription) error { return nil }
func (f *fakeSubStore) Save(context.Context, *billing.Subscription) error   { return nil }
func (f *fakeSubStore) Delete(context.Context, int64) error               { return nil }

func (f *fakeCustomerStore) FindByProviderID(context.Context, string) (*billing.Customer, error) {
	return nil, nil
}
func (f *fakeCustomerStore) FindByBillable(context.Context, string, int64) (*billing.Customer, error) {
	return f.customer, f.findErr
}
func (f *fakeCustomerStore) Create(_ context.Context, c *billing.Customer) error {
	f.customer = c

	return f.createErr
}
func (f *fakeCustomerStore) Save(_ context.Context, c *billing.Customer) error {
	f.savedCalls++
	f.savedLast = c

	return f.saveErr
}

func (f *fakeTransactionStore) FindByProviderID(context.Context, string) (*billing.Transaction, error) {
	return nil, nil
}
func (f *fakeTransactionStore) FindByBillable(context.Context, string, int64, int) ([]billing.Transaction, error) {
	return f.txns, f.err
}
func (f *fakeTransactionStore) Create(context.Context, *billing.Transaction) error { return nil }
func (f *fakeTransactionStore) Save(context.Context, *billing.Transaction) error   { return nil }

func (f *fakePaymentReporter) LastPayment(context.Context, *billing.Subscription) (*billing.Payment, error) {
	return f.last, f.lastErr
}
func (f *fakePaymentReporter) NextPayment(context.Context, *billing.Subscription) (*billing.Payment, error) {
	return f.next, f.nextErr
}

func TestRemainingDays(t *testing.T) {
	now := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)

	future := now.Add(48 * time.Hour)
	past := now.Add(-1 * time.Hour)

	tests := []struct {
		name    string
		expires *time.Time
		want    int
	}{
		{name: "nil", expires: nil, want: 0},
		{name: "past", expires: &past, want: 0},
		{name: "future", expires: &future, want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := remainingDays(now, tt.expires); got != tt.want {
				t.Fatalf("remainingDays() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPlanForSubscription(t *testing.T) {
	planA := &billing.Plan{ID: "pri_a", Name: "A"}
	planB := &billing.Plan{ID: "pri_b", Name: "BName", Options: map[string]any{"slug": "b-slug"}}

	plans := []*billing.Plan{planA, planB}

	tests := []struct {
		name string
		sub  *billing.Subscription
		want *billing.Plan
	}{
		{name: "nil subscription", sub: nil, want: nil},
		{
			name: "match by price id",
			sub: &billing.Subscription{
				Items: []billing.SubscriptionItem{{PriceID: "pri_a"}},
			},
			want: planA,
		},
		{
			name: "match by slug",
			sub:  &billing.Subscription{Plan: "b-slug"},
			want: planB,
		},
		{
			name: "match by name",
			sub:  &billing.Subscription{Plan: "BName"},
			want: planB,
		},
		{
			name: "no match",
			sub:  &billing.Subscription{Plan: "nope"},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := planForSubscription(plans, tt.sub); got != tt.want {
				t.Fatalf("planForSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPlansToMaps_SkipsNils(t *testing.T) {
	plans := []*billing.Plan{nil, {ID: "pri_a", Name: "A"}, nil}

	got := plansToMaps(plans)

	if len(got) != 1 {
		t.Fatalf("plansToMaps length = %d, want 1", len(got))
	}

	if got[0]["id"] != "pri_a" {
		t.Fatalf("plan id = %v, want pri_a", got[0]["id"])
	}
}

func TestPlanToMap(t *testing.T) {
	if got := planToMap(nil); got != nil {
		t.Fatalf("planToMap(nil) = %v, want nil", got)
	}

	plan := &billing.Plan{ID: "pri", Name: "Plan"}

	if got := planToMap(plan); got["id"] != "pri" {
		t.Fatalf("planToMap id = %v, want pri", got["id"])
	}
}

func TestPaymentToMap(t *testing.T) {
	if got := paymentToMap(nil, "2006-01-02"); got != nil {
		t.Fatalf("paymentToMap(nil) = %v, want nil", got)
	}

	payment := &billing.Payment{Amount: 1200, Currency: "USD", Date: time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)}

	got := paymentToMap(payment, "2006-01-02")

	if got["currency"] != "USD" {
		t.Fatalf("currency = %v", got["currency"])
	}

	if got["date"] != "2026-05-02" {
		t.Fatalf("date = %v", got["date"])
	}

	if got["amount"] != payment.FormattedAmount() {
		t.Fatalf("amount = %v", got["amount"])
	}
}

func TestBrandColor_AppName_DashboardURL_DateFormat_Defaults(t *testing.T) {
	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.brand_color":   "",
		"billing.app_name":      "",
		"billing.dashboard_url": "",
		"billing.date_format":   "",
	})

	f := NewFrontendState(billing.NewManager(), cfg, &fakeSubStore{})

	if f.brandColor() != "bg-gray-800" {
		t.Fatalf("brandColor fallback = %q", f.brandColor())
	}

	if f.appName() != "Upstream" {
		t.Fatalf("appName fallback = %q", f.appName())
	}

	if f.dashboardURL() != "/" {
		t.Fatalf("dashboardURL fallback = %q", f.dashboardURL())
	}

	if f.dateFormat() != "January 2, 2006" {
		t.Fatalf("dateFormat fallback = %q", f.dateFormat())
	}
}

func TestBrandColor_AppName_DashboardURL_DateFormat_Overrides(t *testing.T) {
	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.brand_color":   "bg-blue",
		"billing.app_name":      "Acme",
		"billing.dashboard_url": "/dash",
		"billing.date_format":   "2006-01-02",
	})

	f := NewFrontendState(billing.NewManager(), cfg, &fakeSubStore{})

	if f.brandColor() != "bg-blue" {
		t.Fatalf("brandColor = %q", f.brandColor())
	}

	if f.appName() != "Acme" {
		t.Fatalf("appName = %q", f.appName())
	}

	if f.dashboardURL() != "/dash" {
		t.Fatalf("dashboardURL = %q", f.dashboardURL())
	}

	if f.dateFormat() != "2006-01-02" {
		t.Fatalf("dateFormat = %q", f.dateFormat())
	}
}

func TestGenericTrialEndsAt(t *testing.T) {
	f := NewFrontendState(billing.NewManager(), nil, &fakeSubStore{})

	if f.genericTrialEndsAt(nil) != nil {
		t.Fatalf("nil customer should return nil")
	}

	if f.genericTrialEndsAt(&billing.Customer{}) != nil {
		t.Fatalf("no trial should return nil")
	}

	expired := time.Now().Add(-time.Hour)

	if f.genericTrialEndsAt(&billing.Customer{TrialEndsAt: &expired}) != nil {
		t.Fatalf("expired trial should return nil")
	}

	future := time.Now().Add(48 * time.Hour)

	got := f.genericTrialEndsAt(&billing.Customer{TrialEndsAt: &future})

	if _, ok := got.(string); !ok {
		t.Fatalf("active trial should return formatted string, got %T", got)
	}
}

func TestProviderCustomerID(t *testing.T) {
	if providerCustomerID(nil) != nil {
		t.Fatal("nil customer should return nil")
	}

	if providerCustomerID(&billing.Customer{}) != nil {
		t.Fatal("empty PaddleID should return nil")
	}

	if got := providerCustomerID(&billing.Customer{PaddleID: "cus_1"}); got != "cus_1" {
		t.Fatalf("providerCustomerID = %v", got)
	}
}

func TestDefaultInterval(t *testing.T) {
	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.billables": map[string]billing.BillableConfig{
			"team": {DefaultInterval: "yearly"},
		},
	})

	f := NewFrontendState(billing.NewManager(), cfg, &fakeSubStore{})

	if got := f.defaultInterval("team"); got != "yearly" {
		t.Fatalf("defaultInterval(team) = %q, want yearly", got)
	}

	if got := f.defaultInterval("user"); got != "monthly" {
		t.Fatalf("defaultInterval(user) = %q, want monthly", got)
	}
}

func TestSubscriptionState_Nil(t *testing.T) {
	f := NewFrontendState(billing.NewManager(), nil, &fakeSubStore{})
	state := f.subscriptionState(nil, nil)

	if state["status"] != "" || state["plan_code"] != "" {
		t.Fatalf("nil subscription state should be empty: %#v", state)
	}

	if state["pay_now"] != false {
		t.Fatalf("pay_now should be false")
	}
}

func TestSubscriptionState_Populated(t *testing.T) {
	cfg := billing.NewConfigFromValues(map[string]any{"billing.path": "billing"})
	f := NewFrontendState(billing.NewManager(), cfg, &fakeSubStore{})

	plan := &billing.Plan{ID: "pri_a", Name: "Pro"}
	sub := &billing.Subscription{
		Status: billing.StatusAwaitingPayment,
		Plan:   "pro",
		Items:  []billing.SubscriptionItem{{PriceID: "pri_a"}},
	}

	state := f.subscriptionState(sub, plan)

	if state["status"] != string(billing.StatusAwaitingPayment) {
		t.Fatalf("status = %v", state["status"])
	}

	if state["plan_code"] != "pro" {
		t.Fatalf("plan_code = %v", state["plan_code"])
	}

	if state["plan_name"] != "Pro" {
		t.Fatalf("plan_name = %v", state["plan_name"])
	}

	if state["pay_now"] != true {
		t.Fatalf("pay_now should be true when awaiting payment with portal url, got %v", state["pay_now"])
	}

	if state["portal_url"] != "/billing" {
		t.Fatalf("portal_url = %v", state["portal_url"])
	}
}

func TestCtaState(t *testing.T) {
	now := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	future := now.Add(96 * time.Hour)

	cfg := billing.NewConfigFromValues(map[string]any{"billing.path": "billing"})
	f := NewFrontendState(billing.NewManager(), cfg, &fakeSubStore{})
	plan := &billing.Plan{ID: "pri_a"}

	tests := []struct {
		name    string
		sub     *billing.Subscription
		plan    *billing.Plan
		visible bool
	}{
		{name: "nil sub", sub: nil, plan: plan, visible: false},
		{name: "nil plan", sub: &billing.Subscription{Status: billing.StatusAwaitingPayment}, plan: nil, visible: false},
		{name: "not awaiting", sub: &billing.Subscription{Status: billing.StatusActive}, plan: plan, visible: false},
		{
			name: "awaiting payment with matching price",
			sub: &billing.Subscription{
				Status:           billing.StatusAwaitingPayment,
				PendingExpiresAt: &future,
				Items:            []billing.SubscriptionItem{{PriceID: "pri_a"}},
			},
			plan:    plan,
			visible: true,
		},
		{
			name:    "awaiting payment but no portal url (price mismatch)",
			sub:     &billing.Subscription{Status: billing.StatusAwaitingPayment},
			plan:    plan,
			visible: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cta := f.ctaState(tt.sub, tt.plan, now)

			if cta["visible"] != tt.visible {
				t.Fatalf("visible = %v, want %v", cta["visible"], tt.visible)
			}
		})
	}
}

func TestPortalURL(t *testing.T) {
	cfg := billing.NewConfigFromValues(map[string]any{"billing.path": "/billing/"})
	f := NewFrontendState(billing.NewManager(), cfg, &fakeSubStore{})
	plan := &billing.Plan{ID: "pri_a"}

	tests := []struct {
		name string
		sub  *billing.Subscription
		plan *billing.Plan
		want string
	}{
		{name: "nil sub", sub: nil, plan: plan, want: ""},
		{name: "nil plan", sub: &billing.Subscription{Status: billing.StatusActive}, plan: nil, want: ""},
		{
			name: "active",
			sub:  &billing.Subscription{Status: billing.StatusActive},
			plan: plan,
			want: "/billing",
		},
		{
			name: "past due",
			sub:  &billing.Subscription{Status: billing.StatusPastDue},
			plan: plan,
			want: "/billing",
		},
		{
			name: "paused",
			sub:  &billing.Subscription{Status: billing.StatusPaused},
			plan: plan,
			want: "/billing",
		},
		{
			name: "awaiting payment matching price",
			sub: &billing.Subscription{
				Status: billing.StatusAwaitingPayment,
				Items:  []billing.SubscriptionItem{{PriceID: "pri_a"}},
			},
			plan: plan,
			want: "/billing",
		},
		{
			name: "awaiting payment price mismatch",
			sub:  &billing.Subscription{Status: billing.StatusAwaitingPayment},
			plan: plan,
			want: "",
		},
		{
			name: "expired status",
			sub:  &billing.Subscription{Status: billing.StatusExpired},
			plan: plan,
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := f.portalURL(tt.sub, tt.plan); got != tt.want {
				t.Fatalf("portalURL = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPortalURL_EmptyPathReturnsRoot(t *testing.T) {
	cfg := billing.NewConfigFromValues(map[string]any{"billing.path": ""})
	f := NewFrontendState(billing.NewManager(), cfg, &fakeSubStore{})

	sub := &billing.Subscription{Status: billing.StatusActive}

	if got := f.portalURL(sub, &billing.Plan{ID: "pri"}); got != "/" {
		t.Fatalf("portalURL = %q, want /", got)
	}
}

func TestSubscriptionIsActiveOrPastDue(t *testing.T) {
	if subscriptionIsActiveOrPastDue(nil) {
		t.Fatal("nil sub should be false")
	}

	if !subscriptionIsActiveOrPastDue(&billing.Subscription{Status: billing.StatusActive}) {
		t.Fatal("active should be true")
	}

	if !subscriptionIsActiveOrPastDue(&billing.Subscription{Status: billing.StatusPastDue}) {
		t.Fatal("past due should be true")
	}

	if subscriptionIsActiveOrPastDue(&billing.Subscription{Status: billing.StatusCanceled}) {
		t.Fatal("canceled should be false")
	}
}

func TestInvoiceID(t *testing.T) {
	if got := invoiceID(billing.Transaction{PaddleID: "tx_1", InvoiceNumber: "INV-1"}); got != "tx_1" {
		t.Fatalf("invoiceID prefers PaddleID: %q", got)
	}

	if got := invoiceID(billing.Transaction{InvoiceNumber: "INV-1"}); got != "INV-1" {
		t.Fatalf("invoiceID falls back to invoice number: %q", got)
	}
}

func TestCurrentAt_HappyPath_ClearsPendingCheckout(t *testing.T) {
	mgr := billing.NewManager()
	mgr.Plan("team", "Pro", "pri_team").Monthly()

	cfg := billing.NewConfigFromValues(map[string]any{
		"billing.path":          "billing",
		"billing.dashboard_url": "/dash",
	})

	sub := &billing.Subscription{Status: billing.StatusActive, Plan: "pro"}

	customers := &fakeCustomerStore{
		customer: &billing.Customer{PendingCheckoutID: "chk_1", PaddleID: "cus_1"},
	}

	f := NewFrontendState(mgr, cfg, &fakeSubStore{sub: sub}).
		WithCustomerStore(customers).
		WithTransactionStore(&fakeTransactionStore{}).
		WithPaymentReporter(&fakePaymentReporter{})

	billable := &stateBillable{id: 10, btype: "team", name: "Acme"}

	data, err := f.CurrentAt(context.Background(), "team", billable, time.Now())

	if err != nil {
		t.Fatalf("CurrentAt: %v", err)
	}

	if data["state"] != "active" {
		t.Fatalf("state = %v, want active", data["state"])
	}

	if data["billableId"] != int64(10) {
		t.Fatalf("billableId = %v", data["billableId"])
	}

	if customers.savedCalls != 1 {
		t.Fatalf("Save call count = %d, want 1", customers.savedCalls)
	}

	if customers.savedLast.PendingCheckoutID != "" {
		t.Fatalf("pending checkout not cleared: %q", customers.savedLast.PendingCheckoutID)
	}
}

func TestCurrentAt_PendingState_WhenCustomerHasPendingCheckout(t *testing.T) {
	mgr := billing.NewManager()
	cfg := billing.NewConfigFromValues(map[string]any{"billing.path": "billing"})

	customers := &fakeCustomerStore{
		customer: &billing.Customer{PendingCheckoutID: "chk_1"},
	}

	f := NewFrontendState(mgr, cfg, &fakeSubStore{}).WithCustomerStore(customers)

	data, err := f.CurrentAt(context.Background(), "team", &stateBillable{id: 1, btype: "team"}, time.Now())

	if err != nil {
		t.Fatalf("CurrentAt: %v", err)
	}

	if data["state"] != "pending" {
		t.Fatalf("state = %v, want pending", data["state"])
	}

	if customers.savedCalls != 0 {
		t.Fatalf("Save should not be called for pending-only state")
	}
}

func TestCurrentAt_NoneState_WhenNoSubscription(t *testing.T) {
	mgr := billing.NewManager()
	cfg := billing.NewConfigFromValues(nil)

	f := NewFrontendState(mgr, cfg, &fakeSubStore{})

	data, err := f.CurrentAt(context.Background(), "team", &stateBillable{id: 1, btype: "team"}, time.Now())

	if err != nil {
		t.Fatalf("CurrentAt: %v", err)
	}

	if data["state"] != "none" {
		t.Fatalf("state = %v, want none", data["state"])
	}
}

func TestCurrentAt_CustomerStoreError(t *testing.T) {
	mgr := billing.NewManager()

	customers := &fakeCustomerStore{findErr: errors.New("boom")}
	f := NewFrontendState(mgr, nil, &fakeSubStore{}).WithCustomerStore(customers)

	_, err := f.CurrentAt(context.Background(), "team", &stateBillable{id: 1, btype: "team"}, time.Now())

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCurrentAt_TransactionsSurfaced(t *testing.T) {
	mgr := billing.NewManager()
	cfg := billing.NewConfigFromValues(map[string]any{"billing.date_format": "2006-01-02"})

	billedAt := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)

	txns := []billing.Transaction{
		{PaddleID: "tx_1", Total: 1200, Currency: "USD", BilledAt: &billedAt},
		{Total: 0}, // filtered
	}

	f := NewFrontendState(mgr, cfg, &fakeSubStore{}).WithTransactionStore(&fakeTransactionStore{txns: txns})

	data, err := f.CurrentAt(context.Background(), "team", &stateBillable{id: 1, btype: "team"}, time.Now())

	if err != nil {
		t.Fatalf("CurrentAt: %v", err)
	}

	invoices, ok := data["invoices"].([]map[string]any)

	if !ok {
		t.Fatalf("invoices wrong type: %T", data["invoices"])
	}

	if len(invoices) != 1 {
		t.Fatalf("invoices count = %d, want 1", len(invoices))
	}

	if invoices[0]["id"] != "tx_1" {
		t.Fatalf("invoice id = %v", invoices[0]["id"])
	}

	if invoices[0]["billed_at"] != "2026-03-01" {
		t.Fatalf("billed_at = %v", invoices[0]["billed_at"])
	}
}

func TestCurrentAt_TransactionStoreError(t *testing.T) {
	mgr := billing.NewManager()
	txns := &fakeTransactionStore{err: errors.New("boom")}
	f := NewFrontendState(mgr, nil, &fakeSubStore{}).WithTransactionStore(txns)

	_, err := f.CurrentAt(context.Background(), "team", &stateBillable{id: 1, btype: "team"}, time.Now())

	if err == nil {
		t.Fatal("expected error from transactions")
	}
}

func TestCurrentAt_PaymentReporterPopulatesPayments(t *testing.T) {
	mgr := billing.NewManager()
	cfg := billing.NewConfigFromValues(map[string]any{"billing.date_format": "2006-01-02"})

	sub := &billing.Subscription{Status: billing.StatusActive}

	reporter := &fakePaymentReporter{
		last: &billing.Payment{Amount: 1000, Currency: "USD", Date: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
		next: &billing.Payment{Amount: 2000, Currency: "USD", Date: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)},
	}

	f := NewFrontendState(mgr, cfg, &fakeSubStore{sub: sub}).WithPaymentReporter(reporter)

	data, err := f.CurrentAt(context.Background(), "team", &stateBillable{id: 1, btype: "team"}, time.Now())

	if err != nil {
		t.Fatalf("CurrentAt: %v", err)
	}

	last, ok := data["lastPayment"].(map[string]any)

	if !ok || last["date"] != "2026-01-01" {
		t.Fatalf("lastPayment = %#v", data["lastPayment"])
	}

	next, ok := data["nextPayment"].(map[string]any)

	if !ok || next["date"] != "2026-02-01" {
		t.Fatalf("nextPayment = %#v", data["nextPayment"])
	}
}

func TestCurrentAt_PaymentReporterError(t *testing.T) {
	sub := &billing.Subscription{Status: billing.StatusActive}
	reporter := &fakePaymentReporter{lastErr: errors.New("boom")}

	f := NewFrontendState(billing.NewManager(), nil, &fakeSubStore{sub: sub}).WithPaymentReporter(reporter)

	_, err := f.CurrentAt(context.Background(), "team", &stateBillable{id: 1, btype: "team"}, time.Now())

	if err == nil {
		t.Fatal("expected error from payments")
	}
}

func TestWithRoutes_NilIsIgnored(t *testing.T) {
	f := NewFrontendState(billing.NewManager(), nil, &fakeSubStore{})

	originalRoutes := f.routes

	f.WithRoutes(nil)

	if f.routes != originalRoutes {
		t.Fatal("nil routes should be ignored")
	}
}

func TestResolveState_CustomerSaveError(t *testing.T) {
	customers := &fakeCustomerStore{
		customer: &billing.Customer{PendingCheckoutID: "chk"},
		saveErr:  errors.New("boom"),
	}

	f := NewFrontendState(billing.NewManager(), nil, &fakeSubStore{}).WithCustomerStore(customers)

	sub := &billing.Subscription{Status: billing.StatusActive}

	_, err := f.resolveState(context.Background(), sub, customers.customer)

	if err == nil {
		t.Fatal("expected save error to surface")
	}
}

func TestResolveState_Canceled(t *testing.T) {
	f := NewFrontendState(billing.NewManager(), nil, &fakeSubStore{})

	endsAt := time.Now().Add(48 * time.Hour)

	state, err := f.resolveState(context.Background(), &billing.Subscription{Status: billing.StatusCanceled, EndsAt: &endsAt}, nil)

	if err != nil {
		t.Fatalf("resolveState: %v", err)
	}

	if state != "onGracePeriod" {
		t.Fatalf("state = %q, want onGracePeriod", state)
	}
}

func TestResolveState_PastDue(t *testing.T) {
	f := NewFrontendState(billing.NewManager(), nil, &fakeSubStore{})

	state, err := f.resolveState(context.Background(), &billing.Subscription{Status: billing.StatusPastDue}, nil)

	if err != nil {
		t.Fatalf("resolveState: %v", err)
	}

	if state != "past_due" {
		t.Fatalf("state = %q", state)
	}
}

func TestResolveState_ExpiredFallsToNone(t *testing.T) {
	f := NewFrontendState(billing.NewManager(), nil, &fakeSubStore{})

	state, err := f.resolveState(context.Background(), &billing.Subscription{Status: billing.StatusExpired}, nil)

	if err != nil {
		t.Fatalf("resolveState: %v", err)
	}

	if state != "none" {
		t.Fatalf("state = %q", state)
	}
}
