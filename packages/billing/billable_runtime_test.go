package spark_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/billing"
)

type runtimeBillable struct {
	id       int64
	btype    string
	name     string
	email    string
	provider string
}

type fakeRuntimeCustomerStore struct {
	findResult *billing.Customer
	findErr    error
	createErr  error
	saveErr    error
	created    *billing.Customer
	saved      *billing.Customer
}

type fakeProviderOps struct {
	customer  *billing.Customer
	createErr error
	received  billing.CustomerCreateOptions
}

type fakeQuantityUpdater struct {
	err      error
	called   bool
	subArg   *billing.Subscription
	quantity int
	behavior billing.ProrationBehavior
}

func (b *runtimeBillable) BillableID() int64       { return b.id }
func (b *runtimeBillable) BillableType() string    { return b.btype }
func (b *runtimeBillable) BillableName() string    { return b.name }
func (b *runtimeBillable) BillableEmail() string   { return b.email }
func (b *runtimeBillable) PaymentProvider() string { return b.provider }
func (b *runtimeBillable) SetPaymentProvider(p string) error {
	b.provider = p

	return nil
}

func (f *fakeRuntimeCustomerStore) FindByProviderID(context.Context, string) (*billing.Customer, error) {
	return nil, nil
}
func (f *fakeRuntimeCustomerStore) FindByBillable(context.Context, string, int64) (*billing.Customer, error) {
	return f.findResult, f.findErr
}
func (f *fakeRuntimeCustomerStore) Create(_ context.Context, c *billing.Customer) error {
	f.created = c

	return f.createErr
}
func (f *fakeRuntimeCustomerStore) Save(_ context.Context, c *billing.Customer) error {
	f.saved = c

	return f.saveErr
}

func (f *fakeProviderOps) CreateCustomer(_ context.Context, _ billing.Billable, opts billing.CustomerCreateOptions) (*billing.Customer, error) {
	f.received = opts

	return f.customer, f.createErr
}
func (f *fakeProviderOps) PreviewPrices(context.Context, []string, map[string]any) ([]billing.PricePreview, error) {
	return nil, nil
}
func (f *fakeProviderOps) CreateCheckoutSession(context.Context, billing.Billable, *billing.Checkout, map[string]any) (*billing.Checkout, error) {
	return nil, nil
}
func (f *fakeProviderOps) CreatePaymentMethodUpdateTransaction(context.Context, billing.Billable, *billing.Subscription, map[string]any) (*billing.PaymentMethodUpdateTransaction, error) {
	return nil, nil
}
func (f *fakeProviderOps) DownloadInvoice(context.Context, *billing.Transaction) (*billing.InvoiceDownload, error) {
	return nil, nil
}
func (f *fakeProviderOps) LastPayment(context.Context, *billing.Subscription) (*billing.Payment, error) {
	return nil, nil
}
func (f *fakeProviderOps) NextPayment(context.Context, *billing.Subscription) (*billing.Payment, error) {
	return nil, nil
}
func (f *fakeProviderOps) UpdateSubscriptionQuantity(context.Context, *billing.Subscription, int, billing.ProrationBehavior) error {
	return nil
}

func (f *fakeQuantityUpdater) UpdateSubscriptionQuantity(_ context.Context, sub *billing.Subscription, q int, b billing.ProrationBehavior) error {
	f.called = true
	f.subArg = sub
	f.quantity = q
	f.behavior = b

	return f.err
}

func TestSeatProration(t *testing.T) {
	// Can't test directly since it's unexported from inside the package_test,
	// but we exercise it via AddSeats.
	updater := &fakeQuantityUpdater{}
	sub := &billing.Subscription{Items: []billing.SubscriptionItem{{Quantity: 2}}}

	if err := billing.AddSeats(context.Background(), updater, sub, 1, true); err != nil {
		t.Fatal(err)
	}

	if updater.behavior != billing.ProrateNextBilling {
		t.Fatalf("prorates=true -> %q, want ProrateNextBilling", updater.behavior)
	}

	updater2 := &fakeQuantityUpdater{}

	if err := billing.AddSeats(context.Background(), updater2, sub, 1, false); err != nil {
		t.Fatal(err)
	}

	if updater2.behavior != billing.FullNextBilling {
		t.Fatalf("prorates=false -> %q, want FullNextBilling", updater2.behavior)
	}
}

