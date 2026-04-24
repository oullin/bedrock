package spark_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/action"
	"github.com/bedrock/packages/billing/service"
	"github.com/bedrock/packages/billing/state"
	"github.com/bedrock/packages/billing/webhook"
)

type madoraBillable struct {
	id    int64
	typ   string
	name  string
	email string
	seats int
}

type madoraSubStore struct {
	subs    []*billing.Subscription
	saved   []*billing.Subscription
	deleted []int64
}

type madoraCustomerStore struct {
	customer *billing.Customer
	saved    *billing.Customer
}

type madoraOrderStore struct{}

type madoraProductStore struct{}

type madoraTransactionStore struct {
	transactions map[string]*billing.Transaction
	created      []*billing.Transaction
	saved        []*billing.Transaction
}

type madoraExpirableStore struct {
	subscriptions []*billing.Subscription
	deleted       []int64
}

type madoraDispatcher struct {
	events []any
}

func (b madoraBillable) BillableID() int64     { return b.id }
func (b madoraBillable) BillableType() string  { return b.typ }
func (b madoraBillable) BillableName() string  { return b.name }
func (b madoraBillable) BillableEmail() string { return b.email }

func (s *madoraSubStore) FindByID(_ context.Context, id int64) (*billing.Subscription, error) {
	for _, sub := range s.subs {
		if sub.ID == id {
			return sub, nil
		}
	}

	return nil, nil
}

func (s *madoraSubStore) FindByProviderID(_ context.Context, providerID string) (*billing.Subscription, error) {
	for _, sub := range s.subs {
		if sub.PaddleID == providerID {
			return sub, nil
		}
	}

	return nil, nil
}

func (s *madoraSubStore) CurrentForBillable(_ context.Context, billableType string, billableID int64) (*billing.Subscription, error) {
	for _, sub := range s.subs {
		if sub.BillableType == billableType && sub.BillableID == billableID {
			return sub, nil
		}
	}

	return nil, nil
}

func (s *madoraSubStore) ActiveForBillable(_ context.Context, billableType string, billableID int64) ([]*billing.Subscription, error) {
	var matched []*billing.Subscription

	for _, sub := range s.subs {
		if sub.BillableType == billableType && sub.BillableID == billableID {
			matched = append(matched, sub)
		}
	}

	return matched, nil
}

func (s *madoraSubStore) Create(_ context.Context, sub *billing.Subscription) error {
	s.subs = append(s.subs, sub)

	return nil
}

func (s *madoraSubStore) Save(_ context.Context, sub *billing.Subscription) error {
	s.saved = append(s.saved, sub)

	return nil
}

func (s *madoraSubStore) Delete(_ context.Context, id int64) error {
	s.deleted = append(s.deleted, id)

	for i, sub := range s.subs {
		if sub.ID == id {
			s.subs = append(s.subs[:i], s.subs[i+1:]...)

			return nil
		}
	}

	return nil
}

func (s *madoraCustomerStore) FindByProviderID(context.Context, string) (*billing.Customer, error) {
	return nil, nil
}

func (s *madoraCustomerStore) FindByBillable(context.Context, string, int64) (*billing.Customer, error) {
	return s.customer, nil
}

func (s *madoraCustomerStore) Create(_ context.Context, customer *billing.Customer) error {
	s.customer = customer

	return nil
}

func (s *madoraCustomerStore) Save(_ context.Context, customer *billing.Customer) error {
	s.saved = customer
	s.customer = customer

	return nil
}

func (s madoraOrderStore) FindByID(context.Context, int64) (*billing.Order, error) { return nil, nil }
func (s madoraOrderStore) FindByBillable(context.Context, int64, int) ([]billing.Order, error) {
	return nil, nil
}
func (s madoraOrderStore) Create(context.Context, *billing.Order) error { return nil }
func (s madoraOrderStore) Save(context.Context, *billing.Order) error   { return nil }
func (s madoraOrderStore) HasCompletedForProduct(context.Context, int64, int64) (bool, error) {
	return false, nil
}

func (s madoraProductStore) FindByID(context.Context, int64) (*billing.Product, error) { return nil, nil }
func (s madoraProductStore) Active(context.Context) ([]billing.Product, error)         { return nil, nil }
func (s madoraProductStore) ActiveSubscriptions(context.Context) ([]billing.Product, error) {
	return nil, nil
}
func (s madoraProductStore) ActiveOneTime(context.Context) ([]billing.Product, error) {
	return nil, nil
}

func (s *madoraTransactionStore) FindByProviderID(_ context.Context, providerID string) (*billing.Transaction, error) {
	return s.transactions[providerID], nil
}

func (s *madoraTransactionStore) FindByBillable(context.Context, string, int64, int) ([]billing.Transaction, error) {
	return nil, nil
}

func (s *madoraTransactionStore) Create(_ context.Context, tx *billing.Transaction) error {
	s.created = append(s.created, tx)
	s.transactions[tx.PaddleID] = tx

	return nil
}

func (s *madoraTransactionStore) Save(_ context.Context, tx *billing.Transaction) error {
	s.saved = append(s.saved, tx)
	s.transactions[tx.PaddleID] = tx

	return nil
}

func (s *madoraExpirableStore) ExpirableSubscriptions(context.Context, time.Time) ([]*billing.Subscription, error) {
	return s.subscriptions, nil
}

func (s *madoraExpirableStore) Delete(_ context.Context, id int64) error {
	s.deleted = append(s.deleted, id)

	for i, subscription := range s.subscriptions {
		if subscription.ID == id {
			s.subscriptions = append(s.subscriptions[:i], s.subscriptions[i+1:]...)

			return nil
		}
	}

	return nil
}

func (d *madoraDispatcher) Dispatch(event any) error {
	d.events = append(d.events, event)

	return nil
}

func testBillingManager() *billing.Manager {
	mgr := billing.NewManager()
	mgr.RegisterBillable(billing.BillableConfig{
		Model:           "team",
		DefaultInterval: "monthly",
	})
	mgr.Billable("team").
		ChargePerSeat("seat", func(b billing.Billable) int {
			if team, ok := b.(madoraBillable); ok {
				return team.seats
			}

			return 1
		}).
		Plan("Starter", "pri_starter_monthly").
		SetFeatures([]string{"1 workspace", "5 seats"}).
		SetOptions(map[string]any{"slug": "starter"}).
		SetShortDescription("Small team").
		SetInterval("monthly")
	mgr.Billable("team").
		Plan("Pro", "pri_pro_yearly").
		SetFeatures([]string{"Unlimited projects"}).
		SetOptions(map[string]any{"slug": "pro"}).
		SetInterval("yearly")
	mgr.Billable("team").
		Plan("Archived", "pri_archived_monthly").
		SetInterval("monthly").
		Archive()

	return mgr
}

