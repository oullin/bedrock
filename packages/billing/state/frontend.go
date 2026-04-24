// Package state builds the billing portal frontend state.
package state

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/bedrock/packages/billing"
)

// FrontendState builds the data shared with the billing portal frontend.
// Mirrors Billing\FrontendState.
type FrontendState struct {
	manager       *billing.Manager
	config        *billing.Config
	subscriptions billing.SubscriptionStore
	customers     billing.CustomerStore
}

// NewFrontendState creates a FrontendState builder.
func NewFrontendState(mgr *billing.Manager, cfg *billing.Config, subs billing.SubscriptionStore) *FrontendState {
	return &FrontendState{
		manager:       mgr,
		config:        cfg,
		subscriptions: subs,
	}
}

// WithCustomerStore enables customer-backed portal state, including pending
// checkout detection.
func (f *FrontendState) WithCustomerStore(customers billing.CustomerStore) *FrontendState {
	f.customers = customers

	return f
}

// Current returns the full frontend state for the billing portal.
func (f *FrontendState) Current(ctx context.Context, billableType string, billable billing.Billable) (map[string]any, error) {
	return f.CurrentAt(ctx, billableType, billable, time.Now())
}

// CurrentAt returns the full frontend state using an explicit clock. It is
// useful for billing DTOs that expose relative pending-window values.
func (f *FrontendState) CurrentAt(ctx context.Context, billableType string, billable billing.Billable, now time.Time) (map[string]any, error) {
	sub, _ := f.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())
	customer, err := f.customerForBillable(ctx, billable)

	if err != nil {
		return nil, err
	}

	plans := f.manager.Plans(billableType)

	var monthlyPlans, yearlyPlans []*billing.Plan

	for _, p := range plans {
		if !p.Active {
			continue
		}

		switch p.Interval {
		case "monthly":
			monthlyPlans = append(monthlyPlans, p)
		case "yearly":
			yearlyPlans = append(yearlyPlans, p)
		}
	}

	activePlan := planForSubscription(plans, sub)

	state, err := f.resolveState(ctx, sub, customer)

	if err != nil {
		return nil, err
	}

	subscription := f.subscriptionState(sub, activePlan)
	cta := f.ctaState(sub, activePlan, now)

	data := map[string]any{
		"billableId":      billable.BillableID(),
		"billableName":    billable.BillableName(),
		"billableType":    billableType,
		"brandColor":      f.brandColor(),
		"dashboardUrl":    f.dashboardURL(),
		"defaultInterval": f.defaultInterval(billableType),
		"monthlyPlans":    monthlyPlans,
		"yearlyPlans":     yearlyPlans,
		"plan":            activePlan,
		"seatName":        f.manager.SeatName(billableType),
		"sparkPath":       f.config.Path,
		"state":           state,
		"subscription":    subscription,
		"cta":             cta,
		"termsUrl":        f.config.TermsURL,
	}

	return data, nil
}

func (f *FrontendState) customerForBillable(ctx context.Context, billable billing.Billable) (*billing.Customer, error) {
	if f.customers == nil {
		return nil, nil
	}

	return f.customers.FindByBillable(ctx, billable.BillableType(), billable.BillableID())
}

func planForSubscription(plans []*billing.Plan, sub *billing.Subscription) *billing.Plan {
	if sub == nil {
		return nil
	}

	for _, p := range plans {
		if sub.HasPrice(p.ID) {
			return p
		}
	}

	for _, p := range plans {
		if slug, _ := p.Options["slug"].(string); slug != "" && slug == sub.Plan {
			return p
		}

		if p.Name == sub.Plan {
			return p
		}
	}

	return nil
}

func (f *FrontendState) resolveState(ctx context.Context, sub *billing.Subscription, customer *billing.Customer) (string, error) {
	if subscriptionIsActiveOrPastDue(sub) {
		if customer != nil && customer.PendingCheckoutID != "" {
			customer.PendingCheckoutID = ""

			if err := f.customers.Save(ctx, customer); err != nil {
				return "", err
			}
		}
	} else if customer != nil && customer.PendingCheckoutID != "" {
		return "pending", nil
	}

	if sub == nil {
		return "none", nil
	}

	if sub.OnGracePeriod() {
		return "onGracePeriod", nil
	}

	if sub.Active() || sub.OnTrial() {
		return "active", nil
	}

	if sub.PastDue() {
		return "past_due", nil
	}

	return "none", nil
}

func subscriptionIsActiveOrPastDue(sub *billing.Subscription) bool {
	if sub == nil {
		return false
	}

	return sub.Active() || sub.OnTrial() || sub.PastDue()
}

func (f *FrontendState) brandColor() string {
	if f.config.BrandColor != "" {
		return f.config.BrandColor
	}

	return "bg-gray-800"
}

func (f *FrontendState) dashboardURL() string {
	if f.config.DashboardURL != "" {
		return f.config.DashboardURL
	}

	return "/"
}

func (f *FrontendState) defaultInterval(billableType string) string {
	if f.config.Billables != nil {
		if cfg, ok := f.config.Billables[billableType]; ok && cfg.DefaultInterval != "" {
			return cfg.DefaultInterval
		}
	}

	return "monthly"
}

func (f *FrontendState) subscriptionState(sub *billing.Subscription, plan *billing.Plan) map[string]any {
	state := map[string]any{
		"status":             "",
		"plan_code":          "",
		"plan_name":          "",
		"pending_expires_at": (*time.Time)(nil),
		"payment_ready_at":   (*time.Time)(nil),
		"portal_url":         "",
		"pay_now":            false,
	}

	if sub == nil {
		return state
	}

	state["status"] = string(sub.Status)
	state["plan_code"] = sub.Plan

	if plan != nil {
		state["plan_name"] = plan.Name
	}

	state["pending_expires_at"] = sub.PendingExpiresAt
	state["payment_ready_at"] = sub.PaymentReadyAt

	portalURL := f.portalURL(sub, plan)
	state["portal_url"] = portalURL
	state["pay_now"] = sub.Status == billing.StatusAwaitingPayment && portalURL != ""

	return state
}

func (f *FrontendState) ctaState(sub *billing.Subscription, plan *billing.Plan, now time.Time) map[string]any {
	cta := map[string]any{
		"visible":        false,
		"label":          "",
		"remaining_days": 0,
	}

	if sub == nil || plan == nil || sub.Status != billing.StatusAwaitingPayment {
		return cta
	}

	portalURL := f.portalURL(sub, plan)

	if portalURL == "" {
		return cta
	}

	cta["visible"] = true
	cta["label"] = "Complete Subscription"
	cta["remaining_days"] = remainingDays(now, sub.PendingExpiresAt)

	return cta
}

func (f *FrontendState) portalURL(sub *billing.Subscription, plan *billing.Plan) string {
	if sub == nil || plan == nil {
		return ""
	}

	switch sub.Status {
	case billing.StatusAwaitingPayment:
		if !sub.HasPrice(plan.ID) {
			return ""
		}
	case billing.StatusActive, billing.StatusPastDue, billing.StatusPaused:
	default:
		return ""
	}

	path := strings.Trim(f.config.Path, "/")

	if path == "" {
		return "/"
	}

	return "/" + path
}

func remainingDays(now time.Time, expiresAt *time.Time) int {
	if expiresAt == nil || !expiresAt.After(now) {
		return 0
	}

	return int(math.Ceil(expiresAt.Sub(now).Hours() / 24))
}
