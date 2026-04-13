// Package state builds the billing portal frontend state.
package state

import (
	"context"

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

	var activePlan *spark.Plan

	if sub != nil && (sub.Active() || sub.PastDue()) {
		for _, p := range plans {
			if sub.HasPrice(p.ID) {
				activePlan = p

				break
			}
		}
	}

	state := resolveState(sub)

	data := map[string]any{
		"billableId":   billable.BillableID(),
		"billableName": billable.BillableName(),
		"billableType": billableType,
		"brandColor":   f.brandColor(),
		"dashboardUrl": f.dashboardURL(),
		"monthlyPlans": monthlyPlans,
		"yearlyPlans":  yearlyPlans,
		"plan":         activePlan,
		"seatName":     f.manager.SeatName(billableType),
		"sparkPath":    f.config.Path,
		"state":        state,
		"termsUrl":     f.config.TermsURL,
	}

	return data, nil
}

func resolveState(sub *spark.Subscription) string {
	if sub == nil {
		return "none"
	}

	if sub.OnGracePeriod() {
		return "onGracePeriod"
	}

	if sub.Active() {
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