// BillingPortalTest::test_spark_dashboard_url_defaults_to_the_agreement_page
// BillingPortalTest::test_spark_portal_route_keeps_the_default_billing_path
// CatalogPresentersTest::test_frontend_plan_presenter_indexes_seeded_plans_with_default_periods
// CatalogPresentersTest::test_landing_plan_presenter_serialises_active_prices_per_plan
// BillingPlanRegistrarTest::test_it_registers_spark_plans_from_the_active_catalog_prices
// PlanRepositoryHydrationTest::test_plan_repository_returns_plans_with_inverse_relations_hydrated
// BillingLifecycleTest::test_plan_seeder_persists_plans_and_features
func TestMadoraPlanCatalogAndPortalState(t *testing.T) {
	mgr := testBillingManager()
	cfg := billing.Config{
		Path:         "billing",
		DashboardURL: "/agreement",
		BrandColor:   "bg-gray-800",
		Billables: map[string]billing.BillableConfig{
			"team": {DefaultInterval: "monthly"},
		},
	}
	billable := madoraBillable{id: 10, typ: "team", name: "Acme", seats: 4}
	subscriptions := &madoraSubStore{subs: []*billing.Subscription{{
		ID:           99,
		BillableType: "team",
		BillableID:   10,
		PaddleID:     "sub_123",
		Status:       billing.StatusActive,
		Items: []billing.SubscriptionItem{{
			PriceID:  "pri_starter_monthly",
			Quantity: 4,
		}},
	}}}

	frontend := state.NewFrontendState(mgr, &cfg, subscriptions)
	current, err := frontend.Current(context.Background(), "team", billable)

	if err != nil {
		t.Fatalf("frontend state: %v", err)
	}

	if current["dashboardUrl"] != "/agreement" {
		t.Fatalf("dashboardUrl = %v, want /agreement", current["dashboardUrl"])
	}

	if current["sparkPath"] != "billing" {
		t.Fatalf("sparkPath = %v, want billing", current["sparkPath"])
	}

	if current["defaultInterval"] != "monthly" {
		t.Fatalf("defaultInterval = %v, want monthly", current["defaultInterval"])
	}

	if got := len(current["monthlyPlans"].([]*billing.Plan)); got != 1 {
		t.Fatalf("monthlyPlans = %d, want 1 active monthly plan", got)
	}

	if got := len(current["yearlyPlans"].([]*billing.Plan)); got != 1 {
		t.Fatalf("yearlyPlans = %d, want 1 active yearly plan", got)
	}

	activePlan := current["plan"].(*billing.Plan)

	if activePlan.ID != "pri_starter_monthly" {
		t.Fatalf("active plan = %s, want pri_starter_monthly", activePlan.ID)
	}

	planMap := activePlan.ToMap()

	if planMap["price_includes_vat"] != true {
		t.Fatalf("price_includes_vat = %v, want true", planMap["price_includes_vat"])
	}

	if len(planMap["features"].([]string)) != 2 {
		t.Fatalf("features were not serialised: %#v", planMap["features"])
	}

	if mgr.SeatName("team") != "seat" || mgr.SeatCount("team", billable) != 4 {
		t.Fatalf("seat billing was not registered")
	}
}

// FrontendStateTest::test_subscription_is_inactive_when_missing
// FrontendStateTest::test_active_subscription_is_treated_as_active_when_not_on_grace_period
// FrontendStateTest::test_trialing_subscription_is_treated_as_active_when_not_on_grace_period
// FrontendStateTest::test_subscription_on_grace_period_is_not_treated_as_active
// FrontendStateTest::test_inactive_subscription_is_not_treated_as_active
// BillingLifecycleTest::test_billing_state_dto_for_team_without_subscription
// BillingLifecycleTest::test_billing_state_dto_for_active_subscription
// SubscriptionScopeTest::test_accessible_scope_returns_only_access_granting_non_expired_subscriptions
func TestMadoraFrontendStateSubscriptionStates(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)
	past := now.Add(-24 * time.Hour)

	tests := []struct {
		name string
		sub  *billing.Subscription
		want string
	}{
		{name: "missing", want: "none"},
		{name: "active", sub: &billing.Subscription{BillableType: "team", BillableID: 10, Status: billing.StatusActive}, want: "active"},
		{name: "trialing", sub: &billing.Subscription{BillableType: "team", BillableID: 10, Status: billing.StatusTrialing, TrialEndsAt: &future}, want: "active"},
		{name: "grace", sub: &billing.Subscription{BillableType: "team", BillableID: 10, Status: billing.StatusCanceled, EndsAt: &future}, want: "onGracePeriod"},
		{name: "inactive", sub: &billing.Subscription{BillableType: "team", BillableID: 10, Status: billing.StatusCanceled, EndsAt: &past}, want: "none"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var subs []*billing.Subscription

			if tt.sub != nil {
				subs = append(subs, tt.sub)
			}

			frontend := state.NewFrontendState(testBillingManager(), &billing.Config{Path: "billing"}, &madoraSubStore{subs: subs})
			current, err := frontend.Current(context.Background(), "team", madoraBillable{id: 10, typ: "team", name: "Acme"})

			if err != nil {
				t.Fatalf("frontend state: %v", err)
			}

			if current["state"] != tt.want {
				t.Fatalf("state = %v, want %s", current["state"], tt.want)
			}
		})
	}

	store := &madoraSubStore{subs: []*billing.Subscription{
		{BillableType: "team", BillableID: 10, Status: billing.StatusPaused},
		{BillableType: "team", BillableID: 10, Status: billing.StatusPastDue},
	}}
	billing := service.NewBillingService(store, madoraOrderStore{}, madoraProductStore{})
	active, err := billing.GetActiveSubscription(context.Background(), "team", 10)

	if err != nil {
		t.Fatalf("active subscription: %v", err)
	}

	if active == nil || active.Subscription.Status != billing.StatusPastDue {
		t.Fatalf("accessible subscription = %#v, want past_due", active)
	}
}

