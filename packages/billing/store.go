package billing

import (
	"context"
	"time"
)

// CustomerStore abstracts persistence of customer records.
type CustomerStore interface {
	FindByProviderID(ctx context.Context, providerID string) (*Customer, error)
	FindByBillable(ctx context.Context, billableType string, billableID int64) (*Customer, error)
	Create(ctx context.Context, c *Customer) error
	Save(ctx context.Context, c *Customer) error
}

// SubscriptionStore abstracts persistence of subscription records.
type SubscriptionStore interface {
	FindByID(ctx context.Context, id int64) (*Subscription, error)
	FindByProviderID(ctx context.Context, providerID string) (*Subscription, error)
	CurrentForBillable(ctx context.Context, billableType string, billableID int64) (*Subscription, error)
	LatestForPlan(ctx context.Context, billableType string, billableID int64, plan string) (*Subscription, error)
	Create(ctx context.Context, s *Subscription) error
	Save(ctx context.Context, s *Subscription) error
	Delete(ctx context.Context, id int64) error
	FindExpirable(ctx context.Context, before time.Time) ([]*Subscription, error)
	FindPendingWithoutProvider(ctx context.Context) ([]*Subscription, error)
	ExistsAccessibleNonStarter(ctx context.Context, billableType string, billableID int64) (bool, error)
	ActiveForBillable(ctx context.Context, billableType string, billableID int64) ([]*Subscription, error)
}

// SubscriptionItemStore abstracts persistence of subscription items.
type SubscriptionItemStore interface {
	FindBySubscription(ctx context.Context, subscriptionID int64) ([]SubscriptionItem, error)
	Sync(ctx context.Context, subscriptionID int64, items []SubscriptionItem) error
}

// TransactionStore abstracts persistence of transaction records.
type TransactionStore interface {
	FindByProviderID(ctx context.Context, providerID string) (*Transaction, error)
	FindByBillable(ctx context.Context, billableType string, billableID int64, limit int) ([]Transaction, error)
	Create(ctx context.Context, t *Transaction) error
	Save(ctx context.Context, t *Transaction) error
}

// PlanStore abstracts persistence of the plan catalog.
type PlanStore interface {
	ActivePlans(ctx context.Context) ([]Plan, error)
	FindBySlug(ctx context.Context, slug string) (*Plan, error)
	FindByCode(ctx context.Context, code string) (*Plan, error)
	UpsertPlans(ctx context.Context, plans []Plan) error
	UpsertFeatures(ctx context.Context, features []Feature) error
	SyncPlanFeatures(ctx context.Context, planID int64, features []PlanFeature) error
	FindPlanPeriod(ctx context.Context, planID int64, period BillingPeriod) (*PlanPeriod, error)
	UpsertPlanPeriods(ctx context.Context, periods []PlanPeriod) error
}

// PriceStore abstracts persistence of plan period prices.
type PriceStore interface {
	FindByProviderPriceID(ctx context.Context, providerPriceID string) (*PlanPeriodPrice, error)
	ActiveForPeriod(ctx context.Context, planPeriodID int64) (*PlanPeriodPrice, error)
	Create(ctx context.Context, p *PlanPeriodPrice) error
	Save(ctx context.Context, p *PlanPeriodPrice) error
	DeactivateForPeriod(ctx context.Context, planPeriodID int64) error
}

// SubscriptionFeatureStore abstracts persistence of subscription features.
type SubscriptionFeatureStore interface {
	FeaturesForSubscription(ctx context.Context, subscriptionID int64) ([]SubscriptionFeature, error)
	Sync(ctx context.Context, subscriptionID int64, features []SubscriptionFeature) error
	SetActivationState(ctx context.Context, subscriptionID int64, active bool, at time.Time) error
	DeleteAll(ctx context.Context, subscriptionID int64) error
}
