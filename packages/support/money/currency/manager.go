package currency

import (
	"fmt"
	"strings"
)

// Manager manages a collection of currencies and provides methods to access them.
type Manager struct {
	currencies      *Map
	symbols         *[]Symbol
	defaultCurrency *Currency
}

// NewManager creates a new Manager instance using the DefaultProvider.
func NewManager() *Manager {
	provider := DefaultProvider{}

	symbols := provider.GetSymbols()
	defaultCurrency := provider.Get()
	currencies := NewCurrenciesMap()

	return &Manager{
		symbols:         &symbols,
		currencies:      &currencies,
		defaultCurrency: &defaultCurrency,
	}
}

// NewManagerWith creates a new Manager instance using a custom Provider.
func NewManagerWith(metal Provider) *Manager {
	symbols := metal.GetSymbols()
	defaultCurrency := metal.Get()
	currencies := NewCurrenciesMap()

	return &Manager{
		symbols:         &symbols,
		currencies:      &currencies,
		defaultCurrency: &defaultCurrency,
	}
}

// NewManagerFor creates a new Manager instance with a specific provider and an initial dataset of currencies.
func NewManagerFor(curr Provider, dataset *map[string]*Currency) (*Manager, error) {
	var targetCurrency Currency

	if curr != nil {
		targetCurrency = curr.Get()
	} else {
		targetCurrency = new(DefaultProvider).Get()
	}

	currencies, err := NewCurrenciesMapFrom(dataset)

	if err != nil {
		return nil, fmt.Errorf("failed to create dataset map: %w", err)
	}

	if _, err = currencies.HasInvalidState(); err != nil {
		return nil, err
	}

	return &Manager{
		currencies:      &currencies,
		defaultCurrency: &targetCurrency,
	}, nil
}

// FindByCode searches for a currency by its ISO 4217 code.
func (cm *Manager) FindByCode(code string) *Currency {
	if cm == nil || cm.currencies == nil {
		return nil
	}

	return cm.currencies.FindByCode(code)
}

// FindByNumericCode returns the currency given the numeric code defined in ISO-4271.
func (cm *Manager) FindByNumericCode(code string) *Currency {
	if cm == nil || cm.currencies == nil {
		return nil
	}

	for _, currency := range *cm.currencies.dataset {
		if currency != nil && currency.NumericCode == strings.ToUpper(code) {
			return currency
		}
	}

	return nil
}

// Add adds a new currency to the manager.
func (cm *Manager) Add(currency *Currency) *Currency {
	if cm == nil || currency == nil || cm.currencies == nil || cm.currencies.dataset == nil {
		return nil
	}

	code := strings.TrimSpace(
		strings.ToUpper(currency.Code),
	)

	if code == "" {
		return nil
	}

	(*cm.currencies.dataset)[code] = currency

	return currency
}

// AddFrom creates and adds a new currency to the manager using the provided details.
func (cm *Manager) AddFrom(code, Grapheme, Template, Decimal, Thousand, NumericCode string, Fraction int) *Currency {
	if cm == nil {
		return nil
	}

	currency := &Currency{
		Code:        strings.ToUpper(code),
		Grapheme:    Grapheme,
		Template:    Template,
		Decimal:     Decimal,
		Thousand:    Thousand,
		NumericCode: strings.ToUpper(NumericCode),
		Fraction:    Fraction,
	}

	return cm.Add(currency)
}

// GetDefault returns the default currency of the manager.
func (cm *Manager) GetDefault() *Currency {
	if cm == nil {
		return nil
	}

	return cm.defaultCurrency
}

// GetSymbols returns the list of currency symbols managed by the manager.
func (cm *Manager) GetSymbols() *[]Symbol {
	if cm == nil {
		return nil
	}

	return cm.symbols
}

// Resolve tries to find a currency by code and falls back to the default currency.
func (cm *Manager) Resolve(code string) *Currency {
	if cm == nil {
		return nil
	}

	if currency := cm.FindByCode(code); currency != nil {
		return currency
	}

	return cm.defaultCurrency
}