// FrontendStateTest::test_state_is_pending_when_customer_has_pending_checkout
// FrontendStateTest::test_pending_checkout_is_cleared_when_subscription_becomes_active
func TestMadoraFrontendStatePendingCheckout(t *testing.T) {
	customer := &madoraCustomerStore{
		customer: &billing.Customer{BillableType: "team", BillableID: 10, PendingCheckoutID: "chk_123"},
	}
	frontend := state.NewFrontendState(
		testBillingManager(),
		&billing.Config{Path: "billing"},
		&madoraSubStore{},
	).WithCustomerStore(customer)

	current, err := frontend.Current(context.Background(), "team", madoraBillable{id: 10, typ: "team", name: "Acme"})

	if err != nil {
		t.Fatalf("frontend state: %v", err)
	}

	if current["state"] != "pending" {
		t.Fatalf("state = %v, want pending", current["state"])
	}

	active := &billing.Subscription{BillableType: "team", BillableID: 10, Status: billing.StatusActive}
	customer.customer.PendingCheckoutID = "chk_456"
	frontend = state.NewFrontendState(
		testBillingManager(),
		&billing.Config{Path: "billing"},
		&madoraSubStore{subs: []*billing.Subscription{active}},
	).WithCustomerStore(customer)

	current, err = frontend.Current(context.Background(), "team", madoraBillable{id: 10, typ: "team", name: "Acme"})

	if err != nil {
		t.Fatalf("frontend state: %v", err)
	}

	if current["state"] != "active" {
		t.Fatalf("state = %v, want active", current["state"])
	}

	if customer.saved == nil || customer.saved.PendingCheckoutID != "" {
		t.Fatalf("saved customer = %#v, want cleared pending checkout", customer.saved)
	}
}

// BillingLifecycleTest::test_billing_state_dto_for_team_with_awaiting_payment_subscription
// BillingLifecycleTest::test_billing_state_dto_gates_portal_when_checkout_prices_are_missing
// BillingLifecycleTest::test_billing_state_dto_for_starter_trial_has_upgrade_cta
// BillingLifecycleTest::test_billing_state_dto_serialises_to_array
// BillingLifecycleTest::test_billing_state_dto_for_paused_subscription
// BillingLifecycleTest::test_billing_state_dto_for_canceled_subscription_returns_canceled_state
func TestMadoraBillingStatePortalAndCTA(t *testing.T) {
	now := time.Date(2026, 4, 7, 9, 15, 0, 0, time.UTC)
	expiresAt := now.AddDate(0, 0, billing.DefaultPendingExpiryDays)
	readyAt := now
	trialEndsAt := now.AddDate(0, 0, 14)

	tests := []struct {
		name              string
		sub               *billing.Subscription
		wantStatus        billing.SubscriptionStatus
		wantPlanCode      string
		wantPlanName      string
		wantPortalURL     string
		wantPayNow        bool
		wantCTAVisible    bool
		wantCTALabel      string
		wantRemainingDays int
	}{
		{
			name: "awaiting payment has checkout portal",
			sub: &billing.Subscription{
				BillableType:     "team",
				BillableID:       10,
				Status:           billing.StatusAwaitingPayment,
				Plan:             "pro",
				PendingExpiresAt: &expiresAt,
				PaymentReadyAt:   &readyAt,
				Items:            []billing.SubscriptionItem{{PriceID: "pri_pro_yearly"}},
			},
			wantStatus:        billing.StatusAwaitingPayment,
			wantPlanCode:      "pro",
			wantPlanName:      "Pro",
			wantPortalURL:     "/billing",
			wantPayNow:        true,
			wantCTAVisible:    true,
			wantCTALabel:      "Complete Subscription",
			wantRemainingDays: 14,
		},
		{
			name: "awaiting payment without checkout price is gated",
			sub: &billing.Subscription{
				BillableType:     "team",
				BillableID:       10,
				Status:           billing.StatusAwaitingPayment,
				Plan:             "pro",
				PendingExpiresAt: &expiresAt,
				PaymentReadyAt:   &readyAt,
				Items:            []billing.SubscriptionItem{{PriceID: "pri_missing"}},
			},
			wantStatus:   billing.StatusAwaitingPayment,
			wantPlanCode: "pro",
			wantPlanName: "Pro",
		},
		{
			name: "starter trial has no payment portal",
			sub: &billing.Subscription{
				BillableType: "team",
				BillableID:   10,
				Status:       billing.StatusTrialing,
				Plan:         "starter",
				TrialEndsAt:  &trialEndsAt,
			},
			wantStatus:   billing.StatusTrialing,
			wantPlanCode: "starter",
			wantPlanName: "Starter",
		},
		{
			name: "paused subscription keeps portal access",
			sub: &billing.Subscription{
				BillableType: "team",
				BillableID:   10,
				Status:       billing.StatusPaused,
				Plan:         "pro",
				Items:        []billing.SubscriptionItem{{PriceID: "pri_pro_yearly"}},
			},
			wantStatus:    billing.StatusPaused,
			wantPlanCode:  "pro",
			wantPlanName:  "Pro",
			wantPortalURL: "/billing",
		},
		{
			name: "canceled subscription serialises canceled state",
			sub: &billing.Subscription{
				BillableType: "team",
				BillableID:   10,
				Status:       billing.StatusCanceled,
				Plan:         "pro",
				Items:        []billing.SubscriptionItem{{PriceID: "pri_pro_yearly"}},
			},
			wantStatus:   billing.StatusCanceled,
			wantPlanCode: "pro",
			wantPlanName: "Pro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frontend := state.NewFrontendState(
				testBillingManager(),
				&billing.Config{Path: "billing", Billables: map[string]billing.BillableConfig{"team": {DefaultInterval: "monthly"}}},
				&madoraSubStore{subs: []*billing.Subscription{tt.sub}},
			)
			current, err := frontend.CurrentAt(context.Background(), "team", madoraBillable{id: 10, typ: "team", name: "Acme"}, now)

			if err != nil {
				t.Fatalf("frontend state: %v", err)
			}

			subscription := current["subscription"].(map[string]any)
			cta := current["cta"].(map[string]any)

			for _, key := range []string{"status", "plan_code", "plan_name", "pending_expires_at", "payment_ready_at", "portal_url", "pay_now"} {
				if _, ok := subscription[key]; !ok {
					t.Fatalf("subscription state missing key %q: %#v", key, subscription)
				}
			}

			if subscription["status"] != string(tt.wantStatus) {
				t.Fatalf("status = %v, want %s", subscription["status"], tt.wantStatus)
			}

			if subscription["plan_code"] != tt.wantPlanCode {
				t.Fatalf("plan_code = %v, want %s", subscription["plan_code"], tt.wantPlanCode)
			}

			if subscription["plan_name"] != tt.wantPlanName {
				t.Fatalf("plan_name = %v, want %s", subscription["plan_name"], tt.wantPlanName)
			}

			if subscription["portal_url"] != tt.wantPortalURL {
				t.Fatalf("portal_url = %v, want %s", subscription["portal_url"], tt.wantPortalURL)
			}

			if subscription["pay_now"] != tt.wantPayNow {
				t.Fatalf("pay_now = %v, want %v", subscription["pay_now"], tt.wantPayNow)
			}

			if cta["visible"] != tt.wantCTAVisible {
				t.Fatalf("cta.visible = %v, want %v", cta["visible"], tt.wantCTAVisible)
			}

			if cta["label"] != tt.wantCTALabel {
				t.Fatalf("cta.label = %v, want %s", cta["label"], tt.wantCTALabel)
			}

			if cta["remaining_days"] != tt.wantRemainingDays {
				t.Fatalf("cta.remaining_days = %v, want %d", cta["remaining_days"], tt.wantRemainingDays)
			}
		})
	}
}

