package entitlement

import (
	"context"
	"net/http"
	"strings"

	"github.com/bedrock/packages/billing"
)

// Decision is the result of a subscription gate check.
type Decision struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Granted    bool   `json:"granted"`
}

// IsGranted reports whether access was granted.
func (d Decision) IsGranted() bool { return d.Granted }

// IsDenied reports whether access was denied.
func (d Decision) IsDenied() bool { return !d.Granted }

// DecisionAllowed returns a decision indicating access is granted.
func DecisionAllowed() Decision {
	return Decision{StatusCode: http.StatusOK, Granted: true}
}

// DecisionNoSubscription returns a decision indicating no subscription exists.
func DecisionNoSubscription() Decision {
	return Decision{
		StatusCode: http.StatusPaymentRequired,
		Message:    "An active subscription is required to access this feature.",
	}
}

// DecisionFeatureMissing returns a decision indicating a required feature is missing.
func DecisionFeatureMissing() Decision {
	return Decision{
		StatusCode: http.StatusForbidden,
		Message:    "Your current subscription does not include this feature.",
	}
}

// DecisionPlanRestricted returns a decision indicating the plan does not allow access.
func DecisionPlanRestricted() Decision {
	return Decision{
		StatusCode: http.StatusForbidden,
		Message:    "Your current subscription does not include this plan-restricted feature.",
	}
}

// Gate checks subscription requirements for access control.
type Gate struct {
	entitlements *Access
}

// NewGate creates a Gate.
func NewGate(entitlements *Access) *Gate {
	return &Gate{entitlements: entitlements}
}

// CheckRequirements verifies the billable has a valid subscription meeting the
// given plan and/or feature requirements. Arguments are auto-detected: values
// matching a SubscriptionPlan are treated as plan checks; others are treated
// as feature code checks.
func (g *Gate) CheckRequirements(ctx context.Context, billableType string, billableID int64, requirements ...string) Decision {
	sub, err := g.entitlements.CurrentSubscription(ctx, billableType, billableID)
	if err != nil || sub == nil {
		return DecisionNoSubscription()
	}

	var planArgs []string
	var featureArgs []string

	for _, req := range requirements {
		if billing.SubscriptionPlan(strings.ToLower(req)).Valid() {
			planArgs = append(planArgs, strings.ToLower(req))
		} else {
			featureArgs = append(featureArgs, req)
		}
	}

	if len(planArgs) > 0 {
		found := false

		for _, p := range planArgs {
			if strings.EqualFold(sub.Plan, p) {
				found = true

				break
			}
		}

		if !found {
			return DecisionPlanRestricted()
		}
	}

	for _, featureCode := range featureArgs {
		if !g.entitlements.SubscriptionHasFeature(ctx, sub, featureCode) {
			return DecisionFeatureMissing()
		}
	}

	return DecisionAllowed()
}
