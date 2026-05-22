package billing

// SubscriptionBuilder provides a fluent API for constructing subscription
// checkout parameters.
type SubscriptionBuilder struct {
	billable         Billable
	priceID          string
	name             string
	subscriptionType string
	quantity         int
	interval         SubscriptionInterval
}

// NewSubscriptionBuilder creates a builder for the given billable.
func NewSubscriptionBuilder(billable Billable, priceID string, name string) *SubscriptionBuilder {
	return &SubscriptionBuilder{
		billable:         billable,
		priceID:          priceID,
		name:             name,
		subscriptionType: DefaultSubscriptionType,
		quantity:         1,
		interval:         IntervalMonth,
	}
}

// Type sets the subscription type (e.g. "default").
func (b *SubscriptionBuilder) Type(t string) *SubscriptionBuilder {
	b.subscriptionType = t

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

// Checkout returns a Checkout configured with the builder's parameters.
func (b *SubscriptionBuilder) Checkout() *Checkout {
	items := []CheckoutItem{
		{PriceID: b.priceID, Quantity: b.quantity},
	}

	return GuestCheckout(items)
}

// Build returns the configured checkout items and metadata.
func (b *SubscriptionBuilder) Build() ([]CheckoutItem, string, int) {
	return []CheckoutItem{{PriceID: b.priceID, Quantity: b.quantity}}, b.subscriptionType, b.quantity
}