// BillingLifecycleTest::test_full_lifecycle_create_pending_prepare_payment_activate
// BillingLifecycleTest::test_starter_subscription_is_created_as_a_fourteen_day_trial
// BillingLifecycleTest::test_full_lifecycle_create_pending_and_expire
// BillingLifecycleTest::test_pending_expiry_days_is_fourteen
// SubscriptionStatusTransitionTest::test_mark_payment_ready_uses_the_test_clock_and_is_idempotent
// SubscriptionStatusTransitionTest::test_mark_past_due_is_idempotent
// SubscriptionStatusTransitionTest::test_pause_is_idempotent
// SubscriptionStatusTransitionTest::test_cancel_is_idempotent
func TestMadoraSubscriptionLifecycleStateMachine(t *testing.T) {
	now := time.Date(2026, 4, 7, 9, 15, 0, 0, time.UTC)
	billable := madoraBillable{id: 10, typ: "team"}

	sub := billing.NewPendingSubscription(billable, "pro", now)

	if sub.Status != billing.StatusPending || sub.Plan != "pro" || sub.PendingExpiresAt == nil {
		t.Fatalf("pending subscription = %#v", sub)
	}

	if got := int(sub.PendingExpiresAt.Sub(now).Hours() / 24); got != billing.DefaultPendingExpiryDays {
		t.Fatalf("pending expiry days = %d, want %d", got, billing.DefaultPendingExpiryDays)
	}

	if !sub.MarkPaymentReady(now) {
		t.Fatalf("first payment-ready transition returned false")
	}

	if sub.Status != billing.StatusAwaitingPayment || sub.PaymentReadyAt == nil || !sub.PaymentReadyAt.Equal(now) {
		t.Fatalf("awaiting-payment subscription = %#v", sub)
	}

	if sub.MarkPaymentReady(now.Add(time.Hour)) || !sub.PaymentReadyAt.Equal(now) {
		t.Fatalf("payment-ready transition was not idempotent: %#v", sub)
	}

	if !sub.Activate(now.Add(2*time.Hour)) || sub.Status != billing.StatusActive {
		t.Fatalf("active subscription = %#v", sub)
	}

	expired := billing.NewPendingSubscription(billable, "starter", now)

	if !expired.Expire(now.Add(time.Hour)) || expired.Status != billing.StatusExpired {
		t.Fatalf("expired subscription = %#v", expired)
	}

	trial := billing.NewTrialSubscription(billable, "starter", 14, now)

	if trial.Status != billing.StatusTrialing || trial.TrialEndsAt == nil || trial.PaymentReadyAt != nil {
		t.Fatalf("trial subscription = %#v", trial)
	}

	pastDue := &billing.Subscription{Status: billing.StatusActive, UpdatedAt: now}

	if !pastDue.MarkPastDue(now.Add(time.Hour)) || pastDue.Status != billing.StatusPastDue {
		t.Fatalf("past-due subscription = %#v", pastDue)
	}

	pastDueUpdatedAt := pastDue.UpdatedAt

	if pastDue.MarkPastDue(now.Add(2*time.Hour)) || !pastDue.UpdatedAt.Equal(pastDueUpdatedAt) {
		t.Fatalf("past-due transition was not idempotent: %#v", pastDue)
	}

	paused := &billing.Subscription{Status: billing.StatusActive, UpdatedAt: now}

	if !paused.Pause(now.Add(time.Hour)) || paused.Status != billing.StatusPaused || paused.PausedAt == nil {
		t.Fatalf("paused subscription = %#v", paused)
	}

	pausedAt := paused.PausedAt

	if paused.Pause(now.Add(2*time.Hour)) || paused.PausedAt != pausedAt {
		t.Fatalf("pause transition was not idempotent: %#v", paused)
	}

	canceled := &billing.Subscription{Status: billing.StatusActive, UpdatedAt: now}

	if !canceled.Cancel(now.Add(time.Hour)) || canceled.Status != billing.StatusCanceled || canceled.EndsAt == nil {
		t.Fatalf("canceled subscription = %#v", canceled)
	}

	endsAt := canceled.EndsAt

	if canceled.Cancel(now.Add(2*time.Hour)) || canceled.EndsAt != endsAt {
		t.Fatalf("cancel transition was not idempotent: %#v", canceled)
	}
}