func TestBillingPlan(t *testing.T) {
	mgr := billing.NewManager()
	plan := mgr.Plan("team", "Pro", "pri_team")

	billable := &runtimeBillable{id: 1, btype: "team"}

	if got := billing.BillingPlan(mgr, billable, nil); got != nil {
		t.Fatal("nil subscription should return nil")
	}

	canceled := &billing.Subscription{Status: billing.StatusCanceled}

	if got := billing.BillingPlan(mgr, billable, canceled); got != nil {
		t.Fatal("invalid subscription should return nil")
	}

	sub := &billing.Subscription{
		Status: billing.StatusActive,
		Items:  []billing.SubscriptionItem{{PriceID: "pri_team"}},
	}

	if got := billing.BillingPlan(mgr, billable, sub); got != plan {
		t.Fatalf("want matching plan, got %v", got)
	}

	sub.Items[0].PriceID = "other"

	if got := billing.BillingPlan(mgr, billable, sub); got != nil {
		t.Fatal("no matching price should return nil")
	}
}

func TestCanUpdateSeats(t *testing.T) {
	future := time.Now().Add(48 * time.Hour)

	tests := []struct {
		name     string
		customer *billing.Customer
		sub      *billing.Subscription
		want     bool
	}{
		{name: "generic trial", customer: &billing.Customer{TrialEndsAt: &future}, sub: nil, want: true},
		{name: "nil customer and nil sub", customer: nil, sub: nil, want: false},
		{name: "past due sub", customer: nil, sub: &billing.Subscription{Status: billing.StatusPastDue}, want: false},
		{name: "trial sub not recurring", customer: nil, sub: &billing.Subscription{Status: billing.StatusTrialing}, want: false},
		{name: "active recurring", customer: nil, sub: &billing.Subscription{Status: billing.StatusActive}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := billing.CanUpdateSeats(tt.customer, tt.sub); got != tt.want {
				t.Fatalf("CanUpdateSeats = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateSeats_NilSubscription(t *testing.T) {
	updater := &fakeQuantityUpdater{}

	if err := billing.UpdateSeats(context.Background(), updater, nil, 5, false); err != nil {
		t.Fatal(err)
	}

	if updater.called {
		t.Fatal("updater should not be called for nil sub")
	}
}

func TestUpdateSeats_FloorsAtOne(t *testing.T) {
	updater := &fakeQuantityUpdater{}
	sub := &billing.Subscription{}

	if err := billing.UpdateSeats(context.Background(), updater, sub, 0, false); err != nil {
		t.Fatal(err)
	}

	if updater.quantity != 1 {
		t.Fatalf("quantity = %d, want floored to 1", updater.quantity)
	}
}

func TestAddSeats(t *testing.T) {
	updater := &fakeQuantityUpdater{}
	sub := &billing.Subscription{Items: []billing.SubscriptionItem{{Quantity: 2}}}

	if err := billing.AddSeats(context.Background(), updater, sub, 3, false); err != nil {
		t.Fatal(err)
	}

	if updater.quantity != 5 {
		t.Fatalf("quantity = %d, want 5", updater.quantity)
	}
}

func TestAddSeats_NilSubscription(t *testing.T) {
	updater := &fakeQuantityUpdater{}

	if err := billing.AddSeats(context.Background(), updater, nil, 3, false); err != nil {
		t.Fatal(err)
	}

	if updater.called {
		t.Fatal("updater should not be called")
	}
}

func TestRemoveSeats(t *testing.T) {
	updater := &fakeQuantityUpdater{}
	sub := &billing.Subscription{Items: []billing.SubscriptionItem{{Quantity: 4}}}

	if err := billing.RemoveSeats(context.Background(), updater, sub, 2, false); err != nil {
		t.Fatal(err)
	}

	if updater.quantity != 2 {
		t.Fatalf("quantity = %d, want 2", updater.quantity)
	}
}

func TestRemoveSeats_FloorsAtOne(t *testing.T) {
	updater := &fakeQuantityUpdater{}
	sub := &billing.Subscription{Items: []billing.SubscriptionItem{{Quantity: 1}}}

	if err := billing.RemoveSeats(context.Background(), updater, sub, 5, false); err != nil {
		t.Fatal(err)
	}

	if updater.quantity != 1 {
		t.Fatalf("quantity = %d, want 1 (floor)", updater.quantity)
	}
}

func TestRemoveSeats_NegativeCountClampsToOne(t *testing.T) {
	updater := &fakeQuantityUpdater{}
	sub := &billing.Subscription{Items: []billing.SubscriptionItem{{Quantity: 5}}}

	if err := billing.RemoveSeats(context.Background(), updater, sub, -3, false); err != nil {
		t.Fatal(err)
	}

	if updater.quantity != 4 {
		t.Fatalf("quantity = %d, want 4 (count floored to 1, 5-1=4)", updater.quantity)
	}
}

func TestUpdateSeats_PropagatesError(t *testing.T) {
	updater := &fakeQuantityUpdater{err: errors.New("boom")}
	sub := &billing.Subscription{}

	if err := billing.UpdateSeats(context.Background(), updater, sub, 3, false); err == nil {
		t.Fatal("expected error")
	}
}

func TestBootstrapCustomer_CreatesWhenMissing(t *testing.T) {
	store := &fakeRuntimeCustomerStore{}
	provider := &fakeProviderOps{customer: &billing.Customer{PaddleID: "cus_1"}}
	billable := &runtimeBillable{id: 1, btype: "team", name: "Acme", email: "a@b.c"}
	now := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)

	customer, err := billing.BootstrapCustomer(context.Background(), store, provider, billable, 7, now)

	if err != nil {
		t.Fatalf("BootstrapCustomer: %v", err)
	}

	if customer == nil || customer.PaddleID != "cus_1" {
		t.Fatalf("customer = %#v", customer)
	}

	if store.created == nil {
		t.Fatal("Create not called")
	}

	if store.created.BillableType != "team" || store.created.Name != "Acme" || store.created.Email != "a@b.c" {
		t.Fatalf("billable fields not set: %#v", store.created)
	}

	if customer.TrialEndsAt == nil || !customer.TrialEndsAt.Equal(now.AddDate(0, 0, 7)) {
		t.Fatalf("trial end = %v, want now+7d", customer.TrialEndsAt)
	}
}

func TestBootstrapCustomer_UpdatesWhenExists(t *testing.T) {
	existing := &billing.Customer{ID: 42, PaddleID: "cus_old"}
	store := &fakeRuntimeCustomerStore{findResult: existing}
	provider := &fakeProviderOps{}
	billable := &runtimeBillable{id: 1, btype: "team", name: "Renamed", email: "r@b.c"}

	customer, err := billing.BootstrapCustomer(context.Background(), store, provider, billable, 0, time.Now())

	if err != nil {
		t.Fatalf("BootstrapCustomer: %v", err)
	}

	if customer != existing {
		t.Fatal("should return existing customer reference")
	}

	if store.saved == nil {
		t.Fatal("Save not called")
	}

	if store.saved.Name != "Renamed" || store.saved.Email != "r@b.c" {
		t.Fatalf("update fields not set: %#v", store.saved)
	}

	if store.saved.TrialEndsAt != nil {
		t.Fatal("trial should be nil when days=0")
	}

	if store.created != nil {
		t.Fatal("Create should not have been called")
	}
}

func TestBootstrapCustomer_FindError(t *testing.T) {
	store := &fakeRuntimeCustomerStore{findErr: errors.New("boom")}
	provider := &fakeProviderOps{}
	_, err := billing.BootstrapCustomer(context.Background(), store, provider, &runtimeBillable{}, 0, time.Now())

	if err == nil {
		t.Fatal("expected find error")
	}
}

func TestBootstrapCustomer_CreateCustomerError(t *testing.T) {
	store := &fakeRuntimeCustomerStore{}
	provider := &fakeProviderOps{createErr: errors.New("boom")}
	_, err := billing.BootstrapCustomer(context.Background(), store, provider, &runtimeBillable{}, 0, time.Now())

	if err == nil {
		t.Fatal("expected create error")
	}
}

func TestBootstrapCustomer_StoreCreateError(t *testing.T) {
	store := &fakeRuntimeCustomerStore{createErr: errors.New("boom")}
	provider := &fakeProviderOps{customer: &billing.Customer{}}
	_, err := billing.BootstrapCustomer(context.Background(), store, provider, &runtimeBillable{}, 0, time.Now())

	if err == nil {
		t.Fatal("expected store create error")
	}
}

func TestBootstrapCustomer_StoreSaveError(t *testing.T) {
	existing := &billing.Customer{ID: 1}
	store := &fakeRuntimeCustomerStore{findResult: existing, saveErr: errors.New("boom")}
	provider := &fakeProviderOps{}
	_, err := billing.BootstrapCustomer(context.Background(), store, provider, &runtimeBillable{}, 0, time.Now())

	if err == nil {
		t.Fatal("expected save error")
	}
}
