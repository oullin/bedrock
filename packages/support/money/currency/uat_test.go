package currency

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/gollin/packages/support/money/exception"
)

type stubDefault struct {
	code string
}

func (s stubDefault) Get() Currency {
	return Currency{
		Code:        s.code,
		Fraction:    1,
		Grapheme:    "X",
		Template:    "$1",
		Decimal:     ".",
		Thousand:    ",",
		NumericCode: "000",
	}
}

func (s stubDefault) GetCode() string {
	return s.code
}

func (s stubDefault) GetSymbols() []Symbol {
	return []Symbol{{Id: "X$", Currency: s.code}}
}

func TestManagerWithAndResolve(t *testing.T) {
	def := stubDefault{code: "XYZ"}
	manager := NewManagerWith(def)

	if got := manager.GetDefault(); got.Code != def.code {
		t.Fatalf("GetDefault() = %s, want %s", got.Code, def.code)
	}

	if resolved := manager.Resolve("INVALID"); resolved.Code != def.code {
		t.Fatalf("Resolve() fallback = %s, want %s", resolved.Code, def.code)
	}

	if symbols := manager.GetSymbols(); symbols == nil || len(*symbols) == 0 {
		t.Fatal("GetSymbols() returned nil or empty slice")
	}
}

func TestManagerForValidation(t *testing.T) {
	t.Run("nil dataset", func(t *testing.T) {
		if _, err := NewManagerFor(nil, nil); err == nil {
			t.Fatal("expected error when dataset is nil")
		}
	})

	t.Run("success with custom default", func(t *testing.T) {
		def := stubDefault{code: "TST"}
		data := map[string]*Currency{"FOO": {Code: "FOO"}}

		manager, err := NewManagerFor(def, &data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if found := manager.FindByCode("FOO"); found == nil || found.Code != "FOO" {
			t.Fatalf("FindByCode() = %v, want FOO", found)
		}

		if manager.GetDefault().Code != "TST" {
			t.Fatalf("GetDefault() = %s, want TST", manager.GetDefault().Code)
		}
	})

	t.Run("invalid state", func(t *testing.T) {
		original := NewCurrenciesMapFrom
		NewCurrenciesMapFrom = createCurrenciesMapFromFactory()
		defer func() { NewCurrenciesMapFrom = original }()

		data := map[string]*Currency{"BAD": nil}
		if _, err := NewManagerFor(nil, &data); err == nil {
			t.Fatal("expected error for invalid currency map state")
		}
	})
}

func TestManagerAndMapNilSafety(t *testing.T) {
	var manager *Manager

	if res := manager.Resolve("SGD"); res != nil {
		t.Fatalf("nil manager Resolve() = %v, want nil", res)
	}

	if res := manager.AddFrom("NEW", "", "", "", "", "", 0); res != nil {
		t.Fatalf("nil manager AddFrom() = %v, want nil", res)
	}

	if res := manager.Add(&Currency{Code: "ANY"}); res != nil {
		t.Fatalf("nil manager Add() = %v, want nil", res)
	}

	if res := manager.FindByNumericCode("000"); res != nil {
		t.Fatalf("nil manager FindByNumericCode() = %v, want nil", res)
	}

	empty := &Manager{currencies: &Map{dataset: &map[string]*Currency{}}}
	if res := empty.FindByNumericCode("123"); res != nil {
		t.Fatalf("empty dataset FindByNumericCode() = %v, want nil", res)
	}
}

func TestManagerAddAndFind(t *testing.T) {
	manager := NewManager()

	if got := manager.Add(nil); got != nil {
		t.Fatalf("Add(nil) = %v, want nil", got)
	}

	if got := manager.AddFrom("bar", "$", "$1", ".", ",", "111", 2); got == nil {
		t.Fatal("AddFrom() returned nil")
	}

	if found := manager.Resolve("bar"); found == nil || found.Code != "BAR" {
		t.Fatalf("Resolve(bar) = %v, want BAR", found)
	}
}

func TestCurrencyMapHelpers(t *testing.T) {
	var empty Map
	if empty.Get("SGD") != nil {
		t.Fatal("Get on nil dataset should return nil")
	}

	if !empty.IsEmpty() {
		t.Fatal("IsEmpty() on nil dataset should be true")
	}

	if empty.IsNotEmpty() {
		t.Fatal("IsNotEmpty() on nil dataset should be false")
	}

	data := map[string]*Currency{"SGD": {Code: "SGD"}}
	withData := Map{dataset: &data}

	if withData.IsEmpty() {
		t.Fatal("IsEmpty() on populated dataset should be false")
	}

	if !withData.IsNotEmpty() {
		t.Fatal("IsNotEmpty() on populated dataset should be true")
	}

	if cur := withData.Get("SGD"); cur == nil || cur.Code != "SGD" {
		t.Fatalf("Get() returned %v, want SGD", cur)
	}

	if _, err := (Map{dataset: &map[string]*Currency{"BAD": nil}}).HasInvalidState(); err == nil {
		t.Fatal("HasInvalidState() expected error on nil currency value")
	}

	if invalid, err := withData.HasInvalidState(); invalid || err != nil {
		t.Fatalf("HasInvalidState() on valid map = (%v, %v), want (false, nil)", invalid, err)
	}

	if res := withData.FindByCode("missing"); res != nil {
		t.Fatalf("FindByCode on missing currency = %v, want nil", res)
	}

	if res := empty.FindByCode("SGD"); res != nil {
		t.Fatalf("FindByCode on nil dataset = %v, want nil", res)
	}
}

func TestCurrencyGetNil(t *testing.T) {
	var c *Currency
	got, err := c.Get()
	if err == nil {
		t.Fatal("Expected error for nil Currency Get(), got nil")
	}
	if got != nil {
		t.Fatalf("nil Currency Get() = %v, want nil", got)
	}
}

func TestCurrenciesMapFromSingletonReuse(t *testing.T) {
	original := NewCurrenciesMapFrom
	NewCurrenciesMapFrom = createCurrenciesMapFromFactory()
	defer func() { NewCurrenciesMapFrom = original }()

	data := map[string]*Currency{"ABC": {Code: "ABC"}}

	first, err := NewCurrenciesMapFrom(&data)
	if err != nil {
		t.Fatalf("unexpected error on first map creation: %v", err)
	}

	second, err := NewCurrenciesMapFrom(&map[string]*Currency{"IGNORED": {Code: "IGNORED"}})
	if err != nil {
		t.Fatalf("unexpected error on singleton reuse: %v", err)
	}

	if first.dataset != second.dataset {
		t.Fatal("expected singleton dataset to be reused on subsequent calls")
	}
}

// TestCurrenciesMapFromConcurrency verifies that the NewCurrenciesMapFrom function
// behaves correctly under concurrent access. This test ensures:
//   - Thread safety: Multiple goroutines can safely call NewCurrenciesMapFrom simultaneously
//   - Singleton consistency: All concurrent calls receive the same singleton instance
//   - Data integrity: The singleton contains valid data from one of the input datasets
//
// The test replaces NewCurrenciesMapFrom with a test factory that implements
// singleton behavior, then spawns 100 concurrent goroutines attempting to initialize
// with different datasets. All goroutines should successfully complete and receive
// the same singleton dataset instance.
func TestCurrenciesMapFromConcurrency(t *testing.T) {
	// Create a fresh instance of NewCurrenciesMapFrom for this test
	original := NewCurrenciesMapFrom
	NewCurrenciesMapFrom = createCurrenciesMapFromFactory()
	defer func() { NewCurrenciesMapFrom = original }()

	// Number of concurrent goroutines
	const numGoroutines = 100

	// Prepare test data - each goroutine gets different data to try to initialize with
	testDatasets := make([]*map[string]*Currency, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		code := fmt.Sprintf("C%d", i)
		testDatasets[i] = &map[string]*Currency{code: {Code: code}}
	}

	// Channel to collect results
	results := make(chan Map, numGoroutines)
	errors := make(chan error, numGoroutines)

	// Launch goroutines concurrently
	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			m, err := NewCurrenciesMapFrom(testDatasets[idx])
			if err != nil {
				errors <- err
				return
			}
			results <- m
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)
	close(errors)

	// Check for errors
	for err := range errors {
		t.Fatalf("unexpected error during concurrent initialization: %v", err)
	}

	// Verify all goroutines got the same singleton instance
	var firstDataset *map[string]*Currency
	count := 0
	for m := range results {
		count++
		if firstDataset == nil {
			firstDataset = m.dataset
		} else if m.dataset != firstDataset {
			t.Fatal("concurrent calls returned different singleton instances")
		}
	}

	if count != numGoroutines {
		t.Fatalf("expected %d results, got %d", numGoroutines, count)
	}

	// Verify the singleton contains data from exactly one of the test datasets
	if firstDataset == nil {
		t.Fatal("singleton dataset is nil")
	}

	foundMatch := false
	for _, testData := range testDatasets {
		if firstDataset == testData {
			foundMatch = true
			break
		}
	}

	if !foundMatch {
		t.Fatal("singleton dataset doesn't match any of the input datasets")
	}
}