// BillingLifecycleTest::test_stale_subscription_batch_expiry
// BillingLifecycleTest::test_active_subscription_is_not_expired_by_batch
// SubscriptionScopeTest::test_expirable_scope_returns_only_stale_pending_state_subscriptions
// SubscriptionStatusTransitionTest::test_paused_subscription_is_not_expired_by_batch
// SubscriptionStatusTransitionTest::test_canceled_subscription_is_not_expired_by_batch
func TestMadoraStaleSubscriptionExpiryBatch(t *testing.T) {
	now := time.Date(2026, 4, 8, 10, 0, 0, 0, time.UTC)
	past := now.Add(-24 * time.Hour)
	future := now.Add(24 * time.Hour)

	stalePending := &billing.Subscription{ID: 1, Status: billing.StatusPending, PendingExpiresAt: &past}
	staleAwaiting := &billing.Subscription{ID: 2, Status: billing.StatusAwaitingPayment, PendingExpiresAt: &past}
	freshPending := &billing.Subscription{ID: 3, Status: billing.StatusPending, PendingExpiresAt: &future}
	active := &billing.Subscription{ID: 4, Status: billing.StatusActive, PendingExpiresAt: &past}
	paused := &billing.Subscription{ID: 5, Status: billing.StatusPaused, PendingExpiresAt: &past}
	canceled := &billing.Subscription{ID: 6, Status: billing.StatusCanceled, PendingExpiresAt: &past}
	staleTrial := &billing.Subscription{ID: 7, Status: billing.StatusTrialing, Plan: "starter", PendingExpiresAt: &past}

	store := &madoraExpirableStore{subscriptions: []*billing.Subscription{
		stalePending,
		staleAwaiting,
		freshPending,
		active,
		paused,
		canceled,
		staleTrial,
	}}

	expired, err := billing.ExpireStaleSubscriptions(context.Background(), store, now)

	if err != nil {
		t.Fatalf("expire stale subscriptions: %v", err)
	}

	if expired != 3 {
		t.Fatalf("expired = %d, want 3", expired)
	}

	if !reflect.DeepEqual(store.deleted, []int64{1, 2, 7}) {
		t.Fatalf("deleted subscriptions = %#v, want [1 2 7]", store.deleted)
	}

	if stalePending.Status != billing.StatusExpired || staleAwaiting.Status != billing.StatusExpired || staleTrial.Status != billing.StatusExpired {
		t.Fatalf("stale statuses = %s/%s/%s, want expired", stalePending.Status, staleAwaiting.Status, staleTrial.Status)
	}

	if freshPending.Status != billing.StatusPending || active.Status != billing.StatusActive ||
		paused.Status != billing.StatusPaused || canceled.Status != billing.StatusCanceled {
		t.Fatalf("non-expirable statuses changed: %s/%s/%s/%s", freshPending.Status, active.Status, paused.Status, canceled.Status)
	}
}

// BillingLifecycleTest::test_feature_entitlements_provisioned_on_subscription_creation
// BillingLifecycleTest::test_feature_entitlements_deactivated_on_expiry
// SubscriptionStatusTransitionTest::test_active_to_past_due_keeps_entitlements_active
// SubscriptionStatusTransitionTest::test_past_due_to_active_preserves_entitlements
// SubscriptionStatusTransitionTest::test_active_to_paused_deactivates_entitlements
// SubscriptionStatusTransitionTest::test_paused_to_active_reactivates_entitlements
// SubscriptionStatusTransitionTest::test_active_to_canceled_deactivates_entitlements
// SubscriptionStatusTransitionTest::test_cancel_from_past_due_deactivates_entitlements
// SubscriptionStatusTransitionTest::test_mark_past_due_is_idempotent
// SubscriptionStatusTransitionTest::test_pause_is_idempotent
func TestMadoraSubscriptionAccessTransitions(t *testing.T) {
	if !billing.StatusActive.GrantsAccess() || !billing.StatusPastDue.GrantsAccess() || !billing.StatusTrialing.GrantsAccess() {
		t.Fatalf("active, past_due, and trialing statuses should grant access")
	}

	if billing.StatusPaused.GrantsAccess() || billing.StatusCanceled.GrantsAccess() {
		t.Fatalf("paused and canceled statuses should not grant access")
	}

	sub := &billing.Subscription{BillableType: "team", BillableID: 10, Status: billing.StatusPastDue, PaddleID: "sub_123"}
	store := &madoraSubStore{subs: []*billing.Subscription{sub}}
	billing := service.NewBillingService(store, madoraOrderStore{}, madoraProductStore{})

	active, err := billing.GetActiveSubscription(context.Background(), "team", 10)

	if err != nil {
		t.Fatalf("active subscription: %v", err)
	}

	if active == nil || active.Provider != "paddle" {
		t.Fatalf("active provider = %#v, want paddle subscription", active)
	}

	if err := billing.CancelSubscription(context.Background(), "team", 10); err != nil {
		t.Fatalf("cancel: %v", err)
	}

	if sub.Status != billing.StatusCanceled || sub.EndsAt == nil || len(store.saved) != 1 {
		t.Fatalf("cancel did not persist canceled status: %#v", sub)
	}

	graceEnds := time.Now().Add(24 * time.Hour)
	sub.EndsAt = &graceEnds

	if err := billing.ResumeSubscription(context.Background(), "team", 10); err != nil {
		t.Fatalf("resume: %v", err)
	}

	if sub.Status != billing.StatusActive || sub.EndsAt != nil {
		t.Fatalf("resume did not reactivate grace-period subscription: %#v", sub)
	}
}

// BillingLifecycleTest::test_begin_checkout_is_idempotent_under_sequential_calls
func TestMadoraCheckoutCreationIsStable(t *testing.T) {
	mgr := testBillingManager()
	billable := madoraBillable{id: 10, typ: "team", seats: 3}
	creator := action.NewSubscriptionCreator(&madoraSubStore{}, mgr)
	plan := billing.NewPlan("Pro", "pri_pro_monthly")

	first, err := creator.Execute(context.Background(), billable, plan, map[string]any{"return_url": "/agreement"})

	if err != nil {
		t.Fatalf("first checkout: %v", err)
	}

	second, err := creator.Execute(context.Background(), billable, plan, map[string]any{"return_url": "/agreement"})

	if err != nil {
		t.Fatalf("second checkout: %v", err)
	}

	firstJSON, _ := json.Marshal(first.ToMap())
	secondJSON, _ := json.Marshal(second.ToMap())

	if !bytes.Equal(firstJSON, secondJSON) {
		t.Fatalf("checkout creation is not stable:\n%s\n%s", firstJSON, secondJSON)
	}

	if first.Items[0].Quantity != 3 || first.Options()["return_url"] != "/agreement" {
		t.Fatalf("checkout = %#v, options = %#v", first.Items, first.Options())
	}
}

