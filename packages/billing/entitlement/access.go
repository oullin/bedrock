package entitlement

import (
	"context"
	"math"

	"github.com/bedrock/packages/billing"
)

// Access checks subscription-based access and feature entitlements.
type Access struct {
	subscriptions billing.SubscriptionStore
	features      billing.SubscriptionFeatureStore
	clock         billing.Clock
}

// NewAccess creates an Access service.
func NewAccess(
	subscriptions billing.SubscriptionStore,
	features billing.SubscriptionFeatureStore,
	clock billing.Clock,
) *Access {
	return &Access{
		subscriptions: subscriptions,
		features:      features,
		clock:         clock,
	}
}

// HasAccessibleSubscription reports whether the billable has an active subscription.
func (e *Access) HasAccessibleSubscription(ctx context.Context, billableType string, billableID int64) bool {
	sub, err := e.currentSubscription(ctx, billableType, billableID)

	return err == nil && sub != nil
}

// ActiveRoleCapacity returns the active role capacity for the billable.
// Returns 0 if no subscription, math.MaxInt if unlimited.
func (e *Access) ActiveRoleCapacity(ctx context.Context, billableType string, billableID int64) int {
	sub, err := e.currentSubscription(ctx, billableType, billableID)

	if err != nil || sub == nil {
		return 0
	}

	return e.capacityForSubscription(ctx, sub)
}

// IsUnlimitedCapacity reports whether the billable has unlimited role capacity.
func (e *Access) IsUnlimitedCapacity(ctx context.Context, billableType string, billableID int64) bool {
	return e.ActiveRoleCapacity(ctx, billableType, billableID) == math.MaxInt
}

func (e *Access) capacityForSubscription(ctx context.Context, sub *billing.Subscription) int {
	features, err := e.features.FeaturesForSubscription(ctx, sub.ID)

	if err != nil {
		return 0
	}

	for _, f := range features {
		if !f.IsActive || f.Feature == nil {
			continue
		}

		if billing.SubscriptionFeatureCode(f.Feature.Code) != billing.FeatureActiveRoleCapacity {
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
func (e *Access) HasFeature(ctx context.Context, billableType string, billableID int64, featureCode string) bool {
	sub, err := e.currentSubscription(ctx, billableType, billableID)

	if err != nil || sub == nil {
		return false
	}

	return e.SubscriptionHasFeature(ctx, sub, featureCode)
}

// SubscriptionHasFeature reports whether a subscription includes the given feature.
func (e *Access) SubscriptionHasFeature(ctx context.Context, sub *billing.Subscription, featureCode string) bool {
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
func (e *Access) CurrentSubscription(ctx context.Context, billableType string, billableID int64) (*billing.Subscription, error) {
	return e.currentSubscription(ctx, billableType, billableID)
}

func (e *Access) currentSubscription(ctx context.Context, billableType string, billableID int64) (*billing.Subscription, error) {
	return e.subscriptions.CurrentForBillable(ctx, billableType, billableID)
}
