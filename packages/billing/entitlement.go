package billing

import (
	"context"
	"math"
)

// EntitlementAccess checks subscription-based access and feature entitlements.
type EntitlementAccess struct {
	subscriptions SubscriptionStore
	features      SubscriptionFeatureStore
	clock         Clock
}

// NewEntitlementAccess creates an EntitlementAccess service.
func NewEntitlementAccess(
	subscriptions SubscriptionStore,
	features SubscriptionFeatureStore,
	clock Clock,
) *EntitlementAccess {
	return &EntitlementAccess{
		subscriptions: subscriptions,
		features:      features,
		clock:         clock,
	}
}

// HasAccessibleSubscription reports whether the billable has an active subscription.
func (e *EntitlementAccess) HasAccessibleSubscription(ctx context.Context, billableType string, billableID int64) bool {
	sub, err := e.currentSubscription(ctx, billableType, billableID)

	return err == nil && sub != nil
}

// ActiveRoleCapacity returns the active role capacity for the billable.
// Returns 0 if no subscription, math.MaxInt if unlimited.
func (e *EntitlementAccess) ActiveRoleCapacity(ctx context.Context, billableType string, billableID int64) int {
	sub, err := e.currentSubscription(ctx, billableType, billableID)
	if err != nil || sub == nil {
		return 0
	}

	return e.capacityForSubscription(ctx, sub)
}

// IsUnlimitedCapacity reports whether the billable has unlimited role capacity.
func (e *EntitlementAccess) IsUnlimitedCapacity(ctx context.Context, billableType string, billableID int64) bool {
	return e.ActiveRoleCapacity(ctx, billableType, billableID) == math.MaxInt
}

func (e *EntitlementAccess) capacityForSubscription(ctx context.Context, sub *Subscription) int {
	features, err := e.features.FeaturesForSubscription(ctx, sub.ID)
	if err != nil {
		return 0
	}

	for _, f := range features {
		if !f.IsActive || f.Feature == nil {
			continue
		}

		if SubscriptionFeatureCode(f.Feature.Code) != FeatureActiveRoleCapacity {
			continue
		}

		if f.Value == nil {
			return math.MaxInt
		}

		v := 0
		for _, c := range *f.Value {
			if c < '0' || c > '9' {
				return 0
			}

			v = v*10 + int(c-'0')
		}

		return v
	}

	return 0
}

// HasFeature reports whether the billable's subscription includes the given feature.
func (e *EntitlementAccess) HasFeature(ctx context.Context, billableType string, billableID int64, featureCode string) bool {
	sub, err := e.currentSubscription(ctx, billableType, billableID)
	if err != nil || sub == nil {
		return false
	}

	return e.SubscriptionHasFeature(ctx, sub, featureCode)
}

// SubscriptionHasFeature reports whether a subscription includes the given feature.
func (e *EntitlementAccess) SubscriptionHasFeature(ctx context.Context, sub *Subscription, featureCode string) bool {
	features, err := e.features.FeaturesForSubscription(ctx, sub.ID)
	if err != nil {
		return false
	}

	for _, f := range features {
		if f.IsActive && f.Feature != nil && f.Feature.Code == featureCode {
			return true
		}
	}

	return false
}

// CurrentSubscription returns the current accessible subscription for a billable.
func (e *EntitlementAccess) CurrentSubscription(ctx context.Context, billableType string, billableID int64) (*Subscription, error) {
	return e.currentSubscription(ctx, billableType, billableID)
}

func (e *EntitlementAccess) currentSubscription(ctx context.Context, billableType string, billableID int64) (*Subscription, error) {
	return e.subscriptions.CurrentForBillable(ctx, billableType, billableID)
}