func TestFindCurrencyByNumericCode(t *testing.T) {
	manager := NewManager()

	// Test standard currency (SGD - 702)
	sgd := manager.FindByNumericCode("702")
	if sgd == nil {
		t.Fatal("Expected SGD for 702, got nil")
	}
	if sgd.Code != SGD {
		t.Errorf("Expected SGD, got %s", sgd.Code)
	}

	// Test non-existent currency
	none := manager.FindByNumericCode("999999")
	if none != nil {
		t.Errorf("Expected nil, got %v", none)
	}

	// Test fallback to linear scan
	// We need to add a currency that is NOT in the init-generated map but IS in the main map
	// However, init() runs once. AddCurrency updates the main map.
	// GetCurrencyByNumericCode checks the cache (map) first, then the main map (linear scan).

	newCode := "X_TEST_NUM"
	manager.AddFrom(newCode, "X", "1$", ".", ",", "999888", 2)

	found := manager.FindByNumericCode("999888")
	if found == nil {
		t.Fatal("Expected to find newly added currency by numeric code")
	}
	if found.Code != newCode {
		t.Errorf("Expected %s, got %s", newCode, found.Code)
	}

	// Verify fallback still works when numeric code is set after insertion.
	fallback := manager.AddFrom("X_TEST_FALLBACK", "XF", "1$", ".", ",", "", 2)
	fallback.NumericCode = "111222"

	foundFallback := manager.FindByNumericCode("111222")
	if foundFallback == nil {
		t.Fatal("Expected to find newly added currency by numeric code via fallback scan")
	}
	if foundFallback.Code != fallback.Code {
		t.Errorf("Expected %s, got %s", fallback.Code, foundFallback.Code)
	}
}

