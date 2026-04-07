package billing

import "context"

// FrontendStateView is the complete state object served to the billing portal.
type FrontendStateView struct {
	AppName     string         `json:"app_name"`
	BillableID  int64          `json:"billable_id"`
	BillableName string        `json:"billable_name"`
	BillableType string        `json:"billable_type"`
	Plans       []BillingPlan    `json:"plans"`
	State       string         `json:"state"` // "pending", "onGracePeriod", "past_due", "active", "none"
	DashboardURL string        `json:"dashboard_url"`
	TermsURL    string         `json:"terms_url"`
	Sandbox     bool           `json:"sandbox"`
	DateFormat  string         `json:"date_format"`
}

// FrontendStateBuilder builds the frontend state for the billing portal.
type FrontendStateBuilder struct {
	manager       *Manager
	config        *Config
	subscriptions SubscriptionStore
	customers     CustomerStore
	pricePreviewer ProviderPricePreviewer
	clock         Clock
}

// NewFrontendStateBuilder creates a new builder.
func NewFrontendStateBuilder(
	manager *Manager,
	config *Config,
	subscriptions SubscriptionStore,
	customers CustomerStore,
	pricePreviewer ProviderPricePreviewer,
	clock Clock,
) *FrontendStateBuilder {
	return &FrontendStateBuilder{
		manager:        manager,
		config:         config,
		subscriptions:  subscriptions,
		customers:      customers,
		pricePreviewer: pricePreviewer,
		clock:          clock,
	}
}

// Current builds the complete frontend state for a billable entity.
func (f *FrontendStateBuilder) Current(ctx context.Context, billableType string, billable Billable) (*FrontendStateView, error) {
	plans := f.manager.Plans(billableType)

	sub, _ := f.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())
	state := f.resolveState(ctx, billable, sub)

	return &FrontendStateView{
		AppName:      f.config.BrandLogoPath,
		BillableID:   billable.BillableID(),
		BillableName: billable.BillableName(),
		BillableType: billableType,
		Plans:        plans,
		State:        state,
		DashboardURL: f.config.DashboardURL,
		TermsURL:     f.config.TermsURL,
		Sandbox:      f.config.Sandbox,
		DateFormat:   f.config.DateFormat,
	}, nil
}

func (f *FrontendStateBuilder) resolveState(ctx context.Context, billable Billable, sub *Subscription) string {
	if sub == nil {
		// Check for pending checkout.
		customer, err := f.customers.FindByBillable(ctx, billable.BillableType(), billable.BillableID())
		if err == nil && customer != nil && customer.PendingCheckout != nil {
			return "pending"
		}

		return "none"
	}

	if sub.Canceled() && sub.OnGracePeriod(f.clock) {
		return "onGracePeriod"
	}

	if sub.PastDue() {
		return "past_due"
	}

	if sub.Active() || sub.OnTrial(f.clock) {
		return "active"
	}

	return "none"
}
