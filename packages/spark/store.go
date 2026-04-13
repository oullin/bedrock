package spark

import "context"

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
	ActiveForBillable(ctx context.Context, billableType string, billableID int64) ([]*Subscription, error)
	Create(ctx context.Context, s *Subscription) error
	Save(ctx context.Context, s *Subscription) error
	Delete(ctx context.Context, id int64) error
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

// ProductStore abstracts persistence and querying of products.
type ProductStore interface {
	FindByID(ctx context.Context, id int64) (*Product, error)
	Active(ctx context.Context) ([]Product, error)
	ActiveSubscriptions(ctx context.Context) ([]Product, error)
	ActiveOneTime(ctx context.Context) ([]Product, error)
}

// OrderStore abstracts persistence of order records.
type OrderStore interface {
	FindByID(ctx context.Context, id int64) (*Order, error)
	FindByBillable(ctx context.Context, teamID int64, limit int) ([]Order, error)
	Create(ctx context.Context, o *Order) error
	Save(ctx context.Context, o *Order) error
	HasCompletedForProduct(ctx context.Context, teamID int64, productID int64) (bool, error)
}