func TestFindCurrencyByCode(t *testing.T) {
	manager := NewManager()

	sgd := manager.FindByCode(SGD)
	if sgd == nil {
		t.Fatal("Expected SGD, got nil")
	}
	if sgd.Code != SGD {
		t.Errorf("Expected SGD, got %s", sgd.Code)
	}

	sgdLower := manager.FindByCode("sgd")
	if sgdLower == nil {
		t.Fatal("Expected SGD for 'sgd', got nil")
	}

	// Test non-existent
	none := manager.FindByCode("INVALID_CODE")
	if none != nil {
		t.Errorf("Expected nil, got %v", none)
	}
}

func TestNewCurrency(t *testing.T) {
	manager := NewManager()
	c := manager.FindByCode("sgd")
	if c.Code != SGD {
		t.Errorf("Expected SGD, got %s", c.Code)
	}
}

func TestNewDefaultCurrency(t *testing.T) {
	// Should return SDG usually, or whatever NewCurrency(SDG) returns
	manager := NewManager()
	def := manager.FindByCode(SDG)

	if def == nil {
		t.Fatal("Expected default currency, got nil")
	}

	if def.Code != SDG {
		t.Errorf("Expected SDG, got %s", def.Code)
	}
}

func TestGetDefault(t *testing.T) {
	manager := NewManager()
	def := manager.GetDefault()

	if def == nil {
		t.Fatal("Expected default currency, got nil")
	}

	if def.Code != SGD {
		t.Errorf("Expected default currency %s, got %s", SGD, def.Code)
	}
}

func TestGet(t *testing.T) {
	// Valid currency
	c := &Currency{Code: SGD, NumericCode: "702", Fraction: 2}
	got, err := c.Get()
	if err != nil {
		t.Fatalf("Unexpected error getting currency: %v", err)
	}
	if got == nil {
		t.Fatal("Expected to get currency")
	}
	if got != c {
		t.Errorf("Expected Get() to return same currency instance")
	}
	if got.Code != SGD {
		t.Errorf("Expected code %s, got %s", SGD, got.Code)
	}

	// Currency with custom code
	c2 := &Currency{Code: "CUSTOM", Fraction: 3}
	got2, err := c2.Get()
	if err != nil {
		t.Fatalf("Unexpected error getting custom currency: %v", err)
	}
	if got2.Code != "CUSTOM" {
		t.Errorf("Expected CUSTOM, got %s", got2.Code)
	}
	if got2.Fraction != 3 {
		t.Errorf("Expected fraction 3, got %d", got2.Fraction)
	}
}

