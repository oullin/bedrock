// Package state builds the billing portal frontend state.
package state

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/bedrock/packages/spark"
)

// FrontendState builds the data shared with the billing portal frontend.
// Mirrors Spark\FrontendState.
type FrontendState struct {
	manager       *spark.Manager
	config        *spark.Config
	subscriptions spark.SubscriptionStore
}

// NewFrontendState creates a FrontendState builder.
func NewFrontendState(mgr *spark.Manager, cfg *spark.Config, subs spark.SubscriptionStore) *FrontendState {
	return &FrontendState{
		manager:       mgr,
		config:        cfg,
		subscriptions: subs,
	}
}

// Current returns the full frontend state for the billing portal.
func (f *FrontendState) Current(ctx context.Context, billableType string, billable spark.Billable) (map[string]any, error) {
	return f.CurrentAt(ctx, billableType, billable, time.Now())
}

// CurrentAt returns the full frontend state using an explicit clock. It is
// useful for billing DTOs that expose relative pending-window values.
func (f *FrontendState) CurrentAt(ctx context.Context, billableType string, billable spark.Billable, now time.Time) (map[string]any, error) {
	sub, _ := f.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())

	plans := f.manager.Plans(billableType)

	var monthlyPlans, yearlyPlans []*spark.Plan

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

	state := resolveState(sub)
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

func planForSubscription(plans []*spark.Plan, sub *spark.Subscription) *spark.Plan {
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

func resolveState(sub *spark.Subscription) string {
	if sub == nil {
		return "none"
	}

	if sub.OnGracePeriod() {
		return "onGracePeriod"
	}

	if sub.Active() || sub.OnTrial() {
		return "active"
	}

	if sub.PastDue() {
		return "past_due"
	}

	return "none"
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

func (f *FrontendState) subscriptionState(sub *spark.Subscription, plan *spark.Plan) map[string]any {
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
	state["pay_now"] = sub.Status == spark.StatusAwaitingPayment && portalURL != ""

	return state
}

func (f *FrontendState) ctaState(sub *spark.Subscription, plan *spark.Plan, now time.Time) map[string]any {
	cta := map[string]any{
		"visible":        false,
		"label":          "",
		"remaining_days": 0,
	}

	if sub == nil || plan == nil || sub.Status != spark.StatusAwaitingPayment {
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

func (f *FrontendState) portalURL(sub *spark.Subscription, plan *spark.Plan) string {
	if sub == nil || plan == nil {
		return ""
	}

	switch sub.Status {
	case spark.StatusAwaitingPayment:
		if !sub.HasPrice(plan.ID) {
			return ""
		}
	case spark.StatusActive, spark.StatusPastDue, spark.StatusPaused:
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
