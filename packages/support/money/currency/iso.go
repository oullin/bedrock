package currency

import (
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/gollin/packages/support/money/exception"
)

// ISOCodePattern provides functionality to identify and extract currency ISO codes and symbols.
type ISOCodePattern struct {
	currencies              *Map
	symbols                 *[]Symbol
	defaultCurrency         *Currency
	symbolsLongestFirst     *[]Symbol
	symbolsLongestFirstOnce sync.Once
	isoCodePatternOnce      sync.Once
	isoCodePattern          *regexp.Regexp
}

// NewISOCodePattern creates a new ISOCodePattern instance using default currency data.
func NewISOCodePattern() *ISOCodePattern {
	provider := DefaultProvider{}

	symbols := provider.GetSymbols()
	defaultCurrency := provider.Get()
	currencies := NewCurrenciesMap()

	return &ISOCodePattern{
		currencies:      &currencies,
		symbols:         &symbols,
		defaultCurrency: &defaultCurrency,
	}
}

// NewISOCodePatternWith creates a new ISOCodePattern instance with a custom provider and dataset.
func NewISOCodePatternWith(provider Provider, dataset *map[string]*Currency) (*ISOCodePattern, error) {
	if provider == nil {
		return nil, exception.ErrNoConverterProvided
	}

	if dataset == nil || len(*dataset) == 0 {
		return nil, exception.ErrNoCurrencyMapDataset
	}

	currencies, err := NewCurrenciesMapFrom(dataset)

	if err != nil {
		return nil, exception.ErrNoCurrencyMapDataset
	}

	symbols := provider.GetSymbols()
	defaultCurrency := provider.Get()

	return &ISOCodePattern{
		symbols:         &symbols,
		currencies:      &currencies,
		defaultCurrency: &defaultCurrency,
	}, nil
}

// GetSymbolsLongestFirst returns a list of symbols sorted by length (descending) to ensure greedy matching.
// This is important for parsing logic to avoid matching a substring of a longer symbol (e.g. "$" vs "US$").
func (iso *ISOCodePattern) GetSymbolsLongestFirst() *[]Symbol {
	iso.symbolsLongestFirstOnce.Do(func() {
		symbols := *iso.symbols

		sort.SliceStable(symbols, func(i, j int) bool {
			li, lj := len(symbols[i].Id), len(symbols[j].Id)

			if li == lj {
				return symbols[i].Id < symbols[j].Id
			}

			return li > lj
		})

		iso.symbolsLongestFirst = &symbols
	})

	return iso.symbolsLongestFirst
}

// GetPattern returns a compiled regular expression for matching all known currency ISO codes.
// The pattern matches whole words only (using word boundaries).
func (iso *ISOCodePattern) GetPattern() *regexp.Regexp {
	iso.isoCodePatternOnce.Do(func() {
		codes := iso.currencies.GetCodes()

		// Build pattern: \b(CODE1|CODE2|...)\b
		// Escape is not needed since currency codes are alphanumeric
		pattern := `\b(` + strings.Join(*codes, "|") + `)\b`
		iso.isoCodePattern = regexp.MustCompile(pattern)
	})

	return iso.isoCodePattern
}
