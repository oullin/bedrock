package spark

// PricePreview represents a previewed price calculation including
// subtotals and tax. Mirrors Laravel\Paddle\PricePreview.
type PricePreview struct {
	PriceInfo Price
	Total     int64
	Subtotal  int64
	Tax       int64
	Currency  string
}

// RawTotal returns the total in minor units.
func (pp *PricePreview) RawTotal() int64 {
	return pp.Total
}

// RawSubtotal returns the subtotal in minor units.
func (pp *PricePreview) RawSubtotal() int64 {
	return pp.Subtotal
}

// RawTax returns the tax in minor units.
func (pp *PricePreview) RawTax() int64 {
	return pp.Tax
}

// HasTax reports whether the preview includes a non-zero tax amount.
func (pp *PricePreview) HasTax() bool {
	return pp.Tax > 0
}
