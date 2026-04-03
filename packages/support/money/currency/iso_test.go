package currency

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gollin/packages/support/money/exception"
)

type stubProvider struct {
	code    string
	symbols []Symbol
}

func (s stubProvider) Get() Currency {
	return Currency{
		Code:        s.code,
		Fraction:    2,
		Grapheme:    "X$",
		Template:    "$1",
		Decimal:     ".",
		Thousand:    ",",
		NumericCode: "000",
	}
}

func (s stubProvider) GetCode() string { return s.code }
func (s stubProvider) GetSymbols() []Symbol {
	if s.symbols != nil {
		return s.symbols
	}
	return []Symbol{{Id: "X$", Currency: s.code}}
}

func TestISOCodePattern_NewAndMethods(t *testing.T) {
	iso := NewISOCodePattern()
	if iso == nil {
		t.Fatal("NewISOCodePattern() returned nil")
	}

	symbols := iso.GetSymbolsLongestFirst()
	if symbols == nil || len(*symbols) == 0 {
		t.Fatal("GetSymbolsLongestFirst() returned nil or empty")
	}

	pattern := iso.GetPattern()
	if pattern == nil {
		t.Fatal("GetPattern() returned nil")
	}
	if !pattern.MatchString("SGD") {
		t.Fatal("pattern should match SGD")
	}
}

func TestISOCodePattern_NewWith_Errors(t *testing.T) {
	original := NewCurrenciesMapFrom
	NewCurrenciesMapFrom = createCurrenciesMapFromFactory()
	defer func() { NewCurrenciesMapFrom = original }()

	t.Run("nil provider", func(t *testing.T) {
		_, err := NewISOCodePatternWith(nil, &map[string]*Currency{"SGD": {Code: "SGD"}})
		if !errors.Is(err, exception.ErrNoConverterProvided) {
			t.Fatalf("error = %v, want ErrNoConverterProvided", err)
		}
	})

	t.Run("nil dataset", func(t *testing.T) {
		_, err := NewISOCodePatternWith(stubProvider{code: SGD}, nil)
		if !errors.Is(err, exception.ErrNoCurrencyMapDataset) {
			t.Fatalf("error = %v, want ErrNoCurrencyMapDataset", err)
		}
	})

	t.Run("empty dataset", func(t *testing.T) {
		empty := map[string]*Currency{}
		_, err := NewISOCodePatternWith(stubProvider{code: SGD}, &empty)
		if !errors.Is(err, exception.ErrNoCurrencyMapDataset) {
			t.Fatalf("error = %v, want ErrNoCurrencyMapDataset", err)
		}
	})

	t.Run("dataset factory returns error", func(t *testing.T) {
		orig := NewCurrenciesMapFrom
		NewCurrenciesMapFrom = func(_ *map[string]*Currency) (Map, error) {
			return Map{}, errors.New("boom")
		}
		defer func() { NewCurrenciesMapFrom = orig }()

		data := map[string]*Currency{"SGD": {Code: "SGD"}}
		_, err := NewISOCodePatternWith(stubProvider{code: SGD}, &data)
		if !errors.Is(err, exception.ErrNoCurrencyMapDataset) {
			t.Fatalf("error = %v, want ErrNoCurrencyMapDataset", err)
		}
	})
}

func TestISOCodePatternWith_SortingAndPattern(t *testing.T) {
	original := NewCurrenciesMapFrom
	NewCurrenciesMapFrom = createCurrenciesMapFromFactory()
	defer func() { NewCurrenciesMapFrom = original }()

	data := map[string]*Currency{
		"SGD": {Code: "SGD"},
		"EUR": {Code: "EUR"},
	}

	p := stubProvider{
		code: SGD,
		symbols: []Symbol{
			{Id: "S$", Currency: SGD},
			{Id: "C$", Currency: CAD},
			{Id: "A$", Currency: AUD},
			{Id: "CFP", Currency: XPF},
			{Id: "R$", Currency: BRL},
		},
	}

	iso, err := NewISOCodePatternWith(p, &data)
	if err != nil {
		t.Fatalf("NewISOCodePatternWith() unexpected error: %v", err)
	}

	longest := iso.GetSymbolsLongestFirst()
	if longest == nil {
		t.Fatal("GetSymbolsLongestFirst() returned nil")
	}
	if (*longest)[0].Id != "CFP" {
		t.Fatalf("GetSymbolsLongestFirst()[0] = %q, want %q", (*longest)[0].Id, "CFP")
	}

	// Tie-break for length=2 should be lexicographic by Id (A$, C$, R$)
	gotIDs := []string{(*longest)[1].Id, (*longest)[2].Id, (*longest)[3].Id}
	wantIDs := []string{"A$", "C$", "R$"}
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Fatalf("GetSymbolsLongestFirst() tie-break = %#v, want %#v", gotIDs, wantIDs)
	}

	pattern := iso.GetPattern()
	if pattern == nil {
		t.Fatal("GetPattern() returned nil")
	}
	if !pattern.MatchString("10 SGD") || !pattern.MatchString("EUR") {
		t.Fatal("pattern should match SGD and EUR")
	}
}