func TestFormatter(t *testing.T) {
	manager := NewManager()
	sgd := manager.FindByCode(SGD)
	fmt := sgd.Formatter()
	if fmt == nil {
		t.Fatal("Expected formatter, got nil")
	}
	// We assume formatter works if it's not nil, as formatter package is tested separately.
}

func TestEquals(t *testing.T) {
	manager := NewManager()
	sgd1 := manager.FindByCode(SGD)
	sgd2 := manager.FindByCode(SGD)
	eur := manager.FindByCode(EUR)

	if !sgd1.Equals(sgd2) {
		t.Error("Expected SGD to equal SGD")
	}
	if sgd1.Equals(eur) {
		t.Error("Expected SGD not to equal EUR")
	}
	if sgd1.Equals(nil) {
		t.Error("Expected SGD not to equal nil")
	}
}

func TestDbValue(t *testing.T) {
	manager := NewManager()
	sgd := manager.FindByCode(SGD)
	val, err := sgd.DbValue()
	if err != nil {
		t.Fatalf("DbValue error: %v", err)
	}
	if val != SGD {
		t.Errorf("Expected 'SGD', got %v", val)
	}
}

func TestDbScan(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		wantErr        bool
		wantCode       string
		wantNumeric    string
		wantFraction   int
		wantGrapheme   string
		checkAllFields bool
	}{
		{
			name:           "valid string - SGD",
			input:          SGD,
			wantErr:        false,
			wantCode:       SGD,
			wantNumeric:    "702",
			wantFraction:   2,
			wantGrapheme:   "S$",
			checkAllFields: true,
		},
		{
			name:           "valid string - EUR",
			input:          EUR,
			wantErr:        false,
			wantCode:       EUR,
			wantNumeric:    "978",
			wantFraction:   2,
			wantGrapheme:   "€",
			checkAllFields: true,
		},
		{
			name:           "valid string - GBP",
			input:          GBP,
			wantErr:        false,
			wantCode:       GBP,
			wantNumeric:    "826",
			checkAllFields: true,
		},
		{
			name:           "valid bytes - SGD",
			input:          []byte(SGD),
			wantErr:        false,
			wantCode:       SGD,
			wantNumeric:    "702",
			checkAllFields: true,
		},
		{
			name:           "valid bytes - EUR",
			input:          []byte(EUR),
			wantErr:        false,
			wantCode:       EUR,
			wantNumeric:    "978",
			checkAllFields: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "empty bytes",
			input:   []byte(""),
			wantErr: true,
		},
		{
			name:    "invalid currency code",
			input:   "INVALID_CURRENCY_XYZ",
			wantErr: true,
		},
		{
			name:    "invalid currency bytes",
			input:   []byte("INVALID_CODE"),
			wantErr: true,
		},
		{
			name:    "invalid type - int",
			input:   123,
			wantErr: true,
		},
		{
			name:    "invalid type - float64",
			input:   123.45,
			wantErr: true,
		},
		{
			name:    "invalid type - bool",
			input:   true,
			wantErr: true,
		},
		{
			name:    "invalid type - nil",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "invalid type - struct",
			input:   struct{ Code string }{Code: "SGD"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Currency{}
			err := c.DbScan(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for input %v, but got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if c.Code != tt.wantCode {
				t.Errorf("Code: got %s, want %s", c.Code, tt.wantCode)
			}

			if tt.checkAllFields {
				if c.NumericCode != tt.wantNumeric {
					t.Errorf("NumericCode: got %s, want %s", c.NumericCode, tt.wantNumeric)
				}
				if tt.wantFraction > 0 && c.Fraction != tt.wantFraction {
					t.Errorf("Fraction: got %d, want %d", c.Fraction, tt.wantFraction)
				}
				if tt.wantGrapheme != "" && c.Grapheme != tt.wantGrapheme {
					t.Errorf("Grapheme: got %s, want %s", c.Grapheme, tt.wantGrapheme)
				}
				if c.Decimal == "" {
					t.Error("Decimal should be populated")
				}
				if c.Thousand == "" {
					t.Error("Thousand should be populated")
				}
			}
		})
	}
}

func TestDbScanOverwrite(t *testing.T) {
	// Test that DbScan overwrites existing currency data
	c := &Currency{
		Code:        "OLD",
		NumericCode: "000",
		Fraction:    0,
		Grapheme:    "O",
	}

	err := c.DbScan(SGD)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if c.Code != SGD {
		t.Errorf("Code not overwritten: got %s, want %s", c.Code, SGD)
	}
	if c.NumericCode != "702" {
		t.Errorf("NumericCode not overwritten: got %s, want 702", c.NumericCode)
	}
	if c.Fraction != 2 {
		t.Errorf("Fraction not overwritten: got %d, want 2", c.Fraction)
	}
}

