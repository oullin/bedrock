package billing

// CheckoutItem represents a single price/quantity pair for checkout.
type CheckoutItem struct {
	PriceID  string `json:"priceId"`
	Quantity int    `json:"quantity"`
}

// Checkout represents a checkout session configuration.
type Checkout struct {
	Customer    *Customer
	Items       []CheckoutItem
	Transaction map[string]any
	CustomData  map[string]any
	ReturnURL   string
}

// GuestCheckout creates a checkout without an existing customer.
func GuestCheckout(items []CheckoutItem) *Checkout {
	return &Checkout{Items: items}
}

// CustomerCheckout creates a checkout for an existing customer.
func CustomerCheckout(customer *Customer, items []CheckoutItem) *Checkout {
	return &Checkout{Customer: customer, Items: items}
}

// TransactionCheckout creates a checkout from provider transaction data.
func TransactionCheckout(tx map[string]any, customer *Customer) *Checkout {
	return &Checkout{Customer: customer, Transaction: tx}
}

// WithCustomData attaches custom metadata to the checkout session.
func (c *Checkout) WithCustomData(data map[string]any) *Checkout {
	c.CustomData = data

	return c
}

// WithReturnURL sets the URL the customer is redirected to after checkout.
func (c *Checkout) WithReturnURL(url string) *Checkout {
	c.ReturnURL = url

	return c
}

// Options returns the checkout configuration as a map suitable for
// passing to a payment provider SDK.
func (c *Checkout) Options() map[string]any {
	opts := make(map[string]any)

	if c.Customer != nil {
		opts["customer"] = c.Customer.PaddleID
	}

	if len(c.CustomData) > 0 {
		opts["custom_data"] = c.CustomData
	}

	if c.ReturnURL != "" {
		opts["return_url"] = c.ReturnURL
	}

	return opts
}

// ToMap serialises the checkout for JSON responses.
func (c *Checkout) ToMap() map[string]any {
	m := map[string]any{
		"items": c.Items,
	}

	if c.Customer != nil {
		m["customer"] = c.Customer.PaddleID
	}

	if len(c.CustomData) > 0 {
		m["custom_data"] = c.CustomData
	}

	if c.ReturnURL != "" {
		m["return_url"] = c.ReturnURL
	}

	if len(c.Transaction) > 0 {
		m["transaction"] = c.Transaction
	}

	return m
}
