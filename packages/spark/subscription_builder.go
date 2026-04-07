package spark

// SubscriptionBuilder provides a fluent API for constructing new subscription
// checkout sessions, mirroring Cashier's SubscriptionBuilder.
type SubscriptionBuilder struct {
	billable Billable
	priceID  string
	name     string
	subType  string
	quantity int
	interval SubscriptionInterval
}

// NewSubscriptionBuilder creates a builder for a new subscription.
func NewSubscriptionBuilder(billable Billable, priceID string, name string) *SubscriptionBuilder {
	return &SubscriptionBuilder{
		billable: billable,
		priceID:  priceID,
		name:     name,
		subType:  DefaultSubscriptionType,
		quantity: 1,
	}
}

// Type sets the subscription type.
func (b *SubscriptionBuilder) Type(t string) *SubscriptionBuilder {
	b.subType = t

	return b
}

// Quantity sets the number of seats/units.
func (b *SubscriptionBuilder) Quantity(q int) *SubscriptionBuilder {
	b.quantity = q

	return b
}

// Daily sets the billing interval to daily.
func (b *SubscriptionBuilder) Daily() *SubscriptionBuilder {
	b.interval = IntervalDay

	return b
}

// Weekly sets the billing interval to weekly.
func (b *SubscriptionBuilder) Weekly() *SubscriptionBuilder {
	b.interval = IntervalWeek

	return b
}

// Monthly sets the billing interval to monthly.
func (b *SubscriptionBuilder) Monthly() *SubscriptionBuilder {
	b.interval = IntervalMonth

	return b
}

// Yearly sets the billing interval to yearly.
func (b *SubscriptionBuilder) Yearly() *SubscriptionBuilder {
	b.interval = IntervalYear

	return b
}

// Build returns the configured checkout items and metadata.
func (b *SubscriptionBuilder) Build() ([]CheckoutItem, string, int) {
	return []CheckoutItem{
		{PriceID: b.priceID, Quantity: b.quantity},
	}, b.subType, b.quantity
}