// CanonicalMoneySchemaContractTest::test_usd_currency_is_seeded_for_canonical_billing_schema
// CanonicalMoneySchemaContractTest::test_transactions_accept_canonical_transaction_payload
// CanonicalMoneySchemaContractTest::test_transactions_accept_zero_minor_amount_boundaries
// CanonicalMoneySchemaContractTest::test_transactions_reject_negative_total_minor_amounts
// CanonicalMoneySchemaContractTest::test_transactions_reject_negative_tax_minor_amounts
// CanonicalMoneySchemaContractTest::test_transactions_reject_unknown_canonical_currency_codes
// PlanPeriodPriceSchemaContractTest::test_plan_period_prices_use_only_canonical_money_columns
// PlanPeriodPriceSchemaContractTest::test_plan_period_prices_accept_zero_minor_amount_boundaries
// PlanPeriodPriceSchemaContractTest::test_plan_period_prices_reject_negative_minor_amount_boundaries
// PaddleTransactionModelTest::test_transaction_factory_make_normalises_money_fields_without_saving
// PaddleTransactionModelTest::test_transaction_factory_make_honours_minor_amount_overrides
// PaddleTransactionModelTest::test_transaction_factory_make_accepts_signed_string_amounts
// PaddleTransactionModelTest::test_billable_transactions_allow_signed_string_coercion_before_unsigned_schema_rejection
func TestMadoraCanonicalMoneySchema(t *testing.T) {
	if billing.DefaultCurrency != "USD" || !billing.IsSupportedCurrency("usd") {
		t.Fatalf("default currency = %s, supported USD = %v", billing.DefaultCurrency, billing.IsSupportedCurrency("usd"))
	}

	total, err := billing.ParseMinorAmount("12345")

	if err != nil || total != 12345 {
		t.Fatalf("ParseMinorAmount string = %d, %v; want 12345", total, err)
	}

	signed, err := billing.ParseMinorAmount("-42")

	if err != nil || signed != -42 {
		t.Fatalf("ParseMinorAmount signed string = %d, %v; want -42", signed, err)
	}

	if _, err := billing.ParseMinorAmount("12xyz"); err == nil {
		t.Fatalf("malformed minor amount string accepted")
	}

	if err := billing.ValidateTransactionMoney(12345, 123, "USD"); err != nil {
		t.Fatalf("valid transaction money rejected: %v", err)
	}

	if err := billing.ValidateTransactionMoney(0, 0, "USD"); err != nil {
		t.Fatalf("zero transaction money rejected: %v", err)
	}

	if err := billing.ValidateTransactionMoney(-1, 0, "USD"); err == nil {
		t.Fatalf("negative total accepted")
	}

	if err := billing.ValidateTransactionMoney(0, -1, "USD"); err == nil {
		t.Fatalf("negative tax accepted")
	}

	if err := billing.ValidateTransactionMoney(0, 0, "ZZZ"); err == nil {
		t.Fatalf("unknown currency accepted")
	}

	if err := billing.ValidatePriceMoney(0, "USD"); err != nil {
		t.Fatalf("zero price money rejected: %v", err)
	}

	if err := billing.ValidatePriceMoney(-1, "USD"); err == nil {
		t.Fatalf("negative price money accepted")
	}

	tx := billing.Transaction{Total: total, Tax: 123, Currency: "USD"}

	if tx.TotalFormatted() != "USD 123.45" || tx.TaxFormatted() != "USD 1.23" {
		t.Fatalf("formatted transaction totals = %s / %s", tx.TotalFormatted(), tx.TaxFormatted())
	}
}

// BillingLifecycleTest::test_price_rotation_keeps_history_and_only_one_active_price_per_plan_period
// PlanPeriodPriceSchemaContractTest::test_plan_period_prices_use_only_canonical_money_columns
// PlanPeriodPriceSchemaContractTest::test_plan_period_prices_enforce_single_active_guard_per_period
// PlanPeriodPriceSchemaContractTest::test_plan_period_prices_allow_multiple_inactive_versions_per_period
// RotatePlanPricesTest::test_rotate_persists_valid_money_prices_without_legacy_amount_storage
// RotatePlanPricesTest::test_rotate_rejects_incomplete_money_prices_before_persisting_changes
// RotatePlanPricesTest::test_rotate_allows_free_and_custom_prices_without_money_fields
// RotatePlanPricesTest::test_rotate_rejects_duplicate_active_provider_price_ids_on_other_plan_periods
func TestMadoraPlanPriceRotation(t *testing.T) {
	amount := int64(24900)
	rotatedAmount := int64(29900)
	history := []billing.PlanPeriodPrice{{
		ID:              1,
		AmountMinor:     &amount,
		Currency:        "USD",
		PricingMode:     billing.PlanPricingModeMoney,
		ProviderPriceID: "pri_existing",
		Active:          true,
	}}

	if _, ok := reflect.TypeOf(billing.PlanPeriodPrice{}).FieldByName("Amount"); ok {
		t.Fatalf("PlanPeriodPrice should not expose legacy Amount storage")
	}

	rotated, err := billing.RotatePlanPrice(history, billing.PlanPeriodPrice{
		ID:              2,
		AmountMinor:     &rotatedAmount,
		Currency:        "USD",
		PricingMode:     billing.PlanPricingModeMoney,
		Period:          "per month",
		ProviderPriceID: "pri_rotated",
	})

	if err != nil {
		t.Fatalf("rotate valid money price: %v", err)
	}

	if len(rotated) != 2 || rotated[0].Active || !rotated[1].Active {
		t.Fatalf("rotated history = %#v", rotated)
	}

	if rotated[1].DisplayAmount() != "USD 299.00" {
		t.Fatalf("display amount = %s, want USD 299.00", rotated[1].DisplayAmount())
	}

	if _, err := billing.RotatePlanPrice(rotated, billing.PlanPeriodPrice{
		AmountMinor:     &rotatedAmount,
		Currency:        "USD",
		PricingMode:     billing.PlanPricingModeMoney,
		ProviderPriceID: " pri_rotated ",
	}); err == nil {
		t.Fatalf("duplicate active provider price id accepted")
	}

	if _, err := billing.RotatePlanPrice(rotated, billing.PlanPeriodPrice{
		PricingMode:     billing.PlanPricingModeMoney,
		ProviderPriceID: "pri_invalid",
	}); err == nil {
		t.Fatalf("incomplete money price accepted")
	}

	free, err := billing.RotatePlanPrice(nil, billing.PlanPeriodPrice{PricingMode: billing.PlanPricingModeFree})

	if err != nil || free[0].DisplayAmount() != "Free" || free[0].AmountMinor != nil || free[0].Currency != "" {
		t.Fatalf("free price = %#v, err = %v", free, err)
	}

	custom, err := billing.RotatePlanPrice(free, billing.PlanPeriodPrice{PricingMode: billing.PlanPricingModeCustom})

	if err != nil || custom[1].DisplayAmount() != "Custom" || custom[1].AmountMinor != nil || custom[1].Currency != "" {
		t.Fatalf("custom price = %#v, err = %v", custom, err)
	}

	inactiveHistory := []billing.PlanPeriodPrice{
		{ID: 1, ProviderPriceID: "pri_old_one", PricingMode: billing.PlanPricingModeMoney, AmountMinor: &amount, Currency: "USD", Active: false},
		{ID: 2, ProviderPriceID: "pri_old_two", PricingMode: billing.PlanPricingModeMoney, AmountMinor: &amount, Currency: "USD", Active: false},
	}
	withInactive, err := billing.RotatePlanPrice(inactiveHistory, billing.PlanPeriodPrice{
		AmountMinor:     &rotatedAmount,
		Currency:        "USD",
		PricingMode:     billing.PlanPricingModeMoney,
		ProviderPriceID: "pri_new_active",
	})

	if err != nil || len(withInactive) != 3 {
		t.Fatalf("inactive history = %#v, err = %v", withInactive, err)
	}
}

