package billing

// PricePreview represents a price preview with tax calculations, typically
// fetched from the payment provider based on the customer's location.
type PricePreview struct {
	Price    Price
	Total    int64  // Total amount in minor units (including tax).
	Subtotal int64  // Subtotal before tax.
	Tax      int64  // Tax amount in minor units.
	Currency string // ISO 4217 code.
}

// RawTotal returns the total in minor units.
func (p PricePreview) RawTotal() int64 { return p.Total }

// RawSubtotal returns the subtotal in minor units.
func (p PricePreview) RawSubtotal() int64 { return p.Subtotal }

// RawTax returns the tax in minor units.
func (p PricePreview) RawTax() int64 { return p.Tax }

// HasTax reports whether any tax is applied.
func (p PricePreview) HasTax() bool { return p.Tax > 0 }
