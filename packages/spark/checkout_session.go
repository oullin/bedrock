package spark

// CheckoutItem represents a single item in a checkout session.
type CheckoutItem struct {
	PriceID  string
	Quantity int
}

// CheckoutSession holds the data needed to initiate a provider checkout flow.
type CheckoutSession struct {
	Customer    *Customer
	Items       []CheckoutItem
	CustomData  map[string]any
	ReturnURL   string
	Transaction map[string]any // Provider-specific transaction data.
}

// GuestCheckout creates a checkout session without an existing customer.
func GuestCheckout(items []CheckoutItem) *CheckoutSession {
	return &CheckoutSession{Items: items}
}

// CustomerCheckout creates a checkout session for an existing customer.
func CustomerCheckout(customer *Customer, items []CheckoutItem) *CheckoutSession {
	return &CheckoutSession{Customer: customer, Items: items}
}

// TransactionCheckout creates a checkout session from an existing transaction.
func TransactionCheckout(tx map[string]any, customer *Customer) *CheckoutSession {
	return &CheckoutSession{Transaction: tx, Customer: customer}
}

// WithCustomData sets custom metadata on the checkout session.
func (c *CheckoutSession) WithCustomData(data map[string]any) *CheckoutSession {
	c.CustomData = data

	return c
}

// WithReturnURL sets the URL to redirect to after checkout completes.
func (c *CheckoutSession) WithReturnURL(url string) *CheckoutSession {
	c.ReturnURL = url

	return c
}
