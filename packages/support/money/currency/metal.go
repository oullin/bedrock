package currency

// Provider defines the interface for providing currency information.
type Provider interface {
	// Get returns the currency configuration.
	Get() Currency
	// GetCode returns the currency code.
	GetCode() string
	// GetSymbols returns the list of available currency symbols.
	GetSymbols() []Symbol
}

// DefaultProvider is a basic implementation of the Provider interface, defaulting to SGD.
type DefaultProvider struct{}

// GetCode returns the code for the default currency (SGD).
func (m DefaultProvider) GetCode() string {
	return SGD
}

// Get returns the default currency configuration (SGD).
func (m DefaultProvider) Get() Currency {
	return Currency{
		Decimal:     ".",
		Thousand:    ",",
		Code:        m.GetCode(),
		Fraction:    2,
		NumericCode: "702",
		Grapheme:    "S$",
		Template:    "$1",
	}
}

// GetSymbols returns the list of all available currency symbols.
func (m DefaultProvider) GetSymbols() []Symbol {
	return GetSymbols()
}
