package billing

// BillingPageView is the top-level data passed to the billing page.
type BillingPageView struct {
	Team         BillingTeamView         `json:"team"`
	SelectedPlan BillingSelectedPlanView `json:"selected_plan"`
	Subscription BillingStateView        `json:"subscription"`
}

// BillingTeamView holds team information for the billing page.
type BillingTeamView struct {
	ID   int64  `json:"id"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// BillingSelectedPlanView holds the currently selected plan display data.
type BillingSelectedPlanView struct {
	Code          string            `json:"code"`
	Name          string            `json:"name"`
	Period        string            `json:"period"`
	Price         *BillingPriceView `json:"price"`
	DisplayAmount string            `json:"display_amount"`
}

// BillingPriceView holds price information for a plan period.
type BillingPriceView struct {
	AmountMinor     *int64  `json:"amount_minor"`
	Currency        *string `json:"currency"`
	PricingMode     string  `json:"pricing_mode"`
	ProviderPriceID string  `json:"provider_price_id"`
}

// BillingStateView is the subscription state passed to the frontend.
type BillingStateView struct {
	Subscription SubscriptionStateView `json:"subscription"`
	CTA          SubscriptionCTAView   `json:"cta"`
}

// SubscriptionStateView is the frontend representation of the subscription state.
type SubscriptionStateView struct {
	Status           string  `json:"status"`
	PlanCode         string  `json:"plan_code"`
	PlanName         string  `json:"plan_name"`
	BillingPeriod    string  `json:"billing_period"`
	PayNow           bool    `json:"pay_now"`
	PendingExpiresAt *string `json:"pending_expires_at"`
	PaymentReadyAt   *string `json:"payment_ready_at"`
	PortalURL        *string `json:"portal_url"`
}

// SubscriptionCTAView is the frontend representation of the subscription CTA.
type SubscriptionCTAView struct {
	Visible          bool    `json:"visible"`
	Href             *string `json:"href"`
	Label            string  `json:"label"`
	PlanName         string  `json:"plan_name"`
	Status           *string `json:"status"`
	PendingExpiresAt *string `json:"pending_expires_at"`
	RemainingDays    *int    `json:"remaining_days"`
}

// TransactionView is the frontend representation of a transaction.
type TransactionView struct {
	ID            int64   `json:"id"`
	TotalMinor    *int64  `json:"total_minor"`
	TaxMinor      *int64  `json:"tax_minor"`
	Currency      *string `json:"currency"`
	DisplayTotal  *string `json:"display_total"`
	DisplayTax    *string `json:"display_tax"`
	BilledAt      string  `json:"billed_at"`
	InvoiceNumber string  `json:"invoice_number"`
	Status        string  `json:"status"`
}

// LandingPlanView is the plan data shown on marketing/landing pages.
type LandingPlanView struct {
	Code              string           `json:"code"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	MarketingFeatures []string         `json:"marketing_features"`
	CTALabel          string           `json:"cta_label"`
	CTAStyle          string           `json:"cta_style"`
	Featured          bool             `json:"featured"`
	Badge             string           `json:"badge"`
	Periods           []PlanPeriodView `json:"periods"`
}

// PlanPeriodView is a billing period option for a plan.
type PlanPeriodView struct {
	Period        string `json:"period"`
	DisplayLabel  string `json:"display_label"`
	DisplayAmount string `json:"display_amount"`
	PricingMode   string `json:"pricing_mode"`
}

// FrontendPlanView is a plan entry for the billing portal's plan selector.
type FrontendPlanView struct {
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Features    []PlanFeatureView `json:"features"`
	CTALabel    string            `json:"cta_label"`
	CTAStyle    string            `json:"cta_style"`
	Featured    bool              `json:"featured"`
	Badge       string            `json:"badge"`
	Periods     []PlanPeriodView  `json:"periods"`
}

// PlanFeatureView is a single feature listed on a plan card.
type PlanFeatureView struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Value *string `json:"value"`
}
