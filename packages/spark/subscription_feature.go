package spark

import "time"

// SubscriptionFeature links a subscription to a feature, tracking its
// activation state. Soft-deletable.
type SubscriptionFeature struct {
	ID             int64
	SubscriptionID int64
	FeatureID      int64
	Value          *string
	IsActive       bool
	ActivatedAt    *time.Time
	DeactivatedAt  *time.Time
	DeletedAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time

	// Loaded relations.
	Feature *Feature
}