func TestDbScanMultipleCalls(t *testing.T) {
	// Test multiple consecutive scans
	c := &Currency{}

	// First scan
	err := c.DbScan(SGD)
	if err != nil {
		t.Fatalf("First scan error: %v", err)
	}
	if c.Code != SGD {
		t.Errorf("First scan: got %s, want %s", c.Code, SGD)
	}

	// Second scan with different currency
	err = c.DbScan([]byte(EUR))
	if err != nil {
		t.Fatalf("Second scan error: %v", err)
	}
	if c.Code != EUR {
		t.Errorf("Second scan: got %s, want %s", c.Code, EUR)
	}
	if c.NumericCode != "978" {
		t.Errorf("Second scan NumericCode: got %s, want 978", c.NumericCode)
	}

	// Third scan with invalid code should error
	err = c.DbScan("INVALID")
	if err == nil {
		t.Error("Expected error for invalid code")
	}
	// Currency should still have EUR data after failed scan
	if c.Code != EUR {
		t.Errorf("After failed scan: got %s, want %s (should not change)", c.Code, EUR)
	}
}

func TestAddCurrency(t *testing.T) {
	manager := NewManager()
	c := manager.AddFrom("NEW", "N", "1$", ".", ",", "999001", 2)
	if c.Code != "NEW" {
		t.Errorf("Expected NEW, got %s", c.Code)
	}

	got := manager.FindByCode("NEW")
	if got == nil {
		t.Fatal("Expected to retrieve NEW currency")
	}
	if got.Grapheme != "N" {
		t.Errorf("Expected N, got %s", got.Grapheme)
	}
}

func TestCurrencyManagerAndMap_CoveragePaths(t *testing.T) {
	t.Run("Manager nil receiver accessors", func(t *testing.T) {
		var cm *Manager
		if got := cm.FindByCode("SGD"); got != nil {
			t.Fatalf("(*Manager)(nil).FindByCode() = %#v, want nil", got)
		}
		if got := cm.GetDefault(); got != nil {
			t.Fatalf("(*Manager)(nil).GetDefault() = %#v, want nil", got)
		}
		if got := cm.GetSymbols(); got != nil {
			t.Fatalf("(*Manager)(nil).GetSymbols() = %#v, want nil", got)
		}
	})

	t.Run("Manager Add trims empty code", func(t *testing.T) {
		manager := NewManager()
		if got := manager.Add(&Currency{Code: "   "}); got != nil {
			t.Fatalf("Add(empty code) = %#v, want nil", got)
		}
	})

	t.Run("Map HasInvalidState nil dataset", func(t *testing.T) {
		var m Map
		invalid, err := m.HasInvalidState()
		if !invalid || !errors.Is(err, exception.ErrNoCurrencyMapDataset) {
			t.Fatalf("HasInvalidState() = (%v,%v), want (true, ErrNoCurrencyMapDataset)", invalid, err)
		}
	})

	t.Run("Map GetCodes returns all keys", func(t *testing.T) {
		data := map[string]*Currency{
			"SGD": {Code: "SGD"},
			"EUR": {Code: "EUR"},
		}
		m := Map{dataset: &data}

		codes := m.GetCodes()
		if codes == nil || len(*codes) != 2 {
			t.Fatalf("GetCodes() = %#v, want 2 codes", codes)
		}
	})
}

func TestCurrencyDbValueAndScan_ErrorPaths(t *testing.T) {
	t.Run("DbValue nil receiver", func(t *testing.T) {
		var c *Currency
		_, err := c.DbValue()
		if !errors.Is(err, exception.ErrCurrencyNotFound) {
			t.Fatalf("DbValue() error = %v, want ErrCurrencyNotFound", err)
		}
	})

	t.Run("DbScan nil receiver", func(t *testing.T) {
		var c *Currency
		err := c.DbScan("SGD")
		if !errors.Is(err, exception.ErrCurrencyNotFound) {
			t.Fatalf("DbScan() error = %v, want ErrCurrencyNotFound", err)
		}
	})

	t.Run("DbScan unsupported type", func(t *testing.T) {
		var c Currency
		if err := c.DbScan(123); err == nil {
			t.Fatal("DbScan() expected error for unsupported type")
		}
	})

	t.Run("DbScan invalid currency code", func(t *testing.T) {
		var c Currency
		if err := c.DbScan("INVALID"); err == nil {
			t.Fatal("DbScan() expected error for invalid currency code")
		}
	})
}