// SyncSubscriptionAfterUpdateTest::test_syncs_activation_state_when_subscription_updated
// SyncSubscriptionAfterUpdateTest::test_dispatches_subscription_changed_event
// SyncSubscriptionAfterUpdateTest::test_updates_local_plan_and_entitlements_when_price_id_changes
// ReconcileSubscriptionAfterCheckoutTest::test_subscription_changed_event_is_dispatched
// SubscriptionPersisterTest::test_create_or_retrieve_skips_items_missing_a_price_payload
// SubscriptionPersisterTest::test_create_or_retrieve_skips_malformed_subscription_items
// SubscriptionPersisterTest::test_create_or_retrieve_defaults_missing_item_quantity_to_one
func TestMadoraWebhookSubscriptionUpdateSyncsLocalState(t *testing.T) {
	sub := &billing.Subscription{
		ID:           99,
		BillableType: "team",
		BillableID:   10,
		PaddleID:     "sub_123",
		Status:       billing.StatusPastDue,
		Items:        []billing.SubscriptionItem{{PriceID: "pri_old", Quantity: 2}},
	}
	subs := &madoraSubStore{subs: []*billing.Subscription{sub}}
	events := &madoraDispatcher{}
	handler := webhook.NewHandler(subs, nil, &madoraTransactionStore{transactions: map[string]*billing.Transaction{}}, events)

	payload := map[string]any{
		"event_type": "subscription.updated",
		"data": map[string]any{
			"id":     "sub_123",
			"status": "active",
			"items": []any{
				map[string]any{"price": map[string]any{"id": "pri_new", "product_id": "pro"}},
				map[string]any{"quantity": 5},
				"malformed",
			},
		},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/paddle/webhook", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	if sub.Status != billing.StatusActive {
		t.Fatalf("subscription status = %s, want active", sub.Status)
	}

	if len(sub.Items) != 1 || sub.Items[0].PriceID != "pri_new" || sub.Items[0].Quantity != 1 {
		t.Fatalf("subscription items = %#v", sub.Items)
	}

	if len(subs.saved) != 1 {
		t.Fatalf("saved subscriptions = %d, want 1", len(subs.saved))
	}

	if !containsEvent[billing.WebhookReceivedEvent](events.events) ||
		!containsEvent[billing.SubscriptionUpdatedEvent](events.events) ||
		!containsEvent[billing.WebhookHandledEvent](events.events) {
		t.Fatalf("events = %#v", events.events)
	}
}

// ReconcileSubscriptionAfterCheckoutTest::test_reconciles_pre_paddle_subscription_to_cashier_subscription
// ReconcileSubscriptionAfterCheckoutTest::test_plan_is_copied_from_old_subscription
// ReconcileSubscriptionAfterCheckoutTest::test_subscription_changed_event_is_dispatched
// ReconcileSubscriptionAfterCheckoutTest::test_graceful_noop_when_no_pre_paddle_subscription_exists
func TestMadoraReconcileSubscriptionAfterCheckout(t *testing.T) {
	pending := &billing.Subscription{
		ID:           10,
		BillableType: "team",
		BillableID:   42,
		Type:         billing.DefaultSubscriptionType,
		Plan:         "pro",
		Status:       billing.StatusAwaitingPayment,
		Items:        []billing.SubscriptionItem{{PriceID: "pri_pro_monthly"}},
	}
	cashier := &billing.Subscription{
		ID:           11,
		BillableType: "team",
		BillableID:   42,
		PaddleID:     "sub_real_paddle_id",
		Status:       billing.StatusActive,
		Items:        []billing.SubscriptionItem{{PriceID: "pri_pro_monthly"}},
	}
	store := &madoraSubStore{subs: []*billing.Subscription{pending, cashier}}
	events := &madoraDispatcher{}

	reconciled, err := billing.ReconcileSubscriptionAfterCheckout(context.Background(), store, "team", 42, cashier, events)

	if err != nil {
		t.Fatalf("reconcile checkout: %v", err)
	}

	if !reconciled {
		t.Fatalf("expected checkout reconciliation")
	}

	if cashier.Plan != "pro" || cashier.Type != billing.DefaultSubscriptionType {
		t.Fatalf("cashier metadata = plan %q type %q, want pro/default", cashier.Plan, cashier.Type)
	}

	if len(store.saved) != 1 || store.saved[0].ID != cashier.ID {
		t.Fatalf("saved subscriptions = %#v, want cashier", store.saved)
	}

	found, err := store.FindByID(context.Background(), pending.ID)

	if err != nil {
		t.Fatalf("find pending: %v", err)
	}

	if found != nil {
		t.Fatalf("pending subscription was not removed: %#v", found)
	}

	if !containsEvent[billing.SubscriptionUpdatedEvent](events.events) {
		t.Fatalf("events = %#v, want subscription update", events.events)
	}

	noPendingStore := &madoraSubStore{subs: []*billing.Subscription{cashier}}
	reconciled, err = billing.ReconcileSubscriptionAfterCheckout(context.Background(), noPendingStore, "team", 42, cashier, events)

	if err != nil {
		t.Fatalf("noop reconcile: %v", err)
	}

	if reconciled {
		t.Fatalf("reconciled without a pre-paddle subscription")
	}
}

// PaddleTransactionModelTest::test_transaction_factory_create_persists_a_valid_billable_transaction
// PaddleTransactionModelTest::test_billable_transactions_backfill_minor_amounts_on_create
func TestMadoraWebhookTransactionCompletedPersistsTransaction(t *testing.T) {
	transactions := &madoraTransactionStore{transactions: map[string]*billing.Transaction{}}
	events := &madoraDispatcher{}
	handler := webhook.NewHandler(&madoraSubStore{}, nil, transactions, events)

	payload := map[string]any{
		"event_type": "transaction.completed",
		"data": map[string]any{
			"id":     "txn_123",
			"status": "completed",
			"details": map[string]any{
				"totals": map[string]any{
					"total":         "2500",
					"tax":           "250",
					"currency_code": "USD",
				},
			},
		},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/paddle/webhook", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	if len(transactions.created) != 1 {
		t.Fatalf("created transactions = %d, want 1", len(transactions.created))
	}

	tx := transactions.created[0]

	if tx.PaddleID != "txn_123" || tx.Status != billing.TransactionCompleted || tx.Total != 2500 || tx.Tax != 250 || tx.Currency != "USD" {
		t.Fatalf("transaction = %#v", tx)
	}

	if !containsEvent[billing.TransactionCompletedEvent](events.events) {
		t.Fatalf("events = %#v", events.events)
	}
}

// CanonicalMoneySchemaContractTest::test_transactions_require_explicit_minor_amounts_in_canonical_schema
// PaddleTransactionModelTest::test_billable_transactions_reject_malformed_minor_amount_strings_on_create
// PaddleTransactionModelTest::test_billable_transactions_reject_malformed_minor_amount_strings_on_update
func TestMadoraWebhookTransactionRejectsInvalidMoneyPayloads(t *testing.T) {
	tests := []struct {
		name    string
		payload map[string]any
	}{
		{
			name: "missing totals",
			payload: map[string]any{
				"event_type": "transaction.completed",
				"data": map[string]any{
					"id":            "txn_missing_totals",
					"status":        "completed",
					"currency_code": "SGD",
				},
			},
		},
		{
			name: "malformed total",
			payload: map[string]any{
				"event_type": "transaction.completed",
				"data": map[string]any{
					"id":            "txn_malformed_total",
					"status":        "completed",
					"currency_code": "SGD",
					"details": map[string]any{
						"totals": map[string]any{
							"total": "29xx",
							"tax":   "300",
						},
					},
				},
			},
		},
		{
			name: "unknown currency",
			payload: map[string]any{
				"event_type": "transaction.completed",
				"data": map[string]any{
					"id":            "txn_unknown_currency",
					"status":        "completed",
					"currency_code": "ZZZ",
					"details": map[string]any{
						"totals": map[string]any{
							"total": "2900",
							"tax":   "300",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactions := &madoraTransactionStore{transactions: map[string]*billing.Transaction{}}
			events := &madoraDispatcher{}
			handler := webhook.NewHandler(&madoraSubStore{}, nil, transactions, events)

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/paddle/webhook", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.Handle(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rec.Code)
			}

			if len(transactions.created) != 0 {
				t.Fatalf("created transactions = %#v, want none", transactions.created)
			}

			if containsEvent[billing.TransactionCompletedEvent](events.events) {
				t.Fatalf("transaction completion event dispatched for invalid payload: %#v", events.events)
			}
		})
	}

	existing := &billing.Transaction{PaddleID: "txn_update_invalid", Status: billing.TransactionReady, Total: 2500, Tax: 250, Currency: "SGD"}
	transactions := &madoraTransactionStore{transactions: map[string]*billing.Transaction{"txn_update_invalid": existing}}
	events := &madoraDispatcher{}
	handler := webhook.NewHandler(&madoraSubStore{}, nil, transactions, events)
	payload := map[string]any{
		"event_type": "transaction.updated",
		"data": map[string]any{
			"id":            "txn_update_invalid",
			"status":        "completed",
			"currency_code": "SGD",
			"details": map[string]any{
				"totals": map[string]any{
					"total": "3000",
					"tax":   "bad-tax",
				},
			},
		},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/paddle/webhook", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	if len(transactions.saved) != 0 {
		t.Fatalf("saved transactions = %#v, want none", transactions.saved)
	}

	if existing.Status != billing.TransactionReady || existing.Total != 2500 || existing.Tax != 250 || existing.Currency != "SGD" {
		t.Fatalf("existing transaction changed after invalid update: %#v", existing)
	}
}

// PaddleTransactionModelTest::test_billable_transactions_refresh_minor_amounts_on_update
func TestMadoraWebhookTransactionUpdatedRefreshesTransactionMoney(t *testing.T) {
	existing := &billing.Transaction{PaddleID: "txn_123", Status: billing.TransactionReady, Total: 2500, Tax: 250, Currency: "USD"}
	transactions := &madoraTransactionStore{transactions: map[string]*billing.Transaction{"txn_123": existing}}
	events := &madoraDispatcher{}
	handler := webhook.NewHandler(&madoraSubStore{}, nil, transactions, events)

	payload := map[string]any{
		"event_type": "transaction.updated",
		"data": map[string]any{
			"id":     "txn_123",
			"status": "completed",
			"details": map[string]any{
				"totals": map[string]any{
					"total":         "3000",
					"tax":           "0",
					"currency_code": "USD",
				},
			},
		},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/paddle/webhook", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Handle(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}

	if len(transactions.created) != 0 || len(transactions.saved) != 1 {
		t.Fatalf("created/saved transactions = %d/%d, want 0/1", len(transactions.created), len(transactions.saved))
	}

	if existing.Status != billing.TransactionCompleted || existing.Total != 3000 || existing.Tax != 0 || existing.Currency != "USD" {
		t.Fatalf("updated transaction = %#v", existing)
	}

	if !containsEvent[billing.TransactionCompletedEvent](events.events) {
		t.Fatalf("events = %#v", events.events)
	}
}

func containsEvent[T any](events []any) bool {
	for _, event := range events {
		if _, ok := event.(T); ok {
			return true
		}
	}

	return false
}
